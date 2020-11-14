
> 转载自 https://www.cnblogs.com/sparkdev/

> 转载自 https://go-zh.org/ref/mem golang官网内存模型分析

> 转载自 https://taohuawu.club/high-performance-implementation-of-goroutine-pool 

## 忠告
&emsp;&emsp;程序在修改被多个Go程同时访问的数据时必须序列化该访问。要序列化访问，需要通过Channel操作，或sync,sync/atomic包中的同步原语来保护数据。

## 内存模型
&emsp;&emsp;请参阅[golang官网内存模型分析](https://go-zh.org/ref/mem)，对于其中的信道通信案例解析，请参阅[信道案例解析](http://facedamon.github.io/post/golang/base/8.channel/#信道通信-解析go官网案例)

## 线程模型

### Goroutine & Scheduler

&emsp;&emsp;Goroutine，Go语言基于并发编程给出的自家的解决方案。Goroutine是什么？通常goroutine会被
当作coroutine的golang实现，从比较粗浅的层面来看，这种认知也算合理，但实际上，goroutine并非传统意义上的协程，现在主流的线程模型分为三种：`内核级线程模型`、`用户级线程模型`和`混合型线程模型`，传统的协程库属于用户级线程模型，而gourine和它的`Go Scheduler`在底层实现上属于混合型线程模型，因此，有时候为了方面理解可以简单把goroutine类比成协程，**但心里一定要有个清楚的认知---goroutine并不等同于协程**
&emsp;&emsp;这三种线程模型它们之间最大的差异就在于用户线程与内核调度实体之间的对应关系上。内核调度实体KSE就是指可以被操作系统内核调度器调度的对象实体。对于KSE的详细解释可以参考Linux调度器。

### 用户级线程模型

&emsp;&emsp;用户线程与KSE是多对一(N:1)的映射模型，多个用户线程的一般从属于单个进程并且多线程的调度是由用户自己的线程库来完成：线程的创建、销毁以及对线程之间的协调等操作。无需借助系统调用来实现。

&emsp;&emsp;一个进程中所有创建的线程都只和同一个KSE在运行时动态绑定，也就是说，操作系统只知道用户进程而对其中的线程是无感知的，内核的所有调度都是基于用户进程。许多语言实现的协程库基本上都属于这种(python,gevent).由于线程调度是在用户层面完成的，也就是相较于内核调度`不需要让CPU在用户态和内核态之间切换`，这种实现方式相比内核级线程可以做的很清凉，对系统资源的消耗会小很多，因此可以创建的线程数量与上下文切换所花费的代价也会小很多。但该模型有一个原罪：`并不能做到真正意义上的并发`，假设在某个用户进程上的某个用户线程因为一个阻塞调用而被CPU给中断了，那么该进程内的所有线程都被阻塞(因为单个用户进程内的线程自调度是没有CPU时钟中断的，从而没有轮转调度)。整个进程被刮起。即便是多CPU的机器，也无济于事，因为在用户级线程模型下，一个CPU关联运行的是整个用户进程。所以很多的协程库会把自己一些阻塞的操作重新封装为完全的非阻塞形式，然后在以前要阻塞的点上，主动让出自己，并通过某种方式通知或唤醒其它待执行的用户线程在该KSE上运行，从而避免了内核调度器由于KSE阻塞而做上下文切换，这样整个进程也不会阻塞了。

### 内核级线程模型

&emsp;&emsp;用户线程与内核线程KSE是一对一(1:1)的映射模型，也就是每个用户线程绑定一个实际的内核线程，而`线程的调度则完全交付给操作系统内核`去做: 线程的创建、终止以及同步都基于内核提供的系统调用来完成，大部分编程语言的线程库(Java.Thread、C++11 std:thread)都是对操作系统线程的一层封装，创建出来的每个线程与一个独立的KSE静态绑定，因此其调度完全由操作系统内核调度器去做，也就是说,`一个进程创建出来的多个线程每一个都绑定一个KSE`。这种模型的优势和劣势同样明显：优势是实现简单，直接借助曹旭哦系统内核的线程及调度器，所以CPU可以快速切换调度线程，于是多个线程可以同时运行，因此相较于用户级线程模型它真正做到了并行；但它的劣势是，由于直接借助了操作系统内核来创建、销毁和以及索格线程之间的上下文切换和调度，因此资源成本大幅上涨，且对性能影响很大。

### 混合型线程模型

&emsp;&emsp;混合型线程模型是博采众长之后的产物，充分吸收前两种模型的优点且尽量避免它们的缺点。在此模型下，用户线程与内核KSE是多对多(N:N)的映射模型。首先，区别于用户级线程模型，混合型模型中的一个进程可以与多个内核线程KSE关联，也就是说一个进程内的多个线程可以分别绑定一个自己的KSE，这点和内核级线程模型相似；其次，有区别于内核级线程模型，它的进程里的线程并不与KSE唯一绑定，而是可以多个用户线程映射到同一个KSE，当某个KSE因为其绑定的线程的阻塞操作被内核调度出CPU时，其关联的进程中其余用户线程可以重新与其它KSE绑定运行。所以，`混合型线程模型既不是用户级线程那种完全靠自己调度的也不是内核级线程完全靠操作系统调度的，而是自身调度与系统调度协同工作`，因为这种模型的高度复杂性，操作系统内核开发者一般不会使用，所以更多时候是作为第三方库的形式出现，而Go语言中的runtime调度器就是采用这种实现方案，实现了Goroutine与KSE之间的动态关联，不过Go语言的实现更加高级和优雅；

### G-P-M 模型

&emsp;&emsp;每一个操作系统线程都有一个固定大小的内存块(一般是2MB)来做栈，这个栈会用来存储当前正在被调用或挂起(指在调用其它函数时)的函数内部变量。这个固定大小的栈同时很大又很小。因为2MB的栈对于一个小小的goroutine来说是很大的内存浪费，而对于一些复杂的任务(深度递归)来说又显得太小。因此，Golang做了自己的线程。

&emsp;&emsp;在Golang中，每一个goroutine是一个独立的执行单元，相较于每个操作系统线程固定分配的2MB内存模式，`goroutine的栈采取了动态扩容方式，初始时仅为2KB，随着任务执行按需增长，最大可达1GB`(64位机器最大是1G，32位机器最大时256M)，且`完全由golang自己的调度器Go Scheduler来调度`。此外，GC还会周期性地将不再使用的内存回收，收缩栈空间。因此，Go程序可以同时并发成千上万个goroutine是得益于它强劲的调度器和高效的内存模型。Go的创造者对goroutine的定位本身就是重量级的，因为他们不仅让goroutine作为golang并发编程的核心组建而且golang中的许多标准库的实现也能到处见到goroutine的身影，比如net/http,甚至语言本身的组建runtime运行时和GC垃圾回收器都是运行在goroutine上的，作者对goroutine的厚望可见一斑，

&emsp;&emsp;任何用户线程最终肯定都要交给操作系统线程来执行，goroutine也不例外，但是G并不直接绑定操作系统线程运行，而是由Go Scheduler中的`P - Logical Processor`(逻辑处理器)来作为两者的中介，P可以看作是一个抽象的资源或者上下文，一个P绑定一个操作系统，在golang的实现里把操作系统线程抽象成一个数据结构：M，G实际上是由M通过P来调度运行的，但是在G的层面来看，P提供了G运行所需的一切资源和环境，因此在G看来P就是运行它的”CPU“，由G、P、M这三种由Go抽象出来的实现，最终形成了Go调度器的基本结构。

- G：表示Goroutine，每个Goroutine对应一个G结构体，G存储Goroutine的运行堆栈、状态、以及任务函数，可重用。G并非执行体，每个G需要绑定到P才能被调度执行。
- P：Processor，表示逻辑处理器，对G来说，P相当于CPU核，G只有绑定到P(local runq中)才能被调度。对M来说，P提供了相关的执行环境(Context)，如内存分配状态(macache),任务队列(G)等，P的数量决定了系统内最大可并行的G的数量(前提：物理CPU核数 >= P的数量)，P的数量由`runtime.gomaxprocs`决定，但是不论gomaxprocs设置多大，P的数量最大为256.
- M：Machine，操作系统线程抽象，代表着真正执行计算的资源，在绑定有效的P后，进入schedule循环；而`schedule循环的机制大致是从Global队列、P的Local队列以及wait队列中获取G，切换到G的执行栈上并执行G的函数，调用goexit做清理工作并回到M，如此反复`。M并不保留G的状态，这是G可以跨M调度的基础，M的数量是不定的，由go runtime做调整，为了防止创建过多M，目前默认最大限制为10000个。

&emsp;&emsp;关于P，我们在说一些，在Go1.0发布的时候，它的调度器其实是G-M模型，没有P，调度过程完全由G和M组成，这个模型暴露出一些问题：
1. 单一全局互斥锁(Sched.Lock)和集中状态存储的存在导致所有goroutine相关操作，比如：创建、重新调度等都要上锁；
2. goroutine传递问题：M经常在M之间传递可运行的goroutine，这导致调度延迟增大以及额外的性能损耗；
3. 每个M做内存缓存，导致内存占用过高，数据局部性差；
4. 由于syscall调用而形成的剧烈的worker thread阻塞和解除阻塞，导致额外的性能损耗。

&emsp;&emsp;这些问题是在太扎眼了，导致Go1.0虽然号称原声支持并发，却在并发性能上一直饱受诟病，然后，Go委员会中一个核心大佬实在看不下去了，亲自下场重新设计和实现了Go调度器P，并且实现了一个叫做`work-stealing`的调度算法：

- 每个P维护一个G的本地队列；
- 当一个G被创建出来，或者变为可执行状态时，就把它放到P的可执行队列中；
- 当一个G在M里执行结束后，P会从队列中把该G取出；如果此时P的队列为空，即没有其它G可执行，M就随机选择另外一个P，从其可执行的G队列中取走一半。
- 该算法避免了在goroutine调度时使用全局锁

> 至此， Go调度器基本模型确立

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/G-P-M.png)

