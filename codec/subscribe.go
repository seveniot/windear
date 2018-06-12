package codec

import (
	"github.com/SevenIOT/windear/mqtt/packet"
	"github.com/SevenIOT/windear/ex"
)

func DecodeSubscribe(p *Packet)(*packet.SUBSCRIBE,error){
	subscribePacket := &packet.SUBSCRIBE{}

	subscribePacket.PacketIdentifier = uint16(p.Content[0])<<8+uint16(p.Content[1])

	p.Content = p.Content[2:]

	var err error

	subscribePacket.TopicArray,err = decodeSubTopic(p.Content)

	if err!=nil{
		return nil,err
	}

	return subscribePacket,nil
}

func decodeSubTopic(content []byte) ([]packet.SubTopic,error){
	topicArray := make([]packet.SubTopic,0)

	len,offset := len(content),0

	for offset+2<len{
		l := int(content[offset])<<8+int(content[offset+1])

		if 3+l>len{
			return nil,ex.IllegalPayloadData
		}

		topicArray = append(topicArray,packet.SubTopic{Topic:string(content[offset+2:offset+2+l]),Qos:content[offset+2+l]})

		offset = offset+3+l
	}

	if offset==len{
		return topicArray,nil
	}

	return nil,ex.IllegalPayloadData
}
