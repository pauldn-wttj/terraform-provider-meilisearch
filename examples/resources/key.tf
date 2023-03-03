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

resource "meilisearch_key" "example" {
  name = "test"
  description = "this is a description"
  actions = ["*"]
  indexes = ["*"]
}
