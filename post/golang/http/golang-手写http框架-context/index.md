
 # 摘要

 - 将路由`router`独立出来，方便之后增强。
 - 设计上下文`context`，封装Request和Response，提供对JSON、HTML等返回类型的支持。

# 使用效果

&emsp;&emsp;为了展示第二天的成果，我们看一看在使用时的效果。

		func main() {
		    r := geew.New()
		    r.GET("/", func(c *geew.Context){
		        c.HTML(http.StatucOK, "<h1>Hello Gee</h1>")
		    })
		    r.GET("/hello", func(c *geew.Context){
		        c.String(http.StatusOK, "htllo %s, you are at %s\n", c.Query("name"), c.Path)
		    })
		
		    r.POST("/login", func(c *geew.Context){
		        c.JSON(http.StatusOK, geew.H{
		            "username": c.PostForm("username"),
		            "password": c.PostForm("password")
		        })
		    })
		
		    r.Run(":9999")
		}

- `Handler`的参数变成了`geew.Context`，提供了查询Query/PostForm参数的功能。
- `geew.COntext`封装了`HTML/String/JSON`函数，能够快速构造HTTP响应。

# 设计Context

- 必要性

&emsp;&emsp;对Web服务来说，无非是根据请求`*http.Request`，构造响应`http.ResponseWriter`。但是这两个对象提供的接口粒度太细，比如我们要构造一个完整的响应，需要考虑消息头Header和消息题Body，而Header包含了状态码StatusCode，消息类型ContentType等几乎每次请求都需要设置的信息。因此，如果不进行有效的封装，那么框架的用户将需要写大量重复繁杂的代码，而且容易出错。针对常用场景，能够高效地构造出HTTP响应是一个好框架必须考虑的点。

&emsp;&emsp;用返回JSON数据做比较，感受下封装前后的差距。

