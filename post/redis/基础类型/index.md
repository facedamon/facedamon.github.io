
# 基础类型

|string|list|hash|set|zset|
|:-:|:-:|:-:|:-:|:-:

## string

&emsp;&emsp;一个字符串类型允许存储的数据的最大容量是512MB。

|command|return|
|:-:|:-:|
|set key value|返回ok|
|get key|当key不存在是返回empty|
|incr key|返回递增后的值(整数)，线程安全|
|decr key|返回递减后的值(整数)，线程安全|


### 文章访问量统计

&emsp;&emsp;博客常见功能是统计文章的访问量，我们可以为每篇文章使用一个名为`post:文章ID:page.view`的键来记录文章的访问量，每次访问文章的时候使用`INCR`递增键值。

&emsp;&emsp;Redis对于键的命名并没有强制的要求，但是比较好的实践是用`对象类型:对象ID:对象属性`来命名一个键，如使用键`user:1:friends`来存储ID为1的用户的好友列表。

### 生成自增ID

&emsp;&emsp;再关系型数据库中我们可以通过设置字段属性为`AUTO_INCREMENT`来实现每增加一条记录自动为其生成一个唯一的递增ID，而在Redis中可以通过另一种模式来实现：对于每一类对象使用名为`对象类型:count`的键(如`user:count`)来存储当前类型对象的数量，没增加一个新对象时都使用INCR命令递增该键的值。

### 存储文章数据

&emsp;&emsp;由于每个字符串类型键只能存储一个字符串，而一篇博客文章是由标题、正文、作者与发布时间等多个元素构成。为了存储这些元素，我们采用序列化函数将它们序列化成一个字符串，或者使用二进制存储(MessagePack)效率更高，占用空间更小。

|command|return|
|:-:|:-:|
|incrby key increment|可以指定一次增加的值, 返回相加后的值|
|decrby key decrement|可以指定一次递减的值，返回相减后的值|
|incrbyfloat key increment|递增一个双精度浮点数，返回相加后的值|
|append key value|返回值是追加后字符串的总长度|
|strlen key|返回键值的长度， 不存在则返回0|
|mset k v [k v ...]|设置多个键值对(mset k1 v1 k2 v2) 返回ok|
|mget k [k ...]|获取对各键值对(mget k1 k2) 返回多个值|

### 位操作

|位操作|
|:-:|
|setbit key offset value|
|getbit key offset|
|bitcount key [start] [end]|
|bitop operation destkey key [key...]

&emsp;&emsp;`set foo bar`, bar转换成ASCII码`98 97 114`, 转换成二进制 `01100010 01100001 01110010`

&emsp;&emsp;`getbit`可以获得一个字符串类型指定位置的二进制位的值，索引从0开始。`get bit foo 6`值为1。

&emsp;&emsp;`setbit`可以设置字符串类型指定位置的二进制位的值，**返回值是该位置的旧值**。

    setbit foo 6 0
    setbit foo 7 1
    get foo
    aar

&emsp;&emsp;`bitcount`可以获得字符串类型键中值是1的二进制位个数， `bitcount foo`返回10。如果我们只希望统计前两个字节：`bitcount foo 0 1`返回6    

&emsp;&emsp;`bitop`可以对多个字符串类型键进行位运算，并将结果存储在destkey参数指定的键中。`bitop`支持的运算操作有`and or xor not`:

    set foo1 bar
    ok
    set foo2 aar
    ok
    bitop or res foo1 foo2
    8
    get res
    " car"

&emsp;&emsp;**位操作常用来存储布尔值**。比如某网站的每个用户都有一个递增的整数ID，如果使用一个字符串类型键配合位操作来记录每个用户的性别(用户ID作为索引，二进制位值1和0表示男性和女性)，那么记录100万个用户的性别只需占用100KB的空间，而且二进制读取性能非常块。

## hash

