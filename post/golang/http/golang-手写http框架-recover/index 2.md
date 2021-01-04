
# 摘要

&emsp;&emsp;在go中，比较常见的错误处理方式是返回error，由调用者巨顶后续如何处理。但是如果是无法恢复的错误，可以手动触发panic，当然如果在运行程序过程中出现了类似数组下标越界的错误，panic也会被触发。panic会终止当前执行的程序。

		package main
		
		func main() {
		    fmt.Printf("before panic\n")
		    panic("crash")
		    fmt.Printf("after panic\n")
		}

		$ go run main.go
		before panic
		panic: crash
		exit status 2

&emsp;&emsp;下面是数组下标越界触发的panic

		package main
		
		func main() {
		    arr := []int{1,2,3}
		    fmt.Printf("%d\n", arr[4])
		}

		$ go run main.go
		panic: runtime error: index out of range [4] with length 3

## defer

&emsp;&emsp;panic会导致程序被终止，但是在退出前，会先处理完当前协程上已经defer的任务，执行完后再退出。**遵循LIFO原则**

		package main
		
		func main() {
		    defer func(){
		        fmt.Printf("defer func\n")
		    }()
		
		    arr := []int{1,2,3}
		    fmt.Printf("%d\n", arr[4])
		}

		$ go run main.go
		defer func
		panic: runtime error: index out of range [4] with length 3

## recover

&emsp;&emsp;go提供了recover函数，可以避免因为panic发生而导致整个程序终止，**recover函数只在defer中生效。**


		package main
		
		func rec() {
		    defer func(){
		        fmt.Printf("defer func\n")
		        if err := recover(); err != nil {
		            fmt.Printf("recover success\n")
		        }
		    }()
		
		    arr := []int{1,2,3}
		    fmt.Printf("%d\n", arr[4])
		    fmt.Printf("after panic")
		}
		
		func main() {
		    rec()
		    fmt.Printf("after recover\n")
		}

		$ go run main.go
		defer func
		recover success
		after recover

# Geew的错误处理机制

&emsp;&emsp;对于一个Http框架而言，错误处理机制是非常重要的。可能是框架本身没有晚辈的测试，导致再某些情况下出现空指针异常等情况。也有可能用户不正确的参数，触发了某些异常，例如数组下标越界，空指针等。如果因为这些原因系统宕机，必然是不可接收的。

&emsp;&emsp;我们再geew中添加一个非常简单的错误处理机制，即在此类错误发生时，向用户返回`Internal Server Error`，并且再日志中打印必要的错误信息，方便进行错误定位。之前我们实现了中间件机制，错误处理也可以作为一个中间件，增强geew框架的能力。

		func Recovery() HandlerFunc {
		    return func(c *Context) {
		        defer func(){
		            if err := recover(); err != nil {
		                message := fmt.Sprintf("%s", err)
		                log.Printf("%s\n\n", trace(message))
		                c.JSON(http.StatusInternalServerError, geew.H{
		                    code: http.StatusInternalServerError,
		                    msg: message,
		                    data: nil,
		                })
		            }
		        }()
		        c.Next()
		    }
		}
		
		// print stack trace for debug
		func trace(message string) string {
		    var pcs [32]uintptr
		    n := runtime.Callers(3, pcs[:])
		
		    var str strings.Builder
		    str.WriteString(message + "\nTraceback: ")
		    for _, pc := range pcs[:n] {
		        fn := runtime.FuncForPC(pc)
		        file, line := fn.FileLine(pc)
		        str.WriteString(fmt.Sprintf("\n\t%s:%d", filem line))
		    }
		    return str.String()
		}

&emsp;&emsp;在trace中，调用了`runtime.Callers(3, pcs[:])`，Callers用来返回调用栈的程序计数器，第0个Caller是Callers本身，第1个是上层trace，第2个是再上一层的`defer func`。因此，为了日志简洁，我们跳过了前3个Caller。

## geew改变

		func Default() *Engine {
		    e := New()
		    e.Use(TimeOperation(), Recovery())
		    return e
		}


## Test

		// example_test.go
		package example
		
		func TestRecovery() {
		    r := geew.Default()
		    // index out of range for testing Recovery()
			r.GET("/panic", func(c *geew.Context) {
				names := []string{"geew"}
				c.String(http.StatusOK, names[100])
			})
		
			r.Run(":9999")
		}
