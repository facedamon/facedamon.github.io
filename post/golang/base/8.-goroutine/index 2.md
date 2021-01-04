
# goroutine

		package main
		
		func main() {
		    for i := 0; i < 10; i++ {
		        go func(i int){
		            for {
		                fmt.Printf("Hello from goroutine %d\n", i)
		            }
		        }(i)
		    }
		    time.Sleep(time.Millisecond)
		}
&emsp;&emsp;如果不加go关键字的话，就是外循环进来一次一直执行匿名函数，不会退出。加了go关键字之后就说明是并发的执行该匿名函数，外层循环还在往下跑，然后并发的开启了一个函数，结果就是，开启了10个goroutine，每个goroutine拼命的打印。为了保证main主线程在go协程完成前不退出，开启sleep。

- 轻量级线程
- 非抢占式多任务处理，由协程主动交出控制权(线程是抢占式，有可能是一个任务执行到一半，就被操作系统暂停切换)
- 编译器／解释器／虚拟机层面的多任务
- 多个协程可能在一个或多个线程上运行

&emsp;&emsp;上面的案例中由于fmt.Printf操作是一个IO操作，因为IO操作总会由等待的过程，IO操作本身是会切换协程,我们想办法不让它切换(不交出控制权)

		package main
		
		func main() {
		    var a [10] int
		    for i := 0; i < 10; i++ {
		        go func(i int){
		            for {
		                a[i]++
		                // runtime.Gosched()
		            }
		        }(i)
		    }
		    time.Sleep(time.Millisecond)
		    fmt.Println(a)
		}
&emsp;&emsp;我们采用数组计算，就不会出现协程切换了，那么该程序一次进去，就会一直死循环。除非执行runtime.Gosched(),手动交出控制权。手动交出控制权后，打印结果如下：在一毫秒的时间内，10个协程获得的机会是不一样的
		[1709 809 1563 1564 1627 1628 1549 1584 833 1231]

&emsp;&emsp;可能有人会问，为什么匿名函数要定义一个i，将外部的i传进去来执行，而不是直接引用外部的i呢？那么我们不定义i，试一试：

		package main
		
		func main() {
		    var a [10] int
		    for i := 0; i < 10; i++ {
		        go func(){
		            for {
		                a[i]++
		                runtime.Gosched()
		            }
		        }()
		    }
		    time.Sleep(time.Millisecond)
		    fmt.Println(a)
		}
		panic runtime error: index out of range
		goroutine 9 running
&emsp;&emsp;这里就是函数式编程的概念了，这个go func就是个闭包，引用了外部i，当外部的i结束循环的时候，i=10,那么其中一个协程闭包引用了该i，所以越界了。所以这个地方是不能使用闭包的，必须采用函数调用的形式拷贝一份i来固定i。

&emsp;&emsp;其实上述案例还有一个bug，就是a[i]++和fmt.Println(a),有可能会出现一边写一边读的可能。这个问题需要使用chan来解决。

# 协程Coroutine

`Subroutines are special cases of more general program components, called coroutines. In contrast to the unsymmetric`

`子程序是协程的一个特例`

&emsp;&emsp;所有的函数调用都可以看成是子程序，所有的函数调用都是协程的特例，协程是一个比子程序更宽泛的概念，这句话是谁说的呢？

`Donnald Knuth "The Art of Computer Programming.Vol1"`

&emsp;&emsp;我们来看看它怎样来算一个special cases呢！

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/coroutine.png)

- 普通函数
  - 在一个线程内，这个线程有一个main函数，main调用dowork，等dowork做完了以后，才会把控制权交给main函数，然后main去执行下一个语句。

- 协程
  - 也是main和dowork，但是main和dowork不是单向箭头，而是双向通道。main和dowork之间的数据可以双向流通，控制权也可以双向流通，就相当于并发执行的两个线程，main和dowork有可能运行在一个线程中，也可能是多个线程，这件事不需要开发者管理，调度器去管理。

# 其他语言中的协程

- C++： Boost.Coroutine
- Java: 不支持
- Python：3.5以前使用yield关键字，3.5后加入了async def对协程原声支持

# goroutine

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/goroutine.png)

&emsp;&emsp;go语言程序开始后，有一个调度器，有些协程可能放在一个线程里，有些可能两个放在一个线程里，或者很多个协程放在一个线程里，这件事由调度器来管理。

- 任何函数只需加上go就能送给调度器运行
- 不需要在定义时区分是否异步函数
- 调度器在合适的点进行切换
- 使用 -race来检测数据访问冲突(`go run -race xxx.go`)

# goroutine可能的切换点

