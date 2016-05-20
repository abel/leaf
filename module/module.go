package module

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/abel/leaf/conf"
	"github.com/abel/leaf/log"
)

type Module interface {
	OnInit()
	OnDestroy()
	Run(closeSig chan bool)
}

type module struct {
	mi       Module
	closeSig chan bool
	wg       sync.WaitGroup
}

var mods []*module

func Register(mi Module) {
	m := new(module)
	m.mi = mi
	m.closeSig = make(chan bool, 1)

	mods = append(mods, m)
}

func Init() {
	for i := 0; i < len(mods); i++ {
		mods[i].mi.OnInit()
	}

	for i := 0; i < len(mods); i++ {
		go run(mods[i])
	}
}

func Destroy() {
	for i := len(mods) - 1; i >= 0; i-- {
		m := mods[i]
		m.closeSig <- true
		m.wg.Wait()
		destroy(m)
	}
}

func PanicHandler() {
	if err := recover(); err != nil {
		exeName := os.Args[0]                                             //获取程序名称
		now := time.Now()                                                 //获取当前时间
		pid := os.Getpid()                                                //获取进程ID
		time_str := now.Format("20060102150405")                          //设定时间格式
		fname := fmt.Sprintf("%s-%d-%s-dump.log", exeName, pid, time_str) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）
		fmt.Println("dump to file ", fname)
		f, err := os.Create(fname)
		if err != nil {
			return
		}
		defer f.Close()
		f.WriteString(fmt.Sprintf("%v/r/n", err)) //输出panic信息
		f.WriteString("========/r/n")
		f.WriteString(string(debug.Stack())) //输出堆栈信息
	}
}

func run(m *module) {
	defer PanicHandler()
	m.wg.Add(1)
	m.mi.Run(m.closeSig)
	m.wg.Done()
}

func destroy(m *module) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	m.mi.OnDestroy()
}
