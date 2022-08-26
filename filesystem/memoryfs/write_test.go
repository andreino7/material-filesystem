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
func TestAppendAll(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, []byte, error)
	}{
		{
			CaseName: "Append to existing file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		// TODO: fix this
		// {
		// 	CaseName: "Append to new file - absolute path",
		// 	Path:     "/file1",
		// 	Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
		// 		fs := memoryfs.NewMemoryFileSystem()
		// 		return fs, nil, nil
		// 	},
		// 	Assertions: func(t *testing.T, data []byte, err error) {
		// 		assert.Nil(t, err)
		// 		assert.Equal(t, data, []byte("Hello world!"))
		// 	},
		// },
		// TODO: fix this
		// {
		// 	CaseName: "Append to new file and create intermediate directories - absolute path",
		// 	Path:     "/dir1/dir2/file1",
		// 	Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
		// 		fs := memoryfs.NewMemoryFileSystem()
		// 		return fs, nil, nil
		// 	},
		// 	Assertions: func(t *testing.T, data []byte, err error) {
		// 		assert.Nil(t, err)
		// 		assert.Equal(t, data, []byte("Hello world!"))
		// 	},
		// },
		{
			CaseName: "Append to directory - absolute path",
			Path:     "/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
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
			CaseName: "Append to existing file - relative path",
			Path:     "file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		// TODO: fix this
		// {
		// 	CaseName: "Append to new file - relative path",
		// 	Path:     "./file1",
		// 	Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
		// 		fs := memoryfs.NewMemoryFileSystem()
		// 		return fs, fs.DefaultWorkingDirectory(), nil
		// 	},
		// 	Assertions: func(t *testing.T, data []byte, err error) {
		// 		assert.Nil(t, err)
		// 		assert.Equal(t, data, []byte("Hello world!"))
		// 	},
		// },
		// TODO: fix this
		// {
		// 	CaseName: "Append to new file and create intermediate directories - relative path",
		// 	Path:     "dir1/dir2/file1",
		// 	Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
		// 		fs := memoryfs.NewMemoryFileSystem()
		// 		return fs, fs.DefaultWorkingDirectory(), nil
		// 	},
		// 	Assertions: func(t *testing.T, data []byte, err error) {
		// 		assert.Nil(t, err)
		// 		assert.Equal(t, data, []byte("Hello world!"))
		// 	},
		// },
		{
			CaseName: "Append to directory - relative path",
			Path:     "dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
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
			CaseName: "Append to symlink - absolute path",
			Path:     "/file1-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/file1", nil)

				if _, err := fs.CreateRegularFile(p1); err != nil {
					return nil, nil, err
				}

				p2, _ := fspath.NewFileSystemPath("/file1-link", nil)
				if _, err := fs.CreateSymbolicLink(p1, p2); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
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
		path, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		err = fs.AppendAll(path, []byte("Hello world!"))
		if err != nil {
			testCase.Assertions(t, nil, err)
		} else {
			data, _ := fs.ReadAll(path)
			testCase.Assertions(t, data, err)
		}
	}
}
