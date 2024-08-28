package chatgpt

import (
	"reflect"
	"testing"

	"github.com/ukubenet/metadata-repository/config"
	"github.com/ukubenet/metadata-repository/entity"
	entitysearch "github.com/ukubenet/metadata-repository/entity/search"
	"github.com/ukubenet/metadata-repository/metadata"
)

func TestGetEntityAttributeValuesFromResponse(t *testing.T) {
	tests := []struct {
		name     string
		response map[string]any
		meta     *metadata.EntityMetadata
		want     entity.AttributeValues
		wantErr  bool
	}{
		{
			name: "Simple attributes",
			response: map[string]any{
				"name":  "Test Entity",
				"value": 42,
			},
			meta: &metadata.EntityMetadata{
				Attributes: metadata.Attributes{
					"name":  {"type": "string"},
					"value": {"type": "int"},
				},
			},
			want: entity.AttributeValues{
				"name":  "Test Entity",
				"value": 42,
			},
			wantErr: false,
		},
		{
			name: "With reference attribute",
			response: map[string]any{
				"name": "Test Entity",
				"ref": map[string]any{
					"name": "Referenced Entity ChatTest",
				},
			},
			meta: &metadata.EntityMetadata{
				Attributes: metadata.Attributes{
					"name": {"type": "string"},
					"ref": {
						"type":          metadata.ReferenceType,
						"reference":     "reference",
						"referenceType": "catalog",
						"view":          []any{"name"},
					},
				},
			},
			want: entity.AttributeValues{
				"name": "Test Entity",
				"ref": map[string]any{
					"reference": "id",
					"type":      "Catalog",
					"view":      map[string]any{"name": "Referenced Entity ChatTest"},
				},
			},
			wantErr: false,
		},
		{
			name: "Missing reference attribute",
			response: map[string]any{
				"name": "Test Entity",
			},
			meta: &metadata.EntityMetadata{
				Attributes: metadata.Attributes{
					"name": {"type": "string"},
					"ref": {
						"type":          metadata.ReferenceType,
						"reference":     "reference",
						"referenceType": "catalog",
						"view":          []any{"name"},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	config.LoadConfig("../config", "test")
	entitysearch.Indexes.LoadAllIndexes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEntityAttributeValuesFromResponse(tt.response, metadata.GetStructedAttributes(tt.meta.Attributes))
			if (err != nil) != tt.wantErr {
				t.Errorf("getEntityAttributeValuesFromResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(
				entity.GetStructedAttributeValues(got, metadata.GetStructedAttributes(tt.meta.Attributes)),
				entity.GetStructedAttributeValues(tt.want, metadata.GetStructedAttributes(tt.meta.Attributes)),
			) {
				t.Errorf("getEntityAttributeValuesFromResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
