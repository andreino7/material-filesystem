package main

import (
	"fmt"
	daemon "material/filesystem/daemon/service"
	"material/filesystem/filesystem"
)

func main() {
	// TODO: file system type should be an cli flag
	fmt.Println("Welcome to material file system daemon!")
	fmt.Println("Initializing file system daemon")
	daemon, err := daemon.NewFileSystemDaemon(filesystem.InMemoryFileSystem)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", daemon)

	fmt.Println("Starting daemon")
	err = daemon.Run()
	if err != nil {
		panic(err)
	}
}
