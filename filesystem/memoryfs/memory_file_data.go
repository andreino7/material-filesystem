package memoryfs

import "sync"

type inMemoryFileData struct {
	// TODO: Optimization. make this a slice of slice so that
	// expanding size/inserting is more efficient becuase you just need to copy
	// the pointers to the slices
	data  []byte
	mutex sync.RWMutex
}

func (data *inMemoryFileData) Data() []byte {
	return data.data
}

// writeAt writes the content at the given offset.
// If the offset > len(data) fill the gap with 0s.
func (d *inMemoryFileData) writeAt(content []byte, offset int) int {
	// offset > len
	// offset < len
	// offset == len
	if offset == len(d.data) {
		return d.append(content)
	}

	if offset > len(content) {
		// fill with 0s
		size := offset - len(content)
		empty := make([]byte, size)
		d.append(empty)
		return d.append(content)
	}

	return d.insert(content, offset)
}

// TODO: check what happens when they are out of bound
func (d *inMemoryFileData) readAt(start int, end int) []byte {
	return d.data[start:end]
}

func (d *inMemoryFileData) insert(content []byte, pos int) int {
	newData := d.data[:pos]
	newData = append(newData, content...)
	newData = append(newData, d.data[pos:]...)
	d.data = newData
	return len(content)
}

func (d *inMemoryFileData) append(content []byte) int {
	d.data = append(d.data, content...)
	return len(content)
}
