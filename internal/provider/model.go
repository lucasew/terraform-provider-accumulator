package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type groupModel struct {
	ID    types.String  `tfsdk:"id"`
	Value types.Dynamic `tfsdk:"value"`
}

type itemModel struct {
	ID    types.String  `tfsdk:"id"`
	Group types.String  `tfsdk:"group"`
	Key   types.String  `tfsdk:"key"`
	Value types.Dynamic `tfsdk:"value"`
}
