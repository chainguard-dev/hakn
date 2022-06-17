/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package resources

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	. "knative.dev/pkg/logging/testing"
)

func TestMakeSecret(t *testing.T) {
	ctx := TestContextWithLogger(t)
	secret, err := MakeSecret(ctx, "foo", "ns")
	if err != nil {
		t.Error("MakeSecret() =", err)
	}

	for _, key := range []string{corev1.TLSCertKey, corev1.TLSPrivateKeyKey} {
		if _, ok := secret.Data[key]; !ok {
			t.Errorf("secret.Data[%q] is missing", key)
		}
	}
}
