# 原因

- 有一定的复杂性
- 可以灵活调整项目的复杂性
- 平衡语言/爬虫之间的比重

# 爬虫分类

- 通用爬虫，如baidu，google
- 聚焦爬虫，从互联网获取结构化数据

# 总体结构

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/crawler.png)

# golang的库/框架

- henrylee2cn/pholcus
- colly
- gocrawl
- hu17889/go_spider

> 将不使用现成的库/框架

> 使用ElasticSearch作为数据存储

> 使用Go标准模板库实现http数据展示

# 爬虫主题

- 内容： 如新闻，博客，社区。。。
- 人：QQ空间，人人网，微博，微信，facebook？比较难
  - 相亲网站，求职网站？静态简单，切忌隐私数据

# 尝试人工获取内容

&emsp;&emsp;target: [珍爱网](http://www.zhenai.com/zhenghun)，因为该网站的隐私数据相对求职网站较少，并且目录很明显.对于翻页操作也很方便，直接反应在了url上，还有“猜你喜欢”可以连接跳转，但是需要处理闭环问题，重复的人比较多。每个用户主页的url都自带了有顺序的id，这也很有帮助

## 如何发现用户

- 通过城市列表->城市->(下一页)->用户
- 通过用户->猜你喜欢
- 通过已有用户id+1来猜测用户id

# 爬虫总体算法

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/city.png)

# 步骤

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/bingfa.png)

# 单任务版爬虫
- 目标： 获取并打印所有城市第一页用户的详细信息

		package main
		
		// 猜编码
		func determineEncoding(r io.Reader) encoding.Encoding {
		    bytes, err := bufio.NewReader(r).Peek(1024)
		    if err != nil {
		        panic(err)
		    }
		    e, _, _ := charset.DetermineEncoding(bytes, "")
		    return e
		}
		
		func main() {
		    resp, err := http.Get("http://www.zhenai.com/zhenghun")
		    if err != nil {
		        panic(err)
		    }
		    defer resp.Body.Close()
		
		    if resp.StatusCode != http.StatusOk {
		        fmt.Println("Error: status code", resp.StatusCode)
		        return
		    }
		
		    e := determineEncoding(resp.Body)
		    uft8Reader := transform.NewReader(resp.Body, e.NewDecoder())
		
		    all, err := ioutil.ReadAll(utf8Reader)
		    if err != nil {
		        panic(err)
		    }
		    // fmt.Printf("%s\n", all)发现城市乱码，发现编码是gbk格式, 需要transform  
		}
## 获取城市列表
&emsp;&emsp;我们采用正则表达式来获取信息的位置，该网页中数据都存放在css`$('#cityList>dd>a')`下。

### 正则表达式

		package main 
		
		const text = "My email is ccmouse@gmail.com"
		
		func main() {
		    re := regexp.MustCompile('[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+')
		    match := re.Findstring(text)
		    fmt.Println(match)
		}
		package main 
		
		const text = `
		My email is ccmouse@gmail.com
		email is abc@def.org
		email2 is    kkk@qq.com
		email3 is ddd@abc.com.cn
		`
		
		func main() {
		    re := regexp.MustCompile(`[a-zA-Z0-9]+@[a-zA-Z0-9.]+\.[a-zA-Z0-9]+`)
		    match := re.FindAllString(text, -1)
		    fmt.Println(match)
		    // output
		    // [ccmouse@gmail.com abc@def.org kkk@qq.com ddd@abc.com.cn]
		
		    re := regexp.MustCompile(`([a-zA-Z0-9]+)@([a-zA-Z0-9]+)\(.[a-zA-Z0-9]+)`)
		    match = re.FindAllStringSubmatch(text, -1)
		    for _, m := range match {
		        fmt.Println(m)
		    }
		    // output
		    // [ccmouse@gmail.com ccmouse gmail .com]
		    // [abc@def.org abc def .org]
		    // [kkk@qq.com kkk qq .com]
		    // [ddd@abc.com.cn ddd abc .com.cn]
		}

