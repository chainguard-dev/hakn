/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package resources

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	oneWeek = 7 * 24 * time.Hour
)

// MakeSecret synthesizes a Kubernetes Secret object with the keys specified by
// ServerKey, ServerCert, and CACert populated with a fresh certificate.
// This is mutable to make deterministic testing possible.
var MakeSecret = MakeSecretInternal

// MakeSecretInternal is only public so MakeSecret can be restored in testing.  Use MakeSecret.
func MakeSecretInternal(ctx context.Context, name, namespace string) (*corev1.Secret, error) {
	serverKey, serverCert, _, err := CreateCerts(ctx, "*", time.Now().Add(oneWeek))
	if err != nil {
		return nil, err
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			corev1.TLSCertKey:       serverCert,
			corev1.TLSPrivateKeyKey: serverKey,
		},
	}, nil
}
