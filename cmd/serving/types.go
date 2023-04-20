/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/webhook/resourcesemantics"

	net "knative.dev/networking/pkg/apis/networking/v1alpha1"
	autoscalingv1alpha1 "knative.dev/serving/pkg/apis/autoscaling/v1alpha1"
	v1 "knative.dev/serving/pkg/apis/serving/v1"
	"knative.dev/serving/pkg/apis/serving/v1beta1"
)

var ourTypes = map[schema.GroupVersionKind]resourcesemantics.GenericCRD{
	v1beta1.SchemeGroupVersion.WithKind("DomainMapping"): &v1beta1.DomainMapping{},
	v1.SchemeGroupVersion.WithKind("Revision"):           &v1.Revision{},
	v1.SchemeGroupVersion.WithKind("Configuration"):      &v1.Configuration{},
	v1.SchemeGroupVersion.WithKind("Route"):              &v1.Route{},
	v1.SchemeGroupVersion.WithKind("Service"):            &v1.Service{},

	autoscalingv1alpha1.SchemeGroupVersion.WithKind("PodAutoscaler"): &autoscalingv1alpha1.PodAutoscaler{},
	autoscalingv1alpha1.SchemeGroupVersion.WithKind("Metric"):        &autoscalingv1alpha1.Metric{},

	net.SchemeGroupVersion.WithKind("Certificate"):       &net.Certificate{},
	net.SchemeGroupVersion.WithKind("Ingress"):           &net.Ingress{},
	net.SchemeGroupVersion.WithKind("ServerlessService"): &net.ServerlessService{},
}
