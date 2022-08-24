package memoryfs

type walkFn func(string, *inMemoryFile) error

func (fs *MemoryFileSystem) walk(rootFile *inMemoryFile, walkFn walkFn) error {
	for fileName, file := range rootFile.fileMap {
		// skip special keys to avoid infinite cycle
		if fileName == ".." || fileName == "." || fileName == "/" {
			continue
		}

		if err := walkFn(fileName, file); err != nil {
			return err
		}
	}

	return nil
}
