/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/metrics"
	"knative.dev/pkg/webhook"
	"knative.dev/pkg/webhook/configmaps"
	"knative.dev/pkg/webhook/resourcesemantics/validation"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	extravalidation "knative.dev/serving/pkg/webhook"

	// config validation constructors
	istioconfig "knative.dev/net-istio/pkg/reconciler/ingress/config"
	network "knative.dev/networking/pkg"
	"knative.dev/networking/pkg/config"
	pkgleaderelection "knative.dev/pkg/leaderelection"
	tracingconfig "knative.dev/pkg/tracing/config"
	knsdefaultconfig "knative.dev/serving/pkg/apis/config"
	autoscalerconfig "knative.dev/serving/pkg/autoscaler/config"
	"knative.dev/serving/pkg/deployment"
	gcconfig "knative.dev/serving/pkg/gc"
	domainconfig "knative.dev/serving/pkg/reconciler/route/config"
)

var serviceValidation = validation.NewCallback(
	extravalidation.ValidateService, webhook.Create, webhook.Update)

var configValidation = validation.NewCallback(
	extravalidation.ValidateConfiguration, webhook.Create, webhook.Update)

var callbacks = map[schema.GroupVersionKind]validation.Callback{
	servingv1.SchemeGroupVersion.WithKind("Service"):       serviceValidation,
	servingv1.SchemeGroupVersion.WithKind("Configuration"): configValidation,
}

func newValidationAdmissionController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return validation.NewAdmissionController(ctx,

		// Name of the resource webhook.
		"validation.webhook.hakn.chainguard.dev",

		// The path on which to serve the webhook.
		"/resource-validation",

		// The resources to validate and default.
		ourTypes,

		// A function that infuses the context passed to Validate/SetDefaults with custom metadata.
		newContextDecorator(ctx, cmw),

		// Whether to disallow unknown fields.
		true,

		// Extra validating callbacks to be applied to resources.
		callbacks,
	)
}

func newConfigValidationController(ctx context.Context, _ configmap.Watcher) *controller.Impl {
	return configmaps.NewAdmissionController(ctx,

		// Name of the configmap webhook.
		"config.webhook.hakn.chainguard.dev",

		// The path on which to serve the webhook.
		"/config-validation",

		// The configmaps to validate.
		configmap.Constructors{
			tracingconfig.ConfigName:            tracingconfig.NewTracingConfigFromConfigMap,
			autoscalerconfig.ConfigName:         autoscalerconfig.NewConfigFromConfigMap,
			gcconfig.ConfigName:                 gcconfig.NewConfigFromConfigMapFunc(ctx),
			config.ConfigMapName:                network.NewConfigFromConfigMap,
			deployment.ConfigName:               deployment.NewConfigFromConfigMap,
			metrics.ConfigMapName():             metrics.NewObservabilityConfigFromConfigMap,
			logging.ConfigMapName():             logging.NewConfigFromConfigMap,
			domainconfig.DomainConfigName:       domainconfig.NewDomainFromConfigMap,
			pkgleaderelection.ConfigMapName():   pkgleaderelection.NewConfigFromConfigMap,
			knsdefaultconfig.DefaultsConfigName: knsdefaultconfig.NewDefaultsConfigFromConfigMap,
			istioconfig.IstioConfigName:         istioconfig.NewIstioFromConfigMap,
		},
	)
}
