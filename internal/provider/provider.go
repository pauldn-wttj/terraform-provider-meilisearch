package provider

import (
	"context"
	"os"

	"github.com/meilisearch/meilisearch-go"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure MeilisearchProvider satisfies various provider interfaces.
var _ provider.Provider = &MeilisearchProvider{}

// MeilisearchProvider defines the provider implementation.
type MeilisearchProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MeilisearchProviderModel describes the provider data model.
type MeilisearchProviderModel struct {
	Host   types.String `tfsdk:"host"`
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *MeilisearchProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "meilisearch"
	resp.Version = p.version
}

func (p *MeilisearchProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "Host of Meilisearch server. May also be provided via MEILISEARCH_HOST environment variable.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "Meilisearch master API key. May also be provided via MEILISEARCH_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *MeilisearchProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Meilisearch client")

	var config MeilisearchProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Meilisearch host",
			"The provider cannot create the Meilisearch API client as there is an unknown configuration value for the Meilisearch host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MEILISEARCH_HOST environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Meilisearch API key",
			"The provider cannot create the Meilisearch API client as there is an unknown configuration value for the Meilisearch API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MEILISEARCH_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("MEILISEARCH_HOST")
	apiKey := os.Getenv("MEILISEARCH_API_KEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Meilisearch host",
			"The provider cannot create the Meilisearch API client as there is a missing or empty value for the Meilisearch host. "+
				"Set the host value in the configuration or use the MEILISEARCH_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Meilisearch API key",
			"The provider cannot create the Meilisearch API client as there is a missing or empty value for the Meilisearch API key. "+
				"Set the API key value in the configuration or use the MEILISEARCH_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "meilisearch_host", host)
	ctx = tflog.SetField(ctx, "meilisearch_api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "meilisearch_api_key")

	tflog.Debug(ctx, "Creating Meilisearch client")

	// Create a new Meilisearch client using the configuration values
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   host,
		APIKey: apiKey,
	})

	// Make the Meilisearch client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Meilisearch client", map[string]any{"success": true})
}

func (p *MeilisearchProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKeyResource,
	}
}

func (p *MeilisearchProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewKeyDataSource,
		NewIndexDataSource,
		NewVersionDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MeilisearchProvider{
			version: version,
		}
	}
}
