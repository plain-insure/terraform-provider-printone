// Copyright (c) Plain Technologies Aps

package provider

import (
	"context"
	"fmt"

	"github.com/plain-insure/terraform-provider-printone/internal/client"
	"github.com/plain-insure/terraform-provider-printone/internal/provider/datasource_webhooksecret"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*webhooksecretDataSource)(nil)

func NewWebhooksecretDataSource() datasource.DataSource {
	return &webhooksecretDataSource{}
}

type webhooksecretDataSource struct {
	client *client.Client
}

func (d *webhooksecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooksecret"
}

func (d *webhooksecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_webhooksecret.WebhooksecretDataSourceSchema(ctx)
}

func (d *webhooksecretDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *webhooksecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_webhooksecret.WebhooksecretModel

	// Read Terraform configuration data into the model.
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get webhook secret from API.
	secretResp, err := d.client.GetWebhookSecret(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading webhook secret",
			"Could not read webhook secret: "+err.Error(),
		)
		return
	}

	// Convert API response to Terraform model.
	data.Secret = types.StringValue(secretResp.Secret)

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
