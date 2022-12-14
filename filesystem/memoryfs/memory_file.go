package memoryfs

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"path/filepath"
)

// inMemoryFile implements the FileInfo interface
type inMemoryFileInfo struct {
	absolutePath string
	fileType     file.FileType
}

// inMemoryFile implements the File interface
type inMemoryFile struct {
	info      *inMemoryFileInfo
	data      *inMemoryFileData
	isDeleted bool
	fileMap   map[string]*inMemoryFile
	link      *fspath.FileSystemPath
}

// Implement sort interface
type ByAbsolutePath []file.FileInfo

func (info ByAbsolutePath) Len() int {
	return len(info)
}
func (info ByAbsolutePath) Swap(i, j int) {
	info[i], info[j] = info[j], info[i]
}

func (info ByAbsolutePath) Less(i, j int) bool {
	return info[i].AbsolutePath() < info[j].AbsolutePath()
}

func (f inMemoryFile) Info() file.FileInfo {
	return f.info
}

func (f inMemoryFile) Data() file.FileData {
	return f.data
}

func (info *inMemoryFileInfo) FileType() file.FileType {
	return info.fileType
}

func (info *inMemoryFileInfo) Name() string {
	return filepath.Base(info.absolutePath)
}

func (info *inMemoryFileInfo) AbsolutePath() string {
	return info.absolutePath
}

func newInMemoryFile(absolutePath string, fileType file.FileType) *inMemoryFile {
	info := &inMemoryFileInfo{
		absolutePath: absolutePath,
		fileType:     fileType,
	}

	newFile := &inMemoryFile{
		info:    info,
		data:    &inMemoryFileData{},
		fileMap: map[string]*inMemoryFile{},
	}

	if fileType == file.RegularFile {
		newFile.data = &inMemoryFileData{}
	}

	newFile.fileMap["."] = newFile
	return newFile
}
