
>转载自 https://www.flysnow.org

## 什么是基准测试
基准测试，是一种测试代码性能的方法，比如你有多种不同的方案，都可以解决问题，那么到底那种方案性能更好呢？这个时候基准测试就派上用场了。

基准测试主要是通过测试CPU和内存的效率问题，来评估被测试代码的性能，进而找到更好的解决方案。比如连接池的数量不是越多越好，那么哪个值才是最优值呢，这就需要配合基准测试不断调优了。

## 如何编写基准测试
基准测试代码的编写和单元测试非常相似，它也有一定的规则，我们先看一个示例。
		itoa_test.go
		func BenchmarkSprintf(b *testing.B){
		    num := 10
		    b.ResetTimer()
		    for i:=0; i<b.N;i++{
		        fmt.Sprintf("%d", num)
		    }
		}
这是一个基准测试的例子，从中我们可以看出以下规则：

1. 基准测试的代码文件必须以_test.go结尾
2. 基准测试的函数必须以Benchmark开头，必须是可导出的
3. 基准测试函数必须接受一个指向Benchmark类型的指针作为唯一参数
4. 基准测试函数不能有返回值
5. b.ResetTimer是重置计时器，这样可以避免for循环之前的初始化代码的干扰
6. 最后的for循环很重要，被测试的代码放在循环里
7. b.N是基准测试框架提供的，表示循环的次数，因为需要反复调用测试的代码，才可以评估性能

下面我们运行基准测试，看看效果。
		go test -bench=. -run=none
		BenchmarkSprintf-8      20000000               117 ns/op
		PASS
		ok      flysnow.org/hello       2.474s

运行基准测试也要使用go test命令，不过我们要加上-bench=标记，它接受一个表达式参数，匹配基准测试的函数，`.`表示运行所有基准测试。

因为默认情况下go test会运行单元测试，为了防止单元测试的输出影响我们查看基准测试的结果，可以使用-run匹配一个从来没有的单元测试方法，过滤单元测试的输出，我们这里使用none，因为我们基本上不会创建这个名字的单元测试方法。

下面着重解释下输出结果，看到函数后面的-8了嘛，这个表示运行时对应的GOMAXPROCS的值。接着的20000000表示for循环的次数，也就是调用被测试代码的次数，最后的117 ns/op表示每次需要花费117纳秒。

以上是测试时间默认是1秒，也就是1秒的时间，调用两千万次，每次调用花费117纳秒。如果想让测试运行的时间更长，可以通过-benchtime指定，比如3秒

		go test -bench=. -benchtime=3s -run=none
		BenchmarkSprintf-8      50000000               109 ns/op
		PASS
		ok      flysnow.org/hello       5.628s
可以发现，我们加长了测试时间，测试的次数变多了，但是最终的性能结果：每次执行的时间，并没有太大变化。一般来说这个值最好不要超过3秒，意义不大。
## 性能对比
上面哪个基准测试的例子，其实是一个int类型转为string类型的例子，标准库里还有集中方法，我们看下哪种性能更加。
		func BenchmarkSprintf(b *testing.B){
			num:=10
			b.ResetTimer()
			for i:=0;i<b.N;i++{
				fmt.Sprintf("%d",num)
			}
		}
		
		func BenchmarkFormat(b *testing.B){
			num:=int64(10)
			b.ResetTimer()
			for i:=0;i<b.N;i++{
				strconv.FormatInt(num,10)
			}
		}
		
		func BenchmarkItoa(b *testing.B){
			num:=10
			b.ResetTimer()
			for i:=0;i<b.N;i++{
				strconv.Itoa(num)
			}
		}
运行基准测试，看看结果
		go test -bench=. -run=none              
		BenchmarkSprintf-8      20000000               117 ns/op
		BenchmarkFormat-8       50000000                33.3 ns/op
		BenchmarkItoa-8         50000000                34.9 ns/op
		PASS
		ok      flysnow.org/hello       5.951s
从结果上卡strconv.FormatInt函数最快，其次是strconv.Itoa，最后才是fmt.Sprintf，前两个函数性能达到了最后一个3倍多。那么最后一个函数为什么这么慢，我们再通过-benchmem找到根本原因。
		o test -bench=. -benchmem -run=none
		BenchmarkSprintf-8      20000000               110 ns/op              16 B/op          2 allocs/op
		BenchmarkFormat-8       50000000                31.0 ns/op             2 B/op          1 allocs/op
		BenchmarkItoa-8         50000000                33.1 ns/op             2 B/op          1 allocs/op
		PASS
		ok      flysnow.org/hello       5.610s
-benchmem可以提供每次操作分配内存的次数，以及每次操作分配的字节数。从结果上看，性能高的两个函数，每次操作都是进行1次内存分配，而最慢的那个要分配2次；性能高的每次操作分配2个字节内存，而满的那个函数每次需要分配16个字节内存。
从这个数据我们知道它为什么这么慢了，内存分配都占用都太高。

在代码开发中，对于我们要求性能的地方，编写基准测试非常重要，这有助于我们开发性能更好的代码。不过性能、可用性、复用性等也要有一个相对的取舍，不能为了追求性能而过度优化。
