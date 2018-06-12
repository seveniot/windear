package component

import (
	"errors"
	"github.com/garyburd/redigo/redis"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/6
 * 
*/


type SessionComponent struct{
	redis *redis.Pool
	expiration int
	rpcAddr string
}

func NewSessionComponent(addr string,redisPool *redis.Pool)*SessionComponent{
	return &SessionComponent{
		redis:redisPool,
		expiration:300,
		rpcAddr:addr,
	}
}

func(s *SessionComponent) SaveSession(clientId string)error{
	conn := s.redis.Get()
	defer conn.Close()

	_,err := conn.Do("SETNX",clientId,s.rpcAddr)

	if err!=nil{
		return err
	}

	return nil
}

func(s *SessionComponent) DelSession(clientId string)error{
	conn := s.redis.Get()
	defer conn.Close()

	_,err := conn.Do("DEL",clientId)

	return err
}

func(s *SessionComponent) GetSession(clientId string)(res string,err error){
	conn := s.redis.Get()
	defer conn.Close()

	data,err := conn.Do("GET",clientId)

	if err==nil&&data!=nil{
		res = string(data.([]uint8))
	}

	return
}

func(t *SessionComponent)RetrieveSubscribers(topic string)([]string,error){
	conn := t.redis.Get()
	defer conn.Close()

	list,err := conn.Do("SMEMBERS",topic)

	if err!=nil{
		return nil,err
	}

	data,ok := list.([]interface{})

	if !ok{
		return nil,errors.New("the result type is incorrect")
	}

	var res []string

	for _,item := range data{
		raw,ok := item.([]uint8)

		if ok{
			res = append(res,string(raw))
		}
	}

	return res,nil
}