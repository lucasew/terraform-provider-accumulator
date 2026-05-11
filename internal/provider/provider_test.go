package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestNewInitializesProvider(t *testing.T) {
	providerFactory := New("test")
	instance := providerFactory()

	accumulatorProvider, ok := instance.(*AccumulatorProvider)
	if !ok {
		t.Fatalf("expected *AccumulatorProvider, got %T", instance)
	}

	if accumulatorProvider.store == nil {
		t.Fatal("expected provider instance to initialize its own store")
	}
}

func TestProviderConfigureExposesResourceData(t *testing.T) {
	instance := New("test")()
	accumulatorProvider, ok := instance.(*AccumulatorProvider)
	if !ok {
		t.Fatalf("expected *AccumulatorProvider, got %T", instance)
	}

	var resp provider.ConfigureResponse
	accumulatorProvider.Configure(context.Background(), provider.ConfigureRequest{}, &resp)

	data, ok := resp.ResourceData.(*providerData)
	if !ok {
		t.Fatalf("expected *providerData, got %T", resp.ResourceData)
	}

	if data.store == nil {
		t.Fatal("expected configured provider data to include a store")
	}

	if data.store != accumulatorProvider.store {
		t.Fatal("expected configured provider data to expose the provider instance store")
	}

	dsData, ok := resp.DataSourceData.(*providerData)
	if !ok {
		t.Fatalf("expected *providerData in DataSourceData, got %T", resp.DataSourceData)
	}

	if dsData.store != accumulatorProvider.store {
		t.Fatal("expected configured data source data to expose the provider instance store")
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

func TestProviderResourcesReturnItem(t *testing.T) {
	instance := New("test")()
	accumulatorProvider, ok := instance.(*AccumulatorProvider)
	if !ok {
		t.Fatalf("expected *AccumulatorProvider, got %T", instance)
	}

	resources := accumulatorProvider.Resources(context.Background())
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}

	gotTypes := []resource.Resource{
		resources[0](),
	}

	if _, ok := gotTypes[0].(*ItemResource); !ok {
		t.Fatalf("expected resource to be *ItemResource, got %T", gotTypes[0])
	}

	item := gotTypes[0].(*ItemResource)
	if item.store != nil {
		t.Fatal("expected item resource store to be injected during Configure")
	}
}

func TestProviderDataSourcesReturnGroup(t *testing.T) {
	instance := New("test")()
	accumulatorProvider, ok := instance.(*AccumulatorProvider)
	if !ok {
		t.Fatalf("expected *AccumulatorProvider, got %T", instance)
	}

	dataSources := accumulatorProvider.DataSources(context.Background())
	if len(dataSources) != 1 {
		t.Fatalf("expected 1 data source, got %d", len(dataSources))
	}

	got := dataSources[0]()
	group, ok := got.(*GroupDataSource)
	if !ok {
		t.Fatalf("expected data source to be *GroupDataSource, got %T", got)
	}

	if group.store != nil {
		t.Fatal("expected group data source store to be injected during Configure")
	}

	_ = datasource.DataSource(group)
}
