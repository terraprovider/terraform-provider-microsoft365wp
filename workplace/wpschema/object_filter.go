package wpschema

//
// Heavily inspired by https://github.com/dschnare/jsonfilter/tree/f82bfdf3
//

import (
	"encoding/json"
	"fmt"
)

type filterFunc func(path string, value any) (bool, error)

func Traverse(value any, filter filterFunc) (any, error) {
	return traverseAny(value, "", filter)
}

func TraverseJson(jsonIn []byte, filter filterFunc) (jsonOut []byte, err error) {
	var mapIn map[string]any
	err = json.Unmarshal(jsonIn, &mapIn)
	if err != nil {
		return
	}

	mapOut, err := Traverse(mapIn, filter)
	if err != nil {
		return
	}

	return json.Marshal(mapOut)
}

func traverseAny(value any, path string, filter filterFunc) (any, error) {
	switch switchValue := value.(type) {
	case map[string]any:
		return traverseMap(switchValue, path, filter)
	case []any:
		return traverseSlice(switchValue, path, filter)
	}

	return value, nil
}

func traverseMap(currentMap map[string]any, path string, filter filterFunc) (newMap map[string]any, err error) {
	newMap = make(map[string]any)
	for k, v := range currentMap {
		itemPath := fmt.Sprintf("%s/%s", path, k)
		keep, err := filter(itemPath, v)
		if err != nil {
			break
		}
		if keep {
			if newMap[k], err = traverseAny(v, itemPath, filter); err != nil {
				break
			}
		}
	}
	return
}

func traverseSlice(currentSlice []any, path string, filter filterFunc) (newSlice []any, err error) {
	for k, v := range currentSlice {
		itemPath := fmt.Sprintf("%s[%d]", path, k)
		keep, err := filter(itemPath, v)
		if err != nil {
			break
		}
		if keep {
			if newValue, err := traverseAny(v, itemPath, filter); err == nil {
				newSlice = append(newSlice, newValue)
			} else {
				break
			}
		}
	}
	return
}
