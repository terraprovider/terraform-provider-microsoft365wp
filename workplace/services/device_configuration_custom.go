package services

import (
	"encoding/base64"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/manicminer/hamilton/msgraph"
)

var (
	DeviceConfigurationCustomResource = generic.GenericResource{
		TypeNameSuffix: "device_configuration_custom",
		SpecificSchema: deviceConfigurationCustomResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/deviceManagement/deviceConfigurations",
			ReadOptions:                deviceConfigurationCustomReadOptions,
			WriteSubActions:            deviceConfigurationCustomWriteSubActions,
			TerraformToGraphMiddleware: deviceConfigurationCustomTerraformToGraphMiddleware,
		},
	}

	DeviceConfigurationCustomSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceConfigurationCustomResource)

	DeviceConfigurationCustomPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&DeviceConfigurationCustomSingularDataSource, "device_configurations_custom")
)

var deviceConfigurationCustomReadOptions = generic.ReadOptions{
	ODataFilter: "isof('microsoft.graph.windows10CustomConfiguration')",
	ExtraRequests: []generic.ReadExtraRequest{
		{
			ParentAttribute: "assignments",
			UriSuffix:       "assignments",
		},
	},
	ExtraRequestsCustom: []generic.ReadExtraRequestCustom{
		deviceConfigurationCustomExtraRequestCustomReadSecretValue,
	},
}

var deviceConfigurationCustomWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
		},
	},
}

