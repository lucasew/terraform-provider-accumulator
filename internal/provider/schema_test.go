package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestGroupDataSourceSchemaExposesAccumulatorFields(t *testing.T) {
	group := &GroupDataSource{}
	var resp datasource.SchemaResponse

	group.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	required := []string{"id", "value"}
	for _, name := range required {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("expected group data source schema to include %q", name)
		}
	}
}

func TestItemResourceSchemaExposesAccumulatorFields(t *testing.T) {
	item := &ItemResource{}
	var resp resource.SchemaResponse

	item.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	required := []string{"id", "group", "key", "value"}
	for _, name := range required {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("expected item schema to include %q", name)
		}
	}
}
