package searchadapter

import (
	"strings"

	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
)

type Loop struct {
	list []indexItem.IndexItem
}

func CreateLoop() *Loop {
	return &Loop{
		list: make([]indexItem.IndexItem, 0),
	}
}

func (l *Loop) Delete(item indexItem.IndexItem) error {
	return nil
}

func (l *Loop) Set(item indexItem.IndexItem) error {
	l.list = append(l.list, item)
	return nil
}

func (l *Loop) Get(item indexItem.IndexItem) bool {
	for _, indexIem := range l.list {
		if item.Key == indexIem.Key {
			return true
		}
	}

	return false
}

func (l *Loop) Search(criteria any) (result []indexItem.Key, err error) {
	for _, indexIem := range l.list {
		match := true
		for fieldName, value := range criteria.(map[string]any) {
			if entityAttributeValueString, ok := indexIem.Values[fieldName].(string); ok {
				if !strings.HasPrefix(entityAttributeValueString, value.(string)) {
					match = false
					break
				}
			}
		}
		if match {
			result = append(result, indexItem.Key(indexIem.Key))
		}
	}

	return result, nil
}
