package fspath_test

import (
	"material/filesystem/filesystem/fspath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFsPath(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		WorkingDir string
		Assertions func(t *testing.T, caseName string, pathInfo *fspath.FileSystemPath)
	}{
		{
			CaseName:   "path is ../test",
			Path:       "../",
			WorkingDir: "/root/first",
			Assertions: func(t *testing.T, caseName string, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.AbsolutePath(), "/root")
				assert.Equal(t, pathInfo.Dir(), "/")
				assert.Equal(t, pathInfo.Base(), "root")
			},
		},
		{
			CaseName:   "path is ./test",
			Path:       "./test",
			WorkingDir: "/root/first",
			Assertions: func(t *testing.T, caseName string, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.AbsolutePath(), "/root/first/test")
				assert.Equal(t, pathInfo.Dir(), "/root/first")
				assert.Equal(t, pathInfo.Base(), "test")
			},
		},
		{
			CaseName:   "path is /root2",
			Path:       "/root2",
			WorkingDir: "/root/first",
			Assertions: func(t *testing.T, caseName string, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.AbsolutePath(), "/root2")
				assert.Equal(t, pathInfo.Dir(), "/")
				assert.Equal(t, pathInfo.Base(), "root2")
			},
		},
		{
			CaseName:   "path is ../second/third/../../fourth",
			Path:       "../second/third/../../fourth/",
			WorkingDir: "/root/first/",
			Assertions: func(t *testing.T, caseName string, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.AbsolutePath(), "/root/fourth")
				assert.Equal(t, pathInfo.Dir(), "/root")
				assert.Equal(t, pathInfo.Base(), "fourth")
			},
		},
		{
			CaseName:   "path is ../../../../..",
			Path:       "../../../../..",
			WorkingDir: "/root/first",
			Assertions: func(t *testing.T, caseName string, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.AbsolutePath(), "/")
				assert.Equal(t, pathInfo.Dir(), "/")
				assert.Equal(t, pathInfo.Base(), "/")
			},
		},
	}
	for _, testCase := range cases {
		pathInfo := fspath.NewFileSystemPath(testCase.Path, testCase.WorkingDir)
		testCase.Assertions(t, testCase.CaseName, pathInfo)
	}
}
