package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "meilisearch" {
  host 		= "http://localhost:7700"
  api_key = "T35T-M45T3R-K3Y"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"meilisearch": providerserver.NewProtocol6WithError(New("dev")()),
	}
)
