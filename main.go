package main

import "log"

func main() {

	s := stack{}
	s.Push("s1")
	s.Push("f2")
	s.Push("s2")
	s.Push("a23")

	log.Println(s.Pop())
	log.Println(s.Pop())
	log.Println(s.Pop())
	log.Println(s.Pop())
	log.Println(s.Pop())

log.Println(s)



}
