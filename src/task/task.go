package task

import (
	"encoding/json"
	"sync"
)

var RunnerID int64

const (
	Load_task = iota
)

type TaskType int64

func GetTypeTask(t TaskType) Task {
	switch t {
	case Load_task:
		return &LoadTask{}
	}
	return nil
}

type Task interface {
	Run()
	Split(Num int) []Task            // split to Num subTask and submit to taskServer
	ToJson() ([]byte, error)         // need to add RunnerID to json
	FromJson(json_data []byte) error // need to add RunnerID to json
	SetTaskID(ID int64)
	GetTaskID() int64
	GetTaskType() TaskType
}

type OriginTask struct {
	WG     sync.WaitGroup
	Task   Task
	TaskID int64
}

func NewOriginTask(task Task, TaskID int64) *OriginTask {
	return &OriginTask{Task: task, TaskID: TaskID}
}

type Brand struct {
	pay         float64 //总销售额
	number      int64   //零件个数
	sum         int64   //销售量
	notgood_sum int64   //非畅销销售量
	notgood_pay float64 //非畅销销售额
}

type BasicResult struct {
	Answer      [60]Brand
	All_sum     int64
	Notgood_sum int64
	All_pay     float64
	Notgood_pay float64
}

type Result_t struct {
	Result BasicResult
	lock   sync.Mutex
}

func (brand *Brand) Merge(vic *Brand) {
	brand.pay += vic.pay
	brand.number += vic.number
	brand.sum += vic.sum
	brand.notgood_pay += vic.notgood_pay
	brand.notgood_sum += vic.notgood_sum
}

func (res *Result_t) Merge(result *BasicResult) {
	res.lock.Lock()
	defer res.lock.Unlock()
	for i, value := range result.Answer {
		res.Result.Answer[i].Merge(&value)
	}
}

func (res *BasicResult) ToJson() ([]byte, error) {
	result, err := json.Marshal(&res.Answer)
	return result, err
}

func (res *BasicResult) FromJson(data []byte) error {
	err := json.Unmarshal(data, &res.Answer)
	return err
}

var Result Result_t
