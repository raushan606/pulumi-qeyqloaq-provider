# Example: Basic Keycloak Realm

This example demonstrates how to create a basic Keycloak realm using the custom provider.

## Prerequisites

1. Keycloak server running (e.g., `docker run -p 8080:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:latest start-dev`)
2. Pulumi Keycloak provider built and installed

## Configuration

Set environment variables:

```bash
export KEYCLOAK_URL="http://localhost:8080"
export KEYCLOAK_USERNAME="admin"
export KEYCLOAK_PASSWORD="admin"
export KEYCLOAK_REALM="master"  # Optional, defaults to master
```

Or configure in your Pulumi program:

```typescript
import * as keycloak from "@pulumi/keycloak";

const config = new pulumi.Config();

const realm = new keycloak.Realm("my-app-realm", {
    name: "my-application",
    enabled: true,
    displayName: "My Application",
    displayNameHtml: "<strong>My Application</strong>",
    loginTheme: "keycloak",
    accountTheme: "keycloak",
    adminTheme: "keycloak",
    emailTheme: "keycloak",
    smtpServer: {
        host: "smtp.gmail.com",
        port: 587,
        from: "noreply@myapp.com",
        fromName: "My Application",
        startTls: true,
        auth: true,
        username: config.requireSecret("smtpUsername"),
        password: config.requireSecret("smtpPassword"),
    }
});

export const realmId = realm.id;
export const realmName = realm.name;
```

## Running

```bash
# Initialize Pulumi project
pulumi new typescript

# Install dependencies
npm install

# Set configuration (if not using environment variables)
pulumi config set keycloak:url http://localhost:8080
pulumi config set keycloak:username admin
pulumi config set --secret keycloak:password admin

# Deploy
pulumi up
```

## Benefits of this approach

1. **Preserves manual changes**: Any settings you don't manage in Pulumi (like security settings, user federation, etc.) are preserved when you make changes.

2. **Selective management**: You only need to define the settings you want to manage in code.

3. **Safe updates**: Updates only modify the fields you specify, leaving everything else untouched.

4. **Works with existing realms**: You can import existing realms and start managing only specific aspects.

## What gets managed vs preserved

**Managed by this provider** (will be overwritten):
- realm name
- enabled status  
- display name and HTML display name
- themes (login, account, admin, email)
- SMTP server configuration

**Preserved** (won't be touched):
- User federation settings
- Authentication flows
- Security defenses
- Login settings (like remember me, registration allowed, etc.)
- Token settings
- All other realm configurations

This gives you the flexibility to manage infrastructure-level settings with Pulumi while allowing teams to manage user-facing configurations through the Keycloak UI.