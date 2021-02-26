package y3

import (
	"fmt"
)

// signal is builder for PrimitivePacketEncoder
type signal struct {
	encoder *PrimitivePacketEncoder
}

// createSignal create a signal
func createSignal(key byte) *signal {
	return &signal{encoder: NewPrimitivePacketEncoder(int(key))}
}

// SetString set a string Value for the signal
func (s *signal) SetString(v string) *signal {
	s.encoder.SetStringValue(v)
	return s
}

// SetString set a int64 Value for the signal
func (s *signal) SetInt64(v int64) *signal {
	s.encoder.SetInt64Value(v)
	return s
}

// SetString set a float64 Value for the signal
func (s *signal) SetFloat64(v float64) *signal {
	s.encoder.SetFloat64Value(v)
	return s
}

// ToEncoder return current PrimitivePacketEncoder, and checking legality
func (s *signal) ToEncoder(allow func(key byte) bool) *PrimitivePacketEncoder {
	if allow != nil && !allow(byte(s.encoder.seqID)) {
		panic(fmt.Errorf("it is not allowed to use this key to create a signal: %#x", byte(s.encoder.seqID)))
	}

	return s.encoder
}
