package libgorgel

import ()

type EnvelopeFunc func(float64) float64
type Envelope int

const (
	ENVELOPE_RECTANGULAR Envelope = iota
	ENVELOPE_LINEAR
	ENVELOPE_ADSR
	NUMBER_OF_ENVELOPES
)

func RectangularEnvelope(s float64) float64 {
	return 1
}

func LinearEnvelope(s float64) float64 {
	return 1 - s
}

func AdsrEnvelope(s float64) float64 {
	// See http://en.wikipedia.org/wiki/Synthesizer#ADSR_envelope
	if s < 0.25 {
		return s / 0.25
	}
	if s < 0.5 {
		return 0.5 + 0.5*(1-(s-0.25)/0.25)
	}
	if s < 0.75 {
		return 0.5
	}
	return 0.5 - 0.5*(s-0.75)/0.25
}

func fillEnvelopes(envelopes *[]EnvelopeFunc) {
	*envelopes = make([]EnvelopeFunc, NUMBER_OF_ENVELOPES)
	(*envelopes)[ENVELOPE_RECTANGULAR] = RectangularEnvelope
	(*envelopes)[ENVELOPE_LINEAR] = LinearEnvelope
	(*envelopes)[ENVELOPE_ADSR] = AdsrEnvelope
}
