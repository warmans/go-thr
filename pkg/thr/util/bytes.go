package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func SingleByteInt(num int8) byte {
	bs := &bytes.Buffer{}
	if err := binary.Write(bs, binary.LittleEndian, num); err != nil {
		panic(fmt.Sprintf("failed to encode int as byte: %s", err.Error()))
	}
	return bs.Bytes()[0]
}