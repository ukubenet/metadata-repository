package adapter

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/ukubenet/metadata-repository/metadata"
)


func TestRegisterDeploy(t *testing.T) {
	localDeployer := &LocalDeployer{path: "./testdata/"}
	registerMeta := &metadata.RegisterMetadata{
		RegisterName: "test_register",
		Dimensions: metadata.Attributes{
			"dim1": {"type": metadata.StringType},
		},
		Facts: metadata.Attributes{
			"fact1": {"type": metadata.IntegerType},
		},
		Auxiliaries: metadata.Attributes{
			"aux1": {"type": metadata.BooleanType},
		},
	}

	err := localDeployer.RegisterDeploy(registerMeta)
	assert.NoError(t, err)

	// Verify database and tables
	dbPath := "./testdata/register/test_register.db"
	_, err = os.Stat(dbPath)
	assert.NoError(t, err)

	db, err := sql.Open("sqlite3", dbPath)
	assert.NoError(t, err)
	defer db.Close()

	// Check transactions table
	rows, err := db.Query("PRAGMA table_info(transactions)")
	assert.NoError(t, err)
	defer rows.Close()

	expectedColumns := map[string]string{
		"dim1":  "TEXT",
		"fact1": "INTEGER",
		"aux1":  "BOOLEAN",
	}

	actualColumns := make(map[string]string)
	for rows.Next() {
		var (
			cid        int
			name       string
			typeName   string
			notnull    int
			defaultVal interface{}
			pk         int
		)
		err = rows.Scan(&cid, &name, &typeName, &notnull, &defaultVal, &pk)
		assert.NoError(t, err)
		actualColumns[name] = typeName
	}

	assert.Equal(t, expectedColumns, actualColumns)

	// Clean up
	os.RemoveAll("./testdata/")
}

func TestRegisterDelete(t *testing.T) {
	localDeployer := &LocalDeployer{path: "./testdata/"}
	registerMeta := &metadata.RegisterMetadata{
		RegisterName: "test_register",
	}

	// Create dummy directory
	dirPath := "./testdata/register/test_register"
	err := os.MkdirAll(dirPath, 0755)
	assert.NoError(t, err)

	err = localDeployer.RegisterDelete(registerMeta)
	assert.NoError(t, err)

	// Verify deletion
	_, err = os.Stat(dirPath)
	assert.True(t, os.IsNotExist(err))
}
