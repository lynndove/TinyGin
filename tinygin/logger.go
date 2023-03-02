package tinygin

import (
	"log"
	"time"
)

// Logger 中间件等待用户自己定义的 Handler处理结束后，再做一些额外的操作
// 例如计算本次处理所用时间等
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
