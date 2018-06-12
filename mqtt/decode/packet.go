package decode

import (
	"io"
	"github.com/SevenIOT/windear/mqtt/packet/types"
)

type Packet struct{
	PacketType uint8
	Flag uint8
	RemainLength uint32
	Content []byte
}

//func ParsePacket(reader io.Reader) (*Packet, exception.Exception){
//	var buf = make([]byte,1)
//
//	len,err := reader.Read(buf)
//
//	if err!=nil{
//		return nil,exception.New(err.Error())
//	}
//
//	if len==0{
//		return nil,exception.EmptyPacketException
//	}
//
//	p := &Packet{}
//
//	p.PacketType = buf[0]>>4
//	p.Flag = buf[0]&0x0F
//
//	//illegal packet type
//	if p.PacketType >15{
//		return nil,exception.MalformedPacketType
//	}
//
//	if p.PacketType != types.PUBLISH&&p.Flag != getDefaultFlag(p.PacketType){
//		//err = errors.New(fmt.Sprintf("malformed mqtt packet, type:%v, flag:%v",p.PacketType,p.Flag))
//		return nil,exception.MalformedPacketFlag
//	}
//
//	var ex exception.Exception
//
//	p.RemainLength,ex = RemainLen(reader)
//
//	if ex!=nil{
//		return nil,ex
//	}
//
//	p.Content,err  = getBytesByLength(reader,p.RemainLength)
//
//	if err!=nil{
//		return nil,exception.IllegalPayloadDataException
//	}
//
//	return p,nil
//}

func getDefaultFlag(packetType uint8) uint8{
	//if packetType==types.PUBREL||packetType==types.SUBSCRIBE||packetType==types.UNSUBSCRIBE{
	//	return 2
	//}
	switch packetType{
	case types.PUBREL,types.SUBSCRIBE,types.UNSUBSCRIBE:
		return 2
	}

	return 0
}

func getBytesByLength(reader io.Reader,total uint32) (buf []byte, err error){
	buf = make([]byte,total)
	var offset uint32 = 0

	for{
		var l int

		l,err = reader.Read(buf[offset:])

		if err!=nil{
			return
		}

		offset += uint32(l)

		if offset==total{
			return
		}
	}
}
