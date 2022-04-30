package ggpc_test

import (
	"fmt"
	"os"

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

func ExampleNewStored() {
	// temporary file, use an actual file to place your cache in instead!
	f, err := os.CreateTemp("", "my_cache")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// create a new stored cache
	// NOTE: this will not load the cache's state from the file, you have to use
	//       cache.Load() to load it or use LoadStored() instead
	cache := ggpc.NewStored[string, int](f)

	/*
		do stuff with your cache
	*/

	// write the changes to storage (concurrency safe)
	cache.Save()
}
