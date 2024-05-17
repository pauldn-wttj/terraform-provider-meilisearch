package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
data "meilisearch_index" "test" {
	uid = "test_index"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify all attributes are set
					resource.TestCheckResourceAttr("data.meilisearch_index.test", "uid", "test_index"),
					resource.TestCheckResourceAttr("data.meilisearch_index.test", "primary_key", "test_id"),
					// Verify ID placeholder attribute is set
					resource.TestCheckResourceAttr("data.meilisearch_index.test", "id", "placeholder"),
				),
			},
			// Read testing when no primary key is specified
			{
				Config: providerConfig + `
data "meilisearch_index" "test" {
	uid = "test_index_no_primary_key"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify all attributes are set
					resource.TestCheckResourceAttr("data.meilisearch_index.test", "uid", "test_index_no_primary_key"),
					resource.TestCheckNoResourceAttr("data.meilisearch_index.test", "primary_key"),
					// Verify ID placeholder attribute is set
					resource.TestCheckResourceAttr("data.meilisearch_index.test", "id", "placeholder"),
				),
			},
		},
	})
}
