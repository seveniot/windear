package packet

type PUBLISH struct{
	DupFlag bool
	QosLevel uint8
	Retain bool
	Topic string
	PacketIdentifier uint16
	Payload []byte
}