## 继续

		package main
		
		// 猜编码
		func determineEncoding(r io.Reader) encoding.Encoding {
		    bytes, err := bufio.NewReader(r).Peek(1024)
		    if err != nil {
		        panic(err)
		    }
		    e, _, _ := charset.DetermineEncoding(bytes, "")
		    return e
		}
		
		func printCityList(contents []byte) {
		    re := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun[0-9a-z]+)"[^>]*>([^<]+)</a>`)
		    matches := re.FindAllSubmatch(contents, -1) // return [][][]string
		    for _, m := range matches{
		        fmt.Printf("City: %s, URL: %s\n", m[2], m[1])
		    }
		    fmt.Printf("Matches found: %d\n", len(matches))
		}
		
		func main() {
		    resp, err := http.Get("http://www.zhenai.com/zhenghun")
		    if err != nil {
		        panic(err)
		    }
		    defer resp.Body.Close()
		
		    if resp.StatusCode != http.StatusOk {
		        fmt.Println("Error: status code", resp.StatusCode)
		        return
		    }
		
		    e := determineEncoding(resp.Body)
		    uft8Reader := transform.NewReader(resp.Body, e.NewDecoder())
		
		    all, err := ioutil.ReadAll(utf8Reader)
		    if err != nil {
		        panic(err)
		    }
		    // fmt.Printf("%s\n", all)发现城市乱码，发现编码是gbk格式, 需要transform  
		}

&emsp;&emsp;从上面的程序中已经获取到了城市列表。我们需要一些解析器来专门解析数据。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/parse.png)

## 抽象解析器

- 输入：utf-8编码的文本
- 输出：Request{URL, 对应Parser}列表， Item列表

## 抽象架构

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/engine.png)

&emsp;&emsp;engine驱动，有一些seed种子页面，对种子进行Request包装，在上面程序中就是城市列表的url和城市列表的parser，送给engine之后，加入到任务队列，取出任务之后，将任务的URL送给Fetcher，Fetcher就是一个从网络上获取数据的模块，Fetcher返回text文本给engine，engine给Parser，返回request和item列表，这些requests就是我们将来要做的任务，那么就把他们加入任务队列。目前来说item列表，只是打印操作。

&emsp;&emsp;上面代码中func main里面其实是Fetcher所做的事情，我们把它们摘出来。

		// fetcher.go
		package fetcher
		
		// 猜编码
		func determineEncoding(r io.Reader) encoding.Encoding {
		    bytes, err := bufio.NewReader(r).Peek(1024)
		    if err != nil {
		        log.Printf("fetcher error: %v", err)
		        return unicode.UTF8
		    }
		    e, _, _ := charset.DetermineEncoding(bytes, "")
		    return e
		}
		
		func Fetch(url string) ([]byte, error) {
		   resp, err := http.Get(url)
		    if err != nil {
		        panic(err)
		    }
		    defer resp.Body.Close()
		
		    if resp.StatusCode != http.StatusOk {
		        return nil, 
		            fmt.Errorf("wrong status code: %d", resp.StatusCode)
		    }
		
		    e := determineEncoding(resp.Body)
		    uft8Reader := transform.NewReader(resp.Body, e.NewDecoder())
		
		    return ioutil.ReadAll(utf8Reader)
		}

