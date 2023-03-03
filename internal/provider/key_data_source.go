package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/meilisearch/meilisearch-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &keyDataSource{}
	_ datasource.DataSourceWithConfigure = &keyDataSource{}
)

func NewKeyDataSource() datasource.DataSource {
	return &keyDataSource{}
}

// keyDataSource defines the data source implementation.
type keyDataSource struct {
	client *meilisearch.Client
}

type keyDataSourceModel struct {
	UID         types.String   `tfsdk:"uid"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Key         types.String   `tfsdk:"key"`
	Actions     []types.String `tfsdk:"actions"`
	Indexes     []types.String `tfsdk:"indexes"`
	ExpiresAt   types.String   `tfsdk:"expires_at"`
	CreatedAt   types.String   `tfsdk:"created_at"`
	UpdatedAt   types.String   `tfsdk:"updated_at"`
}

func (d *keyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (d *keyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Meilisearch API key.",
		Attributes: map[string]schema.Attribute{
			"uid": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"key": schema.StringAttribute{
				Computed: true,
			},
			"actions": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"indexes": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"expires_at": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *keyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state keyDataSourceModel

	var identifier types.String

	diags := req.Config.GetAttribute(ctx, path.Root("uid"), &identifier)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	key, err := d.client.GetKey(identifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Meilisearch API key",
			err.Error(),
		)
		return
	}

	// Map response body to model
	keyState := keyDataSourceModel{
		UID:         types.StringValue(key.UID),
		Name:        types.StringValue(key.Name),
		Description: types.StringValue(key.Description),
		Key:         types.StringValue(key.Key),
		ExpiresAt:   types.StringValue(key.ExpiresAt.String()),
		CreatedAt:   types.StringValue(key.CreatedAt.String()),
		UpdatedAt:   types.StringValue(key.UpdatedAt.String()),
	}

	for _, action := range key.Actions {
		keyState.Actions = append(keyState.Actions, types.StringValue(action))
	}

	for _, indexes := range key.Indexes {
		keyState.Indexes = append(keyState.Indexes, types.StringValue(indexes))
	}

	state = keyState

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *keyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*meilisearch.Client)
}
