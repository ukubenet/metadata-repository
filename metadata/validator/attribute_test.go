package metavalidator

import (
	"strings"
	"testing"

	"github.com/ukubenet/metadata-repository/metadata"
	metastorage "github.com/ukubenet/metadata-repository/metadata/storage"
)

func TestMain(m *testing.M) {
	metastorage.SetEnv("test")
	m.Run()

}

func TestValidateAttributeNoType(t *testing.T) {
	attributes := metadata.Attributes{
		"name": {"name": "name"},
	}

	err := ValidateAttributes(attributes)
	if !strings.Contains(err.Error(), "attribute \"name\" doesn't have type") {
		t.Fatal(err)
	}
}

func TestValidateAttributeTypeIsNotString(t *testing.T) {
	attribute := metadata.Attribute{
		"name": "name",
		"type": 1,
	}

	err := ValidateAttribute("name", attribute)
	if !strings.Contains(err.Error(), "type of attribute \"name\" is not a string") {
		t.Fatal(err)
	}
}

func TestValidateIncorrectAttributeType(t *testing.T) {
	attribute := metadata.Attribute{
		"name": "name",
		"type": "unsupported_type",
	}

	err := ValidateAttribute("name", attribute)
	if !strings.Contains(
		err.Error(),
		"type \"unsupported_type\" of attribute \"name\" is not defined",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceAttributeNoReference(t *testing.T) {
	attribute := metadata.Attribute{
		"type": "reference",
	}

	err := ValidateAttribute("name", attribute)
	if !strings.Contains(
		err.Error(),
		"reference attribute \"name\" has missed reference property",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceWrongType(t *testing.T) {
	attribute := metadata.Attribute{
		"type":      "reference",
		"reference": 1,
	}

	err := ValidateAttribute("name", attribute)
	if !strings.Contains(
		err.Error(),
		"reference of attribute \"name\" is not a string",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceWrongReference(t *testing.T) {
	metastorage.SetEnv("test")
	attribute := metadata.Attribute{
		"type":      "reference",
		"reference": "test/Reference2",
	}

	err := validateReference("name", attribute)
	if !strings.Contains(
		err.Error(),
		"error to read reference in attribute \"name\"",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceViewIsNotSlice(t *testing.T) {
	metastorage.SetEnv("test")
	attribute := metadata.Attribute{
		"type":      "reference",
		"reference": "test/Reference",
		"view":      "view",
	}

	err := validateReference("name", attribute)
	if !strings.Contains(
		err.Error(),
		"view of reference attribute \"name\" is not a slice",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceViewUnmatchedAttributes(t *testing.T) {
	metastorage.SetEnv("test")
	attribute := metadata.Attribute{
		"type":      "reference",
		"reference": "test/Reference",
		"view":      []interface{}{"view"},
	}

	err := validateReference("name", attribute)
	if !strings.Contains(
		err.Error(),
		"\"view\" from view of reference attribute \"name\" don't belong to reference entity",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceViewSuccess(t *testing.T) {
	metastorage.SetEnv("test")
	attribute := metadata.Attribute{
		"type":      "reference",
		"reference": "test/Reference",
		"view":      []interface{}{"name"},
	}

	err := validateReference("name", attribute)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateTableNoColumns(t *testing.T) {
	attribute := metadata.Attribute{
		"type": "table",
	}

	err := validateTable("name", attribute)
	if !strings.Contains(
		err.Error(),
		"table attribute \"name\" has missed columns property",
	) {
		t.Fatal(err)
	}
}

func TestValidateTableInvalidColumnsType(t *testing.T) {
	attribute := metadata.Attribute{
		"type":    "table",
		"columns": []interface{}{"name"},
	}

	err := ValidateAttribute("name", attribute)
	if !strings.Contains(
		err.Error(),
		"columns property of attribute \"name\" is malformed",
	) {
		t.Fatal(err)
	}
}

func TestValidateTableSuccess(t *testing.T) {
	attribute := metadata.Attribute{
		"type":    "table",
		"columns": map[string]any{"column": map[string]any{"type": "string"}},
	}

	err := ValidateAttribute("name", attribute)
	if err != nil {
		t.Fatal(err)
	}
}
