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
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "dir1")
				assert.Equal(t, res.Info().AbsolutePath(), "/dir1")
				assert.Equal(t, res.Info().FileType(), file.Directory)
			},
		},
		{
			CaseName: "Create directory in subdir - absolute path, work dir nil",
			Path:     "/dir1/dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "dir2")
				assert.Equal(t, res.Info().AbsolutePath(), "/dir1/dir2")
				assert.Equal(t, res.Info().FileType(), file.Directory)
			},
		},
		{
			CaseName: "Name conflict - absolute path, work dir nil",
			Path:     "/dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2/dir3"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrExist)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Parent directory does not exist - absolute path, work dir nil",
			Path:     "/dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir4"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Parent is not a directory - absolute path, work dir nil",
			Path:     "/file1/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1/"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Create directory following symlink - absolute path",
			Path:     "/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/dir1"), fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "/dir1/dir3", res.Info().AbsolutePath())
			},
		},
		{
			CaseName: "Create directory in root - relative path, work dir not nil",
			Path:     "dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "dir1")
				assert.Equal(t, res.Info().AbsolutePath(), "/dir1")
				assert.Equal(t, res.Info().FileType(), file.Directory)
			},
		},
		{
			CaseName: "Create directory in subdir - relative path, work dir not nil",
			Path:     "../dir1/dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}

				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2"), nil)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "dir2")
				assert.Equal(t, res.Info().AbsolutePath(), "/dir1/dir2")
				assert.Equal(t, res.Info().FileType(), file.Directory)
			},
		},
		{
			CaseName: "Name conflict - relative path, work dir not nil",
			Path:     "./dir2/../dir2/.//dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2/dir3"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrExist)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Parent directory does not exist - relative path, work dir not nil",
			Path:     "./dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir4"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Parent is not a directory - relative path, work dir not nil",
			Path:     "file1/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1/"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Working directory previously deleted",
			Path:     "./dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.RemoveAll(fspath.NewFileSystemPath(workDir.Info().AbsolutePath()), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidWorkingDirectory)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Create directory following symlink - relative path",
			Path:     "../dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/dir1"), fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "/dir1/dir3", res.Info().AbsolutePath())
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		dir, err := fs.Mkdir(fspath.NewFileSystemPath(testCase.Path), workingDir)
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
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "file1")
				assert.Equal(t, res.Info().AbsolutePath(), "/file1")
				assert.Equal(t, res.Info().FileType(), file.RegularFile)
			},
		},
		{
			CaseName: "Create file in subdir - absolute path, work dir nil",
			Path:     "/dir1/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "file1")
				assert.Equal(t, res.Info().AbsolutePath(), "/dir1/file1")
				assert.Equal(t, res.Info().FileType(), file.RegularFile)
			},
		},
		{
			CaseName: "Name conflict in root with directory - absolute path, work dir nil",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrExist)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Name conflict in root with regular file - absolute path, work dir nil",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrExist)
				assert.Equal(t, err.Error(), "file already exists")
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Name conflict in subdirectory - absolute path, work dir nil",
			Path:     "/dir1/dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir2/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrExist)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Parent directory does not exists - absolute path, work dir nil",
			Path:     "/dir1/dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Create file following symlink - absolute path",
			Path:     "/dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/dir1"), fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "/dir1/file1", res.Info().AbsolutePath())
			},
		},
		{
			CaseName: "Create file in root - relative path, work dir not nil",
			Path:     "file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "file1")
				assert.Equal(t, res.Info().AbsolutePath(), "/file1")
				assert.Equal(t, res.Info().FileType(), file.RegularFile)
			},
		},
		{
			CaseName: "Create file in subdir - relative path, work dir nil",
			Path:     "./../dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "file1")
				assert.Equal(t, res.Info().AbsolutePath(), "/dir2/file1")
				assert.Equal(t, res.Info().FileType(), file.RegularFile)
			},
		},
		{
			CaseName: "Name conflict in subdirectory - relative path, work dir not nil",
			Path:     "dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/dir2/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file already exists")
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Parent directory does not exists - relative path, work dir not nil",
			Path:     "./dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Working directory previously deleted",
			Path:     "./dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.RemoveAll(fspath.NewFileSystemPath(workDir.Info().AbsolutePath()), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidWorkingDirectory)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Create file following symlink - relative path",
			Path:     "../dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/dir1"), fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "/dir1/file1", res.Info().AbsolutePath())
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		dir, err := fs.CreateRegularFile(fspath.NewFileSystemPath(testCase.Path), workingDir)
		testCase.Assertions(t, dir, err)
	}
}

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
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "dir3")
				assert.Equal(t, res.Info().AbsolutePath(), "/dir1/dir2/dir3")
				assert.Equal(t, res.Info().FileType(), file.Directory)
			},
		},
		{
			CaseName: "Parent is not a directory - absolute path, work dir nil",
			Path:     "/file1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1/"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Create directories - relative path, work dir not nil",
			Path:     "dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.Info().Name(), "dir3")
				assert.Equal(t, res.Info().AbsolutePath(), "/dir1/dir1/dir2/dir3")
				assert.Equal(t, res.Info().FileType(), file.Directory)
			},
		},
		{
			CaseName: "Parent is not a directory - relative path, work dir not nil",
			Path:     "./file1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1/"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
				assert.Nil(t, res)
			},
		},
		{
			CaseName: "Working directory previously deleted",
			Path:     "./dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir4"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.RemoveAll(fspath.NewFileSystemPath(workDir.Info().AbsolutePath()), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, res file.File, err error) {
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidWorkingDirectory)
				assert.Nil(t, res)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		dir, err := fs.MkdirAll(fspath.NewFileSystemPath(testCase.Path), workingDir)
		testCase.Assertions(t, dir, err)
	}
}

