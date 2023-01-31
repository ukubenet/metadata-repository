package storage

import (
	"testing"

	"github.com/ukubenet/metadata-repository/models"
)

func TestReader(t *testing.T) {
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
