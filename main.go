package main

import (
	"TinyGin/tinygin"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// type Engine struct{}

//func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	switch req.URL.Path {
//	case "/":
//		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
//
//	case "/hello":
//		for k, v := range req.Header {
//			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
//		}
//	default:
//		fmt.Fprintf(w, "404 NOTFOUND: %s\n", req.URL)
//	}
//}

func onlyForV2() tinygin.HandlerFunc {
	return func(c *tinygin.Context) {
		// start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	//http.HandleFunc("/", indexHandler)
	//http.HandleFunc("/hello", helloHandler)
	// 使用New()创建 gin 的实例

	r := tinygin.New()

	r.Use(tinygin.Logger(), tinygin.Recovery()) // global middleware
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Lynn", Age: 20}
	stu2 := &student{Name: "Linxi", Age: 30}
	r.GET("/", func(c *tinygin.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", tinygin.H{
			"title":  "tinygin",
			"strArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/data", func(c *tinygin.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", tinygin.H{
			"title": "tinygin",
			"now":   time.Date(2023, 3, 2, 9, 0, 0, 0, time.UTC),
		})
	})

	//r.GET("/", func(c *tinygin.Context) {
	//	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//})

	// 使用 GET()方法添加路由
	//r.GET("/", func(w http.ResponseWriter, req *http.Request) {
	//	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	//})
	//r.GET("/", func(c *tinygin.Context) {
	//	c.HTML(http.StatusOK, "<h1>Hello Go</h1>")
	//})

	//r.GET("/index", func(c *tinygin.Context) {
	//	c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	//})

	// Engine 拥有 RouterGroup 的属性
	//v1 := r.Group("/v1")
	//v1.GET("/", func(c *tinygin.Context) {
	//	c.HTML(http.StatusOK, "<h1>Hello Go</h1>")
	//})
	//v1.GET("/hello", func(c *tinygin.Context) {
	//	// expect /hello?name=tinygin
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	//})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	// // expect /hello/tinygin
	v2.GET("/hello/:name", func(c *tinygin.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	v2.POST("/login", func(c *tinygin.Context) {
		c.JSON(http.StatusOK, tinygin.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/panic", func(c *tinygin.Context) {
		names := []string{"tinygin"}
		c.String(http.StatusOK, names[100])
	})

	//r.GET("/hello", func(c *tinygin.Context) {
	//	// expect /hello?name=tinggin
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	//})
	//
	//r.GET("/hello/:name", func(c *tinygin.Context) {
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	//})
	//
	//r.GET("/assets/*filepath", func(c *tinygin.Context) {
	//	c.JSON(http.StatusOK, tinygin.H{"filepath": c.Param("filepath")})
	//})
	//
	//r.POST("/login", func(c *tinygin.Context) {
	//	c.JSON(http.StatusOK, tinygin.H{
	//		"username": c.PostForm("username"),
	//		"password": c.PostForm("password"),
	//	})
	//})

	//r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
	//	for k, v := range req.Header {
	//		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	//	}
	//})

	// 使用Run()启动Web服务
	r.Run(":9999")
	// engine := new(Engine)
	// main 函数的最后一行，是用来启动 Web 服务的
	// 第一个参数是地址，:9999表示在 9999 端口监听
	// 第二个参数则代表处理所有的HTTP请求的实例，nil 代表使用标准库中的实例处理
	// 第二个参数，是我们基于net/http标准库实现Web框架的入口
	// log.Fatal(http.ListenAndServe(":9999", engine))
}

// handler echoes r.URL.Path
//func indexHandler(w http.ResponseWriter, req *http.Request) {
//	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
//}
//
//// handler echoes r.URL.Header
//func helloHandler(w http.ResponseWriter, req *http.Request) {
//	for k, v := range req.Header {
//		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
//	}
//}
