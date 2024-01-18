package deployer

import (
	"testing"

	config "github.com/ukubenet/metadata-repository/config"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func TestMain(m *testing.M) {
	config.LoadConfig("../config", "test")
	m.Run()
}

func TestDeployer(t *testing.T) {
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
