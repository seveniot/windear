package codec

import (
	"github.com/SevenIOT/windear/mqtt/packet"
	"github.com/SevenIOT/windear/ex"
	"github.com/SevenIOT/windear/mqtt/packet/types"
)

func DecodePublish(p *Packet)(*packet.PUBLISH,error){
	publishPacket := &packet.PUBLISH{
		DupFlag:p.Flag&0x08>0,
		QosLevel:uint8(p.Flag&0x06)>>1,
		Retain:p.Flag&0x01>0,
	}

	if publishPacket.QosLevel>2{
		return nil,ex.UnsupportedMqttVersion
	}

	var err error

	publishPacket.Topic,p.Content,err = DecodeUtf8String(p.Content)

	if err!=nil{
		return nil,err
	}

	if publishPacket.QosLevel>0{
		if len(p.Content)<2{
			return nil,ex.LackOfMessageIdentifier
		}

		publishPacket.PacketIdentifier = uint16(p.Content[0])<<8+uint16(p.Content[1])
		p.Content = p.Content[2:]
	}

	publishPacket.Payload = p.Content

	return publishPacket,nil
}

func EncodePublish(p *packet.PUBLISH)[]byte{
	var buffer []byte

	buffer = append(buffer,)

	var dupFlag,retainFlag uint8 = 0,0
	if p.DupFlag{
		dupFlag = 1
	}

	if p.Retain{
		retainFlag = 1
	}

	//fixedHeader
	buffer = append(buffer,types.PUBLISH<<4|dupFlag<<3|p.QosLevel<<1|retainFlag)

	var topicLen = len(p.Topic)

	remainBuffer := []byte{uint8(topicLen>>8),uint8(topicLen&0xFF)}
	remainBuffer = append(remainBuffer,[]byte(p.Topic)...)

	if p.QosLevel>0{
		remainBuffer = append(remainBuffer,uint8(p.PacketIdentifier>>8),uint8(p.PacketIdentifier&0xFF))
	}

	remainBuffer = append(remainBuffer,p.Payload...)

	reLen := EncodeRemainLength(uint32(len(remainBuffer)))

	buffer = append(buffer,reLen...)
	buffer = append(buffer,remainBuffer...)

	return buffer
}