package provider

import (
	"fmt"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/raushan606/pulumi-qeyqloaq-provider/provider/version"
)

var Version string = version.Version

const Name string = "keycloak"

// Provider creates a new instance of the provider.
func Provider() p.Provider {
	p, err := infer.NewProviderBuilder().
		WithDisplayName("pulumi-keycloak").
		WithDescription("A Pulumi provider for managing Keycloak resources.").
		WithHomepage("https://github.com/raushan606/pulumi-qeyqloaq-provider").
		WithNamespace("qeyqloaq").
		WithResources(infer.Resource(&Realm{})).
		WithConfig(infer.Config(&ProviderConfig{})).
		WithModuleMap(map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		}).Build()
	if err != nil {
		panic(fmt.Errorf("unable to build provider: %w", err))
	}
	return p
}
