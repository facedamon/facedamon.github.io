
>转载自 https://www.flysnow.org

## 大小写
go的标志符，这个翻译感觉怪怪的，不过还是按这个起了标题，可以理解为go的变量、类型、字段等。这里的可见性，也就是说那些方法、函数】类型或者字段的可见性，比如哪些方法不想让另外一个包访问，我们就可以把它声明为非公开的；如果需要被另外一个包访问，就可以声明为公开的，和java语言里面的作用域类似
在go语言中，没有特别的关键字来声明一个方法、函数或者类型是否为公开的，go语言提供的是以大小写的方式进行区分的，如果一个类型的名字是以大写开头，那么其它包就可以访问；如果以小写开头，其它包就不能访问。
		package common
		type count int
		package main
		import (
		    "common"
		    "fmt"
		)
		
		func main(){
		    c := common.count(10)
		    fmt.Println(c)
		}
这是一个定义在common包里的类型count,因为它的名字以小写开头，所以我们不能在其它包里使用它，否则就会编译报错。
		./main.go 9: count refer to unexported name common.count
因为这个类型没有被导出，如果我们改为大写，就可以正常编译运行了，大家可以自己试试。
现在这个类型没有导出，不能使用，现在我们修改一下例子，增加一个函数，看看是否可行。
		package common
		type count int
		
		func New(v int) count {
		    return count(v)
		}
		
		func main(){
		    c := common.New(100)
		    fmt.Println(c)
		}
这里我们在common包里定义了一个导出的函数New，该函数返回一个count类型的值。New函数可以在其它包访问，但是count类型不可以，现在我们在main包李调用这个New函数，会发现是可以正常调用并且运行的，但是有个前提，必须使用:=这样的操作符才可以，因为它可以推断变量的类型。
这是一种非常好的能力，试想，我们在和其它人进行函数方法通信的时候，只需约定好接口，就可以了，至于内部实现，使用方看不到的，隐藏了实现。
		package common
		import "fmt"
		
		func NewLoginer() Loginer{
		    return defaultLogin(0)
		}
		
		type Loginer interface{
		    Login()
		}
		
		type defaultLogin int
		
		func(d defaultLogin) Login(){
		    fmt.Println("login in ...")
		}
		
		func main(){
		    l := common.NewLoginer()
		    l.Login()
		}
以上例子，我们对于函数间的通信，通过Loginer接口即可，在main函数中，使用者只需返回一个Loginer接口，至于这个接口的实现，使用者是不关心的，所以接口的设计者可以把defaultLogin类型设计为不可见，并让它实现接口Loginer，这样我们就隐藏了具体的实现。如果以后重构这个defaultLogin类型的具体实现时，也不会影响外部的使用者，极为方便，这也就是面向接口编程。
例如导出的结构体类型里，又一个未导出的字段，会出现怎样的问题。
		type User struct{
		    Name string
		    email string
		}
当我们在其它包声明和初始化user的时候，字段email是无法初始化的，因为它没有导出，无法访问。此外，一个导出类型，包含了一个未导出的方法也一样，也是无法访问的。
我们再扩展，导出和未导出的类型相互嵌入，会有什么样的发现呢？
		type user struct{
		    Name string
		}
		type Admin struct{
		    user
		}
被嵌入的user是未导出的，但是它的外部类型Admin是导出的，所以外部可以声明初始化Admin。
		func main(){
		    var ad common.Admin
		    ad.Name = "xx"
		    fmt.Println(ad)
		}
这里因为user是未导出的，所以我们不能再使用字面值直接初始化user了，所以只能先定义一个Admin类型的变量，再对Name字段初始化。这里Name可以访问是因为它是导出的，在user嵌入到Admin中时，它已经被提升为Admin的字段，所以它可以被访问。
如果我们还想使用:=操作符怎么做呢？
		ad := common.Admin{}
字面值初始化的时候什么都不做就好了，因为user未导出，所以我们不能直接使用字面值初始化字段。还有要注意的是，因为user未导出，所以我们不能通过外部类型访问内部类型了，也就是说ad.user这样的操作，都不会编译通过。
最后，我们做个总结，导出还是未导出，是通过名称首字母的大小写决定的，它们决定了是否可以访问，也就是标志符的可见性。
