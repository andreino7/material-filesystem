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
func TestFindFiles(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		FileName   string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(t *testing.T, files []file.FileInfo, err error)
	}{
		{
			CaseName: "Find all files and directory matching names using absolute path",
			Path:     "/",
			FileName: "target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)

				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 5)

				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Find all files and directory matching names using absolute path starting from subdir",
			Path:     "/target",
			FileName: "target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)

				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 2)

				assert.Equal(t, files[0].AbsolutePath(), "/target/target")
				assert.Equal(t, files[1].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "No such files using absolute path",
			Path:     "/",
			FileName: "invalid",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)

				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 0)
			},
		},
		{
			CaseName: "Invalid starting directory using absolute path",
			Path:     "/invalid",
			FileName: "target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)

				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
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
			CaseName: "Find all files and directory matching names using relative path",
			Path:     "../../../.",
			FileName: "target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				workDir, err := fs.MkdirAll(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)

				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 5)

				assert.Equal(t, files[0].AbsolutePath(), "/dir1/dir2/target")
				assert.Equal(t, files[1].AbsolutePath(), "/dir1/target")
				assert.Equal(t, files[2].AbsolutePath(), "/target")
				assert.Equal(t, files[3].AbsolutePath(), "/target/target")
				assert.Equal(t, files[4].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "Find all files and directory matching names using relative path starting from subdir",
			Path:     ".",
			FileName: "target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 1)

				assert.Equal(t, files[0].AbsolutePath(), "/target/target/dir/target")
			},
		},
		{
			CaseName: "No such files using relative path",
			Path:     "..",
			FileName: "dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 0)
			},
		},
		{
			CaseName: "Invalid starting directory using relative path",
			Path:     "../../dir/.",
			FileName: "target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
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
		{
			CaseName: "Working directory previously deleted",
			Path:     "../../dir/.",
			FileName: "target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir", nil)
				workDir, err := fs.MkdirAll(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath(workDir.Info().AbsolutePath(), nil)

				if _, err := fs.RemoveAll(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrInvalidWorkingDirectory)
				assert.NotNil(t, files)
				assert.Len(t, files, 0)
			},
		},
		{
			CaseName: "Find all files and directory matching names using symlink",
			Path:     "/target-link",
			FileName: "target",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/dir1/dir2/target", nil)
				if _, err := fs.MkdirAll(p1); err != nil {
					return nil, nil, err
				}
				p2, _ := fspath.NewFileSystemPath("/dir1/target", nil)
				if _, err := fs.CreateRegularFile(p2); err != nil {
					return nil, nil, err
				}
				p3, _ := fspath.NewFileSystemPath("/target/target/dir/target", nil)
				if _, err := fs.MkdirAll(p3); err != nil {
					return nil, nil, err
				}
				p4, _ := fspath.NewFileSystemPath("/target-link", nil)
				p5, _ := fspath.NewFileSystemPath("/target", nil)

				if _, err := fs.CreateSymbolicLink(p5, p4); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, files []file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, files)
				assert.Len(t, files, 2)

				assert.Equal(t, files[0].AbsolutePath(), "/target/target")
				assert.Equal(t, files[1].AbsolutePath(), "/target/target/dir/target")
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		p, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		dir, err := fs.FindFiles(testCase.FileName, p)
		testCase.Assertions(t, dir, err)
	}
}
