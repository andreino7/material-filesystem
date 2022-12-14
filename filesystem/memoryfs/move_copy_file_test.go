package memoryfs_test

import (
	"errors"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMove(t *testing.T) {
	cases := []struct {
		CaseName   string
		SrcPath    string
		DestPath   string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, *memoryfs.MemoryFileSystem, file.FileInfo, error)
	}{
		{
			CaseName: "Moving root directory is not allowed",
			SrcPath:  "/",
			DestPath: "/dir1/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrOperationNotSupported)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "File not found",
			SrcPath:  "/dir10",
			DestPath: "/dir1/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "Source dir and dest dir are same dir",
			SrcPath:  "/dir1",
			DestPath: "/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrSameFile)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "Move to subdir",
			SrcPath:  "/dir1",
			DestPath: "/dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/dir3/dir4", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrInvalid)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "Move to subdir and rename",
			SrcPath:  "/dir1",
			DestPath: "/dir1/dir2/dir3/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir2/dir3/dir4", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir2/dir3/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrInvalid)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "Rename file, no conflict - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir1/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir1/file3")
				assert.Equal(t, info.Name(), "file3")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
			},
		},
		{
			CaseName: "Rename file, conflict - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir1/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.True(t, strings.HasPrefix(info.AbsolutePath(), "/dir1/file2"))
				assert.NotEqual(t, info.AbsolutePath(), "/dir1/file2")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
				files, _ = fs.FindFiles(info.Name(), p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file2", p)
				assert.Len(t, files, 2)
				assert.NotEqual(t, files[0].Name(), files[1].Name())
			},
		},
		{
			CaseName: "Move file, no conflict, no rename - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/file1")
				assert.Equal(t, info.Name(), "file1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 1)
			},
		},
		{
			CaseName: "Move file, create intermediate directories - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir2/dir4/dir5/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir4/dir5/file1")
				assert.Equal(t, info.Name(), "file1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 1)
			},
		},
		{
			CaseName: "Move file, conflict - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.True(t, strings.HasPrefix(info.AbsolutePath(), "/dir2/file1"))
				assert.NotEqual(t, info.AbsolutePath(), "/dir2/file1")
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles(info.Name(), p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file1", p)
				assert.Len(t, files, 2)
				assert.NotEqual(t, files[0], files[1])
			},
		},
		{
			CaseName: "Move file and rename, no conflict - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/file3")
				assert.Equal(t, info.Name(), "file3")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
				files, _ = fs.FindFiles("file3", p)
				assert.Len(t, files, 1)
			},
		},
		{
			CaseName: "Move file and rename, conflict - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir2/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.True(t, strings.HasPrefix(info.AbsolutePath(), "/dir2/file2"))
				assert.NotEqual(t, info.AbsolutePath(), "/dir2/file2")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles(info.Name(), p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file2", p)
				assert.Len(t, files, 2)
				assert.NotEqual(t, files[0].Name(), files[1].Name())
			},
		},
		{
			CaseName: "Move directory, no conflict - absolute path",
			SrcPath:  "/dir1/",
			DestPath: "/dir2/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)

				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)
			},
		},
		// /**
		// FS structure after initialization:
		// /
		// dir1
		// 	file3
		// 	dir3
		// 		file4
		// 		dir4
		// 	dir5
		// 		file5
		// 		dir6
		// 			file6

		// dir2
		// 	dir1
		// 		file3
		// 		file8
		// 		dir5
		// 			file5
		// 			dir6
		// 				file9

		// expected FS structure after move:
		// /
		// dir2
		// 	dir1
		// 		file3
		// 		file3_uid
		// 		file8
		// 		dir3
		// 			file4
		// 			dir4
		// 		dir5
		// 			file5
		// 			file5_uid
		// 			dir6
		// 				file6
		// 				file9

		// */
		{
			CaseName: "Move directory, conflict - absolute path",
			SrcPath:  "/dir1/",
			DestPath: "/dir2/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir3/dir4", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/dir6", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file3", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir3/file4", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/file5", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/dir6/file6", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/file3", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/file8", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/file5", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6/file9", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			// TODO: ues list files to check the fs structure
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)

				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)

				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6", nil)
				fs.ListFiles(p)
			},
		},
		{
			CaseName: "Moving symlink - absolute path",
			SrcPath:  "/dir2/",
			DestPath: "/dir3/dir4/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p1)
				if err != nil {
					return nil, nil, err
				}
				p2, _ := fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.CreateSymbolicLink(p1, p2); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir3/dir4")
				p, _ := fspath.NewFileSystemPath("/dir3/dir4", nil)

				dir, _ := fs.GetDirectory(p)
				assert.Equal(t, dir.Info().AbsolutePath(), "/dir1")
			},
		},
		{
			CaseName: "Moving root directory is not allowed - relative path",
			SrcPath:  "..",
			DestPath: "/dir1/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrOperationNotSupported)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "File not found - relative path",
			SrcPath:  "dir5",
			DestPath: "/dir1/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
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
			},
		},
		{
			CaseName: "Rename file, no conflict - relative path",
			SrcPath:  "../dir1/file1",
			DestPath: "./file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir1/file3")
				assert.Equal(t, info.Name(), "file3")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
			},
		},
		{
			CaseName: "Rename file, conflict - relative path",
			SrcPath:  "file1",
			DestPath: "file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.True(t, strings.HasPrefix(info.AbsolutePath(), "/dir1/file2"))
				assert.NotEqual(t, info.AbsolutePath(), "/dir1/file2")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
				files, _ = fs.FindFiles(info.Name(), p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file2", p)
				assert.Len(t, files, 2)
				assert.NotEqual(t, files[0].Name(), files[1].Name())
			},
		},
		{
			CaseName: "Move file, no conflict, no rename - relative path",
			SrcPath:  "file1",
			DestPath: "../dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/file1")
				assert.Equal(t, info.Name(), "file1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 1)
			},
		},
		{
			CaseName: "Move file, create intermediate directories - relative path",
			SrcPath:  "file1",
			DestPath: "/dir2/dir4/dir5/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir4/dir5/file1")
				assert.Equal(t, info.Name(), "file1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 1)
			},
		},
		{
			CaseName: "Move file, conflict - relative path",
			SrcPath:  "./file1",
			DestPath: "../dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.True(t, strings.HasPrefix(info.AbsolutePath(), "/dir2/file1"))
				assert.NotEqual(t, info.AbsolutePath(), "/dir2/file1")
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles(info.Name(), p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file1", p)
				assert.Len(t, files, 2)
				assert.NotEqual(t, files[0], files[1])
			},
		},
		{
			CaseName: "Move file and rename, no conflict - relative path",
			SrcPath:  "file1",
			DestPath: "../file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/file3")
				assert.Equal(t, info.Name(), "file3")
				p, _ := fspath.NewFileSystemPath("/", nil)

				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
				files, _ = fs.FindFiles("file3", p)
				assert.Len(t, files, 1)
			},
		},
		{
			CaseName: "Move file and rename, conflict - relative path",
			SrcPath:  "/dir1/file1",
			DestPath: "./../dir2/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.True(t, strings.HasPrefix(info.AbsolutePath(), "/dir2/file2"))
				assert.NotEqual(t, info.AbsolutePath(), "/dir2/file2")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 0)
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles(info.Name(), p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file2", p)
				assert.Len(t, files, 2)
				assert.NotEqual(t, files[0].Name(), files[1].Name())
			},
		},
		{
			CaseName: "Move directory, no conflict - relative path",
			SrcPath:  ".",
			DestPath: "../dir2/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)

				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)
			},
		},
		/**
		FS structure after initialization:
		/
		dir1
			file3
			dir3
				file4
				dir4
			dir5
				file5
				dir6
					file6

		dir2
			dir1
				file3
				file8
				dir5
					file5
					dir6
						file9

		expected FS structure after move:
		/
		dir2
			dir1
				file3
				file3_uid
				file8
				dir3
					file4
					dir4
				dir5
					file5
					file5_uid
					dir6
						file6
						file9

		*/
		{
			CaseName: "Move directory, conflict - relative path",
			SrcPath:  ".",
			DestPath: "../dir2/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir3/dir4", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/dir6", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file3", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir3/file4", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/file5", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/dir6/file6", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/file3", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/file8", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/file5", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6/file9", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			// TODO: ues list files to check the fs structure
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)

				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		srcPath, _ := fspath.NewFileSystemPath(testCase.SrcPath, workingDir)
		destPath, _ := fspath.NewFileSystemPath(testCase.DestPath, workingDir)

		file, err := fs.Move(srcPath, destPath)
		testCase.Assertions(t, fs, file, err)
	}
}

