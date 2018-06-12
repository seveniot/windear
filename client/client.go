package client

import (
	"fmt"
	"github.com/SevenIOT/windear/codec"
	"github.com/SevenIOT/windear/log"
	"github.com/SevenIOT/windear/mqtt/packet"
	"net"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/8
 *
 */

type Client interface {
	PubMsg(p *packet.PUBLISH)
}

func NewClient(hostIP, clientId, userName, password string, port int) Client {
	c := &defaultClient{
		HostIP: hostIP,
		Port:   port,
	}

	c.start(clientId, userName, password)

	return c
}

type defaultClient struct {
	HostIP string
	Port   int

	conn net.Conn
}

func (c *defaultClient) start(clientId, userName, password string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", c.HostIP, c.Port))

	if err != nil {
		log.Fatal("connect server fail,error:%v", err.Error())
	}

	conn.Write(codec.EncodeCONN(&packet.CONNECT{
		ClientId: clientId,
		UserName: userName,
		Password: password,
	}))

	c.conn = conn

	go c.recv()
}

func (c *defaultClient) recv() {
	for {
		//ch.Ctx.Conn.SetReadDeadline(time.Now().Add(time.Second*5))//time.Duration(ch.Ctx.Keepalive)))
		p, err := codec.ReadPacket(c.conn)

		if err != nil {
			c.conn.Close()
			//debug.Stack()
			log.Errorf("channel read packet error:%v", err.Error())
			return
		}

		fmt.Println(string(p.Content))
	}
}

func (c *defaultClient) PubMsg(p *packet.PUBLISH) {
	fmt.Println("pub msg")
	c.conn.Write(codec.EncodePublish(p))
}
