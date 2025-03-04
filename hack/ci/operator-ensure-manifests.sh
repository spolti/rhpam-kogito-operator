#!/bin/bash
# Copyright 2020 Red Hat, Inc. and/or its affiliates
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Script responsible for ensuring and correcting manifests as needed.
set -e

source ./hack/env.sh

tempfolder=$(mktemp -d)
echo "Temporary folder is ${tempfolder}"

version=$(getOperatorVersion)

git clone https://github.com/operator-framework/community-operators.git "${tempfolder}"
mkdir -p "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/"
## copy the latest manifests
cp -r bundle/manifests/ "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/"
cp -r bundle/metadata/ "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/"
cp -r bundle/tests/ "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/"
cp bundle.Dockerfile "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/Dockerfile"

#Edit dockerfile with correct relative path
sed -i "s|bundle/manifests|manifests|g" "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/Dockerfile"
sed -i "s|bundle/metadata|metadata|g" "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/Dockerfile"
sed -i "s|bundle/tests|tests|g" "${tempfolder}/community-operators/rhpam-kogito-operator/${version}/Dockerfile"
