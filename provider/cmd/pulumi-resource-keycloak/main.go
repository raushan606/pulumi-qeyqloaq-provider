package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/raushan606/pulumi-keycloak/provider/pkg/provider"
)

func main() {
	// Create the provider
	prov, err := infer.NewProviderBuilder().
		WithResources(
			infer.Resource(&provider.Realm{}),
		).
		WithNamespace("keycloak").
		Build()
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building provider: %s\n", err.Error())
		os.Exit(1)
	}
	
	// Run the provider
	err = prov.Run(context.Background(), "keycloak", "0.1.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}