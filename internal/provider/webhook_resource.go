package provider

import (
	"context"
	"fmt"

	"github.com/plain-insure/terraform-provider-printone/internal/client"
	"github.com/plain-insure/terraform-provider-printone/internal/provider/resource_webhook"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = (*webhookResource)(nil)

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

type webhookResource struct {
	client *client.Client
}

func (r *webhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *webhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_webhook.WebhookResourceSchema(ctx)
}

func (r *webhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_webhook.WebhookModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API request
	webhookReq, diags := webhookModelToRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create webhook via API
	webhookResp, err := r.client.CreateWebhook(ctx, webhookReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating webhook",
			"Could not create webhook, unexpected error: "+err.Error(),
		)
		return
	}

	// Convert API response to Terraform model
	resp.Diagnostics.Append(webhookResponseToModel(ctx, webhookResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_webhook.WebhookModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get webhook from API
	webhookResp, err := r.client.GetWebhook(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading webhook",
			"Could not read webhook ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Convert API response to Terraform model
	resp.Diagnostics.Append(webhookResponseToModel(ctx, webhookResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_webhook.WebhookModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API request
	webhookReq, diags := webhookModelToRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update webhook via API
	webhookResp, err := r.client.UpdateWebhook(ctx, data.Id.ValueString(), webhookReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating webhook",
			"Could not update webhook ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Convert API response to Terraform model
	resp.Diagnostics.Append(webhookResponseToModel(ctx, webhookResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_webhook.WebhookModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete webhook via API
	err := r.client.DeleteWebhook(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting webhook",
			"Could not delete webhook ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}
}
