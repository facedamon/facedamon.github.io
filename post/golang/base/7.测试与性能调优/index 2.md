
> Debugging Sucks and Testing Rocks

# 测试
## 传统测试 VS 表格驱动测试
### 传统测试

		@Test public void testAdd() {
		    assertEquals(3, add(1, 2));
		    assertEquals(Integer.MIN_VALUE, add(1, Integer.MAX_VALUE));
		}
-  测试数据和测试逻辑混在一起
-  出错信息不明确
-  一旦一个数据出错测试全部结束

### 表格驱动测试

		tests := []struct {
		    a, b, c int32
		}{
		    {1, 2, 3},
		    {0, 2, 2},
		    {0, 0, 0},
		    {-1, 1, 0},
		    {math.MaxInt32, 1, math.MinInt32},
		}
		
		for _, test := range tests {
		    if actual := add(test.a, test.b); actual != test.c {
		        // 出错信息
		    }
		}
- 分离的测试数据和测试逻辑
- 明确的出错信息
- 可以部分失败
- <font style="color:red">golang的语法使得我们更容易实践表格驱动测试</font>

## 测试实践

&emsp;&emsp;**filename: name_test.go**

&emsp;&emsp;**funcname: TestName**
		package triangle_test
		
		import (
		    "testing"
		    "triangle"
		)
		
		func TestTriangle(t *Testing.T) {
		    tests := []strust {a, b, c int32} {
		        {3, 4, 5},
		        {5, 12, 13},
		        {8. 15, 17},
		        {12, 35, 37},
		        {30000, 40000, 50000},
		    }
		
		    for _, tt := range tests {
		        if actual := calcTrangle(tt.a, tt.b); actual != tt.c {
		            t.Errorf("calcTriangle(%d, %d); got %d; expected %d", tt.a, tt.b, actual, tt.c)
		        }
		    }
		}
		package triangle
		
		func calcTriangle(a, b int) int {
		    var c int
		    c = int(math.Sqrt(float64(a*a + b*b)))
		    return c
		}

&emsp;&emsp;除了在Idea中进行测试，我们还可以在命令行中测试

		$ ls 
		triangle_test.go
		$ go test .
		ok learngo/basic    0.015s

## 测试代码覆盖率和性能测试
### 测试代码覆盖率
&emsp;&emsp;之前的章节中，最长不重复字串代码。
		func lengthOfNonRepeatingSubStr(s string) int {
		    lastOccurred := make(map[rune]int)
		    start := 0
		    maxLength := 0
		
		    for i, ch := range []rune(s) {
		        if lastI, ok := lastOccurred[ch]; ok && lastI >= start {
		            start = lastI + 1
		        }
		        if i - start + 1 > maxLength {
		            maxLength = i - start + 1
		        }
		        lastOccurred[ch] = i
		    }
		    return maxLength
		}
		func TestSubstr(t *testing.T) {
		    tests := []struct{
		        s string
		        ans int
		    }{
		        // Normal cases
		        {"abcabcbb", 3},
		        {"pwwkew", 3},
		
		        // Edge cases
		        {"", 0},
		        {"b", 1},
		        {"bbbbbbbb", 1},
		        {"abcabcabcd", 4},
		
		        // Chinese support
		        {"这里是慕课网", 6},
		        {"一二三一二三", 3},
		        {"黑化肥挥发发灰会花飞灰化肥挥发发黑会飞花", 8}
		    }
		
		    for _ , tt := range tests {
		        actial := lengthOfNonRepeatingSubStr(tt.s)
		        if actual != tt.ans {
		            t.Errorf("got %d for input %s; expected %d", actual, tt.s, tt.ans)
		        }
		    }
		}
- 命令行执行
  
		$ ls
		substr_test.go
		$ go test .
		ok learngo/basic    0.015s
		$ go test -converprofile=c.out  // 生成代码覆盖率文件
		coverage: 52.6% of statements
		ok learngo/basic    0.015s
		$ ls
		c.out substr_test.go
		$ go tool conver html=c.out     // 以html形式展示文件覆盖率描述文件c.out
### 性能测试
&emsp;&emsp;**funcname: BenchmarkName**
		func BenchmarkSubstr(b *testing.B) {
		    s := "黑化肥挥发发灰会花飞灰化肥挥发发黑会飞花"
		    ans := 8
		
		    for i := 0; i < b.N; i++{
		       actial := lengthOfNonRepeatingSubStr(s)
		        if actual != ans {
		            b.Errorf("got %d for input %s; expected %d", actual, s, ans)
		        } 
		    }
		}
		$ go test -bench .
		goos: darwin
		pkg: learngo/basic
		BenchmarkTriangle-4     1000000(一百万次)     1453 ns/op(每次operation需1453纳秒)
		PASS
		ok learngo/basi     1.476s

