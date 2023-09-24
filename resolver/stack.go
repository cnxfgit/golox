package resolver

type stack struct {
	data []map[string]bool
}

func newStack() *stack {
	return &stack{
		data: make([]map[string]bool, 0),
	}
}

func (s *stack) peek() (map[string]bool, bool) {
	if len(s.data) == 0 {
		return nil, false
	}

	return s.data[len(s.data)-1], true
}

func (s *stack) push(element map[string]bool) {
	s.data = append(s.data, element)
}

func (s *stack) pop() (map[string]bool, bool) {
	if len(s.data) == 0 {
		return nil, false
	}

	element := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return element, true
}

func (s *stack) isEmpty() bool {
	return len(s.data) == 0
}

func (s *stack) len() int {
	return len(s.data)
}

func (s *stack) get(i int) map[string]bool {
	return s.data[i]
}
