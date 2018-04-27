package main

/**
TODO:
监控queue当前长度，动态伸缩worker

**/

/*curl -XPOST -d '{
    "timestamp": 12233445533,"data": [{ "user_id": 1,"msg": "http://www.google.com" },
        {
            "user_id": 2,
            "msg": "http://www.baidu.com"
        }
    ]
}'  http://127.0.0.1:8400/sendMsg
 go-torch -u http://127.0.0.1:8401 -t 30
 go-wrk -c=400 -t=8 -n=100000 -m="POST" -b='xxxxxx'   http://127.0.0.1:8400/sendMsg

*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
    initWorker=4
	MaxWorker =100 
	MaxQueue  = 20
	SendMsg   = 0
	GetMsg    = 0
)

type MsgCollection struct {
	Timestamp int   `json:"timestamp"`
	Msgs      []Msg `json:"data"`
}

type Msg struct {
	UserId  int    `json:"user_id"`
	Content string `json:"msg"`
}

type Job struct {
	ID      int
	Msg     Msg
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
}

func NewDispatcher(maxWorkers int,initWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, WorkerNumber: initWorkers}
}

func initQueue() {
	JobQueue = make(chan Job, MaxQueue) //定义接收队列最大容量
	dispatcher := NewDispatcher(MaxWorker,initWorker)
	dispatcher.Run()
	log.Println("initqueue finish")
}

func sendMsg(w http.ResponseWriter, r *http.Request) string {
	timeout, _ := strconv.Atoi(r.FormValue("timeout"))
	log.Println("timeout:", timeout)
	var content = &MsgCollection{}
	body, _ := ioutil.ReadAll(r.Body)

	//fmt.Println(string(body))

	err := json.Unmarshal(body, &content)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return "failed"
	}
	contentString, _ := json.Marshal(content)
	fmt.Println(string(contentString))
	for _, msg := range content.Msgs {
		SendMsg++
		work := Job{ID: SendMsg, Msg: msg, TimeOut: timeout}
		JobQueue <- work
	}

	w.WriteHeader(http.StatusOK)
	log.Println("sendMsg:" + strconv.Itoa(SendMsg) + "  GetMsg:" + strconv.Itoa(GetMsg))
	return "ok"
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			log.Printf("worker: %d ready", w.ID)
			select {
			case job := <-w.JobChannel:
				GetMsg++
				//timeChan:=<-time.After(10 * time.Second):
				time.Sleep(time.Second * time.Duration(job.TimeOut))
				log.Println("worker:", w.ID, "| jobID:", job.ID, "|content:", job.Msg.Content)

				log.Printf("leave num: %d", SendMsg-GetMsg)
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
	for i := 0; i < d.WorkerNumber; i++ {
		worker := NewWorker(i+1, d.WorkerPool)
		worker.Start()
		//	worker.Stop()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			//通过启动无数goruntine 来接收数据，可能存在蹦
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
				log.Printf("dispatch job:%d", job.ID)
			}(job)
		}
	}
}
