package models

import (
	"context"
	"database/sql"
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

	query = `select	a.uuid, a.name, t.name, a.entity_uuid
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
			&attr.Type.Name,
			&attr.EntityUuid,
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
func (m *DBModel) All(id int) ([]*Entity, error) {
	return nil, nil
}
