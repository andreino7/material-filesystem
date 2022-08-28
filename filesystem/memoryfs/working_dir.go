package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// DefaultWorkingDirectory returns the root of the filesystem.
func (fs *MemoryFileSystem) DefaultWorkingDirectory() file.File {
	return fs.root
}

// GetDirectory returns the directory located at the specified path.
//
// Returns an error when:
// - the directory does not exist
// - the file is not a directory
func (fs *MemoryFileSystem) GetDirectory(path *fspath.FileSystemPath) (file.File, error) {
	// RLock the fs
	fs.RLock()
	defer fs.RUnlock()

	// Find path starting point
	dir, err := fs.traverseToBase(path)
	if err != nil {
		return nil, err
	}

	if dir.info.fileType != file.Directory {
		return nil, fserrors.ErrInvalidFileType
	}

	return dir, nil
}

// resolveWorkDir returns an error if the working directory has been deleted.
// This can happen is clients have a stale reference.
func (fs *MemoryFileSystem) resolveWorkDir(path *fspath.FileSystemPath) (*inMemoryFile, error) {
	currentDir, ok := path.WorkingDir().(*inMemoryFile)
	if !ok || currentDir.isDeleted {
		return nil, fserrors.ErrInvalidWorkingDirectory
	}

	return currentDir, nil
}
