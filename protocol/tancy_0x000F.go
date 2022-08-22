package protocol

import "time"

type Tancy_0x000F struct {
	//上传时间
	Uptime time.Time

	//标况瞬时流量
	OpmomFlow uint32
	//标况bcd
	OpBCD uint16
	//标况累积流量
	OpcumFlow uint32

	//工况瞬时流量
	WomomFlow uint32
	//工况bcd
	WoBCD uint16
	//工况累积流量
	WocumFlow uint32

	//燃气温度
	TGT uint32
	//燃气压力
	TGP uint32

	//转换系数
	ConversionFactor uint32

	//报警
	AlarmWord uint16

	//状态码
	State byte
}

func (entity *Tancy_0x000F) MsgID() MsgID {
	return Msgtancy_0x000F
}

func (entity *Tancy_0x000F) Encode() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (entity *Tancy_0x000F) Decode(data []byte) (int, error) {
	datalen := len(data)
	if datalen < 44 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	var err error

	entity.Uptime, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	// 标况瞬时流量
	entity.OpmomFlow, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	// 标况BCD
	entity.OpBCD, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}
	// 标况累积流量
	entity.OpcumFlow, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 工况瞬时流量
	entity.WomomFlow, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	// 标况BCD
	entity.WoBCD, err = reader.ReadUint16()
	// 工况累积流量
	entity.WocumFlow, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 燃气温度
	entity.TGT, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	// 燃气压力
	entity.TGP, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	entity.ConversionFactor, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	entity.AlarmWord, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.State, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