func TestCreateHardLink(t *testing.T) {
	cases := []struct {
		CaseName   string
		SrcPath    string
		DestPath   string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, *memoryfs.MemoryFileSystem, *fspath.FileSystemPath, file.FileInfo, error)
	}{
		{
			CaseName: "Hard link to file - absolute path",
			SrcPath:  "/file1",
			DestPath: "/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")

				// Modifying src file should modify hardlinked file
				fs.AppendAll(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Removing original file should not have any effect on hard link - absolute path",
			SrcPath:  "/file1",
			DestPath: "/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")

				// Modifying src file should modify hardlinked file
				fs.AppendAll(src, []byte("Hello world!"), nil)
				fs.Remove(src, nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Removing hard link should not have any effect on original file - absolute path",
			SrcPath:  "/file1",
			DestPath: "/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")

				// Modifying src file should modify hardlinked file
				fs.AppendAll(src, []byte("Hello world!"), nil)
				fs.Remove(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				data, _ := fs.ReadAll(src, nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Hard link to directory should fail - absolute path",
			SrcPath:  "/dir1",
			DestPath: "/dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, res)
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
			},
		},
		{
			CaseName: "Should create intermediate directories - absolute path",
			SrcPath:  "/file1",
			DestPath: "/dir2/dir3/dir4/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/dir2/dir3/dir4/file1")

				// Modifying src file should modify hardlinked file
				fs.AppendAll(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Should fail if target already exists - absolute path",
			SrcPath:  "/file1",
			DestPath: "/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, res)
				assert.Equal(t, err, fserrors.ErrExist)
			},
		},
		{
			CaseName: "Hard link to symlink file - absolute path",
			SrcPath:  "/file2",
			DestPath: "/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/file1"), fspath.NewFileSystemPath("/file2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file3")

				// Modifying src file should modify hardlinked file
				fs.AppendAll(fspath.NewFileSystemPath("/file1"), []byte("Hello world!"), nil)
				fs.Remove(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				data, _ := fs.ReadAll(src, nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Hard link to symlink directory should fail - absolute path",
			SrcPath:  "/dir2",
			DestPath: "/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/dir1"), fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, res)
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
			},
		},
		{
			CaseName: "Hard link to file - relative path",
			SrcPath:  "file1",
			DestPath: "file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")

				// Modifying src file should modify hardlinked file
				fs.AppendAll(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Removing original file should not have any effect on hard link - relative path",
			SrcPath:  "file1",
			DestPath: "file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")

				// Modifying src file should modify hardlinked file
				fs.AppendAll(src, []byte("Hello world!"), nil)
				fs.Remove(src, nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Hard link to directory should fail - relative path",
			SrcPath:  "dir1",
			DestPath: "dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, res)
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
			},
		},
		{
			CaseName: "Should create intermediate directories - relative path",
			SrcPath:  "file1",
			DestPath: "./dir2/dir3/dir4/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/dir2/dir3/dir4/file1")

				// Modifying src file should modify hardlinked file
				fs.AppendAll(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Should fail if target already exists - relative path",
			SrcPath:  "file1",
			DestPath: "file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, res)
				assert.Equal(t, err, fserrors.ErrExist)
			},
		},
		{
			CaseName: "Hard link to symlink file - relative path",
			SrcPath:  "file2",
			DestPath: "file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/file1"), fspath.NewFileSystemPath("/file2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file3")

				// Modifying src file should modify hardlinked file
				fs.AppendAll(fspath.NewFileSystemPath("/file1"), []byte("Hello world!"), nil)
				fs.Remove(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				data, _ := fs.ReadAll(src, nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Hard link to symlink directory should fail - relative path",
			SrcPath:  "dir2",
			DestPath: "file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/dir1"), fspath.NewFileSystemPath("/dir2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, res)
				assert.Equal(t, err, fserrors.ErrInvalidFileType)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		srcPath := fspath.NewFileSystemPath(testCase.SrcPath)
		destPath := fspath.NewFileSystemPath(testCase.DestPath)

		res, err := fs.CreateHardLink(srcPath, destPath, workingDir)
		testCase.Assertions(t, fs, srcPath, res, err)
	}
}

