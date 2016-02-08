package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/satori/go.uuid"
)

// UUID is a nullable uuid.UUID. It supports SQL and JSON serialization.
// It will marshal to null if null.
type UUID struct {
	UUID  uuid.UUID
	Valid bool
}

// Scan implements the Scanner interface.
func (u *UUID) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case uuid.UUID:
		u.UUID = x
	case []int8, []byte, string:
		err = u.UUID.Scan(x)
	case nil:
		u.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.UUID: %v", value, value)
	}
	u.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (u UUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.UUID.Value()
}

// NewUUID creates a new UUID.
func NewUUID(u uuid.UUID, valid bool) UUID {
	if u == uuid.Nil {
		valid = false
	}

	return UUID{
		UUID:  u,
		Valid: valid,
	}
}

// UUIDFrom creates a new UUID that will always be valid.
func UUIDFrom(u uuid.UUID) UUID {
	return NewUUID(u, true)
}

// UUIDFromPtr creates a new UUID that will be null if t is nil.
func UUIDFromPtr(u *uuid.UUID) UUID {
	if u == nil {
		return NewUUID(uuid.UUID{}, false)
	}
	return NewUUID(*u, true)
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this uuid is null.
func (u UUID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return []byte("null"), nil
	}
	return []byte(`"` + u.UUID.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullUUID and friends)
// and null input.
func (u *UUID) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		err = u.UUID.UnmarshalText([]byte(x))
	case map[string]interface{}:
		ui, uiOK := x["UUID"].(string)
		valid, validOK := x["Valid"].(bool)
		if !uiOK || !validOK {
			return fmt.Errorf(`json: unmarshalling object into Go value of type null.UUID requires key "UUID" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["UUID"], x["Valid"])
		}
		err = u.UUID.UnmarshalText([]byte(ui))
		u.Valid = valid
		return err
	case nil:
		u.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.UUID", reflect.TypeOf(v).Name())
	}
	u.Valid = err == nil
	return err
}

func (u UUID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte("null"), nil
	}
	return u.UUID.MarshalText()
}

func (u *UUID) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		u.Valid = false
		return nil
	}
	if err := u.UUID.UnmarshalText(text); err != nil {
		return err
	}
	u.Valid = true
	return nil
}

// SetValid changes this UUID's value and sets it to be non-null.
func (u *UUID) SetValid(v uuid.UUID) {
	u.UUID = v
	u.Valid = true
}

// Ptr returns a pointer to this UUID's value, or a nil pointer if this UUID is null.
func (u UUID) Ptr() *uuid.UUID {
	if !u.Valid {
		return nil
	}
	return &u.UUID
}
