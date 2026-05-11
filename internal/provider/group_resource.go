package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type GroupResource struct {
	store AccumulatorStore
}

// Create implements [resource.Resource].
func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	panic("unimplemented")
}

// Delete implements [resource.Resource].
func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	panic("unimplemented")
}

// Read implements [resource.Resource].
func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	panic("unimplemented")
}

// Update implements [resource.Resource].
func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	panic("unimplemented")
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Human-readable group name.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Optional value type enforced for all items in the group.",
			},
			"value": schema.DynamicAttribute{
				Computed:    true,
				Description: "Computed value assembled from all items in the group.",
			},
		},
	}
}
