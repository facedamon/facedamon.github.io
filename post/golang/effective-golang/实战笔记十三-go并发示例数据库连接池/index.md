
>转载自 https://www.flysnow.org

## 资源池
这篇文章演示使用有缓冲的通道实现一个资源池，这个资源池可以管理在任意多喝goroutine之间共享的资源，比如网络链接、数据库链接等，我们在数据库操作的时候，比较常见的就是数据连接池，也可以基于我们实现的资源池来实现。

好了，老规矩，我们先构建一个资源池结构体，然后再赋予一些方法，这个资源池就可以帮助我们管理资源了。
		type pool struct {
		   m sync.Mutex
		   res chan io.Closer
		   factory func()(io.Closer, error)
		   closed bool
		}
这个结构体pool有四个字段，其中m是一个互斥锁，这主要来保证在多个goroutine访问资源时，池内的值是安全的。

res字段是一个有缓冲的通道，用来保存共享的资源，这通道的大小，在初始化pool的时候就指定的。注意这个通道的类型是io.closer接口，所以实现了这个接口的类型都可以作为资源。

factory这个是一个函数类型，它的作用就是当需要一个新的资源时，可以通过这个函数创建，也就是说它是生成新资源的，至于如何生成、生成何种资源，是使用者决定的，所以这也是这个资源池灵活设计的地方。

closed字段表示资源池是否被关闭，如果被关闭的话，再访问是会有错误的。

现在先这个资源池我们已经定义好了，也知道了每个字段的含义，刚刚我们说到关闭错误，那么我们就先定义一个资源池已经关闭的错误。

		var ErrPoolClosed = errors.New("资源池已经关闭")
非常简介，当我们从资源池获取资源的时候，如果该资源池已经关闭，那么就会返回这个错误。单独定义它的目的，是和其它错误有一个区分，这样需要的时候，我们就可以从重夺的error类型里区分出来这个ErrorPoolClosed。
下面我们为pool专门定制一个new函数。
		func P_new(fn func()(io.Closer, error), n uint) (*pool, error) {
		   if n < 0 {
		      return nil, errors.New("pool size less than zero")
		   }
		   return &pool{factory:fn, res: make(chan io.Closer, n)}, nil
		}
这个函数创建一个资源池，它接收两个参数，一个fn是创建新资源的函数；还有一个n是指定资源池大小。这个函数里，做了n大小的判断，用n创建了一个有缓冲的通道，来保存资源，并且返回一个资源池的指针。

有了资源池，那么我们就可以从中获取资源了。
		func (p *pool) Get() (io.Closer, error) {
		   select {
		      case r, ok := <-p.res:
		         if !ok {
		            return nil, ErrPoolClosed
		         }
		         return r, nil
		   default:
		      return p.factory()
		   }
		}
Get方法可以从资源池获取资源，如果没有资源，则调用factory方法生成一个并返回。这里同样使用了select的多路复用，因为这个函数不能阻塞，可以获取就获取，不能就生成一个。

关闭资源池，意味着整个资源池不能再被使用，然后关闭存放资源的通道，同样是否通道里的资源。
		func (p *pool) Close() {
		   // sync the closed filed
		   p.m.Lock()
		   defer p.m.Unlock()
		
		   if p.closed {return}
		
		   // not write to chan
		   p.closed = true
		   close(p.res)
		
		   // release the io
		   for r := range p.res {
		      r.Close()
		   }
		}
这个方法里，我们使用了互斥锁，因为有个标记资源资源池是否关闭的字段closed需要多个goroutine操作，所以我们必须保证这个字段的同步。

然后我们关闭通道，禁止写入，而且我们前面的Get也可以感知到通道已经关闭了。

关闭方法有了，我们看看资源回收如何实现。
		func (p *pool) Put(r io.Closer) {
		   p.m.Lock()
		   defer p.m.Unlock()
		
		   if p.closed {
		      r.Close()
		      return
		   }
		
		   select {
		   case p.res <- r:
		      log.Println("将资源回收至资源池")
		   default:
		      log.Println("资源池满了，释放该资源")
		      r.Close()
		   }
		}
释放资源本质就是把资源再发送到缓冲通道中，就是这么简单，不过为了更安全的实现这个方法，我们使用了互斥锁，保证closed标志的安全，而且这个互斥锁还有一个好处，就是不会往一个已经关闭的通道发送资源。

这是为什么呢？因为close和release这两个方法是互斥的，close方法里对close标志的修改，put方法是可以感知到的，所以就直接return了，不会执行下面的select代码了，也就不会往一个已经关闭的通道里发送资源了。

如果资源池没有关闭，则继续尝试往资源池通道发送资源，如果可以发送，就等于资源又回到资源池了；如果发送不了，说明资源池满了，该资源就无法重新回到资源池里，那么我们就把这个需要释放的资源关闭，抛弃。