- 封装前

		obj = map[string]interface{}{
		    "name": "gee",
		    "password": "1234",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(obj); err != nil {
		    http.Error(w, err.Error(), 500)
		}

- 封装后

		c.JSON(http.StatucOK, geew.H{
		   "name": c.PostForm("name"),
		    "password": c.PostForm("password"),
		})

&emsp;&emsp;针对使用场景，封装`*http.Request`和`http.ResponseWriter`的方法，简化相关接口的调用，只是设计Context的原因之一。对于框架来说，还需要支撑额外的功能。例如，将来解析动态路由`/hello:name`，参数`:name`的值放在哪呢？再比如，框架需要支持中间件，那中间件产生的信息放在哪？Context随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应该由Context承载。因此，设计Context结构，扩展性和复杂性留在了内部，而对外简化接口。路由的处理函数，以及将要实现的中间件，参数都统一使用Context实例，Context就像一次会话的百宝箱，可以找到任何东西。

## context.go

		// H is aslias of JSON
		type H map[string]interface{}
		
		type Context struct {
		    Writer http.ResponseWriter
		    Req *http.Request
		    Path string
		    Method string
		    StatusCode int
		
		    // queryCache use url.ParseQuery cached the param query result from c.Req.URL.Query
		    queryCache url.Values
		    //formCache use url.ParseQuery cached PostForm contains the parsed form data from POST,
		    //PUT body parameters.
		    formCache url.Values
		}
		
		func newContext()(w http.ResponseWriter, req *http.Request) *Context {
		    return &Context{
		        Writer: w,
		        Req: req,
		        Path: req.URL.Path,
		        Method: req.Method,
		        quaryCache: nil,
		        formCache: nil,
		    }
		}
		
		func (c *Context) initQueryCache() {
		        if c.queryCache == nil {
		                c.queryCache = c.Req.URL.Query()
		        }
		}
		
		func (c *Context) GetQueryArray(key string) ([]string, bool) {
		        c.initQueryCache()
		        if values, ok := c.queryCache[key]; ok && len(values) > 0 {
		                return values, true
		        }
		        return []string{}, false
		}
		
		// QueryArray returns a slice of strings for given query key
		// the length od the slice depends on the number of params with the given key
		func (c *Context) QueryArray(key string) []string {
		        values, _ := c.GetQueryArray(key)
		        return values
		}
		
		func (c *Context) GetQuery(key string) (string, bool) {
		        if values, ok := c.GetQueryArray(key); ok {
		                return values[0], ok
		        }
		        return "", false
		}
		
		func (c *Context) DefaultQuery(key, defaultValue string) string {
		        if value, ok := c.GetQuery(key); ok {
		                return value
		        }
		        return defaultValue
		}
		
		func (c *Context) Query(key string) string {
		        value, _ := c.GetQuery(key)
		        return value
		}
		
		func (c *Context) initFormCache() {
		        if c.formCache == nil {
		                c.formCache = make(url.Values)
		                r := c.Req
		                if err := r.ParseMultipartForm(MaxMultipartMemory); err != nil {
		                        if err != http.ErrNotMultipart {
		                                L.Error("error on parse multipart form array: %v", err)
		                        }
		                }
		                c.formCache = r.PostForm
		        }
		}
		
		// GetPostFormArray returns a slice of strings for a given form key
		func (c *Context) GetPostFormArray(key string) ([]string, bool) {
		        c.initFormCache()
		        if values := c.formCache[key]; len(values) > 0 {
		                return values, true
		        }
		        return []string{}, false
		}
		
		// PostFormArray returns a slice of strings for a given form key
		// the length of slice depends on the number of params with the given key
		func (c *Context) PostFormArray(key string) []string {
		        values, _ := c.GetPostFormArray(key)
		        return values
		}
		
		func (c *Context) GetPostForm(key string) (string, bool) {
		        if values, ok := c.GetPostFormArray(key); ok {
		                return values[0], ok
		        }
		        return "", false
		}
		
		func (c *Context) DefaultPostForm(key, defaultValue string) string {
		        if value, ok := c.GetPostForm(key); ok {
		                return value
		        }
		        return defaultValue
		}
		
		// PostForm returns the specified key from a POST urlencoded from or multipart form
		func (c *Context) PostForm(key string) string {
		        value, _ := c.GetPostForm(key)
		        return value
		}
		
		func (c *Context) Param(key string) string {
		        c.Req.ParseMultipartForm(MaxMultipartMemory)
		        v, _ := c.Params[key]
		        return v
		}
		
		func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
		        if c.Req.MultipartForm == nil {
		                if err := c.Req.ParseMultipartForm(MaxMultipartMemory); err != nil {
		                        return nil, err
		                }
		        }
		        f, fh, err := c.Req.FormFile(name)
		        if err != nil {
		                return nil, err
		        }
		        f.Close()
		        return fh, err
		}
		
		func (c *Context) MultipartForm() (*multipart.Form, error) {
		        err := c.Req.ParseMultipartForm(MaxMultipartMemory)
		        return c.Req.MultipartForm, err
		}
		
		func (c *Context) Status(code int) {
		    c.StatusCode = code
		    c.Writer.WriteHeader(code)
		}
		
		func (c *Context) SetHeader(key string, value string) {
		    c.Writer.Header().Set(key, value)
		}
		
		func (c *Context) String(code int, format string, values ...interface{}) {
		    c.SetHeader("Content-Type", "text/plain")
		    c.Status(code)
		    c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
		}
		
		func (c *Context) JSON(code int, obj interface{}) {
		    c.SetHeader("Content-Type", "application/json")
		    c.Status(code)
		    encoder := json.NewENcoder(c.Writer)
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
		    c.Writer.write([]byte(html))
		}

- 代码最开头，给`map[string]interface{}`起了一个别名`geew.H`，构建JSON数据时，显的更简洁。
- `Context`目前只包含了http.ResponseWriter和*http.Request, 另外提供了对Method和Path这两个常用属性的直接访问.
- 对Form的封装(采用了cache缓存), 供了访问Query和PostForm参数的方法。
- 提供了快速构造String/Data/JSON/HTML响应的方法。

