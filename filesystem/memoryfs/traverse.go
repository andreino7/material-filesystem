package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

const MAX_LINK_DEPTH = 40

// traverseDirs traverses the path.Dir()
func (fs *MemoryFileSystem) traverseDirs(path *fspath.FileSystemPath) (*inMemoryFile, error) {
	return fs.traverseToDirWithCreateParentDirs(path, false)
}

// traverseDirsAndCreateParentDirs traverses the path.Dir() and creates any missing parent directory
func (fs *MemoryFileSystem) traverseDirsAndCreateParentDirs(path *fspath.FileSystemPath) (*inMemoryFile, error) {
	return fs.traverseToDirWithCreateParentDirs(path, true)
}

// traverseToBase traverses the path.Dir() and path.Base().
// If path.Base() is a symlink the link is resolved.
func (fs *MemoryFileSystem) traverseToBase(path *fspath.FileSystemPath) (*inMemoryFile, error) {
	return fs.traverseToBaseWithCreateParentDirsAndSkipLastLink(path, false, false)
}

// traverseToBaseWithSkipLastLink traverses the path.Dir() and path.Base().
// If path.Base() is a symlink and skipLastLink is false the link is not resolved.
func (fs *MemoryFileSystem) traverseToBaseWithSkipLastLink(path *fspath.FileSystemPath, skipLastLink bool) (*inMemoryFile, error) {
	return fs.traverseToBaseWithCreateParentDirsAndSkipLastLink(path, false, skipLastLink)
}

// traverseToBaseWithCreateParentDirsAndSkipLastLink traverses the path.Dir() and path.Base().
// If path.Base() is a symlink and skipLastLink is false the link is not resolved.
// If createParentDirs is true any missing parent directory is created.
func (fs *MemoryFileSystem) traverseToBaseWithCreateParentDirsAndSkipLastLink(path *fspath.FileSystemPath, createParentDirs bool, skipLink bool) (*inMemoryFile, error) {
	_, file, err := fs.traverse(path, createParentDirs, skipLink, 0)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, fserrors.ErrNotExist
	}

	return file, nil
}

// traverseToDirWithCreateParentDirs traverses the path.Dir()
// If createParentDirs is true any missing parent directory is created.
func (fs *MemoryFileSystem) traverseToDirWithCreateParentDirs(path *fspath.FileSystemPath, createParentDirs bool) (*inMemoryFile, error) {
	dir, _, err := fs.traverse(path, createParentDirs, true, 0)
	if err != nil {
		return nil, err
	}

	if dir.info.fileType != file.Directory {
		return nil, fserrors.ErrInvalidFileType
	}

	return dir, nil
}

// traverse moves through every path.Dir() and path.Base()
// If createParentDirs is true any missing parent directory is created.
// Any symbolic link is resolved with the exception of path.Base(), which is
// resolved only if skipLink is false.
func (fs *MemoryFileSystem) traverse(path *fspath.FileSystemPath, createDirs bool, skipLink bool, linkDepth int) (*inMemoryFile, *inMemoryFile, error) {
	// Find path starting point
	pathRoot, err := fs.findPathRoot(path)
	if err != nil {
		return nil, nil, err
	}

	// Move through every dir
	pathDirs := pathDirs(path)
	dir, err := fs.traverseFromRootToLastDir(pathRoot, pathDirs, createDirs, linkDepth)
	if err != nil {
		return nil, nil, err
	}

	// Move to the path.Base()
	targetFile, err := fs.moveToBase(dir, path.Base(), skipLink, linkDepth)
	if err != nil {
		return nil, nil, nil
	}

	return dir, targetFile, nil
}

// traverseFromRootToLastDir moves through every path.Dir()
// If createParentDirs is true any missing parent directory is created.
// Any symbolic link is resolved.
func (fs *MemoryFileSystem) traverseFromRootToLastDir(pathRoot *inMemoryFile, pathDirs []string, createDirs bool, linkDepth int) (*inMemoryFile, error) {
	curr := pathRoot
	for _, nextFileName := range pathDirs {
		// move to next file in path
		next, err := fs.moveToNext(curr, nextFileName, createDirs)
		if err != nil {
			return nil, err
		}

		// resolve symlink if needed
		next, linkErr := fs.resolveSymlink(next, linkDepth)
		if linkErr != nil {
			return nil, linkErr
		}

		if next.info.fileType != file.Directory {
			return nil, fserrors.ErrInvalidFileType
		}

		curr = next
	}

	return curr, nil
}

// moveToBase moves from the last dir in path.Dir() to path.Base()
// If base is a symlink is resolved only if skipLink is false.
func (fs *MemoryFileSystem) moveToBase(dir *inMemoryFile, fileName string, skipLink bool, linkDepth int) (*inMemoryFile, error) {
	targetFile, found := dir.fileMap[fileName]
	if !found {
		return nil, nil
	}

	if skipLink {
		return targetFile, nil
	}

	targetFile, err := fs.resolveSymlink(targetFile, linkDepth+1)
	if err != nil {
		return nil, nil
	}

	return targetFile, nil
}

// moveToNext moves to the nextFileName.
// if createDirs is true creates any missing parent directories.
func (fs *MemoryFileSystem) moveToNext(curr *inMemoryFile, nextFileName string, createDirs bool) (*inMemoryFile, error) {
	next, found := curr.fileMap[nextFileName]
	if found {
		return next, nil
	}

	if !createDirs {
		return nil, fserrors.ErrNotExist
	}

	return fs.create(nextFileName, file.Directory, curr)
}

// resolveSymlink tries to resolve a symlink and returns an error if the link points to a file
// that does not exixt or too many symlink were followed.
func (fs *MemoryFileSystem) resolveSymlink(currentFile *inMemoryFile, linkDepth int) (*inMemoryFile, error) {
	if currentFile.info.fileType != file.SymbolicLink {
		return currentFile, nil
	}

	if linkDepth >= MAX_LINK_DEPTH {
		return nil, fserrors.ErrTooManyLinks
	}

	// Link contains absolute path, so no need to pass working dir
	_, target, err := fs.traverse(currentFile.link, false, false, linkDepth+1)
	if err != nil {
		return nil, err
	}

	if target == nil {
		return nil, fserrors.ErrNotExist
	}

	return target, nil
}

// findPathRoot finds the path starting point
func (fs *MemoryFileSystem) findPathRoot(path *fspath.FileSystemPath) (*inMemoryFile, error) {
	if path.IsAbs() {
		return fs.root, nil
	}

	return fs.resolveWorkDir(path)
}
