package packet

type SubTopic struct {
	Topic string
	Qos   uint8
}

type SUBSCRIBE struct {
	PacketIdentifier uint16
	TopicArray       []SubTopic
}
