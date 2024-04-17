package entitysearch

import (
	searchadapter "github.com/ukubenet/metadata-repository/entity/search/adapter"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
	"github.com/ukubenet/metadata-repository/metadata"
)

type EntityIndexes map[string]*EntityIndex

type (
	Searcher interface {
		Search(any) ([]indexItem.Key, error)
	}

	Remover interface {
		Delete(indexItem.IndexItem) error
	}

	Setter interface {
		Set(indexItem.IndexItem) error
	}

	Getter interface {
		Get(indexItem.IndexItem) bool
	}

	Adapter struct {
		factory *EntityIndex
	}

	EntityIndex struct {
		Searcher Searcher
		Remover  Remover
		Setter   Setter
		Getter   Getter
	}
)

func (f *EntityIndex) UseSearcher(searcher Searcher) {
	f.Searcher = searcher
}

func (f *EntityIndex) UseRemover(remover Remover) {
	f.Remover = remover
}

func (f *EntityIndex) UseSetter(setter Setter) {
	f.Setter = setter
}

func (f *EntityIndex) UseGetter(getter Getter) {
	f.Getter = getter
}

func (f *EntityIndex) Use(i interface{}) {

	if searcher, ok := i.(Searcher); ok {
		f.UseSearcher(searcher)
	}

	if remover, ok := i.(Remover); ok {
		f.UseRemover(remover)
	}

	if setter, ok := i.(Setter); ok {
		f.UseSetter(setter)
	}

	if getter, ok := i.(Getter); ok {
		f.UseGetter(getter)
	}

}

func NewFactory() *EntityIndex {
	f := new(EntityIndex)
	return f
}

func (p *Adapter) Search(criteria any) (list []indexItem.Key, err error) {
	searcher := p.factory.Searcher
	if list, err = searcher.Search(criteria); err != nil {
		return
	}

	return
}

func (p *Adapter) Delete(item indexItem.IndexItem) (err error) {
	remover := p.factory.Remover
	if err = remover.Delete(item); err != nil {
		return
	}

	return
}

func (p *Adapter) Set(item indexItem.IndexItem) (err error) {
	setter := p.factory.Setter
	if err = setter.Set(item); err != nil {
		return
	}

	return
}

func (p *Adapter) Get(item indexItem.IndexItem) (found bool) {
	getter := p.factory.Getter
	return getter.Get(item)
}

func CreateFactory(indexType metadata.IndexType) *EntityIndex {
	factory := NewFactory()
	if indexType == metadata.BTreeG {
		factory.Use(searchadapter.CreateBTreeG())
	}

	return factory
}
