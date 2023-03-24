package storage

import (
	"testing"

	"github.com/ukubenet/metadata-repository/models"
)

func TestMain(m *testing.M) {
	set_env("test")
	m.Run()
}

func TestReader(t *testing.T) {
	set_env("test")
	entity := new(models.Entity)
	name := "test/Test"

	reader := CreateFactory()
	adapter := reader.CreateAdapter()
	err := adapter.Read(name, entity)

	if err != nil {
		t.Fatal(err)
	}

	if entity.EntityName != name {
		t.Fatal("Name", entity.EntityName)
	}

	if len(entity.Attributes) != 2 {
		t.Fatal("Attributes", entity.Attributes)
	}
}

func TestReplacer(t *testing.T) {
	entity := new(models.Entity)
	name := "test/Test"

	factory := CreateFactory()
	adapter := factory.CreateAdapter()
	adapter.Read(name, entity)
	err := adapter.Put(entity)

	if err != nil {
		t.Fatal(err)
	}

	if entity.EntityName != name {
		t.Fatal("Name", entity.EntityName)
	}

	if len(entity.Attributes) != 2 {
		t.Fatal("Attributes", entity.Attributes)
	}
}

func TestEraser(t *testing.T) {
	entity := new(models.Entity)
	name := "test/Test"
	entityToDelete := "test/EntityToDelete"

	factory := CreateFactory()
	adapter := factory.CreateAdapter()
	adapter.Read(name, entity)
	entity.EntityName = entityToDelete
	adapter.Put(entity)

	err := adapter.Read(entityToDelete, entity)
	if err != nil {
		t.Fatal(err)
	}

	adapter.Delete(entityToDelete)

	err = adapter.Read(entityToDelete, entity)
	if err == nil {
		t.Fatal(err)
	}
}

func TestLister(t *testing.T) {
	factory := CreateFactory()
	adapter := factory.CreateAdapter()
	list := []string{}
	adapter.List(&list)

	if len(list) != 2 {
		t.Fatal("List quantity", len(list))
	}

	if list[1] != "Test" {
		t.Fatal("List[1] is ", list[1])
	}
}
