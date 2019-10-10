//
// Package log implements Go to https://github.com/NiuStar/log.
//

package log

import (
	"fmt"
	"github.com/NiuStar/utils"
	"io/ioutil"
	"os"
	//"runtime/debug"
	"github.com/sirupsen/logrus"
	"time"
)

var year int     //上次写日志文件的年
var month string //上次写日志文件的月
var day int      //上次写日志文件的日

var saveDays int = 30 //日志保留天数，默认不删除
//初始化，log日志保留几天，默认不删除
func Init(debug bool) {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
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
	fmt.Println("距离下一次清除log还有: ", next.Sub(now).Seconds(), " 秒")
	timer(second, days)
}

func timer(seconds time.Duration, days int) {

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
			}

		}
	}
}

func clearOneWeekLog(days int) { //清除七天前那天的log
	path := utils.GetCurrPath() + "log"

	list := getNowFiles(path)

	now := time.Now()

	for _, val := range list {
		value := val.(map[string]interface{})

		times := value["modtime"].(time.Time)

		fmt.Println(value["name"], times, times.Add(time.Hour*24*time.Duration(days-1)).Before(now))
		if times.Add(time.Hour * 24 * time.Duration(days-1)).Before(now) {
			name := value["name"].(string)

			os.Remove(name)
		}
	}

}

//根据时间创建日志文件，每天创建不同的日志文件
func createNewLogFile() {

	t1 := time.Now()

	y, m, d := t1.Date()

	if year != y || month != m.String() || d != day { //判断当前的年月日是否和记录的年月日相同，不同的话就换了一天
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

		logrus.SetOutput(logfile)
		//os.Stdout = logfile
		//os.Stderr = logfile
	}

	go startTimer(saveDays)
}

//监听当前线程的崩溃事件，拦截，写入日志文件
func InitListner(str string) {
	if err := recover(); err != nil {
		logrus.WithError(err.(error)).Error(str)
	}
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


// Debug logs a message at level Debug on the standard logger.
func Debug(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Debug(list...)
}

// Print logs a message at level Info on the standard logger.
func Print(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Print(list...)
}

// Info logs a message at level Info on the standard logger.
func Info(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(Tag string,format string, args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(Tag string,format string, args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(Tag string,format string, args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(Tag string,format string, args ...interface{}) {
	logrus.Warnf(Tag + format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(Tag string,format string, args ...interface{}) {
	logrus.Warningf(Tag + format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(Tag string,format string, args ...interface{}) {
	logrus.Errorf(Tag + format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(Tag string,format string, args ...interface{}) {
	logrus.Panicf(Tag + format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(Tag string,format string, args ...interface{}) {
	logrus.Fatalf(Tag + format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Debugln(list...)
}

// Println logs a message at level Info on the standard logger.
func Println(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Println(list...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Infoln(list...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Warnln(list...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Warningln(list...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Errorln(list...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Panicln(list...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(Tag string,args ...interface{}) {
	var list []interface{}
	list = append(list, Tag)
	list = append(list,args...)
	logrus.Fatalln(list...)
}
