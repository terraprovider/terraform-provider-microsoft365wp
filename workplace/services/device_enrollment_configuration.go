package services

import (
	"fmt"
	"regexp"
	"strings"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/manicminer/hamilton/msgraph"
)

var (
	DeviceEnrollmentConfigurationResource = generic.GenericResource{
		TypeNameSuffix: "device_enrollment_configuration",
		SpecificSchema: deviceEnrollmentConfigurationResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/deviceManagement/deviceEnrollmentConfigurations",
			ReadOptions:                deviceEnrollmentConfigurationReadOptions,
			WriteSubActions:            deviceEnrollmentConfigurationWriteSubActions,
			TerraformToGraphMiddleware: deviceEnrollmentConfigurationTerraformToGraphMiddleware,
			GraphToTerraformMiddleware: deviceEnrollmentConfigurationGraphToTerraformMiddleware,
			CreateModifyFunc:           deviceEnrollmentConfigurationCreateModifyFunc,
			DeleteModifyFunc:           deviceEnrollmentConfigurationDeleteModifyFunc,
			SerializeWrites:            true,
			SerialWritesDelay:          time.Second * 3,
		},
	}

	DeviceEnrollmentConfigurationSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceEnrollmentConfigurationResource)

	DeviceEnrollmentConfigurationPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&DeviceEnrollmentConfigurationSingularDataSource, "")
)

var deviceEnrollmentConfigurationOdataToConfigType = map[string]string{
	"windows10EnrollmentCompletionPageConfiguration":       "windows10EnrollmentCompletionPageConfiguration",
	"deviceEnrollmentLimitConfiguration":                   "limit",
	"deviceEnrollmentPlatformRestrictionConfiguration":     "singlePlatformRestriction",
	"deviceEnrollmentPlatformRestrictionsConfiguration":    "platformRestrictions",
	"deviceEnrollmentWindowsHelloForBusinessConfiguration": "windowsHelloForBusiness",
}
var kDefConfigPriority float64 = 0
var kDefConfigDisplayName = "###TfWorkplaceDefault###"
var kDefConfigIdRegex, _ = regexp.Compile("(?i)_Default[a-z0-9]+$")
var kSinglePlatformRestrictionOdataType = deviceEnrollmentConfigurationResourceSchema.Attributes["single_platform_restriction"].(generic.OdataDerivedTypeNestedAttributeRs).DerivedType

var deviceEnrollmentConfigurationReadOptions = generic.ReadOptions{
	// do not use ODataFilter here as MS Graph does not seem to behave fully logical with any filter tried yet
	ODataExpand:           "assignments",
	PluralNoFilterSupport: true,
}

var deviceEnrollmentConfigurationWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			AttributesMap: map[string]string{"assignments": "enrollmentConfigurationAssignments"},
			UriSuffix:     "assign",
		},
	},
	&generic.WriteSubActionAllInOne{
		GraphErrorsAsWarnings: true,
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"priority"},
			UriSuffix:  "setPriority",
		},
	},
}

func deviceEnrollmentConfigurationTerraformToGraphMiddleware(params generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	priorityAny, priorityOk := params.RawVal["priority"]
	if priorityOk {
		if priorityAny == nil {
			// priority must be removed if it is not set in the config as it is non-nullable
			delete(params.RawVal, "priority")
		} else {
			priority := priorityAny.(float64)
			if priority == 0 {
				// id might not yet be available to assert that default config is targeted, so we assert using display name
				displayName, displayNameOk := params.RawVal["displayName"].(string)
				if !(displayNameOk && displayName == kDefConfigDisplayName) {
					return fmt.Errorf("priority is %f but display name does not equal marker (it is \"%s\" instead of expected \"%s\")", priority, displayName, kDefConfigDisplayName)
				}
				// remove unchangeable attributes
				delete(params.RawVal, "displayName")
				delete(params.RawVal, "description")
				delete(params.RawVal, "roleScopeTagIds")
				delete(params.RawVal, "assignments")
				delete(params.RawVal, "priority")
			}
		}
	}
	odataType, odataTypeOk := params.RawVal["@odata.type"].(string)
	if odataTypeOk && odataType == kSinglePlatformRestrictionOdataType && params.IsUpdate {
		// platformType cannot be updated and may not even be written again after creation
		delete(params.RawVal, "platformType")
	}
	return nil
}

