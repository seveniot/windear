package codec

import (
	"github.com/SevenIOT/windear/ex"
	"github.com/SevenIOT/windear/mqtt/packet"
)

func DecodeCONN(p *Packet) (*packet.CONNECT, error) {
	if p.RemainLength < 7 {
		return nil, ex.IncorrectRemainLength
	}

	var protocolName string
	var err error

	protocolName, p.Content, err = DecodeUtf8String(p.Content)

	if err != nil {
		return nil, err
	}

	if protocolName != "MQTT" {
		return nil, ex.UnsupportedMqttVersion
	}

	//if p.Content[2] != 'M' || p.Content[3] != 'Q' || p.Content[4] != 'T' || p.Content[5] != 'T' {
	//	return nil,exception.MqttIdentifierNotFoundException
	//}

	//MQTT version 3.1.1
	if p.Content[0] != 4 {
		return nil, ex.UnsupportedMqttVersion
	}

	connectFlag := packet.ConnectFlags(p.Content[1])
	//fmt.Println("username:", connectFlag.UserNameFlag())
	//fmt.Println("password:", connectFlag.PasswordFlag())
	//fmt.Println("willretain:", connectFlag.WillRetain())
	//fmt.Println("willqos:", connectFlag.WillQoS())
	//fmt.Println("willflag:", connectFlag.WillFlag())
	//fmt.Println("cleansession:", connectFlag.CleanSession())
	//fmt.Println("reserved:", connectFlag.Reserved())
	//var keepalive_msb = p.Content[8]
	//var keepalive_lsb = p.Content[9]
	//
	//var keepalive = uint16(keepalive_msb)*256 + uint16(keepalive_lsb)
	//
	//fmt.Println("keepalive:", keepalive)

	connPacket := &packet.CONNECT{
		Keepalive:    uint16(p.Content[2])*256 + uint16(p.Content[3]),
		WillRetain:   connectFlag.WillRetain(),
		WillFlag:     connectFlag.WillFlag(),
		WillQos:      connectFlag.WillQoS(),
		CleanSession: connectFlag.CleanSession(),
	}

	p.Content = p.Content[4:]

	//var len = uint16(p.Content[10])<<8 + uint16(p.Content[11])
	//
	//payload := p.Content[12:]
	//
	//if p.RemainLength < reLen {
	//	return nil,IllegalPayloadDataException
	//}

	connPacket.ClientId, p.Content, err = DecodeUtf8String(p.Content) //string(p.Content[offset : reLen])
	if err != nil {
		return nil, err
	}

	if connPacket.WillFlag {
		connPacket.WillTopic, p.Content, err = DecodeUtf8String(p.Content)

		if err != nil {
			return nil, err
		}

		connPacket.WillMsg, p.Content, err = DecodeUtf8String(p.Content)

		if err != nil {
			return nil, err
		}
	}

	if connectFlag.UserNameFlag() {
		connPacket.UserName, p.Content, err = DecodeUtf8String(p.Content)

		if err != nil {
			return nil, err
		}
	}

	if connectFlag.PasswordFlag() {
		connPacket.Password, p.Content, err = DecodeUtf8String(p.Content)

		if err != nil {
			return nil, err
		}
	}

	return connPacket, nil
}

func EncodeCONN(p *packet.CONNECT) []byte {
	var buffer []byte

	buffer = append(buffer, uint8(0x10))

	remainBuffer := []byte{uint8(0), uint8(4), 'M', 'Q', 'T', 'T', uint8(4)} //protocol name MQTT v3.1.1
	remainBuffer = append(remainBuffer, p.GetConnectFlag())                  //connect flag
	remainBuffer = append(remainBuffer, uint8(p.Keepalive>>8), uint8(p.Keepalive<<8))
	remainBuffer = append(remainBuffer, EncodeUtf8String(p.ClientId)...)
	remainBuffer = append(remainBuffer, EncodeUtf8String(p.UserName)...)
	remainBuffer = append(remainBuffer, EncodeUtf8String(p.Password)...)

	reLen := EncodeRemainLength(uint32(len(remainBuffer)))

	buffer = append(buffer, reLen...)
	buffer = append(buffer, remainBuffer...)

	return buffer
}
