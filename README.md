# Url-Shortening LDAP Auth Service
[![Go](https://github.com/lakhinsu/urlshortening-authservice/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/lakhinsu/urlshortening-authservice/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/lakhinsu/urlshortening-authservice)](https://goreportcard.com/report/github.com/lakhinsu/urlshortening-authservice)

This repository contains a REST API service for URL shortening project (on going). The basic idea behind this project is to create a free alternative of [go/links](https://www.golinks.io/) or bit.ly for organizations.

This service provides user authentication for this project by binding user with a Active Directory server of an organization and generating JWT token for Hasura GraphQL. Authentication has a dedicated service to provide flexibility in authentication. The project will only require a JWT token for GraphQL engine, so a new service can also be implemented as per requirements to authenticate users and generate the token.

This service uses [Gin framework](https://github.com/gin-gonic/gin) to create a REST API.

## Run locally
This repo does not contain a Dockerfile at the moment so i'm listing down steps for running this service locally.
- Get dependencies
`go get .`
- Set environment variables
This service requires following environment variables
    - GIN_MODE -> Gin server mode - "release" or "debug"
    - GIN_HTTPS -> Whether to run API server with HTTPS or not - "true" or "false"
    - GIN_ADDR -> API server addr string - "0.0.0.0" OR <required IP address>
    - GIN_PORT -> port to listen on - "<port number>""
    - LOG_LEVEL -> Log level - "debug" OR "info"
    - LDAP_TLS -> whether to connect with AD server using TLS or not - "true" or "false"
    - LDAP_CONNECTION_STRING -> LDAP server connection string - "<server>:<port>"
    - LDAP_SERVER_NAME -> AD server name, this should match the server name in TLS certificate if LDAP_TLS is enabled.
    - LDAP_SSL_VERIFY -> Whether to verify the SSL certificate of AD server. This is only required when LDAP_TLS is enabled. "true" or "false"
    - LDAP_USER_DN -> User dn for validating users, example: "ou=users,ou=guests,dc=zflexsoftware,dc=com"
    - LDAP_USER_FQDN_ATTRIBUTE_NAME -> Attribute name by which the user's fqdn is stored in AD. example: mail
    - JWT_SECRET -> Secret string to create JWT tokens. Service uses HS256 algorithm
    - JWT_EXPIRE -> JWT 
    - JWT_ISSUER -> Value to set as issuer in JWT tokens

Execute `go run main.go` command at the repository root to start the service.
## Running in Docker
The repo includes Dockerfile and docker-compose.yaml file as samples.
