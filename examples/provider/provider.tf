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

data "meilisearch_key" "main" {
  uid = "a50aaee3-70dc-49b9-8215-7175cb9eef7a"
}

output "main" {
  value = data.meilisearch_key.main
}
