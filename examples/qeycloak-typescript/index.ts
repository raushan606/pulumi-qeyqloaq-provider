import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as qubecloak from "@qeyqloaq/keycloak"

const namespace = "keycloak-local";
const keycloakHost = "localhost";

// --- NAMESPACE ---
const ns = new k8s.core.v1.Namespace(namespace, { metadata: { name: namespace } });

// --- CUSTOM KEYCLOAK PROVIDER ---
const customKeycloakProvider = new qubecloak.Provider("custom-keycloak-provider", {
    url: pulumi.interpolate`http://${keycloakHost}/`,
    username: "admin",
    password: "admin",
    realm: "master",
}, { dependsOn: [ns] });

// --- EXPORTS ---
export const keycloakAdminUrl = pulumi.interpolate`http://${keycloakHost}/admin`;