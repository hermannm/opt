package opt

import (
	"encoding/json"
	"fmt"
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
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
}

func (option *Option[T]) UnmarshalJSON(bytes []byte) error {
	isNull := len(bytes) == 4 &&
		bytes[0] == 'n' &&
		bytes[1] == 'u' &&
		bytes[2] == 'l' &&
		bytes[3] == 'l'

	if isNull {
		option.hasValue = false
		return nil
	} else {
		option.hasValue = true
		return json.Unmarshal(bytes, &option.Value)
	}
}
