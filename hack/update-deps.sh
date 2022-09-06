#!/usr/bin/env bash

# Copyright 2022 Chainguard, Inc.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

source $(dirname "$0")/../vendor/knative.dev/hack/library.sh

export FLOATING_DEPS=(
    # If we want to float other things.
)

go_update_deps "$@"

# Apply any patches
if [ "$(ls -A hack/patches)" ]; then
  git apply hack/patches/*.patch
fi

rm -rf $(find vendor/knative.dev/ -type l)

function rewrite_annotation() {
  sed -E 's@(serving|eventing).knative.dev/release@knative.dev/release@g'
}

function rewrite_webhook() {
  sed 's@webhook.serving.knative.dev@webhook.hakn.chainguard.dev@g' | \
    sed 's@name: eventing-webhook@name: webhook@g' | \
    sed 's@timeoutSeconds: \d+@timeoutSeconds: 25@g'
}

function rewrite_common() {
  local readonly INPUT="${1}"
  local readonly OUTPUT_DIR="${2}"

  cat "${INPUT}" | rewrite_annotation | \
    rewrite_webhook | rewrite_nobody | sed -e's/[[:space:]]*$//' \
    > "${OUTPUT_DIR}/$(basename ${INPUT})"
}

function list_yamls() {
  find "$1" -type f -name '*.yaml' -mindepth 1 -maxdepth 1
}

function rewrite_nobody() {
  sed -e $'s@65534@65532@g'
}

# Remove all of the imported yamls before we start to do our rewrites.
rm $(find config/ -type f | grep imported) || true

#################################################
#
#
#    Serving
#
#
#################################################

# Do a blanket copy of these resources
for x in $(list_yamls ./vendor/knative.dev/serving/config/core/300-resources); do
  rewrite_common "$x" "./config/serving/200-imported/200-serving/100-resources"
done
for x in $(list_yamls ./vendor/knative.dev/serving/config/core/webhooks | grep -v domainmapping); do
  rewrite_common "$x" "./config/serving/200-imported/200-serving/webhooks"
done

rewrite_common "./vendor/knative.dev/serving/config/core/200-roles/podspecable-bindings-clusterrole.yaml" "./config/serving/200-imported/200-serving"


# We need the Image resource from caching, but used by serving.
rewrite_common "./vendor/knative.dev/caching/config/image.yaml" "./config/serving/200-imported/200-serving/100-resources"

# We need the resources from networking, but used by serving.
rewrite_common "./vendor/knative.dev/networking/config/certificate.yaml" "./config/serving/200-imported/200-serving/100-resources"
rewrite_common "./vendor/knative.dev/networking/config/ingress.yaml" "./config/serving/200-imported/200-serving/100-resources"
rewrite_common "./vendor/knative.dev/networking/config/serverlessservice.yaml" "./config/serving/200-imported/200-serving/100-resources"
rewrite_common "./vendor/knative.dev/networking/config/domain-claim.yaml" "./config/serving/200-imported/200-serving/100-resources"

#################################################
#
#
#    net-istio
#
#
#################################################

rewrite_common "./vendor/knative.dev/net-istio/config/203-local-gateway.yaml" "./config/serving/200-imported/200-serving"


#################################################
#
#
#    Set up kodata
#
#
#################################################

# Make sure that all binaries have the appropriate kodata with our version and license data.
for binary in $(find ./config/ -type f | xargs grep ko:// | sed 's@.*ko://@@g' | sed 's@",$@@g' | sort | uniq); do
  if [[ ! -d ./vendor/$binary ]]; then
    echo Skipping $binary, not in vendor.
    continue
  fi
  mkdir ./vendor/$binary/kodata
  pushd ./vendor/$binary/kodata > /dev/null
  ln -s $(echo vendor/$binary/kodata | sed -E 's@[^/]+@..@g')/.git/HEAD .
  ln -s $(echo vendor/$binary/kodata | sed -E 's@[^/]+@..@g')/.git/refs .
  ln -s $(echo vendor/$binary/kodata | sed -E 's@[^/]+@..@g')/LICENSE .
  ln -s $(echo vendor/$binary/kodata | sed -E 's@[^/]+@..@g')/third_party/VENDOR-LICENSE .
  popd > /dev/null
done
