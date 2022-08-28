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

type FileSystem interface {
	// Mkdir creates a new directory at the specified path
	// If there is an error, it will be of type *FileSystemError.
	Mkdir(path *fspath.FileSystemPath) (file.File, error)
	// MkdirAll creates a directory at the specified path,
	// along with any necessary parents.
	// If there is an error, it will be of type *FileSystemError.
	MkdirAll(path *fspath.FileSystemPath) (file.File, error)
	// CreateRegularFile creates a new file at the specified path
	// If there is an error, it will be of type *FileSystemError.
	CreateRegularFile(path *fspath.FileSystemPath) (file.File, error)
	// DefaultWorkingDirectory returns the default working directory.
	DefaultWorkingDirectory() file.File
	// GetDirectory returns the directory located at the specified path.
	// If there is an error, it will be of type *FileSystemError.
	GetDirectory(path *fspath.FileSystemPath) (file.File, error)
	// Remove removes the file located at the specified path.
	// If there is an error, it will be of type *FileSystemError.
	Remove(path *fspath.FileSystemPath) (file.FileInfo, error)
	// RemoveAll removes the file located at the specified path
	// and any children it contains.
	// If there is an error, it will be of type *FileSystemError.
	RemoveAll(path *fspath.FileSystemPath) (file.FileInfo, error)
	// TODO: handle regex in name
	FindFiles(name string, path *fspath.FileSystemPath) ([]file.FileInfo, error)
	// ListFiles lists the files at the specified path.
	// If there is an error, it will be of type *FileSystemError.
	ListFiles(path *fspath.FileSystemPath) ([]file.FileInfo, error)
	// Move moves (renames) srcPath to destPath.
	// If there is an error, it will be of type *FileSystemError.
	Move(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error)
	// Copy copies srcPath to destPath.
	// If there is an error, it will be of type *FileSystemError.
	Copy(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error)
	// Link creates srcPath as a hard link to the destPath file.
	// If there is an error, it will be of type *FileSystemError.
	CreateHardLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error)
	// Symlink creates srcPath as a symbolic link to destPath.
	// If there is an error, it will be of type *FileSystemError.
	CreateSymbolicLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error)
	// AppendAll writes data to the named file, creating it if necessary.
	// If there is an error, it will be of type *FileSystemError.
	AppendAll(path *fspath.FileSystemPath, content []byte) error
	// ReadAll reads the named file and returns the contents.
	// If there is an error, it will be of type *FileSystemError.
	ReadAll(path *fspath.FileSystemPath) ([]byte, error)
	// Open opens the named file for reading or writing.
	// If there is an error, it will be of type *FileSystemError.
	Open(path *fspath.FileSystemPath) (string, error)
	// Close closes the file associated to the given descriptor
	Close(fileDescriptor string)
	// ReadAt reads endPos - startPos bytes from the file starting at startPos
	// and returns them.
	// If there is an error, it will be of type *FileSystemError.
	ReadAt(fileDescriptor string, startPos int, endPos int) ([]byte, error)
	// WriteAt writes content to the file starting at pos
	// and returns the number of bytes written.
	// If there is an error, it will be of type *FileSystemError.
	WriteAt(fileDescriptor string, content []byte, pos int) (int, error)
	Walk(path *fspath.FileSystemPath, walkFn file.WalkFn, filterFn file.FilterFn, followLinks bool) error
}

// NewFileSystem creates a new filesystem for the given fsType.
// Returns an error if the fsType is not supported.
func NewFileSystem(fsType FileSystemType) (FileSystem, error) {
	switch fsType {
	case InMemoryFileSystem:
		return memoryfs.NewMemoryFileSystem(), nil
	default:
		return nil, fmt.Errorf("unsupported filesystem type")
	}
}
