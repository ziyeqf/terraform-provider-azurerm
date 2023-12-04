// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type ClientCertificateAuthorizerOptions struct {
	// Environment is the Azure environment/cloud being targeted
	Environment environments.Environment

	// Api describes the Azure API being used
	Api environments.Api

	// TenantId is the tenant to authenticate against
	TenantId string

	// AuxTenantIds lists additional tenants to authenticate against, currently only
	// used for Resource Manager when auxiliary tenants are needed.
	// e.g. https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/authenticate-multi-tenant
	AuxTenantIds []string

	// ClientId is the client ID used when authenticating
	ClientId string

	// Pkcs12Data is the binary PKCS#12 archive data containing the certificate and private key
	Pkcs12Data []byte

	// Pkcs12Path is a path to a binary PKCS#12 archive on the filesystem
	Pkcs12Path string

	// Pkcs12Pass is the challenge passphrase to decrypt the PKCS#12 archive
	Pkcs12Pass string
}

// NewClientCertificateAuthorizer returns an authorizer which uses client certificate authentication.
func NewClientCertificateAuthorizer(ctx context.Context, options ClientCertificateAuthorizerOptions) (Authorizer, error) {
	if len(options.Pkcs12Data) == 0 {
		var err error
		b, err := os.ReadFile(options.Pkcs12Path)
		if err != nil {
			return nil, fmt.Errorf("could not read PKCS#12 archive at %q: %s", options.Pkcs12Path, err)
		}

		// Try to base64 decode the content, in case it is encoded
		buf := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
		if n, err := base64.StdEncoding.Decode(buf, b); err == nil {
			b = buf[:n]
		}

		options.Pkcs12Data = b
	}

	certs, key, err := azidentity.ParseCertificates(options.Pkcs12Data, []byte(options.Pkcs12Pass))
	if err != nil {
		return nil, fmt.Errorf("could not decode PKCS#12 archive: %s", err)
	}

	k, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key must be an RSA key")
	}
	var (
		certificate *x509.Certificate
		x5c         []string
	)
	for _, cert := range certs {
		if cert == nil {
			// not returning an error here because certs may still contain a sufficient cert/key pair
			continue
		}
		certKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if ok && k.E == certKey.E && k.N.Cmp(certKey.N) == 0 {
			// We know this is the signing cert because its public key matches the given private key.
			// This cert must be first in x5c.
			certificate = cert
			x5c = append([]string{base64.StdEncoding.EncodeToString(cert.Raw)}, x5c...)
		} else {
			x5c = append(x5c, base64.StdEncoding.EncodeToString(cert.Raw))
		}
	}
	if certificate == nil {
		return nil, fmt.Errorf("key doesn't match any certificate")
	}

	scope, err := environments.Scope(options.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", options.Api.Name(), err)
	}

	conf := clientCredentialsConfig{
		Environment:        options.Environment,
		TenantID:           options.TenantId,
		AuxiliaryTenantIDs: options.AuxTenantIds,
		ClientID:           options.ClientId,
		PrivateKey:         key,
		Certificate:        certificate,
		X5C:                x5c,
		Scopes: []string{
			*scope,
		},
	}
	return conf.TokenSource(ctx, clientCredentialsAssertionType)
}
