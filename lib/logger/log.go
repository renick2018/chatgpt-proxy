package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

func Init(path string) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		//fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile|log.Ldate|log.Ltime)
}

// Info 常规信息的输出，程序运行正常
func Info(msg ...string) {
	print("I", msg...)
}

// Warning 警告信息的输出，重要，需要尽快去查看，但不需要立刻终止程序
func Warning(msg ...string) {
	print("W", msg...)
}

// Error 发生重大错误，程序无法运行下去，会调用os.Exit()终止程序；
// 对于调用的第三方包，若希望进行异常recover，也在recover后进行调用，以确保打印信息后退出
func Error(msg ...string) {
	print("E", msg...)
	os.Exit(1)
}

func print(level string, msg ...string)  {
	_, file, line, _ := runtime.Caller(2)
	text := fmt.Sprintf("[%s]%s \"chatgpt-proxy%s:%d\" %s", level, time.Now().Format("2006-01-02 15:04:05.999999"), strings.Split(file, "chatgpt-proxy")[1], line, strings.Join(msg, ""))
	log.Println(text)
	fmt.Println(text)
}
