package adapter

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ukubenet/metadata-repository/metadata"
	_ "github.com/mattn/go-sqlite3"
)

const (
	RegisterPath = "register/"
)

func (local *LocalDeployer) RegisterDeploy(registerMeta *metadata.RegisterMetadata) (err error) {

	dir := local.path + RegisterPath
	if e := os.MkdirAll(dir, 0755); !os.IsExist(e) {
		fi, _ := os.Stat(dir)
		if !fi.Mode().IsDir() {
			return e
		}
	}
	createRegisterTable(dir, registerMeta)

	return
}

func createRegisterTable(dbPath string, registerMeta *metadata.RegisterMetadata) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	columnDefs := ""

	addAttributes := func(attrs metadata.Attributes) {
		for attrName, attrSpecs := range attrs {
			attrType, ok := attrSpecs["type"]
			if !ok {
				panic("attribute doesn't have type")
			}
			columnDefs += fmt.Sprintf("%s %s,\n", attrName, goTypeToSQLite(attrType.(string)))
		}
	}

	addAttributes(registerMeta.Dimensions)
	addAttributes(registerMeta.Auxiliaries)
	addAttributes(registerMeta.Facts)

	createStmt := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s
	);`, registerMeta.RegisterName, columnDefs)

	fmt.Println("Executing SQL:\n", createStmt)

	_, err = db.Exec(createStmt)
	return err
}

func goTypeToSQLite(goType string) string {
	switch goType {
	case metadata.IntegerType:
		return "INTEGER"
	case metadata.NumberType:
		return "REAL"
	case metadata.StringType:
		return "TEXT"
	case metadata.BooleanType:
		return "BOOLEAN"
	default:
		return "TEXT"
	}
}	

func (local *LocalDeployer) RegisterDelete(registerMeta *metadata.RegisterMetadata) (err error) {

	dir := local.path + RegisterPath + registerMeta.RegisterName
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	return
}
