package main

import (
	"fmt"
	"log"
)

type aaa struct {
	tmp int
}

func (a *aaa) name() {
	fmt.Println(a.tmp)
}
func main() {

	a := []*aaa{&aaa{tmp: 5}, &aaa{tmp: 6}, &aaa{tmp: 7}}

	for i, v := range a {
		v.tmp = i
	}

	for _, v := range a {
		log.Println(v.tmp)
	}
}
