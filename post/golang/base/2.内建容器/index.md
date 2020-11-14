
# 数组

&emsp;&emsp;数量写在类型的前面

		func main() {
		    // var 需要确定大小，可以不用初始化赋值
		    var arr1 [5]int
		    // 简化，需要确定大小，需要初始化赋值
		    arr2 := [3]int{1, 3, 5}
		    // ... 编译器根据初始化赋值的长度确定数组大小
		    arr3 := [...]int{2,4,6,8,10}
		
		    // 二维数组 4个长度为5的int数组
		    grid [4][5]int
		}

&emsp;&emsp;遍历数组

		for _, v:= range arr3 {
		    fmt.Println(v)
		}

> 数组是值类型

		
		func printArray(arr [5]int) {
		    arr[0] = 100
		    for _, v:= range arr {
		    fmt.Println(v)
		    }
		}
		
		func main() {
		    // [100, 0, 0, 0, 0]
		    printArray(arr1)
		    // [100, 4, 6, 8, 10]
		    printArray(arr3)
		    printArray(arr2) // error 长度不一致
		    // [0, 0, 0, 0, 0]
		    // [2, 4, 6, 8, 10]
		    fmt.Println(arr1, arr3)
		}
		
- [10]int 和 [20]int是不同类型
- 调用func f(arr [10] int) 会**拷贝**数组
- 在golang中一般不直接使用数组

# 切片(Slice)

		func main() {
		    arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
		    // [2 3 4 5]
		    fmt.Println("arr[2:6] = ", arr[2:6])
		    // [0 1 2 3 4 5]
		    fmt.Println("arr[:6]", arr[:6])
		    // [2 3 4 5 6 7]
		    fmt.Println("arr[2:]", arr[2:])
		    // [0 1 2 3 4 5 6 7]
		    fmt.Println("arr[:]", arr[:])
		}

&emsp;&emsp;上述arr[]是对arr数组的视图操作。也就是说slice是对数组的视图，那它是如何视图的呢？

		func updateSlice(s []int) {
		    s[0] = 100
		}
		
		func main() {
		    s1 := arr[2:]
		    // [2 3 4 5 6 7]
		    fmt.Println("s1=", s1)
		    s2 := arr[:]
		    // [0 1 2 3 4 5 6 7]
		    fmt.Println("s2=", s2)
		    updateSlice(s1)
		    // [100 3 4 5 6 7]
		    fmt.Println(s1)
		    // [0 1 100 3 4 5 6 7]
		    fmt.Println(arr)
		    updateSlice(s2)
		    // [100 1 100 3 4 5 6 7]
		    fmt.Println(s2)
		    // [100 1 100 3 4 5 6 7]
		    fmt.Println(arr)
		}
> Slice本身没有数据，是对底层array的一个view

## Slice的扩展
&emsp;&emsp;首先我们来看下面这个问题，s1是多少？s2是多少？会报错吗？
		arr := [...]int{0, 1 ,2, 3, 4, 5, 6, 7}
		s1 := arr[2:6]
		s2 := s1[3:5]

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/slice扩展.png)

## Slice 实现

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/slice结构.png)

&emsp;&emsp;Slice的底层有一个Array，有三种颜色，其中ptr指向了Slice开头的那个元素，len说明这个Slice的长度是多少，我们使用[number]直接取值时，只能取到len内的值，number>=len的话，会报错下标越界。cap代表整个Array,从ptr开始到结束整个的长度。在扩展的时候，只要不超过cap都会获取到。

- s1的值为[2 3 4 5], s2的值为[5 6]
- Slice可以向后扩展，不可以向前扩展
- s[i]不可以超越len(s),向后扩展不可以超越底层数组cap(s)