## 使用pprof进行性能调优
&emsp;&emsp;既然我们使用了Benakmark，那么我们看看怎样来提升他的性能，看看究竟慢在哪里，哪个部分消耗时间久

		$ go test -bench . -cpuprofile cpu.out
		goos: darwin
		goarch: amd64
		pkg: learngo/basic/substr
		BenchmarkSubstr-4       2000000     968 ns/op
		PASS
		ok   learngo/basic/substr   3.086s
		$ go tool pprof cpu.out
		Type: cpu
		Time: Mar 2, 2020 at 11:40am(CST)
		Duration: 3.08s, Total samples = 2.44s(79.25%)
		Entering interactive mode (type "help" for commands, "o" for options)
		(pprof)help
		(pprof)web
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/substr_svg1.png)
&emsp;&emsp;从上图可以看出消耗时候最多的有两部分，一是Map操作，而是decode操作，也就是对s的utf8解码的过程。
&emsp;&emsp;既然map消耗了实践，那我们试着优化一下，我们不使用map，改用slice。
		func lengthOfNonRepeatingSubStr(s string) int {
			// stores last occurred pos + 1
			// 0 means not seen
			lastOccurred := make([]int, 0xffff)
			start := 0
			maxLength := 0
		
			for i, ch := range []rune(s) {
				if lastI := lastOccurred[ch]; lastI > start {
					start = lastI
				}
				if i - start + 1 > maxLength {
					maxLength = i - start + 1
				}
				lastOccurred[ch] = i + 1
			}
			return maxLength
		}
		goos: darwin
		goarch: amd64
		pkg: github.com/facedamon
		BenchmarkSubstr-4          30000             44054 ns/op
		PASS
		ok      github.com/facedamon    1.948s
&emsp;&emsp;我们发现程序居然变慢了，我们来通过pprof查看一下。
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/substr_svg2.png)
&emsp;&emsp;可以看到出现了大量的GC垃圾回收，主要是因为lastOccurred数组被执行了多次，空间一直被回收，那么我们把这句话提出来。

		import "testing"
		
		var lastOccurred = make([]int, 0xffff)
		
		func lengthOfNonRepeatingSubStr(s string) int {
			// stores last occurred pos + 1
			// 0 means not seen
			for i := range lastOccurred {
				lastOccurred[i] = 0
			}
			start := 0
			maxLength := 0
		
			for i, ch := range []rune(s) {
				if lastI := lastOccurred[ch]; lastI > start {
					start = lastI
				}
				if i - start + 1 > maxLength {
					maxLength = i - start + 1
				}
				lastOccurred[ch] = i + 1
			}
			return maxLength
		}
		goos: darwin
		goarch: amd64
		pkg: github.com/facedamon
		BenchmarkSubstr-4   	  100000	     16400 ns/op
		PASS
&emsp;&emsp;我们发现快了很多，但是和我们第一次使用map时相比，还是慢了很多，我们使用pprof查看一下
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/substr_svg3.png)

&emsp;&emsp;我们发现花时间最多的是`memclrNoHeapPointers`,其实就是内存初始化lastOccurred这段，可能是数组太大。之所以大多数的时间花在了初始化上面，是因为字符串太短了，不管字符串是长是短，我们每次花的时间是固定的，只要字符串足够长，那我们的大头花费就不再初始化上面。既然这样，那我们把字符串变长。并且把准备输入数据的时间去掉。

		func BenchmarkSubstr(b *testing.B) {
			s := "黑化肥挥发发灰会花飞灰化肥挥发发黑会飞花"
			ans := 8
			for i := 0 ; i < 13; i++ {
				s = s + s
			}
			b.Logf("len(s) = %d", len(s))
		
			b.ResetTimer()
		
			for i := 0; i < b.N; i++{
				actual := lengthOfNonRepeatingSubStr(s)
				if actual != ans {
					b.Errorf("got %d for input %s; expected %d", actual, s, ans)
				}
			}
		}
		goos: darwin
		goarch: amd64
		pkg: github.com/facedamon
		BenchmarkSubstr-4   	     500	   2112526 ns/op
		--- BENCH: BenchmarkSubstr-4
		    substr_test.go:34: len(s) = 491520
		    substr_test.go:34: len(s) = 491520
		    substr_test.go:34: len(s) = 491520
		PASS
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/substr_svg4.png)
&emsp;&emsp;我们发现decode花费了很多时间，主要是因为字符很长，解码需要很久，但是没办法，程序本身就需要支持国际化