func deviceEnrollmentConfigurationGraphToTerraformMiddleware(params generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {
	priority, priorityOk := params.RawVal["priority"].(float64)
	if priorityOk && priority == 0 {
		// displayName cannot be our marker to assert that it's the default config when reading from MS Graph, so we assert using id
		id, idOk := params.RawVal["id"].(string)
		if !(idOk && kDefConfigIdRegex.MatchString(id)) {
			return fmt.Errorf("priority is %f but id does not match assertion regex (id is %s)", priority, id)
		}
		// set unchangeable attributes to convention values
		params.RawVal["displayName"] = kDefConfigDisplayName
		params.RawVal["description"] = ""             // schema default value
		params.RawVal["roleScopeTagIds"] = []any{"0"} // schema default value
		params.RawVal["assignments"] = []any{}        // schema default value
	}
	return nil
}

func deviceEnrollmentConfigurationCreateModifyFunc(params *generic.CreateModifyFuncParams) {
	errorSummary := "Error determining MS Graph id of device enrollment default configuration"

	plan := params.Req.Plan

	var priority float64 = -1
	plan.GetAttribute(params.Ctx, path.Root("priority"), &priority)
	if priority == kDefConfigPriority {
		var displayName string
		plan.GetAttribute(params.Ctx, path.Root("display_name"), &displayName)
		if displayName != kDefConfigDisplayName {
			params.Diags.AddError(errorSummary, fmt.Sprintf("Priority is 0, but display_name does not equal marker (it is \"%s\" instead of expected \"%s\")", displayName, kDefConfigDisplayName))
			return
		}

		// Apparently the default config is targeted, so try to determine the id of the existing entity for it to get
		// updated instead of a new entity to get created.
		odataType := ""
		for k, v := range plan.Schema.GetAttributes() {
			if derivedTypeNested, ok := v.(generic.DerivedTypeNestedAttribute); ok {
				rawTfValue, _, _ := tftypes.WalkAttributePath(plan.Raw, tftypes.NewAttributePath().WithAttributeName(k))
				tfValue, ok := rawTfValue.(tftypes.Value)
				if ok && !tfValue.IsNull() {
					odataType = derivedTypeNested.GetDerivedType()
					break
				}
			}
		}
		if !strings.HasPrefix(odataType, "#microsoft.graph.") {
			params.Diags.AddError(errorSummary, "Failed to determine actual OData type from resource config")
			return
		}

		// isof('microsoft.graph.*') does not work, so we need to use deviceEnrollmentConfigurationType
		configType, ok := deviceEnrollmentConfigurationOdataToConfigType[strings.ReplaceAll(odataType, "#microsoft.graph.", "")]
		if !ok {
			panic(fmt.Errorf("unknown MS graph type %s", odataType))
		}
		odataFilter := fmt.Sprintf("deviceEnrollmentConfigurationType eq '%s' and priority eq 0", configType)
		defaultIdResultRaw := params.R.AccessParams.ReadRaw(params.Ctx, params.Diags, msgraph.Uri{Entity: params.R.AccessParams.BaseUri}, "", odataFilter, []string{"id"}, false)
		if params.Diags.HasError() {
			return
		}

		defConfigId := ""
		if defaultIdItems, ok := defaultIdResultRaw["value"].([]any); ok && len(defaultIdItems) == 1 {
			if id, ok := defaultIdItems[0].(map[string]any)["id"].(string); ok {
				defConfigId = id
			}
		}
		if defConfigId == "" {
			params.Diags.AddError(errorSummary, "Azure query result did not have the expected format")
			return
		}
		if !kDefConfigIdRegex.MatchString(defConfigId) {
			params.Diags.AddError(errorSummary, fmt.Sprintf("id returned from Azure query does not match assertion regex (id is %s)", defConfigId))
		}

		params.UpdateExisting = true
		params.Id = defConfigId
	}
}

func deviceEnrollmentConfigurationDeleteModifyFunc(params *generic.DeleteModifyFuncParams) {
	// If it's a default config we must not delete it (it would probably fail anyway) but just skip deletion instead.
	// We cannot rely on priority here to assert that it's default config as the state might not have been populated
	// completely (e.g. in case of previous problems), so we only check the id here
	if kDefConfigIdRegex.MatchString(params.Id) {
		params.SkipDelete = true
	}
}

var deviceEnrollmentConfigurationResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceEnrollmentConfiguration
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: `Created date time in UTC of the device enrollment configuration`,
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: `The description of the device enrollment configuration`,
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `The display name of the device enrollment configuration`,
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: `Last modified date time in UTC of the device enrollment configuration`,
		},
		"priority": schema.Int64Attribute{
			Computed:            true,
			Optional:            true,
			MarkdownDescription: `Priority is used when a user exists in multiple groups that are assigned enrollment configuration. Users are subject only to the configuration with the lowest priority value.`,
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: `Optional role scope tags for the enrollment restrictions.`,
		},
		"version": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: `The version of the device enrollment configuration`,
		},
		"assignments": deviceAndAppManagementAssignment,
		"device_enrollment_limit": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceEnrollmentLimitConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // deviceEnrollmentLimitConfiguration
					"limit": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: `The maximum number of devices that a user can enroll`,
					},
				},
				Validators: []validator.Object{
					deviceEnrollmentConfigurationDeviceEnrollmentConfigurationValidator,
				},
				MarkdownDescription: `Device Enrollment Configuration that restricts the number of devices a user can enroll`,
			},
		},
		"platform_restrictions": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceEnrollmentPlatformRestrictionsConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // deviceEnrollmentPlatformRestrictionsConfiguration
					"android_for_work_restriction": schema.SingleNestedAttribute{
						Optional:            true,
						Attributes:          deviceEnrollmentConfigurationDeviceEnrollmentPlatformRestrictionAttributes,
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Android for work restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions`,
					},
					"android_restriction": schema.SingleNestedAttribute{
						Optional:            true,
						Attributes:          deviceEnrollmentConfigurationDeviceEnrollmentPlatformRestrictionAttributes,
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Android restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions`,
					},
					"ios_restriction": schema.SingleNestedAttribute{
						Optional:            true,
						Attributes:          deviceEnrollmentConfigurationDeviceEnrollmentPlatformRestrictionAttributes,
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Ios restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions`,
					},
					"macos_restriction": schema.SingleNestedAttribute{
						Optional:            true,
						Attributes:          deviceEnrollmentConfigurationDeviceEnrollmentPlatformRestrictionAttributes,
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						Description:         `macOSRestriction`, // custom MS Graph attribute name
						MarkdownDescription: `Mac restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions`,
					},
					"windows_restriction": schema.SingleNestedAttribute{
						Optional:            true,
						Attributes:          deviceEnrollmentConfigurationDeviceEnrollmentPlatformRestrictionAttributes,
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Windows restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions`,
					},
				},
				Validators: []validator.Object{
					deviceEnrollmentConfigurationDeviceEnrollmentConfigurationValidator,
				},
				MarkdownDescription: `Device Enrollment Configuration that restricts the types of devices a user can enroll`,
			},
		},
		"single_platform_restriction": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceEnrollmentPlatformRestrictionConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // deviceEnrollmentPlatformRestrictionConfiguration
					"platform_restriction": schema.SingleNestedAttribute{
						Required:            true,
						Attributes:          deviceEnrollmentConfigurationDeviceEnrollmentPlatformRestrictionAttributes,
						MarkdownDescription: `Restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions`,
					},
					"platform_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("allPlatforms", "ios", "windows", "windowsPhone", "android", "androidForWork", "mac", "linux", "unknownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: `Type of platform for which this restriction applies.`,
					},
				},
				Validators: []validator.Object{
					deviceEnrollmentConfigurationDeviceEnrollmentConfigurationValidator,
				},
				MarkdownDescription: `Device Enrollment Configuration that restricts the types of devices a user can enroll for a single platform`,
			},
		},
		"windows_hello_for_business": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceEnrollmentWindowsHelloForBusinessConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // deviceEnrollmentWindowsHelloForBusinessConfiguration
					"enhanced_biometrics_state": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Controls the ability to use the anti-spoofing features for facial recognition on devices which support it. If set to disabled, anti-spoofing features are not allowed. If set to Not Configured, the user can choose whether they want to use anti-spoofing.`,
					},
					"enhanced_sign_in_security": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: `Setting to configure Enhanced sign-in security. Default is Not Configured`,
					},
					"pin_expiration_in_days": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: `Controls the period of time (in days) that a PIN can be used before the system requires the user to change it. This must be set between 0 and 730, inclusive. If set to 0, the user's PIN will never expire`,
					},
					"pin_lowercase_characters_usage": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("allowed", "required", "disallowed"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("disallowed")},
						Computed:            true,
						MarkdownDescription: `Controls the ability to use lowercase letters in the Windows Hello for Business PIN.  Allowed permits the use of lowercase letter(s), whereas Required ensures they are present. If set to Not Allowed, lowercase letters will not be permitted.`,
					},
					"pin_maximum_length": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(127)},
						Computed:            true,
						MarkdownDescription: `Controls the maximum number of characters allowed for the Windows Hello for Business PIN. This value must be between 4 and 127, inclusive. This value must be greater than or equal to the value set for the minimum PIN.`,
					},
					"pin_minimum_length": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(4)},
						Computed:            true,
						MarkdownDescription: `Controls the minimum number of characters required for the Windows Hello for Business PIN.  This value must be between 4 and 127, inclusive, and less than or equal to the value set for the maximum PIN.`,
					},
					"pin_previous_block_count": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: `Controls the ability to prevent users from using past PINs. This must be set between 0 and 50, inclusive, and the current PIN of the user is included in that count. If set to 0, previous PINs are not stored. PIN history is not preserved through a PIN reset.`,
					},
					"pin_special_characters_usage": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("allowed", "required", "disallowed"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("disallowed")},
						Computed:            true,
						MarkdownDescription: `Controls the ability to use special characters in the Windows Hello for Business PIN.  Allowed permits the use of special character(s), whereas Required ensures they are present. If set to Not Allowed, special character(s) will not be permitted.`,
					},
					"pin_uppercase_characters_usage": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("allowed", "required", "disallowed"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("disallowed")},
						Computed:            true,
						MarkdownDescription: `Controls the ability to use uppercase letters in the Windows Hello for Business PIN.  Allowed permits the use of uppercase letter(s), whereas Required ensures they are present. If set to Not Allowed, uppercase letters will not be permitted.`,
					},
					"remote_passport_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: `Controls the use of Remote Windows Hello for Business. Remote Windows Hello for Business provides the ability for a portable, registered device to be usable as a companion for desktop authentication. The desktop must be Azure AD joined and the companion device must have a Windows Hello for Business PIN.`,
					},
					"security_device_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Controls whether to require a Trusted Platform Module (TPM) for provisioning Windows Hello for Business. A TPM provides an additional security benefit in that data stored on it cannot be used on other devices. If set to False, all devices can provision Windows Hello for Business even if there is not a usable TPM.`,
					},
					"security_key_for_sign_in": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Security key for Sign In provides the capacity for remotely turning ON/OFF Windows Hello Sercurity Keyl Not configured will honor configurations done on the clinet.`,
					},
					"state": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("enabled")},
						Computed:            true,
						MarkdownDescription: `Controls whether to allow the device to be configured for Windows Hello for Business. If set to disabled, the user cannot provision Windows Hello for Business except on Azure Active Directory joined mobile phones if otherwise required. If set to Not Configured, Intune will not override client defaults.`,
					},
					"unlock_with_biometrics_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: `Controls the use of biometric gestures, such as face and fingerprint, as an alternative to the Windows Hello for Business PIN.  If set to False, biometric gestures are not allowed. Users must still configure a PIN as a backup in case of failures.`,
					},
				},
				Validators: []validator.Object{
					deviceEnrollmentConfigurationDeviceEnrollmentConfigurationValidator,
				},
				MarkdownDescription: `Windows Hello for Business settings lets users access their devices using a gesture, such as biometric authentication, or a PIN. Configure settings for enrolled Windows 10, Windows 10 Mobile and later.`,
			},
		},
		"windows10_esp": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windows10EnrollmentCompletionPageConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windows10EnrollmentCompletionPageConfiguration
					"allow_device_reset_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, allows device reset on installation failure. When false, reset is blocked. The default is false.`,
					},
					"allow_device_use_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, allows the user to continue using the device on installation failure. When false, blocks the user on installation failure. The default is false.`,
					},
					"allow_log_collection_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: `When TRUE, allows log collection on installation failure. When false, log collection is not allowed. The default is false.`,
					},
					"block_device_setup_retry_by_user": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, blocks user from retrying the setup on installation failure. When false, user is allowed to retry. The default is false.`,
					},
					"custom_error_message": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("Setup could not be completed. Please try again or contact your support person for help."),
						},
						Computed:            true,
						MarkdownDescription: `The custom error message to show upon installation failure. Max length is 10000. example: "Setup could not be completed. Please try again or contact your support person for help."`,
					},
					"disable_user_status_tracking_after_first_user": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: `When TRUE, disables showing installation progress for first user post enrollment. When false, enables showing progress. The default is false.`,
					},
					"install_progress_timeout_in_minutes": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(60)},
						Computed:            true,
						MarkdownDescription: `The installation progress timeout in minutes. Default is 60 minutes.`,
					},
					"selected_mobile_app_ids": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Selected applications to track the installation status. It is in the form of an array of GUIDs.`,
					},
					"show_installation_progress": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: `When TRUE, shows installation progress to user. When false, hides installation progress. The default is false.`,
					},
					"track_install_progress_for_autopilot_only": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: `When TRUE, installation progress is tracked for only Autopilot enrollment scenarios. When false, other scenarios are tracked as well. The default is false.`,
					},
				},
				Validators: []validator.Object{
					deviceEnrollmentConfigurationDeviceEnrollmentConfigurationValidator,
				},
				MarkdownDescription: `Windows 10 Enrollment Status Page Configuration`,
			},
		},
	},
	MarkdownDescription: `The Base Class of Device Enrollment Configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-deviceEnrollmentConfiguration?view=graph-rest-beta`,
}

