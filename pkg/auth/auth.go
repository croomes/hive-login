package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/croomes/hive-login/pkg/auth/callback"
	"github.com/croomes/hive-login/pkg/auth/login"
)

type AuthProvider struct {
	clientID     string
	clientSecret string
	redirectURI  string

	provider *oidc.Provider

	client *http.Client

	loginHandler    http.Handler
	callbackHandler http.Handler
}

func New(clientID string, clientSecret string, issuerURL string, redirectURI string, rootCAs string) (*AuthProvider, error) {
	client, err := httpClient(rootCAs)
	if err != nil {
		return nil, err
	}

	ctx := oidc.ClientContext(context.Background(), client)
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query provider %q: %v", issuerURL, err)
	}

	var s struct {
		// What scopes does a provider support?
		//
		// See: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
		ScopesSupported []string `json:"scopes_supported"`
	}
	if err := provider.Claims(&s); err != nil {
		return nil, fmt.Errorf("failed to parse provider scopes_supported: %v", err)
	}

	// Does the provider use "offline_access" scope to request a refresh token
	// or does it use "access_type=offline" (e.g. Google)?
	var offlineAsScope bool
	if len(s.ScopesSupported) == 0 {
		// scopes_supported is a "RECOMMENDED" discovery claim, not a required
		// one. If missing, assume that the provider follows the spec and has
		// an "offline_access" scope.
		offlineAsScope = true
	} else {
		// See if scopes_supported has the "offline_access" scope.
		offlineAsScope = func() bool {
			for _, scope := range s.ScopesSupported {
				if scope == oidc.ScopeOfflineAccess {
					return true
				}
			}
			return false
		}()
	}

	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{},
		RedirectURL:  redirectURI,
	}

	return &AuthProvider{
		clientID:        clientID,
		clientSecret:    clientSecret,
		redirectURI:     redirectURI,
		provider:        provider,
		client:          client,
		loginHandler:    login.New(oauth2Config, offlineAsScope),
		callbackHandler: callback.New(redirectURI, client, oauth2Config, provider.Verifier(&oidc.Config{ClientID: clientID})),
	}, nil
}

func (p *AuthProvider) LoginHandler() http.Handler {
	return p.loginHandler
}

func (p *AuthProvider) CallbackHandler() http.Handler {
	return p.callbackHandler
}
