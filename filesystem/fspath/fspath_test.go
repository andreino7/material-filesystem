package fspath_test

import (
	"fmt"
	"material/filesystem/filesystem/fspath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: redo these tests
func TestFsPath(t *testing.T) {
	cases := []struct {
		CaseName   string
		Path       string
		WorkingDir string
		Assertions func(t *testing.T, pathInfo *fspath.FileSystemPath)
	}{
		{
			CaseName: "path is ../",
			Path:     "../",
			Assertions: func(t *testing.T, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.Dir(), ".")
				assert.Equal(t, pathInfo.Base(), "..")
			},
		},
		{
			CaseName: "path is ./test",
			Path:     "./test",
			Assertions: func(t *testing.T, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.Dir(), ".")
				assert.Equal(t, pathInfo.Base(), "test")
			},
		},
		{
			CaseName: "path is /root2",
			Path:     "/root2",
			Assertions: func(t *testing.T, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.Dir(), "/")
				assert.Equal(t, pathInfo.Base(), "root2")
			},
		},
		{
			CaseName: "path is ../second/third/../../fourth",
			Path:     "../second/third/../../fourth/",
			Assertions: func(t *testing.T, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.Dir(), "..")
				assert.Equal(t, pathInfo.Base(), "fourth")
			},
		},
		{
			CaseName: "path is ../../.././././../..",
			Path:     "../../.././././../..",
			Assertions: func(t *testing.T, pathInfo *fspath.FileSystemPath) {
				assert.Equal(t, pathInfo.Dir(), "../../../..")
				assert.Equal(t, pathInfo.Base(), "..")
			},
		},
	}
	for _, testCase := range cases {
		fmt.Println(testCase.CaseName)
		// pathInfo := fspath.NewFileSystemPath(testCase.Path)
		// testCase.Assertions(t, pathInfo)
	}
}
