# Copyright Red Hat

############################
# STEP 1 build executable binary
############################

# Set the base image as the official Go image that already has all the tools and packages to compile and run a Go application
FROM registry.ci.openshift.org/stolostron/builder:go1.17-linux AS builder

# Set the current working directory inside the container
WORKDIR /app

COPY . .

# Fetch dependencies.
# Using go get requires root.
USER root
RUN go get -d -v

# Compile the application
RUN go build -o /idp-configs-api

# Build the migration binary.
RUN go build -o /idp-configs-api-migrate cmd/migrate/migrate.go

############################
# STEP 2 build a small image
############################

FROM registry.redhat.io/ubi8-minimal:latest

COPY --from=builder /idp-configs-api /usr/bin
COPY --from=builder /idp-configs-api-migrate /usr/bin
COPY --from=builder /app/cmd/spec/openapi.json /var/tmp

USER 1001

# Execute the aplication
CMD [ "idp-configs-api" ]
