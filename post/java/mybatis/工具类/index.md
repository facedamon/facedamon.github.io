
# 工具类

## SQL构造器

&emsp;&emsp;为了解决代码中sql拼接问题，mybatis中提供了一个sql工具类。使用这个工具类，我们可以很方便地在Java代码中动态构建SQL语句。

        @Test
        public void testSelectSql() {
            String sql = new SQL() {{
                SELECT("P.ID, P.USERNAME, P.USERNAME, P.PASSWORD, P.FULL_NAME");
                SELECT("P.LAST_NAME, P.CREATED_CN, P.UPDATED_ON");
                FROM("PERSON P");
                FROM("ACCOUNT A");
                INNER_JOIN("DEPARTMENT D on D.ID = P.DEPARTMENT_ID");
                INNER_JOIN("COMPANY C on D.COMPANY_ID = C.ID");
                WHERE("P.ID = A.ID");
                WHERE("P.FIRST_NAME like ?");
                OR();
                HAVING("P.FIRST_NAME like ?");
                ORDER_BY("P.ID");
                ORDER_BY("P.FULL_NAME");
            }}.toString();
            System.out.println(sql);
        }

        //output:
        SELECT P.ID, P.USERNAME, P.USERNAME, P.PASSWORD, P.FULL_NAME, P.LAST_NAME, P.CREATED_CN, P.UPDATED_ON
        FROM PERSON P, ACCOUNT A
        INNER JOIN DEPARTMENT D on D.ID = P.DEPARTMENT_ID
        INNER JOIN COMPANY C on D.COMPANY_ID = C.ID
        WHERE (P.ID = A.ID AND P.FIRST_NAME like ?) 
        OR ()
        HAVING (P.FIRST_NAME like ?)
        ORDER BY P.ID, P.FULL_NAME
        
&emsp;&emsp;在匿名子类的初始化代码块中，调用SELECT()、FROM()等方法构建SQL语句，这种方式能够很好地避免字符串拼接过程中缺少空格或者偶然间重复出现的AND关键字导致的SQL语句不正确。除了SELECT语句外，SQL工具类也可以用作构建UPDATE、INSERT等语句

&emsp;&emsp;在学习完SQL工具类的使用后，接下来我们简单地了解一下SQL工具类的实现源码。SQL继承至AbstractSQL类，只重写了该类的getSelf()方法，代码如下

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616128050.png)

&emsp;&emsp;所有的功能由AbstractSQL类完成，AbstractSQL类中维护了一个SQLStatement内部类的实例和一系列前面提到过的构造SQL语句的方法，例如SELECT()、UPDATE()等方法。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616129394.jpg)

SQLStatement内部类用于描述一个SQL语句，该类中通过StatementType确定SQL语句的类型。SQLStatement类中还维护了一系列的ArrayList属性，当调用SELECT()、UPDATE()等方法时，这些方法的参数内容会记录在这些ArrayList对象中。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616133619.jpg)

&emsp;&emsp;Appendable是一个追加模式的接口，AbstractSQL默认的实现是StringBuilder.

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616133913.jpg)

&emsp;&emsp;SQLStatement对象的sql()方法实现代码如下

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616134167.jpg)

&emsp;&emsp;如上面的代码所示，该方法中会判断SQL语句的类型，以UPDATE语句为例，会调用SQLStatement对象的updateSql()方法生成UPDATE语句。updateSql()方法代码如下：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616134333.jpg)

&emsp;&emsp;updateSql()方法中，最终会调用sqlCalause()方法完成SQL语句的拼接。sqlCalause()方法实现代码如下:

        /**
        *   @param builder 字符串构建对象
        *   @param keyword sql关键字
        *   @param parts 关键字子句内容
        *   @param open 关键字后开始字符
        *   @param close 关键字后结束字符
        *   @param conjunction sql链接关键字，通常为and或or
        */
        private void sqlClause(AbstractSQL.SafeAppendable builder, String keyword, List<String> parts, String open, String close, String conjunction) {
            if (!parts.isEmpty()) {
                if (!builder.isEmpty()) {
                    builder.append("\n");
                }

                builder.append(keyword);
                builder.append(" ");
                builder.append(open);
                String last = "________";
                int i = 0;
                // 拼接关键字后开始字符
                for(int n = parts.size(); i < n; ++i) {
                    String part = (String)parts.get(i);
                    // 如果sql关键字对应的子句内容不为OR或AND，则追加连接关键字
                    if (i > 0 && !part.equals(") \nAND (") && !part.equals(") \nOR (") && !last.equals(") \nAND (") && !last.equals(") \nOR (")) {
                        builder.append(conjunction);
                    }
                    // 追加子句内容
                    builder.append(part);
                    last = part;
                }

                builder.append(close);
            }
        }