&emsp;&emsp;如何获取len,cap呢？
		// s1=[2 3 4 5], len(s1) = 4, cap(s1) = 6
		fmt.Printf("s1=%v, len(s1)=%d, cap(s1)=%d\n", s1, len(s1), cap(s1))
		
		// s2=[5 6] len(s2) = 2, cap(s2) = 3
		fmt.Printf("s2=%v, len(s2)=%d, cap(s2)=%d\n", s2, len(s2), cap(s2))

## Slice添加元素

		s3 := append(s2, 10)
		s4 := append(s3, 11)
		s5 := append(s4, 12)
		
		// s3, s4, s5 = [5 6 10] [5 6 10 11] [5 6 10 11 12]
		fmt.Println("s3, s4, s5 = ", s3, s4, s5)
		// arr = [0 1 2 3 4 5 6 10]
		// s4 and s5 no longer view arr
		fmt.Println("arr = ", arr)

&emsp;&emsp;我们看到append是可以做到的，cap(s2)=3, s2最后一个元素是7， arr还是一个长度为8的数组，第一次append时，将s2所视图的arr最后一个元素变成10了，所以arr自身变了。后面的append所视图的就不是arr本身了，它们所视图的是一个新的arr，golang会在内存开辟一个新的arr空间，将旧的arr拷贝到新的arr中，并且新的arr的cap设置的更长(长度翻倍，当元素超过1000后，每次增加25%)

- 添加元素时如果超越cap，系统会重新分配更大的底层数组
- 由于值传递的关系，必须接收append的返回值
- s = append(s, val)

## Slice的其它操作

		func main() {
		    var s []int // Zero value for slice is nil and ptr is nil too
		    s1 := []int{} // Zero value for slice is nil and ptr is not nil
		    s1 = []int{2, 4, 6, 8}
		
		    s2 := make([]int, 16)
		    s3 := make([]int, 10, 32)
		
		    fmt.Println("Copying slice")
		    // [2 4 6 8 0 0 0 0 0 0 0 0 0 0 0 0]
		    copy(s2, s1)
		
		    fmt.Println("Deleting elements from slice")// golang slice没有删除的内建函数，我们可以使用reslice实现
		    s2 = append(s2[:3], s2[4:]...)
		
		    fmt.Println("Popping from front")
		    front := s2[0]
		    s2 = s2[1:]
		
		    fmt.Println("Popping from back")
		     tail := s2[len(s2) - 1]
		    s2 = s2[:len(s2) - 1]
		}

# Map

		func main() {
		    m := map[string]string {
		        "name": "ccmouse",
		        "course": "golang",
		        "site": "imooc",
		        "quality": "notbad",
		    }
		
		    m2 := make(map[string]int) // mw == empty map
		
		    var m3 map[string]int // m3 == nil
		
		    fmt.Println("Traversing map")
		    // 无序的
		    for k, v := range m {
		        fmt.Println(k, v)
		    }
		
		    fmt.Println("Getting values")
		    // 如果没获取到，不会报错，是Zero value
		    courseName, ok := m["course"]
		    // golang true
		    fmt.Println(courseName)
		    causeName, ok := m["cause"]
		    // Zero Value false
		    fmt.Println(causeName, ok)
		    // 一般这样写
		    if causeName, ok := m["cause"]; ok {
		        fmt.Println(causeName, ok)
		    } else {
		        fmt.Println("Key does not exist")
		    }
		
		    fmt.Println("Deleting values")
		    delete(m, "name")
		}
- 使用range遍历key, 或者遍历key, value对
- 不保证遍历顺序，如果要顺序，需手动对key排序
- 使用len获取长度

- map的key，在Java中满足key的元素必须实现hashcode&equals方法。在golang中呢
    - 使用哈希表，必须可以比较相等
    - 除了slice，map，function的内建类型都可以作为key
    - struct类型不包含上述字段，也可以作为key

## 寻找最长不含有重复字符的子串
- abcabcbb -> abc
- bbbbb -> b
- pwwkew -> wke


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/map寻找最长不含有重复字符的字串.png)