## G-P-M 模型调度

&emsp;&emsp;Go调度器工作时会维护两种用来保存G的任务队列：一种是一个Global任务队列，一种是每个P维护的Local任务队列。
&emsp;&emsp;当创建一个新的goroutinr的时候，它会优先被放入P的本地队列。为了运行goroutine，M需要持有一个P，接着M会启动一个操作系统线程，循环从P的本地队列里取出一个goroutine并执行。当然还有上下文提及的work-stealing调度算法：当M执行完当前P的Local队列里所有G后，P也不会就这么在那躺着啥也不干，它会先尝试从Global队列寻找G来执行，如果Global队列为空，它会随机挑选另外一个P，从它的队列中拿走一半的G到自己的队列中执行。

## channel

&emsp;&emsp;channel是go语言级别提供的goroutine之间的通信方式，可以使用channel在两个或多个goroutine之间传递消息。使用channel发送和接收所需的共享资源，可以在goroutine之间消除竞争条件。

### 无缓冲channel

&emsp;&emsp;unbuffered channe 是指在接收前没有能力保存任何值的channel。这种类型的channel要求发送goroutine和接收goroutine同时准备好，才能完成发送和接收操作。`如果两个goroutine没有同时准备好，channel会导致先执行发送或接收操作的goroutine阻塞等待`。这种对通道进行发送和接收的交互行为本身就是同步的。其中任意一个操作都无法离开另一个操作单独存在。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/unbuffered-channel.png)

