package main

/**
TODO:
监控queue当前长度，动态伸缩worker

**/

/*
curl -H "Content-Type:application/json" -XPOST -d'{
	 "timestamp": 12233445533,
	 "data": [
			 {
					 "user_id": 1,
					 "msg": "http://www.google.com"
			 },
			 {
					 "user_id": 2,
					 "msg": "http://www.baidu.com"
			 },
			 {
					 "user_id": 3,
					 "msg": "http://www.baidu.com"
			 },
			 {
					 "user_id": 4,
					 "msg": "http://www.baidu.com"
			 },
			 {
					 "user_id": 5,
					 "msg": "http://www.baidu.com"
				 }]}'  http://127.0.0.1:8400/sendMsg?timeout=2&asyn=1


 go-torch -u http://127.0.0.1:8401 -t 30

 ./go-wrk/go-wrk -c 10 -t=10 -n=10 -b='{
    "timestamp": 12233445533,
    "data": [
        {
            "user_id": 1,
            "msg": "http://www.google.com"
        },
        {
            "user_id": 2,
            "msg": "http://www.baidu.com"
        },
        {
            "user_id": 3,
            "msg": "http://www.baidu.com"
        },
        {
            "user_id": 4,
            "msg": "http://www.baidu.com"
        },
        {
            "user_id": 5,
            "msg": "http://www.baidu.com"
          }]}' -m="POST"  http://127.0.0.1:8400/sendMsg?timeout=1

*/

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	InitWorker   = 4
	MaxWorker    = 100
	MaxLoadWorks = 10 //每个工人最大负荷量，超出自动加人
	MaxQueue     = 20
)
var SendMsg uint64 = 0
var GetMsg uint64 = 0

type MsgCollection struct {
	Timestamp int   `json:"timestamp"`
	Msgs      []Msg `json:"data"`
}

type Msg struct {
	UserId  int    `json:"user_id"`
	Content string `json:"msg"`
}

type Job struct {
	ID      uint64
	Msg     Msg
	Wait    *sync.WaitGroup
	TimeOut int
}

var JobQueue chan Job

type Worker struct {
	ID         int
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(id int, workerPool chan chan Job) Worker {
	return Worker{
		ID:         id,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

type Dispatcher struct {
	WorkerPool   chan chan Job
	WorkerNumber int
	InitWorker   int
	MaxWorker    int
	WorkerGroup  []Worker
}

func NewDispatcher(MaxWorkers int, initWorkers int) *Dispatcher {
	pool := make(chan chan Job, MaxWorkers)
	//group := [MaxWorker]Worker{}
	group := make([]Worker, MaxWorkers)
	return &Dispatcher{WorkerPool: pool, WorkerNumber: 0, MaxWorker: MaxWorkers, InitWorker: initWorkers, WorkerGroup: group}
}

func initQueue() {
	JobQueue = make(chan Job, MaxQueue) //定义接收队列最大容量
	dispatcher := NewDispatcher(MaxWorker, InitWorker)
	dispatcher.Run()
	log.Println("initqueue finish")
}

func sendMsg(w http.ResponseWriter, r *http.Request) string {
	timeout, _ := strconv.Atoi(r.FormValue("timeout"))
	timeout = 1
	//默认同步
	asyn, _ := strconv.Atoi(r.FormValue("asyn"))
	log.Println("timeout:", timeout)
	log.Println("asyn:", asyn)
	//通过限制go 数量，发起阻塞
	log.Printf("Current NumGoroutine: %d\n", runtime.NumGoroutine())
	if runtime.NumGoroutine() > 1000 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return "server is busy,please wait"
	}
	var content = &MsgCollection{}
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &content)

	if err != nil {
		log.Println("parse error:", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(503)
		return "failed"
	}

	var waitgroup = new(sync.WaitGroup)

	contentString, _ := json.Marshal(content)
	//
	log.Println(string(contentString))
	for _, msg := range content.Msgs {
		atomic.AddUint64(&SendMsg, 1)

		waitgroup.Add(1)

		work := Job{ID: SendMsg, Msg: msg, TimeOut: timeout, Wait: waitgroup}
		JobQueue <- work
	}
	if asyn == 0 {
		waitgroup.Wait()
	}
	w.WriteHeader(http.StatusOK)

	return "ok"
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			log.Printf("worker: %d ready", w.ID)
			select {
			case job := <-w.JobChannel:
				atomic.AddUint64(&GetMsg, 1)
				//模拟任务处理时间
				//timeChan:=<-time.After(10 * time.Second):
				time.Sleep(time.Second * time.Duration(job.TimeOut))
				log.Println("worker:", w.ID, "| jobID:", job.ID, "|content:", job.Msg.Content)
				job.Wait.Done()
				log.Printf("sendMsg:%d, GetMsg: %d ", atomic.LoadUint64(&SendMsg), atomic.LoadUint64(&GetMsg))
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorker; i++ {
		d.WorkerGroup[i] = NewWorker(i, d.WorkerPool)
		if d.WorkerNumber < d.InitWorker {
			d.WorkerGroup[i].Start()
			d.WorkerNumber++
		}
		//	worker.Stop()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	//initGoroutine := runtime.NumGoroutine()
	log.Printf("init NumGoroutine: %d\n", runtime.NumGoroutine())
	for {

		select {
		case job := <-JobQueue:

			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
				log.Printf("dispatch job:%d", job.ID)
			}(job)
			//根据当前等待消息数 去增加worker

			if SendMsg-GetMsg > uint64(d.WorkerNumber*MaxLoadWorks) && d.WorkerNumber < d.MaxWorker {
				log.Printf("overload works,add workers,waitMsg: %d", SendMsg-GetMsg)
				d.WorkerGroup[d.WorkerNumber].Start()
				d.WorkerNumber++
				log.Printf("worker num:%d", d.WorkerNumber)
			}
		}
		if SendMsg-GetMsg < uint64(d.InitWorker*MaxLoadWorks) && d.WorkerNumber > d.InitWorker {
			log.Printf("decrease workernumber")
			for {
				log.Printf("stop worker:%d", d.WorkerNumber)
				d.WorkerGroup[d.WorkerNumber-1].Stop()
				d.WorkerNumber--
				if d.WorkerNumber == d.InitWorker {
					break
				}
			}
		}

	}
}
