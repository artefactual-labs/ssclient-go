package ssclient

import (
	"context"
	"fmt"
	"net/http"

	kabs "github.com/microsoft/kiota-abstractions-go"
	khttp "github.com/microsoft/kiota-http-go"

	"go.artefactual.dev/ssclient/kiota"
)

func New(httpClient *http.Client, baseURL, username, key string) (*kiota.Client, error) {
	authProvider := NewStorageServiceAuthProvider(username, key)
	adapter, err := khttp.NewNetHttpRequestAdapter(authProvider)
	if err != nil {
		return nil, fmt.Errorf("create client adapter: %v", err)
	}
	adapter.SetBaseUrl(baseURL)

	return kiota.NewClient(adapter), nil
}

type authProvider struct {
	username string
	key      string
}

func NewStorageServiceAuthProvider(username, key string) *authProvider {
	return &authProvider{
		username: username,
		key:      key,
	}
}

func (p *authProvider) AuthenticateRequest(ctx context.Context, request *kabs.RequestInformation, additionalAuthenticationContext map[string]interface{}) error {
	request.Headers.Add("Authorization", fmt.Sprintf("ApiKey %s:%s", p.username, p.key))

	return nil
}
