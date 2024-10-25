package queue

import "errors"

func New() *Queue {
	return &Queue{
		items: []string{},
	}
}

type Queue struct {
	items []string
}

func (q *Queue) Size() int {
	return len(q.items)
}

func (q *Queue) Enqueue(data string) {
	q.items = append(q.items, data)
}

func (q *Queue) Pop() (string, error) {
	if len(q.items) == 0 {
		return "", errors.New(`queue is empty, nothing to pop`)
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}
