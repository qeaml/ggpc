package ggpc_test

import (
	"fmt"

	"github.com/qeaml/ggpc"
)

func ExampleNewMemory() {
	cache := ggpc.NewMemory[string, int]()

	// set some values
	cache.Set("1", 1)
	cache.Set("2", 2)
	cache.Set("3", 3)

	// and retrieve them
	one, ok := cache.Get("1")
	if !ok {
		fmt.Println("The one is missing?!")
	}
	fmt.Println("One + One is", one+one)

	// you can use defaults, too
	four := cache.GetOrDefault("4", 4)
	fmt.Println("Four + Four is", four+four)

	// no need to close it or anything. it just disappears
}
