package utils

type Stack[T any] struct {
	data []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{data: make([]T, 0)}
}

func (q *Stack[T]) Push(item T) {
	q.data = append(q.data, item)
}

func (q *Stack[T]) Pop() T {
	last := len(q.data) - 1
	item := q.data[last]
	q.data = q.data[:last]

	return item
}

func (q *Stack[T]) Peek() T {
	last := len(q.data) - 1
	item := q.data[last]

	return item
}

func (q Stack[T]) Len() int {
	return len(q.data)
}
