// Copyright (c) Plain Technologies Aps

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/plain-insure/terraform-provider-printone/internal/client"
)

var _ provider.Provider = (*printoneProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &printoneProvider{
			version: version,
		}
	}
}

type printoneProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// printoneProviderModel describes the provider data model.
type printoneProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}

func (p *printoneProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The PrintOne provider is used to interact with PrintOne API resources.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The PrintOne API endpoint URL. Defaults to https://api.print.one if not provided.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The PrintOne API key for authentication. Can also be set via the PRINTONE_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *printoneProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data printoneProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create PrintOne client configuration
	config := client.Config{}

	// Set base URL if provided
	if !data.Endpoint.IsNull() && !data.Endpoint.IsUnknown() {
		config.BaseURL = data.Endpoint.ValueString()
	}

	// Set API key if provided
	if !data.APIKey.IsNull() && !data.APIKey.IsUnknown() {
		config.APIKey = data.APIKey.ValueString()
	}

	// Create the PrintOne client
	printoneclient := client.NewClient(config)

	// Make the client available to resources and data sources
	resp.DataSourceData = printoneclient
	resp.ResourceData = printoneclient
}

func (p *printoneProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "printone"
	resp.Version = p.version
}

func (p *printoneProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWebhookDataSource,
	}
}

func (p *printoneProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWebhookResource,
	}
}
