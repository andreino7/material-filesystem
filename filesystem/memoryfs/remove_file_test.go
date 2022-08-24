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

func TestRemove(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error)
	}{
		{
			CaseName: "Remove file in root using absolute path",
			Path:     "/target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target")

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 1)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
			},
		},
		{
			CaseName: "Remove file in subdir using absolute path",
			Path:     "/dir1/target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/dir1/target")

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 4)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Remove directory using absolute path",
			Path:     "/target/target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidFileType)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 5)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Remove missing file using absolute path",
			Path:     "/dir1/dir3/target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 5)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Remove file in root using relative path",
			Path:     "../../../target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3/dir4/di5"), nil)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target")

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 1)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
			},
		},
		{
			CaseName: "Remove file in subdir using relative path",
			Path:     "./target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target/target/dir/target")

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 4)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
			},
		},
		{
			CaseName: "Remove directory using relative path",
			Path:     ".",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidFileType)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 5)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Remove missing file using releative path",
			Path:     "../target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 5)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Working directory previously deleted",
			Path:     "../target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.RemoveAll(fspath.NewFileSystemPath(workDir.Info().AbsolutePath()), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidWorkingDirectory)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 4)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		info, err := fs.Remove(fspath.NewFileSystemPath(testCase.Path), workingDir)
		testCase.Assertions(t, fs, info, err)
	}
}

func TestRemoveAll(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error)
	}{
		{
			CaseName: "Remove file in root using absolute path",
			Path:     "/target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target")

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 1)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
			},
		},
		{
			CaseName: "Remove direcoty in subdir using absolute path",
			Path:     "/target/target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target/target")

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 3)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
			},
		},
		{
			CaseName: "Remove missing directory using absolute path",
			Path:     "/dir1/dir3/target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 5)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Removing root using absolute path",
			Path:     "/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrOperationNotSupported)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 5)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Remove directory in root using relative path",
			Path:     "../../../target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir3/dir4/di5"), nil)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target")

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 1)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
			},
		},
		{
			CaseName: "Remove directory in subdir using relative path",
			Path:     "..",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target/target")

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 3)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
			},
		},
		{
			CaseName: "Remove missing file using releative path",
			Path:     "../target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 5)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Removing root using relative file",
			Path:     "../../../../..",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrOperationNotSupported)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 5)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Working directory previously deleted",
			Path:     "../target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.MkdirAll(fspath.NewFileSystemPath("/dir1/dir2/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/dir1/target"), nil); err != nil {
					return nil, nil, err
				}
				workDir, err := fs.MkdirAll(fspath.NewFileSystemPath("/target/target/dir"), nil)
				if err != nil {
					return nil, nil, err
				}
				if _, err := fs.CreateRegularFile(fspath.NewFileSystemPath("/target/target/dir/target"), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.RemoveAll(fspath.NewFileSystemPath(workDir.Info().AbsolutePath()), nil); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidWorkingDirectory)

				files, _ := fs.FindFiles("target", fspath.NewFileSystemPath("/"), nil)
				assert.Len(t, files, 4)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		info, err := fs.RemoveAll(fspath.NewFileSystemPath(testCase.Path), workingDir)
		testCase.Assertions(t, fs, info, err)
	}
}
