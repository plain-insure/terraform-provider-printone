// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/plain-insure/terraform-provider-printone/internal/client"
	"github.com/plain-insure/terraform-provider-printone/internal/provider/datasource_webhook"
	"github.com/plain-insure/terraform-provider-printone/internal/provider/resource_webhook"
)

// webhookModelToRequest converts a Terraform webhook model to an API request
func webhookModelToRequest(ctx context.Context, model *resource_webhook.WebhookModel) (*client.WebhookRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert events list
	var events []string
	diags.Append(model.Events.ElementsAs(ctx, &events, false)...)

	// Convert headers if not null
	var headers map[string]interface{}
	if !model.Headers.IsNull() && !model.Headers.IsUnknown() {
		headers = make(map[string]interface{})
		// For now, we'll leave headers empty - the generated code has complex nested types
		// This can be enhanced later when needed
	}

	// Convert secret headers if not null
	var secretHeaders map[string]interface{}
	if !model.SecretHeaders.IsNull() && !model.SecretHeaders.IsUnknown() {
		secretHeaders = make(map[string]interface{})
		// For now, we'll leave secret headers empty - similar to headers
	}

	request := &client.WebhookRequest{
		Name:          model.Name.ValueString(),
		URL:           model.Url.ValueString(),
		Active:        model.Active.ValueBool(),
		Events:        events,
		Headers:       headers,
		SecretHeaders: secretHeaders,
	}

	return request, diags
}

// webhookResponseToModel converts an API response to a Terraform webhook resource model
func webhookResponseToModel(ctx context.Context, response *client.WebhookResponse, model *resource_webhook.WebhookModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.Id = types.StringValue(response.ID)
	model.Name = types.StringValue(response.Name)
	model.Url = types.StringValue(response.URL)
	model.Active = types.BoolValue(response.Active)

	// Convert events
	eventsList, d := types.ListValueFrom(ctx, types.StringType, response.Events)
	diags.Append(d...)
	model.Events = eventsList

	// For now, set complex nested types to null/unknown
	// These can be enhanced later when the complete mapping is needed
	model.Headers = resource_webhook.NewHeadersValueNull()
	model.SecretHeaders = resource_webhook.NewSecretHeadersValueNull()
	model.SuccessRate = resource_webhook.NewSuccessRateValueNull()

	return diags
}

// webhookResponseToDataSourceModel converts an API response to a Terraform webhook data source model
func webhookResponseToDataSourceModel(ctx context.Context, response *client.WebhookResponse, model *datasource_webhook.WebhookModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.Id = types.StringValue(response.ID)
	model.Name = types.StringValue(response.Name)
	model.Url = types.StringValue(response.URL)
	model.Active = types.BoolValue(response.Active)

	// Convert events
	eventsList, d := types.ListValueFrom(ctx, types.StringType, response.Events)
	diags.Append(d...)
	model.Events = eventsList

	// For now, set complex nested types to null
	// These can be enhanced later when the complete mapping is needed
	model.Headers = datasource_webhook.NewHeadersValueNull()
	model.SecretHeaders = datasource_webhook.NewSecretHeadersValueNull()
	model.SuccessRate = datasource_webhook.NewSuccessRateValueNull()

	return diags
}
