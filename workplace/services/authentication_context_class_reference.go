package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	AuthenticationContextClassReferenceResource = generic.GenericResource{
		TypeNameSuffix: "authentication_context_class_reference",
		SpecificSchema: authenticationContextClassReferenceResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/identity/conditionalAccess/authenticationContextClassReferences",
			ReadOptions:     authenticationContextClassReferenceReadOptions,
			WriteSubActions: authenticationContextClassReferenceWriteSubActions,
		},
	}

	AuthenticationContextClassReferenceSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&AuthenticationContextClassReferenceResource)

	AuthenticationContextClassReferencePluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&AuthenticationContextClassReferenceResource, "")
)

var authenticationContextClassReferenceReadOptions = generic.ReadOptions{
	ValidStatusCodesExtra: []int{201},
	DataSource: generic.DataSourceOptions{
		NoFilterSupport: true,
		Plural: generic.PluralOptions{
			NoSelectSupport: true,
		},
	},
}

var authenticationContextClassReferenceWriteSubActions = []generic.WriteSubAction{
	&authenticationContextClassReferenceClearIsAvailableBeforeDeleteWsa{},
}

var authenticationContextClassReferenceResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // authenticationContextClassReference
		"id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Identifier used to reference the authentication context class. The ID is used to trigger step-up authentication for the referenced authentication requirements and is the value that will be issued in the `acrs` claim of an access token. This value in the claim is used to verify that the required authentication context has been satisfied. The allowed values are `c1` through `c25`. <br/> Supports `$filter` (`eq`).",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "A short explanation of the policies that are enforced by authenticationContextClassReference. This value should be used to provide secondary text to describe the authentication context class reference when building user facing admin experiences. For example, selection UX. The _provider_ default value is `\"\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "A friendly name that identifies the authenticationContextClassReference object when building user-facing admin experiences. For example, a selection UX.",
		},
		"is_available": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whether the authenticationContextClassReference has been published by the security admin and is ready for use by apps. When it's set to `false`, it shouldn't be shown in selection UX used to tag resources with authentication context class values. It will still be shown in the Conditional Access policy authoring experience. <br/> Supports `$filter` (`eq`). The _provider_ default value is `false`.",
		},
	},
	MarkdownDescription: "Represents a Microsoft Entra authentication context class reference. Authentication context class references are custom values that define a Conditional Access authentication requirement. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationcontextclassreference?view=graph-rest-beta ||| MS Graph: Conditional access",
}

//
// authenticationContextClassReferenceClearIsAvailableBeforeDeleteWsa
//

type authenticationContextClassReferenceClearIsAvailableBeforeDeleteWsa struct{}

func (*authenticationContextClassReferenceClearIsAvailableBeforeDeleteWsa) Initialize() {
}

func (*authenticationContextClassReferenceClearIsAvailableBeforeDeleteWsa) CheckRunAction(wsaOperation generic.OperationType) bool {
	return wsaOperation == generic.OperationDelete
}

func (*authenticationContextClassReferenceClearIsAvailableBeforeDeleteWsa) ExecutePre(ctx context.Context, diags *diag.Diagnostics, wsaReq *generic.WriteSubActionRequest) {
	// will be called on Delete only
	var isAvailableTypes types.Bool // might be null
	diags.Append(wsaReq.ReqState.GetAttribute(ctx, path.Root("is_available"), &isAvailableTypes)...)
	if diags.HasError() {
		return
	}
	if isAvailableTypes.ValueBool() {
		wsaReq.GenRes.AccessParams.UpdateRaw(ctx, diags, "", wsaReq.Id, wsaReq.IdAttributer, map[string]any{"isAvailable": false}, false)
	}
}

func (*authenticationContextClassReferenceClearIsAvailableBeforeDeleteWsa) ExecutePost(ctx context.Context, diags *diag.Diagnostics, wsaReq *generic.WriteSubActionRequest) {
}
