package provider

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/meilisearch/meilisearch-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &keyResource{}
	_ resource.ResourceWithConfigure   = &keyResource{}
	_ resource.ResourceWithImportState = &keyResource{}
)

// NewKeyResource is a helper function to simplify the provider implementation.
func NewKeyResource() resource.Resource {
	return &keyResource{}
}

// keyResource is the resource implementation.
type keyResource struct {
	client *meilisearch.Client
}

type keyResourceModel struct {
	UID         types.String   `tfsdk:"uid"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Key         types.String   `tfsdk:"key"`
	Actions     []types.String `tfsdk:"actions"`
	Indexes     []types.String `tfsdk:"indexes"`
	ExpiresAt   types.String   `tfsdk:"expires_at"`
	CreatedAt   types.String   `tfsdk:"created_at"`
	UpdatedAt   types.String   `tfsdk:"updated_at"`
	ID          types.String   `tfsdk:"id"`
}

// Metadata returns the resource type name.
func (r *keyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

// Schema defines the schema for the resource.
func (r *keyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Meilisearch API key.",
		Attributes: map[string]schema.Attribute{
			"uid": schema.StringAttribute{
				Description: "UID (uuid v4) used by Meilisearch to identify the key.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the key.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the key.",
				Optional:    true,
			},
			"key": schema.StringAttribute{
				Description: "Actual key value.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"actions": schema.ListAttribute{
				Description: "Actions permitted for the key.",
				ElementType: types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"indexes": schema.ListAttribute{
				Description: "Indexes the key is authorized to act on (with the actions specified in the scope of the key).",
				ElementType: types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"expires_at": schema.StringAttribute{
				Description: "Date and time when the key will expire (RFC3339)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Date and time when the key was created (RFC3339)",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

// Configure adds the provider configured client to the resource.
func (r *keyResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var ok bool

	r.client, ok = req.ProviderData.(*meilisearch.Client)

	if !ok {
		tflog.Error(ctx, "Type assertion failed when adding configured client to the resource")
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *keyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan keyResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var actions []string
	var indexes []string
	var expiresAt time.Time

	for _, action := range plan.Actions {
		actions = append(actions, action.ValueString())
	}

	for _, index := range plan.Indexes {
		indexes = append(indexes, index.ValueString())
	}

	if !plan.ExpiresAt.IsNull() && plan.ExpiresAt.ValueString() != "" {
		parsedExpiredAt, err := time.Parse(time.RFC3339, plan.ExpiresAt.ValueString())

		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating key",
				"Could not parse expiresAt attribute",
			)
			return
		}

		expiresAt = parsedExpiredAt
	}

	createKey := meilisearch.Key{
		UID:         plan.UID.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Actions:     actions,
		Indexes:     indexes,
		ExpiresAt:   expiresAt,
	}

	key, err := r.client.CreateKey(&createKey)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating key",
			"Could not create key, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UID = types.StringValue(key.UID)
	plan.Key = types.StringValue(key.Key)
	plan.CreatedAt = types.StringValue(key.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(key.UpdatedAt.Format(time.RFC3339))

	if plan.ExpiresAt.IsNull() {
		plan.ExpiresAt = types.StringNull()
	}

	plan.ID = types.StringValue("placeholder")

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *keyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state keyResourceModel

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed key value from Meilisearch
	key, err := r.client.GetKey(state.UID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "api_key_not_found,") {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error Reading Meilisearch Key",
				"Could not read Meilisearch key ID "+state.UID.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	// Overwrite items with refreshed state
	keyState := keyResourceModel{
		UID:         types.StringValue(key.UID),
		Name:        types.StringValue(key.Name),
		Description: types.StringValue(key.Description),
		Key:         types.StringValue(key.Key),
		CreatedAt:   types.StringValue(key.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:   types.StringValue(key.UpdatedAt.Format(time.RFC3339)),
	}

	for _, action := range key.Actions {
		keyState.Actions = append(keyState.Actions, types.StringValue(action))
	}

	for _, indexes := range key.Indexes {
		keyState.Indexes = append(keyState.Indexes, types.StringValue(indexes))
	}

	if key.ExpiresAt.IsZero() {
		keyState.ExpiresAt = types.StringNull()
	} else {
		keyState.ExpiresAt = types.StringValue(key.ExpiresAt.Format(time.RFC3339))
	}

	state = keyState

	state.ID = types.StringValue("placeholder")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *keyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan keyResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateKey := meilisearch.Key{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	// Update existing key
	key, err := r.client.UpdateKey(plan.UID.ValueString(), &updateKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Meilisearch Key",
			"Could not update key, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UID = types.StringValue(key.UID)
	plan.Key = types.StringValue(key.Key)
	plan.CreatedAt = types.StringValue(key.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(key.UpdatedAt.Format(time.RFC3339))

	if plan.ExpiresAt.IsNull() {
		plan.ExpiresAt = types.StringNull()
	}

	plan.ID = types.StringValue("placeholder")

	// Set refreshed state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *keyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state keyResourceModel

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing key
	_, err := r.client.DeleteKey(state.UID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Meilisearch Key",
			"Could not delete key, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *keyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import UID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uid"), req, resp)
}
