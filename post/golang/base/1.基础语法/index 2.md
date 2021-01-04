
# 变量定义

&emsp;&emsp;在golang中变量名在前变量类型在后

&emsp;&emsp;这样设计的初衷是什么?我们回想一下使用其它语言的变量的时候是先想到变量的名字还是先想到变量类型

&emsp;&emsp;声明之后默认是空值,是一个具体的值

		
		package main
		
		import "fmt"
		
		func variableZeroValue() {
		    var a int
		    var s string
		    // 0 ""
		    fmt.Printf("%d %q\n", a, s)
		}
		
		func main() {
		    fmt.Println("Hello world")
		}
		

&emsp;&emsp;变量当然也可以赋初值，也可以同行赋值多个变量

&emsp;&emsp;变量一旦定义一定要用到，否则就是语法错误(可能在刚开始学习的时候，觉得有点束缚手脚，但是当我们在做一个项目的时候，这是一个很好的习惯，避免了我们定义一些无用的变量)



		
		func variableInitialValue() {
		    var a， b int = 3, 4
		    var s string = "abc"
		    // 3 4 abc
		    fmt.Println(a, b, s)
		}
		

&emsp;&emsp;我们再思考一下，上面的变量声明数据类型，好像也是多余的，因为现在的编译器完全可以从后面的值反向推断出数据类型。在早期的C和C++是不能这样做的，但是后来的C++有了auto关键字，因此golang它完全是可以做到的。

		
		func variableTypeDeduction() {
		    var a, b, c ,s = 3, 4, true, "def"
		    // 3, 4, true, def
		    fmt.Println(a, b, c, s)
		}
		

&emsp;&emsp;上面的声明样式，看上去比较繁琐。所以有一种更加简单的写法。我们省略var关键字，使用":="进行赋值

		
		func variableShorter() {
		    a, b, c, s := 3, 4, true, "def"
		    b = 5
		    // 3, 5, true, def
		    fmt.Println(a, b, c, s)
		}
		

&emsp;&emsp;上面的变量定义示例，它们的作用域在函数里面。那在函数外面也能定义变量，此时就不能使用":="格式进行声明赋值了，必须使用var关键字。

&emsp;&emsp;**它并不是一个全局变量**，它是一个包内部变量，它没有全局变量的说法

		package main
		
		import "fmt"
		
		var aa = 3
		var ss = "kkk"
		bb := true // error
		
		... 
		

&emsp;&emsp;上面的一行一行的声明比较啰嗦，可以这样

		package main
		
		import "fmt"
		
		var (
		    aa = 3
		    ss = "kkk"
		    bb = true
		)
		
		...

## 回顾一下
-  使用var关键字
    - var a, b, c bool
    - var s1, s2 string = "hello world"
    - 可放在函数内， 或直接放在包内
    - 使用var()集中定义变量
- 编译器自动决定类型
    - var a, b, i, s1, s2 = true, false, 3, "hello", "world"
- 使用:=定义变量
    - a, b, i, s1, s2 := true, false, 3, "hello", "world"
    - 在函数内使用

> 函数内 能尽量使用:=就使用:=


# 内建变量类型

- bool, string
- (u)int, (u)int8, (u)int16, (u)int32, (u)int64, (**无符号整数和有符号整数**，有符号整数分为两类：1. 规定长度的，比如int8 int16 int32 int64, 2. 不规定长度的，那它的长度是跟着操作系统来的，在32位系统下是32位的，在64位系统下是64位的。 Golang没有诸如"long/ longlong 这些东西")uintptr (指针长度也是跟着操作系统来的)
- byte, rune(**字符型**。字面意思有点像char，但char是1字节，现在是全球化应用，要应对多个语言，rune的长度是32位，rune is an alias for int32 and is equivalent to int32 in all ways. byte是8位， byte is an alias for uint8 and is equivalent to uint8 in all ways)
- float32, float64, complex64, complex128(浮点数类型)

## 强制类型转换
- **golang类型转换都是强制的**

		
		func triangle() {
		    var a, b int = 3, 4
		    var c int
		    c = int(math.Sqrt(float64(a * a + b * b)))
		
		}
		

# 常量与枚举

&emsp;&emsp;使用关键字const,可以省略常量类型，声明时不确定具体是什么类型。常量只是**替换动作**

&emsp;&emsp;golang的常量并不需要全大写，可能在其它语言中会全大写，但是在golang中，首字母大写代表public,因此在这里不需要全部大写。

		func consts() {
		    //const filename = "abc.txt"
		    //const a, b = 3, 4
		    const (
		        filename = "abc.txt"
		        a, b = 3, 4
		    )
		    var c int
		    // 因为a b 是常量，在这里只是替换动作，所以替换后是float64类型
		    c = int(math.Sqrt(a * a + b * b))
		}

&emsp;&emsp;枚举类型也是采用const关键字，但是每一个常量必须被赋值。

		func enums() {
		    const (
		        cpp = 0
		        java = 1
		        python = 2
		        golang = 3
		    )
		    fmt.Println(cpp, java, python, golang)
		}

