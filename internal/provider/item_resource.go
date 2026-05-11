package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ItemResource struct {
	store AccumulatorStore
}

// Create implements [resource.Resource].
func (r *ItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan itemModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(plan.Group.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	value, err := goValueFromDynamic(plan.Value)
	if err != nil {
		resp.Diagnostics.AddError("invalid item value", err.Error())
		return
	}

	if err := r.store.PutItem(groupID, plan.Key.ValueString(), value); err != nil {
		resp.Diagnostics.AddError("creating item failed", err.Error())
		return
	}

	plan.ID = types.StringValue(itemID(groupID, plan.Key.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements [resource.Resource].
func (r *ItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state itemModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(state.Group.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	if err := r.store.DeleteItem(groupID, state.Key.ValueString()); err != nil {
		resp.Diagnostics.AddError("deleting item failed", err.Error())
	}
}

// Read implements [resource.Resource].
func (r *ItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state itemModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID, err := uuid.Parse(state.Group.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	data, err := r.store.GroupData(groupID)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if _, ok := data[state.Key.ValueString()]; !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update implements [resource.Resource].
func (r *ItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var priorState itemModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan itemModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	priorGroupID, err := uuid.Parse(priorState.Group.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid prior group id", err.Error())
		return
	}

	if err := r.store.DeleteItem(priorGroupID, priorState.Key.ValueString()); err != nil {
		resp.Diagnostics.AddError("updating item failed", err.Error())
		return
	}

	groupID, err := uuid.Parse(plan.Group.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid group id", err.Error())
		return
	}

	value, err := goValueFromDynamic(plan.Value)
	if err != nil {
		resp.Diagnostics.AddError("invalid item value", err.Error())
		return
	}

	if err := r.store.PutItem(groupID, plan.Key.ValueString(), value); err != nil {
		resp.Diagnostics.AddError("updating item failed", err.Error())
		return
	}

	plan.ID = types.StringValue(itemID(groupID, plan.Key.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ItemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_item"
}

func itemID(groupID uuid.UUID, key string) string {
	return fmt.Sprintf("%s:%s", groupID.String(), key)
}

func (r *ItemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Stable item identifier.",
			},
			"group": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the target group.",
			},
			"key": schema.StringAttribute{
				Required:    true,
				Description: "Map key contributed by this item.",
			},
			"value": schema.DynamicAttribute{
				Required:    true,
				Description: "Value contributed by this item.",
			},
		},
	}
}
