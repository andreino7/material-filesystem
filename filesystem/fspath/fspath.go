package fspath

import (
	"path/filepath"
	"strings"
)

type FileSystemPath struct {
	path           string
	absolutePath   string
	workingDirPath string
}

func (p *FileSystemPath) Dir() string {
	return filepath.Dir(p.path)
}

func (p *FileSystemPath) AbsDir() string {
	return filepath.Dir(p.absolutePath)
}

// TODO: add test
func (p *FileSystemPath) Split() ([]string, string) {
	return splitHelper(p.path)
}

func (p *FileSystemPath) SplitAbs() ([]string, string) {
	return splitHelper(p.absolutePath)
}

func splitHelper(path string) ([]string, string) {
	dir, file := filepath.Split(path)
	if dir == "/" {
		return []string{}, file
	}
	dirNames := strings.Split(strings.Trim(dir, "/"), "/")
	return dirNames, file
}

// TODO: add test
func (p *FileSystemPath) IsAbs() bool {
	return filepath.IsAbs(p.path)
}

func (p *FileSystemPath) Base() string {
	return filepath.Base(p.path)
}

func (p *FileSystemPath) AbsolutePath() string {
	return p.absolutePath
}

func (p *FileSystemPath) Path() string {
	return p.path
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
