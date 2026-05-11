package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestNewInitializesProviderWithStore(t *testing.T) {
	providerFactory := New("test")
	instance := providerFactory()

	accumulatorProvider, ok := instance.(*AccumulatorProvider)
	if !ok {
		t.Fatalf("expected *AccumulatorProvider, got %T", instance)
	}

	if accumulatorProvider.store == nil {
		t.Fatal("expected provider store to be initialized")
	}
}

func TestProviderMetadataTypeNameIsAccumulator(t *testing.T) {
	instance := New("test")()
	accumulatorProvider, ok := instance.(*AccumulatorProvider)
	if !ok {
		t.Fatalf("expected *AccumulatorProvider, got %T", instance)
	}

	var resp provider.MetadataResponse
	accumulatorProvider.Metadata(context.Background(), provider.MetadataRequest{}, &resp)

	if resp.TypeName != "accumulator" {
		t.Fatalf("expected type name accumulator, got %q", resp.TypeName)
	}
}

func TestProviderResourcesReturnGroupAndItem(t *testing.T) {
	instance := New("test")()
	accumulatorProvider, ok := instance.(*AccumulatorProvider)
	if !ok {
		t.Fatalf("expected *AccumulatorProvider, got %T", instance)
	}

	resources := accumulatorProvider.Resources(context.Background())
	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}

	gotTypes := []resource.Resource{
		resources[0](),
		resources[1](),
	}

	if _, ok := gotTypes[0].(*GroupResource); !ok {
		t.Fatalf("expected first resource to be *GroupResource, got %T", gotTypes[0])
	}

	if _, ok := gotTypes[1].(*ItemResource); !ok {
		t.Fatalf("expected second resource to be *ItemResource, got %T", gotTypes[1])
	}

	group := gotTypes[0].(*GroupResource)
	if group.store == nil {
		t.Fatal("expected group resource to receive provider store")
	}

	item := gotTypes[1].(*ItemResource)
	if item.store == nil {
		t.Fatal("expected item resource to receive provider store")
	}
}
