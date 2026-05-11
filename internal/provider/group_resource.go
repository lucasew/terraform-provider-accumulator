package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupResource struct {
	store AccumulatorStore
}

// Create implements [resource.Resource].
func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	groupID, err := r.store.CreateGroup()
	if err != nil {
		resp.Diagnostics.AddError("creating group failed", err.Error())
		return
	}

	state := groupModel{
		ID:    types.StringValue(groupID.String()),
		Value: dynamicObjectFromGo(map[string]any{}),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete implements [resource.Resource].
func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	if err := r.store.DeleteGroup(groupID); err != nil {
		resp.Diagnostics.AddError("deleting group failed", err.Error())
	}
}

// Read implements [resource.Resource].
func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	data, err := r.store.GroupData(groupID)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Value = dynamicObjectFromGo(data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update implements [resource.Resource].
func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	data, err := r.store.GroupData(groupID)
	if err != nil {
		resp.Diagnostics.AddError("updating group failed", err.Error())
		return
	}

	state.Value = dynamicObjectFromGo(data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Stable group identifier.",
			},
			"value": schema.DynamicAttribute{
				Computed:    true,
				Description: "Computed value assembled from all items in the group.",
			},
		},
	}
}
