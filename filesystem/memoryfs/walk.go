package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
)

type visitFn func(string, *inMemoryFile) error

func (fs *MemoryFileSystem) Walk(path *fspath.FileSystemPath, walkFn file.WalkFn, filterFn file.FilterFn, followLinks bool) error {
	pathRoot, err := fs.traverseToBase(path)
	if err != nil {
		return err
	}

	return fs.doWalk(pathRoot, walkFn, filterFn, followLinks)
}

func (fs *MemoryFileSystem) doWalk(rootFile *inMemoryFile, walkFn file.WalkFn, filterFn file.FilterFn, followLinks bool) error {
	for fileName, curr := range rootFile.fileMap {
		// skip special keys to avoid infinite cycle
		if fileName == ".." || fileName == "." || fileName == "/" {
			continue
		}

		if followLinks && curr.info.fileType == file.SymbolicLink {
			currLink, err := fs.resolveSymlink(curr, 0)
			if err != nil {
				return err
			}
			curr = currLink
		}

		if !filterFn(curr) {
			continue
		}

		if err := walkFn(curr); err != nil {
			return err
		}

		if curr.info.fileType != file.Directory {
			continue
		}

		err := fs.doWalk(curr, walkFn, filterFn, followLinks)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *MemoryFileSystem) visitDir(rootFile *inMemoryFile, visitFn visitFn) error {
	for fileName, file := range rootFile.fileMap {
		// skip special keys to avoid infinite cycle
		if fileName == ".." || fileName == "." || fileName == "/" {
			continue
		}

		if err := visitFn(fileName, file); err != nil {
			return err
		}
	}

	return nil
}
