package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONB is a custom type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

// Value converts JSONB to a database value
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan converts database value to JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}

	result := make(JSONB)
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID         int       `json:"id"`
	UserID     *int      `json:"user_id"`
	Action     string    `json:"action"` // create, update, delete, approve, reject
	EntityType string    `json:"entity_type"`
	EntityID   int       `json:"entity_id"`
	OldValue   JSONB     `json:"old_value,omitempty"`
	NewValue   JSONB     `json:"new_value,omitempty"`
	IPAddress  *string   `json:"ip_address,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
