package wpvalidator

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

//
// TranslateValuesBase64
//

var _ AttributeValueTranslator = (*TranslateValueEncodeBase64Impl)(nil)

type TranslateValueEncodeBase64Impl struct {
	translateValuesDefault
}

func TranslateValueEncodeBase64() TranslateValueEncodeBase64Impl {
	return TranslateValueEncodeBase64Impl{}
}

func (t TranslateValueEncodeBase64Impl) TranslateTerraformToGraph(tfValue tftypes.Value) (any, bool, error) {
	if !tfValue.Type().Is(tftypes.String) {
		panic(fmt.Errorf("can only handle strings but this is %s", tfValue.Type().String()))
	}

	var plainStringRaw string
	if err := tfValue.As(&plainStringRaw); err != nil {
		panic(err)
	}
	base64StringRaw := base64.StdEncoding.EncodeToString([]byte(plainStringRaw)) // []byte() will do UTF-8 internally

	return base64StringRaw, true, nil
}

func (t TranslateValueEncodeBase64Impl) TranslateGraphToTerraform(graphValue any) (tftypes.Value, bool, error) {
	base64StringRaw, ok := graphValue.(string)
	if !ok {
		panic(fmt.Errorf("can only handle strings but this is %T", graphValue))
	}

	plainBytes, err := base64.StdEncoding.DecodeString(base64StringRaw)
	if err != nil {
		return tftypes.Value{}, false, err // could also happen due to errors in MS Graph - no panic here
	}
	plainStringRaw := string(plainBytes)

	return tftypes.NewValue(tftypes.String, plainStringRaw), true, nil
}
