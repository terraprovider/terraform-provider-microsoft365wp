package services

import (
	"reflect"
	"testing"
)

func TestGroupAssignedLicensesDiff(t *testing.T) {
	const (
		skuA = "aaaaaaaa-0000-0000-0000-000000000001"
		skuB = "bbbbbbbb-0000-0000-0000-000000000002"
		skuC = "cccccccc-0000-0000-0000-000000000003"
		p1   = "11111111-0000-0000-0000-000000000001"
		p2   = "22222222-0000-0000-0000-000000000002"
	)
	lic := func(sku string, plans ...string) map[string]any {
		planList := []any{}
		for _, p := range plans {
			planList = append(planList, p)
		}
		return map[string]any{"skuId": sku, "disabledPlans": planList}
	}
	add := func(sku string, plans ...string) map[string]any {
		if plans == nil {
			plans = []string{}
		}
		return map[string]any{"skuId": sku, "disabledPlans": plans}
	}

	cases := map[string]struct {
		current       []any
		desired       []any
		authoritative bool
		prevManaged   map[string]bool
		wantAdd       []any
		wantRemove    []string
	}{
		"create on empty group adds all": {
			desired:    []any{lic(skuA, p1), lic(skuB)},
			wantAdd:    []any{add(skuA, p1), add(skuB)},
			wantRemove: []string{},
		},
		"create converges: already present is no-op": {
			current:    []any{lic(skuA, p1)},
			desired:    []any{lic(skuA, p1)},
			wantAdd:    []any{},
			wantRemove: []string{},
		},
		"case and order of GUIDs are ignored": {
			current:    []any{lic(skuA, p2, p1)},
			desired:    []any{lic("AAAAAAAA-0000-0000-0000-000000000001", "11111111-0000-0000-0000-000000000001", p2)},
			wantAdd:    []any{},
			wantRemove: []string{},
		},
		"changed disabled plans re-add the SKU": {
			current:    []any{lic(skuA, p1)},
			desired:    []any{lic(skuA, p2)},
			wantAdd:    []any{add(skuA, p2)},
			wantRemove: []string{},
		},
		"non-authoritative create keeps unmanaged SKUs": {
			current:    []any{lic(skuC)},
			desired:    []any{lic(skuA)},
			wantAdd:    []any{add(skuA)},
			wantRemove: []string{},
		},
		"authoritative create removes unmanaged SKUs": {
			current:       []any{lic(skuC), lic(skuA)},
			desired:       []any{lic(skuA)},
			authoritative: true,
			wantAdd:       []any{},
			wantRemove:    []string{skuC},
		},
		"non-authoritative update removes only previously managed SKUs": {
			current:     []any{lic(skuA), lic(skuB), lic(skuC)},
			desired:     []any{lic(skuA)},
			prevManaged: map[string]bool{skuA: true, skuB: true},
			wantAdd:     []any{},
			wantRemove:  []string{skuB},
		},
		"non-authoritative delete removes previously managed SKUs only": {
			current:     []any{lic(skuA), lic(skuC)},
			prevManaged: map[string]bool{skuA: true, skuB: true},
			wantAdd:     []any{},
			wantRemove:  []string{skuA},
		},
		"authoritative delete removes everything": {
			current:       []any{lic(skuB), lic(skuA)},
			authoritative: true,
			wantAdd:       []any{},
			wantRemove:    []string{skuA, skuB},
		},
		"mixed diff in one call": {
			current:     []any{lic(skuA, p1), lic(skuB)},
			desired:     []any{lic(skuA, p2), lic(skuC)},
			prevManaged: map[string]bool{skuA: true, skuB: true},
			wantAdd:     []any{add(skuA, p2), add(skuC)},
			wantRemove:  []string{skuB},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			gotAdd, gotRemove, err := groupAssignedLicensesDiff(c.current, c.desired, c.authoritative, c.prevManaged)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(gotAdd, c.wantAdd) {
				t.Errorf("add: got %#v, want %#v", gotAdd, c.wantAdd)
			}
			if !reflect.DeepEqual(gotRemove, c.wantRemove) {
				t.Errorf("remove: got %#v, want %#v", gotRemove, c.wantRemove)
			}
		})
	}
}

func TestGroupAssignedLicensesDiffInvalidElement(t *testing.T) {
	if _, _, err := groupAssignedLicensesDiff(nil, []any{"not a map"}, false, nil); err == nil {
		t.Fatal("expected error for invalid element")
	}
}
