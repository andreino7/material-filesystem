package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"path/filepath"
)

// TODO: validate file name
func (fs *MemoryFileSystem) Mkdir(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	if err := checkFilePath(path); err != nil {
		return nil, err
	}
	return fs.addFileToFs(path, workingDir, true, false)
}

func (fs *MemoryFileSystem) MkdirAll(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	if err := checkFilePath(path); err != nil {
		return nil, err
	}
	return fs.addFileToFs(path, workingDir, true, true)
}

func (fs *MemoryFileSystem) CreateRegularFile(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	if err := checkFilePath(path); err != nil {
		return nil, err
	}
	return fs.addFileToFs(path, workingDir, false, false)
}

func (fs *MemoryFileSystem) addFileToFs(path *fspath.FileSystemPath, workingDir file.File, isDirectory bool, isRecursive bool) (file.File, error) {
	// RW lock the fs
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// find where file needs to be added
	parent, err := fs.navigateToLastDirInPath(path, workingDir, isRecursive)
	if err != nil {
		return nil, err
	}

	// create the file
	return fs.createFile(path.Base(), isDirectory, parent)
}

func (fs *MemoryFileSystem) createDirectory(fileName string, parent *inMemoryFile) (*inMemoryFile, error) {
	return fs.createFile(fileName, true, parent)
}

func (fs *MemoryFileSystem) createFile(fileName string, isDirectory bool, parent *inMemoryFile) (*inMemoryFile, error) {
	if _, found := parent.fileMap[fileName]; found {
		return nil, fserrors.ErrExist
	}

	absolutePath := filepath.Join(parent.info.AbsolutePath(), fileName)
	newFile := newInMemoryFile(absolutePath, isDirectory)
	fs.attachToParent(newFile, parent)
	return newFile, nil
}

func (fs *MemoryFileSystem) attachToParent(newFile *inMemoryFile, parent *inMemoryFile) {
	parent.fileMap[newFile.info.Name()] = newFile
	newFile.fileMap[".."] = parent
}
