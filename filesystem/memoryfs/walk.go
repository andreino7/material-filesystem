package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
)

type visitFn func(string, *inMemoryFile) error

// Walk walks the file tree rooted at root, calling filterFn for each file or directory in the tree, including root,
// and calls walkFn for each file or directory matching the filter.
// Optionally follow symbolic links.
// If walkFn returns an error, the function stops immediately.
//
// Returns an error when:
// - too many links were followed
// - walkfn returns an error
// - the symbolic link doesn't exist
func (fs *MemoryFileSystem) Walk(path *fspath.FileSystemPath, walkFn file.WalkFn, filterFn file.FilterFn, followLinks bool) error {
	pathRoot, err := fs.traverseToBase(path)
	if err != nil {
		return err
	}

	return fs.doWalk(pathRoot, walkFn, filterFn, followLinks, 0)
}

func (fs *MemoryFileSystem) doWalk(rootFile *inMemoryFile, walkFn file.WalkFn, filterFn file.FilterFn, followLinks bool, linkDepth int) error {
	// check if current path is filtered out
	if !filterFn(rootFile) {
		return nil
	}

	// visit the current file
	if err := walkFn(rootFile); err != nil {
		return err
	}

	// Optionally follow links
	if followLinks && rootFile.info.fileType == file.SymbolicLink {
		currLink, err := fs.resolveSymlink(rootFile, linkDepth)
		if err != nil {
			return err
		}
		rootFile = currLink
		// Invoke walkfn on the link target
		if err := walkFn(rootFile); err != nil {
			return err
		}
	}

	return fs.visitDir(rootFile, func(_ string, curr *inMemoryFile) error {
		return fs.doWalk(curr, walkFn, filterFn, followLinks, linkDepth+1)
	})
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
