package core

import (
	"github.com/SevenIOT/windear/codec"
	"github.com/SevenIOT/windear/mqtt/packet"
	"net"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/4
 *
 */

type ServerContext interface {
	RegisterChannel(clientId string, ch *Channel)
	UnRegisterChannel(clientId string)
	Auth(clientId, userName, password string) bool
	SaveSession(clientId, userName string, cleanSession bool) error
	DelSession(clientId string)
	OnRcv(packet *codec.Packet, ch *Channel)
	OnPub(packet *packet.PUBLISH)
	PubMsg(clientId string, content []byte)
	OnSubTopic(topic, clientId string)

	KickClient(clientId string)
}

type ChannelContext struct {
	Conn      net.Conn
	ClientId  string
	UserName  string
	serverCtx ServerContext
	Keepalive uint16
}

func NewContext(conn net.Conn, clientId, userName string) *ChannelContext {
	return &ChannelContext{
		Conn:     conn,
		ClientId: clientId,
		UserName: userName,
	}
}
