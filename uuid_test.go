package null

import (
	"encoding/json"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/satori/go.uuid"
)

var (
	uuidString     = "e0e8d676-ce47-11e5-ab30-625662870761"
	uuidJSON       = []byte(`"` + uuidString + `"`)
	nullUUIDJSON   = []byte(`null`)
	uuidValue, _   = uuid.FromString(uuidString)
	uuidObject     = []byte(`{"UUID":"e0e8d676-ce47-11e5-ab30-625662870761","Valid":true}`)
	uuidNullObject = []byte(`{"UUID":"e0e8d676-ce47-11e5-ab30-625662870761","Valid":false}`)
	uuidBadObject  = []byte(`{"hello": "world"}`)
)

func TestUnmarshalUUIDJSON(t *testing.T) {
	var ui UUID
	err := json.Unmarshal(uuidJSON, &ui)
	maybePanic(err)
	assertUUID(t, ui, "UnmarshalJSON() json")

	var null UUID
	err = json.Unmarshal(nullUUIDJSON, &null)
	maybePanic(err)
	assertNullUUID(t, null, "null uuid json")

	var fromObject UUID
	err = json.Unmarshal(uuidObject, &fromObject)
	maybePanic(err)
	assertUUID(t, fromObject, "uuid from object json")

	var nullFromObj UUID
	err = json.Unmarshal(uuidNullObject, &nullFromObj)
	maybePanic(err)
	assertNullUUID(t, nullFromObj, "null from object json")

	var invalid UUID
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUUID(t, invalid, "invalid from object json")

	var bad UUID
	err = json.Unmarshal(uuidBadObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullUUID(t, bad, "bad from object json")

	var wrongType UUID
	err = json.Unmarshal(intJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullUUID(t, wrongType, "wrong type object json")
}

func TestUnmarshalUUIDText(t *testing.T) {
	u := UUIDFrom(uuidValue)
	txt, err := u.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, uuidString, "marshal text")

	var unmarshal UUID
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertUUID(t, unmarshal, "unmarshal text")

	u = UUIDFrom(uuid.Nil)
	txt, err = u.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, string(nullUUIDJSON), "marshal text")

	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertNullUUID(t, unmarshal, "unmarshal text")

	var null UUID
	err = null.UnmarshalText(nullJSON)
	maybePanic(err)
	assertNullUUID(t, null, "unmarshal null text")
	txt, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, string(nullJSON), "marshal null text")

	var invalid UUID
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullUUID(t, invalid, "bad string")
}

func TestMarshalUUID(t *testing.T) {
	ti := UUIDFrom(uuidValue)
	data, err := json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(uuidJSON), "non-empty json marshal")

	ti.Valid = false
	data, err = json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
}

func TestUUIDFrom(t *testing.T) {
	ti := UUIDFrom(uuidValue)
	assertUUID(t, ti, "UUIDFrom() uuid.UUID")
}

func TestUUIDFromPtr(t *testing.T) {
	ti := UUIDFromPtr(&uuidValue)
	assertUUID(t, ti, "UUIDFromPtr() uuid")

	null := UUIDFromPtr(nil)
	assertNullUUID(t, null, "UUIDFromPtr(nil)")
}

func TestUUIDSetValid(t *testing.T) {
	var ti uuid.UUID
	change := NewUUID(ti, false)
	assertNullUUID(t, change, "SetValid()")
	change.SetValid(uuidValue)
	assertUUID(t, change, "SetValid()")
}

func TestUUIDPointer(t *testing.T) {
	ti := UUIDFrom(uuidValue)
	ptr := ti.Ptr()
	if *ptr != uuidValue {
		t.Errorf("bad %s uuid: %#v ≠ %v\n", "pointer", ptr, uuidValue)
	}

	var nt uuid.UUID
	null := NewUUID(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s uuid: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUUIDScanValue(t *testing.T) {
	var ui UUID
	err := ui.Scan(uuidValue)
	maybePanic(err)
	assertUUID(t, ui, "scanned uuid")
	if v, err := ui.Value(); v != uuidString || err != nil {
		pp.Println(v)
		pp.Println(uuidValue)
		t.Error("bad value or err:", v, err)
	}

	var null UUID
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUUID(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong UUID
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullUUID(t, wrong, "scanned wrong")
}

func assertUUID(t *testing.T, ui UUID, from string) {
	if ui.UUID != uuidValue {
		t.Errorf("bad %v uuid: %v ≠ %v\n", from, ui.UUID, uuidValue)
	}
	if !ui.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUUID(t *testing.T, ui UUID, from string) {
	if ui.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
