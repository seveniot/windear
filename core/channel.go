package core

import (
	"fmt"
	"github.com/SevenIOT/windear/codec"
	"github.com/SevenIOT/windear/ex"
	"github.com/SevenIOT/windear/log"
	"github.com/SevenIOT/windear/mqtt/encode"
	"github.com/SevenIOT/windear/mqtt/packet"
	"github.com/SevenIOT/windear/mqtt/packet/types"
	"net"
	"time"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/4
 *
 */

type Channel struct {
	Ctx     *ChannelContext
	IsAlive bool
}

func EstablishChannel(conn net.Conn, svrCtx ServerContext) {
	var p *codec.Packet
	var err error

	defer func() {
		if err != nil {
			//debug.Stack()
			log.Errorf("channel connect fail:%v, addr:%v", err.Error(), conn.RemoteAddr().String())
			conn.Close()
		}
	}()

	p, err = codec.ReadPacket(conn)

	if p.PacketType != types.CONNECT {
		return
	}

	connPacket, err := codec.DecodeCONN(p)

	if err != nil {
		return
	}

	if len(connPacket.UserName) == 0 || len(connPacket.Password) == 0 {
		conn.Write(encode.CONNACK_REFUSE_NO_AUTH_FOUND)
		err = ex.ConnNoAuthFound
		return
	}

	if !svrCtx.Auth(connPacket.ClientId, connPacket.UserName, connPacket.Password) {
		conn.Write(encode.CONNACK_REFUSE_AUTH_FAIL)
		err = ex.ConnAuthFail
		return
	}

	err = svrCtx.SaveSession(connPacket.ClientId, connPacket.UserName, connPacket.CleanSession)

	if err != nil {
		return
	}

	if connPacket.CleanSession {
		_, err = conn.Write(encode.CONNACK_ACCEPT_WITH_SESSION)
	} else {
		_, err = conn.Write(encode.CONNACK_ACCEPT_WITHOUT_SESSION)
	}

	if err != nil {
		return
	}

	chCtx := &ChannelContext{
		Conn:      conn,
		ClientId:  connPacket.ClientId,
		UserName:  connPacket.UserName,
		Keepalive: connPacket.Keepalive,
		serverCtx: svrCtx,
	}

	ch := Channel{
		Ctx: chCtx,
	}

	svrCtx.RegisterChannel(connPacket.ClientId, &ch)

	log.Infof("client connected,clientId:%v,channel:{%v}", connPacket.ClientId, &ch)

	go ch.startMsgLoop()
}

//[MQTT-3.1.2-24] If the Keep Alive value is non-zero and the Server does not receive a Control Packet from the Client within one and a half times the Keep Alive time period, it MUST disconnect the Network Connection to the Client as if the network had failed
func (ch *Channel) startMsgLoop() {
	ch.IsAlive = true

	for {
		deadLine := time.Now().Add(time.Millisecond * time.Duration(ch.Ctx.Keepalive*1500))
		log.Infof("wait for keepalive time:%v,deadLine:%v,channel:%v", ch.Ctx.Keepalive, deadLine.Format("2006-01-02 15:04:05"), ch)
		ch.Ctx.Conn.SetReadDeadline(deadLine)
		p, err := codec.ReadPacket(ch.Ctx.Conn)

		if err != nil {
			ch.Close()
			//debug.Stack()
			log.Errorf("channel read packet error:%v, channel:{%v}", err.Error(), ch)
			return
		}

		log.Infof("receive packet:%v,channel:%v", p, ch)
		ch.Ctx.serverCtx.OnRcv(p, ch)
	}
}

func (ch *Channel) Close() {
	//todo pub device offline msg
	ch.IsAlive = false
	ch.Ctx.serverCtx.UnRegisterChannel(ch.Ctx.ClientId)
	ch.Ctx.serverCtx.DelSession(ch.Ctx.ClientId)
	ch.Ctx.Conn.Close()
	return
}

func (ch *Channel) String() string {
	return fmt.Sprintf("clientId:%v,userName:%v,addr:%v", ch.Ctx.ClientId, ch.Ctx.UserName, ch.Ctx.Conn.RemoteAddr().String())
}

func (ch *Channel) Publish(pubMsg *packet.PUBLISH) {
	//ch.Ctx.serverCtx.OnRcv()
	ch.Ctx.serverCtx.OnPub(pubMsg)
}

func (ch *Channel) SubTopic(topic, clientId string) {
	ch.Ctx.serverCtx.OnSubTopic(topic, clientId)
}
