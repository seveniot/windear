package server

import (
	"net"
	"fmt"
	"github.com/SevenIOT/windear/log"
	"github.com/SevenIOT/windear/codec"
	"github.com/SevenIOT/windear/mqtt/packet"
	"github.com/SevenIOT/windear/component"
	"strings"
	"github.com/SevenIOT/windear/core"
	"strconv"
	"github.com/garyburd/redigo/redis"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/4
 * 
*/

type Server interface {
	Start()error
}

type defaultServer struct {
	chTable map[string]*core.Channel
	processor *core.Processor
	port int
	rpcPort int
	hostIP string
	rpcAddr string
	topicComponent *component.TopicComponent
	sessionComponent *component.SessionComponent
	agentComponent *component.AgentComponent
	redis *redis.Pool
}

func NewServer(mqttPort,rpcPort int,redisPool *redis.Pool,hostIp string)Server{
	return &defaultServer{
		port:mqttPort,
		rpcPort:rpcPort,
		hostIP:hostIp,
		rpcAddr:hostIp+":"+strconv.Itoa(rpcPort),
		redis:redisPool,
	}
}

func(s *defaultServer)Start()error{
	ln,err := net.Listen("tcp",fmt.Sprintf(":%d",s.port))

	if err!=nil{
		return err
	}

	s.topicComponent = component.NewTopicComponent(s.rpcAddr,s.redis)
	s.sessionComponent = component.NewSessionComponent(s.rpcAddr,s.redis)
	s.agentComponent = component.NewAgentComponent(s.rpcPort,s)

	go s.agentComponent.Start()

	s.chTable = make(map[string]*core.Channel)
	s.processor = &core.Processor{}
	s.processor.Start()

	log.Info("Server Started...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}

		//go server.onAccept(conn)
		go core.EstablishChannel(conn,s)
	}

	return nil
}

func(s *defaultServer)RegisterChannel(clientId string,channel *core.Channel){
	s.KickClient(clientId)

	s.chTable[clientId] = channel
	s.processor.BindHandler(clientId)
}

func(s *defaultServer) UnRegisterChannel(clientId string){
	s.processor.UnBindHandler(clientId)
	delete(s.chTable,clientId)
}

func(s *defaultServer)Auth(clientId,userName,password string)bool{
	return true
}

func(s *defaultServer) SaveSession(clientId,userName string,cleanSession bool)error{
	res,err := s.sessionComponent.GetSession(clientId)

	if err!=nil{
		log.Printf("error in get session by clientId:%v,errorr:%v",clientId,err.Error())
		return err
	}

	if res!=""{
		if res==s.rpcAddr{
			log.Infof("kick local client,clientId:%v",res)
			s.KickClient(clientId)
		}else{
			err = s.agentComponent.KickRemote(clientId,res)

			if err!=nil{
				log.Errorf("kick remote client error:%v,clientId:%v",err.Error(),res)
				return err
			}

			log.Infof("kick remote client success,clientId:%v",res)
		}
	}

	err = s.sessionComponent.SaveSession(clientId)

	return err
}

func(s *defaultServer)DelSession(clientId string){
	err := s.sessionComponent.DelSession(clientId)

	if err!=nil{
		log.Errorf("del session error:%v,clientId:%v",err.Error(),clientId)
	}
}

func(s *defaultServer) OnRcv(packet *codec.Packet,ch *core.Channel){
	s.processor.DispatchTask(packet,ch)
}

func(s *defaultServer) OnPub(packet *packet.PUBLISH){
	clientList,err := s.topicComponent.RetrieveSubscribers(packet.Topic)

	if err!=nil{
		log.Infof("retriveveSubscribers error on pub,error:%v",err.Error())
		return
	}

	for _,client:=range clientList{
		params := strings.Split(client,"#")
		if s.rpcAddr==params[0]{
			s.PubMsg(params[1],codec.EncodePublish(packet))
		}else{
			s.agentComponent.PubMsgRemote(params[1],codec.EncodePublish(packet),params[0])
		}
	}
}

func(s *defaultServer)PubMsg(clientId string, content []byte){
	if ch,exist := s.chTable[clientId];exist{
		ch.Ctx.Conn.Write(content)
	}
}

func(s *defaultServer) OnSubTopic(topic,clientId string){
	s.topicComponent.SubTopic(topic,clientId)
}

func(s *defaultServer)KickClient(clientId string){
	if ch,exist := s.chTable[clientId];exist{
		log.Infof("kick the former channel,clientId:%v, addr:%v",ch.Ctx.ClientId,ch.Ctx.Conn.RemoteAddr().String())
		ch.Close()
	}
}

//func(s *defaultServer) PubMsg(context.Context, *agent.Request) (*agent.Response, error){
//	return nil,nil
//}