package utils

// MSB 描述了`1000 0000`, 用于表示后续字节仍然是该变长类型值的一部分
const MSB byte = 0x80

// DropMSB 描述了`0111 1111`, 用于去除标识位使用
const DropMSB = 0x3F

// DropMSBArrayFlag 描述了`0011 1111`, 用于去除标识位使用
const DropMSBArrayFlag = 0x3F

// ArrayFlag 描述了`0100 0000`, 用于表示该节点的Value为Slice类型
const ArrayFlag = 0x40

// KeyOfArrayItem 描述数组项的TLV的sid值
const KeyOfArrayItem = 0x00
const KeyStringOfArrayItem = "0x00"
