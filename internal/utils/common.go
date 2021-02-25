package utils

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
