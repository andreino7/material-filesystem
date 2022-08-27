package filesystem

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"material/filesystem/filesystem/user"
)

type FileSystemType int

const (
	InMemoryFileSystem FileSystemType = iota
)

// TODO: validate workdir - deleted working dir
type FileSystem interface {
	Mkdir(path *fspath.FileSystemPath, user user.User) (file.File, error)
	MkdirAll(path *fspath.FileSystemPath, user user.User) (file.File, error)
	CreateRegularFile(path *fspath.FileSystemPath, user user.User) (file.File, error)
	DefaultWorkingDirectory() file.File
	GetDirectory(path *fspath.FileSystemPath, user user.User) (file.File, error)
	Remove(path *fspath.FileSystemPath, user user.User) (file.FileInfo, error)
	RemoveAll(path *fspath.FileSystemPath, user user.User) (file.FileInfo, error)
	// TODO: handle regex in name
	FindFiles(name string, path *fspath.FileSystemPath, user user.User) ([]file.FileInfo, error)
	ListFiles(path *fspath.FileSystemPath, user user.User) ([]file.FileInfo, error)
	Move(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, user user.User) (file.FileInfo, error)
	Copy(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, user user.User) (file.FileInfo, error)
	CreateHardLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, user user.User) (file.FileInfo, error)
	AppendAll(path *fspath.FileSystemPath, content []byte, user user.User) error
	ReadAll(path *fspath.FileSystemPath, user user.User) ([]byte, error)
	Open(path *fspath.FileSystemPath, user user.User) (string, error)
	Close(fileDescriptor string, user user.User)
	ReadAt(fileDescriptor string, startPos int, endPos int, user user.User) ([]byte, error)
	WriteAt(fileDescriptor string, content []byte, pos int, user user.User) (int, error)
}

func NewFileSystem(fsType FileSystemType) (FileSystem, error) {
	switch fsType {
	case InMemoryFileSystem:
		return memoryfs.NewMemoryFileSystem(), nil
	default:
		return nil, fmt.Errorf("unsupported filesystem type")
	}
}
