package audio

import (
	"sync"
	"sync/atomic"
)

const (
	// DefaultBufferSize is a default buffer size to use.
	DefaultBufferSize = 2048
	// DefaultSampleRate is a default sample rate to use.
	DefaultSampleRate = 44100
)

type Processor interface {
	Process(data []float32, channels int)
}

type Connection chan bool

func (c Connection) Disconnect() {
	c <- true
}

func (c Connection) Wrap(conn Connection) Connection {
	go func() {
		<-c
		conn.Disconnect()
	}()
	return c
}

type Writer interface {
	Write(...float32)
}

// Buffer is a buffer of float32 values representing a single channel.
type Buffer struct {
	mtx  sync.Mutex
	data []float32
}

// MakeBuffer returns a new Buffer with the provided size.
func MakeBuffer(size int) Buffer {
	return Buffer{data: make([]float32, size)}
}

func (b Buffer) Resize(size int) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.data = make([]float32, size)
}

func (b Buffer) Clear() {
	for i := range b.data {
		b.data[i] = 0
	}
}

type Constant struct {
	sampleRate   int
	targetOffset atomic.Value
	offset       atomic.Value
	smoothing    float32
}

func NewConstant(offset float32) *Constant {
	var targetOffsetValue, offsetValue atomic.Value
	targetOffsetValue.Store(offset)
	offsetValue.Store(offset)
	return &Constant{
		sampleRate:   DefaultSampleRate,
		targetOffset: targetOffsetValue,
		offset:       offsetValue,
		smoothing:    0.1,
	}
}

func (c *Constant) SetOffset(offset float32) {
	c.targetOffset.Store(offset)
}

func (c *Constant) SetGlideMs(glideMs float32) {
	c.smoothing = (glideMs / 1000)
}

func (c *Constant) Process(data []float32, channels int) {
	for i := 0; i < len(data); i += channels {
		current := c.offset.Load().(float32)
		target := c.targetOffset.Load().(float32)
		current += (target - current) / (c.smoothing * float32(c.sampleRate))
		c.offset.Store(current)
		for ch := 0; ch < channels; ch++ {
			data[ch+i] = current
		}
	}
}

type Amplifier struct {
	Gain       Processor
	gainBuffer []float32
}

func NewAmplifier() *Amplifier {
	return &Amplifier{}
}

func (a *Amplifier) Process(data []float32, channels int) {
	if len(a.gainBuffer) < len(data) {
		a.gainBuffer = make([]float32, len(data))
	}
	if a.Gain != nil {
		a.Gain.Process(a.gainBuffer, channels)
	}
	for i := 0; i < len(data); i += channels {
		for ch := 0; ch < channels; ch++ {
			data[ch+i] *= a.gainBuffer[ch+i]
		}
	}
}
