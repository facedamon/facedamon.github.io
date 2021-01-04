
# 流程回顾

    接收key --> 检查是否被缓存 ---是---> 返回缓存值(1)
                    |否
                    |------>是否应当从远程节点获取 ------> 与远程节点交互 ---> 返回缓存值(2)
                                ｜否
                                ｜------> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值(3)

&emsp;&emsp;我们在第二节中描述了geec的流程。在这之前已经实现了流程1和3，今天实现流程2，从远程节点获取缓存值。

    使用一致性哈希选择节点
        |------> 是否从远程节点 ---是---> HTTP客户端访问远程节点 --> 成功？---是---> 服务端返回返回值                                                    ｜否
                    ｜否
                    ｜-------------------------------------> 回退倒本地节点处理

# 抽象PeerPicker

    package geec

    // PeerPicker is the interface that must be impl to locate
    // the peer that owns a specific key.
    type PeerPicker interface {
                PickPeer(key string) (peer PeerGetter, ok bool)
    }
    // PeerGetter is the interface that must be impl by a peer
    type PeerGetter interface {
            Get(group string, key string) ([]byte, error)
    }

- 在这里，抽象出2个接口，PeerPicker的PickPeer方法用于根据传入的key选择相应节点PeerGetter
- 接口PeerGetter的Get方法用于从对应group查找缓存值。PeerGetter就对应于上述流程中的HTTP客户端

# 节点选择与HTTP客户端

&emsp;&emsp;在geec中我们为HTTPPool实现了服务端功能，通信不仅需要服务端还需要客户端，因此，我们接下来要为HTTPPool实现客户端的功能。

&emsp;&emsp;首先创建具体的HTTP客户端httpGetter，实现PeerGetter接口。

    type httpGetter struct {
        baseURL string
    }

    func (h *httpGetter) Get(group string, key string) ([]byte, error) {
        u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))
        res, err := http.Get(u)
        if err != nil {
            return nil, err
        }
        defer res.Body.Close()

        if res.StatusCode != http.StatusOK {
            return nil, fmt.Errorf("server returned: %v", res.Status)
        }
        bytes, err := ioutil.ReadAll(res.Body)
        if err != nil {
            return nil, fmt.Errorf("reading response body: %v", err)
        }

        return bytes, nil
    }

- baseURL表示将要访问的远程节点的地址，例如`http://example.com/geec`
- 使用`http.Get()`方法获取返回值，并转换为`[]byte`类型

&emsp;&emsp;第二步，为HTTPPool添加节点选择的功能

    const defaultReplicas = 50

    type HTTPPool struct {
        // this peer`s base url, eg: "https://example.net:8080"
        self     string
        basePath string

        // guards peers and httpGetters
        mu    sync.Mutex
        peers *consistenthash.Map
        // keyed by eg. "http://127.0.0.1:9999"
        httpGetters map[string]*httpGetter
    }

- 新增成员变量peers，类型是一致性哈希算法的Map，用来根据具体的key选择节点
- 新增成员变量httpGetters，映射远程节点与对应的httpGetter。每一个远程节点对应一个httpGetter，因为httpGetter与远程节点的地址baseURL有关。

&emsp;&emsp;第三步，实现PeerPicker接口

    func (p *HTTPPool) Set(peers ...string) {
            p.mu.Lock()
            defer p.mu.Unlock()
            p.peers = consistenthash.New(defaultReplicas, nil)
            p.peers.Add(peers...)
            p.httpGetters = make(map[string]*httpGetter, len(peers))
            for _, peer := range peers {
                    p.httpGetters[peer] = &httpGetter{
                            baseURL: peer + p.basePath,
                    }
            }
    }

    // PickPeer picks a peer according to key
    func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
            p.mu.Lock()
            defer p.mu.Unlock()
            if peer := p.peers.Get(key); peer != "" && peer != p.self {
                    p.Log("Pick peer %s", peer)
                    return p.httpGetters[peer], true
            }
            return nil, false
    }

