package metadata

type IndexType int
type IndexTypeMapType map[IndexType]interface{}

const (
	Loop   IndexType = 1
	BTreeG IndexType = 2
	Bleve  IndexType = 3
)

func (e IndexType) String() string {
	return [...]string{"Loop", "BTreeG", "Bleve"}[e-1]
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
