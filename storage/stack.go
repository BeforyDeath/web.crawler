package storage

var Stack stack

type stack struct {
	nodes []string
}

func (s *stack) Init() {
	for k, v := range Items {
		if v.Status.Code == 0 || v.Status.Code == 504 {
			s.Push(k)
		}
	}
}

func (s *stack) Push(key string) {
	s.nodes = append(s.nodes, key)
}

func (s *stack) Pop() (hash string) {
	if len(s.nodes) == 0 {
		return
	}
	hash, s.nodes = s.nodes[len(s.nodes)-1], s.nodes[:len(s.nodes)-1]
	return
}
