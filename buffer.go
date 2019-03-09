package audio

import (
	"sync/atomic"
	"time"
)

const (
	// DefaultBufferSize is a default buffer size to use.
	DefaultBufferSize = 2048
	// DefaultSampleRate is a default sample rate to use.
	DefaultSampleRate = 44100
)

// Processor is an interface that processes an audio buffer.
type Processor interface {
	// Process processes the provided audio data. The number
	// of channels is specified using the channels argument.
	// Channel data is interlaced in the data slice.
	Process(data []float32, channels int)
}

// Constant is a constant source that writes its offset value to the output buffer.
type Constant struct {
	sampleRate         int
	targetOffset       atomic.Value
	offset             atomic.Value
	transitionDuration time.Duration
}

// NewConstant retusn a new Constant with the provided initial offset.
func NewConstant(initialOffset float32) *Constant {
	var targetOffsetValue, offsetValue atomic.Value
	targetOffsetValue.Store(initialOffset)
	offsetValue.Store(initialOffset)
	return &Constant{
		sampleRate:         DefaultSampleRate,
		targetOffset:       targetOffsetValue,
		offset:             offsetValue,
		transitionDuration: 0,
	}
}

func (c *Constant) SetOffset(offset float32) {
	c.targetOffset.Store(offset)
}

func (c *Constant) SetTransitionTime(duration time.Duration) *Constant {
	c.transitionDuration = duration
	return c
}

func (c *Constant) getValue() float32 {
	if c.transitionDuration == 0 {
		return c.targetOffset.Load().(float32)
	}
	current := c.offset.Load().(float32)
	target := c.targetOffset.Load().(float32)
	current += (target - current) / (0.1 * float32(c.sampleRate))
	c.offset.Store(current)
	return current
}

func (c *Constant) Process(data []float32, channels int) {
	for i := 0; i < len(data); i += channels {
		val := c.getValue()
		for ch := 0; ch < channels; ch++ {
			data[ch+i] = val
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

// ProcessorChain is a chain of Processors.
type ProcessorChain []Processor

// NewProcessorChain returns a ProcessorChain using the provided Processors.
func NewProcessorChain(processors ...Processor) ProcessorChain {
	return ProcessorChain(processors)
}

// Process calls each Processor sequentially on the provided data.
func (c ProcessorChain) Process(data []float32, channels int) {
	for i := range c {
		c[i].Process(data, channels)
	}
}

// ProcessorFunc is a func that processes an audio input.s
type ProcessorFunc func(data []float32, channels int)

// Process calls ProcessorFunc passing in the provided data.
func (p ProcessorFunc) Process(data []float32, channels int) {
	p(data, channels)
}

// Add reads a into the provided buffer and then adds b.
func Add(a, b Processor) Processor {
	var bBuffer []float32
	return ProcessorFunc(func(data []float32, channels int) {
		if len(bBuffer) < len(data) {
			bBuffer = make([]float32, len(data))
		}
		a.Process(data, channels)
		b.Process(bBuffer, channels)
		for i := range data {
			data[i] += bBuffer[i]
		}
	})
}

// Multiply reads a into the provided buffer and then multiples by b.
func Multiply(a, b Processor) Processor {
	var bBuffer []float32
	return ProcessorFunc(func(data []float32, channels int) {
		if len(bBuffer) < len(data) {
			bBuffer = make([]float32, len(data))
		}
		a.Process(data, channels)
		b.Process(bBuffer, channels)
		for i := range data {
			data[i] *= bBuffer[i]
		}
	})
}
