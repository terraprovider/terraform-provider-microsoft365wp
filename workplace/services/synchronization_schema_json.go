package services

import (
	"strings"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpjsontypes"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	SynchronizationSchemaJsonResource = generic.GenericResource{
		TypeNameSuffix: "synchronization_schema_json",
		SpecificSchema: synchronizationSchemaJsonResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/servicePrincipals",
			ParentEntities: generic.ParentEntities{
				{
					ParentIdField: path.Root("service_principal_id"),
					UriSuffix:     "synchronization/jobs",
				},
			},
			UriSuffix: "schema",
			ReadOptions: generic.ReadOptions{
				DataSource: generic.DataSourceOptions{
					Plural: generic.PluralOptions{
						NoDataSource: true,
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				UpdateInsteadOfCreate: true,
				UsePutForUpdate:       true,
				// http DELETE is valid here (will reset synchronizationSchema to default values)
			},
		},
	}

	SynchronizationSchemaJsonSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&SynchronizationSchemaJsonResource)
)

var synchronizationSchemaJsonResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // synchronizationSchema
		"service_principal_id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Provider Note: ID of the service principal that this synchronization schema belongs to. Required.",
		},
		"id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Unique identifier for the schema.",
		},
		"synchronization_rules_json": schema.StringAttribute{
			Required:            true,
			CustomType:          wpjsontypes.NormalizedType{},
			Description:         `synchronizationRules`, // custom MS Graph attribute name
			MarkdownDescription: "A collection of synchronization rules configured for the [synchronizationJob](synchronization-synchronizationjob.md) or [synchronizationTemplate](synchronization-synchronizationtemplate.md). / Defines how the synchronization should be performed for the synchronization engine, including which objects to synchronize and in which direction, how objects from the source directory should be matched with objects in the target directory, and how attributes should be transformed when they're synchronized from the source to the target directory.\n\nSynchronization rules are updated as part of the [synchronization schema](synchronization-synchronizationschema.md). / https://learn.microsoft.com/en-us/graph/api/resources/synchronization-synchronizationrule?view=graph-rest-beta",
		},
		"version": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The version of the schema, updated automatically with every schema change.",
		},
		"directories_json": schema.StringAttribute{
			Required:            true,
			CustomType:          wpjsontypes.NormalizedType{WpObjectFilterFunc: synchronizationSchemaJsonDirectoriesFilter},
			Description:         `directories`, // custom MS Graph attribute name
			MarkdownDescription: "Contains the collection of directories and all of their objects. / Provides the synchronization engine information about a directory and its objects. This resource tells the synchronization engine, for example, that the directory has objects named **user** and **group**, which attributes are supported for those objects, and the types for those attributes. In order for the object and attribute to participate in [synchronization rules](synchronization-synchronizationrule.md) and [object mappings](synchronization-objectmapping.md), they must be defined as part of the directory definition.\n\nIn general, the default [synchronization schema](synchronization-synchronizationschema.md) provided as part of the [synchronization template](synchronization-synchronizationtemplate.md) defines the most commonly used objects and attributes for that directory. However, if the directory supports the addition of custom attributes, you might want to expand the default definition with your own custom objects or attributes. For more information, see the following articles.\n\nDirectory definitions are updated as part of the [synchronization schema](synchronization-synchronizationschema.md). / https://learn.microsoft.com/en-us/graph/api/resources/synchronization-directorydefinition?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "Defines what objects are synchronized and how they are synchronized. The synchronization schema contains most of the setup information for a particular synchronization job. Typically, you customize some of the [attribute mappings](synchronization-attributemapping.md), or add a [scoping filter](synchronization-filter.md) to synchronize only objects that satisfy a certain condition.\n\nThe following sections describe the high-level components of the synchronization schema. / https://learn.microsoft.com/en-us/graph/api/resources/synchronization-synchronizationschema?view=graph-rest-beta\n\nProvider Note: To import this resource, an ID consisting of `service_principal_id` and `id` being joined by a forward slash (`/`) must be used. ||| MS Graph: Synchronization",
}

func synchronizationSchemaJsonDirectoriesFilter(path string, value any) (pass bool, err error) {
	pass = !strings.HasSuffix(path, "/version")
	return
}
