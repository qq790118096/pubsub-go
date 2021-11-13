package pubsub

import (
	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

type (
	subscriber struct {
		Topic string           // 订阅主题
		Name  string           // 订阅者名称
		Msg   chan interface{} // 消息通道
	}
	publisher struct {
		mx       sync.Mutex
		topicMsg map[string][]chan interface{} // 对映topic下面可能有多个subscriber
	}
)

func NewPublisher() *publisher {
	return &publisher{
		topicMsg: make(map[string][]chan interface{}),
	}
}

func NewSubscriber(name string, t interface{}) *subscriber {
	topic := ""

	if t == nil {
		topic = ""
	} else {
		topic = t.(string)
	}

	s := &subscriber{
		Topic: topic,
		Name:  name,
		Msg:   make(chan interface{}),
	}

	go func() {
		for {
			fmt.Printf("subscriber %s: %s\n", s.Name, <-s.Msg)
		}
	}()

	return s
}

func (p *publisher) Regiser(s *subscriber) {
	p.mx.Lock()
	defer p.mx.Unlock()

	if len(p.topicMsg[s.Topic]) == 0 {
		p.topicMsg[s.Topic] = make([]chan interface{}, 0, 10)
		p.topicMsg[s.Topic] = append(p.topicMsg[s.Topic], s.Msg)
	} else {
		p.topicMsg[s.Topic] = append(p.topicMsg[s.Topic], s.Msg)
	}
}

func (p *publisher) Publish(msg, t interface{}) {
	p.mx.Lock()
	defer func() {
		p.mx.Unlock()
	}()

	if t == nil || t == "" {
		wg.Add(len(p.topicMsg))
		for _, msgChanList := range p.topicMsg {
			mcl := msgChanList
			go func() {
				for _, ch := range mcl {
					ch <- msg
				}
			}()
			wg.Done()
		}
		wg.Wait()
	} else {
		topic := t.(string)
		if msgChanList := p.topicMsg[topic]; len(msgChanList) > 0 {
			for _, msgChan := range msgChanList {
				msgChan <- msg
			}
		}

	}
}
