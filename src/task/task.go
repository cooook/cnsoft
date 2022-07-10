package task

import (
	"encoding/json"
	"fmt"
	"sync"
)

var RunnerID int64

const (
	Load_task = iota
	Select_task
)

type TaskType int64

func GetTypeTask(t TaskType) Task {
	switch t {
	case Load_task:
		return &LoadTask{}
	case Select_task:
		return &SelectTask{}
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
	MergeAnwer() bool
	AllDoneCallBack()
	Initial()
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
	Pay         float64
	Number      int64
	Sum         int64
	Notgood_sum int64
	Notgood_pay float64
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
	Lock   sync.RWMutex
}

func (brand *Brand) Merge(vic *Brand, t int64) {
	if t == 0 {
		brand.Pay += vic.Pay
		brand.Sum += vic.Sum
		if brand.Number == 0 {
			brand.Number = vic.Number
		}
	} else {
		brand.Notgood_pay += vic.Notgood_pay
		brand.Notgood_sum += vic.Notgood_sum
	}
}

func (res *Result_t) Merge(result *BasicResult, t int64) {
	res.Lock.Lock()
	defer res.Lock.Unlock()
	for i, value := range result.Answer {
		res.Result.Answer[i].Merge(&value, t)
	}
	if t == 0 {
		res.Result.All_pay += result.All_pay
		res.Result.All_sum += result.All_sum
	} else {
		res.Result.Notgood_pay += result.Notgood_pay
		res.Result.Notgood_sum += result.Notgood_sum
	}
}

func (res *BasicResult) ToJson() ([]byte, error) {
	result, err := json.Marshal(res)
	return result, err
}

func (res *BasicResult) FromJson(data []byte) error {
	err := json.Unmarshal(data, res)
	if err != nil {
		panic(err)
	}
	return err
}

var Result Result_t
var TotalResult Result_t // only For master

func (result *Result_t) Print() {
	result.Lock.RLock()
	defer result.Lock.RUnlock()
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 5; j++ {
			fmt.Printf("品牌:#%d 一共有%d个零件\n", i*10+j, result.Result.Answer[i*10+j].Number)
			fmt.Printf("品牌:#%d 总销售额为%f\n", i*10+j, result.Result.Answer[i*10+j].Pay)
			fmt.Printf("品牌:#%d 总销售量为%d\n", i*10+j, result.Result.Answer[i*10+j].Sum)
			fmt.Printf("品牌:#%d 非畅销品的销售额为%d\n", i*10+j, result.Result.Answer[i*10+j].Notgood_sum)
			fmt.Printf("品牌:#%d 非畅销品的销售额为%f\n", i*10+j, result.Result.Answer[i*10+j].Notgood_pay)
			fmt.Printf("\n")
		}
	}
	fmt.Printf("总的销售额为%f\n总的销售量为%d\n总的平均销售额为%f\n", result.Result.All_pay, result.Result.All_sum, float64(result.Result.All_pay/float64(result.Result.All_sum)))
	fmt.Printf("非畅销总额为%f\n非畅销的总销售量为%d\n", result.Result.Notgood_pay, result.Result.Notgood_sum)
}
