# 接口
## duck typing

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/大黄鸭.png)

- 大黄鸭是不是一个鸭子?
    - 传统类型系统：脊索动物，脊椎动物亚门，鸟纲雁形目。。。:no
    - duck typing :yes，像鸭子走路，像鸭子叫，那么就是鸭子

> 描述事物的外部行为而非内部结构

&emsp;&emsp;不管你是如何组成的，有没有生命，它是不是一个鸭子，要从使用者的角度来看，小朋友就认为它是一个鸭子，而吃货就认为它不是一个鸭子。

- 严格说golang属于结构化类型系统，类似duck typing

1. python中的duck typing
    - 运行时才知道传入的retriever有没有get
    - 需要注释来说明接口

		def download(retriever):
		    return retriever.get("xxx.com")
		// retriever 就是duck typing 对象 download函数是一个使用者，retriever可以get。

2. C++中的 duck typing
    - 编译时才知道传入的retriever有没有get
    - 需要注释来说明

		template <class R>
		string download(const R& retriever) {
		    return retriever.get("xxx.com")
		}

3. Java中类似代码
    - 传入的参数必须实现Retriever接口
    - 不是duck typing

		<R extends Retriever>
		String download(R r) {
		    return r.get("xxx.com")
		}

> 同时需要实现Readable, Appendable接口怎么办？(可读可写)

&emsp;&emsp; 这个东西是可读可写的，但这件事Java做不了，不能同时是可读可写，不能同时实现俩接口，把所有的方法抄过来，调用者会全部更改

4. golang的 duck typing
    - 同时具有python, c++的duck typing的灵活性
    - 又具有java的类型检查

## 接口的定义
> 上述接口角色：使用者download, 实现者reteiever

- <font style='color:red'>golang接口由使用者定义></font>

		package main
		
		type Retriever interface {
		    Get(url string) string
		}
		
		func download(r Retriever) string {
		    return r.Get("xxx.com")
		}
		
		func main() {
		    var r Retriever
		    r = mock.Retriever{"this is fake imooc.com"}
		    r = real.Retriever{}
		    fmt.Println(download(r))
		}
		package mock
		
		type Retriever struct {
		    Contents string
		}
		
		func (r Retriever) Get(url string) string {
		    return r.Contents
		}
		package real
		
		type Retriever struct {
		    UserAgent string
		    TimeOut time.Duration
		}
		
		func(r Retriever) Get(url string) string {
		    resp, err := http.Get(url)
		    if err != nil {
		        panic(err)
		    }
		    result, err := httputil.DumpResponse(resp, true)
		
		    defer resp.Body.Close()
		
		    if err != nil {
		        panic(err)
		    }
		    return string(result)
		}
- 接口的实现是隐式的
- 接口变量自带指针
- 接口变量同样采用值传递，几乎不需要使用接口的指针
- 指针接收者实现只能以指针方式使用;值接收者都可

		    realRetriever := r.(real.Retriever)// 断言 Type assertion
		    fmt.Println(realRetriever.TimeOut)
- interface{} 表示任何类型

&emsp;&emsp;上一节的Queue队列，限定了元素必须为int,可以修改为interface{}

		package queue
		
		type Queue interface{}
		
		func (q *Queue) Push(v interface{}) {
		    *q = append(*q, v)
		}
		
		func (q *Queue) Pop() interface{} {
		    head := (*q)[0]
		    *q = (*q)[1:]
		    return head
		}
## 接口组合

&emsp;&emsp;既实现Retriever又实现Poster
		type Poster interface {
		    Post(url string, form map[string]string) string
		}
		
		func post( poster poster) {
		    poster.Post("http://www.imooc.com",
		        map[string]string{
		            "name": "ccmouse",
		            "course": "golang",
		        })
		}
		
		type RetrieverPoster interface {
		    Retriever
		    Poster
		}
		
		func session(poster RetrieverPoster) {
		    s.Get()
		    s.Post()
		}
