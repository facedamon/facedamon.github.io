# 资源管理与出错处理
## defer调用
- 确保在函数结束时调用

		package main
		
		func tryDefer() {
		    defer fmt.Println(1)
		    defer fmt.Println(2)
		    fmt.Println(3)
		    panic("error occurr") //触发return
		    fmt.Println(4)
		}
		
		func writeFile(filename string) {
		    file, err := os.OpenFile(filename, os.O_EXCL|os.O_CREATE, 0666) 
		
		    if err != nil {
		        panic(err)
		    }
		    defer file.Close()
		
		    writer := bufio.NewWriter(file)
		    defer writer.Flush()
		
		    f := fib.Fibonacci()
		    for i := 0; i < 20; i++ {
		        fmt.Fprintln(writer, f())
		    }
		}
		
		func main() {
		    tryDefer() // 3 2 1 panic
		    writeFile("fib.txt")
		}
> defer

- 确保调用在函数结束时发生
- 参数在defer语句时计算
- defer先进后出

> 何时使用defer

1. Open/Close
2. Lock/UnLock
3. PrintHeader/PrintFooter

## 错误处理

		if err != nil {
		    if pathError, ok := err.(*os.PathError); !ok {
		        panic(err)
		    } else {
		        fmt.Printf("%s, %s, %s\n",pathError.Op,pathError.Path,pathError.Err)
		    }
		    return
		}

## 服务器统一的错误处理逻辑
&emsp;&emsp;以打印文件的web server为例
		package main
		
		import (
		    "net/http"
		)
		
		//  让handler返回一个error, 往上抛
		type appHandler func(writer http.ResponseWriter, request *http.Request) error
		
		// 然后包装该error，返回一个http需要的接口类型
		func errWrapper(handler appHandler) func(http.ResponseWriter, *http.Request) {
		    return func(writer httpResponseWriter, request *http.Request) {
		        err := handler(writer, request)
		        if err != nil {
		            log.Printf("Error handling request: %s", err.Error)
		            code := http.StatusOk
		            switch {
		                case os.IsNotExists(err): code = http.StatusNotFound
		                case os.IsPermission(err): code = http.StatusForbidden
		                default: code = http.StatusInternalServerError
		            }
		            http.Error(writer, http.StatusText(code), code)
		        }
		    }
		}
		
		func main() {
		    http.HandleFunc("/list/", errWrapper(handleFileLisr))
		
		    err := http.ListenAndServe(":8888", nil)
		    if err != nil {
		        panic(err)
		    }
		}
		package handle
		
		import (
		    "os"
		    "io/ioutil"
		)
		
		func HandleFileList(writer http.ResponseWriter, request *httpRequest) error{
		    path := request.URL.Path[len("/list/"):]
		    file, err := os.Open(path)
		    if err != nil {
		        return err
		    }
		    defer file.Close()
		
		    all, err := ioutil.ReadAll(file)
		    if err != nil {
		        return err
		    }
		
		    writer.Write(all)
		    return nil
		}

## panic和recover

### panic (抓狂)
- 停止当前函数执行
- 一直向上返回，执行每一层的defer
- 如果没有遇见recover，程序退出
  
### recover
- 仅在defer调用中使用
- 获取panic的值
- 如果无法处理，可重新panic

		package main
		
		func main() {
		    tryRecover()
		}
		
		func tryRecover() {
		    // 调用匿名函数
		    defer func() {
		        r := recover()
		        if err, ok := r.(error); ok {
		            fmt.Println("Error occurred:", err)
		        } else {
		            panic(fmt.Sprintf("I dont`t know to do: %v", r))
		        }
		    }()
		    panic(errors.New("this is an error"))
		}

## 服务器统一出错处理

		const prefix = "/list/"
		
		type ue string
		
		func (e ue) Error() string {
		    return e.Message()
		}
		
		func (e ue) Message() string {
		    return string(e)
		}
		
		func HandleFileList(writer http.ResponseWriter, request *httpRequest) error{
		    if strings.Index(request.URL.Path, prefix) != 0 {
		        // 返回给用户看的信息
		        return ue("path must with " + prefix)
		    }
		    path := request.URL.Path[len("/list/"):]
		    file, err := os.Open(path)
		    if err != nil {
		        return err
		    }
		    defer file.Close()
		
		    all, err := ioutil.ReadAll(file)
		    if err != nil {
		        return err
		    }
		
		    writer.Write(all)
		    return nil
		}
		
		func errWrapper(handler appHandler) func(http.ResponseWriter, *http.Request) {
		    return func(writer httpResponseWriter, request *http.Request) {
		        // 捕获一下panic
		        defer func(){
		            if r := recover(); r != nil {
		                 log.Printf("Panic: %v", r)
		                 http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		            }    
		        }()
		
		        err := handler(writer, request)
		        if err != nil {
		            // User Error
		            log.Printf("Error handling request: %s", err.Error)
		            if userError, ok := err.(userError); ok {
		                http.Error(writer, userError.Message(), http.StatusBadRequest)
		                return
		            }
		            // System Error
		            code := http.StatusOk
		            switch {
		                case os.IsNotExists(err): code = http.StatusNotFound
		                case os.IsPermission(err): code = http.StatusForbidden
		                default: code = http.StatusInternalServerError
		            }
		            http.Error(writer, http.StatusText(code), code)
		        }
		    }
		}
		
		type userError interface {
		    error
		    Message() string    // 给用户看
		}
