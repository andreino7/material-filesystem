package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/user"
)

func (fs *MemoryFileSystem) DefaultWorkingDirectory() file.File {
	return fs.root
}

func (fs *MemoryFileSystem) GetDirectory(path *fspath.FileSystemPath, user user.User) (file.File, error) {
	// RLock the fs
	fs.RLock()
	defer fs.RUnlock()

	// Find path starting point
	dir, err := fs.traverseToBase(path, user)
	if err != nil {
		return nil, err
	}

	if dir.info.fileType != file.Directory {
		return nil, fserrors.ErrInvalidFileType
	}

	return dir, nil
}

func (fs *MemoryFileSystem) resolveWorkDir(path *fspath.FileSystemPath) (*inMemoryFile, error) {
	currentDir, ok := path.WorkingDir().(*inMemoryFile)
	if !ok || currentDir.isDeleted {
		return nil, fserrors.ErrInvalidWorkingDirectory
	}

	return currentDir, nil
}
