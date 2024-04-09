package entitysearch

import (
	"github.com/ukubenet/metadata-repository/entity"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
	entitystorage "github.com/ukubenet/metadata-repository/entity/storage"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

type IndexList map[metadata.EntityType]map[string]EntityIndexes

var Indexes = CreateIndexes()

func CreateIndexes() IndexList {
	return make(IndexList)
}

func (indexes IndexList) LoadAllIndexes() {
	indexes[metadata.Catalog] = LoadEntityTypeIndexes(metadata.Catalog)
	indexes[metadata.Event] = LoadEntityTypeIndexes(metadata.Event)
}

func LoadEntityTypeIndexes(entityType metadata.EntityType) map[string]EntityIndexes {
	entityIndexMap := make(map[string]EntityIndexes)
	list, _ := entitystorage.ReadEntityTypeList(entityType)
	for _, entityName := range list {
		meta, _ := metaapi.ReadMetadata(entityType, entityName)
		entityIndexMap[entityName] = LoadEntityIndexes(entityType, *meta)
	}

	return entityIndexMap
}

func LoadEntityIndexes(entityType metadata.EntityType, entityMeta metadata.EntityMetadata) EntityIndexes {
	var list, _ = entitystorage.ReadEntities(entityType, entityMeta.EntityName)
	entityIndexes := make(EntityIndexes)

	for indexName, indexMeta := range entityMeta.SearchCriteria {
		entityIndexes[indexName] = LoadIndex(list, indexMeta)
	}

	return entityIndexes
}

func LoadIndex(list []entity.Entity, indexMeta metadata.Index) *EntityIndex {
	indexer := CreateFactory(indexMeta.IndexType)
	setter := indexer.Setter

	for _, entity := range list {
		if entity.GetAttributes() == nil {
			continue
		}
		item := GetIndexItem(indexMeta, entity)
		setter.Set(item)
	}

	return indexer
}

func GetIndexItem(indexMeta metadata.Index, entity entity.Entity) indexItem.IndexItem {
	values := GetIndexItemValues(indexMeta, entity)
	return indexItem.IndexItem{Key: indexItem.Key(entity.GetID()), Values: values}
}

func GetIndexItemValues(indexMeta metadata.Index, entity entity.Entity) indexItem.ValueMap {
	values := make(indexItem.ValueMap)
	for _, attrName := range indexMeta.Attributes {
		values[attrName] = entity.GetAttributes()[attrName]
	}

	return values
}

func (indexes IndexList) DeleteEntityIndexItems(
	entityType metadata.EntityType,
	metadata metadata.EntityMetadata,
	entity entity.Entity,
) {
	for indexName, indexMeta := range metadata.SearchCriteria {
		item := GetIndexItem(indexMeta, entity)
		indexes[entityType][metadata.EntityName][indexName].Remover.Delete(item)
	}
}

func (indexes IndexList) SetEntityIndexItems(
	entityType metadata.EntityType,
	metadata metadata.EntityMetadata,
	entity entity.Entity,
) {
	for indexName, indexMeta := range metadata.SearchCriteria {
		item := GetIndexItem(indexMeta, entity)
		indexes[entityType][metadata.EntityName][indexName].Setter.Set(item)
	}
}
