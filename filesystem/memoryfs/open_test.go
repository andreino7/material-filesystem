package memoryfs_test

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, *memoryfs.MemoryFileSystem, string, error)
	}{
		{
			CaseName: "Open a file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, fd string, err error) {
				assert.Nil(t, err)
				assert.NotEmpty(t, fd)
			},
		},
		{
			CaseName: "Open a directory should give an error - relative path",
			Path:     "dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, fd string, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fserrors.ErrInvalidFileType, err)
				assert.Empty(t, fd)
			},
		},
		{
			CaseName: "Open a symilink should open original file - relative path",
			Path:     "file1-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p1, _ := fspath.NewFileSystemPath("/file1-link", nil)
				if _, err := fs.CreateSymbolicLink(p, p1); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, fd string, err error) {
				assert.Nil(t, err)
				assert.NotEmpty(t, fd)
			},
		},
		{
			CaseName: "Open a file that does not exist should give an error - absolute path",
			Path:     "/fileeeee",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, fd string, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fserrors.ErrNotExist, err)
				assert.Empty(t, fd)
			},
		},
		{
			CaseName: "Opening the same file multiple time should return different fd - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, fd string, err error) {
				assert.Nil(t, err)
				assert.NotEmpty(t, fd)

				p, _ := fspath.NewFileSystemPath("/file1", nil)
				fd1, err := fs.Open(p)
				assert.Nil(t, err)
				assert.NotEmpty(t, fd1)
				assert.NotEqual(t, fd, fd1)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		path, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		fd, err := fs.Open(path)
		testCase.Assertions(t, fs, fd, err)
	}
}
