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
	name := "test"
	identifier := "Test"

	reader := CreateFactory()
	adapter := reader.CreateAdapter()
	e, err := adapter.Read(metadata.Catalog, name, identifier)

	if err != nil {
		t.Fatal(err)
	}

	if e.EntityName != name {
		t.Fatal("Name", e.EntityName)
	}

	if e.Identifier != identifier {
		t.Fatal("Identifier", e.Identifier)
	}

	if len(e.Attributes) != 2 {
		t.Fatal("Attributes", e.Attributes)
	}
}

func TestReplacer(t *testing.T) {
	name := "test"
	identifier := "Test"

	factory := CreateFactory()
	adapter := factory.CreateAdapter()
	entity, err := adapter.Read(metadata.Catalog, name, identifier)
	if err != nil {
		t.Fatal(err)
	}

	err = adapter.Put(metadata.Catalog, entity)

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
	name := "test"
	identifier := "EntityToDelete"

	var entity = &entity.Entity{
		EntityName: name,
		Identifier: identifier,
	}

	factory := CreateFactory()
	adapter := factory.CreateAdapter()
	entity.EntityName = name
	entity.Identifier = identifier
	adapter.Put(metadata.Catalog, entity)

	_, err := adapter.Read(metadata.Catalog, name, identifier)
	if err != nil {
		t.Fatal(err)
	}

	adapter.Delete(metadata.Catalog, name, identifier)

	_, err = adapter.Read(metadata.Catalog, name, identifier)
	if err == nil {
		t.Fatal(err)
	}
}

func TestLister(t *testing.T) {
	factory := CreateFactory()
	adapter := factory.CreateAdapter()
	name := "test"
	list, _ := adapter.List(metadata.Catalog, name)

	if len(list) != 1 {
		t.Fatal("List quantity", len(list))
	}

	if list[0].Identifier != "Test" {
		t.Fatal("List[0].Identifier is ", list[0].Identifier)
	}
}
