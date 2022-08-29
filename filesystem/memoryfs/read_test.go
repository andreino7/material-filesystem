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
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
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
			CaseName: "Read from existing file - relative path",
			Path:     "file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
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
			CaseName: "Read from symbolic link - relative path",
			Path:     "file1-link",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p1); err != nil {
					return nil, nil, err
				}
				if err := fs.AppendAll(p1, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				p2, _ := fspath.NewFileSystemPath("/file1-link", nil)

				if _, err := fs.CreateSymbolicLink(p1, p2); err != nil {
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
		data, err := fs.ReadAll(path)
		testCase.Assertions(t, data, err)
	}
}

func TestReadAt(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Offset     int
		BuffSize   int
		Assertions func(*testing.T, int, []byte, error)
	}{
		{
			CaseName: "Read from non empty file - absolute path",
			Path:     "/file1",
			Offset:   3,
			BuffSize: 5,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, buff []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("lo wo"), buff)
			},
		},
		{
			CaseName: "Read empty file should return empty bytes - absolute path",
			Path:     "/file1",
			Offset:   3,
			BuffSize: 5,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, buff []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, nBytes, 0)
				assert.Equal(t, make([]byte, 5), buff)
			},
		},
		{
			CaseName: "Read from non empty file, both out of bound - absolute path",
			Path:     "/file1",
			Offset:   15,
			BuffSize: 5,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, buff []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, nBytes, 0)
				assert.Equal(t, make([]byte, 5), buff)
			},
		},
		{
			CaseName: "Read from non empty file, right out of bound - absolute path",
			Path:     "/file1",
			Offset:   9,
			BuffSize: 20,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, buff []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, nBytes, 3)
				expected := make([]byte, 20)
				expected[0] = 'l'
				expected[1] = 'd'
				expected[2] = '!'
				assert.Equal(t, expected, buff)
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

		buff := make([]byte, testCase.BuffSize)
		nBytes, err := fs.ReadAt(fd, buff, testCase.Offset)
		testCase.Assertions(t, nBytes, buff, err)
	}
}

func TestReadAtClosedFile(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, int, error)
	}{
		{
			CaseName: "Read closed file - absolute path",
			Path:     "/file1",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fserrors.ErrNotOpen, err)
				assert.Equal(t, 0, nBytes)
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

		nBytes, err := fs.ReadAt(fd, make([]byte, 15), 5)
		testCase.Assertions(t, nBytes, err)
	}
}

func TestReadAtMovedOrRemovedFile(t *testing.T) {
	cases := []struct {
		CaseName    string
		Path        string
		FsOperation func(*memoryfs.MemoryFileSystem) error
		Initialize  func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions  func(*testing.T, int, []byte, error)
	}{
		{
			CaseName: "Read from renamed file after opening it should still work - absolute path",
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
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, content []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("Hello"), content)
				assert.Equal(t, nBytes, 5)
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
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, content []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("Hello"), content)
				assert.Equal(t, nBytes, 5)
			},
		},
		{
			CaseName: "Read from removed file after opening it should still work - absolute path",
			Path:     "/file1",
			FsOperation: func(fs *memoryfs.MemoryFileSystem) error {
				p1, _ := fspath.NewFileSystemPath("/file1", nil)
				_, err := fs.Remove(p1)
				return err
			},
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, content []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("Hello"), content)
				assert.Equal(t, nBytes, 5)
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

		buff := make([]byte, 5)
		n, err := fs.ReadAt(fd, buff, 0)
		testCase.Assertions(t, n, buff, err)
	}
}

func TestRead(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		BuffSize   int
		Assertions func(*testing.T, int, []byte, error)
	}{
		{
			CaseName: "Read from non empty file - absolute path",
			Path:     "/file1",
			BuffSize: 5,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, buff []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, []byte("Hello"), buff)
			},
		},
		{
			CaseName: "Read empty file should return empty bytes - absolute path",
			Path:     "/file1",
			BuffSize: 5,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, buff []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, nBytes, 0)
				assert.Equal(t, make([]byte, 5), buff)
			},
		},
		{
			CaseName: "Read from non empty file, right out of bound - absolute path",
			Path:     "/file1",
			BuffSize: 20,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/file1", nil)
				if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, nBytes int, buff []byte, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 12, nBytes)
				expected := make([]byte, 20)
				copy(expected, []byte("Hello world!")[0:12])
				assert.Equal(t, expected, buff)
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

		buff := make([]byte, testCase.BuffSize)
		nBytes, err := fs.Read(fd, buff)
		testCase.Assertions(t, nBytes, buff, err)
	}
}

func TestReadInChuncks(t *testing.T) {
	fs := memoryfs.NewMemoryFileSystem()
	p, _ := fspath.NewFileSystemPath("/file1", nil)
	if err := fs.AppendAll(p, []byte("Hello world!")); err != nil {
		t.Fatal("error initializing file system")
	}
	fd, err := fs.Open(p)
	if err != nil {
		t.Fatal("error opening file")
	}

	// First chunk
	buff := make([]byte, 5)
	nBytes, err := fs.Read(fd, buff)
	assert.Nil(t, err)
	assert.Equal(t, 5, nBytes)
	assert.Equal(t, []byte("Hello"), buff)

	// Second chunk
	buff = make([]byte, 5)
	nBytes, err = fs.Read(fd, buff)
	assert.Nil(t, err)
	assert.Equal(t, 5, nBytes)
	assert.Equal(t, []byte(" worl"), buff)

	// Third chunk
	buff = make([]byte, 5)
	nBytes, err = fs.Read(fd, buff)
	assert.Nil(t, err)
	epected := []byte("d!")
	epected = append(epected, 0, 0, 0)
	assert.Equal(t, 2, nBytes)
	assert.Equal(t, epected, buff)
}
