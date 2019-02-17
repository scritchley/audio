package audio

type Envelope struct {
	Gate          Processor
	Attack        Processor
	AttackTimeMs  float32
	attackBuffer  []float32
	Decay         Processor
	DecayTimeMs   float32
	decayBuffer   []float32
	Sustain       Processor
	sustainBuffer []float32
	Release       Processor
	ReleaseTimeMs float32
	releaseBuffer []float32
	*Constant
}

func NewEnvelope() *Envelope {
	return &Envelope{
		Attack:        NewConstant(0).SetGlideMs(10),
		AttackTimeMs:  50,
		Decay:         NewConstant(0).SetGlideMs(10),
		DecayTimeMs:   1000,
		Sustain:       NewConstant(1).SetGlideMs(10),
		Release:       NewConstant(0).SetGlideMs(10),
		ReleaseTimeMs: 1000,
		Constant:      NewConstant(0).SetGlideMs(1000),
	}
}

func (e *Envelope) Process(data []float32, channels int) {
	if e.Gate != nil {
		e.Gate.Process(data, channels)
	}
	if len(e.attackBuffer) < len(data) {
		e.attackBuffer = make([]float32, len(data))
	}
	if e.Attack != nil {
		e.Attack.Process(e.attackBuffer, channels)
	}
	if len(e.decayBuffer) < len(data) {
		e.decayBuffer = make([]float32, len(data))
	}
	if e.Decay != nil {
		e.Decay.Process(e.decayBuffer, channels)
	}
	if len(e.sustainBuffer) < len(data) {
		e.sustainBuffer = make([]float32, len(data))
	}
	if e.Sustain != nil {
		e.Sustain.Process(e.sustainBuffer, channels)
	}
	if len(e.releaseBuffer) < len(data) {
		e.releaseBuffer = make([]float32, len(data))
	}
	if e.Release != nil {
		e.Release.Process(e.releaseBuffer, channels)
	}
	for i := 0; i < len(data); i += channels {
		for ch := 0; ch < channels; ch++ {
			if data[ch+i] > 0.05 {
				e.Constant.SetGlideMs(e.AttackTimeMs)
			} else {
				e.Constant.SetGlideMs(e.ReleaseTimeMs)
			}
			e.Constant.SetOffset(data[ch+i])
			data[ch+i] = e.Constant.getValue()
		}
	}
}
