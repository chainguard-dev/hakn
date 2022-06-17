#!/usr/bin/env bash

# Copyright 2022 Chainguard, Inc.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/../vendor/knative.dev/hack/codegen-library.sh
export PATH="$GOBIN:$PATH"

function run_yq() {
  run_go_tool github.com/mikefarah/yq/v4@v4.23.1 yq "$@"
}

echo "=== Update Codegen for ${MODULE_NAME}"

group "Update deps post-codegen"

# Make sure our dependencies are up-to-date
${REPO_ROOT_DIR}/hack/update-deps.sh
