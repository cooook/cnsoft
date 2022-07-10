package utils

import "fmt"

const StartPort = 9000
const URL = "http://localhost"

type HandlerType int

const (
	Get_task = iota
	Create_task
	Post_answer
	Sync_answer
	Broadcast
	Sync_done
	Sync_request
	Slave_sync_loadPart
	Sync_loadPart
)

func GetHandlerString(Type HandlerType) string {
	switch Type {
	case Get_task:
		return "/get_task"
	case Create_task:
		return "/create_task"
	case Post_answer:
		return "/post_answer"
	case Sync_answer:
		return "/sync_answer"
	case Broadcast:
		return "/broadcast"
	case Sync_done:
		return "/sync_done"
	case Sync_request:
		return "/sync_request"
	case Slave_sync_loadPart:
		return "/slave_sync_loadpart"
	case Sync_loadPart:
		return "/sync_loadpart"
	}
	panic(fmt.Errorf("don't support type:%d", Type))
}

func GetServerURL(Type HandlerType) string {
	return fmt.Sprintf("%s:%d%s", URL, StartPort, GetHandlerString(Type))
}

func GetSlaveURL(Type HandlerType, ID int) string {
	return fmt.Sprintf("%s:%d%s", URL, StartPort+ID, GetHandlerString(Type))
}
