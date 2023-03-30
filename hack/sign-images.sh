#!/usr/bin/env bash

# Copyright 2023 Chainguard, Inc.
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

if [[ ! -f serving.images ]]; then
    echo "serving.images not found"
    exit 1
fi

echo "Signing cosign images using Keyless..."

readarray -t serving < <(cat serving.images || true)
cosign sign --yes --timeout 5m "${serving[@]}"
