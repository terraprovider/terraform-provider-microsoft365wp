package services

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	DeviceManagementIntentResource = generic.GenericResource{
		TypeNameSuffix: "device_management_intent",
		SpecificSchema: deviceManagementIntentResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/deviceManagement/Intents",
			ReadOptions:                deviceManagementIntentReadOptions,
			WriteSubActions:            deviceManagementIntentWriteSubActions,
			TerraformToGraphMiddleware: deviceManagementIntentTerraformToGraphMiddleware,
			CreateModifyFunc:           deviceManagementIntentCreateModifyFunc,
		},
	}

	DeviceManagementIntentSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceManagementIntentResource)

	DeviceManagementIntentPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&DeviceManagementIntentSingularDataSource, "")
)

var deviceManagementIntentReadOptions = generic.ReadOptions{
	ODataExpand: "assignments,settings",
}

var deviceManagementIntentWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
		},
	},
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"settings"},
			UriSuffix:  "updateSettings",
		},
	},
}

func deviceManagementIntentTerraformToGraphMiddleware(ctx context.Context, params generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if params.IsUpdate {
		// templateId cannot be updated and may not even be written again after creation
		delete(params.RawVal, "templateId")
	}
	return nil
}

func deviceManagementIntentCreateModifyFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.CreateModifyFuncParams) {
	// entity needs to be created using `createInstance` action on template
	// https://learn.microsoft.com/en-us/graph/api/intune-deviceintent-devicemanagementtemplate-createinstance?view=graph-rest-beta
	// it does not seem to be required to rename `settings` to `settingsDelta` though
	var templateId string
	diags.Append(params.Req.Plan.GetAttribute(ctx, path.Root("template_id"), &templateId)...)
	if diags.HasError() {
		return
	}
	params.BaseUri = fmt.Sprintf("/deviceManagement/templates/%s/createInstance", templateId)
}

var deviceManagementIntentResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceManagementIntent
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The intent ID",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The user given description",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The user given display name",
		},
		"is_assigned": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "Signifies whether or not the intent is assigned to users",
		},
		"is_migrating_to_configuration_policy": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "Signifies whether or not the intent is being migrated to the configurationPolicies endpoint",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "When the intent was last modified",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tags for this Entity instance. The _provider_ default value is `[\"0\"]`.",
		},
		"template_id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The ID of the template this intent was created from (if any)",
		},
		"assignments": deviceAndAppManagementAssignment,
		"settings": schema.SetNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // deviceManagementSettingInstance
					"definition_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The ID of the setting definition for this instance",
					},
					"abstract_complex": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.deviceManagementAbstractComplexSettingInstance",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // deviceManagementAbstractComplexSettingInstance
								"implementation_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The definition ID for the chosen implementation of this complex setting",
								},
								"value": schema.SetNestedAttribute{
									Required: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // deviceManagementSettingInstance
											"definition_id": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The ID of the setting definition for this instance",
											},
											"abstract_complex": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementAbstractComplexSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementAbstractComplexSettingInstance
														"implementation_id": schema.StringAttribute{
															Optional:            true,
															MarkdownDescription: "The definition ID for the chosen implementation of this complex setting",
														},
														"value": schema.SetNestedAttribute{
															Required: true,
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{ // deviceManagementSettingInstance
																	"definition_id": schema.StringAttribute{
																		Required:            true,
																		MarkdownDescription: "The ID of the setting definition for this instance",
																	},
																	"boolean": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementBooleanSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementBooleanSettingInstance
																				"value": schema.BoolAttribute{
																					Required:            true,
																					MarkdownDescription: "The boolean value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing a boolean value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementbooleansettinginstance?view=graph-rest-beta",
																		},
																	},
																	"integer": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementIntegerSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementIntegerSettingInstance
																				"value": schema.Int64Attribute{
																					Required:            true,
																					MarkdownDescription: "The integer value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing an integer value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementintegersettinginstance?view=graph-rest-beta",
																		},
																	},
																	"string": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementStringSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementStringSettingInstance
																				"value": schema.StringAttribute{
																					Required:            true,
																					MarkdownDescription: "The string value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing a string value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementstringsettinginstance?view=graph-rest-beta",
																		},
																	},
																},
															},
															MarkdownDescription: "The values that make up the complex setting",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a complex value for an abstract setting / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementabstractcomplexsettinginstance?view=graph-rest-beta",
												},
											},
											"boolean": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementBooleanSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementBooleanSettingInstance
														"value": schema.BoolAttribute{
															Required:            true,
															MarkdownDescription: "The boolean value",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a boolean value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementbooleansettinginstance?view=graph-rest-beta",
												},
											},
											"collection": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementCollectionSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementCollectionSettingInstance
														"value_json": schema.StringAttribute{
															Required:            true,
															MarkdownDescription: "JSON representation of the value",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a collection of values / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementcollectionsettinginstance?view=graph-rest-beta",
												},
											},
											"complex": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementComplexSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementComplexSettingInstance
														"value": schema.SetNestedAttribute{
															Required: true,
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{ // deviceManagementSettingInstance
																	"definition_id": schema.StringAttribute{
																		Required:            true,
																		MarkdownDescription: "The ID of the setting definition for this instance",
																	},
																	"boolean": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementBooleanSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementBooleanSettingInstance
																				"value": schema.BoolAttribute{
																					Required:            true,
																					MarkdownDescription: "The boolean value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing a boolean value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementbooleansettinginstance?view=graph-rest-beta",
																		},
																	},
																	"integer": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementIntegerSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementIntegerSettingInstance
																				"value": schema.Int64Attribute{
																					Required:            true,
																					MarkdownDescription: "The integer value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing an integer value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementintegersettinginstance?view=graph-rest-beta",
																		},
																	},
																	"string": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementStringSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementStringSettingInstance
																				"value": schema.StringAttribute{
																					Required:            true,
																					MarkdownDescription: "The string value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing a string value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementstringsettinginstance?view=graph-rest-beta",
																		},
																	},
																},
															},
															MarkdownDescription: "The values that make up the complex setting",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a complex value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementcomplexsettinginstance?view=graph-rest-beta",
												},
											},
											"integer": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementIntegerSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementIntegerSettingInstance
														"value": schema.Int64Attribute{
															Required:            true,
															MarkdownDescription: "The integer value",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing an integer value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementintegersettinginstance?view=graph-rest-beta",
												},
											},
											"string": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementStringSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementStringSettingInstance
														"value": schema.StringAttribute{
															Required:            true,
															MarkdownDescription: "The string value",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a string value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementstringsettinginstance?view=graph-rest-beta",
												},
											},
										},
									},
									MarkdownDescription: "The values that make up the complex setting",
								},
							},
							Validators: []validator.Object{
								deviceManagementIntentDeviceManagementSettingInstanceValidator,
							},
							MarkdownDescription: "A setting instance representing a complex value for an abstract setting / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementabstractcomplexsettinginstance?view=graph-rest-beta",
						},
					},
					"boolean": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.deviceManagementBooleanSettingInstance",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // deviceManagementBooleanSettingInstance
								"value": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: "The boolean value",
								},
							},
							Validators: []validator.Object{
								deviceManagementIntentDeviceManagementSettingInstanceValidator,
							},
							MarkdownDescription: "A setting instance representing a boolean value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementbooleansettinginstance?view=graph-rest-beta",
						},
					},
					"collection": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.deviceManagementCollectionSettingInstance",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // deviceManagementCollectionSettingInstance
								"value_json": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "JSON representation of the value",
								},
							},
							Validators: []validator.Object{
								deviceManagementIntentDeviceManagementSettingInstanceValidator,
							},
							MarkdownDescription: "A setting instance representing a collection of values / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementcollectionsettinginstance?view=graph-rest-beta",
						},
					},
					"complex": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.deviceManagementComplexSettingInstance",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // deviceManagementComplexSettingInstance
								"value": schema.SetNestedAttribute{
									Required: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // deviceManagementSettingInstance
											"definition_id": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The ID of the setting definition for this instance",
											},
											"abstract_complex": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementAbstractComplexSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementAbstractComplexSettingInstance
														"implementation_id": schema.StringAttribute{
															Optional:            true,
															MarkdownDescription: "The definition ID for the chosen implementation of this complex setting",
														},
														"value": schema.SetNestedAttribute{
															Required: true,
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{ // deviceManagementSettingInstance
																	"definition_id": schema.StringAttribute{
																		Required:            true,
																		MarkdownDescription: "The ID of the setting definition for this instance",
																	},
																	"boolean": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementBooleanSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementBooleanSettingInstance
																				"value": schema.BoolAttribute{
																					Required:            true,
																					MarkdownDescription: "The boolean value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing a boolean value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementbooleansettinginstance?view=graph-rest-beta",
																		},
																	},
																	"integer": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementIntegerSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementIntegerSettingInstance
																				"value": schema.Int64Attribute{
																					Required:            true,
																					MarkdownDescription: "The integer value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing an integer value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementintegersettinginstance?view=graph-rest-beta",
																		},
																	},
																	"string": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementStringSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementStringSettingInstance
																				"value": schema.StringAttribute{
																					Required:            true,
																					MarkdownDescription: "The string value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing a string value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementstringsettinginstance?view=graph-rest-beta",
																		},
																	},
																},
															},
															MarkdownDescription: "The values that make up the complex setting",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a complex value for an abstract setting / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementabstractcomplexsettinginstance?view=graph-rest-beta",
												},
											},
											"boolean": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementBooleanSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementBooleanSettingInstance
														"value": schema.BoolAttribute{
															Required:            true,
															MarkdownDescription: "The boolean value",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a boolean value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementbooleansettinginstance?view=graph-rest-beta",
												},
											},
											"collection": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementCollectionSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementCollectionSettingInstance
														"value_json": schema.StringAttribute{
															Required:            true,
															MarkdownDescription: "JSON representation of the value",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a collection of values / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementcollectionsettinginstance?view=graph-rest-beta",
												},
											},
											"complex": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementComplexSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementComplexSettingInstance
														"value": schema.SetNestedAttribute{
															Required: true,
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{ // deviceManagementSettingInstance
																	"definition_id": schema.StringAttribute{
																		Required:            true,
																		MarkdownDescription: "The ID of the setting definition for this instance",
																	},
																	"boolean": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementBooleanSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementBooleanSettingInstance
																				"value": schema.BoolAttribute{
																					Required:            true,
																					MarkdownDescription: "The boolean value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing a boolean value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementbooleansettinginstance?view=graph-rest-beta",
																		},
																	},
																	"integer": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementIntegerSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementIntegerSettingInstance
																				"value": schema.Int64Attribute{
																					Required:            true,
																					MarkdownDescription: "The integer value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing an integer value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementintegersettinginstance?view=graph-rest-beta",
																		},
																	},
																	"string": generic.OdataDerivedTypeNestedAttributeRs{
																		DerivedType: "#microsoft.graph.deviceManagementStringSettingInstance",
																		SingleNestedAttribute: schema.SingleNestedAttribute{
																			Optional: true,
																			Attributes: map[string]schema.Attribute{ // deviceManagementStringSettingInstance
																				"value": schema.StringAttribute{
																					Required:            true,
																					MarkdownDescription: "The string value",
																				},
																			},
																			Validators: []validator.Object{
																				deviceManagementIntentDeviceManagementSettingInstanceValidator,
																			},
																			MarkdownDescription: "A setting instance representing a string value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementstringsettinginstance?view=graph-rest-beta",
																		},
																	},
																},
															},
															MarkdownDescription: "The values that make up the complex setting",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a complex value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementcomplexsettinginstance?view=graph-rest-beta",
												},
											},
											"integer": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementIntegerSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementIntegerSettingInstance
														"value": schema.Int64Attribute{
															Required:            true,
															MarkdownDescription: "The integer value",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing an integer value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementintegersettinginstance?view=graph-rest-beta",
												},
											},
											"string": generic.OdataDerivedTypeNestedAttributeRs{
												DerivedType: "#microsoft.graph.deviceManagementStringSettingInstance",
												SingleNestedAttribute: schema.SingleNestedAttribute{
													Optional: true,
													Attributes: map[string]schema.Attribute{ // deviceManagementStringSettingInstance
														"value": schema.StringAttribute{
															Required:            true,
															MarkdownDescription: "The string value",
														},
													},
													Validators: []validator.Object{
														deviceManagementIntentDeviceManagementSettingInstanceValidator,
													},
													MarkdownDescription: "A setting instance representing a string value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementstringsettinginstance?view=graph-rest-beta",
												},
											},
										},
									},
									MarkdownDescription: "The values that make up the complex setting",
								},
							},
							Validators: []validator.Object{
								deviceManagementIntentDeviceManagementSettingInstanceValidator,
							},
							MarkdownDescription: "A setting instance representing a complex value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementcomplexsettinginstance?view=graph-rest-beta",
						},
					},
					"integer": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.deviceManagementIntegerSettingInstance",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // deviceManagementIntegerSettingInstance
								"value": schema.Int64Attribute{
									Required:            true,
									MarkdownDescription: "The integer value",
								},
							},
							Validators: []validator.Object{
								deviceManagementIntentDeviceManagementSettingInstanceValidator,
							},
							MarkdownDescription: "A setting instance representing an integer value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementintegersettinginstance?view=graph-rest-beta",
						},
					},
					"string": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.deviceManagementStringSettingInstance",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // deviceManagementStringSettingInstance
								"value": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The string value",
								},
							},
							Validators: []validator.Object{
								deviceManagementIntentDeviceManagementSettingInstanceValidator,
							},
							MarkdownDescription: "A setting instance representing a string value / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementstringsettinginstance?view=graph-rest-beta",
						},
					},
				},
			},
			MarkdownDescription: "Collection of all settings to be applied / Base type for a setting instance / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementsettinginstance?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "Entity that represents an intent to apply settings to a device / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceintent-devicemanagementintent?view=graph-rest-beta ||| MS Graph: Device management",
}

var deviceManagementIntentDeviceManagementSettingInstanceValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("abstract_complex"),
	path.MatchRelative().AtParent().AtName("boolean"),
	path.MatchRelative().AtParent().AtName("collection"),
	path.MatchRelative().AtParent().AtName("complex"),
	path.MatchRelative().AtParent().AtName("integer"),
	path.MatchRelative().AtParent().AtName("string"),
)
