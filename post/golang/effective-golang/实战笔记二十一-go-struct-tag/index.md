
>转载自 https://www.flysnow.org

## JSON字符串对象转换
		type User struct{
		    Name string `name`
		    Age int `age`
		}
		func main(){
		    var u User
		    h := `{"name": "张三", "age":15}`
		    err := json.Unmarshall([]byte(h), &u)
		    if err != nil {
		        fmt.Println(err)
		    }else{
		        fmt.Println(u)
		    }
		}
上面的例子就是json字符串转User对象的例子，这里主要利用的就是User这个结构体对应的字段Tag，json解析的原理就是通过反射获得每个字段的tag，然后把解析的json对应的值赋给它们。

利用字段Tag不光可以把json字符串转为结构体对象，还可以把结构体对象转为json字符串。
		newJson,err := json.Marshall(&u)
		fmt.Println((string(newJson)))
这样就可以转为一个新的字符串了，通过打印，可以看到和输入的字符串一样。
## 反射获取字段Tag
字段的Tag是标记到字段上的，所以我们可以通过先获取字段，然后再获取字段上的Tag。
		func main(){
		 var u User
		  t := reflect.TypeOf(u)
		    for i:=0; i < t.NumField(); i++{
		        sf := t.Field(i)
		        fmt.Println(sf.Tag)
		    }
		}
获取字段上一篇我们提到过，获取字段后，调用`.Tag`就获取到对应的Tag字段了。
## 字段Tag的键值对
很多时候我们的一个struct不止具有一个功能，比如我们需要json的互转、还需要BSON以及ORM解析的互转，所以一个字段可能对应多个不同的Tag，以便满足不同的功能场景。

Go Struct为我们提供了键值对的Tag，来满足我们以上的需求。
		type User struct{
		    Name string `json: "name"`
		    Age int `json: "age"`
		}
		func main(){
		    var u User
		    t := reflect.TypeOf(u)
		    for i:=0; i < t.NumField();i++{
		        sf := t.Field(i)
		        fmt.Println(sf.Tag.Get("json"))
		    }
		}
以上例子，使用了键值对的方式配置tag，key-value以冒号分开，这里的key为json，所以我们可以通过这个key获取对应的值，也就是通过`.Tag.Get("json")`方法。Get方法就是通过一个key获取对应的tag设置。除此之前，我们还可以设置多个key，来满足我们上面说的场景。
		type User struct{
		    Name string `json: "name" bson: "b_name"`
		    Age int `json: "age" bson: "b_age"`
		}
		
		func main(){
		    var u User
		    t := reflect.TypeOf(u)
		    for i:=0; i < t.NumField();i++{
		        sf := t.Field(i)
		        fmt.Println(sf.Tag.Get("json"), ",", sf.Tag.Get("bson"))
		    }
		}
多个key使用空格进行分开，然后使用Get方法获取不同Key的值。

struct tag可以提供字符串到struct的映射能力，以便我们作转换，除此之外，还可以为字段的元数据的配置，提供我们需要的配置，比如生成Swagger文档。
