package models

import (
	"database/sql/driver"
	"encoding/json"
)

// NullString est un type personnalisé pour gérer les chaînes nullables en JSON
type NullString struct {
	String string
	Valid  bool
}

// Scan implémente l'interface sql.Scanner
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	switch v := value.(type) {
	case string:
		ns.String = v
	case []byte:
		ns.String = string(v)
	}
	return nil
}

// Value implémente l'interface driver.Valuer
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON implémente l'interface json.Marshaler
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implémente l'interface json.Unmarshaler
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.String)
}

