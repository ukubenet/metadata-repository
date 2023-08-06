package adapter

import (
	"os"
	"testing"

	"github.com/ukubenet/metadata-repository/models"
)

var path string

func TestMain(m *testing.M) {
	path, _ = os.Getwd()
	path += "/"
	m.Run()
}

func TestJsonReader(t *testing.T) {
	reader := JSON(path)

	entity := new(models.Entity)

	err := reader.Read("test", "Test", entity)
	if err != nil {
		t.Fatal(err)
	}

	if entity.EntityName != "test" {
		t.Fail()
	}

	if entity.Identifier != "Test" {
		t.Fail()
	}

	if entity.Attributes == nil {
		t.Fail()
	}
}

func TestJsonReplacer(t *testing.T) {
	inserter := JSON(path)
	var attributes []models.AttributeValue = []models.AttributeValue{
		{"string_attribute": "string"},
		{"number_attribute": 100},
	}

	var candidate *models.Entity = &models.Entity{
		EntityName: "test",
		Attributes: attributes,
		Identifier: "Test",
	}

	inserter.Put(candidate)
}

func TestJsonLister(t *testing.T) {
	lister := JSON(path)
	list := []models.Entity{}

	lister.List("test", &list)

	if len(list) != 1 {
		t.Fail()
	}

	if list[0].Identifier != "Test" {
		t.Fail()
	}
}

func TestJsonTypeLister(t *testing.T) {
	lister := JSON(path)
	list := []string{}

	lister.TypeList(&list)

	if len(list) != 1 {
		t.Fail()
	}

	if list[0] != "test" {
		t.Fail()
	}
}
