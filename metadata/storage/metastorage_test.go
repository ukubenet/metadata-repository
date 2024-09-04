package metastorage

import (
	"testing"

	config "github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
)

func TestMain(m *testing.M) {
	config.LoadConfig("../../config", "test")
	global.SetAppName("test")
	m.Run()
}

func TestReader(t *testing.T) {
	entity := new(metadata.EntityMetadata)
	name := "Test"

	reader := CreateFactory(metadata.Catalog)
	adapter := reader.CreateAdapter()
	err := adapter.Read(name, entity)

	if err != nil {
		t.Fatal(err)
	}

	if entity.EntityName != name {
		t.Fatal("Name", entity.EntityName)
	}

	if len(entity.Attributes) != 1 {
		t.Fatal("Attributes", entity.Attributes)
	}
}

func TestReplacer(t *testing.T) {
	entity := new(metadata.EntityMetadata)
	name := "Test"

	factory := CreateFactory(metadata.Catalog)
	adapter := factory.CreateAdapter()
	adapter.Read(name, entity)
	err := adapter.Put(entity)

	if err != nil {
		t.Fatal(err)
	}

	if entity.EntityName != name {
		t.Fatal("Name", entity.EntityName)
	}

	if len(entity.Attributes) != 1 {
		t.Fatal("Attributes", entity.Attributes)
	}
}

func TestEraser(t *testing.T) {
	entity := new(metadata.EntityMetadata)
	name := "Test"
	entityToDelete := "EntityToDelete"

	factory := CreateFactory(metadata.Catalog)
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
	factory := CreateFactory(metadata.Catalog)
	adapter := factory.CreateAdapter()
	list := []string{}
	adapter.List(&list)

	if len(list) != 1 {
		t.Fatal("List quantity", len(list))
	}

	if list[0] != "Test" {
		t.Fatal("List[1] is ", list[1])
	}
}
