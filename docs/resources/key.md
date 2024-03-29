---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meilisearch_key Resource - meilisearch"
subcategory: ""
description: |-
  Manages a Meilisearch API key.
---

# meilisearch_key (Resource)

Manages a Meilisearch API key.

## Example Usage

```terraform
# Create a Meilisearch API key
resource "meilisearch_key" "example" {
	uid = "11111111-2222-3333-4444-555555555555"
	name = "example"
	description = "Example description"
	actions = ["search"]
  indexes = ["example_index_1", "example_index_2"]
	expires_at = "2042-04-02T00:42:42Z"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `actions` (List of String) Actions permitted for the key.
- `indexes` (List of String) Indexes the key is authorized to act on (with the actions specified in the scope of the key).

### Optional

- `description` (String) Description of the key.
- `expires_at` (String) Date and time when the key will expire (RFC3339)
- `name` (String) Name of the key.
- `uid` (String) UID (uuid v4) used by Meilisearch to identify the key.

### Read-Only

- `created_at` (String) Date and time when the key was created (RFC3339)
- `id` (String) Placeholder identifier attribute.
- `key` (String, Sensitive) Actual key value.
- `updated_at` (String) Date and time when the key was last updated (RFC3339)

## Import

Import is supported using the following syntax:

```shell
# Keys can be imported by specifying the UID used by Meilisearch.
terraform import meilisearch_key.example 11111111-2222-3333-4444-555555555555
```
