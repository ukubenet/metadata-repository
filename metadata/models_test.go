package metadata

import (
	"reflect"
	"testing"
)

func TestGetStructedAttributes(t *testing.T) {
	tests := []struct {
		name       string
		attributes Attributes
		want       StructedAttributes
	}{
		{
			name: "Reference type attribute",
			attributes: Attributes{
				"refAttr": {
					"type":       ReferenceType,
					"reference":  "someRef",
					"entityType": Catalog,
					"view":       []string{"view1", "view2"},
				},
			},
			want: StructedAttributes{
				"refAttr": {
					Type: ReferenceType,
					Specs: ReferenceSpecs{
						Reference:  "someRef",
						EntityType: Catalog,
						View:       []string{"view1", "view2"},
					},
				},
			},
		},
		{
			name: "Table type attribute",
			attributes: Attributes{
				"tableAttr": {
					"type": TableType,
					"columns": Attributes{
						"col1": {"type": "string"},
						"col2": {"type": "int"},
					},
				},
			},
			want: StructedAttributes{
				"tableAttr": {
					Type: TableType,
					Specs: TableSpecs{
						Columns: StructedAttributes{
							"col1": {Type: "string", Specs: PrimitiveType{}},
							"col2": {Type: "int", Specs: PrimitiveType{}},
						},
					},
				},
			},
		},
		{
			name: "Primitive type attribute",
			attributes: Attributes{
				"primAttr": {"type": "string"},
			},
			want: StructedAttributes{
				"primAttr": {Type: "string", Specs: PrimitiveType{}},
			},
		},
		{
			name: "Mixed attribute types",
			attributes: Attributes{
				"refAttr": {
					"type":       ReferenceType,
					"reference":  "someRef",
					"entityType": Catalog,
					"view":       []string{"view1"},
				},
				"tableAttr": {
					"type": TableType,
					"columns": Attributes{
						"col1": {"type": "int"},
					},
				},
				"primAttr": {"type": "bool"},
			},
			want: StructedAttributes{
				"refAttr": {
					Type: ReferenceType,
					Specs: ReferenceSpecs{
						Reference:  "someRef",
						EntityType: Catalog,
						View:       []string{"view1"},
					},
				},
				"tableAttr": {
					Type: TableType,
					Specs: TableSpecs{
						Columns: StructedAttributes{
							"col1": {Type: "int", Specs: PrimitiveType{}},
						},
					},
				},
				"primAttr": {Type: "bool", Specs: PrimitiveType{}},
			},
		},
		{
			name:       "Empty attributes",
			attributes: Attributes{},
			want:       StructedAttributes{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStructedAttributes(tt.attributes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStructedAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}
