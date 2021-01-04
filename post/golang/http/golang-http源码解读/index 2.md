
# 重要接口

|接口名称|简略含义|
|:-:|:-:|
|Handler|处理请求接口|
|ResponseWriter|返回响应接口|
|ServeMux|路由|
|Conn|tcp连接|

## Handler

        type Handler interface {
            ServeHTTP(ResponseWriter, *Request)
        }

&emsp;&emsp;Handler接口是一个逻辑接口，所有的请求都会进入该接口的实例进行逻辑处理。从它的名称就能猜出它的大概含义。

## ResponseWriter

        // ResponseWriter 的作用是被Handler作为形参调用来组装返回的Response
        type ResponseWriter interface {
            // 改写返回的Head
            Header() Header

            // 写Body体
            Write([]byte) (int, error)

            // 写入Head Code
            WriteHeader(int)
        }

## response 

&emsp;&emsp;实现ResponseWriter接口的默认结构体。

        type response struct {
            // 保存http连接信息
            conn    *conn

            // 请求信息
            req     *Request

            // 是否启用chunk
            chunking    bool

            // header是否已经执行过写操作
            wroteHeader bool

            // 100 Continue response was written
            wroteContinue    bool

            // 返回的header对象
            header Header

            // Body字节数
            written int64

            // Content 长度
            contentLength   int64

            // http code
            status  int

            // 当没有设置Content-Type的时候，开始sniff能根据HTTP body来确定Content-Type
            needSniff   bool

            // 是否保持长连接，如果请求中connection设置了keep-alive,该字段就是false
            closeAfterReply bool

            // 当requestBody太大时，response会返回411，并关闭连接
            requestBodyLimitHit bool
        }

## HandlerFunc

        type HandlerFunc func(ResponseWriter, *Request)

        func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
            f(w, r)
        }

&emsp;&emsp;这就是往期文章中所描述的典型的[接口型函数](https://github.com/facedamon/blob/post/golang/golang接口型函数/)。如果你已经理解了接口型函数，那么你在http包中会看到很多这样的使用技巧。

        func NotFound(w ResponseWriter, r *Request) {
            Error(w, "404 page not found", StatusNotFound)
        }

        func NotFoundHandler() Handler {
            return HandlerFunc(NotFound)
        }

> 接下来我们来看看Server.go

# ServeMux 路由

        type ServeMux struct {
            // 由于请求涉及到并发，所以这里有了读写锁
            mu sync.RWMutx
            // 简单的map路由规则 string就是路由表达式
            m map[string]muxEntry
        }

        type muxEntry struct {
            // 是否精确匹配
            explicit bool
            // 映射到哪个Handler接口
            h Handler
        }

&emsp;&emsp;它就是http包的路由规则器，你可以在ServeMux中注册你的路由规则，所以开源社区出现了一大批的http rounter框架，它们就是做的这件事。
&emsp;&emsp;当一个request请求进来时，server会依次根据ServeMux.m中的string路由表达式来匹配，如果找到了匹配的muxEntry，就取出h， 调用h的S而v俄HTTP来u装Response。

        // 根据path获取Handler
        func (mux *ServeMux) match(path string) Handler
        // 根据Request获取Handler，内部实现调用match
        func (mux *ServeMux) handler(r *Request) Handler
        // 说明ServeMux也是一个Handler，内部实现调用handler
        func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
        // 注册handler方法
        func (mux *ServeMux) Handle(pattern string, handler Handler)
        // 接口型函数
        func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))

## DefaultServeMux 

&emsp;&emsp;如果我们没有自定义的ServeMux, 系统默认使用该ServeMux。在http包外层中提供的几个方法, 内部就是调用ServeMux结构方法。

        func Handle(pattern string, handler Handler) {DefaultServeMux.Handle(pattern, handler)}
        func HandleFunc(pattern string, handler func(ResponseWriter, 8Request)) {
            DefaultServeMux.HandleFunc(pattern, handler)
        }

> 来看看Server结构

# Server

        type Server struct {
            // 监听地址和端口
            Addr    string

            // 路由，如果为空，则使用默认路由
            Handler Handler
            ReadTimeout time.Duration
            WriteTimeout    time.Duration

            // 请求头的最大长度
            MaxHeaderBytes int

            // 配置TLS
            TLSConfig   *tls.Config
        }

&emsp;&emsp;Server提供的方法：

        // 对某个端口监听，内部就是使用for死循环开启协程进行accept处理
        func (srv *Server) Serve(l net.Listener) error

        // 开启http server服务，内部调用Serve
        func (srv *Server) ListenAndServe() error

        // 开启https server， 内部调用Serve
        func (srv *Server) LisAndServeTLS(certFile, keyFile string) error

&emsp;&emsp;Http包也直接提供了对外方法，实际上内部就是实例化一个Server，然后调用ListenAndServe方法。

        func ListenAndServe(addr string, handler Handler) error
        func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler) error

# 案例解析

        func main() {
            http.HandleFunc("/hello", HelloServer)
            err := http.ListenAndServe(":8080", nil)
            if err != nil {
                log.Fatal("ListenAndServe: ", err)
            }
        }
        
        func HelloServer(w http.ResponseWriter, req *http.Request) {
            w.Write([]byte("hello\n"))
        }

- 流程分析
- 首先调用Http.HandleFunc
  - 1. 调用DefaultServeMux.HandleFunc
  - 2. 调用DefaultServeMux.Handle
  - 3. 往map[string]muxEntry中注册handle路由规则
- 其次调用http.ListenAndServe
  - 1. 实例话Server
  - 2. 调用Server.ListenAndServe()
  - 3. 调用net.Listen("tcp", addr)监听端口
  - 4. 启动for循环，在循环体中accept请求
  - 5. 对每个请求水花一个Conn，并开启goroutine为这个请求进行服务go c.serve()
  - 6. 读取每个请求的内容w, err := c.readRequest()
  - 7. 判断header是否为空，如果没有设置handler，handler就设置为DefaultServeMux
  - 8. 调用handler的ServeHttp
  - 9. 调用mux.handler(r).ServeHTTP
  - 10. 选择handler
    - A 判断是否有路由能满足这个request(循环遍历ServerMux的muxEntry)
    - B 如果有路由满足，调用该路由handler.ServeHttp
    - C 如果没有路由满足，调用NotFoundHandler.ServeHttp
