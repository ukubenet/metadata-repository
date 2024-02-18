package entitystorage

import (
	"testing"

	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/entity"
	"github.com/ukubenet/metadata-repository/metadata"
)

func TestMain(m *testing.M) {
	config.LoadConfig("../../config", "test")
	m.Run()
}

func TestReader(t *testing.T) {
	entity := new(entity.CatalogEntity)
	name := "test"
	identifier := "Test"

	reader := CreateFactory(metadata.Catalog)
	adapter := reader.CreateAdapter()
	err := adapter.Read(name, identifier, entity)

	if err != nil {
		t.Fatal(err)
	}

	if entity.EntityName != name {
		t.Fatal("Name", entity.EntityName)
	}

	if entity.Identifier != identifier {
		t.Fatal("Identifier", entity.EntityName)
	}

	if len(entity.Attributes) != 2 {
		t.Fatal("Attributes", entity.Attributes)
	}
}

func TestReplacer(t *testing.T) {
	entity := new(entity.CatalogEntity)
	name := "test"
	identifier := "Test"

	factory := CreateFactory(metadata.Catalog)
	adapter := factory.CreateAdapter()
	adapter.Read(name, identifier, entity)
	err := adapter.Put(entity)

	if err != nil {
		t.Fatal(err)
	}

	if entity.EntityName != name {
		t.Fatal("Name", entity.EntityName)
	}

	if entity.Identifier != identifier {
		t.Fatal("Identifier", entity.Identifier)
	}

	if len(entity.Attributes) != 2 {
		t.Fatal("Attributes", entity.Attributes)
	}
}

func TestEraser(t *testing.T) {
	entity := new(entity.CatalogEntity)
	name := "test"
	identifier := "EntityToDelete"

	factory := CreateFactory(metadata.Catalog)
	adapter := factory.CreateAdapter()
	entity.EntityName = name
	entity.Identifier = identifier
	adapter.Put(entity)

	err := adapter.Read(name, identifier, entity)
	if err != nil {
		t.Fatal(err)
	}

	adapter.Delete(name, identifier)

	err = adapter.Read(name, identifier, entity)
	if err == nil {
		t.Fatal(err)
	}
}

func TestLister(t *testing.T) {
	factory := CreateFactory(metadata.Catalog)
	adapter := factory.CreateAdapter()
	list := []entity.CatalogEntity{}
	name := "test"
	adapter.List(name, &list)

	if len(list) != 1 {
		t.Fatal("List quantity", len(list))
	}

	if list[0].Identifier != "Test" {
		t.Fatal("List[0].Identifier is ", list[0].Identifier)
	}
}
