package extra

import (
	"fmt"
	"github.com/liveej/go808/errors"
)

// WIFI信号
type Extra_0x54 struct {
	serialized []byte
	value      Extra_0x54_Value
}

type Extra_0x54_Value struct {
	WifiInfos []WifiInfo
}

type WifiInfo struct {
	Mac  string
	RSSI int
}

func NewExtra_0x54(val byte) *Extra_0x54 {
	extra := Extra_0x54{}
	extra.serialized = []byte{val}
	return &extra
}

func (Extra_0x54) ID() byte {
	return byte(TypeExtra_0x54)
}

func (extra Extra_0x54) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x54) Value() interface{} {
	return extra.value
}

func (extra *Extra_0x54) Decode(data []byte) (int, error) {
	n := int(data[0])
	if len(data) < (n*7 + 1) {
		return 0, errors.ErrInvalidExtraLength
	}
	for i := 0; i < n; i++ {
		x := i * 7
		w := WifiInfo{
			Mac:  fmt.Sprintf("%X:%X:%X:%X:%X:%X", data[x+1], data[x+2], data[x+3], data[x+4], data[x+5], data[x+6]),
			RSSI: int(data[x+7]),
		}
		extra.value.WifiInfos = append(extra.value.WifiInfos, w)
	}

	return len(data), nil
}
