package protocol

type Tancy_0x0017 struct {
	//信号强度bcd
	Signal uint32
}

func (entity *Tancy_0x0017) MsgID() MsgID {
	return Msgtancy_0x0017
}

func (entity *Tancy_0x0017) Encode() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (entity *Tancy_0x0017) Decode(data []byte) (int, error) {

	return 0, nil
}