## router.go

&emsp;&emsp;我们将路由相关的方法和结构提取出来，放到一个新的文件中的router.go，方便我们下一次对router的功能进行增强，例如提供动态路由的支持。router的handle方法做了一个细微的调整，即handler的参数，变成了Context。

		package geew
		
		import (
		        "log"
		        "net/http"
		)
		
		type router struct {
		        handlers map[string]HandlerFunc
		}
		
		func newRouter() *router {
		        return &router{handlers: make(map[string]HandlerFunc)}
		}
		
		func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
		        log.Printf("Route %4s - %s", method, pattern)
		        key := method + "-" + pattern
		        r.handlers[key] = handler
		}
		
		func (r *router) handle(c *Context) {
		        key := c.Method + "-" + c.Path
		        if handler, ok := r.handlers[key]; ok {
		                handler(c)
		        } else {
		                c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		        }
		}

## geew.go (入口)

		package geew
		
		import (
		        "net/http"
		)
		
		// HandlerFunc defines the request handler used by geew
		//type HandlerFunc func(http.ResponseWriter, *http.Request)
		type HandlerFunc func(*Context)
		
		// Engine implement the interface of ServeHTTP
		type Engine struct {
		        //router map[string]HandlerFunc
		        router *router
		}
		
		// New is the constructor of geew.Engine
		func New() *Engine {
		        return &Engine{router: newRouter()}
		}
		
		func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
		        e.router.addRoute(method, pattern, handler)
		}
		
		// GET defines the method to add GET request
		func (e *Engine) GET(pattern string, handler HandlerFunc) {
		        e.addRoute("GET", pattern, handler)
		}
		
		// POST defines the method to add POST request
		func (e *Engine) POST(pattern string, handler HandlerFunc) {
		        e.addRoute("POST", pattern, handler)
		}
		
		// Run defines the method to start up a http server
		func (e *Engine) Run(addr string) error {
		        return http.ListenAndServe(addr, e)
		}
		
		// Implement the HTTP ServeHTTP
		func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
		        c := newContext(w, req)
		        e.router.handle(c)
		}

&emsp;&emsp;将`router`相关的代码独立后，`geew.go`简单了不少。最终要的还是通过实现`ServeHTTP`接口，接管了所有的HTTP请求。相比第一天的代码，这个方法也有细微的调整，在调用router.handle之前，构造了一个`Context`对象。这个对象目前还非常简单，仅仅是包装了原来的两个参数，之后我们会慢慢地给Context插上翅膀。

&emsp;&emsp;如何使用，运行`go test -v -run TestDay2 day2_test.go`

		package example
		
		import (
		        "geew"
		        "net/http"
		        "testing"
		)
		
		func TestDay2(t *testing.T) {
		        r := geew.New()
		        r.GET("/", func(c *geew.Context) {
		                c.HTML(http.StatusOK, "<h1>Hello Geew</h1>")
		        })
		        r.GET("/hello", func(c *geew.Context) {
		                c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Query("name"), c.Path)
		        })
		        r.POST("/login", func(c *geew.Context) {
		                c.JSON(http.StatusOK, geew.H{
		                        "username": c.PostForm("username"),
		                        "password": c.PostForm("password"),
		                })
		        })
		
		        r.Run(":9999")
		}

		$ curl -i http://localhost:9999/
		HTTP/1.1 200 OK
		Date: Mon, 12 Aug 2019 16:52:52 GMT
		Content-Length: 18
		Content-Type: text/html; charset=utf-8
		<h1>Hello Gee</h1>
		
		$ curl "http://localhost:9999/hello?name=geektutu"
		hello geektutu, you're at /hello
		
		$ curl "http://localhost:9999/login" -X POST -d 'username=geektutu&password=1234'
		{"password":"1234","username":"geektutu"}
		
		$ curl "http://localhost:9999/xxx"
		404 NOT FOUND: /xxx
