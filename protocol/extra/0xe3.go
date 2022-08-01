package extra

import (
	"encoding/binary"
	"github.com/liveej/go808/errors"
)

// 里程
type Extra_0xe3 struct {
	serialized []byte
	value      float32
}

func NewExtra_0xe3(val float32) *Extra_0xe3 {
	extra := Extra_0xe3{
		value: val,
	}

	var temp [4]byte
	binary.BigEndian.PutUint16(temp[:4], uint16(val))
	extra.serialized = temp[:4]
	return &extra
}

func (Extra_0xe3) ID() byte {
	return byte(TypeExtra_0xe3)
}

func (extra Extra_0xe3) Data() []byte {
	return extra.serialized
}

func (extra Extra_0xe3) Value() interface{} {
	return extra.value
}

func (extra *Extra_0xe3) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value = float32(binary.BigEndian.Uint16(data[2:4])) / 100
	return 6, nil
}
