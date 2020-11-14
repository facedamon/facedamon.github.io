
# 摘要

&emsp;&emsp;在上一篇博客中，我们实现了LRU缓存淘汰策略。接下来我们使用`sync.RWMutex`封装LRU的几个方法，使之支持并发的读写。(**对于并发编程的详细介绍详见[golang并发编程]**)。在这之前，我们抽象了一个只读数据结构`ByteView` **用来表示缓存值**，是geec主要的数据结构之一。

		// byteview.go
		package geec
		
		// ByteView holds an immutable view of bytes
		type ByteView struct {
		        b []byte
		}
		
		// Len returns that view`s length
		func (v ByteView) Len() int {
		        return len(v.b)
		}
		
		// ByteSlice returns a copy of data as a byte slice
		func (v ByteView) ByteSlice() []byte {
		        return cloneBytes(v.b)
		}
		
		// Implement String interface
		func (v ByteView) String() string {
		        return string(v.b)
		}
		
		func cloneBytes(b []byte) []byte {
		        c := make([]byte, len(b))
		        copy(c, b)
		        return c
		}

- ByteView只有一个数据成员，`b []byte`，b将会存储真实的缓存值。选择byte类型是为了能够支持任意的数据类型的存储，例如字符串、图片等。
- 实现`Len() int`方法，我们在lru.Cache的实现中，要求被缓存对象必须实现Value接口，返回其所占的内存的大小。
- `b`是只读的，使用`ByteSlice()`方法返回一个拷贝，防止缓存值被外部程序修改。

&emsp;&emsp;接下来就可以为lru.Cache添加并发特性了。

		// cache.go
		package geec
		
		import (
		        "geec/lru"
		        "sync"
		)
		
		type cache struct {
		        mu         sync.Mutex
		        lru        *lru.Cache
		        cacheBytes int64
		}
		
		func (c *cache) add(key string, value ByteView) {
		        c.mu.Lock()
		        defer c.mu.Unlock()
		        if c.lru == nil {
		                c.lru = lru.New(c.cacheBytes, nil)
		        }
		        c.lru.Add(key, value)
		}
		
		func (c *cache) get(key string) (value ByteView, ok bool) {
		        c.mu.Lock()
		        defer c.mu.Unlock()
		        if c.lru == nil {
		                return
		        }
		        if v, ok := c.lru.Get(key); ok {
		                return v.(ByteView), ok
		        }
		        return
		}
- `cache.go`的实现非常简单，实例化lru，邓庄get和add方法，并添加互斥锁mu。
- 在`add`方法中，判断了`c.lru`是否为nil，如果不等于nil再创建实例。这种方法称之为**延迟初始化**`Lazy Initialization`，一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。主要用于提高性能，并减少程序内存要求。

# 主体结构 Group

&emsp;&emsp;Group是geec最核心的数据结构，负责与用户的交互，并且控制缓存值存储和获取的流程。

		接收key --> 检查是否被缓存 ---是---> 返回缓存值(1)
						|否
						|------>是否应当从远程节点获取 ------> 与远程节点交互 ---> 返回缓存值(2)
									｜否
									｜------> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值(3)

&emsp;&emsp;接下来我们将实现流程(1)和(3)，远程交互的部分后续再实现。

## 回调Getter

&emsp;&emsp;我们思考一下，如果缓存不存在，应该从数据源(文件、数据库等)获取数据并添加到缓存中。geec是否应该支持多种数据源的配置呢？不应该，一是数据源的种类太多，没办法实现；二是扩展性不好。**如果从源头获取数据，应该是用户决定的事情，**我们就把这件事交给用户好了。因此，我们设计了一个回调函数(callback)，在缓存不在时，调用这个函数，得到源数据。
		
		```go
		package geec
		
		// Getter loads data for a key
		type Getter interface {
		        Get(key string) ([]byte, error)
		}
		
		// GetterFunc impl Getter with a function
		type GetterFunc func(key string) ([]byte, error)
		
		func (f GetterFunc) Get(key string) ([]byte, error) {
		        return f(key)
		}
		

- 定义接口Getter和回调函数`Get(key string)([]byte, error)`,参数是key，返回值[]byte。
- GetterFunc函数型接口。

## Group定义

		// Group is a cache namespace and associated data loaded spread over
		type Group struct {
		        name      string
		        getter    Getter
		        mainCache cache
		}
		
		var (
		        mu     sync.RWMutex
		        groups = make(map[string]*Group)
		)
		
		// NewGroup create a new instance of Group
		func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
		        if getter == nil {
		                panic("nil Getter")
		        }
		        mu.Lock()
		        defer mu.Unlock()
		        g := &Group{
		                name:      name,
		                getter:    getter,
		                mainCache: cache{cacheBytes: cacheBytes},
		        }
		        groups[name] = g
		        return g
		}
		
		// GetGroup returns the named group previously created with NewGroup,
		// oe nil if there`s no such group.
		func GetGroup(name string) *Group {
		        mu.RLock()
		        defer mu.RUnlock()
		        g := groups[name]
		        return g
		}

