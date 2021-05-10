package main

import "fmt"

type Student struct {
	Name string
}

func (s *Student) GetName() string {
	return s.Name
}

func NewStudent(name string) *Student {
	return &Student{
		Name: name,
	}
}

func main() {
	s := NewStudent("sakura")
	name := s.GetName()

	fmt.Println(name)
}
