package framework

import (
	"container/list"
	"src/src/client"
	"src/src/task"
	"src/src/utils"
	"sync"
	"time"
)

var taskList = list.New()
var taskLock sync.Mutex

const MAX_ORIGIN_TASK = 10

var originTask = make([]*task.OriginTask, 0, MAX_ORIGIN_TASK)
var originLock sync.Mutex
var originSize int64 = 0

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

func createNewTask(task task.Task) bool {
	taskLock.Lock()
	defer taskLock.Unlock()
	taskList.PushBack(task)
	return true
}

func CreateOriginTask(tsk task.Task) int64 {
	originLock.Lock()
	defer originLock.Unlock()
	if originSize == MAX_ORIGIN_TASK {
		return -1
	}
	oTask := task.NewOriginTask(tsk, originSize)
	originTask = append(originTask, oTask)
	
	taskPoll := tsk.Split(utils.GetServerNum())
	for _, t := range taskPoll {
		t.SetTaskID(originSize)
		createNewTask(t)
		oTask.WG.Add(1)
	}
	
	originSize += 1
	return originSize - 1
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
