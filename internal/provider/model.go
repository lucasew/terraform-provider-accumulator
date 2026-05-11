package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var Values = map[uuid.UUID]map[string]any{}

type GroupResource struct{}

// Create implements [resource.Resource].
func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var groupID uuid.UUID
	if resp.State.Get(ctx, &groupID).HasError() {
		groupID = uuid.New()
		resp.State.Set(ctx, uuid.New())
	}
	Values[groupID] = map[string]any{}
}

// Delete implements [resource.Resource].
func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var groupID uuid.UUID
	if resp.State.Get(ctx, &groupID).HasError() {
		return
	}

	delete(Values, groupID)
}

// Read implements [resource.Resource].
func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var groupID uuid.UUID
	if resp.State.Get(ctx, &groupID).HasError() {
		return
	}

	if val, ok := Values[groupID]; ok {
		// Set the state for the read operation
		resp.State.Set(ctx, val)
	}
}

// Update implements [resource.Resource].
func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *GroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the group.",
			},
		},
	}
}

type ItemResource struct{}

// Create implements [resource.Resource].
func (r *ItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	panic("unimplemented")
}

// Delete implements [resource.Resource].
func (r *ItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	panic("unimplemented")
}

// Read implements [resource.Resource].
func (r *ItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	panic("unimplemented")
}

// Schema implements [resource.Resource].
func (r *ItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the item.",
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the group this item belongs to.",
			},
		},
	}
}

// Update implements [resource.Resource].
func (r *ItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	panic("unimplemented")
}

func (r *ItemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_item"
}
