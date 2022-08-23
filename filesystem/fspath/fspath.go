package fspath

import (
	"path/filepath"
	"strings"
)

type FileSystemPath struct {
	path string
}

func (p *FileSystemPath) Dir() string {
	return filepath.Dir(p.path)
}

// TODO: add test
func (p *FileSystemPath) Split() ([]string, string) {
	return splitHelper(p.path)
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

func (p *FileSystemPath) Path() string {
	return p.path
}

func NewFileSystemPath(path string) *FileSystemPath {
	return &FileSystemPath{
		path: filepath.Clean(path),
	}
}
