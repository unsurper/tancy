package extra

import (
	"encoding/binary"
	"github.com/liveej/go808/errors"
)

// 4G基站信息
type Extra_0x5d struct {
	serialized []byte
	value      Extra_0x5d_Value
}

type Extra_0x5d_Value struct {
	Cnt  int8
	MCC  int
	MNC  int
	LAC  int
	CID  uint32
	RSSI int8
}

func NewExtra_0x5d(val byte) *Extra_0x5d {
	extra := Extra_0x5d{}
	extra.serialized = []byte{val}
	return &extra
}

func (Extra_0x5d) ID() byte {
	return byte(TypeExtra_0x5d)
}

func (extra Extra_0x5d) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x5d) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x5d) Decode(data []byte) (int, error) {
	if len(data) < 11 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value.Cnt = int8(data[0])
	extra.value.MCC = int(binary.BigEndian.Uint16(data[1:3]))
	extra.value.MNC = int(data[3])
	extra.value.LAC = int(binary.BigEndian.Uint16(data[4:6]))
	extra.value.CID = binary.BigEndian.Uint32(data[6:10])
	extra.value.RSSI = int8(data[10])
	return 11, nil
}
