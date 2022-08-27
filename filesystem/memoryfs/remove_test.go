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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target")

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
				assert.Len(t, files, 1)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
			},
		},
		{
			CaseName: "Remove file in subdir using absolute path",
			Path:     "/dir1/target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/dir1/target")

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
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

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
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

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
				assert.Len(t, files, 5)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Remove symlink using absolute path",
			Path:     "/target-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p1, _ := fspath.NewFileSystemPath("/target", nil)
				if _, err := fs.CreateRegularFile(p1, rootUser); err != nil {
					return nil, nil, err
				}
				p2, _ := fspath.NewFileSystemPath("/target-link", nil)

				if _, err := fs.CreateSymbolicLink(p1, p2, rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target-link")

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
				assert.Len(t, files, 2)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/target")

				p, _ = fspath.NewFileSystemPath("/", nil)
				files, _ = fs.FindFiles("target-link", p, rootUser)
				assert.Len(t, files, 0)
			},
		},
		{
			CaseName: "Remove file in root using relative path",
			Path:     "../../../target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir3/dir4/di5", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target")

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
				assert.Len(t, files, 1)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
			},
		},
		{
			CaseName: "Remove file in subdir using relative path",
			Path:     "./target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target/target/dir/target")

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
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

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
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

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath(workDir.Info().AbsolutePath(), nil)
				if _, err := fs.RemoveAll(p, rootUser); err != nil {
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

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
		p, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		info, err := fs.Remove(p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target")

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
				assert.Len(t, files, 1)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
			},
		},
		{
			CaseName: "Remove direcoty in subdir using absolute path",
			Path:     "/target/target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target/target")

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
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

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, info)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrOperationNotSupported)

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target", nil)
				if _, err := fs.Mkdir(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir3/dir4/di5", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target")

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
				assert.Len(t, files, 1)
				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
			},
		},
		{
			CaseName: "Remove directory in subdir using relative path",
			Path:     "..",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, info.AbsolutePath(), "/target/target")

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
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

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
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

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p, rootUser)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath(workDir.Info().AbsolutePath(), nil)
				if _, err := fs.RemoveAll(p, rootUser); err != nil {
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

				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("target", p, rootUser)
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
		p, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		info, err := fs.RemoveAll(p, rootUser)
		testCase.Assertions(t, fs, info, err)
	}
}
