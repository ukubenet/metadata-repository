package adapter

import (
	"os"
	"testing"

	"github.com/ukubenet/metadata-repository/metadata"
)

var path string

func TestMain(m *testing.M) {
	path, _ = os.Getwd()
	path += "/"
	m.Run()
}

func TestJsonReader(t *testing.T) {
	reader := JSON(path)

	candidate := new(metadata.EntityMetadata)

	err := reader.Read("test/Test", candidate)
	if err != nil {
		t.Fatal(err)
	}

	if candidate.EntityName != "test/Test" {
		t.Fail()
	}

	if candidate.Attributes == nil {
		t.Fail()
	}
}

func TestJsonReplacer(t *testing.T) {
	inserter := JSON(path)
	var attributes []metadata.Attribute = []metadata.Attribute{
		{"name": "string_attribute", "type": "string"},
		{"name": "number_attribute", "type": "number"},
	}

	var candidate *metadata.EntityMetadata = &metadata.EntityMetadata{
		EntityName: "test/Test",
		Attributes: attributes,
	}

	inserter.Put(candidate)
}

func TestJsonLister(t *testing.T) {
	lister := JSON(path)
	list := []string{}

	lister.List(&list)

	if len(list) != 1 {
		t.Fail()
	}

	if list[0] != "Test" {
		t.Fail()
	}
}
