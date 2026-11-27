package generic

import (
	"reflect"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// The framework derives the type of a nested attribute from its nested attributes on every GetType call, i.e.
// recursively for the whole subtree. It does so for every value node while planning (e.g. in TransformDefaults and
// ReifyNullCollectionBlocks via fromtftypes.AttributePath -> Schema.TypeAtTerraformPath). For the deep schemas with
// many OData derived types (e.g. the settings catalog of device_management_configuration_policy) this made a plan of a
// policy with 90 settings take about 27 s of CPU (38 GB allocated). The types never change once the schema has been
// built, so they get cached per nested attribute map (#248).
var nestedAttributeTypeCache sync.Map // map pointer (uintptr) -> attr.Type

// cachedNestedAttributeType returns the cached type for the nested attributes map attributes, computing it with
// compute on first use. Without a map (nil) the type is not cached.
func cachedNestedAttributeType(attributes any, compute func() attr.Type) attr.Type {
	v := reflect.ValueOf(attributes)
	if v.Kind() != reflect.Map || v.IsNil() {
		return compute()
	}
	key := v.Pointer()
	if t, ok := nestedAttributeTypeCache.Load(key); ok {
		return t.(attr.Type)
	}
	t := compute()
	nestedAttributeTypeCache.Store(key, t)
	return t
}
