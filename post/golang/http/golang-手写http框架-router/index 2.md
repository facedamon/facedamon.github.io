
# 摘要

- 使用trie树实现`动态路由`解析
- 支持两种模式`:name`和`*filepath`

## trie树

&emsp;&emsp;上一节，我们用了一个非常简单的`map`结果来存储了路由表，使用`map`存储键值对，索引非常高效，但是有一个弊端，键值对的存储的方式，只能用来索引静态路由。如果我们想支持类似于`/hello/:name`这样的动态路由该怎么办呢？所谓动态路由，即一条路由规则可以匹配某一类型而非某一条固定的路由。例如`/hello/:name`，可以匹配`/hello/geew`, `/hello/zack`等。

&emsp;&emsp;动态路由有很多种实现方式，支持的规则、性能等有很大的差异。例如开源的路由[mux](https://github.com/gorilla/mux)，支持在路由中嵌入正则表达式，例如`/p[0-9A-Za-z]+`，即路径中的参数仅匹配数字和字母；另一个开源实现`httprouter`就不支持正则表达式。著名的Web开源框架`gin`在早期的版本，并没有实现自己的路由，而是直接使用了`httprouter`，后来gin作者放弃了`httprouter`，自己实现了一个版本。

&emsp;&emsp;实现动态路由最常见的数据结构，就是前缀树`trie`。前缀树每一个节点的所有子节点都拥有相同的前缀。这种结构非常适合适用于路由匹配，比如我们定义了如下路由规则：

- **/:lang/doc**
- **/:lang/tutorial**
- **/:lang/intro**
- **/about**
- **/p/blog**
- **/p/related**


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/trie_router.jpg)

&emsp;&emsp;HTTP请求的路径恰好是由`/`分割的多端构成，因此，每一段可以作为前缀树的一个节点。我们通过树结构查询，如果中间某一层的节点都不满足条件，那么就说明没有匹配到路由。

&emsp;&emsp;接下来我们实现的动态路由具备以下两个功能。

- **参数匹配：** 例如 `/p/:lang/doc`，可以匹配 `/p/c/doc`和`/p/go/doc`.
- **通配：** 例如 `/static/*filepath`, 可以匹配 `/static/fav.ico`,也可以匹配`/static/js/jQuery.js`.这种模式常用于静态服务器，能够递归地匹配子路径.

## trie前缀树实现

		// trie.go
		package geew
		
		type node struct {
		    // 待匹配路由 eg: /p/:lang/doc
		    pattern string
		    //路由的一部分 eg: :lang
		    part string
		    //子节点    eg: [doc, intro]
		    children []*node
		    //是否精确匹配 part含有:或*时为true
		    isWild bool
		}

&emsp;&emsp;与普通树不同，为了实现动态匹配，加上了`isWild`这个参数。即当我们匹配`/p/go/doc`这个路由时，第一层节点，`p`精准匹配到了`p`，第二层节点，`go`模糊匹配到`:lang`，那么将会把`lang`这个参数赋值为`go`，继续下一层匹配。我们将匹配的逻辑包装为一个**辅助函数**。

		// 第一次匹配成功的节点
		func (n *node) matchChild(part string) *node {
		    for _, c := range n.children {
		        if c.part == part || c.isWild {
		            return c
		        }
		    }
		    return nil
		}
		
		// 匹配所有子节点
		func (n *node) matchChildren(part string) *node {    
		    nodes := make([]*node, 0)
		    for _, c := range n.children {
		        if c.part == part || c.isWild {
		            nodes = append(nodes, c)
		        }
		    }
		    return nodes
		}

&emsp;&emsp;对于路由来说，最重要的当然是注册与匹配了。开发服务时，注册路由规则，映射handler。访问时，匹配路由规则，查找对应的handler。因此，trie树需要支持节点的插入与查询。**插入功能:递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个，有一点需要注意，`/p/:lang/doc`只有在第三层节点(叶子节点)，即`doc`节点，`pattern`才会被设置为`/p/:lang/doc`。**

		func (n *node) insert(pattern string, parts[]string, height int) {
		    //只有在叶子节点才会设置pattern
		    if len(parts) == height {
		        n.pattern = pattern
		        return
		    }
		
		    p := parts[height]
		    c := n.matchChild(p)
		    if c == nil {
		        c = &node{part: p, isWild: p[0] == ':'|| p[0] == '*'}
		        n.children = append(n.children, c)
		    }
		    c.insert(pattern, parts, height+1)
		}

