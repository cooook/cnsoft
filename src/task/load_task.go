package task

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"src/src/utils"
	"sync"
	"time"
	"unsafe"
)

const BUFFER_SIZE = 100000

type Part struct {
	Flag   int64   //是否非畅销
	Number int64   //编号
	Brand  int64   //品牌
	Sum    int64   //销售量
	Pay    float64 //销售额
	// Sync   bool
}

var Sync bool

type Line struct {
	Number   int64   //编号
	Part_sum int64   //零件数量
	Pay      float64 //零件价值
}

var LoadPart [2e6 + 5]Part
var LoadPartLock sync.RWMutex
var Maxn, Cnt int

type LoadTask struct {
	PartFile     string
	LineItemFile string
	Start        int64
	End          int64
	TaskID       int64
}

func (task *LoadTask) AllDoneCallBack() { // only call on master
	var wg sync.WaitGroup
	for i := 1; i < utils.GetServerNum(); i += 1 {
		wg.Add(1)
		go func(instanceID int) {
			client := &http.Client{}
			req, _ := http.NewRequest("POST", utils.GetSlaveURL(utils.Slave_sync_loadPart, int(instanceID)), nil)
			req.Header.Set("Content-Type", "application/json")
			response, err := client.Do(req)
			if err != nil {
				log.Fatal(err.Error())
				return
			}

			if response.StatusCode == 400 {
				return
			}

			var Sum int64
			var Pay float64
			LoadPartLock.Lock()
			for i := 0; i < Maxn; i += 1 {
				binary.Read(response.Body, binary.LittleEndian, &Sum)
				binary.Read(response.Body, binary.LittleEndian, &Pay)
				LoadPart[i].Sum += Sum
				LoadPart[i].Pay += Pay
			}
			response.Body.Close()
			LoadPartLock.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()

	var buffer bytes.Buffer
	for i := 1; i < Maxn; i += 1 {
		binary.Write(&buffer, binary.LittleEndian, LoadPart[i].Sum)
		binary.Write(&buffer, binary.LittleEndian, LoadPart[i].Pay)
	}
	for i := 1; i < utils.GetServerNum(); i += 1 {
		var tmp_buffer bytes.Buffer
		tmp_buffer.Write(buffer.Bytes())
		wg.Add(1)
		go func(instanceID int) {
			client := &http.Client{}
			req, _ := http.NewRequest("POST", utils.GetSlaveURL(utils.Sync_loadPart, instanceID), bytes.NewReader(tmp_buffer.Bytes()))
			req.Header.Set("Content-Type", "application/json")

			response, err := client.Do(req)
			if err != nil {
				log.Fatal(err.Error())
				return
			}

			if response.StatusCode == 400 {
				return
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	CreateTaskToServer(NewSelectTask(1, 2000000), Select_task)
}

func (task *LoadTask) MergeAnwer() bool {
	Result.Lock.RLock()
	ok := PostAnswerToServer(&Result.Result, task.GetTaskID(), 0)
	Result.Lock.RUnlock()
	return ok
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

func (task *LoadTask) Initial() {
	Result.Lock.Lock()
	defer Result.Lock.Unlock()

	for i := range Result.Result.Answer {
		Result.Result.Answer[i].Number = 0
		Result.Result.Answer[i].Pay = 0
		Result.Result.Answer[i].Sum = 0
	}
	Result.Result.All_pay = 0
	Result.Result.All_sum = 0
}

func (task *LoadTask) Run() {
	//读取所有的part存到结构体
	start := time.Now()

	fpart, err := os.Open(task.PartFile)
	if err != nil {
		panic(err)
	}
	defer fpart.Close()

	part_offset := unsafe.Sizeof(LoadPart[0])
	buffer1 := make([]byte, BUFFER_SIZE*part_offset)
	fdreader := bufio.NewReader(fpart)
	var buffer1_reader *bytes.Reader

	i := 1
	Maxn = 0
	Cnt = 0
	cur := 0
	n := 0

	for {
		if cur == n {
			n, err = fdreader.Read(buffer1)
			cur = 0
			if err == io.EOF {
				break
			}
			buffer1_reader = bytes.NewReader(buffer1)
		}

		err = binary.Read(buffer1_reader, binary.LittleEndian, &LoadPart[i])
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		i += 1
		Maxn += 1
		cur += int(part_offset)
	}

	// 如果要按偏移量start和end
	var line_item Line
	offset := int64(unsafe.Sizeof(line_item))

	cost := time.Since(start)
	fmt.Printf("cost=[%s]\n\n", cost)

	fline, err := os.Open(task.LineItemFile)
	fline.Seek((offset)*(task.Start-1), 0)
	if err != nil {
		panic(err)
	}
	defer fline.Close()

	buffer := make([]byte, BUFFER_SIZE*offset)
	r := bufio.NewReader(fline)
	var buffer_reader *bytes.Reader

	cur = 0
	n = 0

	for i := task.Start; i <= task.End; i++ {
		if cur == n {
			n, err = r.Read(buffer)
			cur = 0
			if err == io.EOF {
				break
			}
			buffer_reader = bytes.NewReader(buffer)
		}

		if err != nil {
			panic(err)
		}

		err = binary.Read(buffer_reader, binary.LittleEndian, &line_item)
		if err == io.EOF {
			break
		}

		i := line_item.Number
		LoadPart[i].Pay += line_item.Pay
		LoadPart[i].Sum += line_item.Part_sum

		Cnt += 1
		cur += int(offset)
	}

	Result.Lock.Lock()
	for i := 1; i <= Maxn; i++ {
		Result.Result.Answer[LoadPart[i].Brand].Sum += LoadPart[i].Sum
		Result.Result.Answer[LoadPart[i].Brand].Pay += LoadPart[i].Pay
		Result.Result.Answer[LoadPart[i].Brand].Number += 1
		Result.Result.All_pay += LoadPart[i].Pay
		Result.Result.All_sum += LoadPart[i].Sum
	}
	Result.Lock.Unlock()

	cost = time.Since(start)
	fmt.Printf("cost=[%s]\n\n", cost)

}