&emsp;&emsp;接下来隔离parser和engine

		// engine/types.go
		package engine
		
		type Request struct {
		    Url         string
		    ParserFunc  func([]byte) ParseResult
		}
		
		type ParseResult struct {
		    Requests    []Request
		    Items       []interface{}
		}
		
		func NilParser([]byte) ParseResult {
		    return ParseResult{}
		}

		// zhenai/parser/citylist.go
		package parser
		
		const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`
		
		func ParseCityList(contents []byte) engine.ParseResult{
		    re := regexp.MustCompile(cityListRe)
		        matches := re.FindAllSubmatch(contents, -1) // return [][][]string
		
		        // 每一个URL就是一个ParseResult
		        result := engine.ParseResult{}
		        for _, m := range matches{
		            // 城市名字
		            result.Items = append(result.Items, "City " + string(m[2]))
		            result.Requests = append(result.Requests, engine.Request{
		                Url: string(m[1]),
		                // 这里的Parser其实已经是第二层的城市解析器了
		                ParserFunc: engine.NilParser,
		            })
		        }
		    return result
		}

&emsp;&emsp;上面出现的NilParser,主要是起到一个临时不会运行出错的作用，因为我们还没有实现其它的Parser。

		// engine/engine.go
		package engine
		
		func Run(seeds ...Request) {
		    var requests []Request
		    for _, r := range seeds {
		        requests = append(requests, r)
		    }
		    for len(requests) > 0 {
		        r := requests[0]
		        requests = requests[1:]
		
		        log.Printf("fetching %s", r.Url)
		        body, err := fetcher.Fetch(r.Url)
		        if err != nil {
		            log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
		            continue
		        }
		
		        // 开始解析
		        parseResult := r.ParserFunc(body)
		        // 加入队列
		        resuests = append(requests, parseResult.Requests...) 
		
		        for _, item := range parseResult.Items {
		            log.Printf("got item %v", item)
		        }
		        
		    }
		}
		package main
		
		func main() {
		    engine.Run(engine.Request{
		        Url:        "http://www.zhenai.com/zhenghun",
		        ParserFunc: parser.ParseCityList,
		    })
		}
&emsp;&emsp;我们看到main现在已经很简洁了，Run启动后的流程和图片中的抽象架构图一致。养成一个好习惯，编写测试案例，我们来测试一下parsecitylist。

		// zhenai/parser/citylist_test.go
		package parser
		
		import (
		    "testing"
		)
		
		func TestParseCityList(t *testing.T) {
		    /*
		        contents, err := fetcher.Fetch(
		            "http://www.zhenai.com/zhenghun"
		        )
		        if err != nil {
		            panic(err)
		        }
		        fmt.Printf("%s\n", contens)
		        我们发现，打印下来的html文件头部少了一些符号。
		        猜测一下，是哪里除了问题？
		        回顾一下我们编写的猜测byte编码的程序，我们使用Reader peek了1024字节
		        应该就是这1024字节不见了。为什么呢？
		        bufio.NewReader(r).Peek(1024) 
		        bufio对r的1024字节读取后存了起来，下次读的时候从1025开始读，这就是bufio的作用。
		        既然如此，那我们就不传Reader对象了，直接传bufio.Reader，这样就算读了1024字节bufio自身不会变
		
		        func determineEncoding(r *bufio.Reader) encoding.Encoding {
		            bytes, err := r.Peek(1024)
		            ....
		        }
		
		        func Fetch(url string) ([]byte, error) {
		            ....
		            bodyReader := bufio.NewReader(resp.Body)
		            e := determineEncoding(bodyReader)
		            utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
		            return ioutil.ReadAll(utf8Reader)
		        }
		    */
		    // 将网站主页的内容提前存储在html文件中
		    contents, err := ioutil.ReadFile("citylist_test_data.html")
		    if err != nil {
		        panic(err)
		    }
		    
		    result := ParseCityList(contents)
		
		    const resultSize = 470
		    expectedUrls := []string{
		        "http://www.zhenai.com/zhenghun/aba",
		        "http://www.zhenai.com/zhenghun/akesu",
		        "http://www.zhenai.com/zhenghun/alashanmeng",
		    }
		    expectedCities := []string {
		        "City 阿坝", "City 阿克苏", "City 阿拉善盟",
		    }
		
		    for i, url := range expectedUirls {
		        if result.Requests[i].Url != url {
		            t.Errorf("expected url #%d: %s; but was %s", i, url, result.Requests[i].Url)
		        }
		    }
		    for i, city := range expectedCities {
		        if result.Items[i].(string) != city {
		            t.Errorf("expected city #%d: %s; but was %s", i, city, result.Items[i].(string))
		        }
		    }
		
		    if len(res.Requests) != resultSize {
		        t.Errorf("result should have %d requests; but had %d", resultSize, len(result.Requests))
		    }
		    
		    if len(res.Items) != resultSize {
		        t.Errorf("result should have %d items; but had %d", resultSize, len(result.Items))
		    }
		}

## 城市解析器

		// zhanai/parser/city.go
		
		package parser
		
		const cityRe = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`
		
		func ParseCity(contents []byte) engine.ParseResult {
		    re := regexp.MustCompile(cityRe)
		    matches := re.FindAllSubmatch(contents, -1)
		
		    result := engine.ParseResult{}
		    for _, m := range matches {
		        result.Items = append(result.Items, "User " + string(m2))
		        result.Requests = append(result.Requests, engine.Request{
		            Url: string(m[1]),
		            ParserFunc: engine.NilParser,
		        })
		    }
		    return result
		}
&emsp;&emsp;这时候就可以把ParserCityList中的NilParser修改为ParseCity了
		func ParseCityList(contents []byte) engine.ParseResult {
		    ...
		    result.Requests = append(result.Requests, engine.Request{
		        Url: string(m[1]),
		        parserFunc: ParseCity,
		    })
		}

## 用户解析器