func deviceConfigurationCustomTerraformToGraphMiddleware(params generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {

	omaSettings := params.RawVal["omaSettings"].([]any)
	if omaSettings == nil {
		return nil
	}

	for i, omaSettingRaw := range omaSettings {
		omaSetting, ok := omaSettingRaw.(map[string]any)
		if !ok {
			return fmt.Errorf("omaSettings[%d] not of type map[string]any", i)
		}
		odataType, ok := omaSetting["@odata.type"].(string)
		if !ok {
			return fmt.Errorf("omaSettings[%d]: @odata.type not found or not of type string", i)
		}
		if odataType == "#microsoft.graph.omaSettingStringXml" {
			value, ok := omaSetting["value"].(string)
			if !ok {
				return fmt.Errorf("omaSettings[%d]: value not found or not of type string", i)
			}
			omaSetting["value"] = base64.StdEncoding.EncodeToString([]byte(value)) // []byte() will do UTF-8 internally
		}
	}

	return nil
}

func deviceConfigurationCustomExtraRequestCustomReadSecretValue(params generic.ReadExtraRequestCustomFuncParams) {
	warningDetail := "Skipping retrieval of secret value(s)"

	omaSettings := params.RawVal["omaSettings"].([]any)
	if omaSettings == nil {
		params.Diags.AddWarning("omaSettings not found or not an array", warningDetail)
		return
	}

	for i, omaSettingRaw := range omaSettings {
		omaSetting, ok := omaSettingRaw.(map[string]any)
		if !ok {
			params.Diags.AddWarning(fmt.Sprintf("omaSettings[%d] not of type map[string]any", i), warningDetail)
			continue
		}
		isEncrypted, ok := omaSetting["isEncrypted"].(bool)
		if !ok {
			params.Diags.AddWarning(fmt.Sprintf("omaSettings[%d]: isEncrypted not found or not of type bool", i), warningDetail)
			continue
		}
		if isEncrypted {
			secretReferenceValueId, ok := omaSetting["secretReferenceValueId"].(string)
			if !ok {
				params.Diags.AddWarning(fmt.Sprintf("omaSettings[%d]: secretReferenceValueId not found or not of type string", i), warningDetail)
				continue
			}

			getSecretUri := msgraph.Uri{
				Entity: fmt.Sprintf("%s/getOmaSettingPlainTextValue(secretReferenceValueId='%s')", params.Uri.Entity, secretReferenceValueId),
			}
			plainValueResponse := generic.ReadRaw(params.Ctx, params.Diags, params.Client, getSecretUri, "", "", nil, params.TolerateNotFound)
			if params.Diags.HasError() {
				return
			}

			plainValue, ok := plainValueResponse["value"].(string)
			if !ok {
				params.Diags.AddWarning(fmt.Sprintf("omaSettings[%d] getOmaSettingPlainTextValue response: value not found or not of type string", i), warningDetail)
				continue
			}

			omaSetting["value"] = plainValue
		}
	}
}

var deviceConfigurationCustomResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceConfiguration
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "DateTime the object was created.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Admin provided description of the Device Configuration.",
		},
		"device_management_applicability_rule_device_mode": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementApplicabilityRuleDeviceMode
				"device_mode": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("standardConfiguration", "sModeConfiguration"),
					},
					MarkdownDescription: "Applicability rule for device mode. / Windows 10 Device Mode type; possible values are: `standardConfiguration` (Standard Configuration), `sModeConfiguration` (S Mode Configuration)",
				},
				"name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Name for object.",
				},
				"rule_type": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
					MarkdownDescription: "Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)",
				},
			},
			MarkdownDescription: "The device mode applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceManagementApplicabilityRuleDeviceMode?view=graph-rest-beta",
		},
		"device_management_applicability_rule_os_edition": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementApplicabilityRuleOsEdition
				"name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Name for object.",
				},
				"os_edition_types": schema.SetAttribute{
					ElementType: types.StringType,
					Required:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("windows10Enterprise", "windows10EnterpriseN", "windows10Education", "windows10EducationN", "windows10MobileEnterprise", "windows10HolographicEnterprise", "windows10Professional", "windows10ProfessionalN", "windows10ProfessionalEducation", "windows10ProfessionalEducationN", "windows10ProfessionalWorkstation", "windows10ProfessionalWorkstationN", "notConfigured", "windows10Home", "windows10HomeChina", "windows10HomeN", "windows10HomeSingleLanguage", "windows10Mobile", "windows10IoTCore", "windows10IoTCoreCommercial"),
						),
					},
					MarkdownDescription: "Applicability rule OS edition type. / Windows 10 Edition type; possible values are: `windows10Enterprise` (Windows 10 Enterprise), `windows10EnterpriseN` (Windows 10 EnterpriseN), `windows10Education` (Windows 10 Education), `windows10EducationN` (Windows 10 EducationN), `windows10MobileEnterprise` (Windows 10 Mobile Enterprise), `windows10HolographicEnterprise` (Windows 10 Holographic Enterprise), `windows10Professional` (Windows 10 Professional), `windows10ProfessionalN` (Windows 10 ProfessionalN), `windows10ProfessionalEducation` (Windows 10 Professional Education), `windows10ProfessionalEducationN` (Windows 10 Professional EducationN), `windows10ProfessionalWorkstation` (Windows 10 Professional for Workstations), `windows10ProfessionalWorkstationN` (Windows 10 Professional for Workstations N), `notConfigured` (NotConfigured), `windows10Home` (Windows 10 Home), `windows10HomeChina` (Windows 10 Home China), `windows10HomeN` (Windows 10 Home N), `windows10HomeSingleLanguage` (Windows 10 Home Single Language), `windows10Mobile` (Windows 10 Mobile), `windows10IoTCore` (Windows 10 IoT Core), `windows10IoTCoreCommercial` (Windows 10 IoT Core Commercial)",
				},
				"rule_type": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
					MarkdownDescription: "Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)",
				},
			},
			MarkdownDescription: "The OS edition applicability for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceManagementApplicabilityRuleOsEdition?view=graph-rest-beta",
		},
		"device_management_applicability_rule_os_version": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementApplicabilityRuleOsVersion
				"max_os_version": schema.StringAttribute{
					Optional:            true,
					Description:         `maxOSVersion`, // custom MS Graph attribute name
					MarkdownDescription: "Max OS version for Applicability Rule.",
				},
				"min_os_version": schema.StringAttribute{
					Optional:            true,
					Description:         `minOSVersion`, // custom MS Graph attribute name
					MarkdownDescription: "Min OS version for Applicability Rule.",
				},
				"name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Name for object.",
				},
				"rule_type": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
					MarkdownDescription: "Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)",
				},
			},
			MarkdownDescription: "The OS version applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceManagementApplicabilityRuleOsVersion?view=graph-rest-beta",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Admin provided name of the device configuration.",
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
			MarkdownDescription: "List of Scope Tags for this Entity instance.",
		},
		"supports_scope_tags": schema.BoolAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
		},
		"version": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "Version of the device configuration.",
		},
		"assignments": deviceAndAppManagementAssignment,
		"windows10": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windows10CustomConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windows10CustomConfiguration
					"oma_settings": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // omaSetting
								"description": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Description.",
								},
								"display_name": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Display Name.",
								},
								"oma_uri": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "OMA.",
								},
								"base64": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.omaSettingBase64",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // omaSettingBase64
											"file_name": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "File name associated with the Value property (*.cer | *.crt | *.p7b | *.bin).",
											},
											"value_base64": schema.StringAttribute{
												Required:            true,
												Description:         `value`, // custom MS Graph attribute name
												MarkdownDescription: "Value. (Base64 encoded string)",
											},
										},
										Validators:          []validator.Object{deviceConfigurationCustomOmaSettingValidator},
										MarkdownDescription: "OMA Settings Base64 definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingBase64?view=graph-rest-beta",
									},
								},
								"boolean": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.omaSettingBoolean",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // omaSettingBoolean
											"value": schema.BoolAttribute{
												Required:            true,
												MarkdownDescription: "Value.",
											},
										},
										Validators:          []validator.Object{deviceConfigurationCustomOmaSettingValidator},
										MarkdownDescription: "OMA Settings Boolean definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingBoolean?view=graph-rest-beta",
									},
								},
								"date_time": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.omaSettingDateTime",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // omaSettingDateTime
											"value": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Value.",
											},
										},
										Validators:          []validator.Object{deviceConfigurationCustomOmaSettingValidator},
										MarkdownDescription: "OMA Settings DateTime definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingDateTime?view=graph-rest-beta",
									},
								},
								"floating_point": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.omaSettingFloatingPoint",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // omaSettingFloatingPoint
											"value": schema.NumberAttribute{
												Required:            true,
												MarkdownDescription: "Value.",
											},
										},
										Validators:          []validator.Object{deviceConfigurationCustomOmaSettingValidator},
										MarkdownDescription: "OMA Settings Floating Point definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingFloatingPoint?view=graph-rest-beta",
									},
								},
								"integer": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.omaSettingInteger",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // omaSettingInteger
											"value": schema.Int64Attribute{
												Required:            true,
												MarkdownDescription: "Value.",
											},
										},
										Validators:          []validator.Object{deviceConfigurationCustomOmaSettingValidator},
										MarkdownDescription: "OMA Settings Integer definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingInteger?view=graph-rest-beta",
									},
								},
								"string": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.omaSettingString",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // omaSettingString
											"value": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Value.",
											},
										},
										Validators:          []validator.Object{deviceConfigurationCustomOmaSettingValidator},
										MarkdownDescription: "OMA Settings String definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingString?view=graph-rest-beta",
									},
								},
								"string_xml": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.omaSettingStringXml",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // omaSettingStringXml
											"file_name": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "File name associated with the Value property (*.xml).",
											},
											"value": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Value. (UTF8 encoded byte array)",
											},
										},
										Validators:          []validator.Object{deviceConfigurationCustomOmaSettingValidator},
										MarkdownDescription: "OMA Settings StringXML definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingStringXml?view=graph-rest-beta",
									},
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "OMA settings. This collection can contain a maximum of 1000 elements. / OMA Settings definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSetting?view=graph-rest-beta",
					},
				},
				Validators: []validator.Object{
					deviceConfigurationCustomDeviceConfigurationValidator,
				},
				MarkdownDescription: "This topic provides descriptions of the declared methods, properties and relationships exposed by the windows10CustomConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10CustomConfiguration?view=graph-rest-beta",
			},
		},
	},
	MarkdownDescription: "Device Configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceConfiguration?view=graph-rest-beta",
}

var deviceConfigurationCustomDeviceConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("windows10"),
)

var deviceConfigurationCustomOmaSettingValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("base64"),
	path.MatchRelative().AtParent().AtName("boolean"),
	path.MatchRelative().AtParent().AtName("date_time"),
	path.MatchRelative().AtParent().AtName("floating_point"),
	path.MatchRelative().AtParent().AtName("integer"),
	path.MatchRelative().AtParent().AtName("string"),
	path.MatchRelative().AtParent().AtName("string_xml"),
)
