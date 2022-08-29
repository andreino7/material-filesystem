package memoryfs

type fileDescriptor struct {
	data   *inMemoryFileData
	offset int
}

func (fd *fileDescriptor) Read(buff []byte) (int, error) {
	nRead := fd.data.read(fd.offset, buff)
	fd.offset += nRead
	return nRead, nil
}

func (fd *fileDescriptor) ReadAt(buff []byte, offset int) (int, error) {
	nRead := fd.data.read(offset, buff)
	return nRead, nil
}
