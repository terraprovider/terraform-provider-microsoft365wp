package services

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestDeviceConfigurationCustomSkipSecretDecryption(t *testing.T) {
	cases := map[string]struct {
		value    string
		set      bool
		expected bool
	}{
		"unset":        {set: false, expected: false},
		"empty":        {value: "", set: true, expected: false},
		"one":          {value: "1", set: true, expected: true},
		"true":         {value: "true", set: true, expected: true},
		"TRUE":         {value: "TRUE", set: true, expected: true},
		"zero":         {value: "0", set: true, expected: false},
		"false":        {value: "false", set: true, expected: false},
		"other truthy": {value: "yes", set: true, expected: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.set {
				t.Setenv(deviceConfigurationCustomSkipSecretDecryptionEnvVar, tc.value)
			} else {
				t.Setenv(deviceConfigurationCustomSkipSecretDecryptionEnvVar, "")
				// t.Setenv cannot truly unset, but empty is treated as unset by the function.
			}
			if got := deviceConfigurationCustomSkipSecretDecryption(); got != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestDeviceConfigurationCustomOmaValuesFromState(t *testing.T) {
	stringType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"value": tftypes.String}}
	base64Type := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"value_base64": tftypes.String, "file_name": tftypes.String}}
	stringXmlType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"value": tftypes.String, "file_name": tftypes.String}}

	omaSettingType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"oma_uri":    tftypes.String,
		"string":     stringType,
		"base64":     base64Type,
		"string_xml": stringXmlType,
	}}
	omaSettingsType := tftypes.Set{ElementType: omaSettingType}
	windows10Type := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"oma_settings": omaSettingsType}}
	rootType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"windows10": windows10Type}}

	newOmaSetting := func(omaUri string, str, b64, xml *string) tftypes.Value {
		strVal := tftypes.NewValue(stringType, nil)
		if str != nil {
			strVal = tftypes.NewValue(stringType, map[string]tftypes.Value{"value": tftypes.NewValue(tftypes.String, *str)})
		}
		b64Val := tftypes.NewValue(base64Type, nil)
		if b64 != nil {
			b64Val = tftypes.NewValue(base64Type, map[string]tftypes.Value{
				"value_base64": tftypes.NewValue(tftypes.String, *b64),
				"file_name":    tftypes.NewValue(tftypes.String, nil),
			})
		}
		xmlVal := tftypes.NewValue(stringXmlType, nil)
		if xml != nil {
			xmlVal = tftypes.NewValue(stringXmlType, map[string]tftypes.Value{
				"value":     tftypes.NewValue(tftypes.String, *xml),
				"file_name": tftypes.NewValue(tftypes.String, nil),
			})
		}
		return tftypes.NewValue(omaSettingType, map[string]tftypes.Value{
			"oma_uri":    tftypes.NewValue(tftypes.String, omaUri),
			"string":     strVal,
			"base64":     b64Val,
			"string_xml": xmlVal,
		})
	}

	strPtr := func(s string) *string { return &s }

	stateWith := func(windows10 tftypes.Value) *tfsdk.State {
		root := tftypes.NewValue(rootType, map[string]tftypes.Value{"windows10": windows10})
		return &tfsdk.State{Raw: root}
	}

	t.Run("nil state", func(t *testing.T) {
		if got := deviceConfigurationCustomOmaValuesFromState(nil); len(got) != 0 {
			t.Fatalf("expected empty, got %v", got)
		}
	})

	t.Run("windows10 null", func(t *testing.T) {
		got := deviceConfigurationCustomOmaValuesFromState(stateWith(tftypes.NewValue(windows10Type, nil)))
		if len(got) != 0 {
			t.Fatalf("expected empty, got %v", got)
		}
	})

	t.Run("mixed oma settings", func(t *testing.T) {
		omaSettings := tftypes.NewValue(omaSettingsType, []tftypes.Value{
			newOmaSetting("uri/string", strPtr("plain-value"), nil, nil),
			newOmaSetting("uri/base64", nil, strPtr("YmFzZTY0"), nil),
			newOmaSetting("uri/xml", nil, nil, strPtr("<xml/>")),
		})
		windows10 := tftypes.NewValue(windows10Type, map[string]tftypes.Value{"oma_settings": omaSettings})

		got := deviceConfigurationCustomOmaValuesFromState(stateWith(windows10))
		want := map[string]string{
			"uri/string": "plain-value",
			"uri/base64": "YmFzZTY0",
			"uri/xml":    "<xml/>",
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})
}
