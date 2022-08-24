package memoryfs_test

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: add tests for working dir deleted
// TODO: add tests for symlink
func TestReadFile(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, []byte, error)
	}{
		{
			CaseName: "Read from existing file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendToFile(fspath.NewFileSystemPath("/file1"), []byte("Hello world!"), nil); err != nil {
					return nil, nil, err
				}

				return fs, nil, nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Read from missing file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
		{
			CaseName: "Read directory - absolute path",
			Path:     "/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
			},
		},
		{
			CaseName: "Read from existing file - relative path",
			Path:     "file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendToFile(fspath.NewFileSystemPath("/file1"), []byte("Hello world!"), nil); err != nil {
					return nil, nil, err
				}

				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Read from missing file - relative path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
		{
			CaseName: "Read directory - relative path",
			Path:     "./dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
			},
		},
		{
			CaseName: "Read from symbolic link - relative path",
			Path:     "file1-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendToFile(fspath.NewFileSystemPath("/file1"), []byte("Hello world!"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/file1"), fspath.NewFileSystemPath("/file1-link"), nil); err != nil {
					return nil, nil, err
				}

				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		path := fspath.NewFileSystemPath(testCase.Path)
		data, err := fs.ReadFile(path, workingDir)
		testCase.Assertions(t, data, err)
	}
}
