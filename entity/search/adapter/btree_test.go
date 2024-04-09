package searchadapter

import (
	"testing"

	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
)

func TestSetter(t *testing.T) {
	item := indexItem.IndexItem{Key: "1", Values: map[string]any{"name": "name", "email": "name@name.com"}}

	btree := CreateBTreeG()
	btree.Set(item)

	criteria := indexItem.ValueMap{"name": "name"}
	result, _ := btree.Search(criteria)
	if result == nil {
		t.Fatal()
	}

	if result[0] != "1" {
		t.Fatal()
	}
}

func TestRemover(t *testing.T) {
	item := indexItem.IndexItem{Key: "1", Values: map[string]any{"name": "name", "email": "name@name.com"}}

	btree := CreateBTreeG()
	btree.Set(item)
	criteria := indexItem.ValueMap{"name": "name"}
	result, _ := btree.Search(criteria)
	if len(result) != 1 {
		t.Fatal()
	}

	itemToRemove := indexItem.IndexItem{Key: "1", Values: map[string]any{"name": "name", "email": "name@name.com"}}
	btree.Delete(itemToRemove)

	found := btree.Get(itemToRemove)
	if found == true {
		t.Fatal()
	}
}

func TestSearcher(t *testing.T) {
	item1 := indexItem.IndexItem{Key: "1", Values: map[string]any{"name": "name1", "email": "name1@name.com1"}}
	item2 := indexItem.IndexItem{Key: "2", Values: map[string]any{"name": "another name1", "email": "name@name.com"}}
	item3 := indexItem.IndexItem{Key: "3", Values: map[string]any{"name": "name1", "email": "name1@name.com"}}
	item4 := indexItem.IndexItem{Key: "4", Values: map[string]any{"name": "another name2", "email": "name@name.com"}}
	item5 := indexItem.IndexItem{Key: "5", Values: map[string]any{"name": "name3", "email": "name3@name.com"}}

	btree := CreateBTreeG()
	btree.Set(item1)
	btree.Set(item2)
	btree.Set(item3)
	btree.Set(item4)
	btree.Set(item5)

	criteria := indexItem.ValueMap{"name": "name", "email": "name"}
	result, _ := btree.Search(criteria)

	if len(result) != 3 {
		t.Fatal()
	}
	if result[0] != "3" {
		t.Log("1st search: " + result[0])
		t.Fatal()
	}
	if result[1] != "1" {
		t.Log("2nd search: " + result[0])
		t.Fatal()
	}
	if result[2] != "5" {
		t.Log("3rd search: " + result[0])
		t.Fatal()
	}
}
