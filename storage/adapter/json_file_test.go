package adapter

import (
	"testing"

	"github.com/ukubenet/metadata-repository/models"
)

func TestJsonReader(t *testing.T) {
	reader := JSON()

	candidate := new(models.Entity)

	err := reader.Read("test/Test", candidate)
	if err != nil {
		t.Fatal(err)
	}

	if candidate.EntityName != "test/Test" {
		t.Fail()
	}

	if candidate.Attributes == nil {
		t.Fail()
	}
}

func TestJsonReplacer(t *testing.T) {
	inserter := JSON()

	var candidate *models.Entity = &models.Entity{
		EntityName: "test/Test",
		Attributes: []models.Attribute{
			{Name: "string_attribute", Type: "string"},
			{Name: "number_attribute", Type: "number"},
		},
	}

	inserter.Put(candidate)
}
