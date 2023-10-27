package adapter

import (
	"os"

	"github.com/ukubenet/metadata-repository/metadata"
)

type (
	LocalDeployer struct {
		path string
	}
)

func Local(path string) *LocalDeployer {
	return &LocalDeployer{
		path: path,
	}
}

func (local *LocalDeployer) Deploy(entitymeta *metadata.EntityMetadata) (err error) {

	dir := local.path + entitymeta.EntityName
	if e := os.Mkdir(dir, 0755); !os.IsExist(e) {
		fi, _ := os.Stat(dir)
		if !fi.Mode().IsDir() {
			return e
		}
	}

	return
}
