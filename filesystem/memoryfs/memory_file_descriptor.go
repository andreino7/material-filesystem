package memoryfs

type fileDescriptor struct {
	data   *inMemoryFileData
	offset int
}

// Read reads at most len(buff) bytes from the file
// starting at the current offset
func (fd *fileDescriptor) Read(buff []byte) (int, error) {
	nRead := fd.data.read(fd.offset, buff)
	fd.offset += nRead
	return nRead, nil
}

// ReadAt reads at most len(buff) bytes from the file
// starting at the given offset.
func (fd *fileDescriptor) ReadAt(buff []byte, offset int) (int, error) {
	nRead := fd.data.read(offset, buff)
	return nRead, nil
}

// Write writes len(buff) bytes to the file
// starting at the current offset
func (fd *fileDescriptor) Write(buff []byte) (int, error) {
	nWrite := fd.data.write(buff, fd.offset)
	fd.offset += nWrite
	return nWrite, nil
}

// WriteAt writes len(buff) bytes to the file
// starting at the given offset
func (fd *fileDescriptor) WriteAt(buff []byte, offset int) (int, error) {
	nWrite := fd.data.write(buff, offset)
	return nWrite, nil
}
