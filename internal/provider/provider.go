package provider

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// AccumulatorProvider defines the provider implementation.
type AccumulatorProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	store   AccumulatorStore
	id      uint64
}

type providerData struct {
	store AccumulatorStore
}

var nextProviderID atomic.Uint64

// Configure implements [provider.Provider].
func (p *AccumulatorProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	slog.Debug("provider.Configure", "provider_id", p.id, "store_type", fmt.Sprintf("%T", p.store))
	resp.DataSourceData = &providerData{
		store: p.store,
	}
	resp.ResourceData = &providerData{
		store: p.store,
	}
}

// DataSources implements [provider.Provider].
func (p *AccumulatorProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &GroupDataSource{}
		},
	}
}

// Schema implements [provider.Provider].
func (p *AccumulatorProvider) Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse) {
}

func (p *AccumulatorProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "accumulator"
}

func (p *AccumulatorProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &ItemResource{}
		},
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		providerID := nextProviderID.Add(1)
		return &AccumulatorProvider{
			version: version,
			store:   NewLoggingStore(fmt.Sprintf("provider-%d", providerID), NewStore()),
			id:      providerID,
		}
	}
}
