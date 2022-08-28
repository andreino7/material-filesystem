package main

import (
	"log"
	daemon "material/filesystem/daemon/service"
	"material/filesystem/filesystem"
)

func main() {
	// TODO: file system type should be an cli flag
	log.Println("Initializing file system daemon")
	daemon, err := daemon.NewFileSystemDaemon(filesystem.InMemoryFileSystem)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Println("Starting daemon")
	err = daemon.Run()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
