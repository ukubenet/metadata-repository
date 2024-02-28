package adapter

import (
	"os"
	"testing"
	"time"

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

func TestJsonCatalogReplacer(t *testing.T) {
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

func TestJsonCatalogLister(t *testing.T) {
	lister := JSON(path)

	list, _ := lister.List(metadata.Catalog, "test")

	if len(list) != 1 {
		t.Fail()
	}

	if list[0].GetID() != "Test" {
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

	if e.GetName() != "test" {
		t.Fail()
	}

	if e.GetID() != "Test" {
		t.Fail()
	}

	if e.GetAttributes() == nil {
		t.Fail()
	}

	eventEntity, ok := e.(*entity.EventEntity)
	if !ok {
		t.Fail()
	}

	if eventEntity.EventTime.IsZero() {
		t.Fail()
	}
}

func TestJsonEventReplacer(t *testing.T) {
	inserter := JSON(path)
	var attributes entity.AttributeValues = entity.AttributeValues{
		"string_attribute": "string",
		"number_attribute": 100,
	}

	var candidate = entity.EventEntity{
		Metadata: entity.Metadata{
			EntityName: "test",
			Attributes: attributes,
			Identifier: "Test",
		},
		EventTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	inserter.Put(candidate)
}

func TestJsonEventLister(t *testing.T) {
	lister := JSON(path)

	list, _ := lister.List(metadata.Event, "test")

	if len(list) != 1 {
		t.Fail()
	}

	if list[0].GetID() != "Test" {
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