- 一个Group可以认为是一个缓存的命名空间，每个Group用拥有一个唯一的名称name。比如可以创建三个Group，缓存学生的成绩命名为scores，缓存学生信息的命名为info，缓存学生课程的命名为courses。
- 第二个属性是`getter Getter`，即缓存未命中时获取源数据的回调。
- 第三个属性是`mainCache cache`，即一开始实现的并发缓存。
- 构建函数`NewGroup`实例化Group，并将group存储在全局变量groups中。
- `GetGroup`用来特定名称的Group，这里使用了只读锁`RLock()`，因为不涉及任何冲突变量的写操作。

## Group的Get方法

&emsp;&emsp;这是geec最为核心的方法。

		// Get value for a key from cache
		func (g *Group) Get(key string) (ByteView, error) {
		        if key == "" {
		                return ByteView{}, fmt.Errorf("key is required")
		        }
		
		        if v, ok := g.mainCache.get(key); ok {
		                log.Println("[geec] hit")
		                return v, nil
		        }
		        return g.load(key)
		}
		
		func (g *Group) load(key string) (value ByteView, err error) {
		        return g.getLocally(key)
		}
		
		func (g *Group) getLocally(key string) (ByteView, error) {
		        bytes, err := g.getter.Get(key)
		        if err != nil {
		                return ByteView{}, err
		        }
		        value := ByteView{b: cloneBytes(bytes)}
		        g.populateCache(key, value)
		        return value, nil
		}
		
		func (g *Group) populateCache(key string, value ByteView) {
		        g.mainCache.add(key, value)
		}

- Get方法实现了上述所说的流程(1)和(3)。
- 流程(1): 从mainCache中查找缓存，如果存在则返回缓存值。
- 流程(3): 缓存不存在，则调用load方法，load调用getLocally(分布式场景下会调用getFromPeer从其它节点获取)，getLocally调用用户回调函数`g.getter.Get()`获取源数据，并且将源数据添加到缓存mainCache中。

# 测试

		// geec_test.go
		
		package geec
		
		import (
		        "fmt"
		        "log"
		        "testing"
		)
		
		var db = map[string]string{
		        "Tom":  "630",
		        "Jack": "589",
		        "Sam":  "567",
		}
		
		func TestGet(t *testing.T) {
		        loadCounts := make(map[string]int, len(db))
		        g := NewGroup("scores", 2<<10, GetterFunc(
		                func(key string) ([]byte, error) {
		                        log.Println("[SlowDB] search key:", key)
		                        if v, ok := db[key]; ok {
		                                if _, ok := loadCounts[key]; !ok {
		                                        loadCounts[key] = 0
		                                }
		                                loadCounts[key] += 1
		                                return []byte(v), nil
		                        }
		                        return nil, fmt.Errorf("%s not exist", key)
		                },
		        ))
		
		        for k, v := range db {
		                if view, err := g.Get(k); err != nil || view.String() != v {
		                        t.Fatal("failed to get value of Tom")
		                }
		                // load from callback func
		                if _, err := g.Get(k); err != nil || loadCounts[k] > 1 {
		                        t.Fatalf("cache %s miss", k)
		                }
		                // cache hit
		        }
		
		        if view, err := g.Get("unknown"); err == nil {
		                t.Fatalf("the value of unknow should be empty, but %s got", view)
		        }
		}

- 在这个测试用例中，我们主要测试了2中情况
1. 在缓存为空的情况下，能够通过回调函数获取到源数据。
2. 在缓存已经存在的情况下，是否直接从缓存中获取，为了实现这一点，使用loadCounts统计某个键调用回调函数的次数，如果次数大于1，则表示调用了对此回调函数，没有缓存。

> 测试结果

		$ go test -run TestGet
		2020/07/24 14:59:52 [SlowDB] search key: Tom
		2020/07/24 14:59:52 [geec] hit
		2020/07/24 14:59:52 [SlowDB] search key: Jack
		2020/07/24 14:59:52 [geec] hit
		2020/07/24 14:59:52 [SlowDB] search key: Sam
		2020/07/24 14:59:52 [geec] hit
		2020/07/24 14:59:52 [SlowDB] search key: unknown
		PASS
		ok      geec    0.148s