1. 两个goroutine都到达通道，但两个都没有开始执行数据的发送或接收。
2. 左侧的goroutine将它的手伸进了通道，这模拟向通道发送数据的行为。这时，这个goroutine会在通道中被锁住，直到交换完成。
3. 右侧的goroutine将它的手放入通道，这模拟了从通道里接收数据。这个goroutine一样也会在通道中被锁住，直到交换完成。
4. 交换数据
5. 两个goroutine都将它们手从通道里拿出来，这模拟了被锁住的goroutine得到释放。两个goroutine现在都可以去做别的事情了。

&emsp;&emsp;下面的例子模拟异常网球比赛。在网球比赛中，两位选手会把球在两个人之间来回传递。选手总是处在以下两种状态之一：要么在等待，要么将球打出。可以使用两个goroutine来模拟网球比赛，并使用无缓冲的通道来模拟球的来回：

		package main
		
		import (
		    "math/rand"
		    "sync"
		    "time"
		    "fmt"
		)
		
		var wg sync.WaitGroup
		
		func init() {
		    rand.Seed(time.Now().UnixNano())
		}
		
		func main() {
		    // 创建一个无缓冲的通道
		    court := make(chan int)
		
		    // 添加两个任务
		    wg.Add(2)
		
		    go player("Nick", court)
		    go player("Jack", court)
		
		    // 发球 ，在通道接收准备之前会阻塞
		    court <- 1
		
		    // 等待游戏结束
		    wg.Wait()
		}
		
		func player(name string, court int) {
		    defer wg.Done()
		
		    for {
		        // 在球未被打过来时，会阻塞等待
		        ball, ok := <- court
		        if !ok {
		            // 如果通道被关闭，我们就赢了
		            fmt.Printf("Player %s Won\n", name)
		            return
		        }
		        //  选随机数，然后用这个数来判断我们是否丢球
		        n := rand.Intn(100)
		        if n % 5 == 0 {
		            fmt.Printf("Player %s Missed\n", name)
		            // 关闭通道，表示我们输了
		            close(court)
		            return
		        }
		
		        // 显示击球数，并将击球数加1
		        fmt.Printf("Player %s Hit %d\n", name, ball)
		        ball++
		
		        // 将球打向对手
		        court <- ball
		    }
		}
		

		Player Zack Hit 1
		Player Nick Hit 2
		Player Jack Hit 3
		Player Nick Hit 4
		Player Jack Missed
		Player Nick Won

