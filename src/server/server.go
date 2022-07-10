package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"src/src/framework"
	"src/src/task"
	"src/src/utils"
	"strconv"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, ok := framework.GetTask()
	if !ok {
		http.Error(w, "Can't get task", http.StatusBadRequest)
	} else {
		w.Header().Set("Task-Type", strconv.FormatInt(int64(task.GetTaskType()), 10))
		json_data, err := task.ToJson()
		if err != nil {
			http.Error(w, "Can't get task", http.StatusBadRequest)
		}
		w.Write(json_data)
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 创建新的任务，提交给master，由master 分割之后再由runner运行
	// var task task.LoadTask // [TODO]
	// 将请求体中的 JSON 数据解析到结构体中
	// 发生错误，返回400 错误码
	json_data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		panic(err)
	}
	t, err := strconv.ParseInt(r.Header.Get("Task-Type"), 10, 64)
	if err != nil {
		panic(err)
	}

	task := task.GetTypeTask(task.TaskType(t))

	err = task.FromJson(json_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	framework.CreateOriginTask(task)
}

func getClinetAnswerHandler(w http.ResponseWriter, r *http.Request) {
	var result task.BasicResult

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result.FromJson(data)

	taskIDtext := r.Header.Get("TaskID")
	taskID, err := strconv.ParseInt(taskIDtext, 10, 64)
	if err != nil {
		panic(err)
	}

	// RunnerIDtext := r.Header.Get("RunnerID")
	// runnerID, err := strconv.ParseInt(RunnerIDtext, 10, 64)
	// if err != nil {
	// 	panic(err)
	// }

	// if taskID == 0 {
	// 	fmt.Printf("\n\n\n\nReceive from %d:\n", runnerID)
	// 	fmt.Println(result)
	// }
	// fmt.Println("\n\n\nSelf:")
	// fmt.Println(task.TotalResult.Result)

	MergeTypeText := r.Header.Get("Merge-Type")
	mergeType, err := strconv.ParseInt(MergeTypeText, 10, 64)
	if err != nil {
		panic(err)
	}

	framework.End(taskID)
	task.TotalResult.Merge(&result, mergeType)
	// fmt.Print("\n\n\nAfter merge:\n")
	// fmt.Println(task.TotalResult.Result)
}

func syncAnswerHandler(w http.ResponseWriter, r *http.Request) {
	var result task.BasicResult

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// RunnerIDtext := r.Header.Get("RunnerID")
	// runnerID, err := strconv.ParseInt(RunnerIDtext, 10, 64)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("sync runnerID:", runnerID)

	err = result.FromJson(data)
	if err != nil {
		panic(err)
	}

	task.Result.Lock.Lock()
	task.Result.Result = result
	task.Result.Lock.Unlock()
}

func broadcastHandler(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	// print Msg from server
	fmt.Println(string(data))
}

func syncDoneHandler(w http.ResponseWriter, r *http.Request) {
	task.LoadDone = true
}

func syncRequestHandler(w http.ResponseWriter, r *http.Request) {
	RunnerIDtext := r.Header.Get("RunnerID")
	runnerID, err := strconv.ParseInt(RunnerIDtext, 10, 64)
	if err != nil {
		panic(err)
	}

	task.TotalResult.Lock.RLock()
	task.SyncAnswer(&task.TotalResult.Result, int(runnerID))
	task.TotalResult.Lock.RUnlock()
}

func getHeader(r *http.Request, s string) int64 {
	Text := r.Header.Get(s)
	num, err := strconv.ParseInt(Text, 10, 64)
	if err != nil {
		panic(err)
	}
	return num
}

func syncLoadPartHandler(w http.ResponseWriter, r *http.Request) {
	var Sum int64
	var Pay float64
	task.LoadPartLock.Lock()
	for i := 0; i < task.Maxn; i += 1 {
		binary.Read(r.Body, binary.LittleEndian, &Sum)
		binary.Read(r.Body, binary.LittleEndian, &Pay)
		task.LoadPart[i].Sum = Sum
		task.LoadPart[i].Pay = Pay
	}
	r.Body.Close()
	task.LoadPartLock.Unlock()
}

func readLoadPartFromSlave(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer
	for i := 0; i < task.Maxn; i += 1 {
		binary.Write(&buffer, binary.LittleEndian, task.LoadPart[i].Sum)
		binary.Write(&buffer, binary.LittleEndian, task.LoadPart[i].Pay)
	}
	w.Write(buffer.Bytes())
}

func startMaster(port int64) {
	http.HandleFunc(utils.GetHandlerString(utils.Get_task), getTaskHandler)
	http.HandleFunc(utils.GetHandlerString(utils.Create_task), createTaskHandler)
	http.HandleFunc(utils.GetHandlerString(utils.Post_answer), getClinetAnswerHandler)
	http.HandleFunc(utils.GetHandlerString(utils.Sync_request), syncRequestHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("StartTaskServer: ", err)
	}
}

func startSlaver(port int64) {
	http.HandleFunc(utils.GetHandlerString(utils.Sync_done), syncDoneHandler)
	http.HandleFunc(utils.GetHandlerString(utils.Sync_loadPart), syncLoadPartHandler)
	http.HandleFunc(utils.GetHandlerString(utils.Slave_sync_loadPart), readLoadPartFromSlave)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Printf("Error: %s, Port: %d", err.Error(), port)
	}
}

func StartTaskServer() {
	task.RunnerID = utils.GetRunnerID()
	fmt.Printf("RunnerID = %d\n", task.RunnerID)
	http.HandleFunc(utils.GetHandlerString(utils.Sync_answer), syncAnswerHandler)
	http.HandleFunc(utils.GetHandlerString(utils.Broadcast), broadcastHandler)
	if task.RunnerID == 0 { // is Master
		go startMaster(utils.StartPort)
	} else {
		go startSlaver(task.RunnerID + utils.StartPort)
	}

}

func DeleteTaskServer() {
	utils.UnsetRunnerID(task.RunnerID)
}
