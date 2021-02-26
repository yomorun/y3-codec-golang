package y3

import (
	"io"
	"sync"

	"github.com/yomorun/y3-codec-golang/pkg/common"
)

// Iterable iterate through and get the data of observe
type Iterable interface {
	Observe() <-chan interface{}
}

// Observable provide subscription and notification processing
type Observable interface {
	Iterable
	Subscribe(key byte) Observable
	OnObserve(function func(v []byte) (interface{}, error)) chan interface{}
}

type observableImpl struct {
	iterable Iterable
}

type iterableImpl struct {
	next                   chan interface{}
	subscribers            []chan interface{}
	mutex                  sync.RWMutex
	producerAlreadyCreated bool
}

func (i *iterableImpl) Observe() <-chan interface{} {
	ch := make(chan interface{})
	i.mutex.Lock()
	i.subscribers = append(i.subscribers, ch)
	i.mutex.Unlock()
	i.connect()
	return ch
}

func (i *iterableImpl) connect() {
	i.mutex.Lock()
	if !i.producerAlreadyCreated {
		go i.produce()
		i.producerAlreadyCreated = true
	}
	i.mutex.Unlock()
}

func (i *iterableImpl) produce() {
	defer func() {
		i.mutex.RLock()
		for _, subscriber := range i.subscribers {
			close(subscriber)
		}
		i.mutex.RUnlock()
	}()

	for {
		select {
		case item, ok := <-i.next:
			if !ok {
				return
			}
			i.mutex.RLock()
			for _, subscriber := range i.subscribers {
				subscriber <- item
			}
			i.mutex.RUnlock()
		}
	}
}

func (o *observableImpl) Observe() <-chan interface{} {
	return o.iterable.Observe()
}

//FromStream reads data from reader
func FromStream(reader io.Reader) Observable {

	f := func(next chan interface{}) {
		defer close(next)
		for {
			buf := make([]byte, 3*1024)
			n, err := reader.Read(buf)

			if err != nil {
				break
			} else {
				value := buf[:n]
				//fmt.Printf("%v:\t $1 on y3 value=%#v\n", time.Now().Format("2006-01-02 15:04:05"), value)
				next <- value
			}
		}
	}

	return createObservable(f)
}

//Processing callback function when there is data
func (o *observableImpl) OnObserve(function func(v []byte) (interface{}, error)) chan interface{} {
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

//Get the value of the subscribe key from the stream
func (o *observableImpl) Subscribe(key byte) Observable {

	f := func(next chan interface{}) {
		defer close(next)

		resultBuffer := make([]byte, 0)
		var (
			flow   int32 = 0 //0 start，1 get length，2 get value
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
				reject = false

				for item := range nextitem.(chan interface{}) {

					buf := item.([]byte)

					for _, b := range buf {
						if reject {
							break
						}

						if flow == 0 && (b<<2)>>2 == key { //T
							flow = 1
							resultBuffer = append(resultBuffer, b)
							continue
						}

						if flow == 1 { // L
							resultBuffer = append(resultBuffer, b)
							l, e := common.DecodeLength(resultBuffer[1 : length+1]) //l 是value占字节，s是l占字节

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

		}

	}

	return createObservable(f)

}

func createObservable(f func(next chan interface{})) Observable {
	next := make(chan interface{})
	subscribers := make([]chan interface{}, 0)

	go f(next)
	return &observableImpl{iterable: &iterableImpl{next: next, subscribers: subscribers}}
}

//filter root data from the stream
func filterRoot(observe <-chan interface{}) <-chan interface{} {
	next := make(chan interface{})

	f := func(observe <-chan interface{}) {
		defer close(next)
		rootBuffer := make([]byte, 0)
		var send chan interface{}

		var (
			rootflow   int32 = 0 //0 start，1 get length，2 get value
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
				i := 0

				for {
					if i == len(buf) {
						break
					}

					b := buf[i]

					if rootflow == 0 {
						if (b<<2)>>2 == rootkey {
							rootflow = 1
							rootBuffer = append(rootBuffer, b)
						}
						i++
						continue
					}

					if rootflow == 1 { // L
						rootBuffer = append(rootBuffer, b)
						l, e := common.DecodeLength(rootBuffer[1 : rootlength+1])

						if e != nil {
							rootlength++
						} else {
							rootvalue = l
							rootflow = 2
							send = make(chan interface{})
							next <- send
						}
						i++
						continue
					}

					if rootflow == 2 {
						if (rootvalue - index) > int32(len(buf[i:])) {
							send <- buf[i:]
							index = index + int32(len(buf[i:]))
							i = len(buf)
							continue
						} else if (rootvalue - index) == int32(len(buf[i:])) {
							send <- buf[i:]
							rootflow = 0
							rootlength = 1
							rootvalue = 0
							index = 0
							rootBuffer = make([]byte, 0)
							close(send)
							i = len(buf)
							continue
						} else if (rootvalue - index) < int32(len(buf[i:])) {
							send <- buf[i:(int32(i) + (rootvalue - index))]
							i = i + int(rootvalue-index)

							rootflow = 0
							rootlength = 1
							rootvalue = 0
							index = 0
							rootBuffer = make([]byte, 0)
							close(send)
							continue
						}
					}
				}

			}
		}
	}

	go f(observe)

	return next
}
