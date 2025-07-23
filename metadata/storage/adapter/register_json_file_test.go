package adapter

import (
	"testing"

	"github.com/ukubenet/metadata-repository/metadata"
)


func TestJsonRegisterReader(t *testing.T) {
	reader := JSON(path)

	candidate := new(metadata.RegisterMetadata)

	err := reader.RegisterRead("test/RegisterTest", candidate)
	if err != nil {
		t.Fatal(err)
	}

	if candidate.RegisterName != "test/RegisterTest" {
		t.Fail()
	}

	if candidate.Dimensions == nil {
		t.Fail()
	}
	if candidate.Facts == nil {
		t.Fail()
	}
}
func TestJsonRegisterReplacer(t *testing.T) {
	inserter := JSON(path)
	var dimensions metadata.Attributes = metadata.Attributes{
		"string_attribute": {"type": "string"},
		"number_attribute": {"type": "number"},
	}
	var facts metadata.Attributes = metadata.Attributes{
		"number_fact": map[string]any{"type": "number"},
	}

	var candidate *metadata.RegisterMetadata = &metadata.RegisterMetadata{
		RegisterName: "test/RegisterTest",
		Dimensions:  dimensions,
		Facts:       facts,
	}

	inserter.RegisterPut(candidate)

}

func TestRegisterJsonLister(t *testing.T) {
	lister := JSON(path)
	list := []string{}

	lister.List(&list)

	if len(list) != 2 {
		t.Fail()
	}

	if list[0] != "RegisterTest" {
		t.Fail()
	}
}