&emsp;&emsp;在main函数中创建了一个int类型的无缓冲通道，使用该通道让两个goroutine在击球时能够相互同步。然后创建了参与比赛的两个goroutine。在这个时候，两个goroutine都阻塞住等待击球。court <- 1模拟发球，将球发到通道里，程序开始执行这个比赛，直到某个goroutine输掉比赛。

&emsp;&emsp;在player函数里，主要运行一个无限循环的for语句。在这个循环里，是玩游戏的过程。goroutine从通道接收数据，用来表示等待接球。这个接收动作会阻塞住goroutine，直到有数据发送到通道里。通道的接收动作返回时，会检测ok标志是否为false。如果为false，表示通道已经关闭，游戏结束。在这个模拟程序中，使用随机数来决定goroutine是否击中了球，如果击中了球，就把ball的值递增1，并将ball作为球重新放入通道，发送给另一位选手。在这个时刻，两个goroutine都被阻塞，直到交换完成。最终，因某个goroutine没有打中球会把通道关闭。之后两个goroutine都会返回，通过defer声明Done，程序终止

### buffered channel

&emsp;&emsp;带缓冲的channel是一种在被接收前能存储一个或多个值的通道。这种类型的通道并不强制要求goroutine之间必须同时完成发送和接收。`通道会阻塞发送和接收动作的条件也会不同。只有在通道中没有要接收的值时，接收动作才会阻塞。只有在通道没有可用缓冲区容纳被发送的值时，发送动作才会阻塞。`这导致有缓冲的通道和无缓冲的通道之间的一个很大的不同：**无缓冲的通道保证进行发送和接收的goroutine会在同一时间进行数据交换；有缓冲的通道没有这种保证。**可以通过下面的图示形象地理解两个goroutine分别向带缓冲的通道里增加一个值和从带换种的通道里移除一个值。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/buffered-channel.png)

1. 右侧的goroutine正在从通道接收一个值
2. 右侧的这个goroutine独立完成了接收值的动作，而左侧的goroutine正在发送一个新值到通道里。
3. 左侧的goroutine还向通道发送新值，而右侧的goroutine正在从通道接收另外一个值。这个步骤里两个操作既不是同步的，也不会互相阻塞。
4. 所有的发送和接收都完成，而通道里还有几个值，也有一些空间可以存更多的值。

