package provider

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/meilisearch/meilisearch-go"
)

func TestAccIndexResource(t *testing.T) {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: "T35T-M45T3R-K3Y",
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{

				PreConfig: func() {
					_, err := client.DeleteIndex("abcdef")
					if err != nil {
						return
					}
				},

				Config: providerConfig + `
resource "meilisearch_index" "test" {
	uid = "abcdef"
	primary_key = "toto"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify all attributes are set
					resource.TestCheckResourceAttr("meilisearch_index.test", "uid", "abcdef"),
					resource.TestCheckResourceAttr("meilisearch_index.test", "primary_key", "toto"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("meilisearch_index.test", "created_at"),
					resource.TestCheckResourceAttrSet("meilisearch_index.test", "updated_at"),
				),
			},
		},
	})
}
