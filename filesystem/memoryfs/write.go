package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// TODO: create missing directories is an option
// TODO: handle keeping the file open for one single writer and multipe readers
func (fs *MemoryFileSystem) AppendAll(path *fspath.FileSystemPath, content []byte) error {
	fs.Lock()

	parent, err := fs.traverseToDirAndCreateParentDirs(path)
	if err != nil {
		fs.Unlock()
		return err
	}

	fileToWrite, err := fs.createFileToWriteIfMissing(parent, path.Base())
	if err != nil {
		fs.Unlock()
		return err
	}

	if fileToWrite.info.fileType != file.RegularFile {
		fs.Unlock()
		return fserrors.ErrInvalidFileType
	}

	fileToWrite.data.Lock()
	defer fileToWrite.data.Unlock()
	fs.Unlock()

	fileToWrite.data.append(content)
	return nil
}

func (fs *MemoryFileSystem) WriteAt(fileDescriptor string, content []byte, pos int) (int, error) {
	// Read lock the open file table
	fs.openFiles.RLock()

	data, found := fs.openFiles.table[fileDescriptor]
	if !found {
		fs.openFiles.RUnlock()
		return 0, fserrors.ErrNotOpen
	}

	// Write lock file
	data.data.Lock()
	defer data.data.Unlock()

	// Unlock file table
	fs.openFiles.RUnlock()

	return data.data.writeAt(content, pos), nil
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
