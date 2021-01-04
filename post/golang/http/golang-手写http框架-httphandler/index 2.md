
# 摘要

- 简单介绍`net/http`库以及`http.Handler`接口
- 搭建Gee框架的雏形

## 标准库启动Web服务

&emsp;&emsp;Go语言内置了`net/http`库，封装了HTTP网络编接口的基础接口，我们实现的Gee Web框架便是基于`net/http`的。我们接下来通过一个例子，简单介绍一下这个库的使用。

		package main
		
		import (
		    "fmt"
		    "log"
		    "net/http"
		)
		
		// handler echoes r.URL.Path
		func indexHandler(w http.ResponseWriter, req *http.Request) {
		    fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
		}
		
		// handler echoes r.URL.Header
		func helloHandler(w http.ResponseWriter, req *http.Request) {
		    for k, v := range req.Header {
		        fmt.Fprintf(w, "Header[%q] = %q\n, k, v")
		    }
		}
		
		func main(){
		    http.HandleFunc("/", indexHandler)
		    http.HandleFunc("/hello", helloHandler)
		    log.Fatal(http.ListenAndServer(":9999", nil))
		}

&emsp;&emsp;我们设置了2两个路由，`/`和`/hello`，分别绑定`indexHandler`和`helloHandler`，根据不同的HTTP请求调用不同的处理函数。

&emsp;&emsp;用curl测试一下：

		$ curl http://127.0.0.1:9999/
		URL.Path = "/"
		
		$ curl http://127.0.0.1:9999/hello
		Header[Accept] = ["*/*"]
		Header["User-Agent"] = ["curl/7.54.0"]

&emsp;&emsp;`main`函数的最后一行，是用来启动Web服务的，第一个参数是地址，`:9999`表示在9999端口监听。二第二个参数则代表处理所有的HTTP请求的实例，`nil`代表使用标准库中的实例处理。第二个参数，则是我们基于`net/http`标准库实现Web框架的入口。

## 实现http.Handler接口

		package http
		
		type Handler interface {
		    ServeHTTP(w ResponseWriter, r *Request)
		}
		
		func ListenAndServe(address string, h Handler) error

&emsp;&emsp;第二个参数的类型是什么？通过查看`net/http`的源码可以发现，Handler是一个接口，需要实现方法ServeHTTP，也就是说，只要传入任何实现了ServeHTTP接口的实例，所有的HTTP请求，就都交给了该实例处理。马上来试一试。

		package main
		
		import (
		    "fmt"
		    "log"
		    "net/http"
		)
		
		// Engine is the uni handler for all request
		type Engine struct{}
		
		func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		    switch r.URL.Path {
		        case "/":
		            fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
		        case "/hello":
		            for k, v := range r.Header {
		                fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		            }
		        default:
		            fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
		    }
		}
		
		func main(){
		    e := new(Engine)
		    log.Fatal(http.ListenAndServe(":9999", e))
		}

&emsp;&emsp;我们定义了一个空的结构体`Engine`,实现了方法`ServeHTTP`。这个方法有两个参数，第二个参数是`Request`，该对象包含了该HTTP请求的所有信息，比如请求地址、Header和Body等信息；第一个参数是`ResponseWriter`，利用`ResponseWriter`可以构造针对该请求的响应。

&emsp;&emsp;在`main`函数中，我们给`ListenAndServe`方法的第二个参数传入了刚才创建的`Engine`实例。至此，我们走出了实现Web框架的第一步，即，将所有的HTTP请求转向了我们自己的处理逻辑。还记得吗，在实现`Engine`之前，我们调用`http.HandleFunc`实现了路由和Handler的映射，也就是只能针对具体的路由写处理逻辑。比如`/hello`。但是在实现Engine知乎，我们拦截了所有的HTTP请求，拥有了统一的控制入口。在这里我们可以自由定义路由映射的规则，也可以统一添加一些处理逻辑，例如：日志、异常处理等。

## geew雏形

