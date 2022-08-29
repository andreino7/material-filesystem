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

	descriptor, err := fs.doOpen(fileToRead)
	fs.Unlock()
	if err != nil {
		return nil, err
	}
	defer fs.Close(descriptor)

	var buff []byte
	fs.doRead(descriptor, func(fd *fileDescriptor) (int, error) {
		buff = make([]byte, fd.data.Size())
		return fd.Read(buff)
	})

	return buff, nil
}

// Read reads up to len(buff) bytes into buff.
// This implementation is thread safe.
//
// Returns an error when:
// - the file is not open
func (fs *MemoryFileSystem) Read(descriptor string, buff []byte) (int, error) {
	return fs.doRead(descriptor, func(fd *fileDescriptor) (int, error) {
		return fd.Read(buff)
	})
}

// ReadAt reads up to len(buff) bytes starting at offset into buff.
// This implementation is thread safe.
//
// Returns an error when:
// - the file is not open
func (fs *MemoryFileSystem) ReadAt(descriptor string, buff []byte, offset int) (int, error) {
	return fs.doRead(descriptor, func(fd *fileDescriptor) (int, error) {
		return fd.ReadAt(buff, offset)
	})
}

func (fs *MemoryFileSystem) doRead(fileDescriptor string, readFn func(fd *fileDescriptor) (int, error)) (int, error) {
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

	return readFn(fd)
}
