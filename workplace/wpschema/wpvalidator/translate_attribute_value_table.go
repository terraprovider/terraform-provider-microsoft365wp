package wpvalidator

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

//
// TranslateValueTable
//

var _ AttributeValueTranslator = (*TranslateValueTable)(nil)

type TranslateValueTable struct {
	translateValuesDefault
	translationTable []TranslationTuple
	toGraph          bool
	toTerraform      bool
}

type TranslationTuple struct {
	TerraformValue tftypes.Value
	GraphValue     any
}

func TranslateValue(translationTable []TranslationTuple) TranslateValueTable {
	return TranslateValueTable{
		translationTable: translationTable,
		toGraph:          true,
		toTerraform:      true,
	}
}

func TranslateValueToGraphOnly(translationTable []TranslationTuple) TranslateValueTable {
	return TranslateValueTable{
		translationTable: translationTable,
		toGraph:          true,
	}
}

func TranslateValueToTerraformOnly(translationTable []TranslationTuple) TranslateValueTable {
	return TranslateValueTable{
		translationTable: translationTable,
		toTerraform:      true,
	}
}

func TranslateGraphNullWithTfDefault(tftype tftypes.Type) TranslateValueTable {
	return TranslateValueTable{
		translationTable: []TranslationTuple{
			{
				GraphValue:     nil,
				TerraformValue: translateAttrValueGetTfDefault(tftype),
			},
		},
		toGraph:     true,
		toTerraform: true,
	}
}

func TranslateGraphMissingWithTfDefault(tftype tftypes.Type) TranslateValueTable {
	return TranslateValueTable{
		translationTable: []TranslationTuple{
			{
				GraphValue:     AttributeMising,
				TerraformValue: translateAttrValueGetTfDefault(tftype),
			},
		},
		toTerraform: true,
	}
}

func translateAttrValueGetTfDefault(tftype tftypes.Type) tftypes.Value {
	var defaultForNewValue any
	switch {
	case tftype.Is(tftypes.Bool):
		defaultForNewValue = false
	case tftype.Is(tftypes.Number):
		defaultForNewValue = 0
	case tftype.Is(tftypes.String):
		defaultForNewValue = ""
	case tftype.Is(tftypes.List{}) || tftype.Is(tftypes.Set{}):
		defaultForNewValue = []tftypes.Value{}
	case tftype.Is(tftypes.Map{}) || tftype.Is(tftypes.Object{}):
		defaultForNewValue = make(map[string]tftypes.Value)
	default:
		panic(fmt.Errorf("dont' know how to handle %s", tftype.String()))
	}
	return tftypes.NewValue(tftype, defaultForNewValue)
}

func (t TranslateValueTable) TranslateTerraformToGraph(tfValue tftypes.Value) (any, bool, error) {
	if t.toGraph {
		for _, t := range t.translationTable {
			if t.TerraformValue.Equal(tfValue) {
				return t.GraphValue, true, nil
			}
		}
	}
	return nil, false, nil
}

func (t TranslateValueTable) TranslateGraphToTerraform(graphValue any) (tftypes.Value, bool, error) {
	if t.toTerraform {
		for _, t := range t.translationTable {
			if t.GraphValue == graphValue {
				return t.TerraformValue, true, nil
			}
		}
	}
	return tftypes.Value{}, false, nil
}
