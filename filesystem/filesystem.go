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
	Mkdir(path *fspath.FileSystemPath, workingDir file.File) (file.File, error)
	MkdirAll(path *fspath.FileSystemPath, workingDir file.File) (file.File, error)
	CreateRegularFile(path *fspath.FileSystemPath, workingDir file.File) (file.File, error)
	DefaultWorkingDirectory() file.File
	GetDirectory(path *fspath.FileSystemPath, workingDir file.File) (file.File, error)
	Remove(path *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error)
	RemoveAll(path *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error)
	// TODO: handle regex in name
	FindFiles(name string, path *fspath.FileSystemPath, workingDir file.File) ([]file.FileInfo, error)
	ListFiles(path *fspath.FileSystemPath, workingDir file.File) ([]file.FileInfo, error)
	Move(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error)
}

func NewFileSystem(fsType FileSystemType) (FileSystem, error) {
	switch fsType {
	case InMemoryFileSystem:
		return memoryfs.NewMemoryFileSystem(), nil
	default:
		return nil, fmt.Errorf("unsupported filesystem type")
	}
}
