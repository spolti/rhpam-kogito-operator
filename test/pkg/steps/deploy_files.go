// Copyright 2020 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package steps

import (
	"fmt"

	"github.com/kiegroup/kogito-operator/api"

	"github.com/cucumber/godog"
	communityFramework "github.com/kiegroup/kogito-operator/test/pkg/framework"
	v1 "github.com/kiegroup/rhpam-kogito-operator/api/v1"
	"github.com/kiegroup/rhpam-kogito-operator/test/pkg/framework"
)

const sourceLocation = "src/main/resources"

func registerKogitoDeployFilesSteps(ctx *godog.ScenarioContext, data *Data) {
	// Deploy steps
	ctx.Step(`^Deploy (quarkus|springboot) file "([^"]*)" from example service "([^"]*)"$`, data.deployFileFromExampleService)
	ctx.Step(`^Deploy (quarkus|springboot) folder from example service "([^"]*)"$`, data.deployFolderFromExampleService)
}

// Deploy steps

func (data *Data) deployFileFromExampleService(runtimeType, file, serviceName string) error {
	sourceFilePath := fmt.Sprintf(`%s/%s/%s/%s`, data.KogitoExamplesLocation, serviceName, sourceLocation, file)
	return deploySourceFilesFromPath(data.Namespace, runtimeType, serviceName, sourceFilePath)
}

func (data *Data) deployFolderFromExampleService(runtimeType, serviceName string) error {
	sourceFolderPath := fmt.Sprintf(`%s/%s/%s`, data.KogitoExamplesLocation, serviceName, sourceLocation)
	return deploySourceFilesFromPath(data.Namespace, runtimeType, serviceName, sourceFolderPath)
}

func deploySourceFilesFromPath(namespace, runtimeType, serviceName, path string) error {
	communityFramework.GetLogger(namespace).Info("Deploying example with source files", "runtimeType", runtimeType, "serviceName", serviceName, "path", path)

	buildHolder, err := getKogitoBuildConfiguredStub(namespace, runtimeType, serviceName, nil)
	if err != nil {
		return err
	}

	buildHolder.KogitoBuild.GetSpec().SetType(api.LocalSourceBuildType)
	buildHolder.KogitoBuild.GetSpec().GetGitSource().SetURI(path)

	err = framework.DeployKogitoBuild(namespace, buildHolder)
	if err != nil {
		return err
	}

	// In case of OpenShift the ImageStream needs to be patched to allow insecure registries
	if communityFramework.IsOpenshift() {
		if err := makeImageStreamInsecure(namespace, framework.GetKogitoBuildS2IImage()); err != nil {
			return err
		}
		if err := makeImageStreamInsecure(namespace, framework.GetKogitoBuildRuntimeImage(buildHolder.KogitoBuild.(*v1.KogitoBuild))); err != nil {
			return err
		}
	}

	// If we don't use Kogito CLI then upload target folder using OC client
	return communityFramework.WaitForOnOpenshift(namespace, fmt.Sprintf("Build '%s-builder' to start", serviceName), defaultTimeoutToStartBuildInMin,
		func() (bool, error) {
			_, err := communityFramework.CreateCommand("oc", "start-build", serviceName+"-builder", "--from-file="+path, "-n", namespace).WithLoggerContext(namespace).Execute()
			return err == nil, err
		})
}
