package tinygin

import (
	"net/http"
	"strings"
)

// 我们将和路由相关的方法和结构提取出来
// 方便下一次对 router 的功能进行增强，例如提供动态路由的支持
// router 的 handle 方法作了一个细微的调整，即 handler 的参数，变成了 Context

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
// 一个部分只能有一个*
func parsePattern(pattern string) []string {
	s := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range s {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	// log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
} // 2.28

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
