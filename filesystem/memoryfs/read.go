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

// Read reads up to len(buff) bytes into buff.
// This implementation is thread safe.
//
// Returns an error when:
// - the file is not open
func (fs *MemoryFileSystem) Read(fileDescriptor string, buff []byte) (int, error) {
	// Read lock the open file table
	fs.openFiles.RLock()

	// Get the file from the open files table
	fd, found := fs.openFiles.table[fileDescriptor]
	if !found {
		fs.openFiles.RUnlock()
		return 0, fserrors.ErrNotOpen
	}

	// Read lock file
	fd.data.RLock()
	defer fd.data.RUnlock()

	// Unlock file table
	fs.openFiles.RUnlock()
	return fd.Read(buff)
}

// ReadAt reads up to len(buff) bytes starting at offset into buff.
// This implementation is thread safe.
//
// Returns an error when:
// - the file is not open
func (fs *MemoryFileSystem) ReadAt(fileDescriptor string, buff []byte, offset int) (int, error) {
	// Read lock the open file table
	fs.openFiles.RLock()

	// Get the file from the open files table
	fd, found := fs.openFiles.table[fileDescriptor]
	if !found {
		fs.openFiles.RUnlock()
		return 0, fserrors.ErrNotOpen
	}

	// Read lock file
	fd.data.RLock()
	defer fd.data.RUnlock()

	// Unlock file table
	fs.openFiles.RUnlock()

	return fd.ReadAt(buff, offset)
}
