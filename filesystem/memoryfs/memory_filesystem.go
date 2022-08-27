package memoryfs

import (
	"material/filesystem/filesystem/file"
	"sync"
)

const (
	defaultUser  = "root"
	defaultGroup = "root"
	defaultDir   = "/"
)

type MemoryFileSystem struct {
	sync.RWMutex
	root      *inMemoryFile
	openFiles *fileTable
	users     *userTable
	groups    *groupTable
}

func NewMemoryFileSystem() *MemoryFileSystem {
	// TODO: make root configurable
	root := newInMemoryFile("/", file.Directory, defaultUser, defaultGroup)
	root.fileMap[".."] = root
	root.fileMap["."] = root
	root.fileMap["/"] = root

	return &MemoryFileSystem{
		root:      root,
		openFiles: newFileTable(),
		users:     newUserTable(defaultUser, defaultGroup),
		groups:    newGroupTable(defaultGroup),
	}
}
