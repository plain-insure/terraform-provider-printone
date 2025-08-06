package provider

import (
	"context"
    "github.com/plain-insure/terraform-provider-printone/internal/provider/datasource_webhook"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/diag"
)

var _ datasource.DataSource = (*webhookDataSource)(nil)

func NewWebhookDataSource() datasource.DataSource {
	return &webhookDataSource{}
}

type webhookDataSource struct{}


func (d *webhookDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (d *webhookDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_webhook.WebhookDataSourceSchema(ctx)
}

func (d *webhookDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_webhook.WebhookModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
    resp.Diagnostics.Append(callWebhookAPI(ctx, &data)...)

	// Example data value setting
	data.Id = types.StringValue("example-id")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Typically this method would contain logic that makes an HTTP call to a remote API, and then stores
// computed results back to the data model. For example purposes, this function just sets computed Order
// values to mock values to avoid data consistency errors.
func callWebhookAPI(ctx context.Context, webhook *datasource_webhook.WebhookModel) diag.Diagnostics {
    webhook.Id = types.StringValue("1")
    webhook.Name = types.StringValue("active")
    webhook.Url = types.StringValue("https://example.com/webhook")

    return nil
}
