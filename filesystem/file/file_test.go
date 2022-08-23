package file_test

type TestFileInfo struct {
	name         string
	absolutePath string
	isDirectory  bool
}

type TestFileData struct {
	data []byte
}

type TestFile struct {
	info TestFileInfo
	data TestFileData
}

func (info *TestFileInfo) Name() string {
	return info.name
}

func (info *TestFileInfo) AbsolutePath() string {
	return info.absolutePath
}

func (info *TestFileInfo) IsDirectory() bool {
	return info.isDirectory
}

func (data *TestFileData) Data() []byte {
	return data.data
}

func (f *TestFile) Info() TestFileInfo {
	return f.info
}

func (f *TestFile) Data() TestFileData {
	return f.data
}
