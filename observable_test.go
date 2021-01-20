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
	var err1 error = nil
	var err2 error = nil
	var count1 int = 0
	var count2 int = 0

	callback1 := func(v []byte) (interface{}, error) {
		if (v[0] == 17) && (v[1] == 2) && (v[2] == 67) && (v[3] == 228) {
			count1++
			return "ok1", nil
		} else {
			err1 = errors.New("fail")
			return nil, errors.New("fail")
		}

	}

	callback2 := func(v []byte) (interface{}, error) {
		if (v[0] == 19) && (v[1] == 2) && (v[2] == 65) && (v[3] == 240) {
			count2++
			return "ok2", nil
		} else {
			err2 = errors.New("fail")
			return nil, errors.New("fail")
		}

	}

	reader := bytes.NewReader(buf)

	source := FromStream(reader)

	consumer1 := source.Subscribe(0x11).OnObserve(callback1)

	consumer2 := source.Subscribe(0x13).OnObserve(callback2)

	for range consumer1 {
		if count1 == 3 || err1 != nil {
			break
		}
	}

	for range consumer2 {
		if count2 == 3 || err2 != nil {
			break
		}
	}

	assert.NoError(t, err1, fmt.Sprintf("subscribe2 error:%v", err1))
	assert.Equal(t, 3, count1, fmt.Sprintf("testing observable1 %v: %v", 3, count1))
	assert.NoError(t, err2, fmt.Sprintf("subscribe2 error:%v", err2))
	assert.Equal(t, 3, count2, fmt.Sprintf("testing observable2 %v: %v", 3, count2))
	testPrintf("count1=%v,count2=%v, observable_result=%v, err1=%v,err2=%v\n", 3, 3, "ok", err1, err2)

}