var deviceEnrollmentConfigurationDeviceEnrollmentConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("device_enrollment_limit"),
	path.MatchRelative().AtParent().AtName("platform_restrictions"),
	path.MatchRelative().AtParent().AtName("single_platform_restriction"),
	path.MatchRelative().AtParent().AtName("windows_hello_for_business"),
	path.MatchRelative().AtParent().AtName("windows10_esp"),
)

var deviceEnrollmentConfigurationDeviceEnrollmentPlatformRestrictionAttributes = map[string]schema.Attribute{ // deviceEnrollmentPlatformRestriction
	"blocked_manufacturers": schema.SetAttribute{
		ElementType:         types.StringType,
		Optional:            true,
		PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
		Computed:            true,
		MarkdownDescription: `Collection of blocked Manufacturers.`,
	},
	"blocked_skus": schema.SetAttribute{
		ElementType:         types.StringType,
		Optional:            true,
		PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
		Computed:            true,
		MarkdownDescription: `Collection of blocked Skus.`,
	},
	"os_maximum_version": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			wpvalidator.TranslateGraphNullWithTfDefault(tftypes.String),
		},
		PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
		Computed:            true,
		MarkdownDescription: `Max OS version supported`,
	},
	"os_minimum_version": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			wpvalidator.TranslateGraphNullWithTfDefault(tftypes.String),
		},
		PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
		Computed:            true,
		MarkdownDescription: `Min OS version supported`,
	},
	"personal_device_enrollment_blocked": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		MarkdownDescription: `Block personally owned devices from enrolling`,
	},
	"platform_blocked": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		MarkdownDescription: `Block the platform from enrolling`,
	},
}
