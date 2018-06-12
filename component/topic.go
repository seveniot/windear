package component

import (
	"github.com/garyburd/redigo/redis"
	"errors"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/5
 * 
*/

type TopicComponent struct{
	redis *redis.Pool
	expiration int
	rpcAddr string
}

func NewTopicComponent(addr string,redisPool *redis.Pool)*TopicComponent{
	return &TopicComponent{
		redis:redisPool,
		expiration:300,
		rpcAddr:addr,
	}
}

func(t *TopicComponent) SubTopic(topic, clientId string)error{
	conn := t.redis.Get()
	defer conn.Close()

	_,err := conn.Do("SADD",topic,t.rpcAddr+"#"+clientId)

	if err!=nil{
		return err
	}

	return nil
}

func(t *TopicComponent)RetrieveSubscribers(topic string)([]string,error){
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
