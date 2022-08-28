package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// ReadAll reads the named file and returns the contents.
// This implementation is thread safe.
//
// Returns an error when:
// - the file does not exist
// - the file is not a regular file
func (fs *MemoryFileSystem) ReadAll(path *fspath.FileSystemPath) ([]byte, error) {
	fs.Lock()

	fileToRead, err := fs.traverseToBase(path)
	if err != nil {
		fs.Unlock()
		return nil, err
	}

	if fileToRead.info.fileType != file.RegularFile {
		fs.Unlock()
		return nil, fserrors.ErrInvalidFileType
	}

	// Read lock file
	fileToRead.data.RLock()
	defer fileToRead.data.RUnlock()
	fs.Unlock()

	return fileToRead.data.data, nil
}

// ReadAt reads endPos - startPos bytes from the file starting at startPos
// and returns them.
// This implementation is thread safe.
//
// Returns an error when:
// - the file is not open
// - startPos or endPos are invalid
func (fs *MemoryFileSystem) ReadAt(fileDescriptor string, startPos int, endPos int) ([]byte, error) {
	if err := validatePos(startPos, endPos); err != nil {
		return nil, err
	}

	// Read lock the open file table
	fs.openFiles.RLock()

	// Get the file from the open files table
	data, found := fs.openFiles.table[fileDescriptor]
	if !found {
		fs.openFiles.RUnlock()
		return nil, fserrors.ErrNotOpen
	}

	// Read lock file
	data.data.RLock()
	defer data.data.RUnlock()

	// Unlock file table
	fs.openFiles.RUnlock()

	return data.data.readAt(startPos, endPos), nil
}

func validatePos(start int, end int) error {
	if end < start {
		return fserrors.ErrInvalid
	}

	if start < 0 {
		return fserrors.ErrInvalid
	}

	return nil
}
