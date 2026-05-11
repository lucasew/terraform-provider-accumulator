package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestGroupResourceSchemaExposesAccumulatorFields(t *testing.T) {
	group := &GroupResource{}
	var resp resource.SchemaResponse

	group.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	required := []string{"id", "name", "type", "value"}
	for _, name := range required {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Fatalf("expected group schema to include %q", name)
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
