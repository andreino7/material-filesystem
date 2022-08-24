package memoryfs_test

import (
	"errors"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: add tests for working dir deleted
func TestListFiles(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, []file.FileInfo, error)
	}{
		{
			CaseName: "List files in / - absolute path",
			Path:     "/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir3/dir4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir3/file4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/dir6/file6"), nil); err != nil {
					return nil, nil, err
				}

				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file8"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6/file9"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 2)
				assert.Equal(t, files[0].Name(), "dir1")
				assert.Equal(t, files[1].Name(), "dir2")
			},
		},
		{
			CaseName: "List files in subdirectory - absolute path",
			Path:     "/dir2/dir1/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir3/dir4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir3/file4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/dir6/file6"), nil); err != nil {
					return nil, nil, err
				}

				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file8"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6/file9"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 3)
				assert.Equal(t, files[0].Name(), "dir5")
				assert.Equal(t, files[1].Name(), "file3")
				assert.Equal(t, files[2].Name(), "file8")
			},
		},
		{
			CaseName: "No such file or directory - absolute path",
			Path:     "/dir2/dir10/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir3/dir4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)
				assert.NotNil(t, files)
				assert.Len(t, files, 0)
			},
		},
		{
			CaseName: "Listing regular file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
				assert.NotNil(t, files)
				assert.Len(t, files, 0)
			},
		},
		{
			CaseName: "List files in / - relative path",
			Path:     "../../../../../../..",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir3/dir4"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir3/file4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/dir6/file6"), nil); err != nil {
					return nil, nil, err
				}

				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file8"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6/file9"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 2)
				assert.Equal(t, files[0].Name(), "dir1")
				assert.Equal(t, files[1].Name(), "dir2")
			},
		},
		{
			CaseName: "List files in subdirectory - relative path",
			Path:     "../../../dir2/dir1/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir3/dir4"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir3/file4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/dir6/file6"), nil); err != nil {
					return nil, nil, err
				}

				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file8"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6/file9"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 3)
				assert.Equal(t, files[0].Name(), "dir5")
				assert.Equal(t, files[1].Name(), "file3")
				assert.Equal(t, files[2].Name(), "file8")
			},
		},
		{
			CaseName: "List files in symlink - relative path",
			Path:     "../../../dir2/dir1/dir5/dir6",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir3/dir4"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir3/file4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir5/dir6/file6"), nil); err != nil {
					return nil, nil, err
				}

				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir2/dir1/dir5/"), nil); err != nil {
					return nil, nil, err
				}

				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/dir2/dir1/"), fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/file8"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/file5"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6/file9"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 4)
				assert.Equal(t, files[0].Name(), "dir5")
				assert.Equal(t, files[1].Name(), "file3")
				assert.Equal(t, files[2].Name(), "file8")
				assert.Equal(t, files[3].Name(), "file9")
			},
		},
		{
			CaseName: "No such file or directory - relative path",
			Path:     "dir10/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir3/dir4"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)
				assert.NotNil(t, files)
				assert.Len(t, files, 0)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		files, err := fs.ListFiles(fspath.NewFileSystemPath(testCase.Path), workingDir)
		testCase.Assertions(t, files, err)
	}
}
