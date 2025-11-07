package provider

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// ProviderConfig holds the configuration for the Keycloak provider
type ProviderConfig struct {
	URL      string  `pulumi:"url"`               // Keycloak server URL (required)
	Username string  `pulumi:"username"`          // Keycloak admin username (required)
	Password string  `pulumi:"password"`          // Keycloak admin password (required)
	Realm    *string `pulumi:"realm,optional"`    // Keycloak admin realm (optional, defaults to "master")
	BasePath *string `pulumi:"basePath,optional"` // Base path for Keycloak (optional, defaults to "/")
	Insecure *bool   `pulumi:"insecure,optional"` // Whether to use insecure connections (optional, defaults to false)
}

// Annotate provides schema documentation for ProviderConfig
func (config *ProviderConfig) Annotate(a infer.Annotator) {
	a.Describe(&config.URL, "Keycloak server URL (e.g., http://localhost:8080)")
	a.Describe(&config.Username, "Keycloak admin username")
	a.Describe(&config.Password, "Keycloak admin password")
	a.Describe(&config.Realm, "Keycloak admin realm")
	a.Describe(&config.BasePath, "Base path for Keycloak API")
	a.Describe(&config.Insecure, "Whether to allow insecure connections")

	a.SetDefault(&config.Realm, "master")
	a.SetDefault(&config.BasePath, "/")
	a.SetDefault(&config.Insecure, false)
}

// KeycloakProvider represents the main provider struct
type KeycloakProvider struct {
	config *ProviderConfig
	client *gocloak.GoCloak
	token  *gocloak.JWT
}

// configKey is used to store the provider config in context
type configKey struct{}

// clientKey is used to store the Keycloak client in context
type clientKey struct{}

// Configure sets up the provider with the given configuration
func (p *KeycloakProvider) Configure(ctx context.Context, config ProviderConfig) error {
	if config.URL == "" {
		return fmt.Errorf("keycloak URL is required")
	}
	if config.Username == "" {
		return fmt.Errorf("keycloak username is required")
	}
	if config.Password == "" {
		return fmt.Errorf("keycloak password is required")
	}

	// Set defaults
	if config.Realm == nil {
		defaultRealm := "master"
		config.Realm = &defaultRealm
	}
	if config.BasePath == nil {
		defaultBasePath := "/"
		config.BasePath = &defaultBasePath
	}
	if config.Insecure == nil {
		defaultInsecure := false
		config.Insecure = &defaultInsecure
	}

	p.config = &config

	client := gocloak.NewClient(config.URL)

	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Keycloak: %w", err)
	}

	p.client = client
	p.token = token

	return nil
}

func getKeycloakClient(ctx context.Context) *gocloak.GoCloak {
	if client, ok := ctx.Value(clientKey{}).(*gocloak.GoCloak); ok {
		return client
	}
	panic("Keycloak client not found in context")
}

func getProviderConfig(ctx context.Context) *ProviderConfig {
	if config, ok := ctx.Value(configKey{}).(*ProviderConfig); ok {
		return config
	}
	panic("Provider config not found in context")
}
