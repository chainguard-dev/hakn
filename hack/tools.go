//go:build tools
// +build tools

/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package tools

// This package imports things required by this repository, to force `go mod` to see them as dependencies
import (
	_ "knative.dev/hack"

	// codegen: hack/generate-knative.sh
	_ "knative.dev/pkg/hack"

	// networking resources
	_ "knative.dev/networking/config"

	// net-istio config
	_ "knative.dev/net-istio/config"

	// All of the binary entrypoints from our config
	_ "knative.dev/serving/cmd/activator"
	_ "knative.dev/serving/cmd/autoscaler"
	_ "knative.dev/serving/cmd/queue"

	// config directories
	_ "knative.dev/caching/config"
	_ "knative.dev/serving/config/core/300-resources"
	_ "knative.dev/serving/config/core/deployments"
	_ "knative.dev/serving/config/core/webhooks"
)
