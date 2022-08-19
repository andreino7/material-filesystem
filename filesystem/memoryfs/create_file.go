package memoryfs

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
)

// TODO: validate file name
func (fs *MemoryFileSystem) Mkdir(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	return fs.addFileToFs(path, workingDir, true, false)
}

func (fs *MemoryFileSystem) MkdirAll(path *fspath.FileSystemPath, workingDir file.File) (file.File, error) {
	return fs.addFileToFs(path, workingDir, true, true)
}

func (fs *MemoryFileSystem) addFileToFs(path *fspath.FileSystemPath, workingDir file.File, isDirectory bool, isRecursive bool) (file.File, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	parent, err := fs.lookupDirWithCreateMissing(path, workingDir, isRecursive)
	if err != nil {
		return nil, err
	}

	return fs.createFile(path.Base(), isDirectory, parent)
}

func (fs *MemoryFileSystem) createFile(fileName string, isDirectory bool, parent *inMemoryFile) (*inMemoryFile, error) {
	if _, found := parent.children[fileName]; found {
		return nil, fmt.Errorf("file already exists")
	}
	newFile := newInMemoryFile(fileName, isDirectory)
	fs.linkToParent(newFile, parent)
	return newFile, nil
}

func (fs *MemoryFileSystem) linkToParent(newFile *inMemoryFile, parent *inMemoryFile) {
	parent.children[newFile.info.Name()] = newFile
	newFile.children[".."] = parent
}

func (fs *MemoryFileSystem) lookupDir(path *fspath.FileSystemPath, workingDir file.File) (*inMemoryFile, error) {
	return fs.lookupDirWithCreateMissing(path, workingDir, false)
}

// TODO: refactor
func (fs *MemoryFileSystem) lookupDirWithCreateMissing(path *fspath.FileSystemPath, workingDir file.File, createMissing bool) (*inMemoryFile, error) {
	pathRoot, pathDirs, err := fs.findPathRoot(path, workingDir)
	if err != nil {
		return nil, err
	}

	for _, currentDir := range pathDirs {
		tmp, found := pathRoot.children[currentDir]
		if !found {
			if createMissing {
				tmp, err = fs.createFile(currentDir, true, pathRoot)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("no such file or directory")
			}
		}

		if !tmp.info.IsDirectory() {
			return nil, fmt.Errorf("file is not a directory")
		}

		pathRoot = tmp
	}

	return pathRoot, nil
}

func (fs *MemoryFileSystem) findPathRoot(path *fspath.FileSystemPath, workingDir file.File) (*inMemoryFile, []string, error) {
	if path.IsAbs() || workingDir == nil {
		pathDirs, _ := path.SplitAbs()
		return fs.root, pathDirs, nil
	}

	pathRoot, err := fs.resolveWorkDir(path, workingDir)
	if err != nil {
		return nil, nil, err
	}
	pathDirs, _ := path.Split()
	return pathRoot, pathDirs, nil

}

// TODO: Test directory not attached to fs
func (fs *MemoryFileSystem) resolveWorkDir(path *fspath.FileSystemPath, workingDir file.File) (*inMemoryFile, error) {
	currentDir, ok := workingDir.(inMemoryFile)
	if !ok {
		return nil, fmt.Errorf("invalid working directory")
	}

	if currentDir.isDeleted {
		return nil, fmt.Errorf("working directory deleted")
	}

	return &currentDir, nil
}
