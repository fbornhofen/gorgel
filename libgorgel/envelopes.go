package libgorgel

import (
	"math"
)

type EnvelopeFunc func(float64) float64
type Envelope int

const (
	ENVELOPE_DEFAULT Envelope = iota
	ENVELOPE_RECTANGULAR
	ENVELOPE_LINEAR
	ENVELOPE_ADSR
	ENVELOPE_POLY
	NUMBER_OF_ENVELOPES
)

func StringToEnvelope(s string) Envelope {
	switch s {
	case "lin":
		return ENVELOPE_LINEAR
	case "rect":
		return ENVELOPE_RECTANGULAR
	case "adsr":
		return ENVELOPE_ADSR
	case "poly":
		return ENVELOPE_POLY
	default:
		return ENVELOPE_DEFAULT
	}
}

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

func PolyEnvelope(s float64) float64 {
	return -(math.Pow(1-s, 5) - math.Pow(1-s, 4)) * 12
}

func fillEnvelopes(envelopes *[]EnvelopeFunc) {
	*envelopes = make([]EnvelopeFunc, NUMBER_OF_ENVELOPES)
	(*envelopes)[ENVELOPE_DEFAULT] = RectangularEnvelope
	(*envelopes)[ENVELOPE_RECTANGULAR] = RectangularEnvelope
	(*envelopes)[ENVELOPE_LINEAR] = LinearEnvelope
	(*envelopes)[ENVELOPE_ADSR] = AdsrEnvelope
	(*envelopes)[ENVELOPE_POLY] = PolyEnvelope
}
