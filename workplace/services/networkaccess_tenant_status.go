package services

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/generic"
	msgraphutils "terraform-provider-microsoft365wp/workplace/services/msgraph_utils"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	NetworkaccessTenantStatusResource = generic.GenericResource{
		TypeNameSuffix: "networkaccess_tenant_status",
		SpecificSchema: networkaccessTenantStatusResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/networkaccess/tenantStatus",
			IsSingleton:                true,
			GraphToTerraformMiddleware: networkaccessTenantStatusGraphToTerraformMiddleware,
			UpdateReplaceFunc:          networkaccessTenantStatusUpdateReplaceFunc,
		},
	}

	NetworkaccessTenantStatusSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&NetworkaccessTenantStatusResource)
)

func networkaccessTenantStatusGraphToTerraformMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {

	// On reading, infer the value for `activate` from `onboardingStatus`
	// (for TF to be able to detect the necessity for an update even if importing, but do not set in ambiguous situations)
	switch params.RawVal["onboardingStatus"] {
	case "onboarded":
		params.RawVal["activate"] = true
	case "offboarded":
		params.RawVal["activate"] = false
	}

	return msgraphutils.SingletonSyntheticIdGraphToTerraformMiddleware(ctx, diags, params)
}

func networkaccessTenantStatusUpdateReplaceFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.UpdateReplaceFuncParams) {
	errorSummary := "Error updating resource"

	// we need to recheck here since this a singleton and TF might have lead us here on resource creation without having compared with prior state itself
	targetState, ok := params.RawVal["activate"].(bool)
	if !ok {
		diags.AddError(errorSummary, "Unable to read target value of activate")
	}
	tflog.Info(ctx, "networkaccessTenantStatusUpdateReplaceFunc", map[string]any{"targetState": targetState})

	currentState := networkaccessTenantStatusReadOnboardingState(ctx, diags, &params.R.AccessParams, errorSummary, true)
	if diags.HasError() {
		return
	}

	//lint:ignore S1002 Keep for clarity
	if targetState == true && currentState == "offboarded" {
		// initiate onboarding
		params.R.AccessParams.CreateRaw2(ctx, diags, "/networkAccess/microsoft.graph.networkaccess.onboard", params.IdAttributer, make(map[string]any), true, false)

		for {
			newState := networkaccessTenantStatusReadOnboardingState(ctx, diags, &params.R.AccessParams, errorSummary, false)
			if diags.HasError() {
				return
			}
			if newState == "onboarded" {
				// onboarding has completed successfully
				break
			}
			if newState != "offboarded" && newState != "onboardingInProgress" {
				diags.AddError(errorSummary, fmt.Sprintf("onboardingStatus does neither indicate progress nor success: `%s`", newState))
				return
			}

			time.Sleep(time.Millisecond * 500)
		}

		//lint:ignore S1002 Keep for clarity
	} else if (targetState == true && currentState == "onboarded") || (targetState == false && currentState == "offboarded") {
		// TF might have lead us here on resource creation without having compared with prior state itself
		tflog.Info(ctx, "networkaccessTenantStatusUpdateReplaceFunc: Nothing to do, target equals current state")

		//lint:ignore S1002 Keep for clarity
	} else if targetState == false && currentState == "onboarded" {
		diags.AddError(errorSummary, "Network access has already been activated and cannot be deactivated anymore")

		//lint:ignore S1002 Keep for clarity
	} else if (targetState == true && currentState == "onboardingInProgress") || (targetState == false && currentState == "offboardingInProgress") {
		diags.AddWarning("Resource is busy", "Onboarding or offboarding already is in progress")

	} else {
		diags.AddError(errorSummary, fmt.Sprintf("Target state is `%t` but current state is `%s`", targetState, currentState))
	}
}

func networkaccessTenantStatusReadOnboardingState(ctx context.Context, diags *diag.Diagnostics, accessParams *generic.AccessParams, errorSummary string, logAsInfo bool) string {

	onboardingStatusResultRaw := accessParams.ReadRaw(ctx, diags, accessParams.BaseUri, false)
	if diags.HasError() {
		return ""
	}

	onboardingStatus, ok := onboardingStatusResultRaw["onboardingStatus"].(string)
	if !ok {
		diags.AddError(errorSummary, "Unable to read current value of onboardingStatus")
	}

	logMsg := "networkaccessTenantStatusReadOnboardingState"
	logFields := map[string]any{"onboardingStatus": onboardingStatus}
	if logAsInfo {
		tflog.Info(ctx, logMsg, logFields)
	} else {
		tflog.Trace(ctx, logMsg, logFields)
	}

	return onboardingStatus
}

var networkaccessTenantStatusResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // networkaccess.tenantStatus
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Identifier.",
		},
		"onboarding_error_message": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Reflects a message to the user if there's an error.",
		},
		"onboarding_status": schema.StringAttribute{
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf("offboarded", "offboardingInProgress", "onboardingInProgress", "onboarded", "onboardingErrorOccurred", "offboardingErrorOccurred", "unknownFutureValue"),
			},
			MarkdownDescription: "Reflects the tenant onboarding status. / The onboarding status of the tenant. <br/> _Provider_ allowed values are: `offboarded`, `offboardingInProgress`, `onboardingInProgress`, `onboarded`, `onboardingErrorOccurred`, `offboardingErrorOccurred`, `unknownFutureValue`.",
		},
		"activate": schema.BoolAttribute{
			Required: true,
		},
	},
	MarkdownDescription: "Represents the status of the Global Secure Access services for the tenant. <br/> Also see [Microsoft docs for networkaccess.tenantStatus](https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-tenantstatus?view=graph-rest-beta). ||| MS Graph: Network access (preview)",
}
