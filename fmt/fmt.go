//
// fmt函数，为多线程输出所用
//

package fmt

import (
	"github.com/NiuStar/log"
	"github.com/NiuStar/utils"
	"github.com/petermattis/goid"
	"sync"
	"time"
)

type logInfo struct {
	text  string
	start time.Time
}

var list map[int64]interface{} = make(map[int64]interface{})
var mutex sync.Mutex

//Start ...函数线程开始时调用
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

			log.Println("[fmt]:",logs.text, "\n", "该请求服务器处理时间为：", t.Sub(logs.start).Seconds(), "s，开始时间：",
				utils.FormatTimeAll(logs.start.Unix()), "结束时间:", utils.FormatTimeAll(t.Unix()), "\n")
		}
	}
}

//Over ...函数线程结束时调用
func Over() {

	mutex.Lock()
	defer mutex.Unlock()

	writeLog(goid.Get())

	delete(list, goid.Get())
}

//Println ...函数输出函数
func Println(a ...interface{}) {

	log.Println("[fmt]:",a, "\n")

}
