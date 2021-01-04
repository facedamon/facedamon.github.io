
> 转载自 https://www.flysnow.org/2016/12/30/golang-function-interface.html


# 简介

&emsp;&emsp;接口型函数，指的是用函数实现接口，这样在调用的时候就会非常简便，这种方式使用于只有一个函数的接口。

## 原始接口

&emsp;&emsp;我们以迭代一个map为例，演示这一技巧.

		type Handler interface {
			Do(k, v interface{})
		}
		
		func Each(m map[interface{}]interface{}, h Handler) {
			if m != nil && len(m) > 0 {
				for k, v := range m {
					h.Do(k, v)
				}
			}
		}

&emsp;&emsp;上述实现方式，和Java的多态如出一辙。不同的接口实现对象，实现不同的Do逻辑，上层Each代码不需要改变。

&emsp;&emsp;首先定义一个Handler接口，只有一个Do方法，接收k, v 两个参数。然后定义一个Each函数，该函数的功能就是迭代传递过来的map参数，然后把map的每个k,v值传递给Handler，去做具体的事情，可以是计算，具体由这个Handler的实现来决定，这也是面向接口编程。

&emsp;&emsp;现在我们就以新学期开学，大家自我介绍为例，演示使用Each方法和Handler接口。

		type welcome string
		func (w welcome) Do(k, v interface{}) {
			fmt.Printf("%s, 我叫%s, 今年%d岁\n", w, k, v)
		}
		
		func main() {
			persons := make(map[interface{}]interface{})
			persons["张三"] = 20
			persons["李四"] = 23
			persons["王五"] = 26
		
			var w welcome = "大家好"
			
			Each(persons, w)
		}

## 接口型函数

&emsp;&emsp;以上实现，有两点不好：

1. 因为必须要实现Handler接口，Do这个方法名不能修改，不能定义一个更有意义的名字
2. 必须要新定义一个类型，才可以实现Handler接口，才能使用Each函数

&emsp;&emsp;首先我们解决第一个问题，根据我们具体做的事情定义一个更有意义的方法名。如果调用者改了方法名，那么就不能实现Handler接口了，但是还要继续使用Each方法，那么就提供一个Adapter包装器，将其转化为Handler接口的样子.

		type HandlerFunc func(k, v interface{})
		
		func (f HandlerFunc) Do(k, v interface{}){
			f(k, v)
		}

&emsp;&emsp;以上代码，我们定义了一个新的类型HandlerFunc,它是一个func(k, v interface{})类型，并实现了Handler接口，那么Do方法就是调用HandlerFunc自身，因为HandlerFunc本身就是一个方法。

		func (w welcome) selfInfo(k, v interface{}) {
		          fmt.Printf("%s, my name is %s, and age is %d\n", w, k, v)
		}
		
		Each(persons, HandlerFunc(w.selfInfo))

&emsp;&emsp;还是差不多原来的实现，只是把方法名Do改为selfInfo。`HandlerFunc(w.selfInfo)`不是方法的调用，而是转型，因为selfInfo和HandlerFunc是同一种类型，所以可以强制转型。转型后，因为HandlerFunc实现了Handler接口，所以我们可以继续使用原来的Each方法。

## 进一步

&emsp;&emsp;现在解决了命名问题，但是每次强制转型不太好，我们继续重构，采用新定义函数帮助调用者强制转型。

		func EachFunc(m map[interface{}]interface{}, f func(k, v interface{})) {
			Each(m, HandlerFunc(f))
		}
		
		func main() {
			persons := make(map[interface{}]interface{})
		        persons["one"] = 20
		        persons["two"] = 23
		        persons["three"] = 26
		
		        var w welcome = "hello everyone!"
		
		        //      Each(persons, w)
		        //Each(persons, HandlerFunc(w.selfInfo))
		        EachFunc(persons, w.selfInfo)
		}

&emsp;&emsp;新增EachFunc函数，帮助调用者强制转型。现在我们发现EachFunc函数接收的是一个`func(k, v interface{})`类型的函数，没有必要实现Handler接口，所以新的自定义类型可以去掉不用。

		func selfInfo(k, v interface{}) {
			fmt.Printf("hi, everyone, my name is %s, and age is %d\n", k, v)
		}
		
		func main() {
			...
			EachFunc(persons, selfInfo)
		}

# 延伸

&emsp;&emsp;以上就是关于接口型函数，我们发现这和`http.Handle`方法很像，其实接口`http.Handler`就是这么实现的。

		type Handler interface {
			ServeHTTP(ResponseWriter, *Request)
		}
		
		func Handle(pattern string, handler Handler) {
			DefaultServeMux.Handle(pattern, handler)
		}
		
		func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
			DefaultServeMux.HandleFunc(pattern, handler)
		}

&emsp;&emsp;这是一种非常好的技巧，提供两种函数,既可以以接口的方式使用，也可以以方法的方式，对应上述代码的Each和EachFUnc这两个函数，灵活方便。
