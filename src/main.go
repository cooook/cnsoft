package main

import (
	"fmt"
	"log"
	"src/src/client"
	"src/src/framework"
	"src/src/server"
	"src/src/task"
	"strings"
)

func handlerUserInput() {
	fmt.Print("Plz input command!\n")
	var command string
	for {
		fmt.Scan(&command)
		if strings.EqualFold(command, "load") {
			ok := client.CreateTaskToServer(task.NewLoadTask("../test/part.bin", "../test/lineitem.bin", 1, 5000001), task.Load_task)
			if !ok {
				log.Fatal("Create Load Task Error!")
			}
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
