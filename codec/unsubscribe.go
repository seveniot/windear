package codec

import (
	"github.com/SevenIOT/windear/ex"
	"github.com/SevenIOT/windear/mqtt/packet"
)

func DecodeUnSubscribe(p *Packet) (*packet.UNSUBSCRIBE, error) {
	unsubscribePacket := &packet.UNSUBSCRIBE{}

	unsubscribePacket.PacketIdentifier = uint16(p.Content[0])<<8 + uint16(p.Content[1])

	p.Content = p.Content[2:]

	var err error

	unsubscribePacket.TopicArray, err = decodeUnSubTopic(p.Content)

	if err != nil {
		return nil, err
	}

	return unsubscribePacket, nil
}

func decodeUnSubTopic(content []byte) ([]string, error) {
	topicArray := make([]string, 0)

	len, offset := len(content), 0

	for offset+2 < len {
		l := int(content[offset])<<8 + int(content[offset+1])

		if 2+l > len {
			return nil, ex.IllegalPayloadData
		}

		topicArray = append(topicArray, string(content[offset+2:offset+2+l]))

		offset = offset + 2 + l
	}

	if offset == len {
		return topicArray, nil
	}

	return nil, ex.IllegalPayloadData
}