&emsp;&emsp;**当插入完成时，只有叶子节点的pattern是被赋值的，叶子节点的前缀节点的pattern是空的，因此，查询功能：当查询结束时，可以使用`n.pattern == ''`来判断路由是否匹配成功，递归查询每一层的节点，退出条件是，匹配到了`*`，或者匹配到了叶子节点**

		func (n *node) search(parts []string, height int) *node {
		    if len(parts) == height || strings.HasPrefix(n.part, "*") {
		        if n.pattern == "" {
		            return nil
		        }
		        // 匹配到叶子节点
		        return n
		    }
		
		    p := parts[height]
		    cs := n.matchChildren(p)
		
		    for _, c := range cs {
		        rs := c.search(parts, height+1)
		        if rs != nil {
		            return rs
		        }
		    }
		    return nil
		}

## router

&emsp;&emsp;接下来我们将trie树应用到路由中去吧。我们使用`roots`来存储每种请求方式的trie树根节点。使用`handlers`存储每种请求方式的`HandlerFunc`。getRoute函数解析了`:`和`*`两种通配符，返回一个map。例如`/p/go/doc`匹配到`/p/:lang/doc`，解析结果为`{lang: "go"}`。

		// router.go
		package geew
		
		type router struct {
		    roots map[string]*node
		    handlers map[string]HandlerFunc
		}
		
		func newRouter() *router {
		    return &router{
		        roots: make(map[string]*node),
		        handlers: make(map[string]HandlerFunc)
		    }
		}
		
		func parsePattern(pattern string) []string {
		    vs := strings.Split(pattern, "/")
		
		    ps := make([]string, 0)
		    for _, item := range vs {
		        if item != "" {
		            ps = append(ps, item)
		            if item[0] == '*' {
		                break
		            }
		        }
		    }
		    return ps
		}
		
		func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
		    // register router
		    // constains ':' or '*'
		    ps := parsePattern(pattern)
		
		    key := method + "-" + pattern
		    _, ok := r.roots[method]
		    if !ok {
		        r.roots[method] = &node{}
		    }
		    r.roots[method].insert(pattern, parts, 0)
		    r.handlers[key] = handler
		}
		
		func (r *router) getRoute(method string, path string) (*node, map[string]string) {
		    // search node
		    // not constains ':' and '*'
		    ps := parsePattern(pattern)
		    params = make(map[string]string)
		    root, ok := r.roots[method]
		
		    if !ok {
		        return nil, nil
		    }
		
		    n := root.search(ps, 0)
		
		    if n != nil {
		        parts := parsePattern(n.pattern)
		        for index, part := range parts {
		            if part[0] == ':' {
		                // :lang
		                //[1:]=lang
		                params[part[1:]] = ps[index]
		            }
		            if part[0] == '*' && len(part) > 1 {
		                // *filepath
		                // [1:]=filepath
		                params[part[1:]] = strings.Join(ps[index:], "/")
		                break
		            }
		        }
		        return n, params
		    }
		    return nil, nil
		}

## context变化

&emsp;&emsp;在HandlerFunc中，希望能够访问到解析的参数，因此需要对context对象增加一个属性和方法来提供对路由参数(Restful)的访问。我们将解析后的参数存储到`Params`中，通过`c.Param("xxx")`的方式获取对应的值。

		// context.go
		package geew
		
		type Context struct {
		    Writer http.ResponseWriter
		    Req *http.Request
		    Path string
		    Method string
		    Params map[string]string
		    StatusCode int
		}
		
		func (c *Context) Param(key string) string {
		    value, _ := c.Params[key]
		    return value
		}

## router变化

&emsp;&emsp;router胡变化比较小，比较重要的一点是，在调用匹配到的handler前，将解析出来的路由参数赋值给Params.

		// router.go
		package geew
		
		func (r *router) handle(c *Context) {
		    n, params := r.getRoute(c.Method, c.Path)
		    if n != nil {
		        c.Params = params
		        key := c.Method + "-" + n.pattern
		        r.handlers[key](c)
		    } else {
		        c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		    }
		}

## 测试

		// router_test.go
		package geew
		
		func TestParsePattern(t *testing.T) {
		    ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
			ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
			ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
			if !ok {
				t.Fatal("test parsePattern failed")
			}
		}
		
		func newTestRouter() *router {
			r := newRouter()
			r.addRoute("GET", "/", nil)
			r.addRoute("GET", "/hello/:name", nil)
			r.addRoute("GET", "/hello/b/c", nil)
			r.addRoute("GET", "/hi/:name", nil)
			r.addRoute("GET", "/assets/*filepath", nil)
			return r
		}
		
		func TestGetRoute(t *testing.T) {
			r := newTestRouter()
			n, ps := r.getRoute("GET", "/hello/geew")
		
			if n == nil {
				t.Fatal("nil shouldn't be returned")
			}
		
			if n.pattern != "/hello/:name" {
				t.Fatal("should match /hello/:name")
			}
		
			if ps["name"] != "geew" {
				t.Fatal("name should be equal to 'geew'")
			}
		
			fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
		
		}
