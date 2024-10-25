package set_test

import (
	"testing"

	"github.com/aacanakin/dag/set"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		t.Run("should return a new set", func(t *testing.T) {
			set := set.New()
			assert.Equal(t, 0, set.Size())
		})
	})

	t.Run("Add", func(t *testing.T) {
		t.Run("should add an item to the set", func(t *testing.T) {
			set := set.New()
			set.Add("X")
			assert.Equal(t, 1, set.Size())
		})

		t.Run("should not add an item to the set if it already exists", func(t *testing.T) {
			set := set.New()
			set.Add("X")
			set.Add("X")
			assert.Equal(t, 1, set.Size())
		})
	})

	t.Run("Size", func(t *testing.T) {
		t.Run("should return the size of the set", func(t *testing.T) {
			set := set.New()
			set.Add("X")
			set.Add("Y")
			assert.Equal(t, 2, set.Size())
		})
	})

	t.Run("Has", func(t *testing.T) {
		t.Run("should return true if the item exists in the set", func(t *testing.T) {
			set := set.New()
			set.Add("X")
			assert.Equal(t, true, set.Has("X"))
		})

		t.Run("should return false if the item does not exist in the set", func(t *testing.T) {
			set := set.New()
			assert.Equal(t, false, set.Has("X"))
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Run("should remove an item from the set", func(t *testing.T) {
			set := set.New()
			set.Add("X")
			set.Remove("X")
			assert.Equal(t, 0, set.Size())
		})

		t.Run("should not remove an item from the set if it does not exist", func(t *testing.T) {
			set := set.New()
			set.Remove("X")
			assert.Equal(t, 0, set.Size())
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("should return a list of items in the set", func(t *testing.T) {
			set := set.New()
			set.Add("X")
			set.Add("Y")
			assert.ElementsMatch(t, []string{"X", "Y"}, set.List())
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Run("should return a string representation of the set", func(t *testing.T) {
			set := set.New()
			set.Add("X")
			set.Add("Y")

			expected := []string{"<Set size=2 data=X,Y />", "<Set size=2 data=Y,X />"}
			assert.Contains(t, expected, set.String())
		})
	})
}
