package codec

import (
	"github.com/SevenIOT/windear/ex"
	"github.com/SevenIOT/windear/mqtt/packet/types"
	"io"
)

type Packet struct {
	PacketType   uint8
	Flag         uint8
	RemainLength uint32
	Content      []byte
}

func ReadPacket(reader io.Reader) (*Packet, error) {
	var err error
	var len int

	var buf = make([]byte, 1)

	len, err = reader.Read(buf)

	if err != nil {
		return nil, err
	}

	if len == 0 {
		return nil, ex.EmptyPacket
	}

	p := &Packet{}

	p.PacketType = buf[0] >> 4
	p.Flag = buf[0] & 0x0F

	//illegal packet type
	if p.PacketType > 15 {
		return nil, ex.MalformedPacketType
	}

	if p.PacketType != types.PUBLISH && p.Flag != getDefaultFlag(p.PacketType) {
		//err = errors.New(fmt.Sprintf("malformed mqtt packet, type:%v, flag:%v",p.PacketType,p.Flag))
		return nil, ex.MalformedPacketFlag
	}

	p.RemainLength, err = DecodeRemainLen(reader)

	if err != nil {
		return nil, err
	}

	p.Content, err = getBytesByLength(reader, p.RemainLength)

	if err != nil {
		return nil, ex.IllegalPayloadData
	}

	return p, nil
}

func getDefaultFlag(packetType uint8) uint8 {
	//if packetType==types.PUBREL||packetType==types.SUBSCRIBE||packetType==types.UNSUBSCRIBE{
	//	return 2
	//}
	switch packetType {
	case types.PUBREL, types.SUBSCRIBE, types.UNSUBSCRIBE:
		return 2
	}

	return 0
}

func getBytesByLength(reader io.Reader, total uint32) (buf []byte, err error) {
	buf = make([]byte, total)
	var offset uint32 = 0

	for {
		var l int

		l, err = reader.Read(buf[offset:])

		if err != nil {
			return
		}

		offset += uint32(l)

		if offset == total {
			return
		}
	}
}
