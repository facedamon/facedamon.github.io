
>此Golang面试题摘自网络，已经找不到作者了，所以使用someone代替。

## 1. 写出下面代码输出内容

		package main
		
		import (
			"fmt"
		)
		
		func main() {
		   defer_call()
		}
		
		func defer_call() {
		   defer func() { fmt.Println("打印前") }()
		   defer func() { fmt.Println("打印中") }()
		   defer func() { fmt.Println("打印后") }()
		   panic("触发异常")
		}

- 考察对defer的理解，defer函数延迟执行，延迟到调用者函数执行return明朗前被执行。多个defer之间LIFO先进后出顺序。故考题中，在Panic触发时结束函数运行，在return前先依次打印：打印后、打印中、打印钱。最后由runtime运行时抛出打印panic异常。

		打印后
		打印中
		打印前
		panic: 触发异常

- 需要注意的是，函数的return value不是原子操作，而是在编译器中人皆为两部分：返回值赋值和return。而defer刚好被插入到末尾的return前执行。故可以在defer函数中修改返回值。如下示例：

		package main
		
		import (
		   "fmt"
		)
		
		func main() {
		       fmt.Println(doubleScore(0)) //0
		       fmt.Println(doubleScore(20.0)) //40
		       fmt.Println(doubleScore(50.0)) //50
		}
		
		func doubleScore(source float32) (score float32) {
		       defer func() {
		        if score < 1 || score >= 100 {
		        //将影响返回值
		        score = source
		        }
		 }()
		 score = source * 2
		 return
		 //或者
		 //return source * 2
		}

## 2.以下代码有什么问题，请说明原因

		package main
		
		import (
		   "fmt"
		)
		
		type student struct {
		   Name string
		   Age  int
		}
		
		func pase_student() map[string]*student {
		   m := make(map[string]*student)
		   stus := []student{
		       {Name: "zhou", Age: 24},
		       {Name: "li", Age: 23},
		       {Name: "wang", Age: 22},
		   }
		   for _, stu := range stus {
		       m[stu.Name] = &stu
		   }
		   return m
		}
		
		func main() {
		   students := pase_student()
		   for k, v := range students {
		      fmt.Printf("key=%s,value=%v \n", k, v)
		   }
		}
		

- 输出的均是相同的值:&{wang 22}
- 因为for遍历时，变量<font color="red">stu指针</font>不变，每次遍历仅进行struct值拷贝，故<font color="red">m[stu.name]=&stu</font>实际上一直指向同一个指针，最终该指针的值为遍历的最后一个struct的值拷贝。
- 修正方案：取数组中原始址的指针

		for i, _ := range stus {
		      stu := stus[i]
		      m[stu.name] = &stu
		}

## 3.下面的代码会输出什么，并说明原因

		func main() {
		   runtime.GOMAXPROCS(1)
		   wg := sync.WaitGroup{}
		   wg.Add(20)
		   for i := 0; i < 10; i++ {
		   	go func() {
		   		fmt.Println("i: ", i)
		   		wg.Done()
		   	}()
		   }
		   for i := 0; i < 10; i++ {
		   	go func(i int) {
		   		fmt.Println("i: ", i)
		   		wg.Done()
		   	}(i)
		   }
		   wg.Wait()
		}

- 将随机输出数字，单前面一个循环中并不会输出所有值
- GOMAXPROCS=1，将影响goroutine的并发，后续代码中的go func相当于串行执行
- 两个for循环内部go func调用参数i的方式是不同的，导致结果完全不同。这也是新手容易遇到的坑。第一个go func中的i是外部for的一个变量，地址不变化。遍历完成后，最终i=10.故go func执行时，i的值始终是10.第二个go func中i是函数参数，与外部for中的i完全是两个变量。为不i奖发生值拷贝，go func内部指向值拷贝地址。

