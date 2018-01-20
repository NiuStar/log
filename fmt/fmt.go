//
// fmt函数，为多线程输出所用
//

package fmt

import (
	sfmt "fmt"
	"github.com/petermattis/goid"
	"nqc.cn/log"
	"nqc.cn/utils"
	"sync"
	"time"
)

type logInfo struct {
	text  string
	start time.Time
}

var list map[int64]interface{} = make(map[int64]interface{})
var mutex sync.Mutex

//线程开始时调用
func Start() {

	mutex.Lock()
	defer mutex.Unlock()

	writeLog(goid.Get())

	list[goid.Get()] = &logInfo{text: "", start: time.Now()}

	//sfmt.Println("goid.Get():",goid.Get())

}

func writeLog(gid int64) {
	if list[gid] != nil {
		logs := list[gid].(*logInfo)
		if len(logs.text) > 0 {

			t := time.Now()

			log.WritePrintln(logs.text, "\n", "该请求服务器处理时间为：", t.Sub(logs.start).Seconds(), "s，开始时间：",
				utils.FormatTimeAll(logs.start.Unix()), "结束时间:", utils.FormatTimeAll(t.Unix()), "\n")
		}
	}
}

//线程结束时调用
func Over() {

	mutex.Lock()
	defer mutex.Unlock()

	writeLog(goid.Get())

	delete(list, goid.Get())
}

//输出函数
func Println(a ...interface{}) {

	mutex.Lock()
	defer mutex.Unlock()

	if list[goid.Get()] != nil {
		logs := list[goid.Get()].(*logInfo)
		logs.text += "\n" + sfmt.Sprint(a)
	} else {
		log.WritePrintln(a, "\n")
	}

}
