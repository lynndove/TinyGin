package tinygin

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
// trace 调用了 runtime.Callers(3, pcs[:])，Callers 用来返回调用栈的程序计数器
// 第 0 个 Caller 是 Callers 本身
// 第 1 个是上一层 trace
// 第 2 个是再上一层的 defer func
// 因此，为了日志简洁一点，跳过了前 3 个 Caller。
// 通过 runtime.FuncForPC(pc) 获取对应的函数
// 通过 fn.FileLine(pc) 获取到调用该函数的文件名和行号，打印在日志中
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// 使用 defer 挂载上错误恢复的函数
// 在这个函数中调用 *recover()*，捕获 panic
// 并且将堆栈信息打印在日志中，向用户返回 Internal Server Error
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