&emsp;&emsp;用户自身的信息和择偶条件，有重复的字段，我们采用正则和css选择器来获取。

		// model/profile.go
		
		package model
		
		type Profile struct {
		    Name            string      // 姓名
		    Gender          string      // 性别
		    Age             int         // 年龄
		    Height          int         // 身高
		    Weight          int         // 体重
		    Income          string      // 收入
		    Marriage        string      // 婚姻状况
		    Education       string      // 教育
		    Occupation      string      // 职业
		    Hokou           string      // 户口
		    Xinzuo          string      // 星座
		    House           string      // 房子
		    Car             string      // 车子
		}

		// zhenai/parser/profile.go
		
		package parser
		
		// 由于正则比较多，不能每次调用func去编译，比较耗时，这里直接先编译
		var ageRe = regexp.MustCompile(`<td><span class="label">年龄: </span>([\d]+)岁</td>`)
		var heightRe = regexp.MustCompile(`<td><span class="label">身高: </span>(\d+)CM</td>`)
		var incomeRe = regexp.MustCompile(`<td><span class="label">月收入: </span>([^<]+)</td>`)
		var weightRe = regexp.MustCompile(`<td><span class="label">体重: </span><span field="">`)
		var marriageRe = regexp.MustCompile(`<td><span class="label">婚况: </span>([^<]+)</td>`)
		
		// 这里为什么要加name呢，是因为在parsecity中已经获取了name，就没必要再获取了
		// 这里加了name后，那上层的ParseCity的ParserFunc就需要函数式编程包装一下了
		func ParseProfile(contents []byte, name string) engine.ParseResult {
		    // 建立人的结构model
		    profile := model.Profile{}
		    profile.Name = name
		
		    if age, err := strconv.Atoi(extractString(contents, ageRe)); err == nil {
		        profile.Age = age
		    }
		    if height, err = strconv.Atoi(extractString(contents, heightRe)); err == nil {
		        profile.Height = height
		    }
		    if weight, err := strconv.Atoi(extractString(contents, weightRe)); err == nil {
		        profile.Weight = weight
		    }
		    profile.Income = extractString(contents, incomeRe)
		    profile.Gender = extractString(contents. genderRe)
		    profile.Car = extractString(contents, carRe)
		    profile.Education = extractString(contents, educationRe)
		    profile.Hukou = extractString(contents, hukouRe)
		    profile.House = extractString(contents, houseRe)
		    profile.Marriage = extractString(contents, marriageRe)
		    profile.Occupation = extractString(contents, occupationRe)
		    profile.Xinzuo = extractString(contents, xinzuoRe)
		
		    result := engine.ParseResult {
		        Items :[]interface{}{profile},
		    }
		
		    return result
		}
		
		func extractString(contents []byte, re *regexp.Regexp) string {
		    match := re.FindSubmatch(contents)
		
		    if len(match) >= 2 {
		        return string(match[1])
		    }
		    return ""
		}
		// zhenai/parser/city.go 修改parsecity
		
		result.Requests = append(result.Requests, engine.Request{
		            Url: string(m[1]),
		            ParserFunc: func(c []byte) engine.ParseResult {
		                return ParseProfile(c, string(m[2]))
		            },
		        })

# 并发版

&emsp;&emsp;在单任务版中速度特别慢，每个页面都要去请求然后等待它返回以后再请求下一个页面，这些动作完全是可以并发处理的。要把其中耗时最长等待实践最久的模块，当然是Fetcher模块，拿到URL后，到互联网上获得文本，这部分的时间是比较慢的，有远程的网络传输。Parser模块毕竟是对本地内存里文本的解析，这当然是快一些的。我们再看一步，这个Fetcher和Parser可以合成一个较大的模块，Fetcher的输出就是Parser的输入，完全可以把Fetcher和Parser合并起来，这个模块叫做Worker。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/crawl_worker.png)

&emsp;&emsp;这个Worker包含了Fetcher和Parser，还包含了一部分Engine。因为engine拿到了Fetcher的文本给了Parser。我们提取一下engine。

		// engine/engine.go
		
		func worker(r Request) (ParseResult, error) {
		    log.Printf("Fetching %s", r.Url)
		    body, err := fetcher.Fetch(r.Url)
		    if err != nil {
		        log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
		        return ParseResult{}, err
		    }
		    return r.ParserFunc(body), nil
		}
		
		func Run(seeds ...Request) {
			var requests []Request
			for _, r := range seeds {
				requests = append(requests, r)
			}
			for len(requests) > 0 {
				r := requests[0]
				requests = requests[1:]
		
				parsrResult, err := worker(r)
				if err != nil {
					continue
				}
		
				//加入队列
				requests = append(requests, parseResult.Requests...)
		
				for _, item := range parseResult.Items {
					log.Printf("got item %v", item)
				}
			}
		}
## 并发结构

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/crawl_worker_bingfa.png)

