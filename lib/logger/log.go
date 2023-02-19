package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func Init(path string) {

}

// Info 常规信息的输出，程序运行正常
func Info(msg ...string) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("[I]%s \"%s:%d\" %s\n", time.Now().Format("2006-01-02 15:04:05.999999"), file, line, strings.Join(msg, ""))
}

// Warning 警告信息的输出，重要，需要尽快去查看，但不需要立刻终止程序
func Warning(msg ...string) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("[W]%s \"%s:%d\" %s\n", time.Now().Format("2006-01-02 15:04:05.999999"), file, line, strings.Join(msg, ""))
}

// Error 发生重大错误，程序无法运行下去，会调用os.Exit()终止程序；
// 对于调用的第三方包，若希望进行异常recover，也在recover后进行调用，以确保打印信息后退出
func Error(msg ...string) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("[E]%s \"%s:%d\" %s\n", time.Now().Format("2006-01-02 15:04:05.999999"), file, line, strings.Join(msg, ""))
	os.Exit(1)
}
