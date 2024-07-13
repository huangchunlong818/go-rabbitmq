package errorx

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// Errors 接口定义了错误处理方法
type Errors interface {
	DoError(string)
}

// 默认的错误处理器
type defaultError struct{}

func (d *defaultError) DoError(msg string) {
	fmt.Println("Error :", msg)
}

// 全局变量保存当前的错误处理器
var errorHandler Errors = &defaultError{}

// SetErrorHandler 设置自定义的错误处理器
func SetErrorHandler(e Errors) {
	if e != nil {
		errorHandler = e
	}
}

// HandleError 调用当前的错误处理器处理错误
func HandleError(msg string) {
	errorHandler.DoError(msg)
}

// 推送错误消息，记录堆栈消息，
func ErrorPush(errMsg string) {
	// 获取堆栈跟踪信息
	buf := make([]byte, 4096)
	runtime.Stack(buf, false)
	stackTrace := string(buf)

	// 分析堆栈跟踪信息，提取错误文件信息
	tranceMsg := ""
	var tranceMsgs []string
	lines := strings.Split(stackTrace, "\n")
	i := 1
	for _, line := range lines {
		if strings.Contains(line, "rabbitmq") {
			tranceMsgs = append(tranceMsgs, strconv.Itoa(i)+"： "+line)
			i++
		}
	}
	if len(tranceMsgs) > 0 {
		tranceMsg = strings.Join(tranceMsgs, "\n")
	}

	HandleError("错误信息: " + errMsg + "\n堆栈信息：\n" + tranceMsg)
}
