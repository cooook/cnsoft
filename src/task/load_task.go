package task

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	"unsafe"
)

type Part struct {
	Flag   int64   //是否非畅销
	Number int64   //编号
	Brand  int64   //品牌
	Sum    int64   //销售量
	Pay    float64 //销售额
}

type Line struct {
	Number   int64   //编号
	Part_sum int64   //零件数量
	Pay      float64 //零件价值
}

var LoadPart [2e6 + 5]Part
var Maxn, Cnt int

type LoadTask struct {
	Result       Result_t
	partFile     string
	lineItemFile string
	start        int64 // 区间左闭右开
	end          int64
	taskID       int64
}

func (task *LoadTask) SetTaskID(ID int64) {
	task.taskID = ID
}

func (task *LoadTask) GetTaskID() int64 {
	return task.taskID
}

func (task *LoadTask) ToJson() ([]byte, error) {
	return json.Marshal(task)
}

func (task *LoadTask) FromJson(json_data []byte) error {
	fmt.Print("There!!!\n")
	return json.Unmarshal(json_data, task)
}

func NewLoadTask(partFile, lineItemFile string, start, end int64) *LoadTask {
	return &LoadTask{partFile: partFile, lineItemFile: lineItemFile, start: start, end: end}
}

func (task *LoadTask) Split(Num int) []Task {
	size := task.end - task.start
	subSize := size/int64(Num) + 1
	result := make([]Task, 0, Num)
	for start := task.start; start <= task.end; start += subSize {
		result = append(result, NewLoadTask(task.partFile, task.lineItemFile, start, start+subSize-1))
	}
	return result
}

func (task *LoadTask) Run() {
	//读取所有的part存到结构体
	start := time.Now()

	fpart, err := os.Open(task.partFile)
	if err != nil {
		panic(err)
	}
	defer fpart.Close()
	i := 1
	Maxn = 0
	Cnt = 0
	for {
		err = binary.Read(fpart, binary.LittleEndian, &LoadPart[i])
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		// fmt.Println(LoadPart[i])
		i += 1
		Maxn += 1
	}
	// 如果要按偏移量start和end
	offset := int64(unsafe.Sizeof(LoadPart[0]))
	fpart.Seek((offset)*(task.start-1), 0)
	for i := task.start; i <= task.end; i++ {
		_ = binary.Read(fpart, binary.LittleEndian, &LoadPart[i])
	}

	cost := time.Since(start)
	fmt.Printf("cost=[%s]\n\n", cost)

	fline, err := os.Open(task.lineItemFile)
	if err != nil {
		panic(err)
	}
	defer fline.Close()

	var line_item Line

	for {
		err = binary.Read(fline, binary.LittleEndian, &line_item)
		if err == io.EOF {
			break
		}

		i := line_item.Number
		LoadPart[i].Pay += line_item.Pay
		LoadPart[i].Sum += line_item.Part_sum
		Cnt += 1
	}

	for i := 1; i <= Maxn; i++ {
		task.Result.Result.Answer[LoadPart[i].Brand].sum += LoadPart[i].Sum
		task.Result.Result.Answer[LoadPart[i].Brand].pay += LoadPart[i].Pay
		task.Result.Result.Answer[LoadPart[i].Brand].number += 1
		task.Result.Result.All_pay += LoadPart[i].Pay
		task.Result.Result.All_sum += LoadPart[i].Sum
	}
}
