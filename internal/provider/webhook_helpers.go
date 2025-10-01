// Copyright (c) Plain Technologies Aps

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/plain-insure/terraform-provider-printone/internal/client"
	"github.com/plain-insure/terraform-provider-printone/internal/provider/datasource_webhook"
	"github.com/plain-insure/terraform-provider-printone/internal/provider/resource_webhook"
)

// webhookModelToRequest converts a Terraform webhook model to an API request.
func webhookModelToRequest(ctx context.Context, model *resource_webhook.WebhookModel) (*client.WebhookRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert events list.
	var events []string
	diags.Append(model.Events.ElementsAs(ctx, &events, false)...)

	// Convert headers if not null.
	var headers map[string]interface{}
	if !model.Headers.IsNull() && !model.Headers.IsUnknown() {
		headers = make(map[string]interface{})
		// For now, we'll leave headers empty - the generated code has complex nested types.
		// This can be enhanced later when needed.
	}

	// Convert secret headers if not null.
	var secretHeaders map[string]interface{}
	if !model.SecretHeaders.IsNull() && !model.SecretHeaders.IsUnknown() {
		secretHeaders = make(map[string]interface{})
		// For now, we'll leave secret headers empty - similar to headers.
	}

	// Convert filters if not null.
	var filters []client.WebhookFilter
	if !model.Filters.IsNull() && !model.Filters.IsUnknown() {
		var filterModels []resource_webhook.WebhookFilterModel
		diags.Append(model.Filters.ElementsAs(ctx, &filterModels, false)...)

		for _, filterModel := range filterModels {
			filter := client.WebhookFilter{
				Key:   filterModel.Key.ValueString(),
				Event: filterModel.Event.ValueString(),
				Type:  filterModel.Type.ValueString(),
			}

			// Set value or values based on filter type.
			if !filterModel.Value.IsNull() && !filterModel.Value.IsUnknown() {
				value := filterModel.Value.ValueString()
				filter.Value = &value
			}

			if !filterModel.Values.IsNull() && !filterModel.Values.IsUnknown() {
				var values []string
				diags.Append(filterModel.Values.ElementsAs(ctx, &values, false)...)
				filter.Values = values
			}

			filters = append(filters, filter)
		}
	}

	request := &client.WebhookRequest{
		Name:          model.Name.ValueString(),
		URL:           model.Url.ValueString(),
		Active:        model.Active.ValueBool(),
		Events:        events,
		Headers:       headers,
		SecretHeaders: secretHeaders,
		Filters:       filters,
	}

	return request, diags
}

// webhookResponseToModel converts an API response to a Terraform webhook resource model.
func webhookResponseToModel(ctx context.Context, response *client.WebhookResponse, model *resource_webhook.WebhookModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.Id = types.StringValue(response.ID)
	model.Name = types.StringValue(response.Name)
	model.Url = types.StringValue(response.URL)
	model.Active = types.BoolValue(response.Active)

	// Convert events.
	eventsList, d := types.ListValueFrom(ctx, types.StringType, response.Events)
	diags.Append(d...)
	model.Events = eventsList

	// Convert filters.
	if len(response.Filters) > 0 {
		filterModels := make([]resource_webhook.WebhookFilterModel, len(response.Filters))
		for i, filter := range response.Filters {
			filterModel := resource_webhook.WebhookFilterModel{
				Key:   types.StringValue(filter.Key),
				Event: types.StringValue(filter.Event),
				Type:  types.StringValue(filter.Type),
			}

			if filter.Value != nil {
				filterModel.Value = types.StringValue(*filter.Value)
			} else {
				filterModel.Value = types.StringNull()
			}

			if len(filter.Values) > 0 {
				valuesList, d := types.ListValueFrom(ctx, types.StringType, filter.Values)
				diags.Append(d...)
				filterModel.Values = valuesList
			} else {
				filterModel.Values = types.ListNull(types.StringType)
			}

			filterModels[i] = filterModel
		}

		filtersList, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":    types.StringType,
				"event":  types.StringType,
				"type":   types.StringType,
				"value":  types.StringType,
				"values": types.ListType{ElemType: types.StringType},
			},
		}, filterModels)
		diags.Append(d...)
		model.Filters = filtersList
	} else {
		model.Filters = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":    types.StringType,
				"event":  types.StringType,
				"type":   types.StringType,
				"value":  types.StringType,
				"values": types.ListType{ElemType: types.StringType},
			},
		})
	}

	// For now, set complex nested types to null/unknown.
	// These can be enhanced later when the complete mapping is needed.
	model.Headers = resource_webhook.NewHeadersValueNull()
	model.SecretHeaders = resource_webhook.NewSecretHeadersValueNull()
	model.SuccessRate = resource_webhook.NewSuccessRateValueNull()

	return diags
}

// webhookResponseToDataSourceModel converts an API response to a Terraform webhook data source model.
func webhookResponseToDataSourceModel(ctx context.Context, response *client.WebhookResponse, model *datasource_webhook.WebhookModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.Id = types.StringValue(response.ID)
	model.Name = types.StringValue(response.Name)
	model.Url = types.StringValue(response.URL)
	model.Active = types.BoolValue(response.Active)

	// Convert events.
	eventsList, d := types.ListValueFrom(ctx, types.StringType, response.Events)
	diags.Append(d...)
	model.Events = eventsList

	// For now, set complex nested types to null.
	// These can be enhanced later when the complete mapping is needed.
	model.Headers = datasource_webhook.NewHeadersValueNull()
	model.SecretHeaders = datasource_webhook.NewSecretHeadersValueNull()
	model.SuccessRate = datasource_webhook.NewSuccessRateValueNull()

	return diags
}
