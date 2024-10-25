package queue_test

import (
	"testing"

	"github.com/aacanakin/dag/queue"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	t.Run("Size", func(t *testing.T) {
		t.Run("should return 0 when queue is empty", func(t *testing.T) {
			queue := queue.New()
			assert.Equal(t, 0, queue.Size())
		})
		t.Run("should return the number of items in the queue", func(t *testing.T) {
			queue := queue.New()
			queue.Enqueue("X")
			queue.Enqueue("Y")
			assert.Equal(t, 2, queue.Size())
		})
	})

	t.Run("Enqueue", func(t *testing.T) {
		t.Run("should add an item to the queue", func(t *testing.T) {
			queue := queue.New()
			queue.Enqueue("X")
			queue.Enqueue("Y")
			assert.Equal(t, 2, queue.Size())
		})
	})

	t.Run("Pop", func(t *testing.T) {
		t.Run("should return an error when queue is empty", func(t *testing.T) {
			queue := queue.New()
			_, err := queue.Pop()
			assert.Error(t, err)
		})
		t.Run("should return the first item when queue is not empty", func(t *testing.T) {
			queue := queue.New()
			queue.Enqueue("X")
			queue.Enqueue("Y")
			item, _ := queue.Pop()
			assert.Equal(t, "X", item)
		})
	})
}
