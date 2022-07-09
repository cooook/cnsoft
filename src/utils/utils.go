package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

const ENV_NAME = "CNSOFT_SERVER_NUM"
const MAX_SERVER_NUM = 64

var lock sync.RWMutex

const FILE_NAME = "./ServerState"

func getServerState() int64 {
	lock.RLock()
	defer lock.RUnlock()

	if _, err := os.Stat(FILE_NAME); err != nil {
		os.Create(FILE_NAME)
	}

	file, err := os.Open(FILE_NAME)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	number, _ := strconv.ParseInt(string(content), 10, 32)
	return number
}

func setServerState(target int64) {
	lock.Lock()
	defer lock.Unlock()

	if err := ioutil.WriteFile(FILE_NAME, []byte(strconv.FormatInt(target, 10)), 0644); err != nil {
		panic(err)
	}
}

func GetRunnerID() int64 {
	number := getServerState()

	for i := 0; i < MAX_SERVER_NUM; i += 1 {
		if (int(number) >> i & 1) == 0 {
			setServerState(number | (1 << i))
			return int64(i)
		}
	}
	panic("Can't support so much server")
}

func UnsetRunnerID(id int64) {
	serverState := getServerState()
	if serverState&(1<<id) == 0 {
		panic(fmt.Errorf("can't stop server, id = %d", id))
	}
	setServerState(serverState ^ (1 << id))
}

func GetServerNum() int {
	servetState := getServerState()
	Answer := 0
	for servetState != 0 {
		Answer += 1
		servetState >>= 1
	}
	return Answer
}
