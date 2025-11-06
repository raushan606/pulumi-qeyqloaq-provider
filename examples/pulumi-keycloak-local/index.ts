import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as keycloak from "@pulumi/keycloak";
import * as customKeycloak from "@qeyqloaq/keycloak"
import { RandomPassword } from "@pulumi/random";

// --- CONFIGURATION ---
const config = new pulumi.Config();
const namespace = config.get("namespace") || "keycloak-local";
const dbHost = config.require("dbHost");
const dbPort = config.getNumber("dbPort") || 5432;
const dbUser = config.require("dbUser");
const dbPassword = config.requireSecret("dbPassword");
const dbName = config.require("dbName");
const keycloakHost = config.get("keycloakHost") || "localhost";
const imagePullSecret = config.get("imagePullSecret");

// --- NAMESPACE ---
const ns = new k8s.core.v1.Namespace(namespace, { metadata: { name: namespace } });

// --- DB PASSWORD SECRET ---
const dbPasswordSecret = new k8s.core.v1.Secret("keycloak-db-password", {
    metadata: { namespace },
    stringData: {
        "db-password": dbPassword,
    },
});

// --- CUSTOM KEYCLOAK PROVIDER ---
const customKeycloakProvider = new customKeycloak.Provider("custom-keycloak-provider", {
    url: pulumi.interpolate`http://${keycloakHost}/`,
    username: "admin",
    password: "admin",
    realm: "master",
}, { dependsOn: [ns] });

// --- REALM MANAGEMENT WITH CUSTOM PROVIDER ---
const realmName = "test-realm";
const testRealm = new customKeycloak.Realm(realmName, {
    name: realmName,
    enabled: true,
    displayName: "Test Realm",
    loginTheme: "keycloak",
    accountTheme: "keycloak",
    adminTheme: "keycloak",
    emailTheme: "keycloak",
    // SMTP configuration can be added here if needed
    // smtpServer: {
    //     host: "smtp.example.com",
    //     port: 587,
    //     from: "noreply@example.com",
    //     fromName: "Test Realm",
    //     startTls: true,
    //     auth: false,
    // },
}, { provider: customKeycloakProvider });

// --- STANDARD KEYCLOAK PROVIDER FOR CLIENT MANAGEMENT ---
// The provider will implicitly wait for the realm through output dependencies
const keycloakProvider = new keycloak.Provider("keycloak", {
    url: pulumi.interpolate`http://${keycloakHost}/`,
    clientId: "admin-cli",
    username: "admin",
    password: "admin",
    initialLogin: false,
    clientTimeout: 30,
}, { dependsOn: [testRealm] });

// --- CLIENTS ---
const appMgmtClientSecret = new RandomPassword("app-mgmt-client-secret", {
    length: 16,
    special: true,
    overrideSpecial: "!@$%*()-_=+[]{}<>:?",
});

const appMgmtClient = new keycloak.openid.Client("app-mgmt", {
    realmId: testRealm.name, // Use the name from our custom realm - this creates implicit dependency
    enabled: true,
    accessType: "CONFIDENTIAL",
    name: "Application Management",
    clientId: "app-mgmt",
    clientSecret: appMgmtClientSecret.result,
    description: "Test Application Management",
    rootUrl: pulumi.interpolate`http://${keycloakHost}:8080`,
    adminUrl: pulumi.interpolate`http://${keycloakHost}:8080`,
    baseUrl: "/",
    webOrigins: [pulumi.interpolate`http://${keycloakHost}:8080`],
    validRedirectUris: ["/oauth"],
    validPostLogoutRedirectUris: ["/"],
    standardFlowEnabled: true,
    clientAuthenticatorType: "client-secret",
    frontchannelLogoutEnabled: true,
    useRefreshTokens: true,
    fullScopeAllowed: true,
    loginTheme: "keycloak",
}, { provider: keycloakProvider, dependsOn: [testRealm] });

// --- CLIENT SCOPE ---
const appMgmtClientScope = new keycloak.openid.ClientScope("app-mgmt-client-scope", {
    name: "app-mgmt",
    description: "Test Application Management",
    realmId: testRealm.name,
    includeInTokenScope: true,
    consentScreenText: "Access to Manage applications",
}, { provider: keycloakProvider });

// --- CLIENT ROLES ---
const appMgmtRoles: { [name: string]: keycloak.Role } = {};
const appMgmtRoleDefs = [
    { name: "cloud-user", description: "" },
    { name: "admin", description: "Cloud Admin" },
    { name: "cloud-admin", description: "Internal admin role" },
    { name: "read:deployments", description: "Read Deployments" },
];
for (const roleDef of appMgmtRoleDefs) {
    appMgmtRoles[roleDef.name] = new keycloak.Role(`app-mgmt-${roleDef.name}`, {
        realmId: testRealm.name,
        name: roleDef.name,
        description: roleDef.description,
        clientId: appMgmtClient.id,
    }, { provider: keycloakProvider, dependsOn: [appMgmtClient] });
}

// --- GROUPS ---
const cloudAdmins = new keycloak.Group("cloud-admins", {
    realmId: testRealm.name,
    name: "cloud-admins",
}, { provider: keycloakProvider });

const cloudUsers = new keycloak.Group("cloud-users", {
    realmId: testRealm.name,
    name: "cloud-users",
}, { provider: keycloakProvider });

// --- ASSIGN ROLES TO GROUPS ---
const cloudAdminAppMgmtGR = new keycloak.GroupRoles("cloud-admin-app-mgmt", {
    realmId: testRealm.name,
    groupId: cloudAdmins.id,
    roleIds: [appMgmtRoles["admin"].id],
}, { provider: keycloakProvider, dependsOn: [cloudAdmins, appMgmtRoles["admin"]] });

const cloudUsersAppMgmtGR = new keycloak.GroupRoles("cloud-users-app-mgmt", {
    realmId: testRealm.name,
    groupId: cloudUsers.id,
    roleIds: [appMgmtRoles["cloud-user"].id],
}, { provider: keycloakProvider, dependsOn: [cloudUsers, appMgmtRoles["cloud-user"]] });

// --- DEFAULT GROUP ASSIGNMENT ---
const defaultGroups = new keycloak.DefaultGroups("default-groups", {
    realmId: testRealm.name,
    groupIds: [cloudUsers.id],
}, { provider: keycloakProvider, dependsOn: [cloudUsers] });

// --- EXPORTS ---
export const keycloakAdminUrl = pulumi.interpolate`http://${keycloakHost}/admin`;
export const testRealmName = testRealm.name;
export const testRealmId = testRealm.realmId;
export const appMgmtClientId = appMgmtClient.clientId;
export const appMgmtClientSecretResult = appMgmtClientSecret.result; 