package wpvalidator

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExactlyOneOfSiblingsValidateObject(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		schemaSiblingNames []string
		validatorSiblings  []string
		currentName        string
		configValue        types.Object
		config             tfsdk.Config
		expectError        bool
		expectDetail       string
	}{
		"succeeds when exactly one sibling is set": {
			schemaSiblingNames: []string{"first", "second"},
			validatorSiblings:  []string{"first", "second"},
			currentName:        "first",
			configValue:        testSiblingObjectValue("one"),
			config: testSiblingConfig(t, []string{"first", "second"}, map[string]types.Object{
				"first":  testSiblingObjectValue("one"),
				"second": testSiblingObjectNull(),
			}),
		},
		"errors when no siblings are set": {
			schemaSiblingNames: []string{"first", "second"},
			validatorSiblings:  []string{"first", "second"},
			currentName:        "first",
			configValue:        testSiblingObjectNull(),
			config: testSiblingConfig(t, []string{"first", "second"}, map[string]types.Object{
				"first":  testSiblingObjectNull(),
				"second": testSiblingObjectNull(),
			}),
			expectError:  true,
			expectDetail: "Exactly one of `first`, `second` must be set.",
		},
		"errors when multiple siblings are set": {
			schemaSiblingNames: []string{"first", "second"},
			validatorSiblings:  []string{"first", "second"},
			currentName:        "first",
			configValue:        testSiblingObjectValue("one"),
			config: testSiblingConfig(t, []string{"first", "second"}, map[string]types.Object{
				"first":  testSiblingObjectValue("one"),
				"second": testSiblingObjectValue("two"),
			}),
			expectError:  true,
			expectDetail: "Exactly one of `first`, `second` must be set; 2 were provided.",
		},
		"ignores sibling names that are not present in the parent object": {
			schemaSiblingNames: []string{"first", "second"},
			validatorSiblings:  []string{"first", "second", "third"},
			currentName:        "first",
			configValue:        testSiblingObjectValue("one"),
			config: testSiblingConfig(t, []string{"first", "second"}, map[string]types.Object{
				"first":  testSiblingObjectValue("one"),
				"second": testSiblingObjectNull(),
			}),
		},
		"delays validation when current attribute is unknown": {
			schemaSiblingNames: []string{"first", "second"},
			validatorSiblings:  []string{"first", "second"},
			currentName:        "first",
			configValue:        testSiblingObjectUnknown(),
			config: testSiblingConfig(t, []string{"first", "second"}, map[string]types.Object{
				"first":  testSiblingObjectUnknown(),
				"second": testSiblingObjectNull(),
			}),
		},
		"delays validation when another sibling is unknown": {
			schemaSiblingNames: []string{"first", "second"},
			validatorSiblings:  []string{"first", "second"},
			currentName:        "first",
			configValue:        testSiblingObjectNull(),
			config: testSiblingConfig(t, []string{"first", "second"}, map[string]types.Object{
				"first":  testSiblingObjectNull(),
				"second": testSiblingObjectUnknown(),
			}),
		},
		"delays validation when the parent object is null": {
			schemaSiblingNames: []string{"first", "second"},
			validatorSiblings:  []string{"first", "second"},
			currentName:        "first",
			configValue:        testSiblingObjectNull(),
			config:             testNullParentSiblingConfig(t, []string{"first", "second"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			validator := ExactlyOneOfSiblings(tc.validatorSiblings...)
			req := frameworkvalidator.ObjectRequest{
				Path:        path.Root("container").AtName(tc.currentName),
				Config:      tc.config,
				ConfigValue: tc.configValue,
			}
			resp := &frameworkvalidator.ObjectResponse{}

			validator.ValidateObject(context.Background(), req, resp)

			if tc.expectError {
				if !resp.Diagnostics.HasError() {
					t.Fatalf("expected validation error, got diagnostics: %#v", resp.Diagnostics)
				}

				errs := resp.Diagnostics.Errors()
				if len(errs) != 1 {
					t.Fatalf("expected exactly 1 error diagnostic, got %d", len(errs))
				}

				if errs[0].Summary() != "Invalid Attribute Combination" {
					t.Fatalf("unexpected diagnostic summary: %q", errs[0].Summary())
				}

				if !strings.Contains(errs[0].Detail(), tc.expectDetail) {
					t.Fatalf("unexpected diagnostic detail\nexpected to contain: %q\nactual: %q", tc.expectDetail, errs[0].Detail())
				}

				withPath, ok := errs[0].(diag.DiagnosticWithPath)
				if !ok {
					t.Fatalf("expected diagnostic with path, got %T", errs[0])
				}

				expectedPath := path.Root("container").AtName(tc.currentName)
				if !withPath.Path().Equal(expectedPath) {
					t.Fatalf("unexpected diagnostic path\nexpected: %s\nactual:   %s", expectedPath, withPath.Path())
				}

				return
			}

			if len(resp.Diagnostics) != 0 {
				t.Fatalf("expected no diagnostics, got %#v", resp.Diagnostics)
			}
		})
	}
}

