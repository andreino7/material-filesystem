package memoryfs

import (
	"material/filesystem/filesystem/file"
	"sync"
)

type MemoryFileSystem struct {
	sync.RWMutex
	root      *inMemoryFile
	openFiles *fileTable
}

func NewMemoryFileSystem() *MemoryFileSystem {
	// TODO: make root configurable
	root := newInMemoryFile("/", file.Directory)
	root.fileMap[".."] = root
	root.fileMap["."] = root
	root.fileMap["/"] = root

	return &MemoryFileSystem{
		root:      root,
		openFiles: newFileTable(),
	}
}
