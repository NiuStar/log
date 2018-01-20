//
//a example
//

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NiuStar/log"
	"net/http"
)

func main() {
	log.Init()
	fmt.Println("fmt")

	j1 := make(map[string]interface{})
	log.WriteString("1.Println log with log.LstdFlags ...")

	go func() {
		log.Write(errors.New("this is a BUG"))
	}()

	_, err := http.Get("cninct.com")
	if err != nil {
		fmt.Println(err)
		log.Write(err)
	}
	fmt.Println("123456")
	body, err := json.Marshal(j1)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

}
