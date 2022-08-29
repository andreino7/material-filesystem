package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// AppendAll writes data to the named file, creating it if necessary along
// with any missing parent directories.
//
// Returns an error when:
// - the file is not a regular file
//
// TODO: create parent directories is an option
func (fs *MemoryFileSystem) AppendAll(path *fspath.FileSystemPath, content []byte) error {
	fs.Lock()

	parent, err := fs.traverseDirsAndCreateParentDirs(path)
	if err != nil {
		fs.Unlock()
		return err
	}

	fileToWrite, err := fs.createFileToWriteIfMissing(parent, path.Base())
	if err != nil {
		fs.Unlock()
		return err
	}

	descriptor, err := fs.doOpen(fileToWrite)
	fs.Unlock()
	if err != nil {
		return err
	}
	defer fs.Close(descriptor)

	_, err = fs.doWrite(descriptor, content, func(fd *fileDescriptor) (int, error) {
		return fd.Write(content)
	})

	if err != nil {
		return err
	}

	return nil
}

// Write writes content to the file and
// returns the number of bytes written.
//
// Returns an error when:
// - the file is not open
func (fs *MemoryFileSystem) Write(descriptor string, content []byte) (int, error) {
	return fs.doWrite(descriptor, content, func(fd *fileDescriptor) (int, error) {
		return fd.Write(content)
	})
}

// WriteAt writes content to the file starting at offset
// and returns the number of bytes written.
//
// Returns an error when:
// - the file is not open
func (fs *MemoryFileSystem) WriteAt(descriptor string, content []byte, offset int) (int, error) {
	if offset < 0 {
		return 0, fserrors.ErrInvalid
	}

	return fs.doWrite(descriptor, content, func(fd *fileDescriptor) (int, error) {
		return fd.WriteAt(content, offset)
	})
}

func (fs *MemoryFileSystem) doWrite(fileDescriptor string, content []byte, writeFn func(fd *fileDescriptor) (int, error)) (int, error) {
	// Read lock the open file table
	fs.openFiles.RLock()

	fd, found := fs.openFiles.table[fileDescriptor]
	if !found {
		fs.openFiles.RUnlock()
		return 0, fserrors.ErrNotOpen
	}

	// Write lock file
	fd.data.Lock()
	defer fd.data.Unlock()

	// Unlock file table
	fs.openFiles.RUnlock()
	return writeFn(fd)
}

func (fs *MemoryFileSystem) createFileToWriteIfMissing(parent *inMemoryFile, name string) (*inMemoryFile, error) {
	fileToWrite, err := fs.moveToBase(parent, name, false, 0)
	if err != nil {
		return nil, err
	}

	if fileToWrite == nil {
		return fs.create(name, file.RegularFile, parent)
	}
	return fileToWrite, nil
}
