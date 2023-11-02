package deployer

import (
	"testing"

	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
	metastorage "github.com/ukubenet/metadata-repository/metadata/storage"
)

func TestMain(m *testing.M) {
	SetEnv("test")
	m.Run()
}

func TestDeployer(t *testing.T) {

	metastorage.SetEnv("test")
	entity, err := metaapi.ReadMetadata("test")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateFactory()
	adapter := deployer.CreateAdapter()
	err = adapter.Deploy(entity)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {

	metastorage.SetEnv("test")
	entity, err := metaapi.ReadMetadata("test")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateFactory()
	adapter := deployer.CreateAdapter()
	err = adapter.Delete(entity)

	if err != nil {
		t.Fatal(err)
	}
}
