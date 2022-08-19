package protocol

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unsurper/tancy/errors"
	"reflect"
	"strconv"
)

// 消息包
type Message struct {
	Header Header
	Body   Entity
}

type DHeader struct {
	MsgID MsgID
}

// 协议编码
func (message *Message) Encode(key ...*rsa.PublicKey) ([]byte, error) {
	// 编码消息体
	count := 0
	var err error
	var body []byte
	checkSum := byte(0x00)
	if message.Body != nil && !reflect.ValueOf(message.Body).IsNil() {
		body, err = message.Body.Encode()
		if err != nil {
			return nil, err
		}

		if len(key) > 0 && key[0] != nil {
			message.Header.Property.enableEncrypt()
			body, err = EncryptOAEP(sha1.New(), key[0], body, nil)
			if err != nil {
				log.WithFields(log.Fields{
					"id":     fmt.Sprintf("0x%x", message.Header.MsgID),
					"reason": err,
				}).Warn("[JT/T808] encrypt body failed")
				return nil, err
			}
		}
	}
	checkSum, count = message.computeChecksum(body, checkSum, count)

	// 编码消息头
	message.Header.MsgID = message.Body.MsgID()
	err = message.Header.Property.SetBodySize(uint16(len(body)))
	if err != nil {
		return nil, err
	}
	header, err := message.Header.Encode()
	if err != nil {
		return nil, err
	}
	checkSum, count = message.computeChecksum(header, checkSum, count)

	// 二进制转义
	buffer := bytes.NewBuffer(nil)
	buffer.Grow(count + 2)
	buffer.WriteByte(PrefixID)
	message.write(buffer, header).write(buffer, body).write(buffer, []byte{checkSum})
	buffer.WriteByte(PrefixID)
	return buffer.Bytes(), nil
}

// 协议解码
func (message *Message) Decode(data []byte, key ...*rsa.PrivateKey) error {
	// 检验标志位
	fmt.Println(data[0], RegisterByte)
	if len(data) < 2 || (data[0] != ReceiveByte && data[0] != RegisterByte) {
		return errors.ErrInvalidMessage
	}
	if len(data) == 0 {
		return errors.ErrInvalidMessage
	}

	var header Header
	var err error

	//处理注册包
	if data[0] == RegisterByte {
		i := 2
		for ; i < len(data); i++ {
			if data[i] == IPByte {
				break
			}
		}
		IccID, err := strconv.ParseUint(string(data[2:i]), 16, 0)
		fmt.Println(IccID)
		if err != nil {
			return err
		}
		header.IccID = uint64(IccID) //用户名唯一标识码

		log.WithFields(log.Fields{
			"DTU": fmt.Sprintf("user: %s online", data[2:i]),
		}).Info("Register DTU")

		message.Header = header
		return nil
	}

	header.MsgID = MsgID(data[2]) //消息ID

	DecID, err := strconv.Atoi(bcdToString(data[3:11]))
	if err != nil {
		return err
	}
	header.DecID = uint64(DecID) //燃气表唯一标识码

	header.LocID = hex.EncodeToString(data[11:19]) //远传位置号

	IccID, err := strconv.Atoi(bcdToString(data[19:25]))
	if err != nil {
		return err
	}
	header.IccID = uint64(IccID) //用户名唯一标识码

	header.Uptime, err = fromBCDTime(data[25:31]) //打包上传时间
	if err != nil {
		return err
	}

	entity, _, err := message.decode(uint16(header.MsgID), data[31:]) //解析实体对象 entity     buffer : 为消息标识

	if err == nil {
		message.Body = entity
	} else {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%x", header.MsgID),
			"reason": err,
		}).Warn("failed to decode message")
	}
	message.Header = header
	return nil
}

//--->
func (message *Message) decode(typ uint16, data []byte) (Entity, int, error) {
	creator, ok := entityMapper[typ]
	if !ok {
		return nil, 0, errors.ErrTypeNotRegistered
	}

	entity := creator()
	entityPacket, ok := interface{}(entity).(EntityPacket)
	if !ok {
		count, err := entity.Decode(data) //解析data数据
		if err != nil {
			return nil, 0, err
		}
		return entity, count, nil
	}
	fmt.Println()
	err := entityPacket.DecodePacket(data)
	if err != nil {
		return nil, 0, err
	}
	return entityPacket, len(data), nil
}

// 写入二进制数据
func (message *Message) write(buffer *bytes.Buffer, data []byte) *Message {
	for _, b := range data {
		if b == PrefixID {
			buffer.WriteByte(EscapeByte)
			buffer.WriteByte(EscapeByteSufix2)
		} else if b == EscapeByte {
			buffer.WriteByte(EscapeByte)
			buffer.WriteByte(EscapeByteSufix1)
		} else {
			buffer.WriteByte(b)
		}
	}
	return message
}

// 校验和累加计算
func (message *Message) computeChecksum(data []byte, checkSum byte, count int) (byte, int) {
	for _, b := range data {
		checkSum = checkSum ^ b
		if b != PrefixID && b != EscapeByte {
			count++
		} else {
			count += 2
		}
	}
	return checkSum, count
}
