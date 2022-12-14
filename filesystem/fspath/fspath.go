package fspath

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"path/filepath"
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

// NewFileSystemPath creates a new filesystem path from
// the given path and workingDir
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
