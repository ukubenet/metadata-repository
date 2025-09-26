package deployer

import (
	"testing"

	config "github.com/ukubenet/metadata-repository/config"
	global "github.com/ukubenet/metadata-repository/global"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func TestMain(m *testing.M) {
	config.LoadConfig("../config", "test")
	global.SetAppName("test")
	m.Run()
}

func TestDeployer(t *testing.T) {
	entity, err := metaapi.ReadMetadata(metadata.Catalog, "test")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateFactory(metadata.Catalog)
	adapter := deployer.CreateAdapter()
	err = adapter.Deploy(entity)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	entity, err := metaapi.ReadMetadata(metadata.Catalog, "test")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateFactory(metadata.Catalog)
	adapter := deployer.CreateAdapter()
	err = adapter.Delete(entity)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDeploy(t *testing.T) {
	register, err := metaapi.ReadRegisterMetadata(metadata.State, "RegisterTest")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateRegisterFactory(metadata.State)
	adapter := deployer.CreateAdapter()
	err = adapter.RegisterDeploy(register)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDelete(t *testing.T) {
	register, err := metaapi.ReadRegisterMetadata(metadata.State, "RegisterTest")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateRegisterFactory(metadata.State)
	adapter := deployer.CreateAdapter()
	err = adapter.RegisterDelete(register)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDeployBalance(t *testing.T) {
	register, err := metaapi.ReadRegisterMetadata(metadata.State, "RegisterTest")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateRegisterFactory(metadata.Balance)
	adapter := deployer.CreateAdapter()
	err = adapter.RegisterDeploy(register)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDeleteBalance(t *testing.T) {
	register, err := metaapi.ReadRegisterMetadata(metadata.State, "RegisterTest")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateRegisterFactory(metadata.Balance)
	adapter := deployer.CreateAdapter()
	err = adapter.RegisterDelete(register)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDeployAccumulator(t *testing.T) {
	register, err := metaapi.ReadRegisterMetadata(metadata.State, "RegisterTest")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateRegisterFactory(metadata.Accumulator)
	adapter := deployer.CreateAdapter()
	err = adapter.RegisterDeploy(register)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDeleteAccumulator(t *testing.T) {
	register, err := metaapi.ReadRegisterMetadata(metadata.State, "RegisterTest")
	if err != nil {
		t.Fatal(err)
	}

	deployer := CreateRegisterFactory(metadata.Accumulator)
	adapter := deployer.CreateAdapter()
	err = adapter.RegisterDelete(register)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterDeployError(t *testing.T) {
	// Test with non-existent register metadata to ensure error handling
	_, err := metaapi.ReadRegisterMetadata(metadata.State, "NonExistentRegister")
	if err == nil {
		t.Fatal("Expected error when reading non-existent register metadata, got nil")
	}
}
