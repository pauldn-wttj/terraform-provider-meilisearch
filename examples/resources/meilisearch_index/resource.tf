# Create a Meilisearch Index
resource "meilisearch_index" "example" {
	uid = "index-name"
	primary_key = "key-name"
}
