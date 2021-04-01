
# 摘要

&emsp;&emsp;SqlSession是MyBatis提供的面向用户的操作数据库API。那么MyBatis底层是如何工作的呢？为了解开MyBatis的神秘面纱，我们需要了解一下MyBatis的其他几个比较核心的组件及这些组件的作用。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/epub_27563388_132.jpg)

- **Configuration**：用于描述MyBatis的主配置信息，其他组件需要获取配置信息时，直接通过Configuration对象获取。除此之外，MyBatis在应用启动时，将Mapper配置信息、类型别名、TypeHandler等注册到Configuration组件中，其他组件需要这些信息时，也可以从Configuration对象中获取。
- **MappedStatement**：MappedStatement用于描述Mapper中的SQL配置信息，是对Mapper XML配置文件中<select|update|delete|insert>等标签或者@Select/@Update等注解配置信息的封装
- **SqlSession**：SqlSession是MyBatis提供的面向用户的API，表示和数据库交互时的会话对象，用于完成数据库的增删改查功能。SqlSession是Executor组件的外观，目的是对外提供易于理解和使用的数据库操作接口。
- **Executor**：Executor是MyBatis的SQL执行器，MyBatis中对数据库所有的增删改查操作都是由Executor组件完成的。
- **StatementHandler**：StatementHandler封装了对JDBC Statement对象的操作，比如为Statement对象设置参数，调用Statement接口提供的方法与数据库交互，等等。
- **ParameterHandler**：当MyBatis框架使用的Statement类型为CallableStatement和PreparedStatement时，ParameterHandler用于为Statement对象参数占位符设置值。
- **ResultSetHandler**：ResultSetHandler封装了对JDBC中的ResultSet对象操作，当执行SQL类型为SELECT语句时，ResultSetHandler用于将查询结果转换成Java对象。
- **TypeHandler**：TypeHandler是MyBatis中的类型处理器，用于处理Java类型与JDBC类型之间的映射。它的作用主要体现在能够根据Java类型调用PreparedStatement或CallableStatement对象对应的setXXX()方法为Statement对象设置值，而且能够根据Java类型调用ResultSet对象对应的getXXX()获取SQL执行结果。

