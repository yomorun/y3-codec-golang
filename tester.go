package y3

import (
	"fmt"
	"time"
)

var (
	enabledTestPrintf = true
)

type observableTester struct {
	observe       byte
	sourceChannel chan interface{}
	source        Observable
}

func newObservableTester(observe byte) *observableTester {
	return &observableTester{observe: observe}
}

func testDecoder(observe byte, buf []byte, callback func(v []byte) (interface{}, error)) {
	newObservableTester(observe).
		Init(callback).
		Write(buf).
		CloseWith(100)
}

func (t *observableTester) Init(callback func(v []byte) (interface{}, error)) *observableTester {
	t.sourceChannel = make(chan interface{})

	t.source = &ObservableImpl{iterable: &IterableImpl{channel: t.sourceChannel}}

	consumer := t.source.Subscribe(t.observe).OnObserve(callback)

	go func() {
		for c := range consumer {
			if c != 0 {
			}
		}
	}()

	return t
}

func (t *observableTester) Write(buf []byte) *observableTester {
	t.sourceChannel <- buf
	return t
}

func (t *observableTester) Close() {
	close(t.sourceChannel)
}

func (t *observableTester) CloseWith(millisecond int64) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
	close(t.sourceChannel)
}

func testPrintf(format string, a ...interface{}) {
	if enabledTestPrintf {
		fmt.Printf(format, a...)
	}
}
