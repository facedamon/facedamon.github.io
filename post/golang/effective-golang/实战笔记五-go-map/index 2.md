
>转载自 https://www.flysnow.org

## 1. 内部实现(<a href="/post/golang/interview/golang-1/">关于Hash散列细枝末节详见Java下的散列表)</a>
Map是给予散列表来实现的，就是我们常说的Hash表，所以我们每次迭代map的时候，打印的key和value是无序的，每次迭代的都不一样，即使我们按照一定的顺序存在也不行。
map的散列表包含一组桶，每次存储和查找键值对的时候，都要先选择一个桶。如何选择桶呢？就是把指定的键传给散列函数，就可以索引到相应的桶了，进而找到对应的键值。
这种方式的好处在于，存储的数据越多，索引分布越均匀，所以我们访问键值对的速度也就越快，当然存储的细节还有很多，大家可以参考hash相关的知识，这里我们只要记住map存储的是无序的键值对集合。
## 2. 声明和初始化
map的创建有make函数，map字面量。make函数我们用它创建过切片，除此之外，它还可以用来创建map。
		dict := make(map[string]int)
示例中创建了一个键类型为string的，值类型为int的map。现在创建好之后，这个map是空的，里面什么都没有，我们给存储一个键值对。
		dict := make([string]int)
		dict["张三"] = 43
存储了一个key为张三的，value为43的键值对数据
此外还有一种使用map字面量的方式创建和初始化map，对于上面的例子，我们可以等同实现。
		dict := map[string]int{"张三":43}
使用一个大括号进行初始化，键值对通过:分开，如果要同时初始化多个键值对，使用逗号分割。
		dict := map[string]int{"张三":43,"莉丝":50}
当然我们可以不指定任何键值对，也即是一个空map
		dict := map[string]int{}
不管怎么样，**使用map字面量创建一定要带上大括号。**如果我们要创建一个nil的map怎么做呢？nil的map是未初始化的，所以我们可以只声明一个变量，既不能使用map字面量，也不能使用make函数分配内存。
		var dict map[string]int
这样就好了，但是这样我们是不能操作存储键值对的，必须要初始化后才可以，比如使用make函数为其开启一块可以存储数据的内存，也就是初始化。
		var dict map[string]int
		dict = make(map[string]int)
		dict["张三"] = 43
		fmt.Println(dict)
**map的键可以是任何值，键的类型可以是内置的类型，也可以是结构类型，但是不管怎么样，这个键可以使用==运算符进行比较，所以像切片、函数以及含有切片的结构类型就不能用于map的键了，因为它们具有引用的语义，不可比较。**
对于map的值来说，就没有什么限制了，切片这种在键里不能用的，完全可以用在值里。
## 3. 使用map
map的使用很简单，和数组切片差不多，数组切片是使用索引，map是通过键
		dict := make(map[string]int)
		dict["张三"] = 43
以上示例，如果张三存在，则对其值修改，不过不存在，则新增这个键值对。
获取一个map键的值也很简单，和存储差不多，还是给予上面的例子。
		age := dict["张三"]
在Go Map中，如果我们获取了一个不存在的键的值，也是可以的，返回的是值类型的零值，这样就会导致我们不知道是真的存在一个为零值的键值对呢，还是说这个键值对本身就不存在。对此，map为我们提供了检测一个键值对是否存在的方法。
		age, exists := dict["莉丝"]
看这个例子，和获取键的值没有太大区别，只是多了一个返回值。第一个返回值是键的值；第二个返回值标记这个键是否存在，这是一个boolean类型的变量，我们判断它就知道该键是否存在了。这也是go多值返回的好处。
如果我们像删除一个map中的键值对，可以使用go内置的`delete`函数
		delete(dict,"张三")
delete函数接受两个参数，第一个是要操作的map，第二个是要删除的map的键
**delete函数删除不存在的键也是可以的，只是没有任何作用**
## 4. 遍历map
使用for range风格的循环，和遍历切片一样。
		dict := map[string]int{"张三"：43}
		for key, value := range dict {
		    fmt.Println(key, value)
		}
这里的range返回两个值，第一个是map的键，第二个是map的键对应的值。这里再次强调，这种遍历是无序的，也就是键值对不会按既定的数据出现，如果想按顺序遍历，可以先对map中的键排序，然后遍历排序好的键，把对应的值取出来。
		func main(){
		    dict := map[string]int{"王五"：60，"张三"： 43}
		    var names []string
		    for name := range dict {
		        names = append(names, name)
		    }
		    sort.Strins(names)
		    for _, key := range names {
		        fmt.Println(key, dict[key])
		    }
		}
这个例子里有何技巧，range一个map的时候，也可以使用一个返回值，这个默认的返回值就是map的键。
## 5. 在函数间传递Map
`函数间传递map是不会拷贝一个该map的副本的`，也就是说如果一个map传递给一个函数，该函数对这个map做了修改，那么这个map的所有引用，都会感知到这个修改。
		func main(){
		    dict := map[string]int{"王五"：60，"张三"：43}
		    modify(dict)
		    fmt.Println(dict["张三"])
		}
		func modify(dict map[string]int) {
		    dict["张三"] = 10
		}
上面这个例子输出的结果是10，也就是说已经被函数给修改了，可以证明传递的并不是一个map的副本。这个特性和切片类似，如果复制整个map的代价太大了。
