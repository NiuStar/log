package main

import (
	"fmt"
	"net/http"
	"nqc.cn/log"
	"encoding/json"
)

func main() {
	log.Init()
	fmt.Println("fmt")
	j1 := make(map[string]interface{})

	log.WriteString("1.Println log with log.LstdFlags ...")
	_, err := http.Get("cninct.com")
	if err != nil {
		fmt.Println(err)
		log.Write(err)
	}
	fmt.Println("123456")
	body ,err := json.Marshal(j1)
	fmt.Println(string(body))
}
