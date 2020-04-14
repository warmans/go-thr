package encoding

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNext_DecodeMultipleResponses(t *testing.T) {

	buff, err := hex.DecodeString(`f07e7f060200010c240002006b001f01f7000000f000010c24027e7f06024c36496d616765547970653a6d61696e004c36496d61676556657273696f6e3a312e332e312e302e6b00f7000000f000010c24024d000c00000b000100000004000000006b0031010000f7000000f000010c24024d000d00000b00010000000400000000000000000000f7000000f000010c24024d000e00000800010000000100000000000000000000f7000000f000010c24024d011400000b00010000000400000000000000000000f7000000f000010c24024d011500000b00010000000400000000000000000000f7000000f000010c24024d011600000b00010000000400000000000000000000f7000000f000010c24024d011700000b00010000000400000000000000000000f7000000f000010c24024d011800000b00010000000400000000000000000000f7000000`)
	require.NoError(t, err)

	msgs := make([]*Message, 0)
	for len(buff) > 0 {
		var msg *Message
		msg, buff = Next(buff)
		if msg != nil {
			msgs = append(msgs, msg)
		}
	}

	require.EqualValues(t, 8, len(msgs))
	require.EqualValues(t, "f000010c24024d000c00000b000100000004000000006b0031010000f7", msgs[0].Hex())
	require.EqualValues(t, "f000010c24024d000d00000b00010000000400000000000000000000f7", msgs[1].Hex())
	require.EqualValues(t, "f000010c24024d000e00000800010000000100000000000000000000f7", msgs[2].Hex())

}

func TestNext_SingleMessage(t *testing.T) {
	msg, buff := Next(hexMustDecode("f000010c24024d000c00000b000100000004000000006b0031010000f7"))
	require.EqualValues(t, 0, len(buff))
	require.EqualValues(t, [3]byte{0x00, 0x01, 0x0c}, msg.ManufacturerCode)
	require.EqualValues(t, [3]byte{0x24, 0x02, 0x4d}, msg.Preamble)
	require.EqualValues(t, 0x00, msg.MessageType)
	require.EqualValues(t, 0x0c, msg.SequenceNum)
	require.EqualValues(t, [2]byte{0x00, 0x00}, msg.Reserved1)
	require.EqualValues(t, 0x0b, msg.PayloadType)
	require.EqualValues(t, "000100000004000000006b0031010000", hex.EncodeToString(msg.Payload))
}

func hexMustDecode(str string) []byte {
	b, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return b
}