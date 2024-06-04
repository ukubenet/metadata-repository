package searchadapter

import (
	"strings"

	"github.com/google/btree"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
)

func GetBtree() *btree.BTreeG[indexItem.IndexItem] {
	return &btree.BTreeG[indexItem.IndexItem]{}
}

// byKeys is a comparison function that compares item keys and returns true
// when a is less than b.
func byKeys(a, b indexItem.IndexItem) bool {
	return a.Key < b.Key
}

// byVals is a comparison function that compares item values and returns true
// when a is less than b.
func byVals(a, b indexItem.IndexItem) bool {
	for index := range a.Values {
		_, ok := a.Values[index].(string)
		if ok {
			if a.Values[index].(string) < b.Values[index].(string) {
				return true
			}
			if a.Values[index].(string) > b.Values[index].(string) {
				return false
			}
		} else {
			panic("not implemented")
		}
	}

	// Both vals are equal so we should fall though
	// and let the key comparison take over.
	return byKeys(a, b)
}

type BTreeGIndex struct {
	btree *btree.BTreeG[indexItem.IndexItem]
}

func CreateBTreeG() *BTreeGIndex {
	return &BTreeGIndex{
		btree: btree.NewG[indexItem.IndexItem](3, byVals),
	}
}

func (b *BTreeGIndex) Delete(item indexItem.IndexItem) error {
	b.btree.Delete(item)

	return nil
}

func (b *BTreeGIndex) Set(item indexItem.IndexItem) error {
	b.btree.ReplaceOrInsert(item)

	return nil
}

func (b *BTreeGIndex) Get(item indexItem.IndexItem) bool {
	_, found := b.btree.Get(item)

	return found
}

func (b *BTreeGIndex) Search(criteria any) (result []indexItem.Key, err error) {
	values := criteria.(indexItem.ValueMap)
	item := indexItem.IndexItem{Key: "", Values: values}

	b.btree.AscendGreaterOrEqual(item, func(current indexItem.IndexItem) bool {
		for indexName, value := range criteria.(indexItem.ValueMap) {
			if !strings.HasPrefix(current.Values[indexName].(string), value.(string)) {
				return false
			}
		}

		result = append(result, current.Key)

		return true
	})

	return result, nil
}
