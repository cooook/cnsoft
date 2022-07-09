package client

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"src/src/task"
	"src/src/utils"
	"strconv"
)

const (
	get_task = iota
	create_task
	post_answer
	print_answer
)

func CreateTaskToServer(task task.Task) bool {
	client := &http.Client{}
	json_data, err := task.ToJson()
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	req, _ := http.NewRequest("POST", utils.GetServerURL(create_task), bytes.NewBuffer(json_data))
	req.Header.Set("Content-Type", "application/json")
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

func GetTaskFromServer() (task.Task, bool) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", utils.GetServerURL(get_task), nil)
	req.Header.Set("RunnerId", strconv.FormatInt(task.RunnerID, 10))
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
		return nil, false
	}
	if response.StatusCode == 400 {
		return nil, false
	}

	var task task.Task
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

func PostAnswerToServer(result task.BasicResult) bool {
	client := &http.Client{}
	json_data, err := result.ToJson()
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	req, _ := http.NewRequest("POST", utils.GetServerURL(post_answer), bytes.NewBuffer(json_data))
	req.Header.Set("Content-Type", "application/json")
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

func PrintCommandToClient(result task.BasicResult) bool {
	client := &http.Client{}
	json_data, err := result.ToJson()
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	req, _ := http.NewRequest("POST", utils.GetServerURL(print_answer), bytes.NewBuffer(json_data))
	req.Header.Set("Content-Type", "application/json")
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
