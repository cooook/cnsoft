package task

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"src/src/utils"
	"strconv"
)

func CreateTaskToServer(task Task, t TaskType) bool {
	client := &http.Client{}
	json_data, err := task.ToJson()
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	req, _ := http.NewRequest("POST", utils.GetServerURL(utils.Create_task), bytes.NewBuffer(json_data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Task-Type", strconv.FormatInt(int64(t), 10))
	req.Header.Set("RunnerID", strconv.FormatInt(RunnerID, 10))
	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	if response.StatusCode == 400 {
		return false
	}
	return true
}

func GetTaskFromServer() (Task, bool) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", utils.GetServerURL(utils.Get_task), nil)
	req.Header.Set("RunnerId", strconv.FormatInt(RunnerID, 10))
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
		return nil, false
	}
	if response.StatusCode == 400 {
		return nil, false
	}

	t, err := strconv.ParseInt(response.Header.Get("Task-Type"), 10, 64)
	if err != nil {
		panic(err)
	}
	task := GetTypeTask(TaskType(t))

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
		return nil, false
	}

	err = task.FromJson(data)
	if err != nil {
		log.Fatal(err.Error())
		return nil, false
	}
	return task, true

}

func PostAnswerToServer(result *BasicResult, taskID int64, mergeType int64) bool {
	client := &http.Client{}
	json_data, err := result.ToJson()
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	req, _ := http.NewRequest("POST", utils.GetServerURL(utils.Post_answer), bytes.NewReader(json_data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("RunnerID", strconv.FormatInt(RunnerID, 10))
	req.Header.Set("TaskID", strconv.FormatInt(taskID, 10))
	req.Header.Set("Merge-Type", strconv.FormatInt(mergeType, 10))
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	if response.StatusCode == 400 {
		return false
	}
	return true
}

func SyncAnswer(result *BasicResult, instanceID int) bool {
	client := &http.Client{}
	json_data, err := result.ToJson()
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	req, _ := http.NewRequest("POST", utils.GetSlaveURL(utils.Sync_answer, instanceID), bytes.NewBuffer(json_data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("RunnerID", strconv.FormatInt(RunnerID, 10))
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	if response.StatusCode == 400 {
		return false
	}
	return true
}

func SyncDone(instanceID int) bool {
	client := &http.Client{}

	req, _ := http.NewRequest("POST", utils.GetSlaveURL(utils.Sync_done, instanceID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("RunnerID", strconv.FormatInt(RunnerID, 10))
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	if response.StatusCode == 400 {
		return false
	}
	return true
}

func Broadcast(s string) {
	client := &http.Client{}

	for i := 0; i < utils.GetServerNum(); i += 1 {
		req, _ := http.NewRequest("POST", utils.GetSlaveURL(utils.Broadcast, i), bytes.NewBuffer([]byte(s)))
		req.Header.Set("RunnerID", strconv.FormatInt(RunnerID, 10))
		response, err := client.Do(req)
		if err != nil {
			log.Fatal(err.Error())
			break
		}

		if response.StatusCode == 400 {
			break
		}
	}
}

func SyncRequest() bool {
	client := &http.Client{}

	req, _ := http.NewRequest("POST", utils.GetServerURL(utils.Sync_request), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("RunnerID", strconv.FormatInt(RunnerID, 10))
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	if response.StatusCode == 400 {
		return false
	}
	return true
}

func SyncLoadPartRequest(idx int64) bool {
	client := &http.Client{}

	req, _ := http.NewRequest("POST", utils.GetServerURL(utils.Sync_loadPart), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("RunnerID", strconv.FormatInt(RunnerID, 10))
	req.Header.Set("idx", strconv.FormatInt(idx, 10))
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	if response.StatusCode == 400 {
		return false
	}

	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		panic(err)
	}

	var load Part
	if err := json.Unmarshal(data, &load); err != nil {
		panic(err)
	}

	LoadPartLock.Lock()
	LoadPart[idx].Pay += load.Pay
	LoadPart[idx].Sum += load.Sum
	LoadPartLock.Unlock()

	return true
}
