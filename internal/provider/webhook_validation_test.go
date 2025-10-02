// Copyright (c) Plain Technologies Aps

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/plain-insure/terraform-provider-printone/internal/provider/resource_webhook"
)

func TestValidateFilters(t *testing.T) {
	tests := []struct {
		name          string
		filters       []resource_webhook.WebhookFilterModel
		expectError   bool
		errorContains string
	}{
		{
			name: "valid equals filter",
			filters: []resource_webhook.WebhookFilterModel{
				{
					Key:    types.StringValue("status"),
					Event:  types.StringValue("order_status_update"),
					Type:   types.StringValue("equals"),
					Value:  types.StringValue("order_created"),
					Values: types.ListNull(types.StringType),
				},
			},
			expectError: false,
		},
		{
			name: "valid in filter",
			filters: []resource_webhook.WebhookFilterModel{
				{
					Key:   types.StringValue("status"),
					Event: types.StringValue("order_status_update"),
					Type:  types.StringValue("in"),
					Value: types.StringNull(),
					Values: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("order_created"),
						types.StringValue("order_updated"),
					}),
				},
			},
			expectError: false,
		},
		{
			name: "invalid event type",
			filters: []resource_webhook.WebhookFilterModel{
				{
					Key:    types.StringValue("status"),
					Event:  types.StringValue("invalid_event"),
					Type:   types.StringValue("equals"),
					Value:  types.StringValue("test"),
					Values: types.ListNull(types.StringType),
				},
			},
			expectError:   true,
			errorContains: "Invalid Filter Event",
		},
		{
			name: "invalid filter type",
			filters: []resource_webhook.WebhookFilterModel{
				{
					Key:    types.StringValue("status"),
					Event:  types.StringValue("order_status_update"),
					Type:   types.StringValue("invalid_type"),
					Value:  types.StringValue("test"),
					Values: types.ListNull(types.StringType),
				},
			},
			expectError:   true,
			errorContains: "Invalid Filter Type",
		},
		{
			name: "equals filter missing value",
			filters: []resource_webhook.WebhookFilterModel{
				{
					Key:    types.StringValue("status"),
					Event:  types.StringValue("order_status_update"),
					Type:   types.StringValue("equals"),
					Value:  types.StringNull(),
					Values: types.ListNull(types.StringType),
				},
			},
			expectError:   true,
			errorContains: "Missing Filter Value",
		},
		{
			name: "equals filter with values set",
			filters: []resource_webhook.WebhookFilterModel{
				{
					Key:   types.StringValue("status"),
					Event: types.StringValue("order_status_update"),
					Type:  types.StringValue("equals"),
					Value: types.StringValue("test"),
					Values: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("value1"),
					}),
				},
			},
			expectError:   true,
			errorContains: "Invalid Filter Values",
		},
		{
			name: "in filter missing values",
			filters: []resource_webhook.WebhookFilterModel{
				{
					Key:    types.StringValue("status"),
					Event:  types.StringValue("order_status_update"),
					Type:   types.StringValue("in"),
					Value:  types.StringNull(),
					Values: types.ListNull(types.StringType),
				},
			},
			expectError:   true,
			errorContains: "Missing Filter Values",
		},
		{
			name: "in filter with value set",
			filters: []resource_webhook.WebhookFilterModel{
				{
					Key:   types.StringValue("status"),
					Event: types.StringValue("order_status_update"),
					Type:  types.StringValue("in"),
					Value: types.StringValue("test"),
					Values: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("value1"),
					}),
				},
			},
			expectError:   true,
			errorContains: "Invalid Filter Value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := validateFilters(tt.filters)

			if tt.expectError {
				if !diags.HasError() {
					t.Errorf("Expected error but got none")
					return
				}

				found := false
				for _, diag := range diags {
					if diag.Summary() == tt.errorContains {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error containing '%s' but got: %v", tt.errorContains, diags)
				}
			} else {
				if diags.HasError() {
					t.Errorf("Expected no error but got: %v", diags)
				}
			}
		})
	}
}
