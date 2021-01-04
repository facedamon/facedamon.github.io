# 结构体和方法
- golang仅支持封装，不支持继承和多态
- golang没有class，只有struct
- golang结构体没有构造方法这一说法

		type treeNode struct {
		    value int
		    left, right *treeNode
		}
		
		// 有时候我们想要实现类似的构造方法
		// 那完全可以使用普通的工厂方法实现
		func createNode(value int) *treeNode {
		    //返回了一个局部变量的地址给外面用
		    //这如果在C++中，是个典型的错误
		    return &treeNode{value}
		}
		
		func main() {
		    var root treeNode
		    fmt.Println(root) // {0, nil, nil}
		    root = treeNode{value: 3}
		    root.left = &treeNode{}
		    root.right = &treeNode{5, nil, nil}
		    root.right.left = new(treeNode)
		    root.left.right = createNode(2)
		}

> 如上，createNode结构体创建在堆上还是栈上?

&emsp;&emsp;在C++中，函数局部变量都是分配在栈上，函数一旦退出，局部变量就会立刻被销毁，如果要传递出去，必须在堆上分配，然后手动释放。

&emsp;&emsp;在Java中，new都是分配在堆上的，这样才会有垃圾回收机制。

&emsp;&emsp;答案是不需要知道。编译器运行环境自己决定的，如果这里的createNode没有取地址并返回，编译器认为这个变量不需要给外面用，就直接在栈上分配。如果取地址返回给外面用，那编译器就会在堆上分配，参与垃圾回收

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/tree结构.png)

## 方法(接收者)

&emsp;&emsp;print是给treeNode来接收的。
		func (node treeNode) print() {
		    fmt.Print(node.value, " ")
		}
		// 值传递
		func (node *treeNode) setValue(value int) {
		    node.value = value
		}
		
		// 中序遍历
		func (node *treeNode) traverse() {
		    if node == nil {
		        return
		    }
		
		    node.left.traverse()
		    node.print()
		    node.right.traverse()
		}
		
		var pRoot *treeNode
		pRoot.setValue(200)// nil 调用方法
		pRoot = &root
		pRoot.setValue(300)
		pRoot.print()
		
		
- 只有使用指针才可以改变结构内容
- nil指针也可以调用方法

> 值接收者 VS 指针接收者

- 要改变内容必须使用指针接收者
- 结构过大也考虑使用指针接收者
- 一致性：如果有指针接收者，最好都是指针接收者

# 包和封装

> 名字一般使用CamelCase

- 首字母大写：public    不同包可以引用访问
- 首字母小写: private   同一包可以访问

> 每个目录一个包

- main包包含可执行入口
- 为结构体定义的方法必须放在同一个包内
- 可以是不同文件

> 如何扩充系统类型或者别人的类型

1. 定义别名
2. 使用组合

- 使用组合
		type myNode struct {
		    node *tree.Node
		}
		
		func (md *myNode) postOrder() {
		    if md == nil || md.node == nil{
		        return
		    }
		    left := myNode{md.node.Left}
		    left.postOrder()
		    right := myNode{md.node.Right}
		    right.postOrder()
		    myNode.node.Print()
		}

- 定义别名
		package queue
		
		type Queue []int
		
		func (q *Queue) Push(v int) {
		    *q = append(*q, v)
		}
		
		func (q *Queue) Pop() int {
		    head := (*q)[0]
		    *q = (*q)[1:]
		    return head
		}
		
		func (q *Queue) IsEmpty() bool {
		    return len(*q) == 0
		}
		package main
		
		import "queue"
		
		func main() {
		    q := queue.Queue{1}
		    q.push(2)
		    q.push(3)
		    fmt.Println(q.Pop())
		    fmt.Println(q.Pop())
		    fmt.Println(q.IsEmpty())
		    fmt.Println(q.Pop())
		    fmt.Println(q.IsEmpty())
		    // 1 2 false 3 true
		    // 这里的q 在多次操作后，已经不是最开始的q了，地址已经改变了
		}
# GOPATH
&emsp;&emsp;默认在~/go(unix,linux)，%USERPROFILE%\go(windows)

> 官方推荐：所有项目和第三方库都放在同一个GOPATH下

&emsp;&emsp;也可以将每个项目放在不同的GOPATH
