package log

import (
	"fmt"
	"log"
	"os"
	"nqc.cn/utils"
	"time"
	"syscall"

)

func GetLogFile() *os.File {
	t3 := time.Now().Format("2006_01_02_15_04_05")
	fmt.Println(t3)
	path := utils.GetCurrPath() + "log/log_" + t3 + ".txt"
	logfile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}

	//syscall.Dup2(int(logfile.Fd()), 1)
	syscall.Dup2(int(logfile.Fd()), 2)
	defer logfile.Close()
	return logfile
}

var errList []error
var strList []string
var ch chan bool

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}
func writeAll(logger *log.Logger) {

		for {
			select {
			case <-ch:
				{
					for _,value := range errList {
						panic(value)
					}
					logger.SetPrefix("[Info]")
					for _,value := range strList {

						logger.Println(value)

					}
					errList=append(errList[:0],errList[len(errList):]...)
					strList=append(strList[:0],strList[len(strList):]...)
				}
				break
			}
		}
	//<-ch//阻塞协程
	//ch <- 1//释放协程
}

func Init() {
	ch = make(chan bool)
	t3 := time.Now().Format("2006_01_02_15_04_05")
	fmt.Println(t3)
	path := utils.GetCurrPath() + "log/log_" + t3 + ".txt"
	logfile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	//syscall.Dup2(int(logfile.Fd()), 1)
	syscall.Dup2(int(logfile.Fd()), 2)
	defer logfile.Close()
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	logger.SetFlags(logger.Flags() | log.LstdFlags)
	logger.Println(int(time.Now().Unix()))

	logger.Println("log初始化完成")
	//var ch chan int
	go writeAll(logger)

}

func InitListner() {
}

func Write(err error) {

	fmt.Println("error: ",err)
	ch1 := make(chan bool)
	go func() {
		ch1 <- true
		panic(err)

	} ()
	select {
	case <- ch1:
		{
			return
		}

	}
}
func WriteString(info string) {
	strList = append(strList,info)
	ch <- true
}
