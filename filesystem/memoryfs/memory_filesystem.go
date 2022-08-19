package memoryfs

import "sync"

type MemoryFileSystem struct {
	mutex sync.RWMutex
	root  *inMemoryFile
}

func NewMemoryFileSystem() *MemoryFileSystem {
	// TODO: make root configurable
	root := newInMemoryFile("/", true)
	root.children[".."] = root

	return &MemoryFileSystem{
		root: root,
	}
}
