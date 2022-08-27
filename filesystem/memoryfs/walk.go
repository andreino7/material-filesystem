package memoryfs

import "material/filesystem/filesystem/user"

type walkFn func(string, *inMemoryFile) error

func (fs *MemoryFileSystem) walk(rootFile *inMemoryFile, user user.User, walkFn walkFn) error {
	for fileName, file := range rootFile.fileMap {
		// skip special keys to avoid infinite cycle
		if fileName == ".." || fileName == "." || fileName == "/" {
			continue
		}

		if err := checkReadPermissions(file, user); err != nil {
			return err
		}

		if err := walkFn(fileName, file); err != nil {
			return err
		}
	}

	return nil
}
