package opt

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

func (option Option[T]) HasValue() bool {
	return option.hasValue
}

func (option Option[T]) IsEmpty() bool {
	return !option.hasValue
}

func (option Option[T]) Get() (value T, ok bool) {
	return option.Value, option.hasValue
}