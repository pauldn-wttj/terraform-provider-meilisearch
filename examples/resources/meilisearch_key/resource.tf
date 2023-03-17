# Create a Meilisearch API key
resource "meilisearch_key" "example" {
	uid = "11111111-2222-3333-4444-555555555555"
	name = "example"
	description = "Example description"
	actions = ["search"]
  indexes = ["example_index_1", "example_index_2"]
	expires_at = "2042-04-02T00:42:42Z"
}
