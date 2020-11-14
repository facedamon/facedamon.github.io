
# 摘要

- **缓存雪崩**：缓存在同一时刻全部失效，造成瞬间DB请求量大、压力骤增，引起雪崩。缓存雪崩通常因为缓存服务器宕机、缓存的key设置了相同的过期时间引起。
- **缓存击穿**：一个存在的key，在缓存过期的一刻，同时有大量的请求，这些请求都会击穿到DB，造成瞬间DB请求量巨大、压力骤增
- **缓存穿透**：查询一个不存在的数据，因为不存在则不会写到缓存中，所以每次都会去请求DB，如果瞬间流量过大，穿透到DB，导致宕机。

# 如何解决

&emsp;&emsp;在上节中，我们并发了N个请求，8003节点向8001同时发起了N次请求。假设对数据库的访问没哦1没有任何限制，很可能向数据库发起N次请求，容易导致缓存击穿。即使对数据库做了防护，HTTP请求是非常耗费资源的操作，针对相同的key，8003节点向8001发起三次请求也是没有必要的。那这种情况下，我们如何做到只向远端节点发起一次请求呢？

&emsp;&emsp;geec实现了一个名为singleflight的package来解决这个问题。

    package singleflight
    
    import "sync"

    type call struct {
        wg  sync.WaitGroup
        // 返回值
        val interface{}
        err error
    }

    type Group struct {
        mu sync.Mutex
        // key，call
        m  map[string]*call
    }

- call代表正在进行中，或已经结束的请求。使用WaitGroup锁避免重入
- Group是singleflight的主数据结构，管理不同key的请求call

&emsp;&emsp;实现Do方法

    func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{},     error) {
        g.mu.Lock()
        if g.m == nil {
            g.m = make(map[string]*call)
        }
        if c, ok := g.m[key]; ok {
            g.mu.Unlock()
            // 如果请求正在进行中，则等待
            c.wg.Wait()
            // 请求结束，返回结果
            return c.val, c.err
        }
        c := new(call)
        // 发起请求钱加锁
        c.wg.Add(1)
        // 添加到g.m, 表明key已经有对应的请求在处理
        g.m[key] = c
        g.mu.Unlock()

        // 调用fn, 发起请求
        c.val, c.err = fn()
        // 请求结束
        c.wg.Done()

        g.mu.Lock()
        // 更新g.m
        delete(g.m, key)
        g.mu.Unlock()

        // 返回结果
        return c.val, c.err
    }

- Do方法，接收2个参数，第一个参数是key，第二个参数是一个函数fn。Do的作用就是针对相同的key，无论Do被调用多少次，函数fn都只会被调用一次，等待fn调用结束，返回返回值或错误。

&emsp;&emsp;`g.mu`是保护Group的成员变量`m`不被并发读写而加上的锁。为了便于理解Do函数，我们将`g.mu`暂时去掉。并且把`g.m`延迟初始化的不分去掉，延迟初始化的目的很简单，提高内存使用效率。

    func (g *Group) Do(key string, fn func()) (interface{}. error)) (interface{}, error) {
        if c, ok := g.m[key]; ok {
            // 如果请求正在进行中，则等待
            c.wg.Wait()
            // 请求结束，返回结果
            return c.val, c.err
        }
        c := new(call)
        // 发起请求前加锁
        c.wg.Add(1)
        // 添加到g.m,表明key已经有对应的请求在处理
        g.m[key] = c

        // 调用烦fn， 发起请求
        c.val, c.err = fn()
        c.wg.Done()

        // 更新g.m
        delete(g.m, key)

        return c.val, c.err
    }

&emsp;&emsp;并发协程之间不需要消息传递，非常适合`sync.WaitGroup`

- wg.Add(1)锁加1
- wg.Wait()阻塞，知道锁被释放
- wg.Done()锁减1

# 使用

    type Group struct {
        name      string
        getter    Getter
        mainCache cache
        peers     PeerPicker
        // use singleflight.Group to make sure that
        // each key is only fetched once
        loader *singleflight.Group
    }

    func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
        // ...
        g := &Group{
            // ...
            loader:    &singleflight.Group{},
        }
        return g
    }

    func (g *Group) load(key string) (value ByteView, err error) {
        // each key is only fetched once (either locally or remotely)
        // regardless of the number of concurrent callers.
        viewi, err := g.loader.Do(key, func() (interface{}, error) {
            if g.peers != nil {
                if peer, ok := g.peers.PickPeer(key); ok {
                    if value, err = g.getFromPeer(peer, key); err == nil {
                        return value, nil
                    }
                    log.Println("[GeeCache] Failed to get from peer", err)
                }
            }

            return g.getLocally(key)
        })

        if err == nil {
            return viewi.(ByteView), nil
        }
        return
    }

- 修改geec.go中的Group，添加成员变量loader，并更新构建函数NewGroup
- 修改load函数，将原来的load的逻辑，使用`g.loader.Do`包裹起来即可，这样确保了并发场景下针对相同的key，load过程只会调用一次。

# 测试

    [Server http://127.0.0.1:8003] Pick peer http://127.0.0.1:8001
    [Server http://127.0.0.1:8001] GET /_geec/scores/Tom
    [SlowDB] search key tom

&emsp;&emsp;可以看到，向API发起了三次并发请求，但8003只向8001发起了一次请求，就搞定了。

&emsp;&emsp;如果并发度不高，可能仍会看到向8001请求三次的场景。这种情况下三次请求是串行的，并没有触发singleflight的锁机制工作，可以加大并发数量再测试。