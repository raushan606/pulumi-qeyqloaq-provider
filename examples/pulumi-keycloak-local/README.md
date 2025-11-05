# Pulumi Keycloak Local Example

This example demonstrates how to use the custom Pulumi Keycloak provider to manage Keycloak realms alongside the standard Keycloak provider for client management.

## Architecture

- **Custom Keycloak Provider**: Used for realm management with merge strategy
- **Standard Keycloak Provider**: Used for client, role, and group management
- **Kubernetes Deployment**: Deploys Keycloak using Helm chart
- **External Database**: Connects to external PostgreSQL database

## Key Changes from Previous Version

1. **Removed Dynamic Provider**: Replaced the custom dynamic provider with the proper custom Keycloak provider
2. **Cleaner Code**: Removed all the manual HTTP calls and complex state management
3. **Better Type Safety**: Leverages generated TypeScript types from the provider
4. **Merge Strategy**: The custom provider only manages specified fields, preserving manual changes

## Configuration

Set the following configuration values in your Pulumi stack:

```bash
pulumi config set dbHost "your-postgres-host"
pulumi config set dbUser "your-postgres-user"
pulumi config set --secret dbPassword "your-postgres-password"
pulumi config set dbName "your-postgres-database"
pulumi config set keycloakHost "localhost"  # Optional, defaults to localhost
pulumi config set namespace "keycloak-local"  # Optional
pulumi config set imagePullSecret "your-image-pull-secret"  # Optional
```

## Usage

1. Build the custom provider and generate SDKs:
   ```bash
   make generate-sdk-typescript
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Deploy the stack:
   ```bash
   pulumi up
   ```

## Resources Created

### Infrastructure
- Kubernetes namespace
- Keycloak Helm release with external PostgreSQL
- Database password secret

### Keycloak Configuration
- **Realm**: `test-realm` (managed by custom provider)
- **Client**: `app-mgmt` (confidential client)
- **Roles**: Various application roles
- **Groups**: `cloud-admins`, `cloud-users`
- **Default Group Assignment**: New users automatically added to `cloud-users`

## Provider Benefits

The custom provider offers several advantages:

1. **Merge Strategy**: Only manages fields you specify, preserving manual changes
2. **Simplified Configuration**: No need for manual HTTP calls or token management
3. **Better Error Handling**: Built-in retry logic and proper error messages
4. **Type Safety**: Full TypeScript support with generated types

## Outputs

- `keycloakAdminUrl`: Admin console URL
- `testRealmName`: Name of the created realm
- `testRealmId`: ID of the created realm
- `appMgmtClientId`: Client ID for the application management client
- `appMgmtClientSecretResult`: Client secret (sensitive)

## Development

To modify the example:

1. Update the realm configuration in the custom provider section
2. Add additional clients, roles, or groups using the standard provider
3. The custom provider handles realm-level settings while the standard provider manages realm contents