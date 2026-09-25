package wpjsontypes

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestNormalizedStringSemanticEquals(t *testing.T) {
	ignoreVersion := func(path string, value any) (bool, error) {
		return !strings.HasSuffix(path, "/version"), nil
	}

	cases := map[string]struct {
		a, b             string
		filter           func(path string, value any) (bool, error)
		ignoreArrayOrder bool
		expected         bool
	}{
		"property order": {
			a:        `{"a":1,"b":2}`,
			b:        `{"b":2,"a":1}`,
			expected: true,
		},
		"array order significant by default": {
			a:        `[1,2,3]`,
			b:        `[3,2,1]`,
			expected: false,
		},
		"array order ignored": {
			a:                `[1,2,3]`,
			b:                `[3,2,1]`,
			ignoreArrayOrder: true,
			expected:         true,
		},
		"nested reordered arrays with differing property order": {
			a:                `[{"name":"u","attributeMappings":[{"target":"mail","params":[{"key":"a"},{"key":"b"}]},{"target":"otherMails","params":[]}]},{"name":"g","attributeMappings":[]}]`,
			b:                `[{"attributeMappings":[],"name":"g"},{"name":"u","attributeMappings":[{"params":[],"target":"otherMails"},{"target":"mail","params":[{"key":"b"},{"key":"a"}]}]}]`,
			ignoreArrayOrder: true,
			expected:         true,
		},
		"changed value still differs": {
			a:                `[{"target":"mail","flowType":"Always"},{"target":"otherMails"}]`,
			b:                `[{"target":"otherMails"},{"target":"mail","flowType":"ObjectAddOnly"}]`,
			ignoreArrayOrder: true,
			expected:         false,
		},
		"duplicate elements are counted": {
			a:                `[1,1,2]`,
			b:                `[1,2,2]`,
			ignoreArrayOrder: true,
			expected:         false,
		},
		"filter and ignored array order combined": {
			a:                `[{"name":"x","version":"1","objects":[{"n":"a"},{"n":"b"}]}]`,
			b:                `[{"name":"x","version":"2","objects":[{"n":"b"},{"n":"a"}]}]`,
			filter:           ignoreVersion,
			ignoreArrayOrder: true,
			expected:         true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			typ := NormalizedType{WpObjectFilterFunc: tc.filter, WpIgnoreArrayOrder: tc.ignoreArrayOrder}
			a, _ := typ.ValueFromString(context.Background(), basetypes.NewStringValue(tc.a))
			b, _ := typ.ValueFromString(context.Background(), basetypes.NewStringValue(tc.b))

			got, diags := a.(Normalized).StringSemanticEquals(context.Background(), b)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if got != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
