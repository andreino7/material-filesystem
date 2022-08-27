package filesystem

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
)

type FileSystemType int

const (
	InMemoryFileSystem FileSystemType = iota
)

// TODO: validate workdir - deleted working dir
type FileSystem interface {
	Mkdir(path *fspath.FileSystemPath) (file.File, error)
	MkdirAll(path *fspath.FileSystemPath) (file.File, error)
	CreateRegularFile(path *fspath.FileSystemPath) (file.File, error)
	DefaultWorkingDirectory() file.File
	GetDirectory(path *fspath.FileSystemPath) (file.File, error)
	Remove(path *fspath.FileSystemPath) (file.FileInfo, error)
	RemoveAll(path *fspath.FileSystemPath) (file.FileInfo, error)
	// TODO: handle regex in name
	FindFiles(name string, path *fspath.FileSystemPath) ([]file.FileInfo, error)
	ListFiles(path *fspath.FileSystemPath) ([]file.FileInfo, error)
	Move(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error)
	Copy(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error)
	CreateHardLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error)
	AppendAll(path *fspath.FileSystemPath, content []byte) error
	ReadAll(path *fspath.FileSystemPath) ([]byte, error)
	Open(path *fspath.FileSystemPath) (string, error)
	Close(fileDescriptor string)
	ReadAt(fileDescriptor string, startPos int, endPos int) ([]byte, error)
	WriteAt(fileDescriptor string, content []byte, pos int) (int, error)
}

func NewFileSystem(fsType FileSystemType) (FileSystem, error) {
	switch fsType {
	case InMemoryFileSystem:
		return memoryfs.NewMemoryFileSystem(), nil
	default:
		return nil, fmt.Errorf("unsupported filesystem type")
	}
}
