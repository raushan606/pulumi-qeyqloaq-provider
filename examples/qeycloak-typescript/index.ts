import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as qubecloak from "@qeyqloaq/keycloak"

const namespace = "keycloak-local";
const keycloakHost = "localhost:8080";

const ns = new k8s.core.v1.Namespace(namespace, { metadata: { name: namespace } });

const customKeycloakProvider = new qubecloak.Provider("custom-keycloak-provider", {
    url: pulumi.interpolate`http://${keycloakHost}/`,
    username: "admin",
    password: "admin",
    realm: "master",
}, { dependsOn: [ns] });

const realm = new qubecloak.Realm("qubecloak-realm", {
    name: "payara-qube",
    enabled: true,
    displayName: "Payara Qube",
    displayNameHtml: "<div class=\"kc-logo-text\"><span>Payara Qube</span></div>",
    loginTheme: "payara",
    accountTheme: "payara",
    adminTheme: "payara",
    emailTheme: "payara",
}, { provider: customKeycloakProvider, dependsOn: [ns] });


// --- EXPORTS ---
export const keycloakAdminUrl = pulumi.interpolate`http://${keycloakHost}/admin`;
export const realmOutput = realm;