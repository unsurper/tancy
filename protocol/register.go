package protocol

// 消息ID枚举
type MsgID uint16

const (
	//  响应间隔或小时记录打包报文
	Msgtancy_0x0008 MsgID = 0x0008
)

// 消息实体映射
var entityMapper = map[uint16]func() Entity{
	uint16(Msgtancy_0x0008): func() Entity {
		return new(tancy_0x0008)
	},
}

// 类型注册
func Register(typ uint16, creator func() Entity) {
	entityMapper[typ] = creator
}