&emsp;&emsp;设想在博客文章列表页，每个文章只显示标题部分，可是使用上述方法，若想取得文章的标题，必须把整个文章数据字符串取出来反序列化，而其中占用空间最大的文章内容部分却是不需要的。在修改数据时，当你只想更改文章的标题时不得不把整个文章数据字符串更新一遍。

&emsp;&emsp;**散列类型适合存储对象：**使用对象类别和ID构成键，使用字段表示对象的属性，而字段值则存储属性值。例如要存储ID为2的汽车对象，可以分别使用名为color、name和price的3个字段来存储该辆汽车的颜色、名称和价格：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/redis/epub_731187_28.jpg)

&emsp;&emsp;回想在关系数据库中如果要存储汽车对象，存储结构如下图：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/redis/epub_731187_29.jpg)

&emsp;&emsp;数据是以二维表的行式存储的，这就要求**所有的记录都拥有同样的属性**,无法单独为某一条记录增减属性。如果想为ID为1的汽车增加上产日期属性，就需要把数据表更改为如下图所示的结构：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/epub_731187_30.jpg)  

&emsp;&emsp;对于ID为2和3的两条记录而言date字段是冗余的。可想而知当不同的记录需要不同的属性时，表的字段数量会越来越多以至于难以维护。而且当使用ORM将关系数据库中的对象实体映射成程序中的实体时，修改表的结构往往意味着要中断服务。在关系型数据库中存储这种半结构化数据还需要额外的表才行。

### 命令

|Hash|
|:-:|
|hset key field value|
|hget key field|
|hmset key field value [field value ...]|
|hmget key field [field ...]
|hgetall key|

    hset car:1 price 500
    (integer) 1
    hset car:1 name BMW
    (integer) 1
    hget car:1 name
    "BMW"

&emsp;&emsp;`hset`的方便之处在于不区分插入和更新操作，这意味修改数据时不用事先判断字段是否存在来决定要执行的是update还是insert。**当执行的是插入操作时hset返回1，当执行的是更新操作时hset返回0**。

&emsp;&emsp;当需要同时设置多个字段时，可以使用`HMSET`命令。

    hmset key field1 value1 field2 value2
    # hmset car:1 price 500 name bmw

&emsp;&emsp;相应地，`hmget`命令可以同时获得多个字段的值：

    hmget key field1 field2
    # hmget car:1 price name

&emsp;&emsp;如果像获取键中所有字段和字段值却不知键中有哪些字段时，可以使用`hgetall`  :

    hgetall car:1
    1) "price"
    2) "500"
    3) "name"
    4) "BMW"

- 判断字段是否存在

|hexists key field|存在返回，否则0|
|:-:|:-:|

    hexists car:1 model
    0
    hset car:1 model C200
    1
    hexists car:1 model
    1

- 当字段不存在时赋值(线程安全)

|hsetnx key field value|与hset类似，区别在于如果字段存在，hsetnx不执行任何操作
|:-:|:-:|

- 增加数字(线程安全)

|hincrby key field increment|返回增值后的字段值
|:-:|:-:|

    hincrby person:1 score 60
    60

- 删除字段(线程安全)

|hdel key field [field ...]|返回被删除的字段个数|
|:-:|:-:|

    hdel car:1 price
    1
    hdel car:1 price
    0

### 实践

1. 存储文章数据

&emsp;&emsp;上节介绍了可以将文章对象序列化后使用一个字符串类型键存储，可是这种方法无法提供对单个字段的原子读写操作支持，从而产生竞态条件。此外，即使只需要文章标题，程序也不得不将包括文章内容在内的所有文章数据取出并反序列化，比较消耗资源。

&emsp;&emsp;此外，还有一种方法时组合使用多个字符串类型键来存储文章的数据，如下图：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/redis/epub_731187_34.jpg)

&emsp;&emsp;使用这种方法的好处在于无论获取还是修改文章数据，都可以只针对某一个属性进行操作每十份方便，缺点是，对象过于松散且随着对象属性增加，键值也会增加，内存消耗非常可观且不易维护。

