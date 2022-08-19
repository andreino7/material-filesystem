package fspath

import "path/filepath"

type FileSystemPath struct {
	absolutePath   string
	workingDirPath string
}

func (p *FileSystemPath) Dir() string {
	return filepath.Dir(p.absolutePath)
}

func (p *FileSystemPath) Base() string {
	return filepath.Base(p.absolutePath)
}

func (p *FileSystemPath) AbsolutePath() string {
	return p.absolutePath
}

func normalizePath(path string, workingDirPath string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	return filepath.Clean(filepath.Join(workingDirPath, path))
}

func NewFileSystemPath(path string, workingDirPath string) *FileSystemPath {
	absolutePath := normalizePath(path, workingDirPath)
	return &FileSystemPath{
		absolutePath,
		workingDirPath,
	}
}
