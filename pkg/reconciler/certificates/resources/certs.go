/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package resources

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"time"

	"go.uber.org/zap"

	"knative.dev/pkg/logging"
)

// The logic below here is adapted from:
// https://github.com/knative/pkg/blob/main/webhook/certificates/resources/certs.go

func createCertTemplate(commonName string, notAfter time.Time) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Chainguard, Inc."},
			CommonName:   commonName,
		},
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		NotBefore:             time.Now(),
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		DNSNames:              []string{commonName},
	}
	return &tmpl, nil
}

// Create cert template suitable for CA and hence signing
func createCACertTemplate(commonName string, notAfter time.Time) (*x509.Certificate, error) {
	rootCert, err := createCertTemplate(commonName, notAfter)
	if err != nil {
		return nil, err
	}
	// Make it into a CA cert and change it so we can use it to sign certs
	rootCert.IsCA = true
	rootCert.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	return rootCert, nil
}

// Create cert template that we can use on the server for TLS
func createServerCertTemplate(commonName string, notAfter time.Time) (*x509.Certificate, error) {
	serverCert, err := createCertTemplate(commonName, notAfter)
	if err != nil {
		return nil, err
	}
	serverCert.KeyUsage = x509.KeyUsageDigitalSignature
	serverCert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	return serverCert, err
}

// Actually sign the cert and return things in a form that we can use later on
func createCert(template, parent *x509.Certificate, pub, parentPriv interface{}) (
	cert *x509.Certificate, certPEM []byte, err error) {
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return
	}
	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	}
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}

func createCA(ctx context.Context, commonName string, notAfter time.Time) (*ecdsa.PrivateKey, *x509.Certificate, []byte, error) {
	logger := logging.FromContext(ctx)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Errorw("error generating random key", zap.Error(err))
		return nil, nil, nil, err
	}
	publicKey := privateKey.Public()

	rootCertTmpl, err := createCACertTemplate(commonName, notAfter)
	if err != nil {
		logger.Errorw("error generating CA cert", zap.Error(err))
		return nil, nil, nil, err
	}

	rootCert, rootCertPEM, err := createCert(rootCertTmpl, rootCertTmpl, publicKey, privateKey)
	if err != nil {
		logger.Errorw("error signing the CA cert", zap.Error(err))
		return nil, nil, nil, err
	}
	return privateKey, rootCert, rootCertPEM, nil
}

func CreateCerts(ctx context.Context, commonName string, notAfter time.Time) (serverKey, serverCert, caCert []byte, err error) {
	logger := logging.FromContext(ctx)
	// First create a CA certificate and private key
	caKey, caCertificate, caCertificatePEM, err := createCA(ctx, commonName, notAfter)
	if err != nil {
		return nil, nil, nil, err
	}

	// Then create the private key for the serving cert
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Errorw("error generating random key", zap.Error(err))
		return nil, nil, nil, err
	}
	publicKey := privateKey.Public()

	servCertTemplate, err := createServerCertTemplate(commonName, notAfter)
	if err != nil {
		logger.Errorw("failed to create the server certificate template", zap.Error(err))
		return nil, nil, nil, err
	}

	// create a certificate which wraps the server's public key, sign it with the CA private key
	_, servCertPEM, err := createCert(servCertTemplate, caCertificate, publicKey, caKey)
	if err != nil {
		logger.Errorw("error signing server certificate template", zap.Error(err))
		return nil, nil, nil, err
	}
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		logger.Errorw("error marshaling private key", zap.Error(err))
		return nil, nil, nil, err
	}
	servKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "PRIVATE KEY", Bytes: privKeyBytes,
	})
	return servKeyPEM, servCertPEM, caCertificatePEM, nil
}
