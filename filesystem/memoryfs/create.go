package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"path/filepath"
)

func (fs *MemoryFileSystem) Mkdir(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	// RW lock the fs
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	return fs.createAt(path, workingDir, file.Directory, false)
}

func (fs *MemoryFileSystem) MkdirAll(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	// RW lock the fs
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	return fs.createAt(path, workingDir, file.Directory, true)
}

func (fs *MemoryFileSystem) CreateRegularFile(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	// RW lock the fs
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if err := checkFilePath(path); err != nil {
		return nil, err
	}
	return fs.createAt(path, workingDir, file.RegularFile, false)
}

func (fs *MemoryFileSystem) createAt(path *fspath.FileSystemPath, workingDir file.File, fileType file.FileType, isRecursive bool) (*inMemoryFile, error) {
	// TODO: validate file name (alphanumeric for simplicity)
	if err := checkFilePath(path); err != nil {
		return nil, err
	}

	// find where file needs to be added
	parent, err := fs.navigateToLastDirInPath(path, workingDir, isRecursive, 0)
	if err != nil {
		return nil, err
	}

	// create the file
	return fs.create(path.Base(), fileType, parent)
}

func (fs *MemoryFileSystem) createDirectory(fileName string, parent *inMemoryFile) (*inMemoryFile, error) {
	return fs.create(fileName, file.Directory, parent)
}

func (fs *MemoryFileSystem) createFile(fileName string, parent *inMemoryFile) (*inMemoryFile, error) {
	return fs.create(fileName, file.RegularFile, parent)
}

// TODO: this is the only place that creates files
func (fs *MemoryFileSystem) create(fileName string, fileType file.FileType, parent *inMemoryFile) (*inMemoryFile, error) {
	if _, found := parent.fileMap[fileName]; found {
		return nil, fserrors.ErrExist
	}

	absolutePath := filepath.Join(parent.info.AbsolutePath(), fileName)
	newFile := newInMemoryFile(absolutePath, fileType)
	fs.attachToParent(newFile, parent)
	return newFile, nil
}

func (fs *MemoryFileSystem) attachToParent(newFile *inMemoryFile, parent *inMemoryFile) {
	parent.fileMap[newFile.info.Name()] = newFile
	newFile.fileMap[".."] = parent
}
