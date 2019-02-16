package audio

import (
	"math"
	"math/rand"
	"sync"
)

type Oscillator struct {
	Frequency  Processor
	mtx        sync.Mutex
	sampleRate int
	phase      []float32
	WaveFunc
}

func (o *Oscillator) SetWaveFunc(w WaveFunc) {
	o.WaveFunc = w
}

func (o *Oscillator) Process(data []float32, channels int) {
	if len(o.phase) != channels {
		o.phase = make([]float32, channels)
	}
	if o.Frequency != nil {
		o.Frequency.Process(data, channels)
	}
	for i := 0; i < len(data); i += channels {
		for ch := 0; ch < channels; ch++ {
			increment := float32(2 * float32(math.Pi) * NormalisedCVToFrequency(data[ch+i]) / float32(o.sampleRate))
			data[ch+i] = o.WaveFunc(o.phase[ch])
			o.phase[ch] += increment
			if o.phase[ch] > math.Pi {
				o.phase[ch] -= 2 * math.Pi
			}
		}
	}
}

func NewOscillator(w WaveFunc) *Oscillator {
	return &Oscillator{
		sampleRate: DefaultSampleRate,
		WaveFunc:   w,
	}
}

const (
	SineB = 4.0 / math.Pi
	SineC = -4.0 / (math.Pi * math.Pi)
	Q     = 0.775
	SineP = 0.225
)

type WaveFunc func(float32) float32

// Sine takes an input value from -Pi to Pi
// and returns a value between -1 and 1
func Sine(x32 float32) float32 {
	x := float64(x32)
	y := SineB*x + SineC*x*(math.Abs(x))
	y = SineP*(y*(math.Abs(y))-y) + y
	return float32(y)
}

const TriangleA = 2.0 / math.Pi

// Triangle takes an input value from -Pi to Pi
// and returns a value between -1 and 1
func Triangle(x float32) float32 {
	return (TriangleA * x) - 1.0
}

// Square takes an input value from -Pi to Pi
// and returns -1 or 1
func Square(x float32) float32 {
	if x >= 0.0 {
		return 1
	}
	return -1.0
}

const SawtoothA = 1.0 / math.Pi

// Triangle takes an input value from -Pi to Pi
// and returns a value between -1 and 1
func Sawtooth(x float32) float32 {
	return SawtoothA * x
}

func Noise(x float32) float32 {
	return 2*rand.Float32() - 1
}
