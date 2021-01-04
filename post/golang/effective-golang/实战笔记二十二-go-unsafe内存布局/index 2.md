
>转载自 https://www.flysnow.org

## 简介
unsafe，顾名思义，是不安全的，Go定义这个包名也是这个意思，让我们尽可能不要使用它，如果你使用它，看到了这个名字，也会想到尽可能的不要使用它，或者更小心的使用它。

虽然这个包不安全，但是它也有它的优势，那就是可以`绕过Go的内存安全机制，直接对内存进行读写`，所以有时候因为性能的需要，会冒险使用该包，对内存进行操作。

## SizeOf函数
SizeOf函数可以返回一个类型所占用的内存大小，这个大小只与类型有关，和类型对应的变量存储的内容大小无关，比如bool类型占用一个字节、int8也占用一个字节。
		func main() {
			fmt.Println(unsafe.Sizeof(true))
			fmt.Println(unsafe.Sizeof(int8(0)))
			fmt.Println(unsafe.Sizeof(int16(10)))
			fmt.Println(unsafe.Sizeof(int32(10000000)))
			fmt.Println(unsafe.Sizeof(int64(10000000000000)))
			fmt.Println(unsafe.Sizeof(int(10000000000000000)))
		}
对于整型来说，占用的字节数意味着这个类型存储数字范围的大小，比如int8占用一个字节也即是8bit，所以它可以存储的大小范围是-128～127，也就是-2^(n-1)～2^(n-1)-1,n表示bit，int8表示8bit，int16表示16bit，其它以此类推。
对于和平台有关的int类型，这个要看平台是32位还是64位，会取最大的。比如我自己测试，以上输出，会发现int和int64的大小是一样的，因为我的是64位平台操作系统。
		func SizeOf(x ArbitraryType) uintptr
以上是SizeOf的函数定义，它接受一个ArbitraryType类型的参数，返回一个uintptr类型的值。这里的ArbitraryType不用关心，它知识一个占位符，为了文档的考虑导出了该类型，但是一般不会使用它，我们只需要知道它表示任何类型，也就是我们这个函数可以接受任意类型的数据。

## alignof函数
Alignof返回一个类型的对齐值，也可以叫做对齐数或者对齐倍数。对齐值是一个和内存对齐有关的值，合理的内存对齐可以提高内存读写的性能，关于内存对齐的知识可以灿开相关文档，这里不再介绍。
		func main() {
			var b bool
			var i8 int8
			var i16 int16
			var i64 int64
		
			var f32 float32
		
			var s string
		
			var m map[string]string
		
			var p *int32
		
			fmt.Println(unsafe.Alignof(b))
			fmt.Println(unsafe.Alignof(i8))
			fmt.Println(unsafe.Alignof(i16))
			fmt.Println(unsafe.Alignof(i64))
			fmt.Println(unsafe.Alignof(f32))
			fmt.Println(unsafe.Alignof(s))
			fmt.Println(unsafe.Alignof(m))
			fmt.Println(unsafe.Alignof(p))
		
		}
以上例子的输出，可以看到，对齐值一般是2^n,最大不会超过8.Alignof的函数定义和SizeOf基本一样。这里需要注意的是每个人的操作系统运行的结果可能不大一样，大同小异。
		func Alignof(x ArbitratyType) unitptr
此外，获取对齐值还可以使用反射包的函数，也就是说：unsafe.ALignof(x)等价于reflect.TypeOf(x).Align().
