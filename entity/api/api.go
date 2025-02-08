package entityapi

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/ukubenet/metadata-repository/entity"
	entitysearch "github.com/ukubenet/metadata-repository/entity/search"
	indexItem "github.com/ukubenet/metadata-repository/entity/search/item"
	storage "github.com/ukubenet/metadata-repository/entity/storage"
	entityvalidator "github.com/ukubenet/metadata-repository/entity/validator"
	"github.com/ukubenet/metadata-repository/metadata"
	metaapi "github.com/ukubenet/metadata-repository/metadata/api"
)

func ReadEntity(entityType metadata.EntityType, name string, identifier string) (*entity.Entity, error) {
	return storage.ReadEntity(entityType, name, identifier)
}

func DeleteEventTransactions(e *entity.Entity) error {

	if len(e.Transactions) == 0 {
		return nil // @todo we need to restore when we set transactions for events  errors.New("event transactions not defined")
	}

	for balanceName, balanceSpecs := range e.Transactions {
		balanceMap := balanceSpecs.(map[string]any)

		newChangeString, ok := balanceMap["new_change"].(string)
		if !ok {
			return errors.New("new_change not defined. Balance name: " + balanceName)
		}

		newChange, err := strconv.ParseFloat(newChangeString, 64)
		if err != nil {
			return fmt.Errorf("invalid new_change value for balance name: %s", balanceName)
		}

		newEntityReference, ok := balanceMap["new_entity_reference"].(string)
		if !ok {
			return errors.New("new_entity_reference not defined. Balance name: " + balanceName)
		}

		updateBalance(balanceName, newEntityReference, -newChange)
	}

	return nil
}

func PostEventTransactions(e *entity.Entity) error {

	if len(e.Transactions) == 0 {
		return nil // @todo restore it errors.New("event transactions not defined")
	}

	for balanceName, balanceSpecs := range e.Transactions {
		balanceMap := balanceSpecs.(map[string]any)

		var oldChange float64
		oldChangeString, ok := balanceMap["old_change"].(string)
		if ok {
			var err error
			oldChange, err = strconv.ParseFloat(oldChangeString, 64)
			if err != nil {
				oldChange = 0
			}
		} else {
			oldChange = 0
		}

		newChangeString, ok := balanceMap["new_change"].(string)
		if !ok {
			return errors.New("new_change not defined. Balance name: " + balanceName)
		}

		newChange, err := strconv.ParseFloat(newChangeString, 64)
		if err != nil {
			return fmt.Errorf("invalid new_change value for balance name: %s", balanceName)
		}

		newEntityReference, ok := balanceMap["new_entity_reference"].(string)
		if !ok {
			return errors.New("new_entity_reference not defined. Balance name: " + balanceName)
		}

		oldEntityReference, ok := balanceMap["old_entity_reference"].(string)
		if ok {
			updateBalance(balanceName, oldEntityReference, -oldChange)
		}

		updateBalance(balanceName, newEntityReference, newChange)
	}

	return nil
}

func updateBalance(balanceName string, entityReference string, change float64) error {
	if change == 0 {
		return nil
	}

	balanceEntities, err := ReadEntities(metadata.Balance, balanceName)
	if err != nil {
		return fmt.Errorf("balance reference not found: Balance name: %s, Error: %v", balanceName, err)
	}

	foundBalanceEntity := entity.Entity{}
	for _, balanceEntity := range balanceEntities {
		if balanceEntity.Attributes[balanceName].(map[string]any)["reference"].(string) == entityReference {
			foundBalanceEntity = balanceEntity

		}
	}
	if foundBalanceEntity.Identifier != "" {
		balanceString, ok := foundBalanceEntity.Attributes["balance"].(string)
		if !ok {
			return fmt.Errorf("invalid balance value for balance name: %s", balanceName)
		}
		balance, err := strconv.ParseFloat(balanceString, 64)
		if err != nil {
			return fmt.Errorf("invalid balance value for balance name: %s", balanceName)
		}
		balance += change
		foundBalanceEntity.Attributes["balance"] = strconv.FormatFloat(balance, 'f', -1, 64)

		err = PutEntity(metadata.Balance, &foundBalanceEntity)
		if err != nil {
			return err
		}
	} else {
		balanceEntity := &entity.Entity{
			EntityName: balanceName,
			Identifier: uuid.New().String(),
			Attributes: make(map[string]interface{}),
		}

		meta, err := metaapi.ReadMetadata(metadata.Balance, balanceName)
		if err != nil {
			return err
		}
		metaAttributes := meta.GetStructedAttributes()
		for key := range metaAttributes {
			if metaAttributes[key].Type == metadata.ReferenceType {
				refValue, err := RetrieveReferenceByEntityId(metaAttributes[key], entityReference)
				if err != nil {
					return err
				}
				balanceEntity.Attributes[key] = refValue
			} else {
				balanceEntity.Attributes["balance"] = strconv.FormatFloat(change, 'f', -1, 64)
			}
		}

		err = PutEntity(metadata.Balance, balanceEntity)
		if err != nil {
			return err
		}
	}

	return nil
}

