package libgorgel

import (

)

type EnvelopeFunc func(float64) float64
type Envelope int

const (
	ENVELOPE_RECTANGULAR Envelope = iota
	ENVELOPE_LINEAR
	NUMBER_OF_ENVELOPES
)

func RectangularEnvelope(s float64) float64 {
	return 1
}

func LinearEnvelope(s float64) float64 {
	return 1-s
}

func fillEnvelopes(envelopes *[]EnvelopeFunc) {
	*envelopes = make([]EnvelopeFunc, NUMBER_OF_ENVELOPES)
	(*envelopes)[ENVELOPE_RECTANGULAR] = RectangularEnvelope
	(*envelopes)[ENVELOPE_LINEAR] = LinearEnvelope
}