package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/user"
)

const MAX_LINK_DEPTH = 40

func (fs *MemoryFileSystem) traverseToDir(path *fspath.FileSystemPath, user user.User) (*inMemoryFile, error) {
	return fs.traverseToDirWithCreateIntermediateDirs(path, false, user)
}

func (fs *MemoryFileSystem) traverseToDirAndCreateParentDirs(path *fspath.FileSystemPath, user user.User) (*inMemoryFile, error) {
	return fs.traverseToDirWithCreateIntermediateDirs(path, true, user)
}

func (fs *MemoryFileSystem) traverseToBase(path *fspath.FileSystemPath, user user.User) (*inMemoryFile, error) {
	return fs.traverseToBaseWithCreateParentDirsAndSkipLastLink(path, false, false, user)
}

func (fs *MemoryFileSystem) traverseToBaseWithSkipLastLink(path *fspath.FileSystemPath, skipLastLink bool, user user.User) (*inMemoryFile, error) {
	return fs.traverseToBaseWithCreateParentDirsAndSkipLastLink(path, false, skipLastLink, user)
}

func (fs *MemoryFileSystem) traverseToBaseWithCreateParentDirsAndSkipLastLink(path *fspath.FileSystemPath, createParentDirs bool, skipLink bool, user user.User) (*inMemoryFile, error) {
	_, file, err := fs.traverse(path, createParentDirs, skipLink, 0, user)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, fserrors.ErrNotExist
	}

	return file, nil
}

func (fs *MemoryFileSystem) traverseToDirWithCreateIntermediateDirs(path *fspath.FileSystemPath, createParentDirs bool, user user.User) (*inMemoryFile, error) {
	dir, _, err := fs.traverse(path, createParentDirs, true, 0, user)
	if err != nil {
		return nil, err
	}

	if dir.info.fileType != file.Directory {
		return nil, fserrors.ErrInvalidFileType
	}

	return dir, nil
}

// traverse moves through every directory/file in pathNames until it reaches the end of the array or
// an error occurs
func (fs *MemoryFileSystem) traverse(path *fspath.FileSystemPath, createDirs bool, skipLink bool, linkDepth int, user user.User) (*inMemoryFile, *inMemoryFile, error) {
	// Find path starting point
	pathRoot, err := fs.findPathRoot(path)
	if err != nil {
		return nil, nil, err
	}

	pathDirs := pathDirs(path)
	dir, err := fs.traverseFromRootToLastDir(pathRoot, pathDirs, createDirs, linkDepth, user)
	if err != nil {
		return nil, nil, err
	}

	targetFile, err := fs.moveToBase(dir, path.Base(), skipLink, linkDepth, user)
	if err != nil {
		return nil, nil, nil
	}

	return dir, targetFile, nil
}

func (fs *MemoryFileSystem) traverseFromRootToLastDir(pathRoot *inMemoryFile, pathDirs []string, createDirs bool, linkDepth int, user user.User) (*inMemoryFile, error) {
	curr := pathRoot

	for _, nextFileName := range pathDirs {
		if err := checkReadPermissions(curr, user); err != nil {
			return nil, err
		}

		next, err := fs.moveToNext(curr, nextFileName, createDirs, user)
		if err != nil {
			return nil, err
		}

		next, linkErr := fs.resolveSymlink(next, linkDepth+1, user)
		if linkErr != nil {
			return nil, linkErr
		}

		if next.info.fileType != file.Directory {
			return nil, fserrors.ErrInvalidFileType
		}

		curr = next
	}

	if err := checkReadPermissions(curr, user); err != nil {
		return nil, err
	}

	return curr, nil
}

func (fs *MemoryFileSystem) moveToBase(dir *inMemoryFile, fileName string, skipLink bool, linkDepth int, user user.User) (*inMemoryFile, error) {
	targetFile, found := dir.fileMap[fileName]
	if !found {
		return nil, nil
	}

	if err := checkReadPermissions(targetFile, user); err != nil {
		return nil, err
	}

	if skipLink {
		return targetFile, nil
	}

	targetFile, err := fs.resolveSymlink(targetFile, linkDepth+1, user)
	if err != nil {
		return nil, nil
	}

	return targetFile, nil
}

func (fs *MemoryFileSystem) moveToNext(curr *inMemoryFile, nextFileName string, createDirs bool, user user.User) (*inMemoryFile, error) {
	next, found := curr.fileMap[nextFileName]
	if found {
		return next, nil
	}

	if !createDirs {
		return nil, fserrors.ErrNotExist
	}

	return fs.create(nextFileName, file.Directory, curr, user)
}

// resolveSymlink tries to resolve symlink and returns an error if the link points to a file
// that does not exixt or too many symlink were followed.
func (fs *MemoryFileSystem) resolveSymlink(currentFile *inMemoryFile, linkDepth int, user user.User) (*inMemoryFile, error) {
	if currentFile.info.fileType != file.SymbolicLink {
		return currentFile, nil
	}

	if err := checkReadPermissions(currentFile, user); err != nil {
		return nil, err
	}

	if linkDepth >= MAX_LINK_DEPTH {
		return nil, fserrors.ErrTooManyLinks
	}

	// Link contains absolute path, so no need to pass working dir
	_, target, err := fs.traverse(currentFile.link, false, false, linkDepth+1, user)
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
