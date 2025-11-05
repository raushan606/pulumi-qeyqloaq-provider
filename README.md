# Pulumi Keycloak Provider

A Pulumi provider for managing Keycloak realms.

## Features

- ✅ Realm management (Create, Read, Update, Delete)
- 🔐 Secure authentication with Keycloak Admin API
- 📝 Full Pulumi schema support

## Installation

```bash
# Install the provider binary
make install

# Or build from source
make build
```

## Usage

```typescript
import * as keycloak from "@pulumi/keycloak";

const realm = new keycloak.Realm("my-realm", {
    name: "my-application",
    displayName: "My Application",
    enabled: true,
    description: "Application realm for my service"
});
```

## Configuration

Set these environment variables or configure in your Pulumi program:

- `KEYCLOAK_URL`: Keycloak server URL (e.g., `http://localhost:8080`)
- `KEYCLOAK_USERNAME`: Admin username
- `KEYCLOAK_PASSWORD`: Admin password
- `KEYCLOAK_REALM`: Admin realm (default: `master`)

## Development

```bash
# Install dependencies
go mod tidy

# Build the provider
make build

# Generate schema
make schema

# Run tests
make test
```