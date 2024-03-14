# Terraform Provider for Meilisearch.

[![Release](https://img.shields.io/github/v/release/project0/terraform-provider-podman)](https://github.com/project0/terraform-provider-podman/releases)
[![Registry](https://img.shields.io/badge/registry-doc%40latest-lightgrey?logo=terraform)](https://registry.terraform.io/providers/pauldn-wttj/meilisearch/latest/docs)
[![License](https://img.shields.io/badge/license-Mozilla-blue.svg)](https://github.com/project0/terraform-provider-podman/blob/main/LICENSE)

This Terraform provider implements resource management for Meilisearch.

## Overview

### Using the provider

To use this provider, you must install it and provide authentication credentials:

```hcl
terraform {
  required_providers {
    meilisearch = {
      source = "pauldn-wttj/meilisearch"
      version = "0.0.1"
    }
  }
}

provider "meilisearch" {
  host = "http://localhost:7700"
  api_key = "T35T-M45T3R-K3Y"
}
```

Alternatively, you may use environment variables `MEILISEARCH_API_KEY` and / or `MEILISEARCH_HOST` for authentication.
The `MEILISEARCH_API_KEY` should have admin privileges since it may be used to create all kinds of resources.

### Resources

- `meilisearch_api_key`: create and manage API keys for Meilisearch.

### Data sources

- `meilisearch_api_key`: read API keys for Meilisearch.

## Development

_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)._

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20
- [Docker](https://docs.docker.com/engine/install/) and [docker-compose](https://docs.docker.com/compose/install/) >= 3.7 for development
- [golangci-lint](https://golangci-lint.run/usage/install/) for development

### Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

```shell
make testacc
```

Tests are run against a Meilisearch Docker container to ease development, this will create a container on your machine.

### Run linter

Install [golangci-lint](https://golangci-lint.run/usage/install/) and run the linter:

```shell
golangci-lint run
```

Or using Docker:

```shell
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.52.2 golangci-lint run
```
