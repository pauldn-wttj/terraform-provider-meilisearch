package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/meilisearch/meilisearch-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &versionDataSource{}
	_ datasource.DataSourceWithConfigure = &versionDataSource{}
)

func NewVersionDataSource() datasource.DataSource {
	return &versionDataSource{}
}

// versionDataSource defines the data source implementation.
type versionDataSource struct {
	client meilisearch.ServiceManager
}

type versionDataSourceModel struct {
	CommitSha  types.String `tfsdk:"commit_sha"`
	CommitDate types.String `tfsdk:"commit_date"`
	PkgVersion types.String `tfsdk:"pkg_version"`
	ID         types.String `tfsdk:"id"`
}

func (d *versionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_version"
}

func (d *versionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the meilisearch version.",
		Attributes: map[string]schema.Attribute{
			"commit_sha": schema.StringAttribute{
				Description: "Commit identifier that tagged the pkgVersion release",
				Computed:    true,
			},
			"commit_date": schema.StringAttribute{
				Description: "Date when the commitSha was created",
				Computed:    true,
			},
			"pkg_version": schema.StringAttribute{
				Description: "Meilisearch version",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
		},
	}
}

func (d *versionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state versionDataSourceModel

	version, err := d.client.Version()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Meilisearch Version",
			err.Error(),
		)
		return
	}

	// Map response body to model
	versionState := versionDataSourceModel{
		CommitSha:  types.StringValue(version.CommitSha),
		CommitDate: types.StringValue(version.CommitDate),
		PkgVersion: types.StringValue(version.PkgVersion),
	}

	state = versionState

	state.ID = types.StringValue("placeholder")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *versionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var ok bool

	d.client, ok = req.ProviderData.(meilisearch.ServiceManager)

	if !ok {
		tflog.Error(ctx, "Type assertion failed when adding configured client to the data source")
	}
}
