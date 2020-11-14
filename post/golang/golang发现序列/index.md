
- 发现了一个好用的golang sql package sqlx，有点类似与java中的`dbutils`，可以实现结果集的自动映射 [sqlx](https://github.com/jmoiron/sqlx)它没有传统ORM的繁琐，毕竟我们使用orm的初衷是为了减少开发成本，但是目前市面上的orm都有较高的学习成本维护成本而且并不能解决常见的sql问题。所以我断然不用

- [gin](https://github.com/gin-gonic/gin) 这是一个好用的轻量级的路由器

-  [iris](https://github.com/kataras/iris) 号称是世界上最快的路由 Gorilla is a web toolkit for the Go programming language that provides useful, composable packages for writing HTTP-based applications.

- [mux](https://github.com/gorilla/mux) 号称最好的分发器

- [martini](https://github.com/go-martini/martini/blob/master/translations/README_zh_cn.md) 号称是非侵入式的路由，支持扩展多个中间件，最简洁最`golang style`的框架，废除传统OOP这一点很吸引我,Go从语言层面就已经去除了臃肿的Class金字塔，但是太多的框架受此影响企图复辟。加之我对Express的偏爱，而此框架又受其影响，所以很值得研究.Martini is a powerful package for quickly writing modular web applications/services in Golang

-  [casbin](https://github.com/casbin/casbin) golang版的权限校验控制框架

- [negroni](https://github.com/urfave/negroni/blob/master/translations/README_zh_CN.md) 如果你喜欢用 Martini ，但又觉得它太魔幻，那么 Negroni 就是你很好的选择了。它是一个具备微型、非嵌入式、鼓励使用原生 net/http 库特征的`中间件`。
