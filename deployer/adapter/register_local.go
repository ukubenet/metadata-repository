package adapter

import (
	"os"

	"github.com/ukubenet/metadata-repository/metadata"
)

const (
	RegisterPath = "register/"
)

func (local *LocalDeployer) RegisterDeploy(registerMeta *metadata.RegisterMetadata) (err error) {

	dir := local.path + RegisterPath + registerMeta.RegisterName
	if e := os.MkdirAll(dir, 0755); !os.IsExist(e) {
		fi, _ := os.Stat(dir)
		if !fi.Mode().IsDir() {
			return e
		}
	}

	return
}

func (local *LocalDeployer) RegisterDelete(registerMeta *metadata.RegisterMetadata) (err error) {

	dir := local.path + RegisterPath + registerMeta.RegisterName
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	return
}
