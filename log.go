package log

import (
	"fmt"
	"os"
	"nqc.cn/utils"
	"time"
	"io/ioutil"
	"runtime/debug"
)

var year int//上次写日志文件的年
var month string//上次写日志文件的月
var day int//上次写日志文件的日

var saveDays int = 7//日志保留天数，默认不删除
//初始化，log日志保留几天，默认不删除
func Init() {
	createNewLogFile()


}

func SetSaveDays(days int) {
	saveDays = days
	clearOneWeekLog(saveDays)
	go startTimer(saveDays)
}

func startTimer(days int) {
	now := time.Now()
	// 计算下一个零点
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

	fmt.Println(utils.FormatTimeAll(next.Unix()))

	second := time.Duration(next.Sub(now).Seconds())
	fmt.Println("距离下一次清除log还有: " ,next.Sub(now).Seconds(), " 秒")
	timer(second,days)
}


func timer(seconds time.Duration,days int) {

	timer := time.NewTicker(seconds * time.Second)
	for {
		select {
		case <-timer.C:
			{
				go clearOneWeekLog(days)
				timer2 := time.NewTicker(3600 * time.Second)
				for {
					select {
					case <-timer2.C:
						{
							startTimer(days)

							return
						}

					}
				}
				return
			}

		}
	}
}

func clearOneWeekLog(days int) {//清除七天前那天的log
	path := utils.GetCurrPath() + "log"

	list := getNowFiles(path)

	now := time.Now()

	for _,val := range list {
		value := val.(map[string]interface{})

		times := value["modtime"].(time.Time)


		fmt.Println(value["name"],times,times.Add(time.Hour * 24 * time.Duration(days - 1)).Before(now))
		if times.Add(time.Hour * 24 * time.Duration(days - 1)).Before(now) {
			name := value["name"].(string)

			os.Remove(name)
		}
	}

}

//根据时间创建日志文件，每天创建不同的日志文件
func createNewLogFile() {

	t1 := time.Now()

	y,m,d := t1.Date()

	if year != y || month != m.String() || d != day {//判断当前的年月日是否和记录的年月日相同，不同的话就换了一天
		t3 := t1.Format("2006_01_02_15_04_05")

		path := utils.GetCurrPath() + "log/log_" + t3 + ".txt"
		fmt.Println(path)
		logfile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Printf("%s\r\n", err.Error())
			os.Exit(-1)
		}

		year = y
		month = m.String()
		day = d

		os.Stdout = logfile
		os.Stderr = logfile
	}

	//go startTimer(saveDays)
}
//监听当前线程的崩溃事件，拦截，写入日志文件
func InitListner(str string) {
	if err := recover(); err != nil {
		WriteString(fmt.Sprintln(fmt.Sprintln() + fmt.Sprintln(str) + fmt.Sprintf(`error: %v %v`,fmt.Sprintln(err),string(debug.Stack()))))
	}
}
//写入error
func Write(err error) {
	defer InitListner("")
	panic(err)
}
//写入正常字符串
func WriteString(info string) {
	WritePrintln(info)
}
//写入多个拼接内容
func WritePrintln(a ...interface{}) {

	createNewLogFile()

	str := time.Now().Format("2006-01-02 15:04:05")

	fmt.Fprintln(os.Stdout,"\n/***********************------------------------------------------------------------------------------------------------------------------------")
	fmt.Fprintf(os.Stdout,str)
	fmt.Fprintln(os.Stdout,a,"\n------------------------------------------------------------------------------------------------------------------------***********************/")

}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func getNowFiles(dirPth string) []interface{} {
	var nowFile []interface{}
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nowFile
	}
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}


		name := fi.Name()
		value := make(map[string]interface{})
		value["name"] = dirPth + "/" + name
		value["modtime"] = fi.ModTime()
		nowFile = append(nowFile, value)
	}
	//fmt.Println("getNowFiles end")
	return nowFile
}