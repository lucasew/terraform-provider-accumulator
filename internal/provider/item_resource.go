package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type ItemResource struct {
	store AccumulatorStore
}

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

// Update implements [resource.Resource].
func (r *ItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	panic("unimplemented")
}

func (r *ItemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_item"
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