## ScriptRunner

&emsp;&emsp;该工具类用于读取脚本文件中的SQL语句并执行，使用起来比较简单。下面是一个使用ScriptRunner执行SQL脚本的案例。

        @Test
        public void testScriptRunner() {
            try {
                Connection connection = DriverManager.getConnection("jdbc:mysql://xxx:yy/zz");
                ScriptRunner scriptRunner = new ScriptRunner(connection);
                scriptRunner.runScript(Resources.getResourceAsReader("create-table.sql"));
            } catch (Exception e) {
                e.printStackTrace();
            }
        }

&emsp;&emsp;ScriptRunner工具类的构造方法需要一个java.sql.Connection对象作为参数。创建ScriptRunner对象后，调用该对象的runScript()方法即可，该方法接收一个读取SQL脚本文件的Reader对象作为参数。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616135784.jpg)

&emsp;&emsp;ScriptRunner类的runScript()方法的逻辑比较清晰，具体做了以下几件事情：

1. 调用setAutoCommit()方法，根据autoCommit属性的值设置事务是否自动提交。
2. 判断sendFullScript属性值，如果值为true，则调用executeFullScript()方法一次性读取SQL脚本文件中的所有内容，然后调用JDBC中Statement对象的execute()方法一次性执行脚本中的所有SQL语句。
3. 如果sendFullScript属性值为false，则调用executeLineByLine()方法逐行读取SQL脚本文件，以分号作为一条SQL语句结束的标志，逐条执行SQL语句。

## SqlRunner

&emsp;&emsp;该类对JDBC做了很好的封装，结合SQL工具类，能够很方便地通过Java代码执行SQL语句并检索SQL执行结果。

|method|desc|
|---|---|
|closeConnection|关闭Connection对象|
|selectOne(String sql, Object... args)|执行select语句，sql语句可以使用占位符，该方法只返回一条记录，若查询结果行数不等于一，则会抛出SQLException|
|selectAll(String sql, Object... args)|返回多条记录，返回值为List\<Map\>，每个Map对应数据库中的一行记录|
|insert(String sql, Object... args)|插入一条记录|
|update(String sql, Object... args)|更新若干条记录|
|delete(String sql, Object... args)|删除若干条记录|
|run(String sql)|执行任意一条SQL语句|

        @Test
        public void testDelete() {
            try {
                SqlRunner sqlRunner = new SqlRunner(DriverManager.getConnection("jdbc:mysql://xxx:yy/zz"));
                String deleteUserSql = new SQL() {{
                    DELETE_FROM("user");
                    WHERE("id = ?");
                }}.toString();
                sqlRunner.delete(deleteUserSql, Integer.valueOf(1));
            } catch (Exception e) {
                e.printStackTrace();
            }
        }

        @Test
        public void testInsert() {
            try {
                SqlRunner sqlRunner = new SqlRunner(DriverManager.getConnection("jdbc:mysql://xxx:yy/zz"));
                String insertUserSql = new SQL() {{
                    INSERT_INTO("user");
                    VALUES("create_time, name, password, phone, nick_name",
                            "?,?,?,?,?");
                }}.toString();
                String createTime = LocalDateTime.now().format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));
                sqlRunner.insert(insertUserSql, createTime, "Jane", "test", "187", "J");
            }catch (Exception e) {
                e.printStackTrace();
            }
        }

&emsp;&emsp;熟悉了SqlRunner工具类的使用后，接下来我们了解一下SqlRunner工具类的实现。以selectAll()方法为例，该方法代码如下:

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616137403.png)

1. 调用Connection对象的prepareStatement()方法获取PreparedStatement对象，调用setParamters()方法为SQL中的参数占位符赋值
2. 调用PreparedStatement的executeQuery()方法执行查询操作
3. 调用getResults()方法，将ResultSet对象转换为List对象，其中List对象中的每个Map对象对应数据库中的一条记录

&emsp;&emsp;接下来我们来看一下setParameters()方法的实现.

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616137601.jpg)