&emsp;&emsp;下面使用一组goroutine来接收并完成任务，带缓冲区的通道提供了一种清晰而直观的方式来实现这个功能：

		package main
		
		import (
		    "math/rand"
		    "sync"
		    "time"
		    "fmt"
		)
		
		const (
		    // 要使用的goroutine的数量
		    numberGorutines = 2
		    // 要处理的工作的数量
		    taskLoad = 5
		)
		
		var wg sync.WaitGroup
		
		func init() {
		    rand.Seed(time.Now().UnixNano())
		}
		
		// 从有缓冲的通道传入工作及工作id
		func worker(tasks chan string, worker int) {
		    defer wg.Done()
		
		    for {
		        task, ok := <- tasks
		        if !ok {
		            fmt.Printf("Worker: %d: Shutting Down\n", worker)
		            return
		        }
		        // 显示开始工作了
		        fmt.Printf("Worker: %d: Started %s\n, worker, task)
		
		        // 随机等一段时间来模拟工作
		        sleep := rand.Int63n(100)
		        time.Sleep(time.Duration(sleep) * time.Millsecond)
		
		        // 显示完成工作
		        fmt.Printf("Worker: %d: Completed %s\n", worker, task)
		    }
		}
		
		func main() {
		    // 创建一个有缓冲的通道来管理工作
		    tasks := make(chan string, taskLoad)
		    
		    wg.Add(numberGoroutines)
		
		    // create worker group
		    for gr := 1; gr <= numberGoroutines; gr++ {
		        go worker(tasks, gr)
		    }
		
		    // create worker
		    for post := 1; post <= taskLoad; post++ {
		        tasks <- fmt.Sprintf("Task: %d", post)
		    }
		
		    close(tasks)
		
		    wg.Wait()
		}

&emsp;&emsp;运行上面的程序，输出结果大致如下：

		Worker: 2: Started Task: 1
		Worker: 1: Started Task: 2
		Worker: 1: Completed Task: 2
		Worker: 1: Started Task: 3
		Worker: 1: Completed Task: 3
		Worker: 1: Started Task: 4
		Worker: 2: Completed Task: 1
		Worker: 2: Started Task: 5
		Worker: 1: Completed Task: 4
		Worker: 1: Shutting Down
		Worker: 2: Completed Task: 5
		Worker: 2: Shutting Down
&emsp;&emsp;关闭通道的代码非常重要。`当通道关闭后，goroutine依旧可以从通道接收数据，但是不能再向通道里发送数据。`能够从已经关闭的通道接收数据这一点非常重要，因为这允许通道关闭后依旧能取出其中缓冲的全部值，而不会有数据丢失。从一个已经关闭没有数据的通道里获取数据，总会立刻返回，并返回一个通道类型的零值。如果在获取通道时还加入了可选的标志，就能得到通道的状态信息。

### 处理超时

&emsp;&emsp;使用channel时需要小心，比如对于下面的简单用法：
		i := <- ch

&emsp;&emsp;碰到永远没有往ch中写入数据的情况，那么这个读取动作将永远也无法从ch中读取到数据，导致的结果就是整个goroutine永远阻塞并且没有挽回的机会。如果channel只是被同一个开发者使用，那么出问题的可能性还低一些。但如果一旦对外公开，就必须考虑到最差情况并对程序进行维护。

&emsp;&emsp;Golang没有提供直接的超时处理机制，但可以利用select机制变通地解决。因为select的特点是只要其中一个case已经完成，程序就会继续往下执行，而不会考虑其它的case。基于此特性我们来实现一个channel的超时机制。

		ch := make(chan int)
		// 首先实现并执行一个匿名的超时等待函数
		timeout := make(chan bool, 1)
		go func(){
		    time.Sleep(le9)
		    timeout <- true
		}()
		
		select {
		    case <- ch:
		        // unbuffered channel blocking
		    case <- timeout:
		        fmt.Printf("%s\n", "Timeout occurred.")
		}

### 关闭channel

&emsp;&emsp;关闭channel非常简单，直接调用Golang内置的close函数就可以了：`close(ch)`
&emsp;&emsp;在关闭了channel之后我们要面对的问题是：`如何判断一个channel是否已经关闭？`其实在从channel中读取数据的同时，还可以获得一个布尔类型的值，该值表示channel是否已关闭：

		x, ok := <- ch
&emsp;&emsp;如果ok为fale，则表示ch已经被关闭。

## 信道通信(解析Go官网案例)

### 有缓冲通道的发送操作总在对应的接收操作完成前发生

		var c = make(chan int, 10)
		var a string
		
		func f() {
		    a = "hello, world"
		    c <- 0
		}
		
		func main() {
		    go f()
		    <- c
		    print(a)
		}
&emsp;&emsp;可保证打印出"hello, world".该程序首先对a进行写入，然后在c上发送信号，随后从c接收对应的信号，卒后执行print函数。

> 通过上面的buffered channel的解释, 我们知道了 "通道会阻塞发送和接收动作的条件也会不同。只有在通道中没有要接收的值时，接收动作才会阻塞。只有在通道没有可用缓冲区容纳被发送的值时，发送动作才会阻塞" 细品一下这句话，你就会发现，`<- c`这个I/O操作会在通道还没有数据时阻塞，达到同步的效果。

### 无缓冲信道进行接收，要发生在对该通道进行发送完成之前

		var c = make(chan int)
		var a string
		
		func f() {
		    a = "hello, world"
		    <- c
		}
		
		func main(){
		    go f()
		    c <- 0
		    print(a)
		}

> 通过上面的unbuffered channel的解释，我们知道了 "如果两个goroutine没有同时准备好，channel会导致先执行发送或接收操作的goroutine阻塞等待" ,那么`c <- 0`在f()还没有准备好之前是阻塞状态，也达到了同步效果。

