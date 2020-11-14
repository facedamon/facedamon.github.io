
# 什么是值传递
&emsp;&emsp;值传递的意思是：函数传递总是原来这个东西的一个副本，一幅拷贝。比如我们传递一个int类型的参数，传递的其实是这个参数的一个副本；传递一个指针类型的参数，其实传递的是这个指针的一份拷贝，而不是这个指针指向的值。
对于int类型这类基础类型我们可以很好的理解，它就是一个拷贝，但是指针呢？我们觉得可以通过它修改原来的值，怎么会是一个拷贝呢？下面我们看一个例子。
		func main(){
		    i := 10
		    ip := &i
		    fmt.Printf("原始指针的内存地址是：%p\n", &ip)
		    modify(ip)
		    fmt.Println("int值被修改了，新值为：",i)
		}
		func modify(ip *int){
		    fmt.Printf("函数里接收到的指针的内存地址是: %p\n", &ip)
		    *ip = 1
		}
我们运行，可以看到输入结果如下：
		原始指针的内存地址是： 0xc42000c028
		函数里接收到的指针的内存地址是：0xc42000c038
		int值被修改了，新值为：1
&emsp;&emsp;首先我们要知道，任何存放在内存里的东西都有自己的地址，指针也不例外，它虽然指向别的数据，但是也存放该指针的内存。
所以通过输出我们可以看到，这是一个指针的拷贝，因为存放这两个指针的内存地址是不同的，虽然指针的值相同，但是是两个不同的指针。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/other/golang是值传递.png)

&emsp;&emsp;通过上面的图，可以更好的理解。首先我们可以看到，我们声明了一个变量i，值为10，它的内存存放地址是0xc420018070,通过这个内存地址，我们可以找到变量i，这个内存地址也就是变量i的指针ip。
指针ip也是一个指针类型的变量，它也需要内存存放它，它的内存地址是多少呢？是0xc42000c028。在我们传递指针变量ip给modify函数的适合，是该指针变量的拷贝，所以新拷贝的指针变量ip，它的内存已经变了，是新的0xc42000c038。
不管是0xc42000c028还是0xc42000c038，我们都可以称之为指针的指针，它们指向同一个指针0xc420018070,这个0xc420018070又指向变量i，这也就是为什么我们可以修改变量i的值。
# 迷惑Map
&emsp;&emsp;了解清楚了传值和传引用，但是对于Map类型来说，可能觉得还是迷惑，一来我们可以通过方法修改它的内容，而来它没有明显的指针
		func main(){
		    persons := make(map[string]int)
		    persons["张三"] = 19
		    mp := &persons
		    fmt.Printf("原始map的内存地址是：%p\n", mp)
		    modify(persons)
		    fmt.Println("map值被修改了，新值为：",persons)
		}
		
		func modify(p map[string]int){
		    fmt.Printf("函数里接收到map的内存地址是：%p\n", &p)
		    p["张三"] = 20
		}
运行打印输出：
		原始map的内存地址是：0xc42000c028
		函数里接收到map的内存地址是: 0xc42000c038
		map值被修改了，新值为：map[张三：20]
两个内存地址是不一样的，所以这又是一个值传递，那么为什么我们可以修改Map的内容呢？先不急，我们先看一个自己实现的struct。
		type Person struct{
		    Name string
		}
		func modify(p Person){
		    fmt.Printf("函数里接收到Person的内存地址是：%p\n", &p)
		    p.Name = "李四"
		}
		
		func main(){
		    p := Person{"张三"}
		    fmt.Printf("原始Person的内存地址是：%p\n", &p)
		    modify(p)
		    fmt.Println(p)
		}
运行打印输出：
		原始Person的内存地址是： 0xc4200721b0
		函数里接收到Person的内存地址是：0xc4200721c0
		{张三}
