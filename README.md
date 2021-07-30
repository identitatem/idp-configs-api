# idp-configs-api

A simple HTTP server that serves "Hello world" on http://localhost:8080/api/hello-world-service/v0/ping

---
#### Getting started:
**To run without docker:**

1. Install Golang 1.15
2. Run `make run`

**To run with docker:**
1. Start the docker daemon
2. Export the following environment variables (needed for pulling the base image from redhat.registry.io):
   1. RH_REGISTRY_USER (redhat.registry.io service account user)
   2. RH_REGISTRY_TOKEN (redhat.registry.io service account token)
3. Run `make docker-build`
4. Run `make docker-run`

---

#### Testing and Linting:

**Unit tests:**
```
make test
```
**Lint:**
```
go get -u honnef.co/go/tools/cmd/staticcheck@latest
make lint
```

---
#### Pushing the container image to quay.io:

1. Export the following environment variables
   1. QUAY_USER (quay.io user name)
   2. QUAY_TOKEN (quay.io token/ encrypted CLI password)
   3. IMG="quay.io/${QUAY_USER}/hello-world-service:latest"
2. Run `make docker-push`

---
#### Viewing the OpenAPI 3.0 spec:

* Run `make generate-docs`

The API will serve the docs under a `/docs` endpoint.
