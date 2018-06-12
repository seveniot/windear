package core

import (
	"github.com/SevenIOT/windear/codec"
	"runtime"
	"sync/atomic"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/5
 *
 */

const (
	MQ_SIZE = 10
)

type Task struct {
	Ch     *Channel
	Packet *codec.Packet
}

type Processor struct {
	handlers   []*Handler
	owenership map[string]int
}

func (p *Processor) Start() {
	size := runtime.NumCPU()
	p.handlers = make([]*Handler, size)
	p.owenership = make(map[string]int)

	for i := 0; i < size; i++ {
		h := &Handler{
			id:   i,
			load: 0,
			mq:   make(chan *Task, MQ_SIZE),
		}

		p.handlers[i] = h
		go h.Start()
	}
}

func (p *Processor) BindHandler(clientId string) {
	var maxLoad int32
	assignIdx := 0

	for i, h := range p.handlers {
		if atomic.LoadInt32(&h.load) < maxLoad {
			assignIdx = i
			maxLoad = h.load
		}
	}

	p.owenership[clientId] = assignIdx
}

func (p *Processor) UnBindHandler(clientId string) {
	delete(p.owenership, clientId)
}

func (p *Processor) DispatchTask(packet *codec.Packet, ch *Channel) {
	idx := p.owenership[ch.Ctx.ClientId]
	p.handlers[idx].mq <- &Task{
		Packet: packet,
		Ch:     ch,
	}
}