&emsp;&emsp;到这里程序优化的就差不多了，但是我们还是不知道最初的map版本和现在数组版本比在字符串很长的情况下到底哪个更快，那我们把之前的版本运行一下。

		//Map版本
		goos: darwin
		goarch: amd64
		pkg: github.com/facedamon
		BenchmarkSubstr-4   	     200	   7147757 ns/op
		--- BENCH: BenchmarkSubstr-4
		    substr_test.go:34: len(s) = 491520
		    substr_test.go:34: len(s) = 491520
		    substr_test.go:34: len(s) = 491520
		PASS
&emsp;&emsp;我们发现，在字符串很长的情况下，Map版本确实更慢，因为它要一直计算Hash

## 测试http服务器
&emsp;&emsp;源文件详见[http服务器错误统一处理](https://facedamon.github.io/post/golang/6.错误处理#服务器统一出错处理)
		package errwrapper_test
		
		func errPanic(writer http.ResponseWriter, request *http.Request) error {
		    panic(123)
		}
		
		type testingUserError string
		
		func (e testingUserError) Error() string {
		    return e.Message()
		}
		func (e testingUserError) Message() string {
		    return string(e)
		}
		func errUserError(writer http.ResponseWriter, request *http.Request) error {
		    return testingUserError("user error")
		}
		
		func errNotFound(writer http.ResponseWriter, request *http.Request) error {
		    return os.ErrNotExist
		}
		
		func errNoPermission(writer http.ResponseWriter, request *http.Request) error {
		    return os.ErrPermission
		}
		
		func errUnknown(writer http.ResponseWriter, request *http.Request) error {
		    return errors.New("unknown error")
		}
		
		func noError(writer http.ResponseWriter, request *http.Request) error {
		    fmt.Fprintln(writer, "no error")
		    return nil
		}
		
		var test = []struct {
		    h       appHandler
		    code    int
		    message string
		}{
		    {errPanic, 500, "Internal Server Error"},
		    {errUserError, 400, "user error"},
		    {errNotFound, 400, "Not Found"},
		    {errNoPermission, 403, "Forbidden"},
		    {errUnknown, 500, "Internal Server Error"},
		    {noError, 200, "no error "},
		}
		
		func TestErrWrapper(t *testing.T) {
		    
		
		    for _, tt : range tests {
		        f := errWrapper(tt.h)
		        response := httptest.NewRecoder()
		        request := httptest.NewRequest(
		            http.MethodGet,
		            "http://www.imooc.com", 
		            nil
		        )
		        f(response, request)
		
		        b, _ := ioutil.ReadAll(response.Body)
		        body := strings.Trim(string(b), "\n")
		        if response.Code != tt.code || body != tt.message {
		            t.Errorf("expected (%d, %s); got (%d, %s)", tt.code, tt.message, response.Code, body)
		        }
		    }
		}

## 开启http server服务测试

		func TestErrWrapperInServer(t *testing.T) {
		    for _, tt := range tests {
		        f := errWrapper(tt.h)
		        server := httptest.NewServer(http.HandleFunc(f))
		        resp, _ := http.Get(server.URL)
		 
		        b, _ := ioutil.ReadAll(resp.Body)
		        body := strings.Trim(string(b), "\n")
		        if resp.StatusCode != tt.code || body != tt.message {
		            t.Errorf("expected (%d, %s); got (%d, %s)", tt.code, tt.message, resp.StatusCode, body)
		        } 
		    }
		    
		}

## 生成文档

		package queue
		
		// An FIFO queue
		type Queue []int
		
		// Pushes the element into the queue
		//  e.g q.Push(123)
		func (q *Queue) Push(v int) {
		    *q = append(*q, v)
		}
		
		// Pops element from head
		func (q *Queue) Pop() int {
		    head := (*q)[0]
		    *q = (*q)[1:]
		    return head
		}
		
		// Returns if the queue is empty or not
		func (q *Queue) IsEmpty() bool {
		    return len(*q) == 0
		}

> $ godoc -http :6060  //在网页上查看doc

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/doc.png)

&emsp;&emsp;golang也提供了写示例代码的能力；

		package queue
		
		func ExampleQueue_Pop() {
		    q := Queue{1}
		    q.Push(2)
		    q.Push(3)
		    fmt.Println(q.Pop())
		    fmt.Println(q.Pop())
		    fmt.Println(q.IsEmpty)
		
		    fmt.Println(q.Pop())
		    fmt.Println(q.IsEmpty) 
		
		    // Output:
		    // 1
		    // 2
		    // false
		    // 3
		    // true
		}
&emsp;&emsp;在Idea下是可以直接执行的，并能检测output是否正确

> $ godoc -http :6060

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/docoutput.png)

