
# 摘要

&emsp;&emsp;分布式缓存需要实现节点间通信，建立基于HTTP的通信机制是比较常见的做法。如果一个节点启动了HTTP服务， 那么这个节点就可以被其它节点访问。今天我们就为单机节点搭建HTTP Server.

&emsp;&emsp;不与其它部分耦合，我们将这部分代码放在新的`http.go`文件中，当前的代码结构如下: 


    geec/
        |--lru/
            |--lru.go       // lru缓存淘汰策略
            |--lru_test.go
        |--byteview.go      // 缓存值的抽象
        |--cache.go         // 并发控制
        |--geec.go          // 服务与外部交互，控制缓存存储和获取的主流程
        |--http.go          // 提供被其它节点访问的能力
		

&emsp;&emsp;首先我们创建一个结构体`HTTPPool`，作为承载节点间HTTP通信的核心数据(包括服务端和客户端，今天只实现服务端)。

    package geec

    const defaultBasePath = "/_geec/"
    
    // HTTPPool impl PeerPicker for a pool of HTTP peers.
    type HTTPPool struct {
            // this peer`s base url, eg: "https://example.net:8080"
            self     string
            basePath string
     }
    
     // NewHTTPPool init an HTTP pool of peers.
     func NewHTTPPool(self string) *HTTPPool {
             return &HTTPPool{
                     self:     self,
                     basePath: defaultBasePath,
             }
     }

- `HTTPPool`只有2个参数，一个是self，用来记录自己的地址，包括主机名/IP和端口
- 另一个是basePath，作为节点间通讯地址的前缀，默认是`/_geec/`,那么`http://example.com/_geec/`开头的请求，就用语节点间的访问。因为一个主机上还可能承载其它的服务，加一段Path是一个好习惯。比如，大部分网站的API接口，一般以`/api`作为前缀。

&emsp;&emsp;接下来，实现最为核心的`ServeHTTP`方法
		
    package geec
    // Log info with server name
    func (p *HTTPPool) Log(format string, v ...interface{}) {
            log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
    }
    // ServeHTTP handle all http requests
    func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r http.Request) {
        if !strings.HasPrefix(r.URL.Path, p.basePath) {
                panic("HTTPPool serving unexpected path: " + r.URL.Path)
        }
        p.Log("%s %s", r.Method, r.URL.Path)
        // /<basePath>/groupname/key required
        parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
        if len(parts) != 2 {
                http.Error(w, "bad request", http.StatusBadRequest)
                return
        }
        groupName := parts[0]
        key := parts[1]
        group := GetGroup(groupName)
        if group == nil {
                http.Error(w, "no such group: "+groupName, http.StatusNotFou    nd)
                return
        }
        view, err := group.Get(key)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        w.Header().Set("Content-Type", "application/octet-stream")
        w.Write(view.ByteSlice())
    }

- ServeHTTP的实现逻辑是比较简单的，首先判断访问路径的前缀是否是`basePath`，不是返回错误
- 约定访问路径格式为`/<basePath>/<groupName>/<key>`，通过groupname得到group实例，再使用`group.Get(key)`获取缓存数据
- 最终使用`w.Write()`讲缓存值作为httpResponse的body返回。

&emsp;&emsp;鳄梨，HTTP服务端已经完整地实现了。接下来，我们将在单机上启动HTTP服务，使用curl进行测试。

# 测试

&emsp;&emsp;实现main函数，实例化group，并启动HTTP服务。

        package geec

        import (
                "fmt"
                "log"
                "net/http"
                "testing"
        )

        var dbs = map[string]string{
                "Tom":  "630",
                "Jack": "589",
                "Sam":  "567",
        }

        func TestHttp(t *testing.T) {
                NewGroup("scores", 2<<10, GetterFunc(
                        func(key string) ([]byte, error) {
                                log.Println("[SLowDB] search key", key)
                                if v, ok := dbs[key]; ok {
                                        return []byte(v), nil
                                }
                                return nil, fmt.Errorf("%s not exist", key)
                        }))

                addr := "127.0.0.1:9999"
                peers := NewHTTPPool(addr)
                log.Println("geec is running at", addr)
                log.Fatal(http.ListenAndServe(addr, peers))
        }

- 同样地，我们使用map模拟了数据源db
- 创建一个名为scores的group，若缓存为空，回调函数会从db中获取数据并返回
- 使用http.ListenAndServe在9999端口启动了HTTP服务

        curl http://127.0.0.1:9999/_geec/scores/Tom
        630
        curl http://127.0.0.1:9999/_geec?scores/kkk
        kkk not exist

&emsp;&emsp;节点通信不仅需要HTTP服务端，还需要HTTP客户端，这就是我们下一步需要做的事情。