package memoryfs

import (
	"material/filesystem/filesystem/file"
	"path/filepath"
)

type inMemoryFileInfo struct {
	absolutePath string
	fileType     file.FileType
	userId       string
	groupId      string
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

func (info *inMemoryFileInfo) FileType() file.FileType {
	return info.fileType
}

func (info *inMemoryFileInfo) Name() string {
	return filepath.Base(info.absolutePath)
}

func (info *inMemoryFileInfo) AbsolutePath() string {
	return info.absolutePath
}

func (info *inMemoryFileInfo) UserId() string {
	return info.userId
}

func (info *inMemoryFileInfo) GroupId() string {
	return info.groupId
}
