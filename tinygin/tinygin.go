package tinygin

import (
	"fmt"
	"net/http"
)

// HandlerFunc 这是提供给框架用户的，用来定义路由映射的处理方法
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Engine 添加了一张路由映射表router，key 由请求方法和静态路由地址构成
// 例如GET-/、GET-/hello、POST-/hello
// 这样针对相同的路由，如果请求方法不同,可以映射不同的处理方法(Handler)
// value 是用户映射的处理方法
type Engine struct {
	router map[string]HandlerFunc
}

// New is the constructor of gin.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET defines the method to add GET request
// 当用户调用(*Engine).GET()方法时
// 会将路由和处理方法注册到映射表 router 中

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
// (*Engine).Run()方法，是 ListenAndServe 的包装
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Engine实现的 ServeHTTP 方法的作用
// 解析请求的路径，查找路由映射表
// 如果查到，就执行注册的处理方法
// 如果查不到，就返回 404 NOT FOUND
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		// 设置返回码
		// w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
