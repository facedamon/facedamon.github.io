
# http
- 使用http客户端发送请求
- 使用http.Client控制请求头部等
- 使用httputil简化工作

		package main
		
		func main(){
		    // 访问手机版imooc
		    r, err := http.NewRequest(http.MethodGet, "http://www.imooc.com", nil)
		    r.Header.Add("User-Agent", "Mozilla/5.0 (iphone;CPU iphone OS 10_3 like...)")
		    client := http.Client {
		        // 重定向                                 重定向地址记录via
		        CheckRedirect: func (req *http.Request, via[]*http.Request) error {
		            fmt.Println("Redirect: ", req)
		            return nil 
		        },
		    }
		    resp, err := client.Do(r)
		
		    // 访问网页版imooc
		    resp, err := http.Get("http://www.imooc.com")
		    if err != nil {
		        panic(err)
		    }
		    defer resp.Body.Close()
		
		    s, err := httputil.DumpResponse(resp, true)
		    if err != nil {
		        panic(err)
		    }
		
		    fmt.Printf("%s\n", s)
		}
# http服务器的性能分析
- import _ "net/http/pprof"
- 访问/debug/pprof/
- 使用go tool pprof分析性能

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/debugpprof.png)

&emsp;&emsp;上图是通过页面进入参数查看页面

		$ go tool pprof http://localhost:8888/debug/pprof/profile
		// 等待用户操作
		// 出现响应
		（pprof）web
		// 查看内存图
&emsp;&emsp;更多的选项可以查看pprof.go文件详细注释

# 其他标准库
- bufio
- log
- encoding/json
- regexp
- time
- strings/math/rand
- 查看文档
    - `godoc -http :8888`
    - https://studygolang.com/pkgdoc
