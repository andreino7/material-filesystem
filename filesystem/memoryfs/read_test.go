package memoryfs_test

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadAll(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, []byte, error)
	}{
		{
			CaseName: "Read from existing file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}

				return fs, nil, nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Read from missing file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
		{
			CaseName: "Read directory - absolute path",
			Path:     "/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p, rootUser); err != nil {
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
			CaseName: "Read from existing file - relative path",
			Path:     "file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}

				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Read from missing file - relative path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, fserrors.ErrNotExist)
			},
		},
		{
			CaseName: "Read directory - relative path",
			Path:     "./dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p, rootUser); err != nil {
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
			CaseName: "Read from symbolic link - relative path",
			Path:     "file1-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p1, rootUser); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendAll(p1, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}
				p2, _ := fspath.NewFileSystemPath("/file1-link", nil)

				if _, err := fs.CreateSymbolicLink(p1, p2, rootUser); err != nil {
					return nil, nil, err
				}

				return fs, fs.DefaultWorkingDirectory(), nil
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
		data, err := fs.ReadAll(path, rootUser)
		testCase.Assertions(t, data, err)
	}
}

// read closed file
func TestReadAt(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Start      int
		End        int
		Assertions func(*testing.T, []byte, error)
	}{
		{
			CaseName: "Read from non empty file - absolute path",
			Path:     "/file1",
			Start:    3,
			End:      7,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, content []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("lo wo"), content)
			},
		},
		{
			CaseName: "Read empty file should return empty bytes - absolute path",
			Path:     "/file1",
			Start:    3,
			End:      7,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p, rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, content []byte, err error) {
				assert.Nil(t, err)
				assert.Len(t, content, 0)
			},
		},
		{
			CaseName: "Read from non empty file, both out of bound - absolute path",
			Path:     "/file1",
			Start:    15,
			End:      18,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, content []byte, err error) {
				assert.Nil(t, err)
				assert.Len(t, content, 0)
			},
		},
		{
			CaseName: "Read from non empty file, right out of bound - absolute path",
			Path:     "/file1",
			Start:    9,
			End:      18,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, content []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("ld!"), content)
			},
		},
		{
			CaseName: "End pos < start pos - absolute path",
			Path:     "/file1",
			Start:    5,
			End:      3,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, content []byte, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fserrors.ErrInvalid, err)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		path, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		fd, err := fs.Open(path, rootUser)
		if err != nil {
			t.Fatal("error opening file")
		}

		content, err := fs.ReadAt(fd, testCase.Start, testCase.End, rootUser)
		testCase.Assertions(t, content, err)
	}
}

func TestReadAtClosedFile(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, []byte, error)
	}{
		{
			CaseName: "Read closed file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, content []byte, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fserrors.ErrNotOpen, err)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		path, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		fd, err := fs.Open(path, rootUser)
		if err != nil {
			t.Fatal("error opening file")
		}
		fs.Close(fd, rootUser)

		content, err := fs.ReadAt(fd, 0, 5, rootUser)
		testCase.Assertions(t, content, err)
	}
}

func TestReadMovedOrRemovedFile(t *testing.T) {
	cases := []struct {
		CaseName    string
		Path        string
		FsOperation func(*memoryfs.MemoryFileSystem) error
		Initialize  func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions  func(*testing.T, []byte, error)
	}{
		{
			CaseName: "Read from renamed file after opening it should still work - absolute path",
			Path:     "/file1",
			FsOperation: func(fs *memoryfs.MemoryFileSystem) error {
				p1, _ := fspath.NewFileSystemPath("/file1", nil)
				p2, _ := fspath.NewFileSystemPath("/file1-new-name", nil)
				_, err := fs.Move(p1, p2, rootUser)
				return err
			},
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, content []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("Hello world!"), content)
			},
		},
		{
			CaseName: "Write to moved file after opening it should still work - absolute path",
			Path:     "/dir1/file1",
			FsOperation: func(fs *memoryfs.MemoryFileSystem) error {
				p1, _ := fspath.NewFileSystemPath("/dir1/file1", nil)
				p2, _ := fspath.NewFileSystemPath("/dir2/file1", nil)
				_, err := fs.Move(p1, p2, rootUser)
				return err
			},
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2", nil)
				if _, err := fs.Mkdir(p, rootUser); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, content []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("Hello world!"), content)
			},
		},
		{
			CaseName: "Read from removed file after opening it should still work - absolute path",
			Path:     "/file1",
			FsOperation: func(fs *memoryfs.MemoryFileSystem) error {
				p1, _ := fspath.NewFileSystemPath("/file1", nil)
				_, err := fs.Remove(p1, rootUser)
				return err
			},
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!"), rootUser); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, content []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("Hello world!"), content)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		path, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		fd, err := fs.Open(path, rootUser)
		if err != nil {
			t.Fatal("error opening file")
		}
		if err := testCase.FsOperation(fs); err != nil {
			t.Fatal("error running fs operation")
		}

		content, err := fs.ReadAt(fd, 0, 30, rootUser)
		testCase.Assertions(t, content, err)
	}
}