我们发现，我们自己定义的Person类型，在函数传参的时候也是值传递，但是它的值并没有被修改。
这也即是说，map类型和我们自己定义的Person struct类型是不一样的我们尝试把modify函数的接收参数该为Person的指针。
		func main(){
		    p := Person{"张三"}
		    modify(&p)
		    fmt.Println(p)
		}
		func modify(p *Person)
		    p.Name = "李四"
&emsp;&emsp;在运行查询输出，我们发现，这次被修改了。我们这里省略了内存地址的打印，因为我们上面int类型的例子已经证明了指针类型的参数也是值传递的。指针类型可以修改，非指针类型不行，那么我们可以大胆的猜测，我们使用make函数创建的map是不是一个指针类型呢？看一下源代码：
		// makemap implements a Go map creation make(map[k]v, hint)
		// if the compiler has determined that the map or the first bucket
		// can be created on the stack, h and/or bucket may be non-nil.
		// if h != nil , the map can be created directly in h.
		// if bucket != nil , bucket can be used as the first bucket.
		func makemap(t *maptype, hint int64, h *hmap, bucket unsafe.Pointer) *hmap{}
&emsp;&emsp;通过chakansrc/runtime/hashmap.go源码发现，的确和我们猜测的一样，make函数返回的是一个hmap类型的指针*hmap。也就是说map === *hmap。现在看func modify(p map)这样的函数，其实就等于func modify(p *hmap)，和前面第一节什么是值传递里举的func modify(ip *int)的例子一样，可以参考分析。
所以在这里，Go语言通过make函数，字面量的包装，为我们省去了指针的操作，让我们可以更容易的使用map。这里的map可以理解为引用类型，但是记住引用类型不是传引用。
# chan类型
chan类型本质上和map类型是一致的，这里不做过多的介绍，参考以下源代码：
		func makechan(t *chantype, size int64) *hchan{}
# 和map、chan都不一样的slice
slice和map、chan都不太一样，它也可以在函数中私改对应的内容。
		func main(){
		    ages := []int{6,6,6}
		    fmt.Printf("原始slice的内存地址是%p\n", ages)
		    modify(ages)
		    fmt.Println(ages)
		}
		func modify(ages []int){
		    fmt.Printf("函数里接收到slice的内存地址是%p\n", ages)
		    ages[0] = 1
		}
&emsp;&emsp;运行打印结果，发现的确是被修改了，而且我们这里打印slice的内存地址是可以直接通过%p打印的，不用使用&取地址符转换。
这就可以证明make的slice也是一个指针了嘛？不一定，也可能是fmt.Printf把slice特殊处理了。
		func (p *pp)fmtPointer(value reflect.Value, verb rune){
		    var u uintptr
		    switch value.Kind(){
		        case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer
		            u = value.Pointer()
		        default:
		            p.badVerb(verb)
		            return
		    }
		}
通过源代码发现，对于chan、map、slice等被当成指针处理，通过value.Pointer()获取对应的值的指针。
		// if v`s kind is slice ,the returned pointer is to the first
		// element of the slice, if the slice is nil the returned value is 0
		// if the slice is empty but non-nil the return value is non-zero
		func (v Value) Pointer() uintptr{
		    k := v.Kind()
		    switch k{
		        case Slice: return (*SliceHeader)(v.ptr).Data
		    }
		}
很明显了，当是slice类型的时候，返回是slice这个结构体里，字段Data第一个元素的地址。
		type SliceHeader struct{
		    Data uintptr
		    len int
		    Cap
		}
		type slice struct{
		    array unsafe.Pointer
		    len int
		    cap int
		}
&emsp;&emsp;所以我们通过%p打印的slice变量ages的地址其实就是内部存储数组元素的地址，slice是一种结构体+元素指针的混合类型，通过元素array的指针，可以达到修改slice里存储元素的目的。
所以修改类型的内容的办法有很多种，类型本身作为指针可以，类型里有指针类型的字段也可以
单纯的从slice这个结构体看，我们可以通过modify修改存储元素的内容
