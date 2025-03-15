package util

import (
	"sync"
)

const (
	RingBufferDefaultSize = 1024
	FixedDataLengthHeader = 4
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

func (rb *RingBuffer) check(length int) {
	remain := rb.remain()
	if remain < length {
		rb.reallocateMemory(length)
	}
}

func (rb *RingBuffer) reallocateMemory(size int) {
	multi := ((size + rb.bufferSize) / rb.bufferSize) + 1
	newBufferSize := rb.bufferSize * multi
	tmpBuffer := make([]byte, newBufferSize)

	if !rb.isEmpty {
		offset := rb.header - rb.trailer
		if offset > 0 {
			// 头后尾前
			copy(tmpBuffer, rb.buffer[:rb.trailer])
			copy(tmpBuffer[(newBufferSize-(rb.bufferSize-rb.header)):], rb.buffer[rb.header:])
			rb.header = newBufferSize - rb.bufferSize + rb.header
		} else if offset < 0 {
			// 头前尾后
			copy(tmpBuffer[rb.header:rb.trailer], rb.buffer[rb.header:rb.trailer])
		} else {
			// 先copy头部
			copy(tmpBuffer[0:], rb.buffer[rb.header:])
			if rb.trailer != 0 {
				// 再copy尾部
				copy(tmpBuffer[rb.bufferSize-rb.header:], rb.buffer[:rb.trailer])
			}
			// 更新头尾
			rb.header = 0
			rb.trailer = rb.bufferSize
		}
	}

	rb.buffer = tmpBuffer
	rb.bufferSize = newBufferSize
	if rb.logHook != nil {
		rb.logHook.Infof("ring buffer reallocate size", rb.bufferSize, rb.header, rb.trailer)
	}

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

func (rb *RingBuffer) put(data []byte) {
	// 预检,保证有足够空间
	rb.check(len(data))

	trailRemain := rb.bufferSize - rb.trailer
	if trailRemain < len(data) {
		if trailRemain <= 0 {
			// 尾部无空间
			copy(rb.buffer, data)
		} else {
			// 尾部有空间,分两次copy
			copy(rb.buffer[rb.trailer:], data[:trailRemain])
			copy(rb.buffer[0:], data[trailRemain:])
		}
	} else {
		copy(rb.buffer[rb.trailer:], data)
	}
	rb.trailer = (rb.trailer + len(data)) % rb.bufferSize
	rb.isEmpty = false
}

func (rb *RingBuffer) mark() {
	rb.headerMark = rb.header
	rb.trailerMark = rb.trailer
}

func (rb *RingBuffer) rollback() {
	rb.header = rb.headerMark
	rb.trailer = rb.trailerMark
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

func (rb *RingBuffer) Get(length int) []byte {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	return rb.get(length)
}

func (rb *RingBuffer) Put(data []byte) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.put(data)
}

// 读指定长度数据
func (rb *RingBuffer) Read(dest []byte) (int, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.isEmpty || len(dest) <= 0 {
		return 0, nil
	}

	bufferLen := rb.length()
	length := len(dest)
	if length > bufferLen {
		length = bufferLen
	}

	data := rb.get(length)
	if data == nil || len(data) <= 0 {
		return 0, nil
	}
	copy(dest, data)
	return length, nil
}

// PutFixedData 定长数据,length位固定4byte
func (rb *RingBuffer) PutFixedData(data []byte) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	length := len(data)
	lengthByte := Int32ToBytes(int32(length))
	rb.put(lengthByte)
	rb.put(data)
}

// GetFixedData 获取定长数据(每调用一次返回一段定长数据)
func (rb *RingBuffer) GetFixedData() []byte {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.length() < FixedDataLengthHeader {
		return nil
	}

	rb.mark()
	lenByte := rb.get(FixedDataLengthHeader)
	if lenByte == nil || len(lenByte) <= 0 {
		rb.rollback()
		return nil
	}
	dataLen := ByteToInt32(lenByte)
	if int(dataLen)+FixedDataLengthHeader > rb.length() {
		rb.rollback()
		return nil
	}

	result := rb.get(int(dataLen))
	if result == nil {
		rb.rollback()
		return nil
	}

	return result
}