&emsp;&emsp;还有一种散列类型存储结构如下图：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/redis/epub_731187_35.jpg)

&emsp;&emsp;从上图可以看出使用散列类型存储文章数据比上上图的方法看起来更加直观也更容易维护(比如可以使用`hgetall`获得一个对象的所有字段，删除一个对象时只需要删除一个键)，另外存储同样的庶几乎散列类型往往比字符串类型更加节约空间。

2. 存储文章缩略名

&emsp;&emsp;我们可以使用一个散列类型的键slug.to.id来存储文章缩略名和ID之间的映射关系。其中字段用来记录缩略名，字段值用来记录缩略名对应的ID。**这样就可以使用`hexists`来判断缩略名是否存在，使用`hget`获得缩略名对应的文章ID**。

3. 伪代码

        #获取文章ID
        postId := incr posts:count
        #判断用户输入的slug缩略名是否可用
        is := hsetnx slug.to.id, slug, postId
        if is == 0 {exit}
        hmset post:postId title, $title, content, $content, slug, $slug

        #当用户访问时
        postId := hget slug.to.id slug
        if postId == 0 {exit}

        post := hgetall post:postId
        fmt.Printf("post title: %s\n", post.title)

### 命令拾遗

1. 只获取字段名或字段值

        hkeys key
        hvals key

2. 获得字段数量

        hlen key

## list

&emsp;&emsp;正当我们信心满满地写着文章列表页地代码时，一个很重要地问题阻碍了开发。我们当前获取文章列表(分页)的逻辑是：

1. 读取posts:count键获得博客中最大的文章ID
2. 规定每页显示10片文章，按照ID倒序排列，如第n页的文章ID范围是从 最大的文章`ID - (n-1) * 10 `到 `max(最大文章ID - n * 10 +1, 1)`
3. 对每个ID使用`hmget`来获得文章数据

&emsp;&emsp;可是这种方式要求用户不能删除文章以保证ID的连续性，否则就必须在程序中使用`exists`判断某个ID是否存在，不存是则跳过。都不得不从最大的文章ID开始遍历来获得当前页面应该显示哪些文章。

### 简介

&emsp;&emsp;列表类型可以存储一个有序的字符串列表，常用的操作是向列表两端添加元素，或者获得列表的某一个片段。内部采用双向链表实现，所以向列表两端添加元素的时间复杂度为O(1)。

### 场景

&emsp;&emsp;我们都知道，链表在首尾两端操作数据实践复杂度为O(1)，但是在遍历上很慢，每次都要从头节点往后遍历，因为它没有下标，只有指针。但是，正式由于这种特性，一些场景非常适合列表。

&emsp;&emsp;如社交网站的新鲜事，我们只关心最新的内容，使用列表类型存储，即使新鲜事的总数达到几千万个，获取其中最新的100条数据也是极快的。

&emsp;&emsp;借助列表类型，Redis还可以作为队列使用。

|命令|描述|返回|
|:-:|:-:|:-:|
|lpush key value [value ...]|列表左侧添加元素|增加后列表的长度|
|rpush key value [value ...]|列表右侧添加元素|增加后列表的长度|
|lpop key|列表左侧弹出元素|弹出元素的值|
|rpop key|列表右侧弹出元素|弹出元素的值|
|llen key|获取列表中元素的个数|返回列表中元素的个数|
|lrange key start end|获得列表片段|返回[start,end]包括边界的元素|
|lrem key count value|删除列表中指定的值|实际删除的元素个数|

    lpush num 2 3 1
    (integer) 3
    # [1 3 2]
    rpush num 0 -1
    (integer) 5
    # [1 3 2 0 -1]
    lop num
    "1"
    rpop num
    "-1"
    # [3 2 0]
    lrange num 0 2
    "3"
    "0"

&emsp;&emsp;综合上面提到的4个命令可以使用列表类型来模拟栈和队列的操作：如果像把列表当作栈，则搭配使用`lpush`和`lpop`或`rpush`和`rpop`，如果想当成队列，则搭配使用`lpush`和`rpop`或`rpush`和`lpop`。

