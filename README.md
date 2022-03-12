# ggpc

Caches for Go, generics included!

## About

`ggpc` is a simple cache library which heavily utilises the new Go 1.18
generics. This arose from a need to keep session tokens around in between
restarts of a web server, but now you can use this thing as well.

These caches are thread-safe, so feel free to concurrently *go ham*.

You can view the documentation [here][docs].

## Two types of caches

There's the in-memory cache, which only lives for as long as the program itself
is running, and the stored/persistent cache, which saves and loads it's state
from someplace else when needed. In most cases that place is going to be a file.

Currently, the caches are stored in JSON files. There is no way to change this.
(it's just a cache, not a database)

## Installation

```shell
go get github.com/qeaml/ggpc
```

## Example Usage

In-memory cache:

```go
package main

import (
 "fmt"

 "github.com/qeaml/ggpc"
)

func main() {
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
```

Stored/persistent cache:

```go
// TODO
```

[docs]: https://pkg.go.dev/github.com/qeaml/ggpc
