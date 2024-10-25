// Package opt provides [Option], a container that either has a value or is empty.
package opt

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Option is a container that either has a value, or is empty. You construct an option with [Value],
// [Empty] or [FromPointer]. You can check if an option contains a value or is empty with
// [Option.HasValue] or [Option.IsEmpty], and access the value through the [Option.Value] field.
//
// The zero value of Option is an empty option.
//
// An empty option marshals to `null` in JSON, and a `null` JSON value unmarshals to an empty
// option.
type Option[T any] struct {
	hasValue bool
	// Before accessing Value, you should check if it is present with [Option.HasValue].
	Value T
}

// Value creates an [Option] that contains the given value.
func Value[T any](value T) Option[T] {
	return Option[T]{
		hasValue: true,
		Value:    value,
	}
}

// Empty creates an empty [Option].
func Empty[T any]() Option[T] {
	return Option[T]{
		hasValue: false,
	}
}

// FromPointer creates an [Option] with the value pointed to by the given pointer, dereferencing it.
// If the pointer is nil, an empty option is returned.
func FromPointer[T any](pointer *T) Option[T] {
	if pointer == nil {
		return Option[T]{
			hasValue: false,
		}
	} else {
		return Option[T]{
			hasValue: true,
			Value:    *pointer,
		}
	}
}

// HasValue returns true if the option contains a value.
func (option Option[T]) HasValue() bool {
	return option.hasValue
}

// IsEmpty returns true if the option does not contain a value.
func (option Option[T]) IsEmpty() bool {
	return !option.hasValue
}

// Get returns the value of the option, and an `ok` flag that is true if the option contained a
// value, and false if it is empty. You should only use the returned value if `ok` is true.
func (option Option[T]) Get() (value T, ok bool) {
	return option.Value, option.hasValue
}

// GetOrDefault returns the option's value if present, or the given default value if the option is
// empty.
func (option Option[T]) GetOrDefault(defaultValue T) T {
	if option.hasValue {
		return option.Value
	} else {
		return defaultValue
	}
}

// Put replaces the current value of the option, if any, with the given value. After this call,
// [Option.HasValue] will return true.
func (option *Option[T]) Put(value T) {
	option.hasValue = true
	option.Value = value
}

// Clear removes the current value of the option, if any. After this call, [Option.IsEmpty] will
// return true.
func (option *Option[T]) Clear() {
	*option = Option[T]{hasValue: false}
}

// ToPointer returns nil if the option is empty, otherwise it returns a pointer to the option's
// value.
//
// It is meant to be used for compatibility with libraries that use pointers for optional values.
func (option Option[T]) ToPointer() *T {
	if option.hasValue {
		return &option.Value
	} else {
		return nil
	}
}

// FromSQL creates an [Option] from the given [sql.Null] value. A null SQL value becomes an empty
// option, and a non-null SQL value becomes an option containing the value.
func FromSQL[T any](sql sql.Null[T]) Option[T] {
	if sql.Valid {
		return Option[T]{hasValue: true, Value: sql.V}
	} else {
		return Option[T]{hasValue: false}
	}
}

// ToSQL converts the option to an [sql.Null]. An empty option becomes null.
func (option Option[T]) ToSQL() sql.Null[T] {
	return sql.Null[T]{Valid: option.hasValue, V: option.Value}
}

// String returns the string representation of the option's value. If the option is empty, it
// returns the string `<empty>` (similar to the string representation `<nil>` for nil pointers).
func (option Option[T]) String() string {
	if option.hasValue {
		return fmt.Sprint(option.Value)
	} else {
		return "<empty>"
	}
}

// MarshalJSON implements the [json.Marshaler] interface for [Option]. If the option contains a
// value, it marshals that value. If the option is empty, it marshals to `null`.
func (option Option[T]) MarshalJSON() ([]byte, error) {
	if option.hasValue {
		return json.Marshal(option.Value)
	} else {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
}

// UnmarshalJSON implements the [json.Unmarshaler] interface for [Option]. If the given JSON value
// is `null`, it unmarshals to an empty option. Otherwise, it tries to unmarshal to the value
// contained by the option.
func (option *Option[T]) UnmarshalJSON(jsonValue []byte) error {
	isNull := len(jsonValue) == 4 &&
		jsonValue[0] == 'n' &&
		jsonValue[1] == 'u' &&
		jsonValue[2] == 'l' &&
		jsonValue[3] == 'l'

	if isNull {
		option.hasValue = false
		return nil
	} else {
		option.hasValue = true
		return json.Unmarshal(jsonValue, &option.Value)
	}
}