&emsp;&emsp;`llen`功能类似sql语句`select count(1) from table_name`，但是`llen`的时间复杂度为O(1)，使用时Redis会直接读取现成的值，而不需要像关系型数据库那样需要遍历一遍数据表来统计条目数量。

&emsp;&emsp;`lrange`支持负索引，表示从右边开始计算序数，如"-1"表示最右边第一个元素，"-2"表示最右边第二个元素，以此类推：

    lrange num -2 -1
    "2"
    "0"

&emsp;&emsp;显然，`lrange num 0 -1`可以获取列表中的所有元素。还有一些特殊情况：

1. 如果start的索引位置比end的索引位置靠后，则会返回空列表
2. 如果end大于实际的索引范围，则会返回列表最右边的元素

&emsp;&emsp;`lrem`会删除列表中前`count`个值为`value`的元素，返回值是实际删除的元素个数。根据`count`值的不同，`lrem`的执行方式会略有差异：

- 当**count > 0**时`lrem`会从列表左边开始删除前`count`个值为`value`的元素
- 当**count < 0**时`lrem`会从列表右边开始删除前`|count|`个值为`value`的元素
- 当**count = 0**时`lrem`会删除所有值为`value`的元素

        rpush num 1
        (inetger) 4
        # [3 2 0 1]
        lrange num 0-1
        "3"
        "2"
        "0"
        "1"
        lrem num -1 2
        (integer) 1
        lrange num 0-1
        "3"
        "0"
        "1"

### 实践

&emsp;&emsp;为了解决之前遇到的分页问题，我们使用列表类型键`posts:list`记录文章ID列表。当发布新文章时使用`lpush`把新文章的ID加入到这个列表中，另外删除文章时也要把列表中的文章ID删除，像这样: **lrem posts:list 1 要删除的文章ID**

    # 分页伪代码
    postsPerPage := 10
    start := (currentPage - 1) * postsPerPage
    end := currentPage * postPerPage - 1
    postsId := lrange posts:list start end

    for id := range postId {
        post := hgetall post:id
        fmt.Printf("title: %s\n", post.title)
    }

&emsp;&emsp;这样显示的文章列表是根据加入列表的顺序倒序的，如果想让最旧的文章显示在前面，可以使用`lrange`获取需要的部分并在客户端中将顺序反转显示。

&emsp;&emsp;分页问题就此解决了，美中不足的是一点的散列类型没有类似字符串类型的`mget`那样可以通过一条命令同时获得多个键的键值的版本，所以对于每个文章ID修需要请求一次数据库，之后我们会介绍使用管道和脚本来优化这个问题。

&emsp;&emsp;另外使用列表类型存储文章ID列表有以下两个问题：

1. 文章的发布时间不宜修改：修改文章的发布时间不仅要修改post:文章ID中的time字段，还需要按照实际的发布时间重新排列posts:list中的元素顺序。
2. 当文章数量较多时访问中间的页面性能较差：列表类型时通过链表实现的，所以当列表元素非常多时访问中间的元素效率并不高。

### 命令拾遗

|命令|描述|返回值|
|:-:|:-:|:-:|
|lindex key index|获取指定索引的元素|返回指定索引的元素，从0开始|
|lset key index value|设置指定索引的元素|返回OK|
|ltrim key start end|删除[start, end]索引范围之外的所有元素|OK|
|linsert key before\|after pivot value|指定位置插入元素|列表元素个数|
|rpoplpush source dest|将元素从一个列表转到另一个列表|

    lindex num 0
    "2"
    #如果index是负数则表示从右边开始计算的索引，最后边元素的索引是-1
    lindex num -1
    "0"
    lset num 1 7
    OK
    lindex num 1
    "7"
    lrange num 0-1
    "1"
    "2"
    "7"
    "3"
    ltrim num 1 2
    OK
    lrange num 0-1
    "2"
    "7"

