package fspath

import "path/filepath"

type FileSystemPath struct {
	path           string
	absolutePath   string
	workingDirPath string
}

func (p *FileSystemPath) Dir() string {
	return filepath.Dir(p.path)
}

func (p *FileSystemPath) Base() string {
	return filepath.Base(p.path)
}

func (p *FileSystemPath) AbsolutePath() string {
	return p.absolutePath
}

func toAbsolutePath(path string, workingDirPath string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	return filepath.Clean(filepath.Join(workingDirPath, path))
}

func NewFileSystemPath(path string, workingDirPath string) *FileSystemPath {
	absolutePath := toAbsolutePath(path, workingDirPath)
	return &FileSystemPath{
		path:           filepath.Clean(path),
		absolutePath:   absolutePath,
		workingDirPath: workingDirPath,
	}
}
