package memoryfs

import (
	"material/filesystem/filesystem/file"
	"sync"
)

type fileTable struct {
	table map[string]*fileDescriptor
	sync.RWMutex
}

func newFileTable() *fileTable {
	return &fileTable{
		table: map[string]*fileDescriptor{},
	}
}

// MemoryFileSystem implements the FileSystem interface
// and it's an in-memory file system.
type MemoryFileSystem struct {
	sync.RWMutex
	// root of the file system
	root *inMemoryFile
	// table of open files
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
