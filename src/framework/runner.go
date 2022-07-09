package framework

import (
	"container/list"
	"src/src/client"
	"src/src/task"
	"sync"
	"time"
)

var taskList = list.New()
var taskLock sync.Mutex

func GetTask() (task.Task, bool) {
	taskLock.Lock()
	defer taskLock.Unlock()
	if taskList.Len() == 0 {
		return nil, false
	}

	element := taskList.Front()

	result, ok := element.Value.(task.Task)

	if !ok {
		return nil, false
	}

	taskList.Remove(element)

	return result, true
}

func CreateNewTask(task task.Task) bool {
	taskLock.Lock()
	defer taskLock.Unlock()
	taskList.PushBack(task)
	return true
}

func Start(stopCh <-chan bool) {
	for {
		select {
		case <-stopCh:
			return
		default:
			task, ok := client.GetTaskFromServer()
			if !ok {
				time.Sleep(time.Duration(1) * time.Second)
			} else {
				task.Run()
			}
		}
	}
}
