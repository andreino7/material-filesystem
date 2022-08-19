package memoryfs

type MemoryFileSystem struct{}

func NewMemoryFileSystem() *MemoryFileSystem {
	return &MemoryFileSystem{}
}
