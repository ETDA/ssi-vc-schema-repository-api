package models

import (
	"encoding/json"
	"time"
)

type VCSchemaReference struct {
	ID            string           `json:"id" gorm:"column:id"`
	SchemaID      string           `json:"schema_id" gorm:"column:schema_id"`
	ReferenceBody *json.RawMessage `json:"reference_body" gorm:"column:reference_body"`
	Name          string           `json:"name" gorm:"column:name"`
	CreatedAt     *time.Time       `json:"created_at" gorm:"column:created_at"`
	Version       string           `json:"version" gorm:"column:version"`
}

func (m VCSchemaReference) TableName() string {
	return "vc_schema_references"
}
