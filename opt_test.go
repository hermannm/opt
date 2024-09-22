package opt_test

import (
	"encoding/json"
	"testing"

	"hermannm.dev/opt"
)

func TestValue(t *testing.T) {
	option := opt.Value("test")

	if !option.HasValue() {
		t.Error("HasValue: want true")
	}
	if option.IsEmpty() {
		t.Error("IsEmpty: want false")
	}
	if option.Value != "test" {
		t.Errorf("Value = %s; want 'test'", option.Value)
	}
	if value, ok := option.Get(); !ok {
		t.Errorf("Get() = %s, %t; want 'test', true", value, ok)
	}
}

func TestEmpty(t *testing.T) {
	option := opt.Empty[string]()

	if !option.IsEmpty() {
		t.Error("IsEmpty: want true")
	}
	if option.HasValue() {
		t.Error("HasValue: want false")
	}
	if option.Value != "" {
		t.Errorf("Value = %s; want zero value ''", option.Value)
	}
	if value, ok := option.Get(); ok {
		t.Errorf("Get() = %s, %t; want '', false", value, ok)
	}
}

func TestFromPointerValue(t *testing.T) {
	value := "test"
	pointer := &value
	option := opt.FromPointer(pointer)

	if !option.HasValue() {
		t.Error("HasValue: want true")
	}
	if option.Value != value {
		t.Errorf("Value = %s; want %s", option.Value, value)
	}
}

func TestFromNilPointer(t *testing.T) {
	var pointer *string
	option := opt.FromPointer(pointer)

	if !option.IsEmpty() {
		t.Error("IsEmpty: want true")
	}
	if option.Value != "" {
		t.Errorf("Value = %s; want zero value ''", option.Value)
	}
}

func TestPut(t *testing.T) {
	option := opt.Empty[string]()
	option.Put("test")

	if !option.HasValue() {
		t.Error("HasValue: want true")
	}
	if option.Value != "test" {
		t.Errorf("Value = %s; want 'test'", option.Value)
	}
}

func TestClear(t *testing.T) {
	option := opt.Value("test")
	option.Clear()

	if !option.IsEmpty() {
		t.Error("IsEmpty: want true")
	}
	if option.Value != "" {
		t.Errorf("Value = %s; want zero value ''", option.Value)
	}
}

func TestZeroValue(t *testing.T) {
	var option opt.Option[string]

	if !option.IsEmpty() {
		t.Error("IsEmpty: want true")
	}
}

type stringer struct {
	value string
}

func (stringer stringer) String() string {
	return "Value: " + stringer.value
}

func TestValueString(t *testing.T) {
	option := opt.Value(stringer{"test"})
	string := option.String()

	expected := "Value: test"
	if string != expected {
		t.Errorf("String() = %s; want %s", string, expected)
	}
}

func TestEmptyString(t *testing.T) {
	option := opt.Empty[stringer]()
	string := option.String()

	expected := "<empty>"
	if string != expected {
		t.Errorf("String() = %s; want %s", string, expected)
	}
}

type jsonObject struct {
	Field1 opt.Option[string] `json:"field1"`
	Field2 opt.Option[string] `json:"field2"`
}

func TestMarshalJSON(t *testing.T) {
	object := jsonObject{
		Field1: opt.Value("test"),
		Field2: opt.Empty[string](),
	}

	jsonValue, err := json.Marshal(object)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}

	expected := `{"field1":"test","field2":null}`
	if string(jsonValue) != expected {
		t.Errorf("json.Marshal() = %s; want %s", string(jsonValue), expected)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	jsonValue := []byte(`{"field1":"test","field2":null}`)

	var object jsonObject
	if err := json.Unmarshal(jsonValue, &object); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if !object.Field1.HasValue() {
		t.Error("Field1.HasValue: want true")
	}
	if object.Field1.Value != "test" {
		t.Errorf("Field1.Value = %s; want 'test'", object.Field1.Value)
	}

	if !object.Field2.IsEmpty() {
		t.Error("Field2.IsEmpty: want true")
	}
	if object.Field2.Value != "" {
		t.Errorf("Field2.Value = %s; want zero value ''", object.Field2.Value)
	}
}
