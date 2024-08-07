package netmsg

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

func GetBytesFromIdAndStruct[T any](id byte, s T) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bytes := make([]byte, 1+buf.Len())
	bytes[0] = id
	copy(bytes[1:], buf.Bytes())

	return bytes
}

func GetStructFromBytes[T any](bytesWithoutId []byte) T {
	buf := bytes.NewBuffer(bytesWithoutId)
	dec := gob.NewDecoder(buf)

	var result T
	err := dec.Decode(&result)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return result
}
