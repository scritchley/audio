package audio

import "math"

type Filter struct {
	CutOff                 Processor
	cutOffBuffer           []float32
	Resonance              Processor
	resonanceBuffer        []float32
	sampleRate             float32
	in1, in2, in3, in4     float32
	out1, out2, out3, out4 float32
}

func NewFilter() *Filter {
	return &Filter{sampleRate: DefaultSampleRate}
}

func (f *Filter) Process(data []float32, channels int) {
	if len(f.cutOffBuffer) < len(data) {
		f.cutOffBuffer = make([]float32, len(data))
	}
	if f.CutOff != nil {
		f.CutOff.Process(f.cutOffBuffer, channels)
	}
	if len(f.resonanceBuffer) < len(data) {
		f.resonanceBuffer = make([]float32, len(data))
	}
	if f.Resonance != nil {
		f.Resonance.Process(f.resonanceBuffer, channels)
	}
	// Apply filter to buffer.
	for i := 0; i < len(data); i += channels {
		for ch := 0; ch < channels; ch++ {
			sample := ch + i
			fl := float32(math.Abs(float64(f.cutOffBuffer[sample])) * 1.16)
			fb := float32((2 * math.Abs(float64(f.resonanceBuffer[sample]))) * float64(1.0-0.15*fl*fl))
			data[sample] -= f.out4 * fb
			data[sample] *= 0.35013 * (fl * fl) * (fl * fl)
			f.out1 = data[sample] + 0.3*f.in1 + (1-fl)*f.out1 // Pole 1
			f.in1 = data[sample]
			f.out2 = f.out1 + 0.3*f.in2 + (1-fl)*f.out2 // Pole 2
			f.in2 = f.out1
			f.out3 = f.out2 + 0.3*f.in3 + (1-fl)*f.out3 // Pole 3
			f.in3 = f.out2
			f.out4 = f.out3 + 0.3*f.in4 + (1-fl)*f.out4 // Pole 4
			f.in4 = f.out3
			data[sample] = f.out4
		}
	}
}