## 4.下面代码会输出什么？

		type People struct{}
		
		func (p *People) ShowA() {
		   fmt.Println("showA")
		   p.ShowB()
		}
		func (p *People) ShowB() {
		   fmt.Println("showB")
		}
		
		type Teacher struct {
		   People
		}
		
		func (t *Teacher) ShowB() {
		   fmt.Println("teacher showB")
		}
		func main() {
		   t := Teacher{}
		   t.ShowA()
		}

- 输出：showA showB
- Go中没有继承，叫组合
- 这里People是匿名组合People。被组合的类型People所包含的方法虽然升级成了外部类型Teacher这个组合类型的方法，但是它们的方法showA调用时接受者并没有发生变化。这里仍然是People。毕竟这个People类型并不知道会被什么类型组合，当然也就无法调用方法时去使用未知的组合者teacher类型的功能。因此这里执行showA时，在执行showB时该函数的接受者时people，而非teacher。

## 5.下面代码会出发异常吗，请详细说明

		func main() {
		   runtime.GOMAXPROCS(1)
		   int_chan := make(chan int, 1)
		   string_chan := make(chan string, 1)
		   int_chan <- 1
		   string_chan <- "hello"
		   select {
		   case value := <-int_chan:
		   	 fmt.Println(value)
		   case value := <-string_chan:
		   	 panic(value)
		   }
		}

- 有可能触发异常，是随机事件
- 单个chan如果无缓冲时，将会阻塞。但结合select可以在多个chan间等待执行。有三点原则：
1. select中只要有一个case能return，则立即执行
2. 当如果同一时间有多个case均能return则伪随机方式抽取任意一个执行
3. 如果没有一个case能return则可以执行default块
- 此考题中的两个case中的两个chan均能return，则会随机执行某个case。故在执行程序时，有可能执行第二个case，触发异常。

## 6.下面代码输出什么？

		func calc(index string, a, b int) int {
		   ret := a + b
		   fmt.Println(index, a, b, ret)
		   return ret
		}
		func main() {
		   a := 1                               //line 1
		   b := 2                               //2
		   defer calc("1", a, calc("10", a, b)) //3
		   a = 0                                //4
		   defer calc("2", a, calc("20", a, b)) //5
		   b = 1                                //6
		}

- 输出结果：<br/>10 1 2 3 <br/>20 0 2 2 <br/> 2 0 2 2 <br> 1 1 3 4
- 解析，在解题之前需要明确两个概念：
1. defer是在函数末尾的return前执行，先进后执行
2. 函数调用时int参数发生值拷贝
- 不管代码顺序如何，defer calc func中参数b必须先计算，故会在运行到第三行时，执行calc("10",a,b)输出：10 1 2 3 得到值3，将cal("1",1,3)存放到延迟后执行函数队列中。执行到第五行时，现行计算cal("20",a,b)即cal("20",0,2)输出：20 0 2 2得到值2，将cal("2",0,2)存放到延迟后执行函数队列中。执行到末尾行，按队列先进后出原则依次执行：cal("2",0,2),cal("1",1,3)，一次输出：2 0 2 2，1 1 3 4

## 7. 请写出以下输出内容

		func main(){
		   s := make([]int, 5)
		   s = append(s,1,2,3)
		   fmt.Println(s)
		}

- 输出：[0 0 0 0 0 1 2 3]
- make可用于初始化数组，第二个可选参数表示数组的长度。数组是不可变的。当执行make([]int, 5)时返回的是一个含义默认值(int的默认值为0)的数组:[0,0,0,0,0].而append函数便是在一个数组或slice后面追加新的元素，并返回一个新的数组或slice.这里append(s,1,2,3)是在数组s的继承上最佳三个新元素1，2，3，故返回的系数组为[0 0 0 0 0 1 2 3]

## 8.下面的代码有什么问题?

		type UserAges struct {
		   ages map[string]int
		   sync.Mutex
		}
		
		func (ua *UserAges) Add(name string, age int) {
		   ua.Lock()
		   defer ua.Unlock()
		   ua.ages[name] = age
		}
		func (ua *UserAges) Get(name string) int {
		   if age, ok := ua.ages[name]; ok {
		   	 return age
		   }
		   return -1
		}

