/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/webhook/resourcesemantics/defaulting"
)

func newDefaultingAdmissionController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return defaulting.NewAdmissionController(ctx,

		// Name of the resource webhook.
		"webhook.hakn.chainguard.dev",

		// The path on which to serve the webhook.
		"/defaulting",

		// The resources to validate and default.
		ourTypes,

		// A function that infuses the context passed to Validate/SetDefaults with custom metadata.
		newContextDecorator(ctx, cmw),

		// Whether to disallow unknown fields.
		true,
	)
}
