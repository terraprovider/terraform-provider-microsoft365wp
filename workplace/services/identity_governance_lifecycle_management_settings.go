package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	IdentityGovernanceLifecycleManagementSettingsResource = generic.GenericResource{
		TypeNameSuffix: "identity_governance_lifecycle_management_settings",
		SpecificSchema: identityGovernanceLifecycleManagementSettingsResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/identityGovernance/lifecycleWorkflows/settings",
			IsSingleton:                true,
			GraphToTerraformMiddleware: identityGovernanceLifecycleManagementSettingsGraphToTerraformMiddleware,
		},
	}

	IdentityGovernanceLifecycleManagementSettingsSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&IdentityGovernanceLifecycleManagementSettingsResource)
)

func identityGovernanceLifecycleManagementSettingsGraphToTerraformMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {
	// MS Graph does indeed not return anything for id, so add something helpful
	if params.RawVal["id"] == nil {
		params.RawVal["id"] = "singleton"
	}
	return nil
}

var identityGovernanceLifecycleManagementSettingsResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // identityGovernance.lifecycleManagementSettings
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"email_settings": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // emailSettings
				"sender_domain": schema.StringAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("microsoft.com")},
					Computed:            true,
					MarkdownDescription: "Specifies the [domain](domain.md) that should be used when sending email notifications. This domain must be [verified](../api/domain-verify.md) in order to be used. We recommend that you use a domain that has the appropriate DNS records to facilitate email validation, like SPF, DKIM, DMARC, and MX, because this then complies with the [RFC compliance](https://www.ietf.org/rfc/rfc2142.txt) for sending and receiving email. For details, see [Learn more about Exchange Online Email Routing](/exchange/mail-flow-best-practices/mail-flow-best-practices). The _provider_ default value is `\"microsoft.com\"`.",
				},
				"use_company_branding": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Specifies if the organization’s banner logo should be included in email notifications. The banner logo will replace the Microsoft logo at the top of the email notification. If `true` the banner logo will be taken from the tenant’s [branding settings](organizationalbranding.md). This value can only be set to `true` if the [organizationalBranding](organizationalbranding.md) **bannerLogo** property is set. The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Defines the settings for emails sent out from email-specific [tasks](../resources/identitygovernance-task.md) within workflows. Accepts 2 parameters<br><br>senderDomain- Defines the domain of who is sending the email. <br>useCompanyBranding- A Boolean value that defines if company branding is to be used with the email. / Defines the settings for emails sent from Lifecycle workflow [tasks](identitygovernance-task.md). Allows you to use a verified custom [domain](domain.md) and [organizationalBranding](organizationalbranding.md) with emails sent out via workflow tasks. / https://learn.microsoft.com/en-us/graph/api/resources/emailsettings?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"workflow_schedule_interval_in_hours": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(3)},
			Computed:            true,
			MarkdownDescription: "The interval in hours at which all [workflows](../resources/identitygovernance-workflow.md) running in the tenant should be scheduled for execution. This interval has a minimum value of 1 and a maximum value of 24. The default value is 3 hours. The _provider_ default value is `3`.",
		},
	},
	MarkdownDescription: "The settings of Microsoft Entra Lifecycle Workflows in the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-lifecyclemanagementsettings?view=graph-rest-beta ||| MS Graph: Lifecycle workflows",
}
