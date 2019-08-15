#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
#
# we skip informers and listers for metrics, because we don't quite support the requisite operations yet
# we skip generating the internal clientset as it's not really needed
GO111MODULE=on "${SCRIPT_ROOT}/generate-internal-groups.sh" deepcopy,conversion \
  github.com/kubernetes-incubator/metrics-server/pkg/client github.com/kubernetes-incubator/metrics-server/pkg/apis github.com/kubernetes-incubator/metrics-server/pkg/apis \
  "tidbcluster:v1alpha1" \
  --go-header-file "${SCRIPT_ROOT}/boilerplate.go.txt"
GO111MODULE=on "${SCRIPT_ROOT}/generate-groups.sh" client \
  github.com/kubernetes-incubator/metrics-server/pkg/client github.com/kubernetes-incubator/metrics-server/pkg/apis \
  "tidbcluster:v1alpha1" \
  --go-header-file "${SCRIPT_ROOT}/boilerplate.go.txt"
