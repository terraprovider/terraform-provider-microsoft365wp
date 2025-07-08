package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	WindowsManagementAppResource = generic.GenericResource{
		TypeNameSuffix: "windows_management_app",
		SpecificSchema: windowsManagementAppResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:           "/deviceAppManagement/windowsManagementApp",
			IsSingleton:       true,
			UpdateReplaceFunc: windowsManagementAppUpdateReplaceFunc,
		},
	}

	WindowsManagementAppSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&WindowsManagementAppResource)
)

func windowsManagementAppUpdateReplaceFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.UpdateReplaceFuncParams) {
	errorSummary := "Error updating resource"

	targetState, ok := params.RawVal["managedInstaller"].(string)
	if !ok {
		diags.AddError(errorSummary, "Unable to read target value of managedInstaller")
	}
	tflog.Info(ctx, "windowsManagementAppResourceUpdateReplaceFunc", map[string]any{"targetState": targetState})

	currentStateResultRaw := params.R.AccessParams.ReadRaw(ctx, diags, params.R.AccessParams.BaseUri, false)
	if diags.HasError() {
		return
	}
	currentState, ok := currentStateResultRaw["managedInstaller"].(string)
	if !ok {
		diags.AddError(errorSummary, "Unable to read current value of managedInstaller")
	}
	tflog.Info(ctx, "windowsManagementAppResourceUpdateReplaceFunc", map[string]any{"currentState": currentState})

	if targetState != currentState {
		params.R.AccessParams.CreateRaw2(ctx, diags, params.R.AccessParams.BaseUri+"/setAsManagedInstaller", params.IdAttributer, make(map[string]any), true, false)
	} else {
		tflog.Info(ctx, "windowsManagementAppResourceUpdateReplaceFunc: Nothing to do, target equals current state")
	}
}

var windowsManagementAppResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // windowsManagementApp
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Unique Identifier for the Windows management app",
		},
		"available_version": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Windows management app available version.",
		},
		"managed_installer": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("disabled", "enabled")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("disabled")},
			Computed:            true,
			MarkdownDescription: "Managed Installer Status. / ManagedInstallerStatus; possible values are: `disabled` (Managed Installer is Disabled), `enabled` (Managed Installer is Enabled). The _provider_ default value is `\"disabled\"`.",
		},
		"managed_installer_configured_date_time": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Managed Installer Configured Date Time",
		},
	},
	MarkdownDescription: "Windows management app entity. / https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-windowsmanagementapp?view=graph-rest-beta\n\nProvider Note: Please note that MS Graph does currently (as of 2025-02) not seem to allow to update (or even read) this entity\nusing application permissions but only using work or school account delegated permissions. ||| MS Graph: Device management",
}
