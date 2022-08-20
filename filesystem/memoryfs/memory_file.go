package memoryfs

import "material/filesystem/filesystem/file"

type inMemoryFileInfo struct {
	name        string
	isDirectory bool
}

// TODO: uncomment. Commented to pass static check
// type inMemoryFileData struct {
// 	data []byte
// }

type inMemoryFile struct {
	info      file.FileInfo
	data      file.FileData
	isDeleted bool
	fileMap   map[string]*inMemoryFile
}

func (f inMemoryFile) Info() file.FileInfo {
	return f.info
}

func (f inMemoryFile) Data() file.FileData {
	return f.data
}

func (info inMemoryFileInfo) IsDirectory() bool {
	return info.isDirectory
}

func (info inMemoryFileInfo) Name() string {
	return info.name
}

func newInMemoryFile(name string, isDirectory bool) *inMemoryFile {
	info := &inMemoryFileInfo{
		name,
		isDirectory,
	}

	file := &inMemoryFile{
		info:    info,
		fileMap: map[string]*inMemoryFile{},
	}

	file.fileMap["."] = file
	return file
}