&emsp;&emsp;`linsert`命令首先在列表中从左到右查找值为`pivot`的元素，然后根据第二个参数`before`还是`after`来决定将`value`插入到该元素的前面还是后面。

    lrange num 0-1
    "2"
    "7"
    "0"
    linsert num after 7 3
    (integer) 4
    lrange num 0-1
    "2"
    "7"
    "3"
    "0"
    linsert num before 2 1
    "1"
    "2"
    "7"
    "3"
    "0"

&emsp;&emsp;`rpoplpush`是个很有意思的命令，从名字就可以看出它的功能：限制性`rpop`命令再执行`lpush`命令。`rpoplpush`命令会先从`source`列表类型键的右边弹出一个元素，然后将其加入到`dest`列表类型键的左边，并返回这个元素的值，整个过程是原子的。

&emsp;&emsp;当把列表类型作为队列使用时，`rpoplpush`命令可以很直观地在多个队列中传递数据，当`source`和`dest`相同时，`rpoplpush`命令会不断地将队尾地元素移到队首，借助这个特性我们可以实现一个网站监控系统：使用一个队列存储需要监控地网址，然后监控程序不断地使用`rpoplpush`命令循环取出一个网址来测试可用性。

## set

&emsp;&emsp;到目前位置，博客系统大部分功能还算比较完善，还有一个问题：如何在Redis中存储文章标签。使用散列类型或列表类型存储，虽然都能实现，但还是颇有不妥。Redis有一中数据类型很适合存储文章的标签，它就是集合类型。

### 简介

&emsp;&emsp;集合的概念高中的数学课就学习过。在集合中的每个元素都是不同的，且没有顺序。一个集合类型键可以存储至多2^32 - 1个字符串。集合类型和列表类型有相似之处，但很容易将它们区分开来：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/redis/epub_731187_42.jpg)

&emsp;&emsp;集合类型的常用操作是向集合中加入或删除元素、判断某个元素是否存在，由于集合类型在Redis内部是使用值为空的散列表(hash table)实现的，所以这些操作的时间复杂度都是O(1)。

### 命令

|命令|描述|返回|
|:-:|:-:|:-:|
|sadd key member [member...]|向集合中添加一个或多个元素|成功加入的元素数量|
|srem key member [member...]|从集合中删除一个或多个元素|删除成功的个数|
|smembers key|获得集合中的所有元素|集合中的所有元素|
|sismember key member|判断元素是否在集合中|1/0|
|sdiff key [key ...]|差集|差集结果|
|sinter key [key ...]|交集|交集结果|
|sunion key [key ...]|并集|并集结果|

&emsp;&emsp;`sadd`如果键不存在则会自动创建，因为在一个集合中不能有相同的元素，所以如果要加入的元素已经存在于集合中就会忽略这个元素。

    sadd letters a
    (integer) 1
    sadd letters a b c
    (integer) 2
    srem letters c d
    (integer) 1
    smembers letters
    "b"
    "a"
    sismember letters a
    (integer) 1
    # 布隆过滤器 O(1)
    sismember letters d
    (integer) 0


1. `sdiff`命令用来对多个集合执行差集运算。集合A与集合B的差集表示为A-B，代表所有属于A且不属于B的元素构成的集合。$$A-B= \lbrace x|x \subseteq A and x \notin  B \rbrace$$

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/redis/epub_731187_43.jpg)

$$\lbrace 1, 2, 3 \rbrace - \lbrace 2, 3, 4 \rbrace = \lbrace 1 \rbrace$$
$$\lbrace 2, 3, 4 \rbrace - \lbrace 1, 2, 3 \rbrace = \lbrace 4 \rbrace$$

    sadd setA 1 2 3
    (integer) 3
    sadd setB 2 3 4
    (integer) 3
    sdiff setA setB
    "1"
    sdiff setB setA
    "4"

