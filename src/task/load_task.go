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
	PartFile     string
	LineItemFile string
	Start        int64
	End          int64
	TaskID       int64
}

func (task *LoadTask) SetTaskID(ID int64) {
	task.TaskID = ID
}

func (task *LoadTask) GetTaskID() int64 {
	return task.TaskID
}

func (task *LoadTask) ToJson() ([]byte, error) {
	return json.Marshal(task)
}

func (task *LoadTask) FromJson(json_data []byte) error {
	json.Unmarshal(json_data, task)
	fmt.Println(task)
	return json.Unmarshal(json_data, task)
}

func (task *LoadTask) GetTaskType() TaskType {
	return Load_task
}

func NewLoadTask(partFile, lineItemFile string, start, end int64) *LoadTask {
	return &LoadTask{PartFile: partFile, LineItemFile: lineItemFile, Start: start, End: end}
}

func (task *LoadTask) Split(Num int) []Task {
	size := task.End - task.Start
	subSize := size/int64(Num) + 1
	result := make([]Task, 0, Num)
	for start := task.Start; start <= task.End; start += subSize {
		end := start + subSize - 1
		if end > task.End {
			end = task.End
		}
		result = append(result, NewLoadTask(task.PartFile, task.LineItemFile, start, end))
	}
	return result
}

func (task *LoadTask) Run() {
	//读取所有的part存到结构体
	start := time.Now()

	fpart, err := os.Open(task.PartFile)
	fmt.Println(task.PartFile)
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
	var line_item Line
	offset := int64(unsafe.Sizeof(line_item))

	fpart.Seek((offset)*(task.Start-1), 0)

	cost := time.Since(start)
	fmt.Printf("cost=[%s]\n\n", cost)

	fline, err := os.Open(task.LineItemFile)
	if err != nil {
		panic(err)
	}
	defer fline.Close()

	for i := task.Start; i < task.End; i++ {
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
		Result.Result.Answer[LoadPart[i].Brand].sum += LoadPart[i].Sum
		Result.Result.Answer[LoadPart[i].Brand].pay += LoadPart[i].Pay
		Result.Result.Answer[LoadPart[i].Brand].number += 1
		Result.Result.All_pay += LoadPart[i].Pay
		Result.Result.All_sum += LoadPart[i].Sum
	}

	cost = time.Since(start)
	fmt.Printf("cost=[%s]\n\n", cost)
}
