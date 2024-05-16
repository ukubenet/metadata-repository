package entityvalidator

import (
	"strings"
	"testing"

	config "github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/entity"
	"github.com/ukubenet/metadata-repository/metadata"
)

func TestMain(m *testing.M) {
	config.LoadConfig("../../config", "test")
	m.Run()

}

func TestValidateAttributesBrokenMeta(t *testing.T) {
	values := entity.AttributeValues{
		"name": "string",
	}

	err := ValidateAttributeValues(metadata.Catalog, "missed entity meta", values)
	if !strings.Contains(err.Error(), "error reading meta of entity \"missed entity meta\"") {
		t.Fatal(err)
	}
}

func TestValidateAttributeValueNoMeta(t *testing.T) {
	values := entity.AttributeValues{
		"name": "sample",
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if !strings.Contains(err.Error(), "no such attribute \"name\"") {
		t.Fatal(err)
	}
}

func TestValidateAttributeReferenceMalformed(t *testing.T) {
	values := entity.AttributeValues{
		"sample": "string",
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if !strings.Contains(err.Error(), "reference attribute \"sample\" is malformed") {
		t.Fatal(err)
	}
}

func TestValidateAttributeReferenceViewAttributeMismatch(t *testing.T) {
	values := map[string]any{
		"sample": map[string]any{
			"reference": "Reference",
			"view": map[string]any{
				"missed": "value",
			},
		},
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if !strings.Contains(
		err.Error(),
		"attribute \"missed\" is not present in reference entity \"Reference\". entity attribute: \"sample\"",
	) {
		t.Fatal(err)
	}
}

func TestValidateAttributeReferenceViewValueMismatch(t *testing.T) {
	values := map[string]any{
		"sample": map[string]any{
			"reference": "Reference",
			"view": map[string]any{
				"name": "value",
			},
		},
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if !strings.Contains(
		err.Error(),
		"value \"value\" from view attribute \"name\" of reference attribute \"sample\" don't belong to reference entity. Value in ref entity: \"string\"",
	) {
		t.Fatal(err)
	}
}

func TestValidateAttributeReferenceSuccess(t *testing.T) {
	values := map[string]any{
		"sample": map[string]any{
			"reference": "Reference",
			"view": map[string]any{
				"name": "string",
			},
		},
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateAttributeReferenceReadError(t *testing.T) {
	values := map[string]any{
		"sample": map[string]any{
			"reference": "No Reference",
			"view": map[string]any{
				"missed": "value",
			},
		},
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if !strings.Contains(
		err.Error(),
		"error to read reference \"No Reference\" in attribute \"sample\"",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceAttributeViewNotMap(t *testing.T) {
	values := map[string]any{
		"sample": map[string]any{
			"reference": "No Reference",
			"view":      "not a map",
		},
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if !strings.Contains(
		err.Error(),
		"view of reference attribute \"sample\" is not a map",
	) {
		t.Fatal(err)
	}
}

func TestValidateRefereceAttributeViewNotString(t *testing.T) {
	values := map[string]any{
		"sample": map[string]any{
			"reference": 1,
			"view":      "not a map",
		},
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if !strings.Contains(
		err.Error(),
		"reference of attribute \"sample\" does not exist or not a string",
	) {
		t.Fatal(err)
	}
}

func TestValidateTableAttributeColumnsMismatch(t *testing.T) {
	values := map[string]any{
		"table": map[string]any{
			"columns": map[string]any{
				"name": "sample",
			},
		},
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if !strings.Contains(
		err.Error(),
		"meta property for column \"name\" of table attribute \"table\" does not exist",
	) {
		t.Fatal(err)
	}
}

func TestValidateTableAttributeSuccess(t *testing.T) {
	values := map[string]any{
		"table": map[string]any{
			"columns": map[string]any{
				"title": "sample",
			},
		},
	}

	err := ValidateAttributeValues(metadata.Catalog, "test/TestMeta", values)
	if err != nil {
		t.Fatal(err)
	}
}
