package provider

import (
	"context"
	"slices"
	"testing"

	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestProviderSchemas checks that the framework accepts the schemas of all resources and data sources, and that the
// sets of resources get wpplanmodifier.SetAlignToConfigOrder (issue #246).
func TestProviderSchemas(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(New())()
	resp, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range resp.Diagnostics {
		t.Errorf("%s: %s: %s", d.Severity, d.Summary, d.Detail)
	}
	if len(resp.ResourceSchemas) == 0 || len(resp.DataSourceSchemas) == 0 {
		t.Fatalf("no schemas: %d resources, %d data sources", len(resp.ResourceSchemas), len(resp.DataSourceSchemas))
	}

	hasAlign := func(a rsschema.Attribute) bool {
		set, ok := a.(rsschema.SetNestedAttribute)
		if !ok {
			return false
		}
		return slices.ContainsFunc(set.PlanModifiers, wpplanmodifier.IsSetAlignToConfigOrder)
	}
	want := map[string][]string{
		"microsoft365wp_mobile_app":                   {"assignments", "categories"},
		"microsoft365wp_cloud_pc_provisioning_policy": {"domain_join_configurations"},
		"microsoft365wp_device_configuration":         {"assignments"},
		"microsoft365wp_browser_site_list":            {"sites", "shared_cookies"},
	}
	found := 0
	for _, f := range New().Resources(ctx) {
		r := f()
		var md resource.MetadataResponse
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "microsoft365wp"}, &md)
		names, ok := want[md.TypeName]
		if !ok {
			continue
		}
		found++
		var sr resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &sr)
		for _, name := range names {
			if !hasAlign(sr.Schema.Attributes[name]) {
				t.Errorf("%s.%s: SetAlignToConfigOrder missing", md.TypeName, name)
			}
		}
	}
	if found != len(want) {
		t.Errorf("found %d of %d resources to check", found, len(want))
	}
}
