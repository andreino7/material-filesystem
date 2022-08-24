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
	return fs.addFileToFs(path, workingDir, file.Directory, false)
}

func (fs *MemoryFileSystem) MkdirAll(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	if err := checkFilePath(path); err != nil {
		return nil, err
	}
	return fs.addFileToFs(path, workingDir, file.Directory, true)
}

func (fs *MemoryFileSystem) CreateRegularFile(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	if err := checkFilePath(path); err != nil {
		return nil, err
	}
	return fs.addFileToFs(path, workingDir, file.RegularFile, false)
}

func (fs *MemoryFileSystem) addFileToFs(path *fspath.FileSystemPath, workingDir file.File, fileType file.FileType, isRecursive bool) (file.File, error) {
	// RW lock the fs
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// find where file needs to be added
	parent, err := fs.navigateToLastDirInPath(path, workingDir, isRecursive)
	if err != nil {
		return nil, err
	}

	// create the file
	return fs.createFile(path.Base(), fileType, parent)
}

func (fs *MemoryFileSystem) createDirectory(fileName string, parent *inMemoryFile) (*inMemoryFile, error) {
	return fs.createFile(fileName, file.Directory, parent)
}

func (fs *MemoryFileSystem) createFile(fileName string, fileType file.FileType, parent *inMemoryFile) (*inMemoryFile, error) {
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
