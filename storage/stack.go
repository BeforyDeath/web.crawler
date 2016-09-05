package storage

var Stack stack

type stack struct {
	Nodes []string
}

func (s *stack) Init() {
	for k, v := range Items {
		if v.Status.Code == 0 || v.Status.Code == 504 {
			s.Push(k)
		}
	}
}

func (s *stack) Push(key string) {
	s.Nodes = append(s.Nodes, key)
}

func (s *stack) Pop() (hash string) {
	if len(s.Nodes) == 0 {
		return
	}
	hash, s.Nodes = s.Nodes[len(s.Nodes)-1], s.Nodes[:len(s.Nodes)-1]
	return
}
