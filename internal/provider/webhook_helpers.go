// Copyright (c) Plain Technologies Aps

package provider

import (
	"context"
	"math/big"

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
		headersMap := make(map[string]string)
		diags.Append(model.Headers.ElementsAs(ctx, &headersMap, false)...)

		// Convert string map to interface{} map
		headers = make(map[string]interface{})
		for k, v := range headersMap {
			headers[k] = v
		}
	}

	// Convert secret headers if not null.
	var secretHeaders map[string]interface{}
	if !model.SecretHeaders.IsNull() && !model.SecretHeaders.IsUnknown() {
		secretHeadersMap := make(map[string]string)
		diags.Append(model.SecretHeaders.ElementsAs(ctx, &secretHeadersMap, false)...)

		// Convert string map to interface{} map
		secretHeaders = make(map[string]interface{})
		for k, v := range secretHeadersMap {
			secretHeaders[k] = v
		}
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

	// Convert headers.
	if response.Headers != nil {
		// Convert interface{} map to string map
		headersMap := make(map[string]string)
		for k, v := range response.Headers {
			if strVal, ok := v.(string); ok {
				headersMap[k] = strVal
			}
		}
		headersValue, d := types.MapValueFrom(ctx, types.StringType, headersMap)
		diags.Append(d...)
		model.Headers = headersValue
	} else {
		model.Headers = types.MapNull(types.StringType)
	}

	// Convert secret headers.
	if response.SecretHeaders != nil {
		// Convert interface{} map to string map
		secretHeadersMap := make(map[string]string)
		for k, v := range response.SecretHeaders {
			if strVal, ok := v.(string); ok {
				secretHeadersMap[k] = strVal
			}
		}
		secretHeadersValue, d := types.MapValueFrom(ctx, types.StringType, secretHeadersMap)
		diags.Append(d...)
		model.SecretHeaders = secretHeadersValue
	} else {
		model.SecretHeaders = types.MapNull(types.StringType)
	}

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

	// Convert headers.
	if response.Headers != nil {
		// Convert interface{} map to string map
		headersMap := make(map[string]string)
		for k, v := range response.Headers {
			if strVal, ok := v.(string); ok {
				headersMap[k] = strVal
			}
		}
		headersValue, d := types.MapValueFrom(ctx, types.StringType, headersMap)
		diags.Append(d...)
		model.Headers = headersValue
	} else {
		model.Headers = types.MapNull(types.StringType)
	}

	// Convert secret headers.
	if response.SecretHeaders != nil {
		// Convert interface{} map to string map
		secretHeadersMap := make(map[string]string)
		for k, v := range response.SecretHeaders {
			if strVal, ok := v.(string); ok {
				secretHeadersMap[k] = strVal
			}
		}
		secretHeadersValue, d := types.MapValueFrom(ctx, types.StringType, secretHeadersMap)
		diags.Append(d...)
		model.SecretHeaders = secretHeadersValue
	} else {
		model.SecretHeaders = types.MapNull(types.StringType)
	}

	// Convert success rate.
	if response.SuccessRate != nil {
		if numVal, ok := response.SuccessRate.(float64); ok {
			model.SuccessRate = types.NumberValue(big.NewFloat(numVal))
		} else if intVal, ok := response.SuccessRate.(int); ok {
			model.SuccessRate = types.NumberValue(big.NewFloat(float64(intVal)))
		} else {
			model.SuccessRate = types.NumberNull()
		}
	} else {
		model.SuccessRate = types.NumberNull()
	}

	return diags
}
