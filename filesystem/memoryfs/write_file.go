package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
)

// TODO: create missing directories is an option
// TODO: handle keeping the file open for one single writer and multipe readers
func (fs *MemoryFileSystem) AppendToFile(path *fspath.FileSystemPath, content []byte, workingDir file.File) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	fileToWrite, err := fs.getFileToWrite(path, workingDir, true, 0)
	if err != nil {
		return err
	}

	fileToWrite.data.data = append(fileToWrite.data.data, content...)
	return nil
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
			return fs.createFile(path.Base(), file.RegularFile, parentDir)
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
