package stack_test

import (
	"testing"

	"github.com/aacanakin/dag/stack"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("should return true when stack is empty", func(t *testing.T) {
			stack := stack.New()
			stack.Push("X")
			stack.Pop()
			assert.Equal(t, true, stack.IsEmpty())
		})
		t.Run("should return false when stack is not empty", func(t *testing.T) {
			stack := stack.New()
			stack.Push("X")
			assert.Equal(t, false, stack.IsEmpty())
		})
	})

	t.Run("Pop", func(t *testing.T) {
		t.Run("should return an error when stack is empty", func(t *testing.T) {
			stack := stack.New()
			_, err := stack.Pop()
			assert.Error(t, err)
		})
		t.Run("should return the last item when stack is not empty", func(t *testing.T) {
			stack := stack.New()
			stack.Push("X")
			stack.Push("Y")
			item, _ := stack.Pop()
			assert.Equal(t, "Y", item)
		})
	})

	t.Run("Push", func(t *testing.T) {
		t.Run("should push an item to the stack", func(t *testing.T) {
			stack := stack.New()
			stack.Push("X")
			stack.Push("Y")
			assert.Equal(t, false, stack.IsEmpty())

			item, _ := stack.Pop()
			assert.Equal(t, "Y", item)

			item, _ = stack.Pop()
			assert.Equal(t, "X", item)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Run("should return a string representation of the stack", func(t *testing.T) {
			stack := stack.New()
			stack.Push("X")
			stack.Push("Y")
			assert.Equal(t, "<Stack size=2 data=X,Y />", stack.String())
		})
	})
}
