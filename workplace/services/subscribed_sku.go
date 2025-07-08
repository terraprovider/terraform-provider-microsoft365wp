package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	subscribedSkuResource = generic.GenericResource{
		TypeNameSuffix: "subscribed_sku",
		SpecificSchema: subscribedSkuResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/subscribedSkus",
			ReadOptions: generic.ReadOptions{
				DataSource: generic.DataSourceOptions{
					NoFilterSupport: true,
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"account_id", "account_name", "applies_to", "capability_status", "sku_id", "sku_part_number"},
					},
				},
			},
		},
	}

	SubscribedSkuSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&subscribedSkuResource)

	SubscribedSkuPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&subscribedSkuResource, "")
)

var subscribedSkuResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // subscribedSku
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the subscribed SKU object. Key, not nullable.",
		},
		"account_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The unique ID of the account this SKU belongs to.",
		},
		"account_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The name of the account this SKU belongs to.",
		},
		"applies_to": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The target class for this SKU. Only SKUs with target class `User` are assignable.",
		},
		"capability_status": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "`Enabled` indicates that the **prepaidUnits** property has at least one unit that is enabled. `LockedOut` indicates that the customer canceled their subscription.",
		},
		"consumed_units": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "The number of licenses that have been assigned.",
		},
		"prepaid_units": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // licenseUnitsDetail
				"enabled": schema.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "The number of units that are enabled for the active subscription of the service SKU.",
				},
				"locked_out": schema.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "The number of units that are locked out because the customer canceled their subscription of the service SKU.",
				},
				"suspended": schema.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "The number of units that are suspended because the subscription of the service SKU has been canceled. The units can't be assigned but can still be reactivated before they're deleted.",
				},
				"warning": schema.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "The number of units that are in warning status. When the subscription of the service SKU has expired, the customer has a grace period to renew their subscription before it's canceled (moved to a `suspended` state).",
				},
			},
			MarkdownDescription: "Information about the number and status of prepaid licenses. / The **prepaidUnits** property of the [subscribedSku](subscribedsku.md) entity is of type **licenseUnitsDetail**. For more information on the progression states of a subscription, see [What if my subscription expires?](/microsoft-365/commerce/subscriptions/what-if-my-subscription-expires?view=o365-worldwide&preserve-view=true) / https://learn.microsoft.com/en-us/graph/api/resources/licenseunitsdetail?view=graph-rest-beta",
		},
		"service_plans": schema.SetNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // servicePlanInfo
					"applies_to": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The object the service plan can be assigned to. The possible values are: <br/>`User` - service plan can be assigned to individual users.<br/>`Company` - service plan can be assigned to the entire tenant.",
					},
					"provisioning_status": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The provisioning status of the service plan. The possible values are:<br/>`Success` - Service is fully provisioned.<br/>`Disabled` - Service is disabled.<br/>`Error` - The service plan isn't provisioned and is in an error state.<br/>`PendingInput` - The service isn't provisioned and is awaiting service confirmation.<br/>`PendingActivation` - The service is provisioned but requires explicit activation by an administrator (for example, Intune_O365 service plan)<br/>`PendingProvisioning` - Microsoft has added a new service to the product SKU and it isn't activated in the tenant.",
					},
					"service_plan_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The unique identifier of the service plan.",
					},
					"service_plan_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The name of the service plan.",
					},
				},
			},
			MarkdownDescription: "Information about the service plans that are available with the SKU. Not nullable / Contains information about a service plan associated with a subscribed SKU. The **servicePlans** property of the [subscribedSku](subscribedsku.md) entity is a collection of **servicePlanInfo**. / https://learn.microsoft.com/en-us/graph/api/resources/serviceplaninfo?view=graph-rest-beta",
		},
		"sku_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The unique identifier (GUID) for the service SKU.",
		},
		"sku_part_number": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The SKU part number; for example, `AAD_PREMIUM` or `RMSBASIC`. To get a list of commercial subscriptions that an organization has acquired, see [List subscribedSkus](../api/subscribedsku-list.md).",
		},
		"subscription_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "A list of all [subscription IDs](../resources/companysubscription.md) associated with this SKU.",
		},
	},
	MarkdownDescription: "Represents information about a service SKU that a company is subscribed to. Use the values of **skuId** and **servicePlans** > **servicePlanId** to assign licenses to unassigned users and groups through the [user: assignLicense](../api/user-assignlicense.md) and [group: assignLicense](../api/group-assignlicense.md) APIs respectively.\n\nFor more information about subscriptions and licenses, see [Subscriptions, licenses, accounts, and tenants for Microsoft's cloud offerings](/microsoft-365/enterprise/subscriptions-licenses-accounts-and-tenants-for-microsoft-cloud-offerings). / https://learn.microsoft.com/en-us/graph/api/resources/subscribedsku?view=graph-rest-beta\n\nProvider Note: Please note that MS Graph does currently (as of 2025-07) not support OData filtering on this endpoint. ||| MS Graph: Licenses and subscriptions",
}
