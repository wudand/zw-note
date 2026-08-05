package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// StringSlice is a []string that scans/values as JSON/JSONB columns.
type StringSlice []string

// Value implements driver.Valuer.
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implements sql.Scanner. Supports []byte (MySQL/pgx) and string (pg jsonb).
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("invalid type for StringSlice: %T", value)
	}
	if len(b) == 0 {
		*s = StringSlice{}
		return nil
	}
	if err := json.Unmarshal(b, s); err != nil {
		return errors.New("invalid JSON for StringSlice")
	}
	return nil
}
