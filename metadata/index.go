package metadata

type IndexType int
type IndexTypeMapType map[IndexType]interface{}

const (
	BTreeG IndexType = 1
	Bleve IndexType = 2
)

func (e IndexType) String() string {
	return [...]string{"BTreeG", "Bleve"}[e-1]
}

type (
	Index struct {
		// EntityType EntityType
		// EntityName string
		// IndexName  string
		IndexType  IndexType `json:"indexType"`
		Attributes []string  `json:"attributes"`
	}
)
