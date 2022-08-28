package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"path/filepath"
)

// Mkdir creates a new directory at the specified path.
// This implementation is thead safe.
//
// Returns an error when:
// - the file name is invalid
// - the file already exists
// - any of the directory in the path does not exist
func (fs *MemoryFileSystem) Mkdir(path *fspath.FileSystemPath) (file.File, error) {
	// RW lock the fs
	fs.Lock()
	defer fs.Unlock()
	return fs.createAt(path, file.Directory, false)
}

// MkdirAll creates a directory at the specified path,
// along with any necessary parents.
// This implementation is thead safe.
//
// Returns an error when:
// - the file name is invalid
// - the file already exists
func (fs *MemoryFileSystem) MkdirAll(path *fspath.FileSystemPath) (file.File, error) {
	// RW lock the fs
	fs.Lock()
	defer fs.Unlock()
	return fs.createAt(path, file.Directory, true)
}

// CreateRegularFile creates a new file at the specified path
// This implementation is thead safe.
//
// Returns an error when:
// - the file name is invalid
// - the file already exists
// - any of the directory in the path does not exist
func (fs *MemoryFileSystem) CreateRegularFile(path *fspath.FileSystemPath) (file.File, error) {
	// RW lock the fs
	fs.Lock()
	defer fs.Unlock()

	if err := checkFilePath(path); err != nil {
		return nil, err
	}
	return fs.createAt(path, file.RegularFile, false)
}

// Link creates srcPath along with any parent directories
// as a hard link to the destPath file.
// Only regular files are supported.
// This implementation is thead safe.
//
// Returns an error when:
// - srcPath file name is invalid
// - srcPath already exists
// - destPath does not exist
// - destPath is not a regular file
// TODO: make create parent directories configurable
func (fs *MemoryFileSystem) CreateHardLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error) {
	fs.Lock()
	defer fs.Unlock()

	// Locate file to link
	fileToLink, err := fs.traverseToBase(srcPath)
	if err != nil {
		return nil, err
	}

	// Only hard links to regular file supported
	if fileToLink.info.fileType != file.RegularFile {
		return nil, fserrors.ErrInvalidFileType
	}

	// Create an empty file
	hardLink, err := fs.createAt(destPath, file.RegularFile, true)
	if err != nil {
		return nil, err
	}

	// Point the file to the same underline data
	hardLink.data = fileToLink.data
	return hardLink.info, nil
}

// Symlink creates srcPath along with any parent directories
// as a symbolic link to destPath.
// Symlink can be created to a non-existent destPath.
// If destPath is later created the symlink will start working.
// This implementation is thead safe.
//
// Returns an error when:
// - srcPath file name is invalid
// - srcPath already exists
// TODO: make create intermediate directories configurable
func (fs *MemoryFileSystem) CreateSymbolicLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath) (file.FileInfo, error) {
	fs.Lock()
	defer fs.Unlock()

	// Create an empty file
	symLink, err := fs.createAt(destPath, file.SymbolicLink, false)
	if err != nil {
		return nil, err
	}

	// Point the file to the original file
	pathLink, err := fspath.NewFileSystemPath(srcPath.AbsolutePath(), nil)
	if err != nil {
		return nil, err
	}

	symLink.link = pathLink
	return symLink.info, nil
}

func (fs *MemoryFileSystem) createAt(path *fspath.FileSystemPath, fileType file.FileType, isRecursive bool) (*inMemoryFile, error) {
	if err := checkFilePath(path); err != nil {
		return nil, err
	}

	// find where file needs to be added
	parent, err := fs.traverseToDirWithCreateParentDirs(path, isRecursive)
	if err != nil {
		return nil, err
	}

	// create the file
	return fs.create(path.Base(), fileType, parent)
}

func (fs *MemoryFileSystem) create(fileName string, fileType file.FileType, parent *inMemoryFile) (*inMemoryFile, error) {
	// check if file exists
	if _, found := parent.fileMap[fileName]; found {
		return nil, fserrors.ErrExist
	}

	// create new file and add to fs tree
	absolutePath := filepath.Join(parent.info.AbsolutePath(), fileName)
	newFile := newInMemoryFile(absolutePath, fileType)
	fs.attachToParent(newFile, parent)
	return newFile, nil
}

func (fs *MemoryFileSystem) attachToParent(newFile *inMemoryFile, parent *inMemoryFile) {
	parent.fileMap[newFile.info.Name()] = newFile
	newFile.fileMap[".."] = parent
}
