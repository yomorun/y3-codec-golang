package y3

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

// Signal is builder for PrimitivePacketEncoder
type Signal struct {
	encoder *PrimitivePacketEncoder
}

// CreateSignal create a Signal
func CreateSignal(key byte) *Signal {
	return &Signal{encoder: NewPrimitivePacketEncoder(int(key))}
}

// SetString set a string Value for the Signal
func (s *Signal) SetString(v string) *Signal {
	s.encoder.SetStringValue(v)
	return s
}

// SetString set a int64 Value for the Signal
func (s *Signal) SetInt64(v int64) *Signal {
	s.encoder.SetInt64Value(v)
	return s
}

// SetString set a float64 Value for the Signal
func (s *Signal) SetFloat64(v float64) *Signal {
	s.encoder.SetFloat64Value(v)
	return s
}

// ToEncoder return current PrimitivePacketEncoder, and checking legality
func (s *Signal) ToEncoder() *PrimitivePacketEncoder {
	if !utils.AllowableSignalKey(byte(s.encoder.seqID)) {
		panic(fmt.Errorf("it is not allowed to use this key to create a signal: %#x", byte(s.encoder.seqID)))
	}

	return s.encoder
}
