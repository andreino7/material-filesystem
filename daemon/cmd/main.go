package main

import (
	"log"
	daemon "material/filesystem/daemon/service"
	"material/filesystem/filesystem"
	"os"
)

const defaultPort = "2222"

func main() {
	// TODO: file system type should be an cli flag
	log.Println("Initializing file system daemon")
	daemon, err := daemon.NewFileSystemDaemon(filesystem.InMemoryFileSystem)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Println("Starting daemon")
	port := os.Getenv("FS_DAEMON_PORT")
	if port == "" {
		port = defaultPort
	}
	err = daemon.Run(port)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