2. `sinter`用来对多个集合执行交集运算。集合A与集合B的交集表示为A∩B，代表所有署于A且署于B的元素构成的集合。$$A∩B= \lbrace x|x \subseteq A and x \subseteq B \rbrace$$

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/redis/epub_731187_44.jpg)


$$\lbrace 1, 2, 3 \rbrace ∩ \lbrace 2, 3, 4 \rbrace = \lbrace 2, 3 \rbrace$$

    sinter setA setB
    "2"
    "3"
    #支持多个键
    sinter setA setB setC
    "2"
    "3"

3. `sunion`用来对多个集合执行并集运算。集合A与集合B的并集表示为A∪B，代表所有署于A或属于B的元素构成的集合。$$A∪B= \lbrace x|x \subseteq A or x \subseteq B \rbrace$$

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/redis/epub_731187_45.jpg)

$$\lbrace 1, 2, 3 \rbrace ∪ \lbrace 2, 3, 4 \rbrace = \lbrace 1, 2, 3, 4 \rbrace$$

    sunion setA setB
    "1"
    "2"
    "3"
    "4"
    # 支持多个键
    sunion setA setB setC
    "1"
    "2"
    "3"
    "4"

### 实践

1. 存储标签

&emsp;&emsp;对每篇文章使用名为 "post:ID:tags"的键存储该片文章的标签，伪代码：

    #给ID为42的文章添加标签
    sadd post:42:tags 技术文章，golang，redis
    # 删除标签
    srem post:42:tags 技术文章
    # 显示所有的标签
    tags := smembers post:42:tags
    fmt.Printf("tags=%s\n", tags)

&emsp;&emsp;使用集合类型键存储标签适合需要单独增加或删除标签的场合。如在WordPress博客程序中无论是添加还是删除标签都是针对单个标签的，可以直观地使用`sadd`和`srem`命令完成操作。

&emsp;&emsp;另一方面，有些地方需要用户直接设置所有标签后一起上传修改，用户编辑自己的爱好后提交，程序直接覆盖原来的标签数据，整个过程没有针对单个标签的操作，并利用到集合类型的优势，所以此时也可以直接使用字符串类型键存储标签数据。

&emsp;&emsp;在实践中的差别是想说明对于Redis存储方式的选择并没有绝对的规则，在一些特定的场合下散列类型甚至字符串类型可能更适合。

2. 通过标签搜索文章

&emsp;&emsp;有时我们还需要列出某个标签下的所有文章，甚至需要获得同时属于某几个标签的文章列表，这种需求在传统关系型数据库中实现起来比较复杂。

&emsp;&emsp;假设，现在有三张表，即posts、tags和posts_tags，分别存储文章数据、标签、文章与标签的对应关系。

|字段名|说明|
|:-:|:-:|
|post_id|文章ID|
|post_title|文章标题|

|字段名|说明|
|:-:|:-:|
|tag_id|标签ID|
|tag_name|标签名称|

|字段名|说明|
|:-:|:-:|
|post_id|对应的文章ID|
|tag_id|对应的标签ID|

&emsp;&emsp;为了找到同时属于"Java"、"MySQL"和"Redis"这3个标签文章，需要使用如下的SQL语句:

    select p.post_title
    from posts_tags pt, posts p, tags t
    where pt.tag_id = t.tag_id
        and (t.tag_name IN('Java', 'MySQL', 'Redis'))
        and pt.post_id = p.post_id
        group by p.post_id having count(p.post_id) = 3

&emsp;&emsp;可以很明显看到这样的SQL语句不仅效率相对较低，而且不易阅读和维护。而使用Redis可以很简单直接地实现这一需求。

&emsp;&emsp;具体做法是为每个标签使用一个名为 "tag:标签名称:posts"的集合类型键存储标有该标签的文章ID列表。假设现在又3篇文章，ID分别为1、2、3，其中ID为1的文章标签是"Java"，ID为2的文章标签是"Java"、"MySQL"，ID为3的文章标签是"Java"、"MySQL"和"Redis"，则有关标签的存储结构如下图：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/redis/epub_731187_52.jpg)


