/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"flag"

	"istio.io/api/networking/v1beta1"

	"knative.dev/net-istio/pkg/reconciler/informerfiltering"
	kingress "knative.dev/net-istio/pkg/reconciler/ingress"
	istioss "knative.dev/net-istio/pkg/reconciler/serverlessservice"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
	"knative.dev/pkg/webhook"
	"knative.dev/pkg/webhook/certificates"
	"knative.dev/serving/pkg/reconciler/configuration"
	"knative.dev/serving/pkg/reconciler/domainmapping"
	"knative.dev/serving/pkg/reconciler/gc"
	"knative.dev/serving/pkg/reconciler/labeler"
	"knative.dev/serving/pkg/reconciler/revision"
	"knative.dev/serving/pkg/reconciler/route"
	"knative.dev/serving/pkg/reconciler/serverlessservice"
	"knative.dev/serving/pkg/reconciler/service"

	tlscertificates "github.com/chainguard-dev/hakn/pkg/reconciler/certificates"
)

func main() {
	flag.Parse()
	// Allow unknown fields in Istio API client. This is to be more
	// resilient to clusters containing malformed resources.
	v1beta1.VirtualServiceUnmarshaler.AllowUnknownFields = true
	v1beta1.GatewayUnmarshaler.AllowUnknownFields = true

	ctx := informerfiltering.GetContextWithFilteringLabelSelector(signals.NewContext())

	ctx = webhook.WithOptions(ctx, webhook.Options{
		ServiceName: "webhook",
		Port:        8443,
		SecretName:  "webhook-certs",
	})

	sharedmain.MainWithConfig(ctx, "hakn-serving", injection.ParseAndGetRESTConfigOrDie(),
		certificates.NewController,
		newDefaultingAdmissionController,
		newValidationAdmissionController,
		newConfigValidationController,

		// Serving resource controllers.
		configuration.NewController,
		labeler.NewController,
		revision.NewController,
		route.NewController,
		serverlessservice.NewController,
		service.NewController,
		gc.NewController,
		domainmapping.NewController,

		// KIngress controller.
		kingress.NewController,
		istioss.NewController,

		// Custom controller for building self-signed certs for terminating
		// GCLB's TLS
		tlscertificates.NewController,
	)
}
