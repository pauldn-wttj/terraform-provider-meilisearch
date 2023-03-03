terraform {
  required_providers {
    meilisearch = {
      source = "hashicorp.com/edu/meilisearch"
    }
  }
}

provider "meilisearch" {
  host = "http://localhost:7700"
  api_key = "T35T-M45T3R-K3Y"
}

data "meilisearch_key" "main" {}
