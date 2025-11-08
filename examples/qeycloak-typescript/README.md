## Example of custom keycloak in Typescript

First do some local setup then do pulumi up.

### Start postgres instance

```bash
docker run -d --rm \
  --name keycloak-postgres \
  -e POSTGRES_USER=keycloak \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=keycloak \
  -p 5432:5432 \
  postgres:15
```

### Start Qube Keycloak Instance

```bash
docker run -it --rm \
  -e KC_BOOTSTRAP_ADMIN_USERNAME=admin \
  -e KC_BOOTSTRAP_ADMIN_PASSWORD=admin \
  -e KC_DB_URL_HOST=host.docker.internal \
  -e KC_DB_URL_DATABASE=keycloak \
  -e KC_DB_URL_PORT=5432 \
  -e KC_DB_USERNAME=keycloak \
  -e KC_DB_PASSWORD=secret \
  -e KC_HTTP_ENABLED=true \
  -e KC_HOSTNAME_STRICT=false \
  -p 8080:8080 \
  -p 9000:9000 \
  payara.azurecr.io/payara-qube/keycloak:26.3.2.1 start-dev
```

Now do `pulumi up`