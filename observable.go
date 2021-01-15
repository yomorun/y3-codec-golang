package y3

import (
	"github.com/yomorun/y3-codec-golang/pkg/spec/encoding"
)

type Iterable interface {
	Observe() <-chan interface{}
}

type Observable interface {
	Iterable
	Subscribe(key byte) Observable
	OnObserve(function func(v []byte) (interface{}, error)) chan interface{}
}

type ObservableImpl struct {
	iterable Iterable
}

type IterableImpl struct {
	channel chan interface{}
}

func (i *IterableImpl) Observe() <-chan interface{} {
	return i.channel
}

func (o *ObservableImpl) Observe() <-chan interface{} {
	return o.iterable.Observe()
}

func (o *ObservableImpl) OnObserve(function func(v []byte) (interface{}, error)) chan interface{} {
	_next := make(chan interface{})

	f := func(next chan interface{}) {
		defer close(next)

		observe := o.Observe()

		for {
			select {
			case item, ok := <-observe:
				if !ok {
					return
				}
				buf := item.([]byte)
				value, err := function(buf)
				if err != nil {
					return
				}

				next <- value
			}
		}
	}

	go f(_next)

	return _next
}

func (o *ObservableImpl) Subscribe(key byte) Observable {

	f := func(next chan interface{}) {
		defer close(next)

		resultBuffer := make([]byte, 0)
		var (
			flow   int32 = 0 //0 未监听到，1 判断长度，2 判断v
			length int32 = 1
			value  int32 = 0
			reject bool  = false
		)

		_observe := o.Observe()
		observe := filterRoot(_observe)

		for {
			select {
			case nextitem, ok := <-observe:
				if !ok {
					return
				}
				for item := range nextitem.(chan interface{}) {

					if !reject {

						buf := item.([]byte)

						for _, b := range buf {
							if flow == 0 && (b<<2)>>2 == key { //T
								flow = 1
								resultBuffer = append(resultBuffer, b)
								continue
							}

							if flow == 1 { // L
								resultBuffer = append(resultBuffer, b)
								l, e := decodeLength(resultBuffer[1 : length+1]) //l 是value占字节，s是l占字节

								if e != nil {
									length++
								} else {
									value = l
									flow = 2
								}
								continue
							}

							if flow == 2 {
								l := len(resultBuffer)
								if int32(l) == (length + value) {
									resultBuffer = append(resultBuffer, b)
									next <- resultBuffer
									flow = 0
									length = 1
									value = 0
									resultBuffer = make([]byte, 0)
									reject = true
									break
								} else if int32(l) < (length + value) {
									resultBuffer = append(resultBuffer, b)
									continue
								}
							}
						}
					}
				}
				reject = false
			}

		}

	}

	return createObservable(f)

}

func decodeLength(buf []byte) (length int32, err error) {
	varCodec := encoding.VarCodec{}
	err = varCodec.DecodePVarInt32(buf, &length)
	return
}

func createObservable(f func(next chan interface{})) Observable {
	next := make(chan interface{})
	go f(next)
	return &ObservableImpl{iterable: &IterableImpl{channel: next}}
}

func filterRoot(observe <-chan interface{}) <-chan interface{} {
	next := make(chan interface{})

	f := func(observe <-chan interface{}) {
		defer close(next)
		rootBuffer := make([]byte, 0)
		var send chan interface{}

		var (
			rootflow   int32 = 0 //0 未监听到，1 判断长度，2 判断v
			rootlength int32 = 1
			rootvalue  int32 = 0
			rootkey    byte  = 0x01
			index      int32 = 0
		)

		for {
			select {
			case item, ok := <-observe:
				if !ok {
					return
				}

				buf := item.([]byte)

				for _, b := range buf {
					if rootflow == 0 && ((b<<2)>>2 == rootkey) {
						rootflow = 1
						rootBuffer = append(rootBuffer, b)
						continue
					}

					if rootflow == 1 { // L
						rootBuffer = append(rootBuffer, b)
						l, e := decodeLength(rootBuffer[1 : rootlength+1]) //l 是value占字节，s是l占字节

						if e != nil {
							rootlength++
						} else {
							rootvalue = l
							rootflow = 2
						}
						continue
					}

					if rootflow == 2 {
						index++

						if index == 1 {
							send = make(chan interface{})
							next <- send
						}
						if index < rootvalue {
							send <- []byte{b}
						} else if index == rootvalue {
							send <- []byte{b}
							rootflow = 0
							rootlength = 1
							rootvalue = 0
							index = 0
							rootBuffer = make([]byte, 0)
							close(send)
						}
					}
				}
			}
		}
	}

	go f(observe)

	return next
}