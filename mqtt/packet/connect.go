package packet

type CONNECT struct {
	WillRetain   bool
	WillFlag     bool
	CleanSession bool
	WillQos      uint8
	Keepalive    uint16
	ClientId     string
	UserName     string
	Password     string
	WillTopic    string
	WillMsg      string
}

type CONNACK []byte

var (
	CONNACK_ACCEPT_WITH_SESSION                  = CONNACK{0x20, 0x02, 0x01, 0x00}
	CONNACK_ACCEPT_WITHOUT_SESSION               = CONNACK{0x20, 0x02, 0x00, 0x00}
	CONNACK_REFUSE_UNACCEPTABLE_PROTOCOL_VERSION = CONNACK{0x20, 0x02, 0x00, 0x01}
	CONNACK_REFUSE_ID_REJECTED                   = CONNACK{0x20, 0x02, 0x00, 0x02}
	CONNACK_REFUSE_SERVICE_UNAVAILABLE           = CONNACK{0x20, 0x02, 0x00, 0x03}
	CONNACK_REFUSE_AUTH_FAIL                     = CONNACK{0x20, 0x02, 0x00, 0x04}
	CONNACK_REFUSE_NO_AUTH_FOUND                 = CONNACK{0x20, 0x02, 0x00, 0x05}
)

type ConnectFlags uint8

func (flag ConnectFlags) UserNameFlag() bool {
	return flag&0x80 > 0
}

func (flag ConnectFlags) PasswordFlag() bool {
	return flag&0x40 > 0
}

func (flag ConnectFlags) WillRetain() bool {
	return flag&0x20 > 0
}

func (flag ConnectFlags) WillQoS() uint8 {
	return uint8(flag&0x18) >> 3
}

func (flag ConnectFlags) WillFlag() bool {
	return flag&0x04 > 0
}

func (flag ConnectFlags) CleanSession() bool {
	return flag&0x02 > 0
}

func (c *CONNECT) GetConnectFlag() uint8 {
	b := uint8(0x80 | 0x40)

	if c.WillRetain {
		b |= 0x20
	}

	b |= uint8(c.WillQos << 3)

	if c.WillFlag {
		b |= 0x04
	}

	if c.CleanSession {
		b |= 0x02
	}

	return b
}
