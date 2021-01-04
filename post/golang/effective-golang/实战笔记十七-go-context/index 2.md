
>转载自 https://www.flysnow.org

## 前言
控制并发有两种经典的方式，一种是WaitGroup，另外一种就是Context，今天我就谈谈Context。
### 什么是WaitGroup
WaitGroup以前我们在并发的时候介绍过，它是一种控制并发的方式，它的这种方式是控制多个goroutine同时完成。
		func main(){
		    var wg sync.WaitGroup
		    wg.Add(2)
		    go func() {
		        time.Sleep(2 * time.Second)
		        fmt.Println("1号完成")
		        wg.Done()
		    }()
		    go func() {
		        time.Sleep(2 * time.Second)
		        time.Println("2号完成")
		        wg.Done()
		    }()
		    wg.Wait()
		    fmt.Println("好了，大家都干完了，放工")
		}
这是一个很简答的例子，一定要例子中的2个goroutine同时做完，才算是完成，先做好的就要等着其它未完成的，所有的goroutine都要全部完成才可以。

`适用`：好多个goroutine协同做一件事情的时候，因为每个goroutinr做的都是这件事的一部分，只有全部完成，这件事才算完成，这是就是等待的方式。

在实际的业务中，我们可能会有这么一种场景：需要我们主动的通知某一个goroutine的结束。比如我们开启一个后台goroutine一直做事情，比如监控，现在不需要了，就需要通知这个监控goroutine结束，不然它会一直跑，就泄漏了。

### chan通知
我们都知道一个goroutine启动后，我们是无法控制它的，大部分情况是等待它自己结束，那么如果这个goroutine是一个不会自己结束的后台goroutine呢？比如监控等，会一直运行。

这种情况，一个傻瓜式的办法就是全局变量，其它地方通过修改这个变量完成结束通知，然后后台goroutine不停的检查这个变量，如果发现被通知关闭了，就自我结束。

这种方式也可以，但是首先我们要保证这个变量在多线程下的安全，基于此，有一种更好的方式：`chan + select`

		func main() {
			stop := make(chan bool)
		
			go func() {
				for {
					select {
					case <-stop:
						fmt.Println("监控退出，停止了...")
						return
					default:
						fmt.Println("goroutine监控中...")
						time.Sleep(2 * time.Second)
					}
				}
			}()
		
			time.Sleep(10 * time.Second)
			fmt.Println("可以了，通知监控停止")
			stop<- true
			//为了检测监控过是否停止，如果没有监控输出，就表示停止了
			time.Sleep(5 * time.Second)
		
		}
例子中我们定义了一个stop的chan，通知它结束后台goroutine。实现也非常简单，在后台goroutine中，使用select判断stop是否可以接收到值，如果可以接收到，就表示退出停止了；如果没有接收到，机会执行default里的监控逻辑，继续监控，只收到stop的通知。

有了以上的逻辑，我们就可以在其它goroutine中给stopchan发送消息了，例子中是在main goroutine中发送的，控制让这个监控的groutine结束。

发送了stop<-true结束的指令后，我这里使用time.SLeep(5 * time.Second)故意停顿5秒来检测我们结束监控goroutine是否成功。如果成功的话，不会再有goroutine监控中的输出了；如果没有成功，监控goroutine就会继续打印goroutine监控中输出。

这种chan + select的方式，是比较优雅的结束一个goroutine的方式。不过这种方式也有局限性，如果有很多goroutinr都粗腰控制结束怎么办呢？如果这些goroutine又衍生了其它更多的goroutine怎么办呢？如果一层层的无穷无尽的goroutine呢？这就非常复杂了，即使我们定义很多chan也很难解决这个问题，因为goroutine的关系链导致了这种场景非常复杂。

### 初识Context
上面说的这种场景是存在的，比如一个网络请求Request，每一个Request都需要开启一个goroutine做一些事情，这些goroutine又可能会开启其它的goroutine。所以我们需要一种可以`跟踪goroutine`的方案，才可以达到控制它们的目的，这就是Go为我们提供的Context，称之为上下文非常贴切。

