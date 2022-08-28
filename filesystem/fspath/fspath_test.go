package fspath_test

import (
	"fmt"
	"material/filesystem/filesystem/file"
	"material/filesystem/filesystem/fspath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFileInfo struct {
	name         string
	absolutePath string
	fType        file.FileType
}

type TestFileData struct {
	data []byte
}

type TestFile struct {
	info TestFileInfo
	data TestFileData
}

func (info TestFileInfo) Name() string {
	return info.name
}

func (info TestFileInfo) FileType() file.FileType {
	return info.fType
}

func (info TestFileInfo) AbsolutePath() string {
	return info.absolutePath
}

func (data TestFileData) Data() []byte {
	return data.data
}

func (f TestFile) Info() file.FileInfo {
	return f.info
}

func (f TestFile) Data() file.FileData {
	return f.data
}

func TestFileSystemPath(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		WorkingDir file.File
		Assertions func(*testing.T, *fspath.FileSystemPath, error)
	}{
		{
			CaseName:   "work dir missing and path relative",
			Path:       "../",
			WorkingDir: nil,
			Assertions: func(t *testing.T, p *fspath.FileSystemPath, err error) {
				assert.Nil(t, p)
				assert.NotNil(t, err)
				assert.Equal(t, "invalid path", err.Error())
			},
		},
		{
			CaseName:   "work dir missing and path absolute",
			Path:       "/a/b/c",
			WorkingDir: nil,
			Assertions: func(t *testing.T, p *fspath.FileSystemPath, err error) {
				assert.NotNil(t, p)
				assert.Nil(t, err)
				assert.Equal(t, "/a/b/c", p.AbsolutePath())
			},
		},
		{
			CaseName: "work dir present and path relative",
			Path:     "../..",
			WorkingDir: TestFile{
				info: TestFileInfo{absolutePath: "/a/b/c"},
			},
			Assertions: func(t *testing.T, p *fspath.FileSystemPath, err error) {
				assert.NotNil(t, p)
				assert.Nil(t, err)
				assert.Equal(t, "/a", p.AbsolutePath())
			},
		},
		{
			CaseName: "path with spaces",
			Path:     "../../path with spaces",
			WorkingDir: TestFile{
				info: TestFileInfo{absolutePath: "/a/b/c"},
			},
			Assertions: func(t *testing.T, p *fspath.FileSystemPath, err error) {
				assert.NotNil(t, p)
				assert.Nil(t, err)
				assert.Equal(t, "/a/path with spaces", p.AbsolutePath())
				assert.Equal(t, "path with spaces", p.Base())
				assert.Equal(t, "../..", p.Dir())
			},
		},
	}
	for _, testCase := range cases {
		fmt.Println(testCase.CaseName)
		p, err := fspath.NewFileSystemPath(testCase.Path, testCase.WorkingDir)
		testCase.Assertions(t, p, err)
	}
}
