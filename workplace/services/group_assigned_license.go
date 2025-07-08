package services

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	GroupAssignedLicenseResource = generic.GenericResource{
		TypeNameSuffix: "group_assigned_license",
		SpecificSchema: groupAssignedLicenseResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/groups",
			ParentEntities: generic.ParentEntities{
				{
					ParentIdField: path.Root("group_id"),
				},
			},
			UriNoId: true,
			EntityId: generic.EntityIdOptions{
				AttrNameGraph: "skuId",
			},
			ReadOptions: generic.ReadOptions{
				ODataSelect: []string{"assignedLicenses"},
				DataSource: generic.DataSourceOptions{
					NoFilterSupport: true,
				},
			},
			GraphToTerraformMiddleware:                     groupAssignedLicenseGraphToTerraformMiddleware,
			GraphToTerraformMiddlewareTargetSetRunOnRawVal: true,
			CreateReplaceFunc:                              groupAssignedLicenseCreateReplaceFunc,
			DeleteReplaceFunc:                              groupAssignedLicenseDeleteReplaceFunc,
		},
	}

	GroupAssignedLicenseSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&GroupAssignedLicenseResource)

	GroupAssignedLicensePluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&GroupAssignedLicenseResource, "")
)

func groupAssignedLicenseGraphToTerraformMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {

	assignedLicenses, ok1 := params.RawVal["assignedLicenses"].([]any)
	if !ok1 {
		return errors.New("property 'assignedLicenses' not found or of wrong type")
	}

	clear(params.RawVal)
	if !params.IsTargetSetOnRawVal {
		// resource or singular data source
		for _, assignedLicenseAny := range assignedLicenses {
			if assignedLicenseMap, ok2 := assignedLicenseAny.(map[string]any); ok2 && assignedLicenseMap["skuId"] == params.ExpectedId {
				maps.Copy(params.RawVal, assignedLicenseMap)
			}
		}
		// empty map gets translated to "not found" upstream
		tflog.Info(ctx, fmt.Sprintf("No element found with skuId '%s' in property assignedLicenses", params.ExpectedId))
	} else {
		// plural data source
		params.RawVal["value"] = assignedLicenses
	}

	return nil
}

func groupAssignedLicenseCreateReplaceFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.CreateReplaceFuncParams) {

	groupAssignedLicensePostAssignLicense(ctx, diags, &params.R.AccessParams, params.BaseUri, params.IdAttributer, params.Client, []any{params.RawVal}, nil)
	if diags.HasError() {
		return
	}

	params.Id, _ = params.RawVal["skuId"].(string)
}

func groupAssignedLicenseDeleteReplaceFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.DeleteReplaceFuncParams) {

	id := params.R.AccessParams.GetId(ctx, diags, params.Id, params.IdAttributer)
	if diags.HasError() {
		return
	}

	groupAssignedLicensePostAssignLicense(ctx, diags, &params.R.AccessParams, params.BaseUri, params.IdAttributer, params.Client, nil, []string{id})
}

func groupAssignedLicensePostAssignLicense(ctx context.Context, diags *diag.Diagnostics, aps *generic.AccessParams,
	baseUri string, idAttributer generic.GetAttributer, client *msgraph.Client,
	addLicenses []any, removeLicenses []string) {

	uri := aps.GetBaseUri(ctx, diags, baseUri, idAttributer)
	if diags.HasError() {
		return
	}
	uri.Entity += "/assignLicense"

	postRawVal := map[string]any{
		"addLicenses":    addLicenses,
		"removeLicenses": removeLicenses,
	}
	generic.CreateRaw(ctx, diags, client, uri, postRawVal, []int{http.StatusAccepted}, true, false)
}

var groupAssignedLicenseResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // assignedLicense
		"group_id": schema.StringAttribute{
			Required:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
		"disabled_plans": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			PlanModifiers: []planmodifier.Set{
				wpdefaultvalue.SetDefaultValueEmpty(),
				setplanmodifier.RequiresReplace(),
			},
			Computed:            true,
			MarkdownDescription: "A collection of the unique identifiers for plans that have been disabled. IDs are available in **servicePlans** > **servicePlanId** in the tenant's [subscribedSkus](../resources/subscribedsku.md) or **serviceStatus** > **servicePlanId** in the tenant's [companySubscription](../resources/subscribedsku.md). The _provider_ default value is `[]`.",
		},
		"sku_id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "The unique identifier for the SKU. Corresponds to the **skuId** from [subscribedSkus](../resources/subscribedsku.md) or [companySubscription](../resources/companysubscription.md).",
		},
	},
	MarkdownDescription: "Represents a license assigned to a user or group. The **assignedLicenses** property of the [user](user.md) or [group](group.md) entitity is a collection of **assignedLicense** objects. / https://learn.microsoft.com/en-us/graph/api/resources/assignedlicense?view=graph-rest-beta\n\nProvider Note: To import this resource, an ID consisting of `group_id` and `sku_id` being joined by a forward slash (`/`) must be used. ||| MS Graph: Licenses and subscriptions",
}
