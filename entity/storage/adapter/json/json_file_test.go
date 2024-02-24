package adapter

import (
	"os"
	"testing"

	"github.com/ukubenet/metadata-repository/entity"
	"github.com/ukubenet/metadata-repository/metadata"
)

var path string

func TestMain(m *testing.M) {
	path, _ = os.Getwd()
	path += "/"
	m.Run()
}

func TestJsonReader(t *testing.T) {
	reader := JSON(path)

	entity, err := reader.Read(metadata.Catalog, "test", "Test")
	if err != nil {
		t.Fatal(err)
	}

	if entity.GetName() != "test" {
		t.Fail()
	}

	if entity.GetID() != "Test" {
		t.Fail()
	}

	if entity.GetAttributes() == nil {
		t.Fail()
	}
}

func TestJsonReplacer(t *testing.T) {
	inserter := JSON(path)
	var attributes entity.AttributeValues = entity.AttributeValues{
		"string_attribute": "string",
		"number_attribute": 100,
	}

	var candidate = entity.CatalogEntity{
		Metadata: entity.Metadata{
			EntityName: "test",
			Attributes: attributes,
			Identifier: "Test",
		},
	}

	inserter.Put(candidate)
}

func TestJsonLister(t *testing.T) {
	lister := JSON(path)

	list, _ := lister.List(metadata.Catalog, "test")

	if len(list) != 1 {
		t.Fail()
	}

	if list[0].GetID() != "Test" {
		t.Fail()
	}
}

func TestJsonTypeLister(t *testing.T) {
	lister := JSON(path)

	list, _ := lister.TypeList(metadata.Catalog)

	if len(list) != 1 {
		t.Fail()
	}

	if list[0] != "test" {
		t.Fail()
	}
}
