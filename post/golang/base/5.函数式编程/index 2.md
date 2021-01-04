# 函数式编程

> 函数式编程 VS 函数指针

- 函数是一等公民：参数，变量，返回值都可以是函数
- 闭包

		package main
		
		func adder() func(int) int {
		    sum := 0 // 自由变量
		    return func (v int) int {
		        sum += v //v 局部变量
		        return sum
		    }
		}
		
		func main() {
		    a := adder()
		    for i := 0; i < 10; i++ {
		        fmt.Printf("0 + 1 + ... + %d = %d\n", i, a(i))
		    }
		}
![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/闭包.png)

&emsp;&emsp;所以，当函数返回的时候，返回的是一个闭包。我们return这个func不是说返回这个代码，而是返回了这个函数以及对sum的引用

> golang语言闭包应用

- 斐波那契数列(func 返回值)

		package main
		
		// 1, 1, 2, 3, 5, 8, 13, ...
		//    a, b
		//       a  b
		func fibonacci func() int {
		    a, b := 0, 1 //自由变量
		    return func() int {
		        a, b = b, a + b
		        return a
		    }
		}
		
		func main() {
		    f := fibonacci() //生成器
		    fmt.Println(f()) //1
		    fmt.Println(f()) //1
		    fmt.Println(f()) //2
		    fmt.Println(f()) //3
		    fmt.Println(f()) //5
		    fmt.Println(f()) //8
		    fmt.Println(f()) //13
		    fmt.Println(f()) //21
		}
- 为函数实现接口

&emsp;&emsp;上述数列调用有个缺点，我每次都要手动调用f, 那能不能像读文件一样读取数列内容呢？我们只要让数列生成器实现Reader接口，每扫描一次文件，就在Read方法里调一次数列生成器
		type intGen func() int
		
		func (g intGen) Read(p []byte) (n int, err error) {
		    next := g()
		    if next > 10000 {
		        return 0, io.EOF
		    }
		    s := fmt.Sprintf("%d\n", next)
		    return strings.NewReader(s).Read(p)
		}
		
		func printFileContents(reader io.Reader) {
		    scanner := bufio.NewScanner(reader)
		
		    for scanner.Scan() {
		        fmt.Println(scanner.Text())
		    }
		}
		
		func main() {
		    f := fibonacci()
		    printFileContents(s)
		}
&emsp;&emsp;接口型函数详见->[接口型函数](https://facedamon.github.io/post/golang/standard-lib/标准库database/接口型函数/)

- 使用函数来遍历二叉树(func 参数)

		type Node struct {
		    value int
		    Left, Right *Node
		}
		
		func (node *Node) TraverseFunc(f func(*Node)) {
		    if node == nil {
		        return
		    }
		    node.Left.TraverseFunc(f)
		    f(node)
		    node.Right.TraverseFunc(f)
		}
		
		func (node *Node) Traverse() {
		    node.TraverseFunc(func(n *Node){
		        n.Print()
		    })
		    fmt.Println()
		}
		
		func (node *Node) Print() {
		    fmt.Println(node.value, " ")
		}

- 更为自然，不需要修饰如何访问自由变量
- 没有Lambda表达式
