package memoryfs

import (
	"sync"
)

type MemoryFileSystem struct {
	mutex sync.RWMutex
	root  *inMemoryFile
}

func NewMemoryFileSystem() *MemoryFileSystem {
	// TODO: make root configurable
	root := newInMemoryFile("/", true)
	root.fileMap[".."] = root
	root.fileMap["."] = root
	root.fileMap["/"] = root

	return &MemoryFileSystem{
		root: root,
	}
}
