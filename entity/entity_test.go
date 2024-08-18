package entity

import (
	"testing"

	"github.com/ukubenet/metadata-repository/metadata"
)

func TestGetStructedAttributeValues(t *testing.T) {
	t.Run("ReferenceType", func(t *testing.T) {
		values := AttributeValues{
			"refAttr": map[string]interface{}{
				"reference": "ref1",
				"view":      []interface{}{"view1", "view2"},
			},
		}
		attributesMeta := metadata.StructedAttributes{
			"refAttr": {Type: metadata.ReferenceType},
		}

		result := GetStructedAttributeValues(values, attributesMeta)

		if len(result) != 1 {
			t.Fatalf("Expected 1 attribute, got %d", len(result))
		}

		refValue, ok := result["refAttr"].Value.(ReferenceValue)
		if !ok {
			t.Fatalf("Expected ReferenceValue, got %T", result["refAttr"].Value)
		}

		expectedRef := ReferenceValue{"ref1": []string{"view1", "view2"}}
		if len(refValue) != len(expectedRef) {
			t.Errorf("Expected reference with %d keys, got %d", len(expectedRef), len(refValue))
		}
		if ref, ok := refValue["ref1"]; !ok {
			t.Errorf("Expected reference 'ref1', not found in %v", refValue)
		} else if views, ok := ref.([]any); !ok {
			t.Errorf("Expected []string for views, got %T", ref)
		} else if len(views) != len(expectedRef["ref1"].([]string)) {
			t.Errorf("Expected views %v for 'ref1', got %v", expectedRef["ref1"], views)
		}
	})

	t.Run("TableType", func(t *testing.T) {
		values := AttributeValues{
			"tableAttr": map[string]interface{}{
				"rows": []interface{}{
					map[string]interface{}{"col1": "value1", "col2": 42},
					map[string]interface{}{"col1": "value2", "col2": 84},
				},
			},
		}
		attributesMeta := metadata.StructedAttributes{
			"tableAttr": {
				Type: metadata.TableType,
				Specs: metadata.TableSpecs{
					Columns: metadata.StructedAttributes{
						"col1": {Type: metadata.StringType},
						"col2": {Type: metadata.IntegerType},
					},
				},
			},
		}

		result := GetStructedAttributeValues(values, attributesMeta)

		if len(result) != 1 {
			t.Fatalf("Expected 1 attribute, got %d", len(result))
		}

		tableValue, ok := result["tableAttr"].Value.(TableSpecs)
		if !ok {
			t.Fatalf("Expected TableSpecs, got %T", result["tableAttr"].Value)
		}

		if len(tableValue.Rows) != 2 {
			t.Errorf("Expected 2 rows, got %d", len(tableValue.Rows))
		}

		if tableValue.Rows[0]["col1"].Value != "value1" {
			t.Errorf("Expected 'value1' for first row col1, got '%v'", tableValue.Rows[0]["col1"].Value)
		}

		if tableValue.Rows[1]["col2"].Value != 84 {
			t.Errorf("Expected 84 for second row col2, got %v", tableValue.Rows[1]["col2"].Value)
		}
	})

	t.Run("DefaultType", func(t *testing.T) {
		values := AttributeValues{
			"defaultAttr": map[string]interface{}{
				"key": "value",
			},
		}
		attributesMeta := metadata.StructedAttributes{
			"defaultAttr": {Type: metadata.StringType},
		}

		result := GetStructedAttributeValues(values, attributesMeta)

		if len(result) != 1 {
			t.Fatalf("Expected 1 attribute, got %d", len(result))
		}

		defaultValue, ok := result["defaultAttr"].Value.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", result["defaultAttr"].Value)
		}

		if defaultValue["key"] != "value" {
			t.Errorf("Expected 'value' for key, got '%v'", defaultValue["key"])
		}
	})
}
