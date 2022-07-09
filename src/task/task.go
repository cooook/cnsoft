package task

import (
	"encoding/json"
	"sync"
)

var RunnerID int64

type Task interface {
	Run()
	Split(Num int)                   // split to Num subTask and submit to taskServer
	ToJson() ([]byte, error)         // need to add RunnerID to json
	FromJson(json_data []byte) error // need to add RunnerID to json
}

type BasicResult struct {
	Answer [25][25]int
}

type Result_t struct {
	Result BasicResult
	lock   sync.Mutex
}

func (res *Result_t) Merge(result *BasicResult) {
	res.lock.Lock()
	defer res.lock.Unlock()
	for i, value := range result.Answer {
		for j, value2 := range value {
			res.Result.Answer[i][j] += value2
		}
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
