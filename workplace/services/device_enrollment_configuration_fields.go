package services

import (
	"context"
	"fmt"

	"terraform-provider-microsoft365wp/workplace/generic"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ponytail: extend generated schemas here until the generator includes these fields.
func init() {
	esp := DeviceEnrollmentConfigurationResource.SpecificSchema.Attributes["windows10_esp"].(generic.OdataDerivedTypeNestedAttributeRs)
	// The singular data source was initialized with its own nested attribute map.
	espDS := DeviceEnrollmentConfigurationSingularDataSource.SpecificSchema.Attributes["windows10_esp"].(generic.OdataDerivedTypeNestedAttributeDs)
	for name, description := range map[string]string{
		"install_quality_updates":             "Allows quality updates installation during OOBE.",
		"allow_non_blocking_app_installation": "Allows non-blocking app installation failures during the technician phase.",
	} {
		esp.Attributes[name] = schema.BoolAttribute{
			Optional: true, Computed: true,
			MarkdownDescription: description + " Omit to use the service default on creation or preserve the current service value on update. The provider writes this setting separately from other ESP settings and verifies the result.",
		}
		espDS.Attributes[name] = dsschema.BoolAttribute{Computed: true, MarkdownDescription: description}
	}
	DeviceEnrollmentConfigurationResource.AccessParams.TerraformToGraphMiddleware = deviceEnrollmentConfigurationFieldsMiddleware
	DeviceEnrollmentConfigurationResource.AccessParams.WriteOptions.SubActions = append(
		[]generic.WriteSubAction{&enrollmentESPFieldsAction{}},
		DeviceEnrollmentConfigurationResource.AccessParams.WriteOptions.SubActions...,
	)
}

const enrollmentESPType = "#microsoft.graph.windows10EnrollmentCompletionPageConfiguration"

var enrollmentESPFields = map[string]string{
	"install_quality_updates":             "installQualityUpdates",
	"allow_non_blocking_app_installation": "allowNonBlockingAppInstallation",
}

func deviceEnrollmentConfigurationFieldsMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if err := deviceEnrollmentConfigurationTerraformToGraphMiddleware(ctx, diags, params); err != nil {
		return err
	}
	var priority types.Int64
	diags.Append(params.Config.GetAttribute(ctx, path.Root("priority"), &priority)...)
	if diags.HasError() {
		return nil
	}
	if priority.IsNull() {
		delete(params.RawVal, "priority")
	} else if priority.IsUnknown() {
		return fmt.Errorf("priority must be known before writing the enrollment configuration")
	}
	if params.RawVal["@odata.type"] != enrollmentESPType {
		return nil
	}
	for name, graphName := range enrollmentESPFields {
		var configured types.Bool
		diags.Append(params.Config.GetAttribute(ctx, path.Root("windows10_esp").AtName(name), &configured)...)
		if diags.HasError() {
			return nil
		}
		if configured.IsUnknown() {
			diags.AddError("Unknown ESP setting", fmt.Sprintf("%s must be known before writing the enrollment configuration.", name))
		} else if configured.IsNull() {
			delete(params.RawVal, graphName)
		}
	}
	return nil
}

type enrollmentESPFieldsAction struct{ generic.WriteSubActionBase }

type enrollmentESPFieldsWrite struct {
	desired, before map[string]any
}

func (a *enrollmentESPFieldsAction) ExecutePre(ctx context.Context, diags *diag.Diagnostics, req *generic.WriteSubActionRequest) {
	if req.RawVal["@odata.type"] != enrollmentESPType {
		return
	}
	values := enrollmentESPFieldsWrite{desired: map[string]any{}, before: map[string]any{}}
	for _, graphName := range enrollmentESPFields {
		if value, configured := req.RawVal[graphName]; configured {
			values.desired[graphName] = value
			delete(req.RawVal, graphName)
		}
	}
	// Default-profile creation also updates an existing ID, without a prior Terraform state.
	if req.Id != "" {
		values.before = enrollmentReadESPFields(ctx, diags, req)
		if diags.HasError() {
			return
		}
		for _, graphName := range enrollmentESPFields {
			if _, configured := values.desired[graphName]; !configured {
				values.desired[graphName] = values.before[graphName]
			}
		}
	}
	req.SubActionsData["espFields"] = values
}

func (a *enrollmentESPFieldsAction) ExecutePost(ctx context.Context, diags *diag.Diagnostics, req *generic.WriteSubActionRequest) {
	values, present := req.SubActionsData["espFields"].(enrollmentESPFieldsWrite)
	if !present {
		return
	}
	// The generic writer copies Plan before post-actions. Do not retain unverified desired values on failure.
	for name, graphName := range enrollmentESPFields {
		prior := types.BoolNull()
		if value, known := values.before[graphName].(bool); known {
			prior = types.BoolValue(value)
		}
		diags.Append(req.RespState.SetAttribute(ctx, path.Root("windows10_esp").AtName(name), prior)...)
	}
	if diags.HasError() {
		return
	}
	if len(values.desired) > 0 {
		// Graph can reset these values in a full ESP PATCH; send them separately after the parent write.
		body := map[string]any{"@odata.type": enrollmentESPType}
		for name, value := range values.desired {
			body[name] = value
		}
		req.GenRes.AccessParams.UpdateRaw(ctx, diags, "", req.Id, req.IdAttributer, body)
	}
	// A rejected PATCH can still have changed the profile; retain the observed values even on error.
	var readDiags diag.Diagnostics
	observed := enrollmentReadESPFields(ctx, &readDiags, req)
	diags.Append(readDiags...)
	if readDiags.HasError() {
		return
	}
	for name, graphName := range enrollmentESPFields {
		diags.Append(req.RespState.SetAttribute(ctx, path.Root("windows10_esp").AtName(name), observed[graphName])...)
		if wanted, requested := values.desired[graphName]; requested && observed[graphName] != wanted {
			diags.AddError("ESP setting did not converge", fmt.Sprintf("Graph did not retain %s; state contains the observed value.", name))
		}
	}
}

func enrollmentReadESPFields(ctx context.Context, diags *diag.Diagnostics, req *generic.WriteSubActionRequest) map[string]any {
	uri := req.GenRes.AccessParams.GetUriWithIdForUD(ctx, diags, "", req.Id, req.IdAttributer)
	if diags.HasError() {
		return nil
	}
	observed := req.GenRes.AccessParams.ReadRaw(ctx, diags, uri.Entity, false)
	if diags.HasError() {
		return nil
	}
	if observed["id"] != req.Id {
		diags.AddError("Unable to verify ESP settings", "Graph did not return the requested profile; retaining the last known values.")
		return nil
	}
	for _, graphName := range enrollmentESPFields {
		if _, valid := observed[graphName].(bool); !valid {
			diags.AddError("Unable to verify ESP settings", "Graph did not return both ESP settings as booleans; retaining the last known values.")
			return nil
		}
	}
	return observed
}
