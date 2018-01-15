package fmt

import (
	"nqc.cn/log"
	"github.com/petermattis/goid"
	"sync"
	sfmt "fmt"
	"time"
	"nqc.cn/utils"
)

type logInfo struct {
	text string
	start time.Time
}

var list map[int64]interface{} = make(map[int64]interface{})
var mutex sync.Mutex

func Start() {

	mutex.Lock()
	defer mutex.Unlock()

	writeLog(goid.Get())

	list[goid.Get()] = &logInfo{text:"",start:time.Now()}

	//sfmt.Println("goid.Get():",goid.Get())

}

func writeLog(gid int64) {
	if list[gid] != nil {
		logs := list[gid].(*logInfo)
		if len(logs.text) > 0 {

			t := time.Now()

			log.WritePrintln(logs.text,"\n","该请求服务器处理时间为：",t.Sub(logs.start).Seconds() ,"s，开始时间：",
				utils.FormatTimeAll(logs.start.Unix()),"结束时间:",utils.FormatTimeAll(t.Unix()),"\n")
		}
	}
}

func Over() {


	mutex.Lock()
	defer mutex.Unlock()

	writeLog(goid.Get())

	delete(list,goid.Get())
}


func Println(a ...interface{}) {

	mutex.Lock()
	defer mutex.Unlock()

	if list[goid.Get()] != nil {
		logs := list[goid.Get()].(*logInfo)
		logs.text += "\n" + sfmt.Sprint(a)
	} else {
		log.WritePrintln(a,"\n")
	}

}
