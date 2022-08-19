package memoryfs_test

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: add tests for working dir deleted
// TODO: check file exact location instead of just name
func TestMkdir(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, file.File, error)
	}{
		{
			CaseName: "Create directory in root - absolute path, work dir nil",
			Path:     "/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				return memoryfs.NewMemoryFileSystem(), nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir1")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Create directory in subdir - absolute path, work dir nil",
			Path:     "/dir1/dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir2")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Name conflict - absolute path, work dir nil",
			Path:     "/dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2/dir3", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file already exists")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Parent directory does not exist - absolute path, work dir nil",
			Path:     "/dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir4", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "no such file or directory")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Parent is not a directory - absolute path, work dir nil",
			Path:     "/file1/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1/", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file is not a directory")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Create directory in root - relative path, work dir not nil",
			Path:     "dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir1")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Create directory in subdir - relative path, work dir not nil",
			Path:     "../dir1/dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil); err != nil {
					return nil, nil, err
				}

				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2", ""), nil)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir2")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Name conflict - relative path, work dir not nil",
			Path:     "./dir2/../dir2/.//dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2/dir3", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file already exists")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Parent directory does not exist - relative path, work dir not nil",
			Path:     "./dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir4", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "no such file or directory")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Parent is not a directory - relative path, work dir not nil",
			Path:     "file1/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1/", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file is not a directory")
				assert.Nil(t, file)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		dir, err := fs.Mkdir(fspath.NewFileSystemPath(testCase.Path, ""), workingDir)
		testCase.Assertions(t, dir, err)
	}
}

func TestCreateRegularFile(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, file.File, error)
	}{
		{
			CaseName: "Create file in root - absolute path, work dir nil",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				return memoryfs.NewMemoryFileSystem(), nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "file1")
				assert.False(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Create file in subdir - absolute path, work dir nil",
			Path:     "/dir1/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "file1")
				assert.False(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Name conflict in root with directory - absolute path, work dir nil",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/file1", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file already exists")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Name conflict in root with regular file - absolute path, work dir nil",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file already exists")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Name conflict in subdirectory - absolute path, work dir nil",
			Path:     "/dir1/dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir2/file1", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file already exists")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Parent directory does not exists - absolute path, work dir nil",
			Path:     "/dir1/dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "no such file or directory")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Create file in root - relative path, work dir not nil",
			Path:     "file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "file1")
				assert.False(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Create file in subdir - relative path, work dir nil",
			Path:     "./../dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "file1")
				assert.False(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Name conflict in subdirectory - relative path, work dir not nil",
			Path:     "dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir2/file1", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file already exists")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Parent directory does not exists - relative path, work dir not nil",
			Path:     "./dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "no such file or directory")
				assert.Nil(t, file)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		dir, err := fs.CreateRegularFile(fspath.NewFileSystemPath(testCase.Path, ""), workingDir)
		testCase.Assertions(t, dir, err)
	}
}

// TODO: check that every directory in path exists
func TestMkdirAll(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, file.File, error)
	}{
		{
			CaseName: "Create directories in root - absolute path, work dir nil",
			Path:     "/dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				return memoryfs.NewMemoryFileSystem(), nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir3")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Parent is not a directory - absolute path, work dir nil",
			Path:     "/file1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1/", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file is not a directory")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Create directories - relative path, work dir not nil",
			Path:     "dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir3")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Parent is not a directory - relative path, work dir not nil",
			Path:     "./file1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1/", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file is not a directory")
				assert.Nil(t, file)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		dir, err := fs.MkdirAll(fspath.NewFileSystemPath(testCase.Path, ""), workingDir)
		testCase.Assertions(t, dir, err)
	}
}
