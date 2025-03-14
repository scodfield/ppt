package util

import (
	"sync"
)

const (
	RingBufferDefaultSize = 1024
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type RingBuffer struct {
	buffer      []byte
	bufferSize  int
	header      int
	headerMark  int
	trailer     int
	trailerMark int
	isEmpty     bool
	mu          sync.Mutex
	logHook     Logger
}

func NewRingBuffer(bufferSize int, logHook Logger) *RingBuffer {
	return &RingBuffer{
		buffer:      make([]byte, bufferSize),
		bufferSize:  bufferSize,
		header:      0,
		headerMark:  0,
		trailer:     0,
		trailerMark: 0,
		isEmpty:     true,
		logHook:     logHook,
	}
}

func (rb *RingBuffer) length() int {
	offset := rb.header - rb.trailer
	if offset > 0 {
		return rb.bufferSize - offset
	} else if offset < 0 {
		return -offset
	}
	if rb.isEmpty {
		return 0
	}
	return rb.bufferSize
}

func (rb *RingBuffer) remain() int {
	return rb.bufferSize - rb.length()
}

func (rb *RingBuffer) get(length int) []byte {
	if length <= 0 || length > rb.bufferSize {
		return nil
	}
	bufferLen := rb.length()
	if bufferLen <= 0 || bufferLen < length {
		return nil
	}

	result := make([]byte, length)
	offset := rb.header - rb.trailer
	if offset >= 0 {
		trailRemain := rb.bufferSize - rb.header
		if trailRemain >= length {
			copy(result, rb.buffer[rb.header:])
		} else {
			copy(result, rb.buffer[rb.header:])
			copy(result[trailRemain:], rb.buffer)
		}
	} else {
		copy(result, rb.buffer[rb.header:])
	}
	rb.header = (rb.header + length) % rb.bufferSize
	if rb.header == rb.trailer {
		rb.isEmpty = true
	}
	return result
}

func (rb *RingBuffer) Length() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.length()
}

func (rb *RingBuffer) Remain() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.remain()
}

func (rb *RingBuffer) IsEmpty() bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.isEmpty
}

func (rb *RingBuffer) Reset() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.header = 0
	rb.headerMark = 0
	rb.trailer = 0
	rb.trailerMark = 0
	rb.isEmpty = true
}

func (rb *RingBuffer) Read(v []byte) (int, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.isEmpty || len(v) <= 0 {
		return 0, nil
	}

}