1. 如果参数为null,则抛出异常
2. 如果参数为NULL枚举类型，则通过NULL枚举转换数据库类型
3. 根绝参数类型获取对应的TypeHandler对象，调用TypeHandler对象的setParamter方法为参数赋值

&emsp;&emsp;最后，我们再来看一下getResultSet()方法处理结果集的过程：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/mybatis/1616137987.jpg)

1. 获取ResultSetMetaData对象，通过ResultSetMetaData对象获取所有列名
2. 获取列的JDBC类型，根据类型获取TypeHandler对象
3. 遍历ResultSet对象，将ResultSet对象中的记录行转换为Map对象
4. 通过TypeHandler对象的getResult方法将JDBC类型转换为JAVA类型

## MetaObject

&emsp;&emsp;MetaObject是MyBatis中的反射工具类，该工具类在MyBatis源码中出现的频率非常高。使用MetaObject工具类，我们可以很优雅地获取和设置对象的属性值。

&emsp;&emsp;假设我们有两个实体类：User和Order，User类用于描述用户信息，Order类用于描述订单信息，一个用户可以有多笔订单，User类中通过一个List对象记录用户的订单信息。User和Order实体类的代码如下：

        @Data
        @AllArgsConstructor
        public class Order {
            String orderNo;
            String goodsName;
        }

        @Data
        @AllArgsConstructor
        public class User {
            List<Order> orders;
            String name;
            Integer age;
        }

        @Test
        public void testMetaObject() {
            List<Order> orders = new ArrayList<Order>() {
                {
                    add(new Order("order20210320", "mybatis 源码深度解析"));
                    add(new Order("order20210321", "tx mybatis 入门与进阶"));
                }
            };
            User user = new User(orders, "damon", 3);
            MetaObject metaObject = SystemMetaObject.forObject(user);
            //获取第一笔订单的商品名称
            System.out.println(metaObject.getValue("orders[0].goodsName"));
            //获取第二笔订单的商品名称
            System.out.println(metaObject.getValue("orders[1].goodsName"));
            //为属性设置值
            metaObject.setValue("orders[1].orderNo", "order20181103");
            //判断user对象是否有orderNo属性
            System.out.println(metaObject.hasGetter("orderNo"));
            //判断user对象是否有name属性
            System.out.println(metaObject.hasGetter("name"));
        }

        //output:
        mybatis 源码深度解析
        tx mybatis 入门与进阶
        false
        true

&emsp;&emsp;如上面的代码所示，我们创建了一个User对象并初始化User对象的属性值，接着调用SystemMetaObject类的forObject()静态方法创建一个与User对象关联的MetaObject对象。我们可以通过MetaObject对象的getValue()方法以表达式的方式获取User对象的属性值。我们还可以使用MetaObject对象的setValue()方法以表达式的方式为User对象的属性设置值。**当类的层级比较深时，使用MetaObject工具能够很方便地获取和设置对象的属性值**。除此之外，我们还可以使用MetaObject工具类的hasSetter()和hasGetter()方法通过名称判断对象是否有某个属性且该属性有对应的Getter/Setter方法。

## MetaClass

&emsp;&emsp;MetaClass是MyBatis中的反射工具类，与MetaOjbect不同的是，**MetaObject用于获取和设置对象的属性值，而MetaClass则用于获取类相关的信息**。例如，我们可以使用MetaClass判断某个类是否有默认构造方法，还可以判断类的属性是否有对应的Getter/Setter方法。

        @Test
        public void testMetaClass() {
            MetaClass metaClass = MetaClass.forClass(Order.class, new DefaultReflectorFactory());
            //获取所有有getter方法的属性名
            String[] getterNames = metaClass.getGetterNames();
            System.out.println(JSON.toJSONString(getterNames));
            //某属性是否有getter/setter方法
            System.out.println(String.format("%s:%b", "orderNo属性是否有Getter方法",metaClass.hasGetter("orderNo")));
            System.out.println(String.format("%s:%b", "orderNo属性是否有Setter方法",metaClass.hasSetter("orderNo")));
            System.out.println(String.format("%s:%s", "orderNo属性类型", metaClass.getGetterType("orderNo")));
            //获取属性Getter方法
            Invoker invoker = metaClass.getGetInvoker("orderNo");
            try {
                Object orderNo = invoker.invoke(new Order("order20210319", "mybatis"), null);
                System.out.println(orderNo);
            } catch (IllegalAccessException e) {
                e.printStackTrace();
            } catch (InvocationTargetException e) {
                e.printStackTrace();
            }
        }
        
        //output:
        ["orderNo","goodsName"]
        orderNo属性是否有Getter方法:true
        orderNo属性是否有Setter方法:true
        orderNo属性类型:class java.lang.String
        order20210319

