
>转载自 https://www.flysnow.org

## 读写锁
前面的有篇文章在讲资源竞争的时候，讲互斥锁，互斥锁的根本就是当一个goroutine访问的时候，其它goroutine都不能访问，这样保证了资源的同步，避免了竞争，不过也降低了性能。

仔细剖析我们的场景，当我们读取一个数据的时候，如果这个数据永远也不会被修改，那么其实是不存在资源竞争的问题的，因为数据是不可变的，不管怎么读取，多少goroutine同时读取，都是可以的。

所以其实读取并不是问题，问题是修改，修改的数据要同步，这样其它的goroutine才可以感知到。所以真正的互斥应该是读取和修改、修改和修改之间，读取和读取是没有互斥操作的。

所以这就延伸出来另一种锁，叫做读写锁。

读写锁可以让多个读操作同步并发，同时读取，但是对于写操作是完全互斥的。也就是说，当一个goroutine进行写操作的时候，其它goroutine既不能进行读操作，也不能进行写操作。
		var count int
		var wg sync.WaitGroup
		var rw sync.RWMutex
		
		func read(n int){
		    rw.RLock()
		    fmt.Printf("读goroutine %d 正在读取...\n", n)
		    v := count
		    fmt.Printf("读goroutine %d 读取结束，值为：%d\n", n,v)
		    wg.Done()
		    rw.RUnlock()
		}
		
		func write(n int){
		    rw.Lock()
		    fmt.Printf("写goroutine %d 正在写入...\n",n)
		    v := rand.Intn(1000)
		    count = v
		    fmt.Printf("写goroutine %d写入结束，新值为：%d\n", n,v)
		    wg.Done()
		    rw.Unlock()
		}
		
		func main(){
		    wg.Add(10)
		    for i:=0; i < 5; i++ {
		        go read(i)
		    }
		    for i:=0; i < 5; i++ {
		        go write(i)
		    }
		    wg.Wait()
		}
我们可以使用`go build -race`检测，也没有竞争提示。

我们在做java开发的时候，肯定知道SynchronizedMap这个map，它是一个在多线程下安全的map，我们可以通过Collections.synchronizedMap(Map<K,V>)来获取一个安全的Map，下面我们看看如何使用读写锁，基于go来实现一个线程安全map。
		package common
		
		import "sync"
		
		type smap struct{
		    rw *sync.RWMutex
		    data map[interface{}]interface{}
		}
		
		func New() *smap{
		    return &smap{
		        rw: new(sync.RWMutex),
		        data: make(map[interface{}])interface{}
		    }
		}
		
		func (sm *smap)Store(k,v interface{}){
		    sm.rw.Lock()
		    defer sm.rw.Unlock()
		    sm.data[k] = v
		}
		
		func (sm *smap)Load(k interface{}) interface{}{
		    sm.rw.RLock()
		    defer sm.rw.RUnlock()
		    return sm.data[k]
		}
		
		func (sm *smap)Delete(k interface{}){
		    sm.rw.Lock()
		    defer sm.rw.Unlock()
		    delete(sm.data, k)
		}
		
		func (sm *smap)Each(cd func(interface{},interface{})){
		    sm.rw.RLock()
		    defer sm.rw.RUnlock()
		    for k, v := range sm.data{
		        cb(k, v)
		    }
		}
在go1.9之后的版本中，提供了默认的sync.map线程安全的map，采用的是lock-free方案，摒弃了java中的分段锁方案。
