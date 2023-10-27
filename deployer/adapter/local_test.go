package adapter

import (
	"os"
	"testing"

	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

var path string

func TestMain(m *testing.M) {
	path, _ = os.Getwd()
	path += "/"
	m.Run()
}

func TestLocalDeployer(t *testing.T) {
	deployer := Local(path)

	entity, _ := metaapi.ReadMetadata("test")

	err := deployer.Deploy(entity)
	if err != nil {
		t.Fatal(err)
	}
}
