// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// AccumulatorProvider defines the provider implementation.
type AccumulatorProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	store   AccumulatorStore
}

// Configure implements [provider.Provider].
func (p *AccumulatorProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

// DataSources implements [provider.Provider].
func (p *AccumulatorProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

// Schema implements [provider.Provider].
func (p *AccumulatorProvider) Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse) {
}

func (p *AccumulatorProvider) schema() providerschema.Schema {
	return providerschema.Schema{}
}

func (p *AccumulatorProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "accumulator"
}

func (p *AccumulatorProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &GroupResource{store: p.store}
		},
		func() resource.Resource {
			return &ItemResource{store: p.store}
		},
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AccumulatorProvider{
			version: version,
			store:   NewStore(),
		}
	}
}
