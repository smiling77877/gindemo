package syntax

import (
	"log"
	"testing"
)

type Parent struct {
}

func (p Parent) SayHello() {
	log.Println("Hello, I am " + p.Name())
}

func (p Parent) Name() string {
	return "Parent"
}

type Son struct {
	Parent
}

func (s Son) Name() string {
	return "Son"
}

func Test_comm(t *testing.T) {
	var s Son
	s.SayHello()
}
