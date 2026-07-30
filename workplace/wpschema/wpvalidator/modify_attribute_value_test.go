package wpvalidator

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestModifyValueSortSetImplModifyGraphToTerraform(t *testing.T) {
	t.Parallel()

	validator := ModifyValuesGraphToTerraformReorderByKey("id")

	tests := map[string]struct {
		graphValue    any
		priorState    *tfsdk.State
		path          *tftypes.AttributePath
		expectIDs     []string
		expectLastRaw map[string]any
		expectValue   any
	}{
		"returns input unchanged for non-list values": {
			graphValue:  map[string]any{"id": "a"},
			path:        tftypes.NewAttributePath().WithAttributeName("rules"),
			expectValue: map[string]any{"id": "a"},
		},
		"returns input unchanged for single-element lists": {
			graphValue: []any{
				map[string]any{"id": "a"},
			},
			path: tftypes.NewAttributePath().WithAttributeName("rules"),
			expectValue: []any{
				map[string]any{"id": "a"},
			},
		},
		"reorders using prior state and appends unmatched elements": {
			graphValue: []any{
				map[string]any{"id": "b", "name": "Beta"},
				map[string]any{"id": "c", "name": "Gamma"},
				map[string]any{"id": "a", "name": "Alpha"},
				map[string]any{"name": "no-id"},
			},
			priorState: testPriorStateForOrderedKeys("rules", "id", "a", "b"),
			path:       tftypes.NewAttributePath().WithAttributeName("rules"),
			expectIDs:  []string{"a", "b", "c"},
			expectLastRaw: map[string]any{
				"name": "no-id",
			},
		},
		"falls back to alphabetical ordering without prior state order": {
			graphValue: []any{
				map[string]any{"id": "c"},
				map[string]any{"id": "a"},
				map[string]any{"id": "b"},
			},
			path:      tftypes.NewAttributePath().WithAttributeName("rules"),
			expectIDs: []string{"a", "b", "c"},
		},
		"falls back to alphabetical ordering when path has no attribute name": {
			graphValue: []any{
				map[string]any{"id": "b"},
				map[string]any{"id": "a"},
			},
			priorState: testPriorStateForOrderedKeys("rules", "id", "rules-first"),
			path:       tftypes.NewAttributePath().WithElementKeyInt(0),
			expectIDs:  []string{"a", "b"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := validator.ModifyGraphToTerraform(tc.graphValue, tc.priorState, tc.path)
			if err != nil {
				t.Fatalf("ModifyGraphToTerraform returned error: %v", err)
			}

			if tc.expectValue != nil {
				if !reflect.DeepEqual(got, tc.expectValue) {
					t.Fatalf("unexpected value\nexpected: %#v\nactual:   %#v", tc.expectValue, got)
				}
				return
			}

			gotList, ok := got.([]any)
			if !ok {
				t.Fatalf("expected []any result, got %T", got)
			}

			gotIDs := testExtractIDs(t, gotList, len(tc.expectIDs))
			if !reflect.DeepEqual(gotIDs, tc.expectIDs) {
				t.Fatalf("unexpected ids\nexpected: %#v\nactual:   %#v", tc.expectIDs, gotIDs)
			}

			if tc.expectLastRaw != nil {
				last, ok := gotList[len(gotList)-1].(map[string]any)
				if !ok {
					t.Fatalf("expected final element to be map[string]any, got %T", gotList[len(gotList)-1])
				}

				if !reflect.DeepEqual(last, tc.expectLastRaw) {
					t.Fatalf("unexpected final element\nexpected: %#v\nactual:   %#v", tc.expectLastRaw, last)
				}
			}
		})
	}
}

func testPriorStateForOrderedKeys(attributeName string, keyName string, orderedKeys ...string) *tfsdk.State {
	elementType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{keyName: tftypes.String}}
	elements := make([]tftypes.Value, 0, len(orderedKeys))
	for _, key := range orderedKeys {
		elements = append(elements, tftypes.NewValue(elementType, map[string]tftypes.Value{
			keyName: tftypes.NewValue(tftypes.String, key),
		}))
	}

	return &tfsdk.State{
		Raw: tftypes.NewValue(
			tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				attributeName: tftypes.List{ElementType: elementType},
			}},
			map[string]tftypes.Value{
				attributeName: tftypes.NewValue(tftypes.List{ElementType: elementType}, elements),
			},
		),
	}
}

func testExtractIDs(t *testing.T, list []any, count int) []string {
	t.Helper()

	ids := make([]string, 0, count)
	for i := 0; i < count; i++ {
		elemMap, ok := list[i].(map[string]any)
		if !ok {
			t.Fatalf("expected element %d to be map[string]any, got %T", i, list[i])
		}

		id, ok := elemMap["id"].(string)
		if !ok {
			t.Fatalf("expected element %d id to be string, got %#v", i, elemMap["id"])
		}

		ids = append(ids, id)
	}

	return ids
}