下面我们就使用context重写上面的示例。
		func main() {
			ctx, cancel := context.WithCancel(context.Background())
			go func(ctx context.Context) {
				for {
					select {
					case <-ctx.Done():
						fmt.Println("监控退出，停止了...")
						return
					default:
						fmt.Println("goroutine监控中...")
						time.Sleep(2 * time.Second)
					}
				}
			}(ctx)
		
			time.Sleep(10 * time.Second)
			fmt.Println("可以了，通知监控停止")
			cancel()
			//为了检测监控过是否停止，如果没有监控输出，就表示停止了
			time.Sleep(5 * time.Second)
		
		}
重写表简单，就是把原来的chan stop换成context，使用context跟踪goroutine，以便进行控制，比如结束等。

context.background返回一个空的context，这个空的context一般用于整个context树的根节点。然后我们使用context.withcalcel(parent)函数，创建一个可以取消的子context，然后当作参数传给goroutine使用，这样就可以使用这个子context跟踪这个goroutine。

在goroutine中，使用select调用<-ctx.Done()判断是否要结束，如果接收到值的话，就可以返回结束goroutine了；如果接收不到，就会继续监控。

那么是如何发送结束指令的呢？这就是示例中的cancel函数了，它是我们调用context.withcancel(parent)函数生成子context的时候返回的，第二个返回值就是这个取消函数，它是cancelfunc类型。我们调用它就可以发出取消指令，然后我们的监控goroutine就会收到信号，就会返回结束。

### context控制多个goroutine
使用context控制一个goroutine的例子如上，非常简单，下面我们看看多个goroutine的例子，其实也非常简单。
		func main() {
			ctx, cancel := context.WithCancel(context.Background())
			go watch(ctx,"【监控1】")
			go watch(ctx,"【监控2】")
			go watch(ctx,"【监控3】")
		
			time.Sleep(10 * time.Second)
			fmt.Println("可以了，通知监控停止")
			cancel()
			//为了检测监控过是否停止，如果没有监控输出，就表示停止了
			time.Sleep(5 * time.Second)
		}
		
		func watch(ctx context.Context, name string) {
			for {
				select {
				case <-ctx.Done():
					fmt.Println(name,"监控退出，停止了...")
					return
				default:
					fmt.Println(name,"goroutine监控中...")
					time.Sleep(2 * time.Second)
				}
			}
		}
示例中启动了3个监控goroutine进行不断的监控，每一个都使用了context进行跟踪，当我们使用cancel函数通知取消时，这三个goroutine都会被结束。这是context的控制能力，它就像一个控制器一样，按下开关后，所有基于这个context或者衍生的子context都会收到通知，这时就可以进行清理操作了，最终释放goroutine，这就优雅的解决了goroutine启动后不可控的问题。
### Context接口
context的接口定义的比较简介，我们看下这个接口的方法。
		type context interface {
		    Deadline() (deadline time.TIme, ok bool)
		    Done() <-chan struct{}
		    Err() error
		    Value(key interface{}) interface{}
		}
这个接口共有4个方法，了解这些方法的含义非常重要，这样我们才可以更好的使用它们。

`Deadline`方法是获取设置的截止时间的意思，第一个返回式是截止时间，到了这个时间点，context会自动发起取消请求；第二个返回值ok==false时表示没有设置截止时间，如果需要取消的话，需要调用取消函数进行取消。

`Done`方法返回一个只读的chan，类型为struct{}，我们在goroutine中，如果该方法返回的chan可以读取，则意味着parent context已经发起了取消请求，我们通过Done方法收到这个信号后，就应该清理操作，然后退出goroutine，释放资源。

`Err`方法返回取消的错误原因，因为什么Context被取消。

`Value`方法获取该Context上绑定的值，是一个键值对，所以要通过一个Key才可以获取对应的值，这个值一般是线程安全的。

以上四个方法中常用的就是Done了，如果Context取消的时候，我们就可以得到一个关闭的chan，关闭的chan是可以读取的，所以只要可以去读的时候，就以为这收到context取消的信号了，以下是这个方法的经典用法。
		func Stream(ctx context.Context, out chan<- Value) error {
		  	for {
		  		v, err := DoSomething(ctx)
		  		if err != nil {
		  			return err
		  		}
		  		select {
		  		case <-ctx.Done():
		         //该取消了，返回取消原因
		  		    return ctx.Err()
		  		case out <- v:
		  		}
		  	}
		  }
