package utils

import (
	"bytes"
	"encoding/binary"
	"time"
	"unsafe"
)

func UintToByteArray(num uint) []byte {
	buf := new(bytes.Buffer)
	size := unsafe.Sizeof(num)
	switch size {
	case 4:
		binary.Write(buf, binary.BigEndian, uint32(num))
	case 8:
		binary.Write(buf, binary.BigEndian, uint64(num))
	default:
		panic("unsupported uint size")
	}
	return buf.Bytes()
}

func ByteArrayToUint(b []byte) uint {
	size := len(b)
	buf := bytes.NewReader(b)
	switch size {
	case 4:
		var num uint32
		binary.Read(buf, binary.BigEndian, &num)
		return uint(num)
	case 8:
		var num uint64
		binary.Read(buf, binary.BigEndian, &num)
		return uint(num)
	default:
		panic("unsupported byte array size")
	}
}

func GetCurrentTimeAsFloat() float64 {
	return float64(time.Now().UnixNano()) / 1000000000.0
}
