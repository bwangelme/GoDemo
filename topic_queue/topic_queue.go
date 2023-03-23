package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// TopicMap worker 的注册信息
type TopicMap struct {
	table map[string][]chan *Msg
	sync.Mutex
}

func (t *TopicMap) Add(topic string, ch chan *Msg) {
	t.Lock()
	defer t.Unlock()

	chans, exists := t.table[topic]
	if !exists {
		chans = make([]chan *Msg, 0)
	}
	chans = append(chans, ch)
	t.table[topic] = chans
}

func (t *TopicMap) Get(topic string) []chan *Msg {
	t.Lock()
	defer t.Unlock()

	return t.table[topic]
}

func (t *TopicMap) ForEach(callback func(chan *Msg)) {
	t.Lock()
	defer t.Unlock()

	for _, chs := range t.table {
		for _, ch := range chs {
			callback(ch)
		}
	}
}

func NewTopicMap() *TopicMap {
	return &TopicMap{
		table: make(map[string][]chan *Msg),
		Mutex: sync.Mutex{},
	}
}

// Msg MQ 消息
type Msg struct {
	topic string
	data  []byte
}

// Queue 主消息队列
var Queue chan *Msg
var topicMap *TopicMap

func Init() {
	Queue = make(chan *Msg)
	topicMap = NewTopicMap()
	InitConsumers(3)
}

// InitConsumers 初始化并注册消费者
func InitConsumers(cnt int) {
	for i := 0; i < cnt; i++ {
		consumerChan := make(chan *Msg, 0)
		topic := fmt.Sprintf("topic_%v", i%2 == 0)
		topicMap.Add(topic, consumerChan)

		go Consumer(i, consumerChan)
	}
}

// Consumer MQ 消费者
func Consumer(id int, ch chan *Msg) {
	fmt.Printf("Consumer %v Start\n", id)
	for {
		select {
		case msg, ok := <-ch:
			//localMsg := &Msg{}
			//copy(localMsg, msg)

			if !ok {
				fmt.Printf("Consumer %v exit\n", id)
				return
			}

			fmt.Printf("Consumer %v Receive msg `%v` from topic %v\n", id, string(msg.data), msg.topic)
		}
	}
}

//Clean 结束时的清理动作
func Clean() {
	close(Queue)
	CloseAllConsumer()
}

func CloseAllConsumer() {
	topicMap.ForEach(func(ch chan *Msg) {
		close(ch)
	})
}

// Dispatcher 分发主消息队列中的消息给 worker
func Dispatcher() {
	for {
		select {
		case msg, ok := <-Queue:
			if !ok {
				fmt.Printf("Dispatcher exit\n")
				return
			}

			workerChans := topicMap.Get(msg.topic)
			for _, ch := range workerChans {
				ch <- msg
			}
		}
	}
}

// Producer 每秒钟发送一个随机数字
func Producer() {
	for i := 0; i < 10; i++ {
		val := rand.Intn(10000)
		topic := fmt.Sprintf("topic_%v", val%2 == 0)
		msg := &Msg{
			topic: topic,
			data:  []byte(fmt.Sprintf("data is %d", val)),
		}
		Queue <- msg

		time.Sleep(500 * time.Millisecond)
	}
	Clean()
}

func main() {
	Init()
	go Producer()
	Dispatcher()

	// 等待消费者结束，这里的处理方式不太优雅
	time.Sleep(1 * time.Second)
}
