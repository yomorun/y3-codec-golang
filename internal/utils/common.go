package utils

import "reflect"

// MSB 描述了`1000 0000`, 用于表示后续字节仍然是该变长类型值的一部分
const MSB byte = 0x80

// DropMSB 描述了`0111 1111`, 用于去除标识位使用
const DropMSB = 0x3F

// DropMSBArrayFlag 描述了`0011 1111`, 用于去除标识位使用
const DropMSBArrayFlag = 0x3F

// SliceFlag 描述了`0100 0000`, 用于表示该节点的Value为Slice类型
const SliceFlag = 0x40

// KeyOfSliceItem 描述数组项的TLV的sid值
const KeyOfSliceItem = 0x00

// KeyStringOfSliceItem 描述数组项的TLV的sid值的字符串表示
const KeyStringOfSliceItem = "0x00"

// RootToken 描述根节点的sid值
const RootToken byte = 0x01

// TypeOfByteSlice Type of []byte
var TypeOfByteSlice = reflect.TypeOf([]byte{})

// TyeOfStringSlice Type of []string{}
var TyeOfStringSlice = reflect.TypeOf([]string{})

// TyeOfInt32Slice Type of []int32{}
var TyeOfInt32Slice = reflect.TypeOf([]int32{})

// TyeOfUint32Slice Type of []uint32{}
var TyeOfUint32Slice = reflect.TypeOf([]uint32{})

// TyeOfInt64Slice Type of []int64{}
var TyeOfInt64Slice = reflect.TypeOf([]int64{})

// TyeOfUint64Slice Type of []uint64{}
var TyeOfUint64Slice = reflect.TypeOf([]uint64{})

// TyeOfFloat32Slice Type of []float32{}
var TyeOfFloat32Slice = reflect.TypeOf([]float32{})

// TyeOfFloat64Slice Type of []float64{}
var TyeOfFloat64Slice = reflect.TypeOf([]float64{})

// TyeOfBoolSlice Type of []bool{}
var TyeOfBoolSlice = reflect.TypeOf([]bool{})
