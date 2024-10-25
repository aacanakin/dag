package set

import (
	"fmt"
	"strings"
)

type Set struct {
	items map[string]bool
}

func New() *Set {
	return &Set{
		make(map[string]bool),
	}
}

func (s *Set) Add(id string) bool {
	if s.Has(id) {
		return false
	}
	s.items[id] = true
	return true
}

func (s *Set) Size() int {
	return len(s.items)
}

func (s *Set) Has(id string) bool {
	return s.items[id]
}

func (s *Set) List() []string {
	list := []string{}
	for k := range s.items {
		list = append(list, k)
	}
	return list
}

func (s *Set) Remove(id string) bool {
	if !s.Has(id) {
		return false
	}
	delete(s.items, id)
	return true
}

func (s *Set) String() string {
	return fmt.Sprintf("<Set size=%d data=%s />", s.Size(), strings.Join(s.List(), ","))
}