- I/O, select
- channel
- 锁
- 函数调用
- runtime.Gosched()

# channel

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/channel.png)

&emsp;&emsp;上面说，我们可以在程序中开很多个goroutine，那么多个goroutine之间的通道就是channel。

		package main
		
		func chanDemo() {
		    //var c chan int  // c == nil
		    c := make(chan int)
		    c <- 1  // 发送给chan
		    c <- 2
		    n := <- c   // 拿出chan
		    fmt.Println(n)
		}
		
		func main() {
		    chanDemo()
		}
&emsp;&emsp;来看看程序能不能跑起来：
		fatal error: all goroutines are asleep - deadlock!
		goroutine 1 [chan send]:
&emsp;&emsp;程序死锁了，在c <-1处。我们仔细看上面的图片，channel是goroutine和goroutine之间的通道，这里发了一个1，但是没人收，肯定会报错。因此需要一个goroutine负责收。

		package main
		
		func chanDemo() {
		    //var c chan int  // c == nil
		    c := make(chan int)
		    go func(){
		        n := <-c
		        fmt.Println(n)
		    }()
		    c <- 1  // 发送给chan
		    c <- 2
		}
		
		func main() {
		    chanDemo()
		}
&emsp;&emsp;
		1
&emsp;&emsp;只打印了1.这是因为发送了2，还没来得及收，main就结束了，所以：
		package main
		
		func chanDemo() {
		    //var c chan int  // c == nil
		    c := make(chan int)
		    go func(){
		        n := <-c
		        fmt.Println(n)
		    }()
		    c <- 1  // 发送给chan
		    c <- 2
		    time.Sleep(time.Millisecond)
		}
		
		func main() {
		    chanDemo()
		}
		// output:
		// 1
		// 2

&emsp;&emsp;golang函数是一等公民，那么channel也可以作为参数来传递

		package main
		
		func worker(id int, c chan int) {
		    for {
		        fmt.Printf("Worker %d received %c\n", id, <-c)
		    }
		}
		
		func chanDemo() {
		    var channels [10]chan int
		    for i := 0; i < 10; i++ {
		        channels[i] = make(chan int)
		        go worker(i, channels[i])
		    }
		
		    for i := 0; i < 10; i++ {
		        channels[i] <- 'a' + i
		    }
		    time.Sleep(time.Millsecnod)
		}
		
		func main(){
		    chanDemo()
		}
		// Output
		//  Worker 0 received a
		//  Worker 4 received e
		//  Worker 1 received b
		//  Worker 3 received d
		//  Worker 7 received h
		//  Worker 5 received f
		//  Worker 2 received c
		//  Worker 9 received j
		//  Worker 8 received i
		//  Worker 6 received g
&emsp;&emsp;channel当然也可以作为返回值，那么这个worker其实也就不叫worker了，而是createWorker。而且还能修饰chan的方向，告诉调用者chan是发还是收
		package main
		
		func createWorker(id int) chan<- int {
		    c := make(chan int)
		    // 这才是真正的worker
		    go func(){
		        for {
		            fmt.Println("Worker %d received %c\n", id, <-c)
		        }
		    }()
		    return c
		}
		
		func chanDemo() {
		    var channels[10]chan<- int
		    for i:= 0; i < 10; i++ {
		        // 将返回的chan存储起来，此时的chan还是空的，真正的work还在等待
		        channels[i] = createWorker(i)
		    }
		     
		    for i := 0; i < 10; i++ {
		        channels[i] <- 'a' + i
		    }
		    time.Sleep(time.Millsecnod)
		}

# bufferedchannel
&emsp;&emsp;上述案例中，每发一个chan，就必须有一个goroutine来接收,否则会异常。但是如果一旦发送了数据之后，就一定要进行协程的切换，比较耗资源，那么就可以加入一个缓冲区

		func worker(id int, c chan int) {
		    for {
		        fmt.Printf("Worker %d received %c\n", id, <-c)
		    }
		}
		
		func bufferedChannel() {
		    c := make(chan int， 3)
		    go worker(0, c)
		    c <- 'a'
		    c <- 'b'
		    c <- 'c'
		    c <- 'd'
		    time.Sleep(time.Millisecond)
		}
		
		func main() {
		    bufferedChannel()
		}
