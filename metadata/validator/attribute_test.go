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

func TestValidateAttributesNoName(t *testing.T) {
	attributes := []metadata.Attribute{}
	attribute := metadata.Attribute{
		"type": "string",
	}
	attributes = append(attributes, attribute)

	err := ValidateAttributes(attributes)
	if !strings.Contains(err.Error(), "attribute 0 doesn't have name") {
		t.Fatal(err)
	}

}

func TestValidateAttributesNameNotString(t *testing.T) {
	attributes := []metadata.Attribute{}
	attribute := metadata.Attribute{
		"name": 0,
		"type": "string",
	}
	attributes = append(attributes, attribute)

	err := ValidateAttributes(attributes)
	if !strings.Contains(err.Error(), "name of attribute 0 is not a string") {
		t.Fatal(err)
	}
}

func TestValidateAttributeNoType(t *testing.T) {
	attribute := metadata.Attribute{
		"name": "name",
	}

	err := validateAttribute(attribute)
	if !strings.Contains(err.Error(), "attribute \"name\" doesn't have type") {
		t.Fatal(err)
	}
}

func TestValidateAttributeTypeIsNotString(t *testing.T) {
	attribute := metadata.Attribute{
		"name": "name",
		"type": 1,
	}

	err := validateAttribute(attribute)
	if !strings.Contains(err.Error(), "type of attribute \"name\" is not a string") {
		t.Fatal(err)
	}
}

func TestValidateIncorrectAttributeType(t *testing.T) {
	attribute := metadata.Attribute{
		"name": "name",
		"type": "unsupported_type",
	}

	err := validateAttribute(attribute)
	if !strings.Contains(
		err.Error(),
		"type \"unsupported_type\" of attribute \"name\" is not defined",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceAttributeNoReference(t *testing.T) {
	attribute := metadata.Attribute{
		"name": "name",
		"type": "reference",
	}

	err := validateAttribute(attribute)
	if !strings.Contains(
		err.Error(),
		"reference attribute \"name\" has missed reference property",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceWrongType(t *testing.T) {
	attribute := metadata.Attribute{
		"name":      "name",
		"type":      "reference",
		"reference": 1,
	}

	err := validateAttribute(attribute)
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
		"name":      "name",
		"type":      "reference",
		"reference": "test/Reference2",
	}

	err := validateReference(attribute)
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
		"name":      "name",
		"type":      "reference",
		"reference": "test/Reference",
		"view":      "view",
	}

	err := validateReference(attribute)
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
		"name":      "name",
		"type":      "reference",
		"reference": "test/Reference",
		"view":      []string{"view"},
	}

	err := validateReference(attribute)
	if !strings.Contains(
		err.Error(),
		"some attributes from view of reference attribute \"name\" don't belong to reference entity. Possible attributes: [\"name\"]",
	) {
		t.Fatal(err)
	}
}

func TestValidateReferenceViewSuccess(t *testing.T) {
	metastorage.SetEnv("test")
	attribute := metadata.Attribute{
		"name":      "name",
		"type":      "reference",
		"reference": "test/Reference",
		"view":      []string{"name"},
	}

	err := validateReference(attribute)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateTableNoColumns(t *testing.T) {
	attribute := metadata.Attribute{
		"name": "name",
		"type": "table",
	}

	err := validateTable(attribute)
	if !strings.Contains(
		err.Error(),
		"table attribute \"name\" has missed columns property",
	) {
		t.Fatal(err)
	}
}

func TestValidateTableInvalidColumnsType(t *testing.T) {
	attribute := metadata.Attribute{
		"name":    "name",
		"type":    "table",
		"columns": []string{"name"},
	}

	err := validateAttribute(attribute)
	if !strings.Contains(
		err.Error(),
		"columns property of attribute \"name\" is not a list of attributes",
	) {
		t.Fatal(err)
	}
}

func TestValidateTableSuccess(t *testing.T) {
	attribute := metadata.Attribute{
		"name": "name",
		"type": "table",
		"columns": []metadata.Attribute{
			{"name": "column", "type": "string"},
		},
	}

	err := validateAttribute(attribute)
	if err != nil {
		t.Fatal(err)
	}
}
