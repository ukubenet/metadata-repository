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

func TestJsonCatalogReader(t *testing.T) {
	reader := JSON(path)

	entity, err := reader.Read(metadata.Catalog, "test", "Test")
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

func TestJsonCatalogReplacer(t *testing.T) {
	inserter := JSON(path)
	var attributes entity.AttributeValues = entity.AttributeValues{
		"string_attribute": "string",
		"number_attribute": 100,
	}

	var candidate = &entity.Entity{
		EntityName: "test",
		Attributes: attributes,
		Identifier: "Test",
	}

	inserter.Put(metadata.Catalog, candidate)
}

func TestJsonCatalogLister(t *testing.T) {
	lister := JSON(path)

	list, _ := lister.List(metadata.Catalog, "test")

	if len(list) != 1 {
		t.Fail()
	}

	if list[0].Identifier != "Test" {
		t.Fail()
	}
}

func TestJsonCatalogTypeLister(t *testing.T) {
	lister := JSON(path)

	list, _ := lister.TypeList(metadata.Catalog)

	if len(list) != 1 {
		t.Fail()
	}

	if list[0] != "test" {
		t.Fail()
	}
}

func TestJsonEventReader(t *testing.T) {
	reader := JSON(path)

	e, err := reader.Read(metadata.Event, "test", "Test")
	if err != nil {
		t.Fatal(err)
	}

	if e.EntityName != "test" {
		t.Fail()
	}

	if e.Identifier != "Test" {
		t.Fail()
	}

	if e.Attributes == nil {
		t.Fail()
	}
}

func TestJsonEventReplacer(t *testing.T) {
	inserter := JSON(path)
	var attributes entity.AttributeValues = entity.AttributeValues{
		"string_attribute": "string",
		"number_attribute": 100,
	}

	var candidate *entity.Entity = &entity.Entity{
		EntityName: "test",
		Attributes: attributes,
		Identifier: "Test",
	}

	inserter.Put(metadata.Event, candidate)
}

func TestJsonEventLister(t *testing.T) {
	lister := JSON(path)

	list, _ := lister.List(metadata.Event, "test")

	if len(list) != 1 {
		t.Fail()
	}

	if list[0].Identifier != "Test" {
		t.Fail()
	}
}

func TestJsonEventTypeLister(t *testing.T) {
	lister := JSON(path)

	list, _ := lister.TypeList(metadata.Event)

	if len(list) != 1 {
		t.Fail()
	}

	if list[0] != "test" {
		t.Fail()
	}
}
