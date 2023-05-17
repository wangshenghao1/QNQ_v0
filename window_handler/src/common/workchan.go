package common

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var timeoutValue = 2 * time.Minute
var workerCache = make(map[string]*QWorker)

type Producer interface {
	//Produce Send the message to chan
	Produce(data interface{})
	//Stop  producer
	Stop()
	//StartPump start pump
	StartPump()
	//AddConsumer Add the consumer
	AddConsumer(worker *QWorker)

	RemoveConsumer(sn string)
}

var qmqWaitGroup *sync.WaitGroup

type QProducer struct {
	msg  chan interface{}
	sub  map[string]chan interface{} //key:{SN} value:{worker}
	wg   *sync.WaitGroup
	lock sync.RWMutex
}

// AddConsumer TODO Worker内存释放
func (p *QProducer) AddConsumer(worker *QWorker) {
	var msg chan interface{}
	p.lock.Lock()
	if _, ok := p.sub[worker.SN]; ok {
		p.lock.Unlock()
		return
	} else {
		msg = make(chan interface{}, 1)
		p.sub[worker.SN] = msg
	}
	p.lock.Unlock()
	worker.Sub = msg
	//Listen current channel
	GetCoroutinesPool().Submit(worker.Execute)
}

// RemoveConsumer d
func (p *QProducer) RemoveConsumer(id string) {
	close(p.sub[id])
	delete(p.sub, id)
}

func (p *QProducer) Produce(data interface{}) {
	p.msg <- data
	//log.Printf("send %v to mq", data)
}

func (p *QProducer) Stop() {
	close(p.msg)
}

func (p *QProducer) StartPump() {
	log.Printf("Pumo running...")
	p.wg.Add(1)
	go func() {
		for {
			if m, ok := <-p.msg; ok {
				p.lock.RLock()
				p.distributeMsg(m)
				p.lock.RUnlock()
			} else {
				log.Printf("Pump Stop Done...")
				p.lock.RLock()
				for _, v := range p.sub {
					close(v)
				}
				p.lock.RUnlock()
				break
			}
		}
		p.wg.Done()
	}()
}

func NewProducer() Producer {
	var p Producer
	if qmqWaitGroup != nil {
		return nil
	}
	qmqWaitGroup = &sync.WaitGroup{}
	p = &QProducer{
		msg:  make(chan interface{}, 1),
		sub:  make(map[string]chan interface{}),
		wg:   qmqWaitGroup,
		lock: sync.RWMutex{},
	}
	return p
}

func (p *QProducer) distributeMsg(msg interface{}) {
	msgStr := fmt.Sprintf("%v", msg)
	msgHead := getMsgHeader(msgStr)
	// Worker start flag
	if msgHead == "0x" {
		SN := getMsgSN(msgStr, msgTaskActivation)
		taskFlag := msgStr[9:10]
		if taskFlag == TaskOverFlag {
			if workerCache[SN] != nil {
				workerCache[SN].OverChan <- 1
			}
			delete(workerCache, SN)
			go p.RemoveConsumer(SN)
		} else {
			taskType := msgStr[0:5]
			worker := WorkerFactoryMap[taskType](SN)
			worker.Status = TASK_READY
			workerCache[SN] = worker
			go p.AddConsumer(worker)
		}
		return
	}

	//TODO Worker init
	if msgHead == "01" {

		return
	}

	//Worker content
	if msgHead == "00" {
		SN := getMsgSN(msgStr, msgTaskMsg)
		if p.sub[SN] == nil || len(msgStr)-8 <= 6 {
			log.Printf("worker SN {%v} is nil， msg : %v", SN, msgStr)
		} else {
			p.sub[SN] <- msgStr[6 : len(msgStr)-8]
		}
		return
	}
}

func (w *QWorker) Execute(v ...interface{}) {
	defer w.Deconstruct()
	if w.Sub == nil {
		w.ExecuteFunc(nil, w)
		return
	}
	qmqWaitGroup.Add(1)
	defer qmqWaitGroup.Done()
	ticker := time.NewTicker(timeoutValue)
	for {
		select {
		case m := <-w.Sub:
			if m != nil {
				w.ExecuteFunc(m, w)
			}
		case <-ticker.C:
			log.Printf("Consumer %s timeout", w.SN)
			break
		case <-w.OverChan:
			log.Printf("Consumer %s over", w.SN)
			return
		}
	}
}

func getMsgHeader(msg string) string {
	if len(msg) > 2 {
		return msg[0:2]
	}
	return ""
}

const (
	msgTaskActivation = iota
	msgTaskInit
	msgTaskMsg
	msgTaskLocal
)

func getMsgSN(msg string, msgType int) string {
	switch msgType {
	case msgTaskActivation:
		return msg[5:9]
	case msgTaskInit:
		return msg[5:9]
	case msgTaskMsg:
		return msg[2:6]
	}
	return ""
}