- Set()方法实例化了一致性哈希算法，并且添加了传入的节点
- 并为每一个节点创建了一个HTTP客户端httpGetter
- PickerPeer()包装了一致性哈希算法的Get方法，根据具体的key，选择节点，返回节点对应的HTTP客户端

&emsp;&emsp;至此，HTTPPool既具备了提供HTTP服务的能力，也具备了根据具体的key，创建HTTP客户端从远程节点获取缓存值的能力。

# 实现主流程

&emsp;&emsp;最后，需要将上述新增的功能集成在主流程中。

    func createGroup() *Group {
        return NewGroup("scores", 2<<10, GetterFunc(
            func(key string) ([]byte, error) {
                log.Println("[SlowDB] search key", key)
                if v, ok := db[key]; ok {
                    return []byte(v), nil
                }
                return nil, fmt.Errorf("%s not exist", key)
            }))
    }

    func startCacheServer(addr string, addrs []string, gee *Group) {
        peers := NewHTTPPool(addr)
        peers.Set(addrs...)
        gee.RegisterPeers(peers)
        log.Println("geec is running at", addr)
        log.Fatal(http.ListenAndServe(addr[7:], peers))
    }

    func startAPIServer(apiAddr string, gee *Group) {
        http.Handle("/api", http.HandlerFunc(
            func(w http.ResponseWriter, r *http.Request) {
                key := r.URL.Query().Get("key")
                view, err := gee.Get(key)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
                w.Header().Set("Content-Type", "application/octet-stream")
                w.Write(view.ByteSlice())
            }))
        log.Println("fontend server is running at", apiAddr)
        log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
    }

    func TestHTTPPool_PickPeer(t *testing.T) {
        var port int
        var api bool
        flag.IntVar(&port, "port", 8001, "geec server port")
        flag.BoolVar(&api, "api", true, "Start a api server")
        flag.Parse()

        apiAddr := "http://127.0.0.1:9999"
        addrMap := map[int]string {
            8001: "http://127.0.0.1:8001",
            8002: "http://127.0.0.1:8002",
            8003: "http://127.0.0.1:8003",
        }

        var addrs []string
        for _, v := range addrMap {
            addrs = append(addrs, v)
        }

        gee := createGroup()
        if api {
            go startAPIServer(apiAddr, gee)
        }
        startCacheServer(addrMap[port], []string(addrs), gee)
    }

&emsp;&emsp;main函数的代码比较多，但是逻辑是非常简单的。

- startCacheServer用来启动缓存服务器：创建HTTPPool，添加节点信息，注册到gee中，启动HTTP服务(共3个端口， 8001/8002/8003)，用户不感知。
- startAPIServer用来启动一个API服务(端口9999)，与用户进行交互，用户感知。

        go build -o server 
        ./server -port=8001 &
        ./server -port=8002 &
        ./server -port=8003 -api=1 &

        echo ">>> start test"
        curl "http://127.0.0.1:9999/api?key=Tom" &
        curl "http://127.0.0.1:9999/api?key=Tom" &
        curl "http://127.0.0.1:9999/api?key=Tom" &

        [Server http://127.0.0.1:8003] Picker peer http://127.0.0.1://8001
        [Server http://127.0.0.1:8003] Picker peer http://127.0.0.1://8001
        [Server http://127.0.0.1:8003] Picker peer http://127.0.0.1://8001

&emsp;&emsp;测试的时候，我们并发了3个请求，从日志中看出，三次均选择了节点8001，这是一致性哈希算法的影响。但是有一个问题在于，同时向8001发起了3次请求。试想，加入有10万个请求在并发请求该数据呢？那就会向8001同时发起10万次请求，如果8001又同时向数据库发起10万次查询请求，很容易导致缓存被击穿。

&emsp;&emsp;三次请求的结果是一致的，对于相同的key，能不能只向8001发起一次请求？这个问题下一次解决。