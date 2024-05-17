package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/meilisearch/meilisearch-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &indexDataSource{}
	_ datasource.DataSourceWithConfigure = &indexDataSource{}
)

func NewIndexDataSource() datasource.DataSource {
	return &indexDataSource{}
}

// indexDataSource defines the data source implementation.
type indexDataSource struct {
	client *meilisearch.Client
}

type indexDataSourceModel struct {
	UID        types.String `tfsdk:"uid"`
	PrimaryKey types.String `tfsdk:"primary_key"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
	ID         types.String `tfsdk:"id"`
}

func (d *indexDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_index"
}

func (d *indexDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves details of a Meilisearch index.",
		Attributes: map[string]schema.Attribute{
			"uid": schema.StringAttribute{
				Description: "Unique identifier of the index.",
				Required:    true,
			},
			"primary_key": schema.StringAttribute{
				Description: "Primary key of the index (`null` if not specified and if no documents have been added yet, see [official documentation](https://www.meilisearch.com/docs/learn/core_concepts/primary_key#meilisearch-guesses-your-primary-key) for more details).",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Date and time when the key was created (RFC3339)",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Date and time when the key was last updated (RFC3339)",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
		},
	}
}

func (d *indexDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state indexDataSourceModel

	var identifier types.String

	diags := req.Config.GetAttribute(ctx, path.Root("uid"), &identifier)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	index, err := d.client.GetIndex(identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Meilisearch index",
			err.Error(),
		)
		return
	}

	// Map response body to model
	indexState := indexDataSourceModel{
		UID:       types.StringValue(index.UID),
		CreatedAt: types.StringValue(index.CreatedAt.String()),
		UpdatedAt: types.StringValue(index.UpdatedAt.String()),
	}

	if index.PrimaryKey != "" {
		indexState.PrimaryKey = types.StringValue(index.PrimaryKey)
	}

	state = indexState

	state.ID = types.StringValue("placeholder")

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *indexDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var ok bool

	d.client, ok = req.ProviderData.(*meilisearch.Client)

	if !ok {
		tflog.Error(ctx, "Type assertion failed when adding configured client to the data source")
	}
}