var testSiblingChildAttrTypes = map[string]attr.Type{
	"value": types.StringType,
}

func testSiblingObjectValue(value string) types.Object {
	return types.ObjectValueMust(testSiblingChildAttrTypes, map[string]attr.Value{
		"value": types.StringValue(value),
	})
}

func testSiblingObjectNull() types.Object {
	return types.ObjectNull(testSiblingChildAttrTypes)
}

func testSiblingObjectUnknown() types.Object {
	return types.ObjectUnknown(testSiblingChildAttrTypes)
}

func testSiblingConfig(t *testing.T, siblingNames []string, siblingValues map[string]types.Object) tfsdk.Config {
	t.Helper()

	schemaAttributes := make(map[string]rschema.Attribute, len(siblingNames))
	containerAttrTypes := make(map[string]attr.Type, len(siblingNames))
	containerValues := make(map[string]attr.Value, len(siblingNames))

	for _, name := range siblingNames {
		schemaAttributes[name] = rschema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]rschema.Attribute{
				"value": rschema.StringAttribute{Optional: true},
			},
		}
		containerAttrTypes[name] = types.ObjectType{AttrTypes: testSiblingChildAttrTypes}
		if value, ok := siblingValues[name]; ok {
			containerValues[name] = value
			continue
		}
		containerValues[name] = testSiblingObjectNull()
	}

	containerValue := types.ObjectValueMust(containerAttrTypes, containerValues)
	rootValue := types.ObjectValueMust(map[string]attr.Type{
		"container": types.ObjectType{AttrTypes: containerAttrTypes},
	}, map[string]attr.Value{
		"container": containerValue,
	})

	raw, err := rootValue.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("failed to convert config value to terraform value: %v", err)
	}

	return tfsdk.Config{
		Raw: raw,
		Schema: rschema.Schema{
			Attributes: map[string]rschema.Attribute{
				"container": rschema.SingleNestedAttribute{
					Optional:   true,
					Attributes: schemaAttributes,
				},
			},
		},
	}
}

func testNullParentSiblingConfig(t *testing.T, siblingNames []string) tfsdk.Config {
	t.Helper()

	schemaAttributes := make(map[string]rschema.Attribute, len(siblingNames))
	containerAttrTypes := make(map[string]attr.Type, len(siblingNames))

	for _, name := range siblingNames {
		schemaAttributes[name] = rschema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]rschema.Attribute{
				"value": rschema.StringAttribute{Optional: true},
			},
		}
		containerAttrTypes[name] = types.ObjectType{AttrTypes: testSiblingChildAttrTypes}
	}

	rootValue := types.ObjectValueMust(map[string]attr.Type{
		"container": types.ObjectType{AttrTypes: containerAttrTypes},
	}, map[string]attr.Value{
		"container": types.ObjectNull(containerAttrTypes),
	})

	raw, err := rootValue.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("failed to convert null parent config to terraform value: %v", err)
	}

	return tfsdk.Config{
		Raw: raw,
		Schema: rschema.Schema{
			Attributes: map[string]rschema.Attribute{
				"container": rschema.SingleNestedAttribute{
					Optional:   true,
					Attributes: schemaAttributes,
				},
			},
		},
	}
}