&emsp;&emsp;我们要并发worker，然后engine不动，加入了scheduler模块，因为要面临一个多对多并发目录的分配，有很多的Request，有很多的Worker等待要执行的Request，所以这个scheduler来分配这些任务。首先Worker的输入是Request输出是Requests，Items，上面的代码已经改好了。接下来engine接收之后，打印Items，把Request分给Scheduler。单任务结构图的箭头代表函数的输入输出参数，并发版的箭头就变成了channel.每个方框代表一个goroutine。那目前最核心的事情就是scheduler。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/crawl_worker_1.png)

&emsp;&emsp;scheduler最简单的做法：所有Worker公用一个channel，所有的Worker都是在同一时间抢同一个request，做完了继续抢下一个。

&emsp;&emsp;接下来我们把之前的engine修改为simpleengine，这样做的目的是将来实现的concurrentengine可以来回的切换。并且修改main.go

		$ mv engine.go simple.go
		
		// engine/simple.go
		
		type SimpleEngine struct{}
		
		func worker(r Request) (ParseResult, error) {
		    log.Printf("Fetching %s", r.Url)
		    body, err := fetcher.Fetch(r.Url)
		    if err != nil {
		        log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
		        return ParseResult{}, err
		    }
		    return r.ParserFunc(body), nil
		}
		
		func (e SimpleEngine)Run(seeds ...Request) {
			var requests []Request
			for _, r := range seeds {
				requests = append(requests, r)
			}
			for len(requests) > 0 {
				r := requests[0]
				requests = requests[1:]
		
				parsrResult, err := worker(r)
				if err != nil {
					continue
				}
		
				//加入队列
				requests = append(requests, parseResult.Requests...)
		
				for _, item := range parseResult.Items {
					log.Printf("got item %v", item)
				}
			}
		}
		
		package main
		
		func main(){
		    // 这里将来可以切换ConcurrentEngine
		    engine.SimpleEngine{}.Run(engine.Request{
		        Url:    "http://www.zhenai.com/zhenghun",
		        ParserFunc: parser.ParseCityList,
		    })
		}
&emsp;&emsp;现在开始撰写ConCurrentEngine
		// engine/concurrentEngine.go
		
		package engine
		
		// ConcurrentEngine本身就是一个scheduler，里面的参数给外面的scheduler设置
		type ConcurrentEngine struct {
		    Scheduler Schedler
		    WorkerCount int
		}
		
		// Scheduler 是一个接口，可以有很多实现
		type Scheduler interface {
		    // 提交给worker的操作
		    Submit(Request)
		    // 这就是work chan 送给worker的任务
		    ConfigureMasterWorkerChan(chan Request)
		}
		
		func (e *ConcurrentEngine) Run(seeds ...Request) {
		    // 所有的worker公用一个输入输出
		    in := make(chan Request)
		    out := make(chan ParseResult)
		    // 设置scheduler的输出，也就是worker的输入
		    e.Scheduler.ConfigureMasterWorkerChan(in)
		
		    // 开始建worker
		    for i := 0； i < e.WorkerCount; i++ {
		        createWorker(in, out)
		    }
		    // 创建完之后，再submit给worker chan
		    for _, r := range seeds {
		        e.Scheduler.Submit(r)
		    }
		
		    // 收out
		    for {
		        result := <- out
		        for _, item := range result.Items {
		            log.Printf("Got item: %v", item)
		        }
		        // 把下一批Request送给scheduler装载
		        for _, request := range result.Requests {
		            e.Scheduler.Submit(request)
		        }
		    }
		}
		
		func createWorker(in chan Request, out chan ParseResult) {
		    // 从work chan中抢request
		    // 执行完后还给work chan out
		    go func(){
		        for {
		            request := <- in
		            result, err := worker(request)
		            if err != nil {
		                continue
		            }
		            out <- result
		        }
		    }()
		}

		// scheduler/simple.go
		
		package scheduler
		
		type SimpleScheduler struct {
		    // 这其实就是ConfigureMasterChan
		    workerChan chan engine.Request
		}
		
		func (s *SimpleSchedluer) Submit(r engine.Request) {
		    // 把request送给work chan
		    s.workerChan <- r
		}
		
		func (s *SimpleScheduler) ConfigureMasterWorkerChan(c chan engine.Request) {
		    s.workerChan = c
		}

		package main
		
		func main() {
		    e := engine.ConCurrentEngine{
		        Scheduler: &engine.SimpleScheduler{},
		        WorkerCount: 10,
		    }
		    e.Run(engine.Request{
		        Url: "http://www.zhenai.com/zhenghun",
		        ParserFunc: ParseCityList,
		    })
		}

&emsp;&emsp;运行之后，发现程序卡在了10条数据之后。我们回顾一下
