package memoryfs_test

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDirectory(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, file.File, error)
	}{
		{
			CaseName: "Change working directory using absolute path - absolute path, work dir nil",
			Path:     "/dir5/dir6",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir5/dir6/dir7"), nil); err != nil {
					return nil, nil, err
				}

				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir6")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Change working directory using absolute path - no such directory",
			Path:     "/dir5/dir6/dir8",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir5/dir6/dir7"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, file)
				assert.Equal(t, err.Error(), "no such file or directory")
			},
		},
		{
			CaseName: "Change working directory using absolute path - regular file",
			Path:     "/dir1/dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir5/dir6/dir7"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir2/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, file)
				assert.Equal(t, err.Error(), "file is not a directory")
			},
		},
		{
			CaseName: "Change working directory - relative path (../../), work dir not nil",
			Path:     "../../",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir5/dir6/dir7"), nil)
				if err != nil {
					return nil, nil, err
				}

				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir5")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Change working directory - relative path (../../././.././dir3/../dir1/dir2/), work dir not nil",
			Path:     "../../././.././dir3/../dir1/dir2/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir5/dir6/dir7"), nil)
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
			CaseName: "Change working directory - relative path to before root, work dir not nil",
			Path:     "../../../../../../../../..",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir5/dir6/dir7"), nil)
				if err != nil {
					return nil, nil, err
				}

				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "/")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Change working directory using relative path - no such directory",
			Path:     "../../dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir5/dir6/dir7"), nil)
				if err != nil {
					return nil, nil, err
				}

				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, file)
				assert.Equal(t, err.Error(), "no such file or directory")
			},
		},
		{
			CaseName: "Change working directory using relative path - regular file",
			Path:     "file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workingDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil)
				if err != nil {
					return nil, nil, err
				}

				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir5/dir6/dir7"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir2/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workingDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, file)
				assert.Equal(t, err.Error(), "file is not a directory")
			},
		},
		{
			CaseName: "Working directory previously deleted",
			Path:     "file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workingDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil)
				if err != nil {
					return nil, nil, err
				}

				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir5/dir6/dir7"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir2/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.RemoveAll(fspath.NewFileSystemPath(workingDir.Info().AbsolutePath()), nil); err != nil {
					return nil, nil, err
				}
				return fs, workingDir, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, file)
				assert.Equal(t, err.Error(), "working directory deleted")
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		dir, err := fs.GetDirectory(fspath.NewFileSystemPath(testCase.Path), workingDir)
		testCase.Assertions(t, dir, err)
	}
}

func TestDefaultWorkingDirectory(t *testing.T) {
	cases := []struct {
		CaseName   string
		Assertions func(*testing.T, file.File)
	}{
		{
			CaseName: "Default working directory is /",
			Assertions: func(t *testing.T, file file.File) {
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "/")
				assert.True(t, file.Info().IsDirectory())
			},
		},
	}
	for _, testCase := range cases {
		fs := memoryfs.NewMemoryFileSystem()
		dir := fs.DefaultWorkingDirectory()
		testCase.Assertions(t, dir)
	}
}