// TODO: check file content for equality
func TestCopy(t *testing.T) {
	cases := []struct {
		CaseName   string
		SrcPath    string
		DestPath   string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, *memoryfs.MemoryFileSystem, file.FileInfo, error)
	}{
		{
			CaseName: "Copy root is not allowed",
			SrcPath:  "/",
			DestPath: "/dir1/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrOperationNotSupported)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "File not found",
			SrcPath:  "/dir10",
			DestPath: "/dir1/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrNotExist)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "Source dir and dest dir are same dir",
			SrcPath:  "/dir1",
			DestPath: "/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrSameFile)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "Copy and rename file, no conflict - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir1/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir1/file3")
				assert.Equal(t, info.Name(), "file3")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file3", p)
				assert.Len(t, files, 1)
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				b1, _ := fs.ReadAll(p)
				p, _ = fspath.NewFileSystemPath("/dir1/file3", nil)
				b2, _ := fs.ReadAll(p)
				assert.Equal(t, b1, b2)

			},
		},
		{
			CaseName: "Copy file and rename, conflict - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir1/file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.True(t, strings.HasPrefix(info.AbsolutePath(), "/dir1/file2"))
				assert.NotEqual(t, info.AbsolutePath(), "/dir1/file2")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles(info.Name(), p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file2", p)
				assert.Len(t, files, 2)
				assert.NotEmptyf(t, files[0].Name(), files[1].Name())
			},
		},
		{
			CaseName: "Copy file, no conflict, no rename - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/file1")
				assert.Equal(t, info.Name(), "file1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 2)
				p, _ = fspath.NewFileSystemPath("/dir1", nil)
				filesOldDir, _ := fs.FindFiles("file1", p)
				assert.Len(t, filesOldDir, 1)
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				filesNewDir, _ := fs.FindFiles("file1", p)
				assert.Len(t, filesNewDir, 1)

				// make sure they are not the same file
				assert.NotEqual(t, filesNewDir[0], filesOldDir[0])
			},
		},
		{
			CaseName: "Copy file, create intermediate directories - absolute path",
			SrcPath:  "/dir1/file1",
			DestPath: "/dir2/dir4/dir5/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir4/dir5/file1")
				assert.Equal(t, info.Name(), "file1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 2)

				p, _ = fspath.NewFileSystemPath("/dir1", nil)
				filesOldDir, _ := fs.FindFiles("file1", p)
				assert.Len(t, filesOldDir, 1)
				p, _ = fspath.NewFileSystemPath("/dir2/dir4/dir5", nil)
				filesNewDir, _ := fs.FindFiles("file1", p)
				assert.Len(t, filesNewDir, 1)

				// make sure they are not the same file
				assert.NotEqual(t, filesNewDir[0], filesOldDir[0])
			},
		},
		{
			CaseName: "Copy directory, no conflict - absolute path",
			SrcPath:  "/dir1/",
			DestPath: "/dir2/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir1")
				p, _ := fspath.NewFileSystemPath("/dir2", nil)
				files, _ := fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)

				p, _ = fspath.NewFileSystemPath("/", nil)
				files, _ = fs.FindFiles("dir1", p)
				assert.Len(t, files, 2)

				file1 := files[0]
				file2 := files[1]
				assert.NotEqual(t, file1, file2)
				checkCopy(t, fs, file1.AbsolutePath(), file2.AbsolutePath())
			},
		},
		// /**
		// FS structure after initialization:
		// /
		// dir1
		// 	file3
		// 	dir3
		// 		file4
		// 		dir4
		// 	dir5
		// 		file5
		// 		dir6
		// 			file6

		// dir2
		// 	dir1
		// 		file3
		// 		file8
		// 		dir5
		// 			file5
		// 			dir6
		// 				file9

		// expected FS structure after move:
		// /
		// dir1
		// 	file3
		// 	dir3
		// 		file4
		// 		dir4
		// 	dir5
		// 		file5
		// 		dir6
		// 			file6
		// dir2
		// 	dir1
		// 		file3
		// 		file3_uid
		// 		file8
		// 		dir3
		// 			file4
		// 			dir4
		// 		dir5
		// 			file5
		// 			file5_uid
		// 			dir6
		// 				file6
		// 				file9

		// */
		{
			CaseName: "Copy directory, conflict - absolute path",
			SrcPath:  "/dir1/",
			DestPath: "/dir2/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1/dir3/dir4", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/dir6", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file3", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir3/file4", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/file5", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/dir6/file6", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/file3", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/file8", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/file5", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6/file9", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			// TODO: ues list files to check the fs structure
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("dir1", p)
				assert.Len(t, files, 2)

				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)
			},
		},
		{
			CaseName: "Moving root directory is not allowed - relative path",
			SrcPath:  "..",
			DestPath: "/dir1/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.NotNil(t, err)
				target := &fserrors.FileSystemError{}
				assert.True(t, errors.As(err, &target))
				assert.Equal(t, err, fserrors.ErrOperationNotSupported)
				assert.Nil(t, info)
			},
		},
		{
			CaseName: "File not found - relative path",
			SrcPath:  "dir5",
			DestPath: "/dir1/file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
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
			},
		},
		{
			CaseName: "Copy and rename file, no conflict - relative path",
			SrcPath:  "../dir1/file1",
			DestPath: "./file3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir1/file3")
				assert.Equal(t, info.Name(), "file3")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file1", p)
				assert.Len(t, files, 1)
			},
		},
		{
			CaseName: "Copy file and rename, conflict - relative path",
			SrcPath:  "file1",
			DestPath: "file2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.True(t, strings.HasPrefix(info.AbsolutePath(), "/dir1/file2"))
				assert.NotEqual(t, info.AbsolutePath(), "/dir1/file2")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles(info.Name(), p)
				assert.Len(t, files, 1)
				files, _ = fs.FindFiles("file2", p)
				assert.Len(t, files, 2)
				assert.NotEqual(t, files[0].Name(), files[1].Name())
			},
		},
		{
			CaseName: "Copy file, no conflict, no rename - relative path",
			SrcPath:  "file1",
			DestPath: "../dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/file1")
				assert.Equal(t, info.Name(), "file1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 2)
				p, _ = fspath.NewFileSystemPath("/dir1", nil)
				filesOldDir, _ := fs.FindFiles("file1", p)
				assert.Len(t, filesOldDir, 1)
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				filesNewDir, _ := fs.FindFiles("file1", p)
				assert.Len(t, filesNewDir, 1)

				// make sure they are not the same file
				assert.NotEqual(t, filesNewDir[0], filesOldDir[0])
			},
		},
		{
			CaseName: "Copy file, create intermediate directories - relative path",
			SrcPath:  "./../dir1/file1",
			DestPath: "../dir2/dir4/dir5/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir4/dir5/file1")
				assert.Equal(t, info.Name(), "file1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("file1", p)
				assert.Len(t, files, 2)

				p, _ = fspath.NewFileSystemPath("/dir1", nil)
				filesOldDir, _ := fs.FindFiles("file1", p)
				assert.Len(t, filesOldDir, 1)

				p, _ = fspath.NewFileSystemPath("/dir2/dir4/dir5/", nil)
				filesNewDir, _ := fs.FindFiles("file1", p)
				assert.Len(t, filesNewDir, 1)

				// make sure they are not the same file
				assert.NotEqual(t, filesNewDir[0], filesOldDir[0])
			},
		},
		{
			CaseName: "Copy directory, no conflict - relative path",
			SrcPath:  ".",
			DestPath: "../dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/file2", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir1")

				p, _ := fspath.NewFileSystemPath("/dir2", nil)
				files, _ := fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)

				p, _ = fspath.NewFileSystemPath("/", nil)
				files, _ = fs.FindFiles("dir1", p)
				assert.Len(t, files, 2)

				file1 := files[0]
				file2 := files[1]
				assert.NotEqual(t, file1, file2)
				checkCopy(t, fs, file1.AbsolutePath(), file2.AbsolutePath())
			},
		},
		// /**
		// FS structure after initialization:
		// /
		// dir1
		// 	file3
		// 	dir3
		// 		file4
		// 		dir4
		// 	dir5
		// 		file5
		// 		dir6
		// 			file6

		// dir2
		// 	dir1
		// 		file3
		// 		file8
		// 		dir5
		// 			file5
		// 			dir6
		// 				file9

		// expected FS structure after move:
		// /
		// dir1
		// 	file3
		// 	dir3
		// 		file4
		// 		dir4
		// 	dir5
		// 		file5
		// 		dir6
		// 			file6
		// dir2
		// 	dir1
		// 		file3
		// 		file3_uid
		// 		file8
		// 		dir3
		// 			file4
		// 			dir4
		// 		dir5
		// 			file5
		// 			file5_uid
		// 			dir6
		// 				file6
		// 				file9

		// */
		{
			CaseName: "Copy directory, conflict - relative path",
			SrcPath:  ".",
			DestPath: "/dir2/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				workDir, err := fs.Mkdir(p)
				if err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir3/dir4", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/dir6", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file3", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir3/file4", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/file5", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/dir5/dir6/file6", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/file3", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/file8", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/file5", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/dir1/dir5/dir6/file9", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			// TODO: ues list files to check the fs structure
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir2/dir1")
				p, _ := fspath.NewFileSystemPath("/", nil)
				files, _ := fs.FindFiles("dir1", p)
				assert.Len(t, files, 2)

				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				files, _ = fs.FindFiles("dir1", p)
				assert.Len(t, files, 1)
			},
		},
		{
			CaseName: "Copying symlink should copy the orignal file/directory - absolute path",
			SrcPath:  "/dir2/",
			DestPath: "/dir3/dir4/",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/dir1/dir3/", nil)
				if _, err := fs.MkdirAll(p1); err != nil {
					return nil, nil, err
				}
				p2, _ := fspath.NewFileSystemPath("/dir1/file1/", nil)
				workDir, err := fs.CreateRegularFile(p2)
				if err != nil {
					return nil, nil, err
				}
				p3, _ := fspath.NewFileSystemPath("/dir1/", nil)
				p4, _ := fspath.NewFileSystemPath("/dir2/", nil)

				if _, err := fs.CreateSymbolicLink(p3, p4); err != nil {
					return nil, nil, err
				}
				return fs, workDir, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, info file.FileInfo, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, info)

				assert.Equal(t, info.AbsolutePath(), "/dir3/dir4")
				p, _ := fspath.NewFileSystemPath("/dir3/dir4", nil)
				dir, _ := fs.ListFiles(p)
				assert.Len(t, dir, 2)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		srcPath, _ := fspath.NewFileSystemPath(testCase.SrcPath, workingDir)
		destPath, _ := fspath.NewFileSystemPath(testCase.DestPath, workingDir)

		file, err := fs.Copy(srcPath, destPath)
		testCase.Assertions(t, fs, file, err)
	}
}

func checkCopy(t *testing.T, fs *memoryfs.MemoryFileSystem, srcPath string, destPath string) {
	p1, _ := fspath.NewFileSystemPath(srcPath, nil)
	originalFiles, _ := fs.ListFiles(p1)
	p2, _ := fspath.NewFileSystemPath(destPath, nil)
	newFiles, _ := fs.ListFiles(p2)
	assert.Equal(t, len(originalFiles), len(newFiles))

	for i := range originalFiles {
		assert.Equal(t, originalFiles[i].Name(), newFiles[i].Name())
		assert.NotEqual(t, originalFiles[i], newFiles[i])
		checkCopy(t, fs, originalFiles[i].AbsolutePath(), newFiles[i].AbsolutePath())
	}
}
