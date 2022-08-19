package main

import (
	"fmt"
	"material/filesystem/filesystem"
)

func main() {
	filesystem.NewFileSystem(filesystem.InMemoryFileSystem)
	fmt.Println("Hello, world.")
}
