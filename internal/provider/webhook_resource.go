// Copyright (c) Plain Technologies Aps

package provider

import (
	"context"
	"fmt"
	"slices"

	"github.com/plain-insure/terraform-provider-printone/internal/client"
	"github.com/plain-insure/terraform-provider-printone/internal/provider/resource_webhook"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = (*webhookResource)(nil)

// Supported event types
var supportedEvents = []string{
	"order_status_update",
	"template_preview_rendered",
	"batch_status_update",
	"coupon_code_used",
	"incasso_failed",
	"incasso_reversed",
	"qr_code_scanned",
	"company_onboarding_changed",
	"company_signup",
	"company_neared_postpaid_limit",
}

// Supported filter types
var supportedFilterTypes = []string{
	"equals",
	"not-equals",
	"in",
	"not-in",
}

// validateFilters validates webhook filters according to the business rules
func validateFilters(ctx context.Context, filters []resource_webhook.WebhookFilterModel) diag.Diagnostics {
	var diags diag.Diagnostics

	for i, filter := range filters {
		// Validate event type
		if !slices.Contains(supportedEvents, filter.Event.ValueString()) {
			diags.AddError(
				"Invalid Filter Event",
				fmt.Sprintf("Filter %d has invalid event '%s'. Supported events: %v", i, filter.Event.ValueString(), supportedEvents),
			)
		}

		// Validate filter type
		filterType := filter.Type.ValueString()
		if !slices.Contains(supportedFilterTypes, filterType) {
			diags.AddError(
				"Invalid Filter Type",
				fmt.Sprintf("Filter %d has invalid type '%s'. Supported types: %v", i, filterType, supportedFilterTypes),
			)
			continue
		}

		// Validate value/values constraints based on filter type
		hasValue := !filter.Value.IsNull() && !filter.Value.IsUnknown() && filter.Value.ValueString() != ""
		hasValues := !filter.Values.IsNull() && !filter.Values.IsUnknown() && len(filter.Values.Elements()) > 0

		switch filterType {
		case "equals", "not-equals":
			if !hasValue {
				diags.AddError(
					"Missing Filter Value",
					fmt.Sprintf("Filter %d with type '%s' must have a 'value' set", i, filterType),
				)
			}
			if hasValues {
				diags.AddError(
					"Invalid Filter Values",
					fmt.Sprintf("Filter %d with type '%s' must NOT have 'values' set, only 'value'", i, filterType),
				)
			}
		case "in", "not-in":
			if !hasValues {
				diags.AddError(
					"Missing Filter Values",
					fmt.Sprintf("Filter %d with type '%s' must have 'values' set", i, filterType),
				)
			}
			if hasValue {
				diags.AddError(
					"Invalid Filter Value",
					fmt.Sprintf("Filter %d with type '%s' must NOT have 'value' set, only 'values'", i, filterType),
				)
			}
		}
	}

	return diags
}

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

	// Read Terraform plan data into the model.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate filters if present
	if !data.Filters.IsNull() && !data.Filters.IsUnknown() {
		var filterModels []resource_webhook.WebhookFilterModel
		resp.Diagnostics.Append(data.Filters.ElementsAs(ctx, &filterModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.Append(validateFilters(ctx, filterModels)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Convert Terraform model to API request.
	webhookReq, diags := webhookModelToRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create webhook via API.
	webhookResp, err := r.client.CreateWebhook(ctx, webhookReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating webhook",
			"Could not create webhook, unexpected error: "+err.Error(),
		)
		return
	}

	// Convert API response to Terraform model.
	resp.Diagnostics.Append(webhookResponseToModel(ctx, webhookResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For secret_headers, always set exactly from the plan (no transformation).
	var planInput resource_webhook.WebhookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planInput)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.SecretHeaders = planInput.SecretHeaders

	// Save data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_webhook.WebhookModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get webhook from API.
	webhookResp, err := r.client.GetWebhook(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading webhook",
			"Could not read webhook ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Convert API response to Terraform model.
	resp.Diagnostics.Append(webhookResponseToModel(ctx, webhookResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Always set secret_headers from previous state, never overwrite with API response, and copy exactly.
	var prevState resource_webhook.WebhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prevState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.SecretHeaders = prevState.SecretHeaders

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_webhook.WebhookModel
	var state resource_webhook.WebhookModel

	// Read Terraform plan data into the model.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the current state to get the ID.
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

       // Validate filters if present
       if !plan.Filters.IsNull() && !plan.Filters.IsUnknown() {
	       var filterModels []resource_webhook.WebhookFilterModel
	       resp.Diagnostics.Append(plan.Filters.ElementsAs(ctx, &filterModels, false)...)
	       if resp.Diagnostics.HasError() {
		       return
	       }
	       resp.Diagnostics.Append(validateFilters(ctx, filterModels)...)
	       if resp.Diagnostics.HasError() {
		       return
	       }
       }

	// Convert Terraform model to API request.
	webhookReq, diags := webhookModelToRequest(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update webhook via API using the ID from state.
	webhookResp, err := r.client.UpdateWebhook(ctx, state.Id.ValueString(), webhookReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating webhook",
			"Could not update webhook ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Convert API response to Terraform model.
	resp.Diagnostics.Append(webhookResponseToModel(ctx, webhookResp, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Always set secret_headers from the plan, never from the API response, and copy exactly.
	var planInput resource_webhook.WebhookModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planInput)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.SecretHeaders = planInput.SecretHeaders

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_webhook.WebhookModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete webhook via API.
	err := r.client.DeleteWebhook(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting webhook",
			"Could not delete webhook ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}
}
