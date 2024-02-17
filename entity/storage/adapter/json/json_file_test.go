package adapter

import (
	"os"
	"testing"

	"github.com/ukubenet/metadata-repository/entity"
)

var path string

func TestMain(m *testing.M) {
	path, _ = os.Getwd()
	path += "/"
	m.Run()
}

func TestJsonReader(t *testing.T) {
	reader := JSON(path)

	entity := new(entity.CatalogEntity)

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
	var attributes entity.AttributeValues = entity.AttributeValues{
		"string_attribute": "string",
		"number_attribute": 100,
	}

	var candidate *entity.CatalogEntity = &entity.CatalogEntity{
		EntityName: "test",
		Attributes: attributes,
		Identifier: "Test",
	}

	inserter.Put(candidate)
}

func TestJsonLister(t *testing.T) {
	lister := JSON(path)
	list := []entity.CatalogEntity{}

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
