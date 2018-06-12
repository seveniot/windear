package core

import (
	"github.com/SevenIOT/windear/codec"
	"github.com/SevenIOT/windear/ex"
	"github.com/SevenIOT/windear/log"
	"github.com/SevenIOT/windear/mqtt/encode"
	"github.com/SevenIOT/windear/mqtt/packet"
	"github.com/SevenIOT/windear/mqtt/packet/types"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/5
 *
 */

type Handler struct {
	id   int
	load int32
	mq   chan *Task
}

func (h *Handler) Start() {
	for {
		task := <-h.mq
		h.process(task)
	}
}

func (h *Handler) process(task *Task) {
	var err error

	defer func() {
		if err != nil {
			task.Ch.Close()

			log.Infof("task process error:%v, clientId:%v, addr:%v", err.Error(), task.Ch.Ctx.ClientId, task.Ch.Ctx.Conn.RemoteAddr().String())
		}
	}()

	switch task.Packet.PacketType {
	case types.DISCONNECT:
		err = ex.ClientDisconnect
		return
	case types.PUBLISH:
		publishPacket, err := codec.DecodePublish(task.Packet)

		if err != nil {
			return
		}

		log.Infof("client publish,topic:%v,qos:%v,identifier:%v,channel:{%v}", publishPacket.Topic, publishPacket.QosLevel, publishPacket.PacketIdentifier, task.Ch)

		if publishPacket.QosLevel > 0 {
			task.Ch.Ctx.Conn.Write([]byte{types.PUBACK << 4, 0x02, uint8(publishPacket.PacketIdentifier >> 8), uint8(publishPacket.PacketIdentifier & 0xFF)})
			//dispatcher.handler.OnPublish(publishPacket,ctx)
			//if publishPacket.QosLevel == 1{
			//	dispatcher.handler.OnPublish(publishPacket,ctx)
			//}else if publishPacket.QosLevel == 2 {
			//	dispatcher.msgIdentifierCache[publishPacket.PacketIdentifier] = publishPacket
			//}
		}

		task.Ch.Publish(publishPacket)

		//ch.serverContext.PubTopic(publishPacket.Topic,encode.Publish(publishPacket))
	case types.PUBACK:
		packetIdentifier := task.Packet.Content[0]<<8 + task.Packet.Content[1]
		log.Printf("client puback,packet identifier:%v,channel:{%v}", packetIdentifier, task.Ch)
	case types.PUBREC:
		packetIdentifier := task.Packet.Content[0]<<8 + task.Packet.Content[1]
		log.Printf("client pubrec,packet identifier:%v,channel:{%v}", packetIdentifier, task.Ch)

		task.Ch.Ctx.Conn.Write([]byte{types.PUBREL<<4 | 0x02, 0x02, uint8(packetIdentifier >> 8), uint8(packetIdentifier & 0xFF)})
	case types.PUBCOMP:
		packetIdentifier := task.Packet.Content[0]<<8 + task.Packet.Content[1]
		log.Printf("client pubcomp,packet identifier:%v,channel:{%v}", packetIdentifier, task.Ch)
	case types.SUBSCRIBE:
		var subscribePacket *packet.SUBSCRIBE
		subscribePacket, err = codec.DecodeSubscribe(task.Packet)

		if err != nil {
			return
		}

		log.Infof("client subscribe,topics:%v,channel:{%v}", subscribePacket.TopicArray, task.Ch)

		var buffer []byte
		reLen := encode.RemainLength(uint32(2 + len(subscribePacket.TopicArray)))
		buffer = append(buffer, types.PUBACK<<4)
		buffer = append(buffer, reLen...)
		buffer = append(buffer, uint8(subscribePacket.PacketIdentifier>>8))
		buffer = append(buffer, uint8(subscribePacket.PacketIdentifier&0x0F))

		for _, topic := range subscribePacket.TopicArray {
			task.Ch.SubTopic(subscribePacket.TopicArray[0].Topic, task.Ch.Ctx.ClientId)

			//if err!=nil{
			//	buffer = append(buffer,0x80)
			//	log.Errorf("client subscribe error, clientId:%v, topic:%v, error:%v",task.Ch.Ctx.ClientId, task.Ch.Ctx.Conn.RemoteAddr().String(),topic,err.Error())
			//}else{
			//	buffer = append(buffer,topic.Qos)
			//}

			buffer = append(buffer, topic.Qos)
		}

		task.Ch.Ctx.Conn.Write(buffer)

		//pubPack := packet.PUBLISH{QosLevel:2,Topic:subscribePacket.TopicArray[0].Topic,Payload:[]byte("your sub was accepted"),PacketIdentifier:168}
		//ch.response(encode.Publish(&pubPack))
	case types.UNSUBSCRIBE:
		var unsubscribePacket *packet.UNSUBSCRIBE
		unsubscribePacket, err = codec.DecodeUnSubscribe(task.Packet)

		if err != nil {
			return
		}

		log.Printf("client unsubscribe,topics:%v,channel:{%v}", unsubscribePacket.TopicArray, task.Ch)

		task.Ch.Ctx.Conn.Write([]byte{types.UNSUBACK << 4, 0x02, uint8(unsubscribePacket.PacketIdentifier >> 8), uint8(unsubscribePacket.PacketIdentifier & 0xFF)})
	case types.PINGREQ:
		log.Printf("client ping request,channel:{%v}", task.Ch)
		task.Ch.Ctx.Conn.Write([]byte{types.PINGRESP << 4, 0})

		//ch.response(encode.Publish(&packet.PUBLISH{DupFlag:false,QosLevel:1,Retain:false,Topic:"hello",PacketIdentifier:1,Payload:[]byte("hello world")}))
	default:
		log.Printf("msg received,type:%v,payload:%v,channel:{%v}", task.Packet.PacketType, string(task.Packet.Content), task.Ch)
	}
}
