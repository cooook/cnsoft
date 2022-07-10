package task

import (
	"encoding/json"
	"fmt"
	"src/src/utils"
)

type SelectTask struct {
	Start  int64
	End    int64
	TaskID int64
}

var LoadDone = false

func NewSelectTask(start, end int64) *SelectTask {
	return &SelectTask{Start: start, End: end}
}

func (task *SelectTask) AllDoneCallBack() { // only call on master
	// CreateTaskToServer()
	for i := 1; i < utils.GetServerNum(); i += 1 {
		SyncDone(i)
	}
	go Broadcast("Load done.")
	LoadDone = true
}

func (task *SelectTask) MergeAnwer() bool {
	Result.Lock.RLock()
	ok := PostAnswerToServer(&Result.Result, task.GetTaskID(), 1)
	Result.Lock.RUnlock()
	return ok
}

func (task *SelectTask) SetTaskID(ID int64) {
	task.TaskID = ID
}

func (task *SelectTask) GetTaskID() int64 {
	return task.TaskID
}

func (task *SelectTask) ToJson() ([]byte, error) {
	return json.Marshal(task)
}

func (task *SelectTask) FromJson(json_data []byte) error {
	json.Unmarshal(json_data, task)
	return json.Unmarshal(json_data, task)
}

func (task *SelectTask) GetTaskType() TaskType {
	return Select_task
}

func (task *SelectTask) Split(Num int) []Task {
	size := task.End - task.Start
	subSize := size/int64(Num) + 1
	result := make([]Task, 0, Num)
	for start := task.Start; start <= task.End; start += subSize {
		end := start + subSize - 1
		if end > task.End {
			end = task.End
		}
		result = append(result, NewSelectTask(start, end))
	}
	return result
}

func (task *SelectTask) Initial() {
	SyncRequest()
	Result.Lock.Lock()
	defer Result.Lock.Unlock()
	for i := range Result.Result.Answer {
		Result.Result.Answer[i].Notgood_pay = 0
		Result.Result.Answer[i].Notgood_sum = 0
	}
	Result.Result.Notgood_pay = 0
	Result.Result.Notgood_sum = 0
}

func (task *SelectTask) Run() {
	var average_sum [60]float64

	Result.Lock.RLock()
	for i := 1; i <= 5; i++ {
		for j := 1; j <= 5; j++ {
			id := i*10 + j
			average_sum[id] = float64(1.0 * Result.Result.Answer[id].Sum / Result.Result.Answer[id].Number)
			fmt.Printf("id = %d, avr = %f\n", id, average_sum[id])
		}
	}
	Result.Lock.RUnlock()

	Result.Lock.Lock()
	for i := task.Start; i <= task.End; i++ {
		if float64(LoadPart[i].Sum) < float64(average_sum[LoadPart[i].Brand]*0.3) {
			LoadPart[i].Flag = int64(0)
			Result.Result.Answer[LoadPart[i].Brand].Notgood_sum += LoadPart[i].Sum
			Result.Result.Answer[LoadPart[i].Brand].Notgood_pay += LoadPart[i].Pay
			Result.Result.Notgood_pay += LoadPart[i].Pay
			Result.Result.Notgood_sum += LoadPart[i].Sum
		}
	}
	Result.Lock.Unlock()
}
