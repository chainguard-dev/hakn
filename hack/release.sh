#!/usr/bin/env bash

# Copyright 2022 Chainguard, Inc.
# SPDX-License-Identifier: Apache-2.0

source $(dirname $0)/../vendor/knative.dev/hack/release.sh

declare -A COMPONENTS
COMPONENTS=(
  ["sample.yaml"]="config"
)
readonly COMPONENTS

function build_release() {
   # Update release labels if this is a tagged release
  if [[ -n "${TAG}" ]]; then
    echo "Tagged release, updating release labels to samples.knative.dev/release: \"${TAG}\""
    LABEL_YAML_CMD=(sed -e "s|samples.knative.dev/release: devel|samples.knative.dev/release: \"${TAG}\"|")
  else
    echo "Untagged release, will NOT update release labels"
    LABEL_YAML_CMD=(cat)
  fi

  local all_yamls=()
  for yaml in "${!COMPONENTS[@]}"; do
    local config="${COMPONENTS[${yaml}]}"
    echo "Building Knative Sample Controller - ${config}"
    ko resolve ${KO_FLAGS} -f ${config}/ | "${LABEL_YAML_CMD[@]}" > ${yaml}
    all_yamls+=(${yaml})
  done
  ARTIFACTS_TO_PUBLISH="${all_yamls[@]}"
}

main $@
