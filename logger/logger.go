package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type sLogLine struct {
	lvl LogLvl
	msg string
}

type LogLvl int

const (
	LVL_FATAL LogLvl = iota
	LVL_ERROR
	LVL_WARN
	LVL_APP
	LVL_DEBUG
	LVL_INTERNAL
)

func (d LogLvl) String() string {
	return [...]string{"FATAL", "ERROR", "WARN", "APP", "DEBUG", "INTERNAL"}[d]
}

var (
	Vbs  LogLvl = LVL_ERROR
	done chan struct{}
	logs chan sLogLine
)

func Start(v LogLvl) {

	done = make(chan struct{})
	logs = make(chan sLogLine, 1000)

	log.SetFlags(0)

	Vbs = LogLvl(v)

	go monitorLoop()
}

func Stop() {

	close(logs)
	<-done
}

func Log(l LogLvl, m string) {

	var c string

	pc, _, ln, ok := runtime.Caller(1)
	if ok {
		details := runtime.FuncForPC(pc)
		c, ln = details.FileLine(pc)
		path, _ := os.Getwd()
		c = strings.TrimPrefix(c, path+"/")
	} else {
		c = "unknown"
	}

	if l <= Vbs {
		if Vbs >= LVL_DEBUG {
			m = fmt.Sprintf("[%s:%v] %s", c, ln, m)
		}
		logs <- sLogLine{lvl: l, msg: m}
	}
}

func monitorLoop() {

	for f := range logs {
		log.Print(f.msg)
	}

	close(done)
}
