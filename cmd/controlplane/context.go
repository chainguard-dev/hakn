/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"sync"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/logging"

	knsdefaultconfig "knative.dev/serving/pkg/apis/config"
)

var (
	decorator func(context.Context) context.Context
	decoSetup sync.Once
)

func newContextDecorator(ctx context.Context, cmw configmap.Watcher) func(ctx context.Context) context.Context {
	// We don't need a copy of this for every webhook, so just set it up
	// once, and return that singleton to each of our webhooks.
	decoSetup.Do(func() {
		// Decorate contexts with the current state of the config.
		knsstore := knsdefaultconfig.NewStore(logging.FromContext(ctx).Named("kns-config-store"))
		knsstore.WatchConfigs(cmw)

		decorator = func(ctx context.Context) context.Context {
			ctx = knsstore.ToContext(ctx)
			return ctx
		}
	})

	return decorator
}
