package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"

	"github.com/google/uuid"
)

// Open opens the named file for reading or writing.
// Returns an error if the file was not found or it's not a "regular" file.
// This implementation is thread safe
//
// Returns an error when:
// - path does not exist
// - the file is not a RegularFile
func (fs *MemoryFileSystem) Open(path *fspath.FileSystemPath) (string, error) {
	fs.Lock()
	defer fs.Unlock()

	fileToOpen, err := fs.traverseToBase(path)
	if err != nil {
		return "", err
	}

	if fileToOpen.info.fileType != file.RegularFile {
		return "", fserrors.ErrInvalidFileType
	}

	fs.openFiles.Lock()
	defer fs.openFiles.Unlock()

	// TODO: fd could just be an int
	fd := uuid.NewString()
	fs.openFiles.table[fd] = &fileDescriptor{data: fileToOpen.data, offset: 0}
	return fd, nil
}

// Close closes the file associated to the given descriptor
// If the file is not open, this is a noop.
// This implementation is thread safe.
func (fs *MemoryFileSystem) Close(descriptor string) {
	fs.openFiles.Lock()
	defer fs.openFiles.Unlock()
	delete(fs.openFiles.table, descriptor)
}
