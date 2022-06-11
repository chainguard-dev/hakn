/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package certificates

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"

	certresources "github.com/chainguard-dev/hakn/pkg/reconciler/certificates/resources"
)

const (
	// Time used for updating a certificate before it expires.
	oneDay = 24 * time.Hour
)

type reconciler struct {
	pkgreconciler.LeaderAwareFuncs

	client       kubernetes.Interface
	secretlister corelisters.SecretLister
	key          types.NamespacedName
}

var _ controller.Reconciler = (*reconciler)(nil)
var _ pkgreconciler.LeaderAware = (*reconciler)(nil)

// Reconcile implements controller.Reconciler
func (r *reconciler) Reconcile(ctx context.Context, key string) error {
	if r.IsLeaderFor(r.key) {
		// only reconciler the certificate when we are leader.
		return r.reconcileCertificate(ctx)
	}
	return controller.NewSkipKey(key)
}

func (r *reconciler) reconcileCertificate(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	secret, err := r.secretlister.Secrets(r.key.Namespace).Get(r.key.Name)
	if apierrors.IsNotFound(err) {
		// The secret should be created explicitly by a higher-level system
		// that's responsible for install/updates.  We simply populate the
		// secret information.
		return nil
	} else if err != nil {
		logger.Errorf("Error accessing certificate secret %q: %v", r.key.Name, err)
		return err
	}

	if _, haskey := secret.Data[corev1.TLSPrivateKeyKey]; !haskey {
		logger.Infof("Certificate secret %q is missing key %q", r.key.Name, corev1.TLSPrivateKeyKey)
	} else if _, haskey := secret.Data[corev1.TLSCertKey]; !haskey {
		logger.Infof("Certificate secret %q is missing key %q", r.key.Name, corev1.TLSCertKey)
	} else {
		// Check the expiration date of the certificate to see if it needs to be updated
		cert, err := tls.X509KeyPair(secret.Data[corev1.TLSCertKey], secret.Data[corev1.TLSPrivateKeyKey])
		if err != nil {
			logger.Warnw("Error creating pem from certificate and key", zap.Error(err))
		} else {
			certData, err := x509.ParseCertificate(cert.Certificate[0])
			if err != nil {
				logger.Errorw("Error parsing certificate", zap.Error(err))
			} else if time.Now().Add(oneDay).Before(certData.NotAfter) {
				return nil
			}
		}
	}
	// Don't modify the informer copy.
	secret = secret.DeepCopy()

	// One of the secret's keys is missing, so synthesize a new one and update the secret.
	newSecret, err := certresources.MakeSecret(ctx, r.key.Name, r.key.Namespace)
	if err != nil {
		return err
	}
	secret.Data = newSecret.Data
	_, err = r.client.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
	return err
}
