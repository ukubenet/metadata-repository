package entitysearch

import (
	"testing"

	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/entity"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func TestMain(m *testing.M) {
	config.LoadConfig("../../config", "test")
	m.Run()

}

func TestLoadIndex(t *testing.T) {

	var entities []entity.Entity

	e := entity.CatalogEntity{
		Metadata: entity.Metadata{
			EntityName: "User",
			Identifier: "1",
			Attributes: entity.AttributeValues{
				"name":  "User 1",
				"email": "email@gmail.com",
			},
		},
	}
	entities = append(entities, e)

	e = entity.CatalogEntity{
		Metadata: entity.Metadata{
			EntityName: "User",
			Identifier: "2",
			Attributes: entity.AttributeValues{
				"name":  "User",
				"email": "email2@gmail.com",
			},
		},
	}
	entities = append(entities, e)

	e = entity.CatalogEntity{
		Metadata: entity.Metadata{
			EntityName: "User",
			Identifier: "3",
			Attributes: entity.AttributeValues{
				"name":  "User",
				"email": "email3@gmail.com",
			},
		},
	}
	entities = append(entities, e)

	indexMeta := metadata.Index{
		IndexType:  metadata.BTreeG,
		Attributes: []string{"name", "email"},
	}

	index := LoadIndex(entities, indexMeta)

	indexItem := indexItem.IndexItem{
		Key: "3", Values: indexItem.ValueMap{
			"name":  "User",
			"email": "email3@gmail.com",
		},
	}

	found := index.Getter.Get(indexItem)
	if found == false {
		t.Fail()
	}

	indexItem.Key = "5"
	found = index.Getter.Get(indexItem)
	if found == true {
		t.Fail()
	}
}

func TestLoadIndexes(t *testing.T) {
	meta, _ := metaapi.ReadMetadata(metadata.Catalog, "test")
	entityIndexes := LoadEntityIndexes(metadata.Catalog, *meta)

	indItem := indexItem.IndexItem{
		Key: "test2", Values: indexItem.ValueMap{
			"name":  "name0",
			"email": "name0@name.com",
		},
	}

	found := entityIndexes["btree"].Getter.Get(indItem)
	if found == false {
		t.Fail()
	}

	criteria := indexItem.ValueMap{
		"name":  "name",
		"email": "name",
	}
	result, _ := entityIndexes["btree"].Searcher.Search(criteria)

	if len(result) != 3 {
		t.Fatal()
	}
}
