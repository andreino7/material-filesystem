package memoryfs

type MemoryFileSystem struct {
	root *fileWrapper
}

func NewMemoryFileSystem() *MemoryFileSystem {
	// TODO: make root configurable
	return &MemoryFileSystem{
		root: newFileWrapper("/", true),
	}
}
