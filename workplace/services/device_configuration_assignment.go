package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var DeviceConfigurationAssignmentResource = newDeviceConfigurationAssignmentResource()

func newDeviceConfigurationAssignmentResource() *generic.GenericResource {
	r := GetAssignmentChildResource(&DeviceConfigurationResource, true)
	r.SpecificSchema.MarkdownDescription = "Manages one assignment; assignments created elsewhere remain outside this resource. Do not also manage the same policy's complete assignments list through another resource. Changes to the parent policy, target, or intent replace the assignment." + r.SpecificSchema.MarkdownDescription
	r.SpecificSchema.Attributes["intent"] = schema.StringAttribute{
		Required:            true,
		Validators:          []validator.String{stringvalidator.OneOf("apply")},
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		MarkdownDescription: "Must be `apply`. Remove the assignment by destroying this resource. Changes replace the assignment.",
	}
	return &r
}
