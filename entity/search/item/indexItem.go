package indexItem

type Key string
type ValueMap map[string]any

type IndexItem struct {
	Key    Key
	Values ValueMap
}
