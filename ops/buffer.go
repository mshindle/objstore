package ops

import "sync"

// A WriteBuffer provides a in memory buffer supporting the io.WriterAt and io.Writer interface
// Can be used with the s3manager.Downloader to download content to a buffer
// in memory. Safe to use concurrently.
type WriteBuffer struct {
	buf []byte
	m   sync.Mutex

	// GrowthCoeff defines the growth rate of the internal buffer. By
	// default, the growth rate is 1, where expanding the internal
	// buffer will allocate only enough capacity to fit the new expected
	// length.
	GrowthCoeff float64
}

// NewWriteBuffer creates a WriteAtBuffer with an internal buffer
// provided by buf.
func NewWriteBuffer(buf []byte) *WriteBuffer {
	return &WriteBuffer{buf: buf, GrowthCoeff: 1.4}
}

// WriteAt writes a slice of bytes to a buffer starting at the position provided
// The number of bytes written will be returned, or error. Can overwrite previous
// written slices if the write ats overlap.
func (b *WriteBuffer) WriteAt(p []byte, pos int64) (n int, err error) {
	pLen := len(p)
	expLen := pos + int64(pLen)
	b.m.Lock()
	defer b.m.Unlock()
	if int64(len(b.buf)) < expLen {
		if int64(cap(b.buf)) < expLen {
			if b.GrowthCoeff < 1 {
				b.GrowthCoeff = 1
			}
			newBuf := make([]byte, expLen, int64(b.GrowthCoeff*float64(expLen)))
			copy(newBuf, b.buf)
			b.buf = newBuf
		}
		b.buf = b.buf[:expLen]
	}
	copy(b.buf[pos:], p)
	return pLen, nil
}

func (b *WriteBuffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// Size returns the current size of the buffer in bytes
func (b *WriteBuffer) Size() int {
	b.m.Lock()
	defer b.m.Unlock()
	return len(b.buf)
}

// Bytes returns a slice of bytes written to the buffer.
func (b *WriteBuffer) Bytes() []byte {
	b.m.Lock()
	defer b.m.Unlock()
	return b.buf[:len(b.buf):len(b.buf)]
}

