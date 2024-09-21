package opt

import (
	"encoding/json"
	"fmt"
	"slices"
)

type Option[T any] struct {
	hasValue bool
	Value    T
}

func Value[T any](value T) Option[T] {
	return Option[T]{
		hasValue: true,
		Value:    value,
	}
}

func Empty[T any]() Option[T] {
	return Option[T]{
		hasValue: false,
	}
}

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

func (option Option[T]) HasValue() bool {
	return option.hasValue
}

func (option Option[T]) IsEmpty() bool {
	return !option.hasValue
}

func (option Option[T]) Get() (value T, ok bool) {
	return option.Value, option.hasValue
}

func (option *Option[T]) Put(value T) {
	option.hasValue = true
	option.Value = value
}

func (option *Option[T]) Clear() {
	*option = Option[T]{hasValue: false}
}

func (option Option[T]) String() string {
	if option.hasValue {
		return fmt.Sprint(option.Value)
	} else {
		return "<empty>"
	}
}

func (option Option[T]) MarshalJSON() ([]byte, error) {
	if option.hasValue {
		return json.Marshal(option.Value)
	} else {
		// We don't return the jsonNull slice from below here, since we don't want to risk it being
		// modified
		return []byte("null"), nil
	}
}

var jsonNull = []byte("null")

func (option *Option[T]) UnmarshalJSON(bytes []byte) error {
	if slices.Equal(bytes, jsonNull) {
		option.hasValue = false
		return nil
	}

	option.hasValue = true
	return json.Unmarshal(bytes, &option.Value)
}
