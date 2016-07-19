package log

import (
	"fmt"
	"log"
	"os"
	"nqc.cn/utils"
	"time"
	"sync"
	"runtime/debug"
)
var mutex sync.Mutex
var errList []error
var strList []string
var ch chan bool
func writeAll(logger *log.Logger) {
		for {
			mutex.Lock()
			logger.SetPrefix("[Info]")
			for _,value := range strList {
				logger.Println(value)
			}
			errList=append(errList[:0],errList[len(errList):]...)
			strList=append(strList[:0],strList[len(strList):]...)
			mutex.Unlock()
			time.Sleep(3 * time.Second)
		}
}

func Init() {

	ch = make(chan bool)
	//var ch chan int、
	t3 := time.Now().Format("2006_01_02_15_04_05")
	path := utils.GetCurrPath() + "log/log_" + t3 + ".txt"
	fmt.Println(path)
	logfile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	//go writeAll(logger)
	//defer logfile.Close()
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	logger.SetFlags(logger.Flags() | log.LstdFlags)
	logger.Println(int(time.Now().Unix()))

	logger.Println("log初始化完成")

	go writeAll(logger)
	//flag.Parse()
}

func InitListner() {
	if err := recover(); err != nil {
		WriteString(fmt.Sprintln(fmt.Sprintln() + fmt.Sprintf(`error: %v %v`,fmt.Sprintln(err),string(debug.Stack()))))
	}
}

func Write(err error) {
	defer InitListner()
	panic(err)
}

func WriteString(info string) {
	mutex.Lock()
	strList = append(strList,info)
	mutex.Unlock()
}