&emsp;&emsp;假设从左往右进行扫描，扫描到字母X，首先记录当前找到的最长不含有重复字符的字串的开始start,保证了从start到x的区间是不含有重复字符的，那么接下来，我们就要看从start到x-1位置是否含有x这个字符，我们记录lastOccurred[x]即x最后出现的位置，这当然是个map了，我们问map就行了。问出的结果有两种种，1. x从来没有出现过，或者出现在start之前，那么x就可以拿进来了，那么最长字串就可以+1了。2. 出现在中间的区间，那么x就不能拿进来了，这时候就要更新start为lastOccurred[x]+1的位置。

- 对于每一个字母X
    - lastOccurred[x]不存在，或者 < start -> 无需操作
    - lastOccured[x] >= start -> 更新start
    - 更新lastOccurred[x], 更新maxLength

		func lengthOfNonRepeatingSubStr(s string) int{
		    lastOccurred := make(map[byte]int)
		    start, maxLength := 0, 0
		    //byte不支持中文，国际化需要使用rune
		    for i, ch := range []byte(s) {
		        if lastI, ok := lastOccruued[ch];ok && lastI >= start {
		            start = lastI
		        }
		        // 区间，从start到i
		        if i - start + 1 > maxLength {
		            maxLength = i - start + 1
		        }
		        lastOccurred[ch] = i
		    }
		    return maxLength
		}

# 字符和字符串处理

- rune is an alias for int32 and is equivalent to int32 in all ways

		func main() {
		    s := "Yes我爱慕课网!" // UTF-8 中文三个字节
		    fmt.Println(s)
		    for _, b := range []byte(s) {
		        fmt.Printf("%X ", b)
		    } 
		    // 59 65 73 E6 88 91 E7 88 B1 E6 85 95 E8 AF BE E7 BD 91 21
		    // 所以len(s) = 19
		    fmt.Println()
		
		    for i, ch := range s { // ch is a rune
		        fmt.Printf("(%d %X) ", i, ch)
		    }
		    //(0 59)(1 65)(2 73)(3 6211)(6 7231)(9 6155)(12 8BFE)(15 7F51)(18 21)
		    // 将utf8解码后编码转成了unicode
		    fmt.Println()
		
		    fmt.Println("Rune count:", utf8.RuneCountInString(s))// 9
		
		    //接下来是常规的我们遍历string，由于是byte类型，所以如果直接循环的话，会出现三字节的拆分乱码
		    bytes := []byte(s)
		    for len(bytes) > 0 {
		        ch, size := utf8.DecodeRune(bytes)
		        bytes = bytes[size:]
		        fmt.Printf("%c ", ch)
		    }
		    fmt.Println() // Y e s 我 爱慕 课 网 !
		
		    //上面的操作过于繁琐，层面直接告诉我第几个元素是啥就行了
		    for i, ch := range []rune(s) {
		        fmt.Printf("(%d %c)", i, ch)
		    }
		    //(0 Y)(1 e)(2 s)(3 我)(4 爱)(5 慕)(6 课)（7 网）(8 !)
		    // 所以len([]rune(s)) = 9
		    fmt.Println()
		}
> 那么我们来修改一下上段的寻找最长无重复字符的字串问题，国际版

		func lengthOfNonRepeatingSubStr(s string) int{
		    lastOccurred := make(map[rune]int)
		    start, maxLength := 0, 0
		
		    for i, ch := range []rune(s) {
		        if lastI, ok := lastOccruued[ch];ok && lastI >= start {
		            start = lastI
		        }
		        // 区间，从start到i
		        if i - start + 1 > maxLength {
		            maxLength = i - start + 1
		        }
		        lastOccurred[ch] = i
		    }
		    return maxLength
		}

# 学习strings包
- Fields, Split, Join
- Contains, Index
- ToLower, ToUpper
- Trim, TrimRight, TrimLeft
