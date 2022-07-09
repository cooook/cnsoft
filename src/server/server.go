package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"src/src/framework"
	"src/src/task"
	"src/src/utils"
)

const (
	get_task = iota
	create_task
	post_answer
	print_answer
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, ok := framework.GetTask()
	if !ok {
		http.Error(w, "Can't get task", http.StatusBadRequest)
	} else {
		json_data, err := task.ToJson()
		if err != nil {
			http.Error(w, "Can't get task", http.StatusBadRequest)
		}
		w.Write(json_data)
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 创建新的任务，提交给master，由master 分割之后再由runner运行
	var task task.Task // [TODO]

	// 将请求体中的 JSON 数据解析到结构体中
	// 发生错误，返回400 错误码
	defer r.Body.Close()
	json_data, _ := ioutil.ReadAll(r.Body)
	err := task.FromJson(json_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task.Split(utils.GetServerNum())
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

	task.Result.Merge(&result)
}

func clinetPrintAnswerHandler(w http.ResponseWriter, r *http.Request) {
	var result task.BasicResult

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result.FromJson(data)

	task.Result.Result.Answer = result.Answer
	// [TODO] Print
}

func startMaster(port int64) {
	http.HandleFunc(utils.GetHandlerString(get_task), getTaskHandler)
	http.HandleFunc(utils.GetHandlerString(create_task), createTaskHandler)
	http.HandleFunc(utils.GetHandlerString(post_answer), getClinetAnswerHandler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("StartTaskServer: ", err)
	}
}

func startSlaver(port int64) {
	http.HandleFunc(utils.GetHandlerString(print_answer), clinetPrintAnswerHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Printf("Error: %s, Port: %d", err.Error(), port)
	}
}

func StartTaskServer() {
	task.RunnerID = utils.GetRunnerID()
	fmt.Printf("RunnerID = %d\n", task.RunnerID)
	if task.RunnerID == 0 { // is Master
		go startMaster(utils.StartPort)
	} else {
		go startSlaver(task.RunnerID + utils.StartPort)
	}

}

func DeleteTaskServer() {
	utils.UnsetRunnerID(task.RunnerID)
}