&emsp;&emsp;最简单的，当需要获取标记"MySQL"标签的文章时只需要使用`smembers tags:MySQL:posts`即可。如果要实现找到同时属于Java、MySQL和Redis这3个标签的文章，只需要将`tag:Java:posts`、`tag:MySQL:posts`和`tag:Redis:posts`这3个键**取交集**，借助`sinter`即可完成。


### 命令拾遗

|命令|描述|返回|
|:-:|:-:|:-:|
|scard key|获得集合中元素个数|集合中元素个数|
|sdiffstore dest key [key ...]|取差集并存储到dest中|
|sinterstore dest key [key ...]|取交集并存储到dest中|
|sunionstore dest key [key ...]|取并集并存储到dest中|
|srandmember key [count]|随机获得集合中的元素|随机返回count个元素|
|spop key|从集合中随机弹出一个元素|弹出元素值|

1. 当count为正数时，`srandmember`会随机从集合里获得count个不重复的元素的元素。如果count的值大于集合中的元素个数，则`srandmember`会返回集合中的全部元素。
2. 当count为负数时，`srandmember`会随机从集合里获得`|count|`个元素，这些元素有可能相同。

&emsp;&emsp;其实`srandmember`并不是严格的随机性，而是**伪随机**的。集合类型采用的存储结构时散列表，散列表使用散列函数将元素映射到不同的桶上以实现O(1)实践复杂度的元素查找。当使用散列表存储元素b时，使用散列函数计算出b的散列值是0，所以将b存储编号为0的桶中，下次要查找b时就可以用同样的散列函数再次计算b的散列值并直接接到相应的桶中找到b。当两个不同的元素的散列值相同时会出现冲突，Redis使用拉链法来解决冲突，即将散列值冲突的元素以链表的形式存入同一桶中，查找元素时先找到元素对应的桶，然后再从桶中的链表中找到对应的元素。使用`srandmember`从集合中获得一个随机元素时，Redis首先会从所有桶中随机选择一个桶，然后再从桶中的所有元素中随机选择一个元素，所以元素所在的桶中的元素数量越少，其被随机选中的可能性就越大。


## zset(最高级类型)

&emsp;&emsp;加入想要给文章访问量排序的话，就需要使用一种有序的数据结构了。有序集合类型(sorted set)在集合类型的基础上为集合中的每个元素都关联了一个**分数(score)**，这使得我们不仅可以完成插入、删除和判断元素是否存在等集合类型支持的操作，还能获得分数最高或最低的前N个元素、获得指定分数范围内的元素等于分数有关的操作。

&emsp;&emsp;有序集合在某些方面和列表类型非常相似：
1. 二者都是有序的
2. 二者都可以获得某一范围的元素

&emsp;&emsp;但二者有着很大的区别，这使得它们的应用场景也是不同的。
1. **列表类型是通过链表实现的**，获取靠近两端的数据速度极快，而当元素增多后，访问中间数据的速度会较慢，所以它更加适合实现如"热点"这样很少访问中间元素的应用。
2. **有序集合类型是使用散列表和跳跃表(Skiplist)实现的**，所以即使读取位于中间部分的数据速度也很快。
3. 列表中不能简单地调整某个元素地位置，但是有序集合可以。
4. 有序集合要比列表类型更耗费内存。

### 命令

