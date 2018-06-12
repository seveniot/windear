package packet

type UNSUBSCRIBE struct{
	PacketIdentifier uint16
	TopicArray []string
}
