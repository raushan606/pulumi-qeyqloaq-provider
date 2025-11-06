// Package main runs the provider's gRPC server.
package main

import (
	"context"
	"fmt"
	"os"

	keycloak "github.com/raushan606/pulumi-qeyqloaq-provider/provider"
)

// Serve the provider against Pulumi's Provider protocol.
func main() {
	err := keycloak.Provider().Run(context.Background(), keycloak.Name, keycloak.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