func PutEntity(entityType metadata.EntityType, entity *entity.Entity) error {

	if entity.EntityName == "" {
		return errors.New("entity name not defined")
	}
	if entity.Identifier == "" {
		return errors.New("entity identifier not defined")
	}
	if len(entity.Attributes) == 0 {
		return errors.New("entity attributes not defined")
	}
	if err := entityvalidator.ValidateAttributeValues(entityType, entity.EntityName, entity.Attributes); err != nil {
		return err
	}

	if entityType == metadata.Event {
		err := PostEventTransactions(entity)
		if err != nil {
			return err
		}
	}

	factoryWriter := storage.CreateFactory()
	adapter := factoryWriter.CreateAdapter()
	err := adapter.Put(entityType, entity)

	return err
}

func DeleteEntity(entityType metadata.EntityType, name string, identifier string) error {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()

	if entityType == metadata.Event {
		entity, _ := ReadEntity(entityType, name, identifier)
		err := DeleteEventTransactions(entity)
		if err != nil {
			return err
		}
	}

	err := adapter.Delete(entityType, name, identifier)

	return err
}

func ReadEntities(entityType metadata.EntityType, name string) ([]entity.Entity, error) {
	return storage.ReadEntities(entityType, name)
}

func SearchEntities(entityType metadata.EntityType, name string, indexName string, criteria any) ([]indexItem.Key, error) {
	index, ok := entitysearch.Indexes[entityType][name][indexName]
	if !ok {
		return nil, fmt.Errorf("no such index %q, entity: %q", indexName, name)
	}

	list, err := index.Searcher.Search(criteria)

	return list, err
}

func ReadEntityTypes(entityType metadata.EntityType) ([]string, error) {
	storage := storage.CreateFactory()
	adapter := storage.CreateAdapter()

	list, err := adapter.TypeList(entityType)

	return list, err
}

func ReadReference(refSpecs *metadata.ReferenceSpecs, reference string) (map[string]any, error) {
	refEntity, err := ReadEntity(refSpecs.EntityType, refSpecs.Reference, reference)
	if err != nil {
		return nil, err
	}
	view := make(map[string]any, 0)
	for _, viewFieldName := range refSpecs.View {
		view[viewFieldName] = refEntity.Attributes[viewFieldName]
	}

	return map[string]any{"reference": refEntity.Identifier, "type": refSpecs.EntityType.String(), "view": view}, nil
}

func FindSearchIndex(meta metadata.EntityMetadata, createria map[string]any) (*string, error) {
	var indexName *string
	for indName, indexMeta := range meta.SearchCriteria {
		foundIndex := true
		for _, fieldName := range indexMeta.Attributes {
			if _, ok := createria[fieldName]; !ok {
				foundIndex = false
				break
			}
		}

		if foundIndex {
			indexName = &indName
		}
		break
	}

	if indexName == nil {
		return indexName, fmt.Errorf("no search index for defined criteria. Entity: %s, Criteria: %v", meta.EntityName, createria)
	}

	return indexName, nil
}

func FindReference(refSpecs metadata.ReferenceSpecs, createria map[string]any) (map[string]any, error) {
	refmeta, err := metaapi.ReadMetadata(refSpecs.EntityType, refSpecs.Reference)
	if err != nil {
		return nil, err
	}

	indexName, err := FindSearchIndex(*refmeta, createria)
	if err != nil {
		return nil, err
	}

	// indexValueMap := make(indexItem.ValueMap)
	// for _, fieldName := range refmeta.SearchCriteria[*indexName].Attributes {
	// 	indexValueMap[fieldName] = searchCriteriaMap[fieldName]
	// }

	list, err := SearchEntities(refSpecs.EntityType, refSpecs.Reference, *indexName, createria)
	if err != nil {
		return nil, err
	}
	if len(list) != 1 {
		return nil, fmt.Errorf("search by reference should find only 1 record. entity: %s, criteria: %v", refSpecs.Reference, createria)
	}

	return ReadReference(&refSpecs, string(list[0]))
}

func RetrieveReferenceByEntityId(attribute metadata.StructedAttribute, entityId string) (map[string]any, error) {
	refSpecs, ok := attribute.Specs.(metadata.ReferenceSpecs)
	if !ok {
		return nil, fmt.Errorf("meta type should be reference")
	}

	entityType := refSpecs.EntityType

	entity, err := ReadEntity(entityType, refSpecs.Reference, entityId)
	if err != nil {
		return nil, err
	}
	view := make(map[string]interface{})
	for _, viewFieldName := range refSpecs.View {
		view[viewFieldName] = entity.Attributes[viewFieldName]
	}

	return map[string]any{
		"reference": entityId,
		"type":      "reference",
		"view":      view,
	}, nil
}

func ReadReferences(metaSpecs *metadata.ReferenceSpecs) map[string]map[string]any {
	entities, err := ReadEntities(metaSpecs.EntityType, metaSpecs.Reference)
	if err != nil {
		panic(err)
	}

	result := make(map[string]map[string]any)
	for _, entity := range entities {
		view := make(map[string]any)
		for _, viewName := range metaSpecs.View {
			viewValue, ok := entity.Attributes[viewName]
			if ok {
				view[viewName] = viewValue
			}
		}

		result[entity.Identifier] = view
	}

	return result
}