&emsp;&emsp;我们通过MetaClass获取了Java类的基本信息，包括Java类中所有的Getter方法对应的属性名称。除此之外，我们还可以通过MetaClass获取Getter/Setter方法对应的Invoker对象，然后通过Invoker对象调用Getter/Setter方法。Invoker接口有2个不同的实现，分别为GetFieldInvoker、SetFieldInvoker，对应Java类的Getter方法、Setter方法。

## ObjectFactory

&emsp;&emsp;ObjectFactory是MyBatis中的对象工厂，MyBatis每次创建Mapper映射结果对象的新实例时，都会使用一个对象工厂（ObjectFactory）实例来完成。ObjectFactory接口只有一个默认的实现，即DefaultObjectFactory，默认的对象工厂需要做的仅仅是实例化目标类，要么通过默认构造方法，要么在参数映射存在的时候通过参数构造方法来实例化。

        @Test
        public void testObjectFactory() {
            ObjectFactory objectFactory = new DefaultObjectFactory();
            List<Integer> list = objectFactory.create(List.class);
            Map<String, String> map = objectFactory.create(Map.class);
            list.addAll(Arrays.asList(1, 2, 3));
            map.put("test", "test");
            System.out.println(list);
            System.out.println(map);
        }

        //output:
        [1, 2, 3]
        {test=test}

&emsp;&emsp;DefaultObjectFactory实现类支持通过接口的方式创建对象，例如当我们指定创建java.util.List实例时，实际上创建的是java.util.ArrayList对象。List、Map、Set接口对应的实现分别为ArrayList、HashMap、HashSet。

> 关于工厂模式请参考 [abstract factory](http://facedamon.github.io/post/design-patterns/%E5%88%9B%E5%BB%BA/%E6%8A%BD%E8%B1%A1%E5%B7%A5%E5%8E%82%E6%A8%A1%E5%BC%8F/)

## ProxyFactory

&emsp;&emsp;ProxyFactory是MyBatis中的代理工厂，主要用于创建动态代理对象，ProxyFactory接口有两个不同的实现，分别为CglibProxyFactory和JavassistProxyFactory。从实现类的名称可以看出，MyBatis支持两种动态代理策略，分别为Cglib和Javassist动态代理。ProxyFactory主要用于实现MyBatis的懒加载功能。当开启懒加载后，MyBatis创建Mapper映射结果对象后，会通过ProxyFactory创建映射结果对象的代理对象。当我们调用代理对象的Getter方法获取数据时，会执行CglibProxyFactory或JavassistProxyFactory中定义的拦截逻辑，然后执行一次额外的查询。

        @Test
        public void testProxyFactory() {
            //创建ProxyFactory对象
            ProxyFactory proxyFactory = new JavassistProxyFactory();
            Order order = new Order("gn20210319", "mybatis");
            ObjectFactory objectFactory = new DefaultObjectFactory();
            //调用ProxyFactory对象的createProxy方法创建代理对象
            Object proxyOrder = proxyFactory.createProxy(order,
                    Mockito.mock(ResultLoaderMap.class),
                    Mockito.mock(Configuration.class),
                    objectFactory,
                    Arrays.asList(String.class, String.class),
                    Arrays.asList(order.getOrderNo(), order.getGoodsName()));
            System.out.println(proxyOrder.getClass());
            System.out.println(((Order)proxyOrder).getOrderNo());
        }

        //output:
        class org.facedamon.Order_$$_jvst745_0
        gn20210319

&emsp;&emsp;我们创建了一个Order对象，然后通过JavassistProxyFactory实例创建了一个Order对象的动态代理对象。代理对象创建完毕后，会把原始对象的属性复制到代理对象中，调用代理对象的Getter方法获取属性值时，会执行JavassistProxyFactory中的拦截逻辑。

> 关于代理模式请参考 [proxy pattern](http://facedamon.github.io/post/design-patterns/%E7%BB%93%E6%9E%84/%E4%BB%A3%E7%90%86%E6%A8%A1%E5%BC%8F/)