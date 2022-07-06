package framework

import (
	"container/list"
	"sync"
	"time"
)

var taskList = list.New()
var taskLock sync.Mutex

func GetTask() (Task, bool) {
	taskLock.Lock()
	defer taskLock.Unlock()
	if taskList.Len() == 0 {
		return nil, false
	}

	element := taskList.Front()

	result, ok := element.Value.(Task)

	if !ok {
		return nil, false
	}

	taskList.Remove(element)

	return result, true
}

func CreateNewTask(task Task) bool {
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
			task, ok := GetTask()
			if !ok {
				time.Sleep(time.Duration(5) * time.Second)
			} else {
				task.Run()
			}
		}
	}
}
