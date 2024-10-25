package stack

import (
	"errors"
	"fmt"
	"strings"
)

func New() *Stack {
	return &Stack{items: []string{}}
}

type Stack struct {
	items []string
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Pop() (string, error) {
	if s.IsEmpty() {
		return "", errors.New(`could not pop item, stack is empty`)
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, nil
}

func (s *Stack) Push(data string) {
	s.items = append(s.items, data)
}

func (s *Stack) String() string {
	return fmt.Sprintf("<Stack size=%d data=%s />", len(s.items), strings.Join(s.items, ","))
}