|命令|描述|返回|
|:-:|:-:|:-:|
|zadd key score member [score member ...]|加入一个元素和该元素的分数|新加入地元素个数|
|zscore key member|获得元素的分数|member的分数|
|zrange key start stop [withscores]|获取某个范围的元素列表(从小到大)|元素列表|
|zrevrange key srart stop [withscores]|获取某个范围的元素列表(从大到小)|元素列表|
|zrangebyscore key min max [withscores] [limit offset count]|获得指定分数范围的元素|元素列表|
|zincrby key increment member|增加member的分数|更改后的分数|

    zadd scoreboard 89 tom 67 peter 100 david
    (integer) 3
    zadd scoreboard 76 peter
    (integer) 0
    # 分数还可以是浮点数
    zadd testboard 1.5 b
    zadd testboard 17E+307 a

    zscore scoreboard tom
    "89"

    zrange scoreboard 0 2
    "peter"
    "tom"
    "david"
    zrange scoreboard 1-1
    "tom"
    "david"
    zrange scoreboard 0-1 withscores
    "peter"
    "76"
    "tom"
    "89"
    "david"
    "100"

    zrangebyscore scoreboard 80 100
    "tom"
    "david"
    # 高于80,不包含80， +inf正无穷， -inf负无穷
    zrangebyscore scoreboard (80 +inf
    "tom"
    "david"
    zadd scoreboard 56 jerry 92 wendy 67 yvonne
    zrangebyscore scoreboard 60 +inf limit 1 3
    "peter"
    "tom"
    "wendy"
    zrevrangebyscore scoreboard 100 0 limit 0 3
    "david"
    "wendy"
    "tom"

    zincrby scoreboard 4 jerry
    "60"
    zincrby scoreboard -4 jerry
    "56"

### 实践

&emsp;&emsp;要按照文章的点击量排序，就必须再额外使用一个有序集合类型的键来实现。以文章的ID作为元素，以该文章的点击量作为该元素的分数：`posts:page:view`，每次用户访问一篇文章时，博客程序就通过`zincrby posts:page:view 1 ID`更新访问量。

### 命令拾遗

|命令|描述|返回|
|:-:|:-:|:-:|
|zcard key|获得集合中元素的数量|元素数量|
|zcount key min max|获得指定分数范围内的元素个数|元素数量|
|zrem key member [member ...]|删除一个或多个元素|成功删除的元素数量|
|zremrangebyrank key start stop|按照元素分数从小打到顺序删除处在指定排名范围内(索引)的所有元素|删除的元素数量|
|zremrangebyscore key min max|删除[min,max]分数范围内的所有元素|删除的元素数量|
|zrank key member|获得元素的排名|排名|
|zrevrank key member|获得元素的排名|排名|

    zcard scoreboard
    (integer) 6
    zcount scoreboard 90 100
    (integer) 2
    zcount scoreboard (89 +inf
    (integer) 2
    zrem scoreboard wendy
    (integer) 1
    zcard scoreboard
    (integer) 5
    zadd testRem 1 a 2 b 3 c 4 d 5 e 6 f
    (integer) 6
    zremrangebyrank testRem 0 2
    (integer) 3
    zrange testRem 0-1
    "d"
    "e"
    "f"
    zremrangebyscore testRem (4 5
    (integer) 1
    zrange testRem 0-1
    "d"
    "f"
    zrank scoreboard peter
    (integer) 0
    zrevrank scoreboard peter
    (integer) 4

- **zinterstore dest numkeys key [key ...] [weights weight [weight ...]] [aggregate sum|min|max]**

&emsp;&emsp;`zinterstore`命令用来计算多个有序集合的交集并将结果存储再dest键中，返回值为dest键中的元素个数。

1. 当`aggregate`是`sum`时，dest键中元素的分数是每个参与计算的集合中该元素分数的和。

        zadd srotedSets1 1 a 2 b
        (integer) 2
        zadd sortedSets2 10 a 20 b
        (integer) 2
        zinterstore sortedSetsResult 2 sortedSets1 sortedSets2
        (integer) 2
        zrange sortedSetsREsult 0-1 withscores
        "a"
        "11"
        "b"
        "22"

2. 当`aggregate`是`min`时，dest键中元素的分数是每个参与计算的集合中该元素分数的最小值(max同理)。


        zinterstore sortedSetsResult 2 sortedSets1 sortedSets2 aggregate min
        (integer) 2
        zrange sortedSetsREsult 0-1 withscores
        "a"
        "1"
        "b"
        "2"