&emsp;&emsp;我们发现这个worker，它不知道这个数据是不是发送完了，有一种方法可以告诉接收方数据发完了，那就是Close.当chan close之后worker还是会收到数据的(Millsecond)，int型就是0， 否则就是空串,那么来改进worker。

		func worker(id int, c chan int) {
		    /*for {
		        fmt.Printf("Worker %d received %c\n", id, <-c)
		    }
		    */
		    /*
		    for {
		        n, ok := <-c
		        if !ok {
		            break;
		        }
		        fmt.Printf("Worker %d received %c\n", id, n)
		    }
		    */
		    // 还有一种更简洁高效的方法
		    for n := range c {
		       fmt.Printf("Worker %d received %c\n", id, n) 
		    }
		
		}
		
		func bufferedChannel() {
		    c := make(chan int， 3)
		    go worker(0, c)
		    c <- 'a'
		    c <- 'b'
		    c <- 'c'
		    c <- 'd'
		    time.Sleep(time.Millisecond)
		}
		
		func channelClose() {
		    c := make(chan int)
		    go worker(0, c)
		    c <- 'a'
		    c <- 'b'
		    c <- 'c'
		    c <- 'd'
		    close(c)
		    time.Sleep(time.Millisecond)
		}
		
		func main() {
		    //bufferedChannel()
		    channelClose()
		}

- channel
- buffered channel
- range
- 理论基础：Communication Sequential Process

# 实践

`Dont communication by sharing memory; share memory by communicating`: 不要通过共享内存来通信；通过通信来共享内存
&emsp;&emsp;首先需要解决sleep问题，当时加上这句话是因为worker打印不知道外面什么时候结束，所以加上了一个实践界限。那么我们通过通知来解决这个问题

		package main
		
		func doWorker(id int, c chan int, done chan bool) {
		    for n := range c {
		       fmt.Println("Worker %d received %c\n", id, n) 
		        done <- true
		    }
		}
		
		type worker struct {
		    in chan int
		    done chan bool
		}
		
		func createWorker(id int) worker {
		    w := worker{
		        in : make(chan int),
		        done: make(chan bool),
		    }
		    // 这才是真正的worker
		    go doWorker(id, w.in, w.done)
		    return w
		}
		
		func chanDemo() {
		    var workers[10]worker int
		    for i:= 0; i < 10; i++ {
		        workers[i] = createWorker(i)
		    }
		     
		    for i, worker := range workers {
		        worker.in <- 'a' + i
		    }
		    for i, worker := range workers {
		        worker.in <- 'A' + i
		    }
		    // wait for all of them
		    for _, worker := range workers {
		        // 每个worker发了两遍，收两遍
		        <-worker.done
		        <-worker.done
		    }
		}
&emsp;&emsp;看看运行结果,报错了。第一组worker完整的打印完了，第二组worker循环等待了。
		Worker 9 received j
		Worker 1 received b
		...
		fatal error: all goroutines are asleep - deadlock!
&emsp;&emsp;这是为什么呢。其实chan的发送都是阻塞式的。第一组的worker在发送完done后，还没有人来接收done，紧接着就开启第二组worker了，那么第二组worker就阻塞了。这个问题有一种解决方案，就是我们把done的发送放到另一个goroutine里，这样就是并行的发送done，就不会阻塞了。
		func doWorker(id int, c chan int, done chan bool) {
		    for n := range c {
		        fmt.Printf("WOrker %d received %c\n", id, n)
		        go func(){ done <- true}()
		    }
		}
&emsp;&emsp;那么按照两组worker的通信顺序来讲，我们是要等待两次，如果只等一次的话，就不会出现deadlock了，那么第二种解决方法就是，分次等待。
		func chanDemo() {
		    var workers[10]worker int
		    for i:= 0; i < 10; i++ {
		        workers[i] = createWorker(i)
		    }
		     
		    for i, worker := range workers {
		        worker.in <- 'a' + i
		    }
		    for _, worker := range workers {
		        <-worker.done
		    }
		    for i, worker := range workers {
		        worker.in <- 'A' + i
		    }
		
		    for _, worker := range workers {
		        <-worker.done
		    }
		}
&emsp;&emsp;其实对于多任务的等待，golang自身提供了解决方法，WaitGroup。

		package main
		
		func doWorker(id int, c chan int, *sync.WaitGroup) {
		    for n := range c {
		       fmt.Println("Worker %d received %c\n", id, n) 
		        wg.Done()
		    }
		}
		
		type worker struct {
		    in chan int
		    wg *sync.WaitGroup
		}
		
		func createWorker(id int, wg *sync.WaitGroup) worker {
		    w := worker{
		        in : make(chan int),
		        done: wg,
		    }
		    // 这才是真正的worker
		    go doWorker(id, w.in, wg)
		    return w
		}
		
		func chanDemo() {
		    var wg sync.WaitGroup
		    var workers[10]worker int
		    for i:= 0; i < 10; i++ {
		        workers[i] = createWorker(i, &wg)
		    }
		     
		    wg.Add(20)
		    for i, worker := range workers {
		        worker.in <- 'a' + i
		    }
		    for i, worker := range workers {
		        worker.in <- 'A' + i
		    }
		    // wait for all of them
		    wg.Wait()
		}
