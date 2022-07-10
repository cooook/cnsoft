package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

var RunnerID int64

type Task interface {
	Run()
	Split(Num int)                   // split to Num subTask and submit to taskServer
	ToJson() ([]byte, error)         // need to add RunnerID to json
	FromJson(json_data []byte) error // need to add RunnerID to json
}

type Brand struct {
	pay         float64 //总销售额
	number      int     //零件个数
	sum         int     //销售量
	notgood_sum int     //非畅销销售量
	notgood_pay float64 //非畅销销售额
}

type BasicResult struct {
	Answer [60]Brand
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

func main() {
	var result BasicResult
	result.Answer[0].notgood_pay = 1
	data, _ := result.ToJson()
	result.FromJson(data)
	fmt.Print(result)
}
