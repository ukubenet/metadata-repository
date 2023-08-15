package entitystorage

import (
	"testing"

	"github.com/ukubenet/metadata-repository/entity"
)

func TestMain(m *testing.M) {
	set_env("test")
	m.Run()
}

func TestReader(t *testing.T) {
	set_env("test")
	entity := new(entity.Entity)
	name := "test"
	identifier := "Test"

	reader := CreateFactory()
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
	entity := new(entity.Entity)
	name := "test"
	identifier := "Test"

	factory := CreateFactory()
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
	entity := new(entity.Entity)
	name := "test"
	identifier := "EntityToDelete"

	factory := CreateFactory()
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
	factory := CreateFactory()
	adapter := factory.CreateAdapter()
	list := []entity.Entity{}
	name := "test"
	adapter.List(name, &list)

	if len(list) != 1 {
		t.Fatal("List quantity", len(list))
	}

	if list[0].Identifier != "Test" {
		t.Fatal("List[0].Identifier is ", list[0].Identifier)
	}
}
