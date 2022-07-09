package main

import (
	"fmt"
	"src/src/framework"
	"src/src/server"
	"strings"
)

func handlerUserInput() {
	fmt.Print("Plz input command!\n")
	var command string
	for {
		fmt.Scan(&command)
		if strings.EqualFold(command, "load") {
			fmt.Print("LOAD\n")
			// server.
		} else if strings.EqualFold(command, "select") {
			fmt.Print("Select\n")
			// framework.CreateNewTask()
		} else {
			fmt.Print("No such command")
		}
	}
}

func main() {

	server.StartTaskServer()
	defer server.DeleteTaskServer()
	stopCh := make(chan bool)
	go framework.Start(stopCh)
	defer func() { stopCh <- true }()
	handlerUserInput()
}
