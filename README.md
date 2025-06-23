# OPA Directus

> A tool to configure and adjust OPA policies from auto generated Directus collections.

## Start directus with database and OPA

```bash
mkdir -p database
touch database/data.db
docker compose up
```

## Generate the OPA policies configuration collection

> This reads the bundle/demo/demo.rego, pareses the metadata and generates the Directus collections.

```bash
# build the tool
go build .

# read the .env file
export $(grep -v '^#' .env | xargs -d '\n')

# run the tool
./opa-directus generate
```

Now open http://0.0.0.0:8055/admin/content/pdp-config and adjust the config for the policies.

## Fetch the policies from the OPA server

Now we watch the configuration in directus and download it when it changes:

```bash
./opa-directus fetch --watch
```

## Test the policies

```bash
curl -X POST localhost:8181/v1/data/demo -d '{"input": {"amount": 15}' 
```
