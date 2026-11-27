package generic

import (
	"testing"

	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// TestOdataDerivedTypeNestedAttributeGetTypeCache checks that the cached type equals the uncached one and that
// different nested attribute maps get their own types.
func TestOdataDerivedTypeNestedAttributeGetTypeCache(t *testing.T) {
	inner := map[string]rsschema.Attribute{"x": rsschema.StringAttribute{Optional: true}}
	a := OdataDerivedTypeNestedAttributeRs{SingleNestedAttribute: rsschema.SingleNestedAttribute{Optional: true, Attributes: map[string]rsschema.Attribute{
		"s": rsschema.StringAttribute{Optional: true},
		"n": OdataDerivedTypeNestedAttributeRs{SingleNestedAttribute: rsschema.SingleNestedAttribute{Optional: true, Attributes: inner}},
	}}}
	b := OdataDerivedTypeNestedAttributeRs{SingleNestedAttribute: rsschema.SingleNestedAttribute{Optional: true, Attributes: map[string]rsschema.Attribute{
		"b": rsschema.BoolAttribute{Optional: true},
	}}}

	want := a.SingleNestedAttribute.GetType()
	for i := 0; i < 2; i++ { // second call from the cache
		if got := a.GetType(); !got.Equal(want) {
			t.Fatalf("call %d: got %s, want %s", i, got, want)
		}
	}
	if b.GetType().Equal(a.GetType()) {
		t.Error("different attribute maps must not share a cached type")
	}
	// a nil map is not cached
	empty := OdataDerivedTypeNestedAttributeRs{SingleNestedAttribute: rsschema.SingleNestedAttribute{Optional: true}}
	if !empty.GetType().Equal(empty.SingleNestedAttribute.GetType()) {
		t.Error("nil map: types differ")
	}
}
