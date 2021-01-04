
>转载自 https://www.flysnow.org

## 简介
这篇通过一个例子，演示使用通道来监控程序的执行时间，生命周期，甚至终止程序等等。我们这个程序叫runner，我们可以称之为执行者，它可以在后台执行任何任务，而且我们还可以控制这个执行者，比如强制终止它等。

现在开始吧，运用我们前面十几篇连载的知识，来构建我们的runner，使用一个结构体类型就可以。
		//一个执行者，可以执行任何任务，但是这些任务是限制完成的
		//该执行者可以通过发送终止信号终止它
		type Runner struct{
		    tasks []func(int) // 要执行的任务
		    complete chan err //用于通知任务全部完成
		    timeout <-chan time.Time //这些任务在多久内完成
		    interrupt chan os.Signal    //可以控制强制终止的信息,中断信号
		}
示例中，我们定义了一个结构体类型runner，这个runner包含了要执行哪些任务tasks，然后使用complete通知任务是否全部完成，不过这个执行者是有时间限制的，这既是timeout，如果在限定的时间内没有完成，就会接收到超时的通知，如果完成了就会接收到完成的通知。注意这里的timeout是单向通道，只能接收。

complete定义为error类型的通道，是为了当执行任务出现问题时返回错误的原因，如果没有出现错误，返回是nil。
此外，我们还定义了一个中断的信号，让我们可以随时的终止执行者。

有了结构体，我们需要一个工厂，用于返回我们需要的runner。
		func New(tm time.Duration) *Runner{
		    return &Runner{
		        complete:make(chan error),
		        timeout:time.After(tm),
		        interrupt:make(chan os.Signal, 1)
		    }
		}
它只有一个参数妙用来设置这个执行者的超时时间。这个超时区间被我们传递给了time.after函数，这个函数可以在tm时间后，会伙同一个time.time类型的只能接收的单向通道，来告诉我们已经到时间了。

complete是一个无缓冲通道，也就是同步通道，因为我们要使用它来控制我们整个程序是否终止，所以它必须是同步的，要让main等待，一直到任务完成或者被强制终止。

interrupt是一个有缓冲的通道，这样做是因为，我们可以至少接收到一个操作系统的中断信息，这样g runtime在发送这个信号的时候不会被阻塞，如果是无缓冲通道就会阻塞。

系统信号是什么？比如我们在程序执行的时候按下ctrl-c，这既是中断信号。

我们这里初始化了结构体的三个字段，而执行的任务tasks没有初始化，默认就是nil，因为它是一个切片。定义一个方法，通过方法给执行者添加需要的任务。
		//将需要执行的任务，添加到runner里
		func (r *Runner) Add(tasks ...func(int)){
		    r.tasks = append(t.tasks, tasks...)
		}
到这里我们需要的runner，如何添加任务，如果获取获取一个执行者都有了，下面就开始runner如何运行任务？如何在运行时强制中断任务？在这些处理之前，我们先来定义两个我们的两个错误变量，以便接下来的代码实例使用。
		var ErrTimeOut = errors.New("执行者执行超时")
		var ErrInterrupt = errors.New("执行者被中断")
		//执行任务，执行的过程中接收到中断信号，返回中断错误
		//如果任务全部执行完，还没有接收到中断信号，则返回nil
		func (r *Runner) run() error {
		    for id, task := range r.tasks {
		        if r.isInterrupt(){
		            return ErrInterrupt
		        }
		        task(id)
		    }
		    return nil
		}
		
		//检查是否接收到了中断信号
		func (r *Runner) isInterrupt() bool {
		    select {
		    case <-r.interrupt:
		        signal.Stop(r.interrupt)
		        return true
		    default:
		        return false
		    }
		}
		
		func (r *runner) Run() error {
		   for id, task := range r.tasks {
		      if r.isInterrupt() {
		         return ErrInterrupt
		      }
		      task(id)
		   }
		   return nil
		}
		
		func (r *runner) Start() error {
		   signal.Notify(r.interrupt, os.Interrupt)
		
		   go func() { r.complete <- r.Run() }()
		
		   select {
		   // 阻塞等待
		   case err := <-r.complete:
		      return err
		   case <-r.timeout:
		      return ErrTimeOut
		   }
		}
		func TestRunner_Start(t *testing.T) {
		   log.Println("...开始执行任务...")
		
		   ti := 3 * time.Second
		   r := RunnerNew(ti)
		   r.Add(createTask(), createTask(), createTask())
		   if err := r.Start(); err != nil {
		      switch err {
		      case ErrTimeOut:
		         log.Println(err)
		         os.Exit(1)
		      case ErrInterrupt:
		         log.Println(err)
		         os.Exit(2)
		      }
		   }
		
		   log.Println("执行任务结束")
		}
		
		func createTask() func(int) {
		   return func(i int) {
		      log.Printf("...正在执行任务%d...", i)
		      time.Sleep(time.Duration(i) * time.Second)
		   }
		}
		
