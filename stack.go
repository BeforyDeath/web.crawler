package main

type stack struct {
	nodes []string
}

func (s *stack) Push(key string) {
	s.nodes = append(s.nodes, key)
}

func (s *stack) Pop() (node string) {
	if len(s.nodes) == 0 {
		return node
	}
	node, s.nodes = s.nodes[len(s.nodes) - 1], s.nodes[:len(s.nodes) - 1]
	return node
}