package network

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"src/src/framework"
)

func response(w http.ResponseWriter, r *http.Request) {
	task, ok := framework.GetTask()
	if !ok {
		fmt.Fprintf(w, "Get Task Error")
	} else {
		fmt.Fprintf(w, "%v", task) //这个写入到w的是输出到客户端的
	}
}

func StartTaskServer() bool {
	_, err := os.Stat(".haveServer")
	if err == nil { // err == nil 代表有文件存在，不需要创建server
		return false
	}

	_, err = os.Create(".haveServer")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", response)
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("StartTaskServer: ", err)
	}
	return true
}

func DeleteTaskServer() {
	_, err := os.Stat(".haveServer")
	if err != nil {
		return
	}

	os.Remove(".haveServer")
}
