package y3

import (
	"fmt"
	"time"
)

var (
	// enabledTestPrintf set whether to print debug information in the test
	enabledTestPrintf = true
)

// observableTester use Observable to listen the node
type observableTester struct {
	observe       byte
	sourceChannel chan interface{}
	source        Observable
}

// newObservableTester creat a observableTester
func newObservableTester(observe byte) *observableTester {
	return &observableTester{observe: observe}
}

// testDecoder is a shortcut to perform decoding tests
func testDecoder(observe byte, buf []byte, callback func(v []byte) (interface{}, error)) {
	newObservableTester(observe).
		Init(callback).
		Write(buf).
		CloseWith(100)
}

// Init create a channel for testing
func (t *observableTester) Init(callback func(v []byte) (interface{}, error)) *observableTester {
	t.sourceChannel = make(chan interface{})

	t.source = &observableImpl{iterable: &iterableImpl{channel: t.sourceChannel}}

	consumer := t.source.Subscribe(t.observe).OnObserve(callback)

	go func() {
		for c := range consumer {
			if c != 0 {
			}
		}
	}()

	return t
}

// Write is used to write data to the Channel
func (t *observableTester) Write(buf []byte) *observableTester {
	t.sourceChannel <- buf
	return t
}

// Close is used to close the Channel
func (t *observableTester) Close() {
	close(t.sourceChannel)
}

// Close is used to close the Channel with waiting time
func (t *observableTester) CloseWith(millisecond int64) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
	close(t.sourceChannel)
}

// testPrintf print debug output in test cases
func testPrintf(format string, a ...interface{}) {
	if enabledTestPrintf {
		fmt.Printf(format, a...)
	}
}
