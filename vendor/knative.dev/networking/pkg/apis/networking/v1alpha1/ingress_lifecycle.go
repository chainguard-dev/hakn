/*
Copyright 2019 The Knative Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
)

var ingressCondSet = apis.NewLivingConditionSet(
	IngressConditionNetworkConfigured,
	IngressConditionLoadBalancerReady,
)

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (*Ingress) GetConditionSet() apis.ConditionSet {
	return ingressCondSet
}

// GetGroupVersionKind returns SchemeGroupVersion of an Ingress
func (i *Ingress) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("Ingress")
}

// GetCondition returns the current condition of a given condition type
func (is *IngressStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return ingressCondSet.Manage(is).GetCondition(t)
}

// InitializeConditions initializes conditions of an IngressStatus
func (is *IngressStatus) InitializeConditions() {
	ingressCondSet.Manage(is).InitializeConditions()
}

// MarkNetworkConfigured set IngressConditionNetworkConfigured in IngressStatus as true
func (is *IngressStatus) MarkNetworkConfigured() {
	ingressCondSet.Manage(is).MarkTrue(IngressConditionNetworkConfigured)
}

// MarkResourceNotOwned changes the "NetworkConfigured" condition to false to reflect that the
// resource of the given kind and name has already been created, and we do not own it.
func (is *IngressStatus) MarkResourceNotOwned(kind, name string) {
	ingressCondSet.Manage(is).MarkFalse(IngressConditionNetworkConfigured, "NotOwned",
		fmt.Sprintf("There is an existing %s %q that we do not own.", kind, name))
}

// MarkLoadBalancerReady marks the Ingress with IngressConditionLoadBalancerReady,
// and also populate the address of the load balancer.
func (is *IngressStatus) MarkLoadBalancerReady(publicLbs []LoadBalancerIngressStatus, privateLbs []LoadBalancerIngressStatus) {
	is.PublicLoadBalancer = &LoadBalancerStatus{Ingress: publicLbs}
	is.PrivateLoadBalancer = &LoadBalancerStatus{Ingress: privateLbs}

	ingressCondSet.Manage(is).MarkTrue(IngressConditionLoadBalancerReady)
}

// MarkLoadBalancerNotReady marks the "IngressConditionLoadBalancerReady" condition to unknown to
// reflect that the load balancer is not ready yet.
func (is *IngressStatus) MarkLoadBalancerNotReady() {
	ingressCondSet.Manage(is).MarkUnknown(IngressConditionLoadBalancerReady, "Uninitialized",
		"Waiting for load balancer to be ready")
}

// MarkLoadBalancerFailed marks the "IngressConditionLoadBalancerReady" condition to false.
func (is *IngressStatus) MarkLoadBalancerFailed(reason, message string) {
	ingressCondSet.Manage(is).MarkFalse(IngressConditionLoadBalancerReady, reason, message)
}

// MarkIngressNotReady marks the "IngressConditionReady" condition to unknown.
func (is *IngressStatus) MarkIngressNotReady(reason, message string) {
	ingressCondSet.Manage(is).MarkUnknown(IngressConditionReady, reason, message)
}

// IsReady returns true if the Status condition MetricConditionReady
// is true and the latest spec has been observed.
func (i *Ingress) IsReady() bool {
	is := i.Status
	return is.ObservedGeneration == i.Generation &&
		is.GetCondition(IngressConditionReady).IsTrue()
}
