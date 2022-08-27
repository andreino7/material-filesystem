package memoryfs_test

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendAll(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, []byte, error)
	}{
		{
			CaseName: "Append to existing file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
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
			CaseName: "Append to new file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Append to new file and create intermediate directories - absolute path",
			Path:     "/dir1/dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Append to directory - absolute path",
			Path:     "/dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
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
			CaseName: "Append to existing file - relative path",
			Path:     "file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
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
			CaseName: "Append to new file - relative path",
			Path:     "./file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Append to new file and create intermediate directories - relative path",
			Path:     "dir1/dir2/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				return fs, fs.DefaultWorkingDirectory(), nil
			},
			Assertions: func(t *testing.T, data []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, data, []byte("Hello world!"))
			},
		},
		{
			CaseName: "Append to directory - relative path",
			Path:     "dir1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir1", nil)
				if _, err := fs.Mkdir(p); err != nil {
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
			CaseName: "Append to symlink - absolute path",
			Path:     "/file1-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/file1", nil)

				if _, err := fs.CreateRegularFile(p1); err != nil {
					return nil, nil, err
				}

				p2, _ := fspath.NewFileSystemPath("/file1-link", nil)
				if _, err := fs.CreateSymbolicLink(p1, p2); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
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
		err = fs.AppendAll(path, []byte("Hello world!"))
		if err != nil {
			testCase.Assertions(t, nil, err)
		} else {
			data, _ := fs.ReadAll(path)
			testCase.Assertions(t, data, err)
		}
	}
}

func TestWriteAt(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Position   int
		Text       []byte
		Assertions func(*testing.T, *memoryfs.MemoryFileSystem, *fspath.FileSystemPath, int, error)
	}{
		{
			CaseName: "Write to empty open file at pos 0 - absolute path",
			Path:     "/file1",
			Position: 0,
			Text:     []byte("Hello world!"),
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, p *fspath.FileSystemPath, b int, err error) {
				assert.Nil(t, err)
				assert.Equal(t, b, 12)

				data, _ := fs.ReadAll(p)
				assert.Equal(t, []byte("Hello world!"), data)
			},
		},
		{
			CaseName: "Write to not empty open file at pos 0 - absolute path",
			Path:     "/file1",
			Position: 0,
			Text:     []byte("Hello universe! "),
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, p *fspath.FileSystemPath, b int, err error) {
				assert.Nil(t, err)
				assert.Equal(t, b, 16)

				data, _ := fs.ReadAll(p)
				assert.Equal(t, []byte("Hello universe! Hello world!"), data)
			},
		},
		{
			CaseName: "Write to not empty open file at random pos - absolute path",
			Path:     "/file1",
			Position: 6,
			Text:     []byte("universe! "),
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, p *fspath.FileSystemPath, b int, err error) {
				assert.Nil(t, err)
				assert.Equal(t, b, 10)

				data, _ := fs.ReadAll(p)
				assert.Equal(t, []byte("Hello universe! world!"), data)
			},
		},
		{
			CaseName: "Write file to pos > end - absolute path",
			Path:     "/file1",
			Position: 20,
			Text:     []byte("Hello universe!"),
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, p *fspath.FileSystemPath, b int, err error) {
				assert.Nil(t, err)
				assert.Equal(t, b, 15)

				data, _ := fs.ReadAll(p)
				epected := []byte("Hello world!")
				epected = append(epected, 0, 0, 0, 0, 0, 0, 0)
				epected = append(epected, []byte("Hello universe!")...)
				assert.Equal(t, epected, data)
			},
		},
		{
			CaseName: "Write to empty open file at pos > 0 - absolute path",
			Path:     "/file1",
			Position: 10,
			Text:     []byte("Hello world!"),
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, p *fspath.FileSystemPath, b int, err error) {
				assert.Nil(t, err)
				assert.Equal(t, b, 12)

				data, _ := fs.ReadAll(p)
				epected := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}
				epected = append(epected, []byte("Hello world!")...)
				assert.Equal(t, epected, data)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		path, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		fd, err := fs.Open(path)
		if err != nil {
			t.Fatal("error opening file")
		}

		b, err := fs.WriteAt(fd, testCase.Text, testCase.Position)
		testCase.Assertions(t, fs, path, b, err)
	}
}

func TestWriteAtClosedFile(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, int, error)
	}{
		{
			CaseName: "Write to closed file should fail - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, b int, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fserrors.ErrNotOpen, err)
				assert.Equal(t, 0, b)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		path, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		fd, err := fs.Open(path)
		if err != nil {
			t.Fatal("error opening file")
		}
		fs.Close(fd)

		b, err := fs.WriteAt(fd, []byte("Hello world!"), 0)
		testCase.Assertions(t, b, err)
	}
}

func TestWriteMovedOrRemovedFile(t *testing.T) {
	cases := []struct {
		CaseName    string
		Path        string
		FsOperation func(*memoryfs.MemoryFileSystem) error
		Initialize  func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions  func(*testing.T, *memoryfs.MemoryFileSystem, int, error)
	}{
		{
			CaseName: "Write to renamed file after opening it should still work - absolute path",
			Path:     "/file1",
			FsOperation: func(fs *memoryfs.MemoryFileSystem) error {
				p1, _ := fspath.NewFileSystemPath("/file1", nil)
				p2, _ := fspath.NewFileSystemPath("/file1-new-name", nil)
				_, err := fs.Move(p1, p2)
				return err
			},
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, b int, err error) {
				assert.Nil(t, err)
				assert.Equal(t, b, 12)

				p, _ := fspath.NewFileSystemPath("/file1-new-name", nil)
				data, _ := fs.ReadAll(p)
				assert.Equal(t, []byte("Hello world!"), data)
			},
		},
		{
			CaseName: "Write to moved file after opening it should still work - absolute path",
			Path:     "/dir1/file1",
			FsOperation: func(fs *memoryfs.MemoryFileSystem) error {
				p1, _ := fspath.NewFileSystemPath("/dir1/file1", nil)
				p2, _ := fspath.NewFileSystemPath("/dir2/file1", nil)
				_, err := fs.Move(p1, p2)
				return err
			},
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
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, b int, err error) {
				assert.Nil(t, err)
				assert.Equal(t, b, 12)

				p, _ := fspath.NewFileSystemPath("/dir2/file1", nil)
				data, _ := fs.ReadAll(p)
				assert.Equal(t, []byte("Hello world!"), data)
			},
		},
		{
			CaseName: "Write to removed file after opening it should still work - absolute path",
			Path:     "/file1",
			FsOperation: func(fs *memoryfs.MemoryFileSystem) error {
				p1, _ := fspath.NewFileSystemPath("/file1", nil)
				_, err := fs.Remove(p1)
				return err
			},
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, fs *memoryfs.MemoryFileSystem, b int, err error) {
				assert.Nil(t, err)
				assert.Equal(t, b, 12)

				p, _ := fspath.NewFileSystemPath("/file1", nil)
				_, err = fs.ReadAll(p)
				assert.Equal(t, fserrors.ErrNotExist, err)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		path, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		fd, err := fs.Open(path)
		if err != nil {
			t.Fatal("error opening file")
		}
		if err := testCase.FsOperation(fs); err != nil {
			t.Fatal("error running fs operation")
		}

		b, err := fs.WriteAt(fd, []byte("Hello world!"), 0)
		testCase.Assertions(t, fs, b, err)
	}
}