&emsp;&emsp;我们发现我哥doWork很奇怪，不是很纯粹，不应该接收wg。我们需要包装一下,采用函数式编程。

		package main
		
		func doWorker(id int, w worker) {
		    for n := range w.in {
		       fmt.Println("Worker %d received %c\n", id, n) 
		        w.done()
		    }
		}
		
		type worker struct {
		    in chan int
		    done func()
		}
		
		func createWorker(id int, wg *sync.WaitGroup) worker {
		    w := worker{
		        in : make(chan int),
		        done: func(){
		            wg.done()
		        },
		    }
		    // 这才是真正的worker
		    go doWorker(id, w.in, w)
		    return w
		}
		
		func chanDemo() {
		    var wg sync.WaitGroup
		    var workers[10]worker int
		    for i:= 0; i < 10; i++ {
		        workers[i] = createWorker(i, &wg)
		    }
		     
		    wg.Add(20)
		    for i, worker := range workers {
		        worker.in <- 'a' + i
		    }
		    for i, worker := range workers {
		        worker.in <- 'A' + i
		    }
		    // wait for all of them
		    wg.Wait()
		}
## 使用channel来实现树的遍历

		type treeNode struct {
		    value int
		    left, right *treeNode
		}
		func createNode(value int) *treeNode {
		    //返回了一个局部变量的地址给外面用
		    //这如果在C++中，是个典型的错误
		    return &treeNode{value}
		}
		// 中序遍历
		func (node *treeNode) traverse() {
		    if node == nil {
		        return
		    }
		
		    node.left.traverse()
		    node.print()
		    node.right.traverse()
		}
		
		func (node *treeNode) traverseFunc(f func(*treeNode)) {
		    if node == nil {
		        return
		    }
		    node.left.traverseFunc(f)
		    f(node)
		    node.right.traverseFunc(f)
		}
		
		func (node *treeNode) traverseWithChannel() chan *treeNode {
		    out := make(chan, int)
		    go func(){
		        node.traverseFunc(func(node *treeNode){
		            out <- node
		        })
		        close(out)
		    }()
		    return out
		}
&emsp;&emsp;上述traverseWithChannel的好处是，外部调用可以很直观的序列，比如这样：
		maxNode := 0
		for c := range traverseWithChannel() {
		    if c.value >= maxNode {
		        maxNode = v.value
		    }
		    fmt.Println("maxNode is ", maxNode)
		}
# Select调度
&emsp;&emsp;比如有两个chan，想从c1,c2收数据，如果同时收，谁来的快就收谁。select是非阻塞式的。
		func generator() chan int {
		    out := make(chan int)
		    go func(){
		        i := 0
		        for {
		            time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
		            out <- i
		            i++
		        }
		    }()
		    return out
		}
		func main() {
		    var c1, c2 = generator(), generator()
		    for {
		        select {
		            case n := <- c1:
		            fmt.Println("Received from c1:", n)
		            case n := <- c2:
		            fmt.Println("Received from c2:", n)
		        }
		    }
		}
&emsp;&emsp;上面的c1,c2打印是无序的，谁先发送了，就打印谁。我们再来加入消费者worker，让它更复杂。

		func worker(id, int c chan int) {
		    for n := range c {
		        fmt.Printf("WOrker %d received %d\n", id, n)
		    }
		}
		
		func createWorker(id int) chan<- int {
		    c := make(chan int)
		    go worker(id, c)
		    return c
		}
		
		func main() {
		    var c1, c2 = generator(), generator()
		    w := createWorker(0)
		    for {
		
		        select {
		            case n := <- c1:
		                //w <- n 这不行，这是一个阻塞操作
		            case n := <- c2:
		                //w <- n 这不行，这是一个阻塞操作
		        }
		    }
		}
&emsp;&emsp;上面的程序把生成的chan发送给worker是一个阻塞过程，这不行，必须非阻塞的发送

		func main() {
		    var c1, c2 = generator(), generator()
		    w := createWorker(0)
		    for {
		        n := 0
		        select {
		            case n := <- c1:
		            case n := <- c2:
		            case w <- n:
		        }
		    }
		}
