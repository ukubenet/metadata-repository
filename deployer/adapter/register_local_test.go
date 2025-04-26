package adapter

import (
	"os"
	"testing"

	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)



func TestLocalRegisterDeployer(t *testing.T) {
	path, _ := os.Getwd()

	path = path + "/"
	deployer := Local(path)

	// @todo: there is a bug. Register is empty. JUst struct. It does not throw an error that such register does not exist.
	register, _ := metaapi.ReadRegisterMetadata(metadata.State, "RegisterTest")

	err := deployer.RegisterDeploy(register)
	if err != nil {
		t.Fatal(err)
	}

}