好了，资源池管理写好了，也知道资源池是如何实现的了，现在我们看看如何使用这个资源池，模拟一个数据库连接池吧。
		package pool
		
		import (
		   "io"
		   "log"
		   "math/rand"
		   "sync"
		   "sync/atomic"
		   "testing"
		   "time"
		)
		
		const max = 5
		const pools = 2
		
		type db struct {
		   id int32
		}
		
		func (d db) Close() error {
		   log.Println("db close", d.id)
		   return nil
		}
		
		var idcounter int32
		
		func createDb() (io.Closer, error) {
		   id := atomic.AddInt32(&idcounter, 1)
		   return &db{id}, nil
		}
		
		func query(q int, p *pool) {
		   conn, err := p.Get()
		   if err != nil {
		      log.Println(err)
		      return
		   }
		   defer p.Put(conn)
		   //模拟查询
		   time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		   log.Printf("第%d个查询，使用的是ID为%d的数据库连接", q, conn.(*db).id)
		}
		
		func TestPool_(t *testing.T) {
		   var wg sync.WaitGroup
		   wg.Add(max)
		
		   p, err := p_new(createDb, pools)
		   if err != nil {
		      log.Println(err)
		      return
		   }
		
		   for q := 0; q < max; q++ {
		      go func(t int) {
		         query(t, p)
		         wg.Done()
		      }(q)
		   }
		
		   wg.Wait()
		   log.Println("开始关闭资源池")
		   p.Close()
		}
		GOROOT=/usr/local/go #gosetup
		GOPATH=/Users/facedamon/go #gosetup
		/usr/local/go/bin/go test -c -o /private/var/folders/m4/74p2bmyn7vb6962wkrvr7tnc0000gn/T/___TestPool__in_facedamon_faceawesome_pool facedamon/faceawesome/pool #gosetup
		/usr/local/go/bin/go tool test2json -t /private/var/folders/m4/74p2bmyn7vb6962wkrvr7tnc0000gn/T/___TestPool__in_facedamon_faceawesome_pool -test.v -test.run ^TestPool_$ #gosetup
		=== RUN TestPool_
		2019/05/30 18:00:55 生成新资源
		2019/05/30 18:00:55 生成新资源
		2019/05/30 18:00:55 生成新资源
		2019/05/30 18:00:55 生成新资源
		2019/05/30 18:00:55 生成新资源
		2019/05/30 18:00:55 第2个查询，使用的是ID为4的数据库连接
		2019/05/30 18:00:55 将资源回收至资源池
		2019/05/30 18:00:55 第4个查询，使用的是ID为1的数据库连接
		2019/05/30 18:00:55 将资源回收至资源池
		2019/05/30 18:00:55 第3个查询，使用的是ID为5的数据库连接
		2019/05/30 18:00:55 资源池满了，释放该资源
		2019/05/30 18:00:55 db close 5
		2019/05/30 18:00:56 第1个查询，使用的是ID为3的数据库连接
		2019/05/30 18:00:56 资源池满了，释放该资源
		2019/05/30 18:00:56 db close 3
		2019/05/30 18:00:56 第0个查询，使用的是ID为2的数据库连接
		2019/05/30 18:00:56 资源池满了，释放该资源
		2019/05/30 18:00:56 db close 2
		2019/05/30 18:00:56 开始关闭资源池
		2019/05/30 18:00:56 db close 4
		2019/05/30 18:00:56 db close 1
		--- PASS: TestPool_ (0.89s)
		PASS
		
		Process finished with exit code 0
到这里，我们已经完成了一个资源池的管理，并且进行了使用测试。资源对象池的使用比较频繁，因为我们想把一些对象缓存起来，以便使用，这样就会比较高效，而且不会经常调用GC，为此Go为我们提供了原声的资源池管理，防止我们重复造轮子，这就是sync.Pool，我们看下刚刚我们的例子，如果使用sync.pool实现。
		func createDb2() interface{} {
		   id := atomic.AddInt32(&idcounter, 1)
		   return &db{id}
		}
		func query2(q int, p *sync.Pool)  {
		   conn := p.Get().(*db)
		   defer p.Put(conn)
		
		   //模拟查询
		   time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		   log.Printf("第%d个查询，使用的是ID为%d的数据库连接", q, conn.id)
		}
		func TestSyncPool(t *testing.T)  {
		   var wg sync.WaitGroup
		   wg.Add(max)
		
		   p := &sync.Pool{New: createDb2}
		
		   for q := 0; q < max; q++ {
		      go func(t int) {
		         query2(t, p)
		         wg.Done()
		      }(q)
		   }
		
		   wg.Wait()
		
		}
