package memoryfs_test

import (
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fserrors"
	"material/filesystem/filesystem/fspath"
	"material/filesystem/filesystem/memoryfs"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// walk returns an error

func TestWalk(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		Initialize func() (*memoryfs.MemoryFileSystem, file.File, error)
		Assertions func(*testing.T, []string, error)
		FollowLink bool
	}{
		{
			CaseName:   "Walk starting from root",
			Path:       "/",
			FollowLink: true,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir/to/walk/skip", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/to/skip/walk", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/skip", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/to/walk", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/to/walk/file", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}

				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res []string, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Len(t, res, 10)
				sort.Strings(res)
				assert.Equal(t, []string{"/", "dir", "dir1", "dir2", "file", "to", "to", "to", "walk", "walk"}, res)
			},
		},
		{
			CaseName:   "Walk starting from subdir",
			Path:       "/dir1",
			FollowLink: true,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p, _ := fspath.NewFileSystemPath("/dir/to/walk/skip", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir1/to/skip/walk", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/skip", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/to/walk", nil)
				if _, err := fs.MkdirAll(p); err != nil {
					return nil, nil, err
				}
				p, _ = fspath.NewFileSystemPath("/dir2/to/walk/file", nil)
				if _, err := fs.CreateRegularFile(p); err != nil {
					return nil, nil, err
				}

				return fs, nil, nil
			},
			Assertions: func(t *testing.T, res []string, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Len(t, res, 2)
				sort.Strings(res)
				assert.Equal(t, []string{"dir1", "to"}, res)
			},
		},
		{
			CaseName:   "Follow links infinite loop",
			Path:       "/",
			FollowLink: true,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/dir/to/walk/", nil)
				if _, err := fs.MkdirAll(p1); err != nil {
					return nil, nil, err
				}
				p2, _ := fspath.NewFileSystemPath("/dir/link-to-follow", nil)
				p3, _ := fspath.NewFileSystemPath("/dir", nil)
				if _, err := fs.CreateSymbolicLink(p3, p2); err != nil {
					return nil, nil, err
				}

				return fs, nil, nil

			},
			Assertions: func(t *testing.T, res []string, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fserrors.ErrTooManyLinks, err)
			},
		},
		{
			CaseName:   "Follow links",
			Path:       "/",
			FollowLink: true,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/dir/to/walk/", nil)
				if _, err := fs.MkdirAll(p1); err != nil {
					return nil, nil, err
				}
				p2, _ := fspath.NewFileSystemPath("/dir/link-to-follow", nil)
				p3, _ := fspath.NewFileSystemPath("/dir/to", nil)
				if _, err := fs.CreateSymbolicLink(p3, p2); err != nil {
					return nil, nil, err
				}

				return fs, nil, nil

			},
			Assertions: func(t *testing.T, res []string, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Len(t, res, 7)
				sort.Strings(res)
				assert.Equal(t, []string{"/", "dir", "link-to-follow", "to", "to", "walk", "walk"}, res)
			},
		},
		{
			CaseName:   "Follow links",
			Path:       "/",
			FollowLink: false,
			Initialize: func() (*memoryfs.MemoryFileSystem, file.File, error) {
				fs := memoryfs.NewMemoryFileSystem()
				p1, _ := fspath.NewFileSystemPath("/dir/to/walk/", nil)
				if _, err := fs.MkdirAll(p1); err != nil {
					return nil, nil, err
				}
				p2, _ := fspath.NewFileSystemPath("/dir/link-to-follow", nil)
				p3, _ := fspath.NewFileSystemPath("/dir/to", nil)
				if _, err := fs.CreateSymbolicLink(p3, p2); err != nil {
					return nil, nil, err
				}

				return fs, nil, nil

			},
			Assertions: func(t *testing.T, res []string, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, res)
				assert.Len(t, res, 5)
				sort.Strings(res)
				assert.Equal(t, []string{"/", "dir", "link-to-follow", "to", "walk"}, res)
			},
		},
	}
	for _, testCase := range cases {
		fs, workingDir, err := testCase.Initialize()
		if err != nil {
			t.Fatal("error initializing file system")
		}
		p, _ := fspath.NewFileSystemPath(testCase.Path, workingDir)
		res := []string{}
		err = fs.Walk(p, func(f file.File) error {
			res = append(res, f.Info().Name())
			return nil
		}, func(f file.File) bool {
			return f.Info().Name() != "skip"
		}, testCase.FollowLink)
		testCase.Assertions(t, res, err)
	}
}
