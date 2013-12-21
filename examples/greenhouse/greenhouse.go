package main

import (
	"fmt"
	"bitbucket.com/cswank/gogadgets"
)

type Greenhouse struct {
	gogadgets.GoGadget
}

func (g *Greenhouse) Start(in <-chan gogadgets.Message, out chan<- gogadgets.Message) {
	
}


func main() {
	a := gogadgets.App{}
	g := &Greenhouse{}
	a.AddGadget(g)
	fmt.Println(a)
}
