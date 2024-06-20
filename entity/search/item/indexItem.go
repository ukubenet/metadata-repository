package indexItem

type Key string
type ValueMap map[string]any

// @todo consider IndexItem to be map where Key is an entity identifier and ValueMap
type IndexItem struct {
	Key    Key
	Values ValueMap
}
