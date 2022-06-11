/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package resources

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	. "knative.dev/pkg/logging/testing"
)

func TestCreateCerts(t *testing.T) {
	wantCommonName := "got-the-hook"
	sKey, serverCertPEM, caCertBytes, err := CreateCerts(TestContextWithLogger(t), wantCommonName, time.Now().AddDate(0, 0, 7))
	if err != nil {
		t.Fatal("Failed to create certs", err)
	}

	// Test server private key
	p, _ := pem.Decode(sKey)
	if p.Type != "PRIVATE KEY" {
		t.Fatal("Expected the key to be Private key type")
	}
	key, err := x509.ParsePKCS8PrivateKey(p.Bytes)
	if err != nil {
		t.Fatal("Failed to parse private key", err)
	}
	if _, ok := key.(*ecdsa.PrivateKey); !ok {
		t.Fatalf("Key is not ecdsa format, actually %t", key)
	}

	// Test Server Cert
	sCert, err := validCertificate(serverCertPEM, t)
	if err != nil {
		t.Fatal(err)
	}

	// Test CA Cert
	caParsedCert, err := validCertificate(caCertBytes, t)
	if err != nil {
		t.Fatal(err)
	}

	if caParsedCert.Subject.CommonName != wantCommonName {
		t.Fatalf("Unexpected Cert Common Name %q, wanted %q", caParsedCert.Subject.CommonName, wantCommonName)
	}

	// Verify domain names
	expectedDNSNames := []string{wantCommonName}
	if diff := cmp.Diff(caParsedCert.DNSNames, expectedDNSNames); diff != "" {
		t.Fatal("Unexpected CA Cert DNS Name (-want +got) :", diff)
	}

	// Verify Server Cert is Signed by CA Cert
	if err = sCert.CheckSignatureFrom(caParsedCert); err != nil {
		t.Fatal("Failed to verify that the signature on server certificate is from parent CA cert", err)
	}
}

func validCertificate(cert []byte, t *testing.T) (*x509.Certificate, error) {
	t.Helper()
	const certificate = "CERTIFICATE"
	caCert, _ := pem.Decode(cert)
	if caCert.Type != certificate {
		return nil, fmt.Errorf("cert.Type = %s, want: %s", caCert.Type, certificate)
	}
	parsedCert, err := x509.ParseCertificate(caCert.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse cert %w", err)
	}
	if parsedCert.SignatureAlgorithm != x509.ECDSAWithSHA256 {
		return nil, fmt.Errorf("Failed to match signature. Got: %s, want: %s", parsedCert.SignatureAlgorithm, x509.SHA256WithRSA)
	}
	return parsedCert, nil
}