&emsp;&emsp;上面的话就不会阻塞了，但是有个问题，有可能c1,c2接收太慢，每次都是直接w <- n,那我们使用标志符号来区别一下。
		func main() {
		    var c1, c2 = generator(), generator()
		    w := createWorker(0)
		
		    n := 0
		    hasValue := false
		    for {
		        select {
		            case n := <- c1:
		                hasValue = true
		            case n := <- c2:
		                hasValue = true
		            case w <- n:
		                if hasValue {
		                    // send to w
		                }
		        }
		    }
		}
&emsp;&emsp;上述程序意思是c1,c2拿到值了才发送给worker，但是这样的逻辑是不对的。上面已经case了。那么怎么解决呢，我们可以利用nil 在select中的性质。nil在select中不会被阻塞，但是永远不会被匹配到，除非某一次遍历获取了值。

		func main() {
		    var c1, c2 = generator(), generator()
		    worker := createWorker(0)
		
		    n := 0
		    hasValue := false
		    for {
		        var activeWorker chan<- int
		        if hasValue {
		            activeWorker = worker
		        }
		
		        select {
		            case n := <- c1:
		                hasValue = true
		            case n := <- c2:
		                hasValue = true
		            case activeWorker <- n:
		                hasValue = false
		        }
		    }
		}
&emsp;&emsp;上看程序看起来完美，其实还是有问题。生产chan和消耗chan的速度是不一样的，如果c1生成chan的速度太快，那么n就会被覆盖，而且c2无法获取chan了。怎么办呢？我们把收到的n给存下来，排队。
		func main() {
		    var c1, c2 = generator(), generator()
		    worker := createWorker(0)
		
		    var values []int
		    for {
		        var activeWorker chan<- int
		        var activeValue int
		        if len(values) {
		            activeWorker = worker
		            activeValue = values[0]
		        }
		
		        select {
		            case n := <- c1:
		                values = append(values, n)
		            case n := <- c2:
		                values = append(values, n)
		            case activeWorker <- activeValue:
		                values = values[1:]
		        }
		    }
		}
&emsp;&emsp;上面的程序把chan保存下来了，但是可能因为速度差，导致values形成堆积。我们想看values到底挤压了多少chan，可以这样做。
		func main() {
		    var c1, c2 = generator(), generator()
		    worker := createWorker(0)
		
		    var values []int
		    tm := time.After(10 * time.Second)
		    tick := time.Tick(time.Second)  //定时
		    for {
		        var activeWorker chan<- int
		        var activeValue int
		        if len(values) {
		            activeWorker = worker
		            activeValue = values[0]
		        }
		
		        select {
		            case n := <- c1:
		                values = append(values, n)
		            case n := <- c2:
		                values = append(values, n)
		            case activeWorker <- activeValue:
		                values = values[1:]
		            case <- time.After(800 * time.Millisecond):
		                fmt.Println("TimeOut")
		            case <- tick:
		                fmt.Println("queue len =", len(values))
		            case <- tm:
		                fmt.Println("Bye")
		                return
		        }
		    }
		}
&emsp;&emsp;可以看到queue len一直在增加.

# 传统同步机制

&emsp;&emsp;上面的生产消费案例，我们没有使用共享内存，锁等手段来实现通信。这就是CSP模型，但是golang也有传统同步机制，尽量少用
- WaitGroup
- Mutex
- Cond

&emsp;&emsp;接下来看一下Mutex互斥量的使用
		package main
		
		type atomicInt int
		
		func (a *atomicInt) increment() {
		    *a++
		}
		
		func (a *atomicInt) get() int {
		    return int(*a)
		}
		
		func main() {
		    var a atomicInt
		    a.increment()
		    go func() {
		        a.increment()
		    }()
		    time.Sleep(time.Millsecond)
		    fmt.Println(a)
		}
		$ go run -race atomic.go
		WARNING: DATA RACE
		Read at 0x00c420016060 by main goroutine
		Previous write at 0x00c420016060 by goroutine 6
&emsp;&emsp;上面就出现了数据冲突，一边读一边写。我们加上锁看看

		package main
		
		type atomicInt struct {
		    value int
		    lock  sync.Mutex
		}
		
		func (a *atomicInt) increment() {
		    a.lock.Lock()
		    defer a.lock.Unlock()
		    a.value++
		}
		
		func (a *atomicInt) get() int {
		    a.lock.Lock()
		    defer a.lock.Unlock()
		    return a.value
		}
		
		func main() {
		    var a atomicInt
		    a.increment()
		    go func() {
		        a.increment()
		    }()
		    time.Sleep(time.Millsecond)
		    fmt.Println(a.get())
		}
		$ go run -race atomic.go
		2
&emsp;&emsp;加上锁之后，没有数据冲突了。

> 对于原子化的操作，golang本身提供了atomic包供我们使用，自己就不必再实现一边了。

