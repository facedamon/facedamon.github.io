
>转载自 https://www.flysnow.org

## 什么是单元测试
单元测试一般是用来测试我们的代码逻辑有没有问题，有没有按照我们期望的运行，以保证代码质量。

大多数的单元测试，都是对某一个函数方法进行测试，以尽可能的保证没有问题或者问题可被我们预知。为了达到这个目的，我们可以使用各种手段、逻辑、模拟不同的场景进行测试。

这里我们在package main里定义一个函数Add，求两个数之和的函数，然后我们使用单元测是进行求和逻辑测试。
		main.go
		func Add(a, b int) int {
		    return a+b
		}
		main_test.go
		func TestAdd(t *testing.T){
		    sum := Add(1,2)
		    if sum == 3 {
		        t.Log("the result is ok")
		    } else {
		        t.Fatal("the result is wrong")
		    }
		}
然后我们在中断下运行go test -v 就可以看到测试结果了。

Go语言为我们提供了测试框架，以便帮助我们更容易的进行单元测试，但是要使用这个框架，需要遵循以下几点规则：

1. 含有单元测试代码的go文件必须以`_test.go`结尾，Go语言测试工具只认符合这个规则的文件。
2. 单元测试文件名_test.go前面的部分最好是被测试的昂发所在go文件的文件名，比如例子中main_test.go，因为测试的Add函数，在main.go文件中。
3. 单元测试的函数名比如以`Test`开头，是可导出公开的函数。
4. 测试函数的签名必须接收一个指向testing.T类型的指针，并且不能返回任何值。
5. 函数名最好是Test+要测试测方法函数名，比如例子中是TestAdd，表示测试的是Add这个函数。

遵循以上规则，我们就可以很容易的编写单元测试了，单元测试的重点在于测试代码的逻辑，场景等，以便尽可能的测试全面，保证代码质量逻辑。

## 表组测试
还有一种单元测试方法叫表组测试，这个和基本的单元测试非常相似，只不过它是有好几个不同输入以及输出组成的一组单元测试。

比如上个例子中，我们测试了1+2，如果我们再加上3+4，9+2等，这就有了好几个输入，同时对应也有好几个输出，这种一次性测试很多输入输出场景的测试，就是表组测试。

		func TestAdd(t *testing.T) {
			sum := Add(1,2)
			if sum == 3 {
				t.Log("the result is ok")
			} else {
				t.Fatal("the result is wrong")
			}
			
			sum=Add(3,4)
			if sum == 7 {
				t.Log("the result is ok")
			} else {
				t.Fatal("the result is wrong")
			}
		}
## 模拟调用
单元测试的原则，就是你所测试的函数方法，不要受到依赖环境的影响，比如网络访问，因为有时候我们运行单元测试的时候，并没有联网，那么总不能让单元测试因为这个失败吧？所以这个时候模拟网络访问就有必要了。

针对模拟网络访问，标准库提供了一个`httptest`包，可以让我们模拟http的网络调用，下面举个例子了解使用。
首先我们创建一个处理http请求的函数，并注册路由。
		package common
		import (
		    "net/http"
		    "encoding/json"
		)
		func Routes(){
		    http.HandleFunc("/sendjson", SendJSON)
		}
		
		func SendJSON(rw http.ResponseWriter, r *http.Request) {
		    u := struct {Name string}{Name: "张三"}
		    rw.Header().Set("Content-Type", "application/json")
		    rw.WriteHeader(http.StatusOK)
		    json.NewEncoder(rw).Encode(u)
		}
非常简单，这里是一个/sendjson API，当我们访问这个API时，会返回一个JSON字符串。现在我们对这个API服务进行测试，但是我们又不能时时刻刻都启动着服务，所以这里就用到了外部终端对API的网络访问请求。
		func TestSendJSON(t *testing.T) {
		    req, err := http.NewRequest(http.MethodGet, "/sendjson", nil)
		    if nil != err {
		        t.Fatal("创建Request失败")
		    }
		    rw := httptest.NewRecorder()
		    http.DefaultServeMux.ServeHTTP(rw, req)
		    log.Println("code:", rw.Code)
		    log.Println("body:", rw.Body.String())
		}
运行这个单元测试，就可以看到我们访问/sendjson API的结果，并且我们没有启动任何HTTP服务就达到了目的。这主要利用httptest.NewRecorder创建一个http.ResponseWriter,模拟真实服务端的响应，这种响应是通过调用http.DefaultServeMux.ServeHTTP方法触发的。

还有一种模拟调用的昂视，是真的在测试机上模拟一个服务器，然后进行调用测试。
		func mockServer() *httptest.Server {
			//API调用处理函数
			sendJson := func(rw http.ResponseWriter, r *http.Request) {
				u := struct {
					Name string
				}{
					Name: "张三",
				}
		
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusOK)
				json.NewEncoder(rw).Encode(u)
			}
			//适配器转换
			return httptest.NewServer(http.HandlerFunc(sendJson))
		}
		
		func TestSendJSON(t *testing.T) {
			//创建一个模拟的服务器
			server := mockServer()
			defer server.Close()
			//Get请求发往模拟服务器的地址
			resq, err := http.Get(server.URL)
			if err != nil {
				t.Fatal("创建Get失败")
			}
			defer resq.Body.Close()
		
			log.Println("code:", resq.StatusCode)
			json, err := ioutil.ReadAll(resq.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("body:%s\n", json)
		}
模拟服务器的创建使用的是httptest.NewServer函数，它接收一个http.handler处理API请求的接口。代码示例中使用了Handler的适配器模式，http.handlerfunc是一个函数类型，实现了http.handler接口，这里是强制类型转换，不是函数的调用，关于适配器的更多信息，可以参考另一篇文章:接口型函数。

## 测试覆盖率
我们尽可能的模拟更多的场景来测试我们代码的不同情况，但是有时候的确也又忘记测试的代码，这时候我们就需要测试覆盖率作为参考了。

由单元测试的代码，触发运行到的被测试代码的代码行数占所有代码行数的比例，被称为测试覆盖率，代码覆盖率不一定完全精准，但是可以作为参考，可以帮我们测量和我们预计的覆盖率之间的差距，go test工具，就为我们提供了这么一个度量测试覆盖率的能力。
		main.go
		func Tag(tag int){
			switch tag {
			case 1:
				fmt.Println("Android")
			case 2:
				fmt.Println("Go")
			case 3:
				fmt.Println("Java")
			default:
				fmt.Println("C")
		
			}
		}
		main_test.go
		func TestTag(t *testing.T) {
		    Tag(1)
		    Tag(2)
		}
现在我们使用go test工具运行单元测试，和前几次不一样的是，我们要显示测试覆盖率，所以要多加一个参数`- coverprofile`，所以完整的命令为：`go test -v -coverprofile=c.out`，-coverprofile是指定生成的覆盖率文件，例子中是c.out，这个文件一会我们会用到。现在我们看终端输出，已经有了一个覆盖率。
		=== RUN   TestTag
		Android
		Go
		--- PASS: TestTag (0.00s)
		PASS
		coverage: 60.0% of statements
		ok      flysnow.org/hello       0.005s
coverage: 60.0% of statements, 60%的测试覆盖率，还没有达到100%，那么我们看看还有哪些代码没有被测试到。这就需要我们刚刚生成的测试覆盖率文件c.out生成覆盖率报告了。生成报告有go为我们提供的工具，使用`go tool cover -html=c.out -o=tag.html`即可生成一个名字为tag.html的html格式的测试覆盖率报告，这里有详细的信息告诉我们哪一行代码测试到了，哪一行代码没有测试到。
