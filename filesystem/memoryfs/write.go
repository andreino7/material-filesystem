package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// TODO: create missing directories is an option
// TODO: handle keeping the file open for one single writer and multipe readers
func (fs *MemoryFileSystem) AppendAll(path *fspath.FileSystemPath, content []byte, workingDir file.File) error {
	fs.mutex.Lock()

	fileToWrite, err := fs.getFileToWrite(path, workingDir, true, 0)
	if err != nil {
		fs.mutex.Unlock()
		return err
	}
	fileToWrite.data.mutex.Lock()
	defer fileToWrite.data.mutex.Unlock()
	fs.mutex.Unlock()

	fileToWrite.data.append(content)
	return nil
}

func (fs *MemoryFileSystem) WriteAt(fileDescriptor string, content []byte, pos int) (int, error) {
	// Read lock the open file table
	fs.openFiles.mutex.RLock()

	data, found := fs.openFiles.table[fileDescriptor]
	if !found {
		fs.openFiles.mutex.RUnlock()
		return 0, fserrors.ErrNotOpen
	}

	// Write lock file
	data.data.mutex.Lock()
	defer data.data.mutex.Unlock()

	// Unlock file table
	fs.openFiles.mutex.RUnlock()

	return data.data.writeAt(content, pos), nil
}

func (fs *MemoryFileSystem) getFileToWrite(path *fspath.FileSystemPath, workingDir file.File, createIntermediateDir bool, linkDepth int) (*inMemoryFile, error) {
	parentDir, err := fs.navigateToLastDirInPath(path, workingDir, createIntermediateDir, linkDepth)
	if err != nil {
		return nil, err
	}

	fileToWrite, found := parentDir.fileMap[path.Base()]
	if !found {
		if parentDir.info.fileType != file.Directory {
			return nil, fserrors.ErrInvalidFileType
		} else {
			return fs.createFile(path.Base(), parentDir)
		}
	}

	if fileToWrite.info.FileType() == file.RegularFile {
		return fileToWrite, nil
	}

	if fileToWrite.info.FileType() == file.Directory {
		return nil, fserrors.ErrInvalidFileType
	}

	// symbolic link
	return fs.getFileToWrite(fspath.NewFileSystemPath(fileToWrite.link), nil, false, linkDepth+1)
}
