package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type groupModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	Data types.Map    `tfsdk:"data"`
}

type itemModel struct {
	ID    types.String  `tfsdk:"id"`
	Group types.String  `tfsdk:"group"`
	Key   types.String  `tfsdk:"key"`
	Value types.Dynamic `tfsdk:"value"`
}
