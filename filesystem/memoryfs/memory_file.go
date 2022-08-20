package memoryfs

import "material/filesystem/filesystem/file"

type inMemoryFileInfo struct {
	name         string
	absolutePath string
	isDirectory  bool
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

func (info inMemoryFileInfo) IsDirectory() bool {
	return info.isDirectory
}

func (info inMemoryFileInfo) Name() string {
	return info.name
}

func (info inMemoryFileInfo) AbsolutePath() string {
	return info.absolutePath
}

func newInMemoryFile(name string, absolutePath string, isDirectory bool) *inMemoryFile {
	info := &inMemoryFileInfo{
		name,
		absolutePath,
		isDirectory,
	}

	file := &inMemoryFile{
		info:    info,
		fileMap: map[string]*inMemoryFile{},
	}

	file.fileMap["."] = file
	return file
}
