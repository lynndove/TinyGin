package tinygin

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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
		router        *router
		groups        []*RouterGroup     // store all groups
		htmlTemplates *template.Template // for html render 将所有的模板加载进内存
		funcMap       template.FuncMap   // for html render 所有的自定义模板渲染函数
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

	// 当我们接收到一个具体请求时，要判断该请求适用于哪些中间件
	// 单通过 URL 的前缀来判断
	// 得到中间件列表后，赋值给 c.handlers
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 在调用 router.handle 之前, 构造了一个 Context 对象
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

// Use 将中间件应用到某个 Group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	// StripPrefix返回一个处理程序
	// 该处理程序通过从请求URL的Path（如果设置了RawPath）中删除给定前缀
	// 并调用处理程序h来服务http请求
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static serve static files
// Static这个方法是暴露给用户的
// 用户可以将磁盘上的某个文件夹root映射到路由relativePath
// eg : r.Static("/assets", "/usr/tinygin/blog/static")
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// 框架的模板渲染直接使用了html/template提供的
// 设置自定义渲染函数funcMap
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 加载模板
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
