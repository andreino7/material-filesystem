package memoryfs_test

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				fs.AppendToFile(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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
				fs.AppendToFile(src, []byte("Hello world!"), nil)
				fs.Remove(src, nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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
				fs.AppendToFile(src, []byte("Hello world!"), nil)
				fs.Remove(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				data, _ := fs.ReadFile(src, nil)
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
				fs.AppendToFile(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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
				fs.AppendToFile(fspath.NewFileSystemPath("/file1"), []byte("Hello world!"), nil)
				fs.Remove(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				data, _ := fs.ReadFile(src, nil)
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
				fs.AppendToFile(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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
				fs.AppendToFile(src, []byte("Hello world!"), nil)
				fs.Remove(src, nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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
				fs.AppendToFile(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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
				fs.AppendToFile(fspath.NewFileSystemPath("/file1"), []byte("Hello world!"), nil)
				fs.Remove(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
				data, _ := fs.ReadFile(src, nil)
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

				fs.AppendToFile(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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

				fs.AppendToFile(fspath.NewFileSystemPath(res.AbsolutePath()), []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(src, nil)
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

				fs.AppendToFile(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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

				fs.AppendToFile(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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

				fs.AppendToFile(fspath.NewFileSystemPath(res.AbsolutePath()), []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(src, nil)
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

				fs.AppendToFile(src, []byte("Hello world!"), nil)
				data, _ := fs.ReadFile(fspath.NewFileSystemPath(res.AbsolutePath()), nil)
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
