package adapter

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ukubenet/metadata-repository/metadata"
)

const (
	RegisterPath = "register/"
)

func (local *LocalDeployer) RegisterDeploy(registerMeta *metadata.RegisterMetadata) (err error) {
	dbPath := local.path + RegisterPath + registerMeta.RegisterName + ".db"

	// Ensure the directory exists
	dir := local.path + RegisterPath
	if e := os.MkdirAll(dir, 0755); !os.IsExist(e) {
		fi, _ := os.Stat(dir)
		if !fi.Mode().IsDir() {
			return e
		}
	}

	// Create the SQLite database and tables
	if err := createRegisterTables(dbPath, registerMeta); err != nil {
		return err
	}

	return nil
}

func createRegisterTables(dbPath string, registerMeta *metadata.RegisterMetadata) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Extract auxiliary columns from reference dimensions
	auxiliaryColumns := registerMeta.Auxiliaries
	for name, dimension := range registerMeta.Dimensions {
		if dimension["type"] == metadata.ReferenceType {
			refSpecs := dimension["specs"].(metadata.ReferenceSpecs)
			for _, viewColumn := range refSpecs.View {
				uniqueColumnName := name + "_" + viewColumn
				auxiliaryColumns[uniqueColumnName] = map[string]interface{}{
					"type": metadata.StringType,
				}
			}
		}
	}

	// Create transactions table
	if err := createTable(db, "transactions", registerMeta.Dimensions, registerMeta.Facts, auxiliaryColumns); err != nil {
		return err
	}

	// Create state table (use Facts directly for now)
	if err := createTable(db, "state", registerMeta.Dimensions, registerMeta.Facts, auxiliaryColumns); err != nil {
		return err
	}

	return nil
}

func createTable(db *sql.DB, tableName string, dimensions, factsOrStates, auxiliaries metadata.Attributes) error {
	columnDefs := ""

	addAttributes := func(attrs metadata.Attributes) {
		for attrName, attrSpecs := range attrs {
			attrType, ok := attrSpecs["type"]
			if !ok {
				panic("attribute doesn't have type")
			}

			if attrType == metadata.ReferenceType {
				// ReferenceType columns are strings referring to another entity
				columnDefs += fmt.Sprintf("%s TEXT,\n", attrName)
			} else {
				columnDefs += fmt.Sprintf("%s %s,\n", attrName, goTypeToSQLite(attrType.(string)))
			}
		}
	}

	addAttributes(dimensions)
	addAttributes(factsOrStates)
	addAttributes(auxiliaries)

	// Remove trailing comma from column definitions
	if len(columnDefs) > 0 {
		columnDefs = columnDefs[:len(columnDefs)-2] // Remove last comma and newline
	}

	createStmt := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s
	);`, tableName, columnDefs)

	fmt.Println("Executing SQL:\n", createStmt)

	_, err := db.Exec(createStmt)
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
