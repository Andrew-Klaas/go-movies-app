#!/bin/bash

#transit setup
vault secrets enable transit
vault write -f transit/keys/my-key

##database setup
vault secrets enable database

vault write database/config/my-postgresql-database \
plugin_name=postgresql-database-plugin \
allowed_roles="my-role" \
connection_url="postgresql://{{username}}:{{password}}@localhost:5432/movies?sslmode=disable" \
username="goapp" \
password="password"

vault write database/roles/my-role \
db_name=my-postgresql-database \
creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
ALTER USER \"{{name}}\" WITH SUPERUSER;" \
default_ttl="1h" \
max_ttl="24h"

#vault read database/creds/my-role