&emsp;&emsp;`SqlSession`组件，它是用户层面的API。实际上`SqlSession`是`Executor`组件的外观，目的是为用户提供更友好的数据库操作接口，这是设计模式中[外观模式](http://facedamon.github.io/post/design-patterns/%E7%BB%93%E6%9E%84/%E5%A4%96%E8%A7%82%E6%A8%A1%E5%BC%8F)的典型应用。真正执行SQL操作的是Executor组件，Executor可以理解为SQL执行器，它会使用`StatementHandler`组件对JDBC的`Statement`对象进行操作。当Statement类型为`CallableStatement`和`PreparedStatement`时，会通过`ParameterHandler`组件为参数占位符赋值。ParameterHandler组件中会根据Java类型找到对应的`TypeHandler`对象，TypeHandler中会通过Statement对象提供的setXXX()方法（例如setString()方法）为Statement对象中的参数占位符设置值。`StatementHandler`组件使用JDBC中的Statement对象与数据库完成交互后，当SQL语句类型为SELECT时，MyBatis通过`ResultSetHandler`组件从Statement对象中获取ResultSet对象，然后将ResultSet对象转换为Java对象。

## Configuration

&emsp;&emsp;MyBatis框架的配置信息有两种，**一种是配置MyBatis框架属性的主配置文件；另一种是配置执行SQL语句的Mapper配置文件**。Configuration的作用是描述MyBatis主配置文件的信息。Configuration类中定义了一系列的属性用来控制MyBatis运行时的行为，这些属性的值可以在MyBatis主配置文件中通过\<setting\>标签指定。

|属性|作用|
|---|---|
|cacheEnabled|是都开启Mapper二级缓存|
|lazyLoadingEnabled|延迟加载的全局开关。当开启时，所有关联对象都会延迟加载|
|useColumnLabel|使用列标签代替列名|
|useGeneratedKeys|允许JDBC支持自动生成主键|
|...|...|

&emsp;&emsp;Configuration除了提供了表4-1中的属性控制MyBatis的行为外，还作为容器存放TypeHandler（类型处理器）、TypeAlias（类型别名）、Mapper接口及Mapper SQL配置信息。这些信息在MyBatis框架启动时注册到Configuration组件中。

- **mapperRegistry**: 用于注册Mapper接口信息，建立Mapper接口的Class对象和MapperProxyFactory对象之间的关系，其中MapperProxyFactory对象用于创建Mapper动态代理对象。
- **interceptorChain**: 用于注册MyBatis插件信息，MyBatis插件实际上就是一个拦截器。
- **typeHandlerRegistry**：用于注册所有的TypeHandler，并建立JDBC类型、JDBC类型与TypeHandler之间的对应关系
- **typeAliasRegistry**：用于注册所有的类型别名。
- **languageRegistry**：用于注册LanguageDriver，LanguageDriver用于解析SQL配置，将配置信息转换为SqlSource对象。
- **mappedStatements**：MappedStatement对象描述\<insert|select|update|delete\>等标签或者通过@Select、@Delete、@Update、@Insert等注解配置的SQL信息。MyBatis将所有的MappedStatement对象注册到该属性中，其中**Key为Mapper的Id，Value为MappedStatement对象**。
- **caches**：用于注册Mapper中配置的所有缓存信息，**其中Key为Cache的Id，也就是Mapper的命名空间，Value为Cache对象**。
- **resultMaps**：用于注册Mapper配置文件中通过\<resultMap\>标签配置的ResultMap信息，ResultMap用于建立Java实体属性与数据库字段之间的映射关系，其中Key为ResultMap的Id，该Id是由Mapper命名空间和\<resultMap\>标签的id属性组成的，Value为解析\<resultMap\>标签后得到的ResultMap对象。
- **parameterMaps**：用于注册Mapper中通过<parameterMap>标签注册的参数映射信息。Key为ParameterMap的Id，由Mapper命名空间和<parameterMap>标签的id属性构成，Value为解析<parameterMap>标签后得到的ParameterMap对象。
- **keyGenerators**：用于注册KeyGenerator，KeyGenerator是MyBatis的主键生成器，MyBatis中提供了3种KeyGenerator，即Jdbc3KeyGenerator（数据库自增主键）、NoKeyGenerator（无自增主键）、SelectKeyGenerator（通过select语句查询自增主键，例如oracle的sequence）。
- **loadedResources**：用于注册所有Mapper XML配置文件路径。
- **sqlFragments**：用于注册Mapper中通过<sql>标签配置的SQL片段，Key为SQL片段的Id，Value为MyBatis封装的表示XML节点的XNode对象。
- **incompleteStatements**：用于注册解析出现异常的XMLStatementBuilder对象。
- **incompleteCacheRefs**：用于注册解析出现异常的CacheRefResolver对象。
- **incompleteResultMaps**：用于注册解析出现异常的ResultMapResolver对象。
- **incompleteMethods**：用于注册解析出现异常的MethodResolver对象。

&emsp;&emsp;MyBatis框架启动时，会对所有的配置信息进行解析，然后将解析后的内容注册到Configuration对象的这些属性中。除此之外，Configuration组件还作为Executor、StatementHandler、ResultSetHandler、ParameterHandler组件的工厂类，用于创建这些组件的实例。Configuration类中提供了这些组件的工厂方法，这些工厂方法签名如下:

        //ParameterHandler组件工厂方法
        public ParameterHandler newParameterHandler(
                    MappedStatement mappedStatement, 
                    Object parameterObject, 
                    BoundSql boundSql)
        ;
        //ResultSetHandler组件工厂方法
        public ResultsetHandler newResultSetHandler(
                    Executor executor, 
                    MappedStatement mappedStatement, 
                    RowBounds rowBounds, 
                    ParameterHandler parameterHandler, 
                    ResultHandler resultHandler, 
                    BoundSql boundSql)
        ;
        //StatementHandler组件工厂方法
        public StatementHandler newStatementHandler(
                    Executor executor, 
                    MappedStatement mappedStatement, 
                    RowBounds rowBounds, 
                    ParameterHandler parameterHandler, 
                    ResultHandler resultHandler, 
                    BoundSql boundSql)
        ;
        //Executor组件工厂方法
        public Executor newExecutor(
                    Transaction transaction, 
                    ExecutorType executorType)
        ;

&emsp;&emsp;这些工厂方法会根据MyBatis不同的配置创建对应的实现类。例如，Executor组件有4种不同的实现，分别为`BatchExecutor`、`ReuseExecutor`、`SimpleExecutor`、`CachingExecutor`，当defaultExecutorType的参数值为REUSE时，newExecutor()方法返回的是ReuseExecutor实例，当参数值为SIMPLE时，返回的是SimpleExecutor实例，这是典型的[工厂方法模式](http://facedamon.github.io/post/design-patterns/%E5%88%9B%E5%BB%BA/%E5%B7%A5%E5%8E%82%E6%A8%A1%E5%BC%8F)的应用。

## Executor

&emsp;&emsp;Executor接口中定义了对数据库的增删改查方法，其中query()和queryCursor()方法用于执行查询操作，update()方法用于执行插入、删除、修改操作。Executor接口有几种不同的实现类:

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/epub_27563388_141.jpg)

&emsp;&emsp;MyBatis提供了3种不同的Executor，分别为`SimpleExecutor`、`ResueExecutor`、`BatchExecutor`，这些Executor都继承至`BaseExecutor`，BaseExecutor中定义的方法的执行流程及通用的处理逻辑，具体的方法由子类来实现，是典型的[模板方法模式](http://facedamon.github.io/post/design-patterns/%E8%A1%8C%E4%B8%BA/%E6%A8%A1%E6%9D%BF%E6%A8%A1%E5%BC%8F)的应用。

&emsp;&emsp;`SimpleExecutor`是基础的Executor，能够完成基本的增删改查操作，`ResueExecutor`对JDBC中的Statement对象做了缓存，当执行相同的SQL语句时，直接从缓存中取出Statement对象进行复用，避免了频繁创建和销毁Statement对象，从而提升系统性能，这是[享元思想](http://facedamon.github.io/post/design-patterns/%E8%A1%8C%E4%B8%BA/%E6%A8%A1%E6%9D%BF%E6%A8%A1%E5%BC%8F)的应用。BatchExecutor则会对调用同一个Mapper执行的update、insert和delete操作，调用Statement对象的批量操作功能。

&emsp;&emsp;当MyBatis开启了二级缓存功能时，会使用`CachingExecutor`对SimpleExecutor、ResueExecutor、BatchExecutor进行装饰，为查询操作增加二级缓存功能，这是[装饰器模式](http://facedamon.github.io/post/design-patterns/%E7%BB%93%E6%9E%84/%E8%A3%85%E9%A5%B0%E5%99%A8%E6%A8%A1%E5%BC%8F/)的应用