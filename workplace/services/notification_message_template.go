package services

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	NotificationMessageTemplateResource = generic.GenericResource{
		TypeNameSuffix: "notification_message_template",
		SpecificSchema: notificationMessageTemplateResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/deviceManagement/notificationMessageTemplates",
			ReadOptions: generic.ReadOptions{
				ODataExpand: "localizedNotificationMessages",
			},
			WriteOptions: generic.WriteOptions{
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionIndividual{
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"localizedNotificationMessages"},
							UriSuffix:  "localizedNotificationMessages",
						},
						ComparisonKeyAttribute: "locale",
						SetNestedPath:          tftypes.NewAttributePath().WithAttributeName("localized_notification_messages"),
						IdGetterFunc: func(ctx context.Context, diags *diag.Diagnostics, vRaw map[string]any, parentId string) string {
							return fmt.Sprintf("%s_%s", parentId, vRaw["locale"].(string))
						},
						TerraformToGraphMiddleware: func(_ context.Context, _ *diag.Diagnostics, p *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
							if p.IsUpdate {
								// locale cannot be updated
								delete(p.RawVal, "locale")
							}
							//lint:ignore S1002 Keep for clarity
							if p.RawVal["isDefault"].(bool) == false {
								// isDefault cannot be set to false, one must set another locale's isDefault to true instead
								delete(p.RawVal, "isDefault")
							}
							return nil
						},
					},
				},
			},
		},
	}

	NotificationMessageTemplateSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&NotificationMessageTemplateResource)

	NotificationMessageTemplatePluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&NotificationMessageTemplateResource, "")
)

var notificationMessageTemplateResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // notificationMessageTemplate
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Key of the entity.",
		},
		"branding_options": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("none", "includeCompanyLogo", "includeCompanyName", "includeContactInformation", "includeCompanyPortalLink", "includeDeviceDetails", "unknownFutureValue"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("includeCompanyLogo,includeCompanyName,includeContactInformation"),
			},
			Computed:            true,
			MarkdownDescription: "The Message Template Branding Options. Branding is defined in the Intune Admin Console. / Branding Options for the Message Template. Branding is defined in the Intune Admin Console; possible values are: `none` (Indicates that no branding options are set in the message template.), `includeCompanyLogo` (Indicates to include company logo in the message template.), `includeCompanyName` (Indicates to include company name in the message template.), `includeContactInformation` (Indicates to include contact information in the message template.), `includeCompanyPortalLink` (Indicates to include company portal website link in the message template.), `includeDeviceDetails` (Indicates to include device details in the message template.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.). The _provider_ default value is `\"includeCompanyLogo,includeCompanyName,includeContactInformation\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Display name for the Notification Message Template.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "DateTime the object was last modified.",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tags for this Entity instance. The _provider_ default value is `[\"0\"]`.",
		},
		"localized_notification_messages": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // localizedNotificationMessage
					"is_default": schema.BoolAttribute{
						Optional:            true,
						Default:             booldefault.StaticBool(false),
						Computed:            true,
						MarkdownDescription: `Flag to indicate whether or not this is the default locale for language fallback. This flag can only be set. To unset, set this property to true on another Localized Notification Message.`,
					},
					"locale": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Locale for which this message is destined.",
					},
					"message_template": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Message Template content.",
					},
					"subject": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Message Template Subject.",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "The list of localized messages for this Notification Message Template. / The text content of a Notification Message Template for the specified locale. / https://learn.microsoft.com/en-us/graph/api/resources/intune-notification-localizednotificationmessage?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
	},
	MarkdownDescription: "Notification messages are messages that are sent to end users who are determined to be not-compliant with the compliance policies defined by the administrator. Administrators choose notifications and configure them in the Intune Admin Console using the compliance policy creation page under the “Actions for non-compliance” section. Use the notificationMessageTemplate object to create your own custom notifications for administrators to choose while configuring actions for non-compliance. / https://learn.microsoft.com/en-us/graph/api/resources/intune-notification-notificationmessagetemplate?view=graph-rest-beta ||| MS Graph: Device management",
}
