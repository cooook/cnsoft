package main

import (
	"fmt"
	"log"
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
			ok := task.CreateTaskToServer(task.NewLoadTask("../test/part.bin", "../test/lineitem.bin", 1, 59986052), task.Load_task)
			if !ok {
				log.Fatal("Create Load Task Error!")
			}
		} else if strings.EqualFold(command, "select") {
			if !task.LoadDone {
				fmt.Println("Not load over yet, plz wait...")
			} else {
				task.SyncRequest()
				task.Result.Print()
			}
		} else {
			fmt.Println("No such command")
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
