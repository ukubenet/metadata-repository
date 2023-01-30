package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// Get returns one entity and error, if any
func (m *DBModel) Get(id string) (*Entity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select uuid, entity_name, created_at, updated_at from entity where uuid = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var entity Entity

	err := row.Scan(
		&entity.UUID,
		&entity.EntityName,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	query = `select	a.uuid, a.name, t.name, a.entity_uuid, a.created_at, a.updated_at
			from attribute a
			inner join entity e on a.entity_uuid = e.uuid
			inner join type t on t.uuid = a.type_uuid
			where a.entity_uuid = $1
    `
	rows, _ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close()

	var attributes []Attribute
	for rows.Next() {
		var attr Attribute
		err := rows.Scan(
			&attr.UUID,
			&attr.Name,
			&attr.Type,
			&attr.EntityUuid,
			&attr.CreatedAt,
			&attr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, attr)
	}

	entity.Attributes = attributes

	return &entity, nil
}

// All returns all entities and error, if any
func (m *DBModel) All(attribute ...string) ([]*Entity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(attribute) > 0 {
		where = fmt.Sprintf("where uuid in (select entity_uuid from attribute where uuid = '%s')", attribute[0])
	}
	query := fmt.Sprintf(`select uuid, entity_name, created_at, updated_at from entity %s order by entity_name`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*Entity

	for rows.Next() {
		var entity Entity
		err := rows.Scan(
			&entity.UUID,
			&entity.EntityName,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		attributeQuery := `select a.uuid, a.name, t.name, a.entity_uuid, a.created_at, a.updated_at
			from attribute a
			inner join entity e on a.entity_uuid = e.uuid
			inner join type t on t.uuid = a.type_uuid
			where a.entity_uuid = $1
    `
		attributeRows, _ := m.DB.QueryContext(ctx, attributeQuery, entity.UUID)

		var attributes []Attribute
		for attributeRows.Next() {
			var attr Attribute
			err := attributeRows.Scan(
				&attr.UUID,
				&attr.Name,
				&attr.Type,
				&attr.EntityUuid,
				&attr.CreatedAt,
				&attr.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
			attributes = append(attributes, attr)
		}
		attributeRows.Close()

		entity.Attributes = attributes
		entities = append(entities, &entity)
	}

	return entities, nil
}

func (m *DBModel) AttributesAll() ([]*Attribute, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select uuid, name, created_at from attribute order by name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributes []*Attribute

	for rows.Next() {
		var a Attribute
		err := rows.Scan(
			&a.UUID,
			&a.Name,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, &a)
	}

	return attributes, nil
}
