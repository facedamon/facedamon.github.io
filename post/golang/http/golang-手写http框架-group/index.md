
# 摘要

&emsp;&emsp;**分组控制**是http框架提供的基础功能之一。所谓分组，是指路由的分组。如果没有路由分组，我们需要针对每一个路由进行控制。但是真实的业务场景中，往往**某一组路由需要相似的处理**。例如：

- 以`/post`开头的路由匿名可访问
- 以`/admin`开头的路由需要鉴权
- 以`/api`开头的路由是RESTful接口，可以对接第三方，需要三方平台鉴权

&emsp;&emsp;大部分情况下的路由分组，是以相同的前缀来区分的。我们实现的分组控制也是以前缀来区分，并且支持分组的**嵌套**。例如`/post`是一个分组，`/post/a`和`/post/b`可以是该分组下的子分组。作用在`/post`分组上的中间件(middleware)，也都会作用在子分组，子分组还可以应用自己特有的中间件。

&emsp;&emsp;**中间件可以给框架提供无限的扩展能力**，应用在分组上，可以使得分组控制的收益更为明显，而不是共享相同的路由前缀。例如`/admin`分组，可以应用鉴权中间件，`/`分组应用日志中间件，`/`是默认的最顶层分组，也就意味着给所有的路由，即整个框架增加了记录日志的能力。提供扩展能力支持中间件的内容，我们将在下一节当中介绍。

## 分组嵌套

&emsp;&emsp;一个group对象需要哪些属性？首先想到的肯定是前缀(prefix)，比如`/`，或者`/api`;要支持分组嵌套，那么需要知道当前分组的父节点(parent)是谁。当然了，按照我们一开始的分析，中间件是应用在分组上的，那还需要存储应用在该分组上的中间件(middlewares)。还记得我们之前调用函数`(*Engine).addRoute()`来映射所有的路由规则和Handler。如果group对象需要直接映射路由规则的话，比如我们想在使用框架时，这么调用：

		r := geew.New()
		v1 := r.Group("/v1")
		v1.GET("/", func(c *Context){
		    c.HTML(http.StatusOK,"<h1>Hello Geew</h1>")
		})

&emsp;&emsp;那么**group对象，还需要有访问`router`的能力**，为了方便，我们可以在group中，保存一个指针，指向`Engine`,整个框架的所有资源都是由`Engine`统一协调的，那么就可以通过Engine间接访问各种接口了。

		// group.go
		package geew
		
		type group struct {
		    prefix  string
		    // support middleware
		    middlewares []HandlerFunc
		    // support nesting
		    parent  *group
		    // all groups share a engine instance
		    engine  *Engine
		}

&emsp;&emsp;我们还可以进一步抽象，将`Engine`作为最顶层的分组，也就是说Engine拥有group所有的能力，这样就可以将和路由有关的函数，都交给group实现了。

		// geew.go
		package geew
		
		type Engine struct {
		    *group
		    router  *router
		    // store all group
		    groups  []*group
		}
		
		func New() *Engine {
		    e := &Engine{router: newRouter()}
		    e.group = &group{engine: e}
		    e.groups = []*group{e.group}
		}

		// group.go
		package geew
		
		// Group is defines to create a new group
		// remember all groups share the same engine instance
		func (g *group) Group(prefix string) *group {
		    e := g.engine
		    newg := &group{
		        prefix: g.prefix + prefix,
		        parent: g,
		        engine: e,
		    }
		    e.groups = append(e.groups, newg)
		    return newg
		}
		
		func (g *group) addRoute(method string, comp string, handler HandlerFunc) {
		    pattern := g.prefix + comp
		    // add route to trie
		    g.engine.router.addRoute(method, pattern, handler)
		}
		
		// GET defines the method to add GET req
		func (g *group) GET(pattern string, handler HandlerFunc) {
		    g.addRoute("GET", pattern, handler)
		}
		
		// POST defines the method to add POST req
		func (g *group) POST(pattern string, handler HandlerFunc) {
		    g.addRoute("POST", pattern, handler)
		}

&emsp;&emsp;仔细观察`addRoute`方法，调用了`group.engine.router.addRoute`来实现了路由的映射。由于`Engine`组合了`group`所有的属性和方法，，这样的实现，我们既可以像原来一样添加路由，也可以添加分组路由。

## Test

		// example_test.go
		package example
		
		func TestGroup(t *testing.T) {
		    r := geew.New()
		    r.GET("/index", func(c *geew.Context) {
				c.HTML(http.StatusOK, "<h1>Index Page</h1>")
			})
			v1 := r.Group("/v1")
			{
				v1.GET("/", func(c *geew.Context) {
					c.HTML(http.StatusOK, "<h1>Hello Geew</h1>")
				})
		
				v1.GET("/hello", func(c *geew.Context) {
					// expect /hello?name=geew
					c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
				})
			}
			v2 := r.Group("/v2")
			{
				v2.GET("/hello/:name", func(c *geew.Context) {
					// expect /hello/geew
					c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
				})
				v2.POST("/login", func(c *geew.Context) {
					c.JSON(http.StatusOK, gee.H{
						"username": c.PostForm("username"),
						"password": c.PostForm("password"),
					})
				})
		
			}
		
			r.Run(":9999")
		}

- curl测试

		// shell
		$ curl 'http://127.0.0.1:9999/v1/hello?name=geew'
		hello geew, you`re at /v1/hello
		
		$ curl 'http://127.0.0.1:9999/v2/hello/geew'
		hello geew, you`re at /hello/geew
