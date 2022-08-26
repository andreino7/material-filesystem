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

// TODO: validate file name
// TODO: make create intermediate directories configurable
func (fs *MemoryFileSystem) CreateHardLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Locate file to link
	fileToLink, err := fs.traverseToBase(srcPath, workingDir)
	if err != nil {
		return nil, err
	}

	// Only hard links to regular file supported
	if fileToLink.info.fileType != file.RegularFile {
		return nil, fserrors.ErrInvalidFileType
	}

	// Create an empty file
	hardLink, err := fs.createAt(destPath, workingDir, file.RegularFile, true)
	if err != nil {
		return nil, err
	}

	// Point the file to the same underline data
	hardLink.data = fileToLink.data
	return hardLink.info, nil
}

// TODO: make create intermediate directories configurable
// TODO: document symbolic links to not existing file should work
func (fs *MemoryFileSystem) CreateSymbolicLink(srcPath *fspath.FileSystemPath, destPath *fspath.FileSystemPath, workingDir file.File) (file.FileInfo, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Create an empty file
	symLink, err := fs.createAt(destPath, workingDir, file.SymbolicLink, false)
	if err != nil {
		return nil, err
	}

	// Point the file to the original file
	symLink.link = toAbsolutePath(srcPath, workingDir)
	return symLink.info, nil
}

func (fs *MemoryFileSystem) createAt(path *fspath.FileSystemPath, workingDir file.File, fileType file.FileType, isRecursive bool) (*inMemoryFile, error) {
	// TODO: validate file name (alphanumeric for simplicity)
	if err := checkFilePath(path); err != nil {
		return nil, err
	}

	// find where file needs to be added
	parent, err := fs.traverseToDirWithCreateIntermediateDirs(path, workingDir, isRecursive)
	if err != nil {
		return nil, err
	}

	// create the file
	return fs.create(path.Base(), fileType, parent)
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
