package main

import (
	"fmt"
	"src/src/framework"
	"src/src/network"
	"strings"
)

func handlerUserInput() {
	fmt.Print("Plz input command!")
	var command string
	for {
		fmt.Scan(&command)
		if strings.EqualFold(command, "load") {
			framework.CreateNewTask()
		} else if strings.EqualFold(command, "select") {
			framework.CreateNewTask()
		} else {
			fmt.Print("No such command")
		}
	}
}

func main() {
	start := network.StartTaskServer()
	if start {
		defer network.DeleteTaskServer()
	}
	stopCh := make(chan bool)
	go framework.Start(stopCh)
	defer func() { stopCh <- true }()
	handlerUserInput()
}