&emsp;&emsp;golang为了满足这种需求，做了一个简化。**iota表示这组const是自增的**。iota可以作为一个计算种子，只要你可以写出枚举的公式，就可以实现枚举自增

		func enums() {
		    const (
		        cpp iota
		        _ 
		        python 
		        golang 
		        javascript
		    )
		    fmt.Println(cpp,  python, golang,javascript)
		
		    //b, kb, mb, gb, tb, pb
		    const (
		        b = 1 << (10 * iota)    // 1 << 0 = 1
		        kb                      // 1 << 10 = 2^10 = 1024
		        mb
		        gb
		        tb
		        pb
		
		    )
		    fmt.Println(b, kb, mb, gb, tb, pb)
		}

# 条件语句

## if
&emsp;&emsp;if条件是不需要括号的
		package main
		
		import "io/ioutil"
		
		func main () {
		    const filename = "abc.txt"
		    contents, err := ioutil.ReadFile(filename)
		    if err != nil {
		        // no such file or directory abc.txt
		        fmt.Println(err)
		    } else {
		        // ...
		        fmt.Printf("%s\n", contents)
		    }
		}
		

&emsp;&emsp;上面这样写比较繁琐，golang的if是可以像for一样写的。但是这样做的话，像下面的contents和err的作用域就只能在if结构中访问了。

		package main
		
		import "io/ioutil"
		
		func main () {
		    const filename = "abc.txt"
		    if contents, err := ioutil.ReadFile(filename); err != nil {
		        // no such file or directory abc.txt
		        fmt.Println(err)
		    } else {
		        // ...
		        fmt.Printf("%s\n", contents)
		    }
		    fmt.Println(contents) // error
		}
		

## switch

&emsp;&emsp;switch会自动break，除非使用fallthrough。Golang的switch非常灵活，表达式不鄙视常量或整数，如果switch没有表达式，它会匹配true

		
		func grade(scare int) string {
		    g := ""
		    switch {
		        case score < 0 || score > 100 : 
		        panic(fmt.Sprintf("Wrong score: %d", score))
		        case score < 60: g = "F"
		        case score < 80: g = "C"
		        case score < 90: g = "B"
		        case score < 100: g = "A"
		    }
		}

# 循环

## for

&emsp;&emsp;for条件不需要括号；for的条件可以圣洛初始条件，结束条件，递增表达式

		package main
		
		import (
		    "fmt"
		    "strconv"
		)
		
		// 省略初始条件
		func convertToBin(n int) string{
		    result := ""
		    for ; n > 0; n /= 2 {
		        lsb := n % 2 // 最低位
		        result = strconv.Itoa(lsb) + result
		    }
		    return result
		}
		
		// 一行一行的读
		// 省略起始条件，省略自增表达式
		func printFile(filename string) {
		    file, err := os.Open(filename)
		    if err != nil {
		        panic(err)
		    } 
		    scanner := bufio.NewScanner(file)
		    
		    for scanner.Scan() {
		        fmt.Println(scanner.Text())
		    }
		}
		
		// 全部省略
		func forever() {
		    for {
		        fmt.Println("abc")
		    }
		}
		
		func main() {
		    fmt.Println(convertToBin(5),convertToBin(13)) // 101, 1101
		}

# 函数

&emsp;&emsp;函数名在前，返回值在后，参数也一样，名在前类型在后

		func eval(a, b int, op string) int {
		    switch op {
		        case "+": return a + b
		        case "-": return a - b
		        case "*": return a * b
		        case "/": return a / b
		        default: panic("unsupported operation: " + op)
		    }
		}

&emsp;&emsp;函数可以返回多个值

		// 13 / 3 = 4 ... 1
		func div(a, b int) (q, r int) {
		    return a / b, a % b
		    //q = a / b
		    //r = a % b
		    //return // 自动返回q r
		}

&emsp;&emsp;上面的eval的写法比较繁琐，在golang中函数是一等公民，我们使用函数式编程实现一下。这里只是简单的展示一下，后期会详细讲解函数式编程。

		
		func apply(op func(int, int) int, a, b int) int {
		    p := reflect.ValudOf(op).Pointer()
		    opName := runtime.FuncForPC(p).Name()
		    fmt.Printf("Calling function %s with args (%d, %d)", opName, a, b)
		    return op(a, b)
		}
		
		func main() {
		    fmt.Println(
		        // 匿名函数
		        apply(func(a, b int) int{
		            return int(math.Pow(float64(a), float64(b)))
		        }, 3, 4)
		    )
		}
		

&emsp;&emsp;golang没有默认参数、可选参数、操作符重载、函数重载这些花销的东西，唯一有的是可变参数列表。"..."函数可以接收多个参数，还有就是后期的数组概念，可以将其打散

		func sum(numbers ...int) int {
		    s := 0
		    for i := range numbers {
		        s += numbers[i]
		    }
		    return s
		}

# 指针

&emsp;&emsp;首先我们不用怕，golang的指针比C语言的指针简单多了，虽然它是一个很重要的概念，但是它也是一个挺简单的概念。

- 指针不可以参与运算

> Go语言只有值传递一种方式，详见[杂谈 | golang是值传递](https://facedamon.github.io/post/golang/golang是值传递/)
