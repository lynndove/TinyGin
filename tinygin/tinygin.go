package tinygin

import (
	"log"
	"net/http"
)

// HandlerFunc 这是提供给框架用户的，用来定义路由映射的处理方法
// type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type HandlerFunc func(*Context)

// Engine 添加了一张路由映射表router，key 由请求方法和静态路由地址构成
// 例如GET-/、GET-/hello、POST-/hello
// 这样针对相同的路由，如果请求方法不同,可以映射不同的处理方法(Handler)
// value 是用户映射的处理方法
type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		engine      *Engine
	}
	// go 中的嵌套类型，类似 Java/Python 等语言的继承
	// 这样 Engine 就可以拥有 RouterGroup 的属性了
	Engine struct {
		*RouterGroup
		// router map[string]HandlerFunc
		router *router
		groups []*RouterGroup // store all groups
	}
)

// New is the constructor of gin.Engine
func New() *Engine {
	// return &Engine{router: make(map[string]HandlerFunc)}
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// create a new RouterGroup
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	//key := method + "-" + pattern
	//engine.router[key] = handler
	//  group.prefix + prefix 的方式 group 初始化时已经拼接了完整的 prefix
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
// 当用户调用(*Engine).GET()方法时
// 会将路由和处理方法注册到映射表 router 中
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
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
	//key := req.Method + "-" + req.URL.Path
	//if handler, ok := engine.router[key]; ok {
	//	handler(w, req)
	//} else {
	//	// 设置返回码
	//	// w.WriteHeader(http.StatusNotFound)
	//	fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	//}
	// 在调用 router.handle 之前, 构造了一个 Context 对象
	c := newContext(w, req)
	engine.router.handle(c)
}
