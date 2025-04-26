package metastorage

import (
	"testing"

	"github.com/ukubenet/metadata-repository/metadata"
)

func TestRegisterReader(t *testing.T) {
	register := new(metadata.RegisterMetadata)
	name := "RegisterTest"

	reader := CreateRegisterFactory(metadata.State)
	adapter := reader.CreateRegisterAdapter()
	err := adapter.RegisterRead(name, register)

	if err != nil {
		t.Fatal(err)
	}

	if register.RegisterName != name {
		t.Fatal("Name", register.RegisterName)
	}

	if len(register.Dimensions) != 2 {
		t.Fatal("Dimensions", register.Dimensions)
	}
}

func TestRegisterReplacer(t *testing.T) {
	register := new(metadata.RegisterMetadata)
	name := "RegisterTest"

	factory := CreateRegisterFactory(metadata.State)
	adapter := factory.CreateRegisterAdapter()
	adapter.RegisterRead(name, register)
	err := adapter.RegisterPut(register)

	if err != nil {
		t.Fatal(err)
	}

	if register.RegisterName != name {
		t.Fatal("RegisterName", register.RegisterName)
	}

	if len(register.Dimensions) != 2 {
		t.Fatal("Dimensions", register.Dimensions)
	}
}

func TestRegisterEraser(t *testing.T) {
	register := new(metadata.RegisterMetadata)
	name := "RegisterTest"
	registerToDelete := "RegisterToDelete"

	factory := CreateRegisterFactory(metadata.State)
	adapter := factory.CreateRegisterAdapter()
	adapter.RegisterRead(name, register)
	register.RegisterName = registerToDelete
	adapter.RegisterPut(register)

	err := adapter.RegisterRead(registerToDelete, register)
	if err != nil {
		t.Fatal(err)
	}

	adapter.RegisterDelete(registerToDelete)

	err = adapter.RegisterRead(registerToDelete, register)
	if err == nil {
		t.Fatal(err)
	}
}

func TestRegisterLister(t *testing.T) {
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