- 在执行get方法时可能被panic
- 虽然有使用sync.Mutex做写锁，但是map是兵法读写不安全的。map属于引用类型，兵法读写时多个协程间是通过指针访问同一个地址，即访问共享变量，此时同时读写资源存在竞争关系。会报错信息:fatal error: concurrent map read and map write
- go1.9新版本中提供并发安全的map。首先需要了解两种锁：
1. sync.Mutex互斥锁
2. sync.RWMutex读写锁，基于互斥锁的实现，可以加多个度锁活着一个写锁

		//author: ysqi ,https://yushuangqi.com
		
		package main
		
		import (
			"fmt"
			"sync"
		)
		
		type UserAges struct {
			ages map[string]int
			sync.RWMutex
		}
		
		func (ua *UserAges) Add(name string, age int) {
			ua.Lock()
			defer ua.Unlock()
			ua.ages[name] = age
		}
		
		func (ua *UserAges) Get(name string) int {
			ua.RLock()
			defer ua.RUnlock()
			if age, ok := ua.ages[name]; ok {
				return age
			}
		
			return -1
		}
		
		func main() {
			count := 10000
			gw := sync.WaitGroup{}
			gw.Add(count * 3)
			u := UserAges{ages: map[string]int{}}
			add := func(i int) {
				u.Add(fmt.Sprintf("user_%d", i), i)
				gw.Done()
			}
			for i := 0; i < count; i++ {
				go add(i)
				go add(i)
			}
			for i := 0; i < count; i++ {
				go func(i int) {
					defer gw.Done()
					u.Get(fmt.Sprintf("user_%d", i))
					fmt.Print(".")
				}(i)
			}
			gw.Wait()
			fmt.Println("Done")
		}
		

## 9.下面的迭代会有什么问题？

		func (set *threadSafeSet) Iter() <-chan interface{} {
		   ch := make(chan interface{})
		   go func() {
		   	set.RLock()
		   	for elem := range set.s {
		   		ch <- elem
		   	}
		   	close(ch)
		   	set.RUnlock()
		   }()
		   return ch
		}

- 内部迭代出现阻塞。默认初始化时无缓冲区，需要等待接受者读取后才能继续写入。
- chan在使用make初始化时可附带一个可选参数来设置缓冲区。默认无缓冲，题目中便初始化的时无缓冲的chan，这样直邮写入的元素直到被读取后才能继续写入，不然就一致阻塞。设置缓冲区大小后，写入数据时可连续写入到缓冲区，知道缓冲区被占满。从chan中接收一次便可从缓冲区中释放一次。可以理解为chan是可以设置吞吐量的处理池

		ch := make(chan interface{})
		ch := make(chan interface{},1)是不一样的
		无缓冲的不仅仅时只能向ch通道放一个值而是一直要有人接收，那么ch <- elem才会继续下去，否则就一直阻塞。
		而缓冲为1则即使没有接受者也不会阻塞，因为缓冲大小是1只有当放第二个值的时候第一个还没被人拿走，这时候才会阻塞

## 10.以下代码能编译通过嘛？为什么？

		package main
		
		import (
			"fmt"
		)
		
		type People interface {
			Speak(string) string
		}
		type Stduent struct{}
		
		func (stu *Stduent) Speak(think string) (talk string) {
			if think == "bitch" {
				talk = "You are a good boy"
			} else {
				talk = "hi"
			}
			return
		}
		func main() {
			var peo People = Stduent{}
			think := "bitch"
			fmt.Println(peo.Speak(think))
		}
		

- 编译失败，值类型student{}未实现接口People的方法，不能定义为People类型
- func (stu *Student) Speak(think string) (talk string)是表示结构体类型*Student的指针有提供该方法，但该放啊并不属于结构类型Student的方法。
- 修改方案：
1. 定义为指针 var peo People = &Student
2. 方法定义在值类型上，指针类型本身是包含值类型的方法。func(stu Student) Speak(think string) (talk string)
