package client

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_WhenTLSSkipVerifyIsEnabled_ThenInsecureSkipVerifyIsTrue(t *testing.T) {
	clientID := "someClientID"
	clientSecret := "someClientSecret"
	controlPlane := "someControlPlane"
	tlsSkipVerify := true

	client, err := New(clientID, clientSecret, controlPlane, tlsSkipVerify)

	require.NoError(t, err)

	expectedClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: tlsSkipVerify,
			},
		},
	}

	assert.Equal(t, controlPlane, client.ControlPlane)
	assert.Equal(t, expectedClient, client.client)
}

func TestNewClient_WhenTLSSkipVerifyIsDisabled_ThenInsecureSkipVerifyIsFalse(t *testing.T) {
	clientID := "someClientID"
	clientSecret := "someClientSecret"
	controlPlane := "someControlPlane"
	tlsSkipVerify := false

	client, err := New(clientID, clientSecret, controlPlane, tlsSkipVerify)

	require.NoError(t, err)

	expectedClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: tlsSkipVerify,
			},
		},
	}

	assert.Equal(t, controlPlane, client.ControlPlane)
	assert.Equal(t, expectedClient, client.client)
}

func TestNewClient_WhenClientIDIsEmpty_ThenThrowError(t *testing.T) {
	clientID := ""
	clientSecret := "someClientSecret"
	controlPlane := "someControlPlane"
	tlsSkipVerify := false

	client, err := New(clientID, clientSecret, controlPlane, tlsSkipVerify)

	expectedErrorMessage := "clientID, clientSecret and controlPlane must have non-empty values"

	assert.Nil(t, client)
	assert.EqualError(t, err, expectedErrorMessage)
}

func TestNewClient_WhenClientSecretIsEmpty_ThenThrowError(t *testing.T) {
	clientID := "someClientID"
	clientSecret := ""
	controlPlane := "someControlPlane"
	tlsSkipVerify := false

	client, err := New(clientID, clientSecret, controlPlane, tlsSkipVerify)

	expectedErrorMessage := "clientID, clientSecret and controlPlane must have non-empty values"

	assert.Nil(t, client)
	assert.EqualError(t, err, expectedErrorMessage)
}

func TestNewClient_WhenControlPlaneIsEmpty_ThenThrowError(t *testing.T) {
	clientID := "someClientID"
	clientSecret := "someClientSecret"
	controlPlane := ""
	tlsSkipVerify := false

	client, err := New(clientID, clientSecret, controlPlane, tlsSkipVerify)

	expectedErrorMessage := "clientID, clientSecret and controlPlane must have non-empty values"

	assert.Nil(t, client)
	assert.EqualError(t, err, expectedErrorMessage)
}
