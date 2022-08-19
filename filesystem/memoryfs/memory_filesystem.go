package memoryfs

import "sync"

type MemoryFileSystem struct {
	mutex sync.RWMutex
	root  *fileWrapper
}

func NewMemoryFileSystem() *MemoryFileSystem {
	// TODO: make root configurable
	return &MemoryFileSystem{
		root: newFileWrapper("/", true),
	}
}
