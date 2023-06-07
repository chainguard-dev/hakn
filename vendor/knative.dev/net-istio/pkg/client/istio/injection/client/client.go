/*
Copyright 2020 The Knative Authors

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

// Code generated by injection-gen. DO NOT EDIT.

package client

import (
	context "context"
	json "encoding/json"
	errors "errors"
	fmt "fmt"

	v1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	discovery "k8s.io/client-go/discovery"
	dynamic "k8s.io/client-go/dynamic"
	rest "k8s.io/client-go/rest"
	versioned "knative.dev/net-istio/pkg/client/istio/clientset/versioned"
	typednetworkingv1beta1 "knative.dev/net-istio/pkg/client/istio/clientset/versioned/typed/networking/v1beta1"
	injection "knative.dev/pkg/injection"
	dynamicclient "knative.dev/pkg/injection/clients/dynamicclient"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterClient(withClientFromConfig)
	injection.Default.RegisterClientFetcher(func(ctx context.Context) interface{} {
		return Get(ctx)
	})
	injection.Dynamic.RegisterDynamicClient(withClientFromDynamic)
}

// Key is used as the key for associating information with a context.Context.
type Key struct{}

func withClientFromConfig(ctx context.Context, cfg *rest.Config) context.Context {
	return context.WithValue(ctx, Key{}, versioned.NewForConfigOrDie(cfg))
}

func withClientFromDynamic(ctx context.Context) context.Context {
	return context.WithValue(ctx, Key{}, &wrapClient{dyn: dynamicclient.Get(ctx)})
}

// Get extracts the versioned.Interface client from the context.
func Get(ctx context.Context) versioned.Interface {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		if injection.GetConfig(ctx) == nil {
			logging.FromContext(ctx).Panic(
				"Unable to fetch knative.dev/net-istio/pkg/client/istio/clientset/versioned.Interface from context. This context is not the application context (which is typically given to constructors via sharedmain).")
		} else {
			logging.FromContext(ctx).Panic(
				"Unable to fetch knative.dev/net-istio/pkg/client/istio/clientset/versioned.Interface from context.")
		}
	}
	return untyped.(versioned.Interface)
}

type wrapClient struct {
	dyn dynamic.Interface
}

var _ versioned.Interface = (*wrapClient)(nil)

func (w *wrapClient) Discovery() discovery.DiscoveryInterface {
	panic("Discovery called on dynamic client!")
}

func convert(from interface{}, to runtime.Object) error {
	bs, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("Marshal() = %w", err)
	}
	if err := json.Unmarshal(bs, to); err != nil {
		return fmt.Errorf("Unmarshal() = %w", err)
	}
	return nil
}

// NetworkingV1beta1 retrieves the NetworkingV1beta1Client
func (w *wrapClient) NetworkingV1beta1() typednetworkingv1beta1.NetworkingV1beta1Interface {
	return &wrapNetworkingV1beta1{
		dyn: w.dyn,
	}
}

type wrapNetworkingV1beta1 struct {
	dyn dynamic.Interface
}

func (w *wrapNetworkingV1beta1) RESTClient() rest.Interface {
	panic("RESTClient called on dynamic client!")
}

func (w *wrapNetworkingV1beta1) DestinationRules(namespace string) typednetworkingv1beta1.DestinationRuleInterface {
	return &wrapNetworkingV1beta1DestinationRuleImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1beta1",
			Resource: "destinationrules",
		}),

		namespace: namespace,
	}
}

type wrapNetworkingV1beta1DestinationRuleImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typednetworkingv1beta1.DestinationRuleInterface = (*wrapNetworkingV1beta1DestinationRuleImpl)(nil)

func (w *wrapNetworkingV1beta1DestinationRuleImpl) Create(ctx context.Context, in *v1beta1.DestinationRule, opts v1.CreateOptions) (*v1beta1.DestinationRule, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "DestinationRule",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.DestinationRule{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1DestinationRuleImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapNetworkingV1beta1DestinationRuleImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapNetworkingV1beta1DestinationRuleImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.DestinationRule, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.DestinationRule{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1DestinationRuleImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.DestinationRuleList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.DestinationRuleList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1DestinationRuleImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.DestinationRule, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.DestinationRule{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1DestinationRuleImpl) Update(ctx context.Context, in *v1beta1.DestinationRule, opts v1.UpdateOptions) (*v1beta1.DestinationRule, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "DestinationRule",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.DestinationRule{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1DestinationRuleImpl) UpdateStatus(ctx context.Context, in *v1beta1.DestinationRule, opts v1.UpdateOptions) (*v1beta1.DestinationRule, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "DestinationRule",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.DestinationRule{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1DestinationRuleImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}

func (w *wrapNetworkingV1beta1) Gateways(namespace string) typednetworkingv1beta1.GatewayInterface {
	return &wrapNetworkingV1beta1GatewayImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1beta1",
			Resource: "gateways",
		}),

		namespace: namespace,
	}
}

type wrapNetworkingV1beta1GatewayImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typednetworkingv1beta1.GatewayInterface = (*wrapNetworkingV1beta1GatewayImpl)(nil)

func (w *wrapNetworkingV1beta1GatewayImpl) Create(ctx context.Context, in *v1beta1.Gateway, opts v1.CreateOptions) (*v1beta1.Gateway, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "Gateway",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Gateway{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1GatewayImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapNetworkingV1beta1GatewayImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapNetworkingV1beta1GatewayImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.Gateway, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Gateway{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1GatewayImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.GatewayList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.GatewayList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1GatewayImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Gateway, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Gateway{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1GatewayImpl) Update(ctx context.Context, in *v1beta1.Gateway, opts v1.UpdateOptions) (*v1beta1.Gateway, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "Gateway",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Gateway{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1GatewayImpl) UpdateStatus(ctx context.Context, in *v1beta1.Gateway, opts v1.UpdateOptions) (*v1beta1.Gateway, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "Gateway",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Gateway{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1GatewayImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}

func (w *wrapNetworkingV1beta1) ProxyConfigs(namespace string) typednetworkingv1beta1.ProxyConfigInterface {
	return &wrapNetworkingV1beta1ProxyConfigImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1beta1",
			Resource: "proxyconfigs",
		}),

		namespace: namespace,
	}
}

type wrapNetworkingV1beta1ProxyConfigImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typednetworkingv1beta1.ProxyConfigInterface = (*wrapNetworkingV1beta1ProxyConfigImpl)(nil)

func (w *wrapNetworkingV1beta1ProxyConfigImpl) Create(ctx context.Context, in *v1beta1.ProxyConfig, opts v1.CreateOptions) (*v1beta1.ProxyConfig, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "ProxyConfig",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ProxyConfig{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ProxyConfigImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapNetworkingV1beta1ProxyConfigImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapNetworkingV1beta1ProxyConfigImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.ProxyConfig, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ProxyConfig{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ProxyConfigImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.ProxyConfigList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ProxyConfigList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ProxyConfigImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ProxyConfig, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ProxyConfig{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ProxyConfigImpl) Update(ctx context.Context, in *v1beta1.ProxyConfig, opts v1.UpdateOptions) (*v1beta1.ProxyConfig, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "ProxyConfig",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ProxyConfig{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ProxyConfigImpl) UpdateStatus(ctx context.Context, in *v1beta1.ProxyConfig, opts v1.UpdateOptions) (*v1beta1.ProxyConfig, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "ProxyConfig",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ProxyConfig{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ProxyConfigImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}

func (w *wrapNetworkingV1beta1) ServiceEntries(namespace string) typednetworkingv1beta1.ServiceEntryInterface {
	return &wrapNetworkingV1beta1ServiceEntryImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1beta1",
			Resource: "serviceentries",
		}),

		namespace: namespace,
	}
}

type wrapNetworkingV1beta1ServiceEntryImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typednetworkingv1beta1.ServiceEntryInterface = (*wrapNetworkingV1beta1ServiceEntryImpl)(nil)

func (w *wrapNetworkingV1beta1ServiceEntryImpl) Create(ctx context.Context, in *v1beta1.ServiceEntry, opts v1.CreateOptions) (*v1beta1.ServiceEntry, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "ServiceEntry",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ServiceEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ServiceEntryImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapNetworkingV1beta1ServiceEntryImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapNetworkingV1beta1ServiceEntryImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.ServiceEntry, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ServiceEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ServiceEntryImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.ServiceEntryList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ServiceEntryList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ServiceEntryImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ServiceEntry, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ServiceEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ServiceEntryImpl) Update(ctx context.Context, in *v1beta1.ServiceEntry, opts v1.UpdateOptions) (*v1beta1.ServiceEntry, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "ServiceEntry",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ServiceEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ServiceEntryImpl) UpdateStatus(ctx context.Context, in *v1beta1.ServiceEntry, opts v1.UpdateOptions) (*v1beta1.ServiceEntry, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "ServiceEntry",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.ServiceEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1ServiceEntryImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}

func (w *wrapNetworkingV1beta1) Sidecars(namespace string) typednetworkingv1beta1.SidecarInterface {
	return &wrapNetworkingV1beta1SidecarImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1beta1",
			Resource: "sidecars",
		}),

		namespace: namespace,
	}
}

type wrapNetworkingV1beta1SidecarImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typednetworkingv1beta1.SidecarInterface = (*wrapNetworkingV1beta1SidecarImpl)(nil)

func (w *wrapNetworkingV1beta1SidecarImpl) Create(ctx context.Context, in *v1beta1.Sidecar, opts v1.CreateOptions) (*v1beta1.Sidecar, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "Sidecar",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Sidecar{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1SidecarImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapNetworkingV1beta1SidecarImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapNetworkingV1beta1SidecarImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.Sidecar, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Sidecar{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1SidecarImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.SidecarList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.SidecarList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1SidecarImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Sidecar, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Sidecar{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1SidecarImpl) Update(ctx context.Context, in *v1beta1.Sidecar, opts v1.UpdateOptions) (*v1beta1.Sidecar, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "Sidecar",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Sidecar{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1SidecarImpl) UpdateStatus(ctx context.Context, in *v1beta1.Sidecar, opts v1.UpdateOptions) (*v1beta1.Sidecar, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "Sidecar",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.Sidecar{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1SidecarImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}

func (w *wrapNetworkingV1beta1) VirtualServices(namespace string) typednetworkingv1beta1.VirtualServiceInterface {
	return &wrapNetworkingV1beta1VirtualServiceImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1beta1",
			Resource: "virtualservices",
		}),

		namespace: namespace,
	}
}

type wrapNetworkingV1beta1VirtualServiceImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typednetworkingv1beta1.VirtualServiceInterface = (*wrapNetworkingV1beta1VirtualServiceImpl)(nil)

func (w *wrapNetworkingV1beta1VirtualServiceImpl) Create(ctx context.Context, in *v1beta1.VirtualService, opts v1.CreateOptions) (*v1beta1.VirtualService, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "VirtualService",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.VirtualService{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1VirtualServiceImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapNetworkingV1beta1VirtualServiceImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapNetworkingV1beta1VirtualServiceImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.VirtualService, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.VirtualService{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1VirtualServiceImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.VirtualServiceList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.VirtualServiceList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1VirtualServiceImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.VirtualService, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.VirtualService{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1VirtualServiceImpl) Update(ctx context.Context, in *v1beta1.VirtualService, opts v1.UpdateOptions) (*v1beta1.VirtualService, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "VirtualService",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.VirtualService{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1VirtualServiceImpl) UpdateStatus(ctx context.Context, in *v1beta1.VirtualService, opts v1.UpdateOptions) (*v1beta1.VirtualService, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "VirtualService",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.VirtualService{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1VirtualServiceImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}

func (w *wrapNetworkingV1beta1) WorkloadEntries(namespace string) typednetworkingv1beta1.WorkloadEntryInterface {
	return &wrapNetworkingV1beta1WorkloadEntryImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1beta1",
			Resource: "workloadentries",
		}),

		namespace: namespace,
	}
}

type wrapNetworkingV1beta1WorkloadEntryImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typednetworkingv1beta1.WorkloadEntryInterface = (*wrapNetworkingV1beta1WorkloadEntryImpl)(nil)

func (w *wrapNetworkingV1beta1WorkloadEntryImpl) Create(ctx context.Context, in *v1beta1.WorkloadEntry, opts v1.CreateOptions) (*v1beta1.WorkloadEntry, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "WorkloadEntry",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadEntryImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapNetworkingV1beta1WorkloadEntryImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapNetworkingV1beta1WorkloadEntryImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.WorkloadEntry, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadEntryImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.WorkloadEntryList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadEntryList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadEntryImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.WorkloadEntry, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadEntryImpl) Update(ctx context.Context, in *v1beta1.WorkloadEntry, opts v1.UpdateOptions) (*v1beta1.WorkloadEntry, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "WorkloadEntry",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadEntryImpl) UpdateStatus(ctx context.Context, in *v1beta1.WorkloadEntry, opts v1.UpdateOptions) (*v1beta1.WorkloadEntry, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "WorkloadEntry",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadEntry{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadEntryImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}

func (w *wrapNetworkingV1beta1) WorkloadGroups(namespace string) typednetworkingv1beta1.WorkloadGroupInterface {
	return &wrapNetworkingV1beta1WorkloadGroupImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "networking.istio.io",
			Version:  "v1beta1",
			Resource: "workloadgroups",
		}),

		namespace: namespace,
	}
}

type wrapNetworkingV1beta1WorkloadGroupImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typednetworkingv1beta1.WorkloadGroupInterface = (*wrapNetworkingV1beta1WorkloadGroupImpl)(nil)

func (w *wrapNetworkingV1beta1WorkloadGroupImpl) Create(ctx context.Context, in *v1beta1.WorkloadGroup, opts v1.CreateOptions) (*v1beta1.WorkloadGroup, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "WorkloadGroup",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadGroup{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadGroupImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapNetworkingV1beta1WorkloadGroupImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapNetworkingV1beta1WorkloadGroupImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.WorkloadGroup, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadGroup{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadGroupImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.WorkloadGroupList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadGroupList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadGroupImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.WorkloadGroup, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadGroup{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadGroupImpl) Update(ctx context.Context, in *v1beta1.WorkloadGroup, opts v1.UpdateOptions) (*v1beta1.WorkloadGroup, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "WorkloadGroup",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadGroup{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadGroupImpl) UpdateStatus(ctx context.Context, in *v1beta1.WorkloadGroup, opts v1.UpdateOptions) (*v1beta1.WorkloadGroup, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "WorkloadGroup",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.WorkloadGroup{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapNetworkingV1beta1WorkloadGroupImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}
