
>转载自 https://www.flysnow.org

## 1. 包的命名
go语言的包的命名，遵循简洁、小写、和go文件所在目录同名的原则，这样便于引用，书写以及快速定位查找。
比如go自带的http包，它这个http目录下的所有go文件都属于http包，所以我们使用http包里面的函数、接口的时候，导入这个http包就可以了。

		package main
		import "net/http"
		func main() {
		    http.ListenAndServer("127.0.0.1:80", handler);
		}
从这个案例可以看到，我们导入的是`net/http`，这在go里叫做全路径，因为http包在net里面，net是最顶级的包，所以必须使用`全路径`导入，go编译程序才能找到http这个包，和我们文件系统的目录路径是一样的。
因为有了全路径，所以命名的包可以和其它库一样，只要它们的全路径不同就可以了，使用全路径的导入，也增加了包名命名的灵活性。
		package main
		import "github.com/rujews/tools"
## 2. main包
当把一个go文件的包名声明为main时，就等于告诉go编译器，我这个是一个可执行的程序，那么go编译程序就会尝试把它编译为一个二进制的可执行文件。
一个main包，一定会包含一个main函数。
在go语言中，同时要满足main包和包含main()函数，才会被编译成一个可执行文件。
我们来看一个Hello World的go语言版本：
		package main
		import "fmt"
		func main(){
		    fmt.Println("Hello World")
		}
假设该go文件叫hello.go放在`$GOPATH/src/hello`目录下，那么我们在这个目录下执行`go build `就会生成二进制文件，在windows系统下生成的是hello.exe，在unix，mac和linux下生成的是hello，我们在cmd或终端里执行它，就可以看待控制台打印的:

## 3. 导入包
要想使用一个包，必须先导入它才可以使用，Go语言提供了import关键字来导入一个包，这个关键字go编译器到磁盘的哪里去找想要导入的包，所以导入的包必须是一个全路径难过的包，也即是包所在的位置。
		import "fmt"
这就表示我们导入了fmt包，也就等于告诉go编译器，我们要使用这个包下面的代码。如果要导入多个包怎么办？go语言还为我们提供了导入块。
		import(
		    "net/http"
		    "fmt"
		)
使用一对括号包含的导入块，每个包独占一行。
现在我们导入了包，那么编译的时候，go编译器去什么位置找它们呢？这里就要介绍go的环境变量了。go有两个很重要的环境变量`GOROOT`和`GOPATH`，这是两个定义路径的环境变量，goroot是安装go的路径，比如/user/local/go；gopath是我们自定义的工作空间，比如/home/facedamon/go。
编译器会使用我们设置的这两个路径，再加上import导入的相对全路径来查找磁盘上的包，比如我们导入放入fmt包，编译器最终找到的是/usr/local/go/fmt这个位置
值得了解的是：对于包的查找是有优先级的，编译器会优先在goroot里搜索，其次是gopath，一旦找到，就会马上停止搜索。如果最终都没找到，就会编译异常
## 3. 远程包导入
		import "github.com/spf13/cobra"
这种导入，前提必须是该包托管在一个分布式的版本控制系统上，比如github、bitbucket等，并且是public的权限，可以让我们直接访问它们。
编译在导入它们的时候会先在gopath下搜索这个包，如果没有找到，就会使用`go get`工具从版本控制系统GitHub获取，并且会把获取到的源代码存储在gopath目录下对应的url的目录里，以供编译使用。
go get工具可以递归获取依赖包。
## 4. 命名导入
在使用import关键字导入包之后，我们就可以在代码中通过包名使用包下相应的函数、接口等。如果我们导入的包名正好有重复的怎么办？针对这种情况，go可以让我们对导入的包重新命名，这就是命名导入。
		package main
		import(
		    "fmt"
		    myfmt "mylib/fmt"
		)
## 5. init函数
每个包都可以有任意多个init函数，这些init函数都会在main函数之前执行。init函数通常来说用来做初始化变量、设置包活着其它需要在程序执行前的引导工作。比如使用_空标志符来导入一个包的目的就是想执行这个包的init函数。
我们以数据库驱动为例，go语言为了统一关于数据库的访问，使用`databases/sql`抽象了一层数据库操作，可以满足我们操作mysql、postgre等数据库，这样不管我们使用这些数据库的哪个驱动，编码操作都是一样的，想换驱动的时候，就直接换掉，而不是修改具体的代码。
这些数据库驱动的实现，就是具体的，可以由任何人实现，它的原理就是定义了init函数，在程序运行之前把实现好的驱动注册到sql包里，这样我们就使用它操作数据库了。
		package mysql
		import(
		    "database/sql"
		)
		func init(){
		    sql.Register("mysql", &MySQLDriver{})
		}
因为我们只想执行这个mysql包的init方法，并不想使用这个包，所以我们在导入这个包的时候，需要使用_重命名包名，避免编译报错
		import (
		    "database/sql"
		   _ "github.com/go-sql-driver/mysql"
		)
		db, err := sql.Open("mysql", "user:password@/dbname")
看起来非常简洁，剩下针对的数据库的操作，都是使用database/sql标准接口。
