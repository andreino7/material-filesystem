package memoryfs

import "sync"

// inMemoryFile implements the FileData interface
type inMemoryFileData struct {
	// TODO: Optimization. make this a slice of slice of fixed size so that
	// expanding size/inserting is more efficient because you just need to copy
	// the pointers to the slices
	data []byte
	sync.RWMutex
}

func (data *inMemoryFileData) Data() []byte {
	return data.data
}

func (data *inMemoryFileData) Size() int {
	return len(data.data)
}

// write writes the content at the given offset.
// If the offset > len(data) fill the gap with 0s.
func (d *inMemoryFileData) write(content []byte, offset int) int {

	if offset == len(d.data) {
		return d.append(content)
	}

	if offset > len(d.data) {
		// fill with 0s
		size := offset - len(d.data)
		empty := make([]byte, size-1)
		d.append(empty)
		return d.append(content)
	}

	return d.insert(content, offset)
}

func (d *inMemoryFileData) read(start int, buff []byte) int {
	if start >= len(d.data) {
		return 0
	}

	end := start + len(buff)
	if end >= len(d.data) {
		end = len(d.data)
	}

	copy(buff, d.data[start:end])
	return end - start
}

// TODO: improve, see optimization above
func (d *inMemoryFileData) insert(content []byte, pos int) int {
	// reslice
	newData := make([]byte, 0, len(d.data)+len(content))
	// insert data
	newData = append(newData, d.data[:pos]...)
	newData = append(newData, content...)
	newData = append(newData, d.data[pos:]...)

	d.data = newData
	return len(content)
}

func (d *inMemoryFileData) append(content []byte) int {
	d.data = append(d.data, content...)
	return len(content)
}
