package fspath

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"path/filepath"
	"strings"
)

type FileSystemPath struct {
	path       string
	workingDir file.File
}

func (p *FileSystemPath) WorkingDir() file.File {
	return p.workingDir
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

func (p *FileSystemPath) AbsolutePath() string {
	if p.IsAbs() {
		return p.path
	}

	return filepath.Clean(filepath.Join(p.workingDir.Info().AbsolutePath(), p.path))
}

func (p *FileSystemPath) Base() string {
	return filepath.Base(p.path)
}

func (p *FileSystemPath) Path() string {
	return p.path
}

func NewFileSystemPath(path string, workingDir file.File) (*FileSystemPath, error) {
	cleanPath := filepath.Clean(path)
	if workingDir == nil && !filepath.IsAbs(path) {
		return nil, fmt.Errorf("invalid path")
	}
	return &FileSystemPath{
		path:       cleanPath,
		workingDir: workingDir,
	}, nil
}
