package tinygin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 给map[string]interface{}起了一个别名tinygin.H，构建JSON数据时，显得更简洁
type H map[string]interface{}

// Context req 是结构体，用指针可以节省内存，Writer 是一个接口类型，不能用指针
// 如何实现类似于gin中的c.Abort()这种中间件的退出机制呢？
// context 中维护一个状态值，调用 c.Abort() 改变状态，循环时检查状态，发现已经中止停止循环即可
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// index是记录当前执行到第几个中间件
// 当在中间件中调用Next方法时，控制权交给了下一个中间件
// 直到调用到最后一个中间件，然后再从后往前，调用每个中间件在Next方法之后定义的部分
func (c *Context) Next() {
	c.index++
	size := len(c.handlers)
	for ; c.index < size; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) PostForm(key string) string {
	// FormValue 返回查询的命名组件的第一个值
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 快速构造String/Data/JSON/HTML响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	// returns a new encoder that writes to w
	encoder := json.NewEncoder(c.Writer)
	// Encode writes the JSON encoding of v to the stream
	// followed by a newline character
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// Fail 短路中间件，如果使用 后续的中间件和handler就直接跳过了
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
