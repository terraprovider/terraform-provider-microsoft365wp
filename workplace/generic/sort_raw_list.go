package generic

import "sort"

// SortRawListByStringKey sorts a list attribute within a raw map[string]any by the given string key.
// The attribute identified by listAttr is expected to be a []any of map[string]any items.
// Items are sorted in ascending order by the string value of the given sortKey field.
func SortRawListByStringKey(rawVal map[string]any, listAttr string, sortKey string) {
	list, ok := rawVal[listAttr].([]any)
	if !ok {
		return
	}
	sort.SliceStable(list, func(i, j int) bool {
		iMap, iOk := list[i].(map[string]any)
		jMap, jOk := list[j].(map[string]any)
		if !iOk || !jOk {
			return false
		}
		iVal, _ := iMap[sortKey].(string)
		jVal, _ := jMap[sortKey].(string)
		return iVal < jVal
	})
}

