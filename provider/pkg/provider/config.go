package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// ProviderConfig holds the configuration for the Keycloak provider
type ProviderConfig struct {
	// Keycloak server URL (required)
	URL string `pulumi:"url"`
	
	// Admin username (required)
	Username string `pulumi:"username"`
	
	// Admin password (required)  
	Password string `pulumi:"password"`
	
	// Admin realm (optional, defaults to "master")
	Realm *string `pulumi:"realm,optional"`
	
	// Base path for Keycloak (optional, defaults to "/")
	BasePath *string `pulumi:"basePath,optional"`
	
	// Whether to use insecure connections (optional, defaults to false)
	Insecure *bool `pulumi:"insecure,optional"`
}

// Annotate provides schema documentation for ProviderConfig
func (config *ProviderConfig) Annotate(a infer.Annotator) {
	a.Describe(&config.URL, "Keycloak server URL (e.g., http://localhost:8080)")
	a.Describe(&config.Username, "Keycloak admin username")
	a.Describe(&config.Password, "Keycloak admin password")
	a.Describe(&config.Realm, "Keycloak admin realm")
	a.Describe(&config.BasePath, "Base path for Keycloak API")
	a.Describe(&config.Insecure, "Whether to allow insecure connections")
	
	// Set default values
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
	// Validate required fields
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
	
	// Create Keycloak client
	client := gocloak.NewClient(config.URL)
	
	// Authenticate and get token
	token, err := client.LoginAdmin(ctx, config.Username, config.Password, *config.Realm)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Keycloak: %w", err)
	}
	
	p.client = client
	p.token = token
	
	return nil
}

// Helper function to get Keycloak client from context
func getKeycloakClient(ctx context.Context) *gocloak.GoCloak {
	if client, ok := ctx.Value(clientKey{}).(*gocloak.GoCloak); ok {
		return client
	}
	panic("Keycloak client not found in context")
}

// Helper function to get provider config from context
func getProviderConfig(ctx context.Context) *ProviderConfig {
	if config, ok := ctx.Value(configKey{}).(*ProviderConfig); ok {
		return config
	}
	panic("Provider config not found in context")
}

// Helper function to create context with client and config
func withKeycloakContext(ctx context.Context, provider *KeycloakProvider) context.Context {
	ctx = context.WithValue(ctx, configKey{}, provider.config)
	ctx = context.WithValue(ctx, clientKey{}, provider.client)
	return ctx
}

// GetProviderConfigFromEnv loads provider configuration from environment variables
func GetProviderConfigFromEnv() ProviderConfig {
	config := ProviderConfig{
		URL:      getEnvOrDefault("KEYCLOAK_URL", ""),
		Username: getEnvOrDefault("KEYCLOAK_USERNAME", ""),
		Password: getEnvOrDefault("KEYCLOAK_PASSWORD", ""),
	}
	
	if realm := getEnvOrDefault("KEYCLOAK_REALM", ""); realm != "" {
		config.Realm = &realm
	}
	
	if basePath := getEnvOrDefault("KEYCLOAK_BASE_PATH", ""); basePath != "" {
		config.BasePath = &basePath
	}
	
	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}