context接口并不需要我们实现，Go内置已经帮我们实现了2个，我们代码中最开始都是以这两个内置的作为最顶层的parent context衍生出更多的子context。
		var (
			background = new(emptyCtx)
			todo       = new(emptyCtx)
		)
		
		func Background() Context {
			return background
		}
		
		func TODO() Context {
			return todo
		}
一个是Background，主要用于main函数、初始化以及测试代码中，作为context这个树结构的最顶层的context，也就是根context。

一个是TODO，它目前还不知道具体的使用场景，如果我们不知道该使用什么context的时候，可以使用这个。

它们两个本质上都是emptyCtx结构体类型，是一个不可取消，没有设置截止时间，没有携带任何值的context
		type emptyCtx int
		
		func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
			return
		}
		
		func (*emptyCtx) Done() <-chan struct{} {
			return nil
		}
		
		func (*emptyCtx) Err() error {
			return nil
		}
		
		func (*emptyCtx) Value(key interface{}) interface{} {
			return nil
		}
这就是emptyCtx实现Context接口的方法，可以看到，这些方法什么都没做，返回的都是nil或者零值。
### Context的衍生
有了如上的根context，那么是如何衍生更多的子context的呢？这就要靠context包为我们提供的with系列的函数了。
		func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
		func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
		func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
		func WithValue(parent Context, key, val interface{}) Context
这四个with函数，接收的都有一个parent参数，就是父context，我们要基于这个父context创建出子context的意思，这中方法可以理解为子context对父context的继承，也可以理解为父context的衍生。

通过这些函数，就创建了一个context树，树的每个节点都可以有任意多个子节点，节点层级可以有任意多个。

`WithCancel`函数，传递一个父Context作为参数，返回一个子Context，以及一个取消函数用来取消Context。

`WIthDeadline`函数，和WithCancel差不多，它会多传递一个截止时间参数，意味着到了这个时间点，会自动取消Conteext，当然我们也可以不等到这个时候，可以提前通过取消函数进行取消。

`WithTimeout`和`WIthDeadline`基本上一样，这个表示是超时自动取消，是多少时间后自动取消Context的意思。

`WithValue`函数和取消Context无关，它是为了生成一个绑定了一个键值对数据的COntext，这个绑定的数据可以通过Context.Value方法访问到，后面我们会专门讲。

大家可以留意到前三个函数都返回一个取消函数CancekFunc，这是一个函数类型，它的定义非常简单。
		type CancelFunc func()
这就是取消函数的类型，该函数可以取消一个context，以及这个节点Context下所有的context，不管有多少层级。
### WithValue传递元数据
通过Context我们也可以传递一些必须的元数据，这些数据会附加在context上以供使用。
		var key string="name"
		
		func main() {
			ctx, cancel := context.WithCancel(context.Background())
			//附加值
			valueCtx:=context.WithValue(ctx,key,"【监控1】")
			go watch(valueCtx)
			time.Sleep(10 * time.Second)
			fmt.Println("可以了，通知监控停止")
			cancel()
			//为了检测监控过是否停止，如果没有监控输出，就表示停止了
			time.Sleep(5 * time.Second)
		}
		
		func watch(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					//取出值
					fmt.Println(ctx.Value(key),"监控退出，停止了...")
					return
				default:
					//取出值
					fmt.Println(ctx.Value(key),"goroutine监控中...")
					time.Sleep(2 * time.Second)
				}
			}
		}
在前面的例子中，我们通过传递参数的方式，把name的值传递给监控函数。在这个例子中，我们实现一样的效果，但是通过contextvalue的方式。

我们可以使用context.withvalue方法附加一对k-v的键值对，这里key必须是等价性的，也就是具有可比性；value值要是线程安全的。

这样我们就生成了一个新的context，这个新的context带的context带有这个键值对，在使用的时候，可以通过value方法读取ctx.value(key)。

记住，使用withvalue传值，一般是必须的值，不要什么值都传递。
### Context使用原则
1. 不要把Context放在结构体中，要以参数的方式传递；
2. 以Context作为参数的函数方法，应该把Context作为第一个参数，放在第一位；
3. 给一个函数方法传递context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO；
4. Context的value相关方法应该传递必须的数据，不要什么数据都使用这个传递；
5. Context是线程安全的，可以放心在多个goroutine中传递。
