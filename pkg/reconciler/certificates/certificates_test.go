/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package certificates

import (
	"context"
	"errors"
	"testing"
	"time"

	kubeclient "knative.dev/pkg/client/injection/kube/client/fake"
	_ "knative.dev/pkg/client/injection/kube/informers/core/v1/secret/fake"
	pkgreconciler "knative.dev/pkg/reconciler"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	clientgotesting "k8s.io/client-go/testing"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/system"
	"knative.dev/pkg/webhook"

	certresources "github.com/chainguard-dev/hakn/pkg/reconciler/certificates/resources"

	. "knative.dev/pkg/reconciler/testing"
	. "knative.dev/pkg/webhook/testing"
)

func TestReconcile(t *testing.T) {
	const secretName = "webhook-secret"
	secret, err := certresources.MakeSecret(context.Background(),
		secretName, system.Namespace())
	if err != nil {
		t.Fatal("MakeSecret() =", err)
	}

	// Mutate the MakeSecret to return our secret deterministically.
	certresources.MakeSecret = func(ctx context.Context, name, namespace string) (*corev1.Secret, error) {
		return secret, nil
	}
	defer func() {
		certresources.MakeSecret = certresources.MakeSecretInternal
	}()

	// The key to use, which for this singleton reconciler doesn't matter (although the
	// namespace matters for namespace validation).
	key := system.Namespace() + "/does not matter"

	table := TableTest{{
		Name:    "well formed secret exists",
		Key:     key,
		Objects: []runtime.Object{secret},
	}, {
		Name: "secret does not exist",
		Key:  key,
	}, {
		Name: "missing server key",
		Key:  key,
		Objects: []runtime.Object{&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: system.Namespace(),
			},
			Type: corev1.SecretTypeTLS,
			Data: map[string][]byte{
				corev1.TLSPrivateKeyKey: []byte("present"),
				corev1.TLSCertKey:       []byte("present"),
			},
		}},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: secret,
		}},
	}, {
		Name: "missing server cert",
		Key:  key,
		Objects: []runtime.Object{&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: system.Namespace(),
			},
			Type: corev1.SecretTypeTLS,
			Data: map[string][]byte{
				corev1.TLSPrivateKeyKey: []byte("present"),
				// corev1.TLSCertKey: []byte("missing"),
			},
		}},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: secret,
		}},
	}, {
		Name: "missing CA cert",
		Key:  key,
		Objects: []runtime.Object{&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: system.Namespace(),
			},
			Type: corev1.SecretTypeTLS,
			Data: map[string][]byte{
				corev1.TLSPrivateKeyKey: []byte("present"),
				corev1.TLSCertKey:       []byte("present"),
			},
		}},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: secret,
		}},
	}, {
		Name: "certificate expiring soon",
		Key:  key,
		// 23 hours  falls inside of the grace period of 1 day so the secret will be updated.
		Objects: []runtime.Object{secretWithCertData(t, time.Now().Add(23*time.Hour))},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: secret,
		}},
	}, {
		Name: "certificate not expiring soon",
		Key:  key,
		// 25 hours falls outside of the grace period of 1 day so the secret will not be updated.
		Objects: []runtime.Object{secretWithCertData(t, time.Now().Add(25*time.Hour))},
	}}

	table.Test(t, MakeFactory(func(ctx context.Context, listers *Listers, cmw configmap.Watcher) controller.Reconciler {
		return &reconciler{
			client:       kubeclient.Get(ctx),
			secretlister: listers.GetSecretLister(),
			key: types.NamespacedName{
				Namespace: system.Namespace(),
				Name:      secretName,
			},
		}
	}))
}

func TestReconcileMakeSecretFailure(t *testing.T) {
	secretName := "webhook-secret"
	secret, err := certresources.MakeSecret(context.Background(),
		secretName, system.Namespace())
	if err != nil {
		t.Fatal("MakeSecret() =", err)
	}

	// Mutate the MakeSecret to return our secret deterministically.
	certresources.MakeSecret = func(ctx context.Context, name, namespace string) (*corev1.Secret, error) {
		return nil, errors.New("this is an error")
	}
	defer func() {
		certresources.MakeSecret = certresources.MakeSecretInternal
	}()

	// The key to use, which for this singleton reconciler doesn't matter (although the
	// namespace matters for namespace validation).
	key := system.Namespace() + "/does not matter"

	table := TableTest{{
		Name:    "would return error, but not called",
		Key:     key,
		Objects: []runtime.Object{secret},
	}, {
		Name:    "malformed secret",
		Key:     key,
		WantErr: true,
		Objects: []runtime.Object{&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: system.Namespace(),
			},
			Data: map[string][]byte{
				// corev1.TLSPrivateKeyKey:  []byte("missing"),
				corev1.TLSCertKey: []byte("present"),
			},
		}},
	}, {
		Name: "missing server key",
		Key:  key,
		Objects: []runtime.Object{&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: system.Namespace(),
			},
			Data: map[string][]byte{
				// corev1.TLSPrivateKeyKey:  []byte("missing"),
				corev1.TLSCertKey: []byte("present"),
			},
		}},
		WantErr: true,
	}}

	table.Test(t, MakeFactory(func(ctx context.Context, listers *Listers, cmw configmap.Watcher) controller.Reconciler {
		return &reconciler{
			client:       kubeclient.Get(ctx),
			secretlister: listers.GetSecretLister(),
			key: types.NamespacedName{
				Namespace: system.Namespace(),
				Name:      secretName,
			},
		}
	}))
}

func TestNew(t *testing.T) {
	ctx, _ := SetupFakeContext(t)
	ctx = webhook.WithOptions(ctx, webhook.Options{})

	c := NewController(ctx, configmap.NewStaticWatcher())
	if c == nil {
		t.Fatal("Expected NewController to return a non-nil value")
	}

	if want, got := 0, c.WorkQueue().Len(); want != got {
		t.Errorf("WorkQueue.Len() = %d, wanted %d", got, want)
	}

	la, ok := c.Reconciler.(pkgreconciler.LeaderAware)
	if !ok {
		t.Fatalf("%T is not leader aware", c.Reconciler)
	}

	if err := la.Promote(pkgreconciler.UniversalBucket(), c.MaybeEnqueueBucketKey); err != nil {
		t.Error("Promote() =", err)
	}

	// Queue has async moving parts so if we check at the wrong moment, this might still be 0.
	if wait.PollImmediate(10*time.Millisecond, 250*time.Millisecond, func() (bool, error) {
		return c.WorkQueue().Len() == 1, nil
	}) != nil {
		t.Error("Queue length was never 1")
	}
}

func secretWithCertData(t *testing.T, expiration time.Time) *corev1.Secret {
	const secretName = "webhook-secret"
	serverKey, serverCert, _, err := certresources.CreateCerts(context.Background(), "webhook-service", expiration)
	if err != nil {
		t.Fatal("Failed to create cert:", err)
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: system.Namespace(),
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			corev1.TLSPrivateKeyKey: serverKey,
			corev1.TLSCertKey:       serverCert,
		},
	}
}
