package y3

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObservable(t *testing.T) {
	buf := []byte{0x81, 0x16, 0xb0, 0x14, 0x10, 0x4, 0x79, 0x6f, 0x6d, 0x6f, 0x11, 0x2, 0x43, 0xe4, 0x92, 0x8, 0x13, 0x2, 0x41, 0xf0, 0x14, 0x2, 0x42, 0x20, 0x81, 0x16, 0xb0, 0x14, 0x10, 0x4, 0x79, 0x6f, 0x6d, 0x6f, 0x11, 0x2, 0x43, 0xe4, 0x92, 0x8, 0x13, 0x2, 0x41, 0xf0, 0x14, 0x2, 0x42, 0x20, 0x81, 0x16, 0xb0, 0x14, 0x10, 0x4, 0x79, 0x6f, 0x6d, 0x6f, 0x11, 0x2, 0x43, 0xe4, 0x92, 0x8, 0x13, 0x2, 0x41, 0xf0, 0x14, 0x2, 0x42, 0x20}
	var err error = nil
	var count int = 0

	callback1 := func(v []byte) (interface{}, error) {
		if (v[0] == 17) && (v[1] == 2) && (v[2] == 67) && (v[3] == 228) {
			count++
			return "ok1", nil
		} else {
			err = errors.New("fail")
			return nil, errors.New("fail")
		}

	}

	callback2 := func(v []byte) (interface{}, error) {
		if (v[0] == 19) && (v[1] == 2) && (v[2] == 65) && (v[3] == 240) {
			count++
			return "ok2", nil
		} else {
			err = errors.New("fail")
			return nil, errors.New("fail")
		}

	}

	reader := bytes.NewReader(buf)

	source := FromStream(reader)

	consumer1 := source.Subscribe(0x11).OnObserve(callback1)

	consumer2 := source.Subscribe(0x13).OnObserve(callback2)

	go func() {
		for range consumer2 {
			if count == 6 || err != nil {
				break
			}
		}

	}()

	for range consumer1 {
		if count == 6 || err != nil {
			break
		}
	}
	assert.NoError(t, err, fmt.Sprintf("subscribe error:%v", err))
	assert.Equal(t, 6, count, fmt.Sprintf("testing observable %v: %v", 6, count))
	testPrintf("count=%v, observable_result=%v, err=%v\n", 6, count, err)

}
