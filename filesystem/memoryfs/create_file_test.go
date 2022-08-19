package memoryfs_test

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: add tests for working dir not nil
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
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir1")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Create directory in subdir - absolute path, work dir nil",
			Path:     "/dir1/dir2",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir2")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		{
			CaseName: "Name conflict - absolute path, work dir nil",
			Path:     "/dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1/dir2/dir3", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "file already exists")
				assert.Nil(t, file)
			},
		},
		{
			CaseName: "Parent directory does not exist - absolute path, work dir nil",
			Path:     "/dir1/dir2/dir3",
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir1", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir2", ""), nil); err != nil {
					return nil, nil, err
				}
				if _, err := fs.Mkdir(fspath.NewFileSystemPath("/dir4", ""), nil); err != nil {
					return nil, nil, err
				}
				return fs, nil, nil
			},
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "no such file or directory")
				assert.Nil(t, file)
			},
		},
		// {
		// 	CaseName: "Parent is not a directory - absolute path, work dir nil",
		// 	Path:     "/file/dir1",
		// 	Initialize: func() (*memapfs.MemMapFs, error) {
		// 		fs := memapfs.NewMemMapFs()
		// 		if _, err := fs.Create("/first_file"); err != nil {
		// 			return nil, err
		// 		}
		// 		return fs, nil
		// 	},
		// 	Assertions: func(t *testing.T, dirPath string, file *filesystem.File, err error) {
		// 		assert.NotNil(t, err)
		// 		assert.Equal(t, err.Error(), "file is not a directory")
		// 		assert.Nil(t, file)
		// 	},
		// },
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		dir, err := fs.Mkdir(fspath.NewFileSystemPath(testCase.Path, ""), workingDir)
		testCase.Assertions(t, dir, err)
	}
}

// func TestCreate(t *testing.T) {
// 	cases := []struct {
// 		CaseName   string
// 		Paths      []string
// 		Initialize func() (*memapfs.MemMapFs, error)
// 		Assertions func(t *testing.T, dirName string, file *filesystem.File, err error)
// 	}{
// 		{

// 			CaseName: "Create multiple files - no name conflicts",
// 			Paths:    []string{"/first_file", "/second_file", "/first_directory/first_file", "/first_directory/second_directory/first_file", "/second_directory/some_file"},
// 			Initialize: func() (*memapfs.MemMapFs, error) {
// 				fs := memapfs.NewMemMapFs()
// 				if _, err := fs.Mkdir("/first_directory"); err != nil {
// 					return nil, err
// 				}
// 				if _, err := fs.Mkdir("/first_directory/second_directory"); err != nil {
// 					return nil, err
// 				}
// 				if _, err := fs.Mkdir("/second_directory"); err != nil {
// 					return nil, err
// 				}
// 				return fs, nil
// 			},
// 			Assertions: func(t *testing.T, dirPath string, file *filesystem.File, err error) {
// 				assert.Nil(t, err)
// 				assert.NotNil(t, file)
// 				assert.Equal(t, file.Path(), dirPath)
// 				assert.False(t, file.IsDirectory())
// 			},
// 		},
// 		{
// 			CaseName: "Name conflict",
// 			Paths:    []string{"/first_file", "/first_directory/some_file"},
// 			Initialize: func() (*memapfs.MemMapFs, error) {
// 				fs := memapfs.NewMemMapFs()
// 				if _, err := fs.Create("/first_file"); err != nil {
// 					return nil, err
// 				}
// 				if _, err := fs.Mkdir("/first_directory"); err != nil {
// 					return nil, err
// 				}
// 				if _, err := fs.Create("/first_directory/some_file"); err != nil {
// 					return nil, err
// 				}
// 				return fs, nil
// 			},
// 			Assertions: func(t *testing.T, dirPath string, file *filesystem.File, err error) {
// 				assert.NotNil(t, err)
// 				assert.Equal(t, err.Error(), "file already exists")
// 				assert.Nil(t, file)
// 			},
// 		},
// 		{
// 			CaseName: "Parent directory does not exists",
// 			Paths:    []string{"/second_directory/some_file"},
// 			Initialize: func() (*memapfs.MemMapFs, error) {
// 				fs := memapfs.NewMemMapFs()
// 				if _, err := fs.Mkdir("/first_directory"); err != nil {
// 					return nil, err
// 				}
// 				return fs, nil
// 			},
// 			Assertions: func(t *testing.T, dirPath string, file *filesystem.File, err error) {
// 				assert.NotNil(t, err)
// 				assert.Equal(t, err.Error(), "file not found")
// 				assert.Nil(t, file)
// 			},
// 		},
// 	}
// 	for _, testCase := range cases {
// 		fs, err := testCase.Initialize()
// 		if err != nil {
// 			t.Fatal("error initializing file system")
// 		}
// 		for _, absPath := range testCase.Paths {
// 			file, err := fs.Create(absPath)
// 			testCase.Assertions(t, absPath, file, err)
// 		}
// 	}
// }

// TODO: check that every directory in path exists
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
			Assertions: func(t *testing.T, file file.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, file)
				assert.Equal(t, file.Info().Name(), "dir3")
				assert.True(t, file.Info().IsDirectory())
			},
		},
		// TODO: uncomment when fixing name conflict for mkdirall
		// {
		// 	CaseName: "Name conflict",
		// 	Paths:    []string{"/first_directory", "/first_directory/third_directory"},
		// 	Initialize: func() (*memapfs.MemMapFs, error) {
		// 		fs := memapfs.NewMemMapFs()
		// 		if _, err := fs.Mkdir("/first_directory"); err != nil {
		// 			return nil, err
		// 		}
		// 		if _, err := fs.Mkdir("/first_directory/third_directory"); err != nil {
		// 			return nil, err
		// 		}
		// 		return fs, nil
		// 	},
		// 	Assertions: func(t *testing.T, dirPath string, file *filesystem.File, err error) {
		// 		assert.NotNil(t, err)
		// 		assert.Equal(t, err.Error(), "file already exists")
		// 		assert.Nil(t, file)
		// 	},
		// },
		// {
		// 	CaseName: "Parent is not a directory",
		// 	Paths:    []string{"/first_file/dir1/dir2"},
		// 	Initialize: func() (*memapfs.MemMapFs, error) {
		// 		fs := memapfs.NewMemMapFs()
		// 		if _, err := fs.Create("/first_file"); err != nil {
		// 			return nil, err
		// 		}
		// 		return fs, nil
		// 	},
		// 	Assertions: func(t *testing.T, dirPath string, file *filesystem.File, err error) {
		// 		assert.NotNil(t, err)
		// 		assert.Equal(t, err.Error(), "file is not a directory")
		// 		assert.Nil(t, file)
		// 	},
		// },
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		dir, err := fs.MkdirAll(fspath.NewFileSystemPath(testCase.Path, ""), workingDir)
		testCase.Assertions(t, dir, err)
	}
}