func TestCreateSymbolicLink(t *testing.T) {
	cases := []struct {
		CaseName   string
		SrcPath    string
		DestPath   string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, *memoryfs.MemoryFileSystem, *fspath.FileSystemPath, file.FileInfo, file.File, error)
	}{
		{
			CaseName: "Symbolic link to file - absolute path",
			SrcPath:  "/file1",
			DestPath: "/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")

				fs.AppendAll(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Symbolic link to directory - absolute path",
			SrcPath:  "/dir1",
			DestPath: "/dir-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/dir-link")

				filesLink, _ := fs.ListFiles(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				filesOriginal, _ := fs.ListFiles(src, nil)

				assert.Equal(t, filesLink, filesOriginal)
			},
		},
		{
			CaseName: "Writing a symbolic link should write the original file - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/file3-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file3-link")

				fs.AppendAll(fspath.NewFileSystemPath(res.AbsolutePath()), []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(src, nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Symbolic link to symbolic link - absolute path",
			SrcPath:  "/dir1/file1-link",
			DestPath: "/file3-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/dir1/file1"), fspath.NewFileSystemPath("/dir1/file1-link"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file3-link")

				fs.AppendAll(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Symbolic to file that does not exist - absolute path",
			SrcPath:  "/file1",
			DestPath: "/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")

				_, err = fs.ListFiles(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
		{
			CaseName: "Removing symlink source should cause ErrNotExist when reading symlink - absolute path",
			SrcPath:  "/file1",
			DestPath: "/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")
				fs.Remove(src, nil)
				_, err = fs.GetDirectory(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
		{
			CaseName: "Moving symlink source should cause ErrNotExist when reading symlink - absolute path",
			SrcPath:  "/file1",
			DestPath: "/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")
				fs.Move(src, fspath.NewFileSystemPath("/file3"), nil)
				_, err = fs.GetDirectory(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
		{
			CaseName: "Symbolic link to file - relative path",
			SrcPath:  "file1",
			DestPath: "file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")

				fs.AppendAll(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Symbolic link to directory - relative path",
			SrcPath:  "..",
			DestPath: "../../dir-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil)
				if err != nil {
					return nil, nil, err
				}

				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file3"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/dir-link")

				filesLink, _ := fs.ListFiles(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				filesOriginal, _ := fs.ListFiles(src, workDir)

				assert.Equal(t, filesLink, filesOriginal)
			},
		},
		{
			CaseName: "Writing a symbolic link should write the original file - relative path",
			SrcPath:  "../dir1/file1",
			DestPath: "file3-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file2"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file3-link")

				fs.AppendAll(fspath.NewFileSystemPath(res.AbsolutePath()), []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(src, nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Symbolic link to symbolic link - relative path",
			SrcPath:  "./dir1/file1-link",
			DestPath: "file3-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file1"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/file2"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateSymbolicLink(fspath.NewFileSystemPath("/dir1/file1"), fspath.NewFileSystemPath("/dir1/file1-link"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file3-link")

				fs.AppendAll(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadAll(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Symbolic link to file that does not exist - relative path",
			SrcPath:  "file1",
			DestPath: "file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")

				_, err = fs.ListFiles(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
		{
			CaseName: "Removing symlink source should cause ErrNotExist when reading symlink - relative path",
			SrcPath:  "./file1",
			DestPath: "file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")
				fs.Remove(src, nil)
				_, err = fs.GetDirectory(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
		{
			CaseName: "Moving symlink source should cause ErrNotExist when reading symlink - relative path",
			SrcPath:  "file1",
			DestPath: "file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/file1"), nil); err != nil {
					return nil, nil, err
				}
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, src *fspath.FileSystemPath, res file.FileInfo, workDir file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res.AbsolutePath(), "/file2")
				fs.Move(src, fspath.NewFileSystemPath("/file3"), nil)
				_, err = fs.GetDirectory(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		srcPath := fspath.NewFileSystemPath(testCase.SrcPath)
		destPath := fspath.NewFileSystemPath(testCase.DestPath)

		res, err := fs.CreateSymbolicLink(srcPath, destPath, workingDir)
		testCase.Assertions(t, fs, srcPath, res, workingDir, err)
	}
}
