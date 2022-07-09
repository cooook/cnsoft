package utils

import "fmt"

const StartPort = 9000
const URL = "http://localhost"

type HandlerType int

const (
	get_task = iota
	create_task
	post_answer
	print_answer
)

func GetHandlerString(Type HandlerType) string {
	switch Type {
	case get_task:
		return "/get_task"
	case create_task:
		return "/create_task"
	case post_answer:
		return "/post_answer"
	case print_answer:
		return "/print_answer"
	}
	panic(fmt.Errorf("don't support type:%d", Type))
}

func GetServerURL(Type HandlerType) string {
	return fmt.Sprintf("%s:%d%s", URL, StartPort, GetHandlerString(Type))
}