&emsp;&emsp;接下来重新组织上面的代码，搭建出整个框架的雏形。最终的代码目录结构是这样的。

		geew/
		    |--example
		        |--example_test.go
		gee.go
		go.mod
		go.mod

### go.mod

		module geew
		
		go 1.14
		
		require geew v0.0.0

## example_test.go

		package example
		
		import (
		    "fmt"
		    "net/http"
		    "geew"
		    "testing"
		)
		
		func TestHandler(t *testing.T){
		    r := geew.New()
		    r.GET("/", func(w http.ResponseWriter, req *http.Request){
		        fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
		    })
		
		    r.GET("/hello", func(w http.ResponseWriter, req *http.Request){
		        for k, v := range req.Header {
		            fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		        }
		    })
		
		    r.Run(":9999")
		}

&emsp;&emsp;看到这里，如果你使用过`gin`框架的话，肯定会觉得无比亲切。`geew`框架的设计以及API均参考了`gin`。使用New创建gee实例，使用`GET`方法添加路由，最后使用`Run`启动Web服务。这里的路由，只是静态路由，不支持`/hello/:name`这样的动态路由，动态路由我们将在下一次实现。

## geew.go

		package geew
		
		import (
		    "fmt"
		    "net/http"
		)
		
		// HandlerFunc defines the request handler used by gee
		type HandlerFunc func(http.ResponseWriter, *http.Request)
		
		// Engine implement the interface of ServeHTTP
		type Engine struct {
		    router map[string]HandlerFunc
		}
		
		// New is the constructor of gee.Engine
		func New() *Engine {
		    return &Engine{router: make(map[string]HandlerFunc)}
		}
		
		func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
		    key := method + "-" + pattern
		    e.router[key] = handler
		}
		
		// GET defines the method to add GET request
		func (e *Engine) GET(pattern string, handler HandlerFunc) {
		    e.addRoute("GET", pattern, handler)
		}
		
		// POST defines the method to add POST request
		func (e *Engine) POST(pattern string, handler HandlerFunc) {
		    e.addRoute("POST", pattern, handler)
		}
		
		// Run defines the method to start a http server
		func (e *Engine) Run(addr string) (err error) {
		    return http.ListenAndServe(addr, e)
		}
		
		func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
		    key := req.Method + "-" + req.URL.Path
		    if handler, ok := e.router[key]; ok {
		        handler(w, req)
		    } else {
		        fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
		    }
		}

&emsp;&emsp;那么`geew.go`就是重头戏了，我们重点介绍一下这部分的实现。

- 首先定义了类型`HandlerFunc`，这是提供给框架用户的，用来定义路由映射的处理方法。我们在`Engine`中，添加了一张路由映射表`router`，key由请求方法和静态路由地址构成，例如`GET-/`、`GET-/hello`、`/POST-/hello`，这样针对相同的路由，如果请求方法不同，可以映射不同的处理方法(Handler)，value是用户映射的处理方法。
- 当用户调用`(*Engine).GET`方法时，会将路由和处理方法注册到映射表router中，`(*Engine).Run`方法是ListenAndServe的包装。
- `Engine`实现的ServeHTTP方法的作用就是，解析请求路径，查找路由映射，如果查到就执行注册处理方法。如果查不到，就返回`404 NOT FOUND`。

&emsp;&emsp;执行`go test -v -run TestHandler example_test.go`.再用curl 工具访问，结果与开始的一致。

		$ curl http://localhost:9999/
		URL.Path = "/"
		
		$ curl http://localhost:9999/hello
		Header["Accept"] = ["*/*"]
		Header["User-Agent"] = ["curl/7.54.0"]
		
		$ curl http://localhost:9999/world
		404 NOT FOUND: /world

&emsp;&emsp;至此，整个`geew`框架的原型已经出来了。实现了路由映射表，提供了用户注册静态路由的方法，包装了启动服务的函数。当然，到目前为止，我们还没有实现比net/http标准库更强大的能力，不用担心，很快就可以将动态路由、中间件等功能添加上去了。
