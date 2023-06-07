/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package certificates

import (
	"context"

	// Injection stuff
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	secretinformer "knative.dev/pkg/client/injection/kube/informers/core/v1/secret"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
)

// NewController constructs a controller for materializing gateway certificates.
// In order for it to bootstrap, an empty secret should be created with the
// expected name (and lifecycle managed accordingly), and thereafter this controller
// will ensure it has the appropriate shape for the webhook.
func NewController(
	ctx context.Context,
	_ configmap.Watcher,
) *controller.Impl {
	client := kubeclient.Get(ctx)

	// TODO(mattmoor): namespace scope this informer
	secretInformer := secretinformer.Get(ctx)

	key := types.NamespacedName{
		Namespace: "istio-system",
		Name:      "tls-cert",
	}

	wh := &reconciler{
		LeaderAwareFuncs: pkgreconciler.LeaderAwareFuncs{
			// Enqueue the key whenever we become leader.
			PromoteFunc: func(bkt pkgreconciler.Bucket, enq func(pkgreconciler.Bucket, types.NamespacedName)) error {
				enq(bkt, key)
				return nil
			},
		},
		key: key,

		client:       client,
		secretlister: secretInformer.Lister(),
	}

	const queueName = "GatewayCertificates"
	c := controller.NewContext(ctx, wh, controller.ControllerOptions{
		WorkQueueName: queueName,
		Logger:        logging.FromContext(ctx).Named(queueName),
	})

	// Reconcile when the cert bundle changes.
	secretInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterWithNameAndNamespace(key.Namespace, key.Name),
		// It doesn't matter what we enqueue because we will always Reconcile
		// the named MWH resource.
		Handler: controller.HandleAll(c.Enqueue),
	})

	return c
}
