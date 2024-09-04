package generic

import (
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// 2023-04 Felix Storm: For now I decided to abuse Attribute.Description to store custom MS Graph attributes (as we only
// use MarkdownDescription which will hide Description anyway).
// While creating a few new composed Attribute types embedding the existing ones plus an additional field for a non-standard
// MS Graph attribute name would work and would probably be a lot more elegant, schema literals would become a lot more
// complicated then also (since the fields of the embedded type can only be filled by explicitely specifying the type
// name of the embedded struct, see e.g. https://stackoverflow.com/questions/44123399/how-do-i-initialize-a-composed-struct-in-go
// or https://github.com/golang/go/issues/9859).
// For OData derived types though, we do create a custom composed Attribute since it only gets used when its specific
// purpose is required.

type DerivedTypeNestedAttribute interface {
	rsschema.NestedAttribute // methods are identical for dsschema.NestedAttribute also
	GetDerivedType() string
}

var (
	_ DerivedTypeNestedAttribute = OdataDerivedTypeNestedAttributeRs{}
	_ DerivedTypeNestedAttribute = OdataDerivedTypeNestedAttributeDs{}
)

type OdataDerivedTypeNestedAttributeRs struct {
	rsschema.SingleNestedAttribute
	DerivedType string
}

func (n OdataDerivedTypeNestedAttributeRs) GetDerivedType() string {
	return n.DerivedType
}

type OdataDerivedTypeNestedAttributeDs struct {
	dsschema.SingleNestedAttribute
	DerivedType string
}

func (n OdataDerivedTypeNestedAttributeDs) GetDerivedType() string {
	return n.DerivedType
}
