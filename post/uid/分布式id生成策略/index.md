
# 美团

## 背景

&emsp;&emsp;在复杂分布式系统中，往往需要对大量的数据和消息进行唯一标识。如在美团点评的金融、支付、餐饮、酒店、猫眼电影等产品的系统中，数据日渐增长，对数据分库分表后需要有一个唯一 ID 来标识一条数据或消息，数据库的自增 ID 显然不能满足需求；特别一点的如订单、骑手、优惠券也都需要有唯一 ID 做标识。此时一个能够生成全局唯一 ID 的系统是非常必要的。概括下来，那业务系统对 ID 号的要求有哪些呢？

- **全局唯一性：**不能出现重复的 ID 号，既然是唯一标识，这是最基本的要求。
- **趋势递增：**在 MySQL InnoDB 引擎中使用的是聚集索引，由于多数 RDBMS 使用 B-tree 的数据结构来存储索引数据，在主键的选择上面我们应该尽量使用有序的主键保证写入性能。
- **单调递增：**保证下一个 ID 一定大于上一个 ID，例如事务版本号、IM 增量消息、排序等特殊需求。
- **信息安全：**如果 ID 是连续的，恶意用户的扒取工作就非常容易做了，直接按照顺序下载指定 URL 即可；如果是订单号就更危险了，竞对可以直接知道我们一天的单量。所以在一些应用场景下，会需要 ID 无规则、不规则。

&emsp;&emsp;`上述 1、2、3 对应三类不同的场景，3 和 4 需求还是互斥的，无法使用同一个方案满足。`

&emsp;&emsp;同时除了对 ID 号码自身的要求，业务还对 ID 号生成系统的可用性要求极高，想象一下，如果 ID 生成系统瘫痪，整个美团点评支付、优惠券发券、骑手派单等关键动作都无法执行，这就会带来一场灾难。

## 常见方法介绍

### UUID

&emsp;&emsp;UUID（Universally Unique Identifier）的标准型式包含 32 个 16 进制数字，以连字号分为五段，形式为8-4-4-4-12的 36 个字符，示例：550e8400-e29b-41d4-a716-446655440000

- **优点：**
  - 性能非常高：本地生成，没有网络消耗。
- **缺点：**
  - 不易于存储：UUID 太长，16 字节 128 位，通常以 36 长度的字符串表示，很多场景不适用。
  - 信息不安全：基于 MAC 地址生成 UUID 的算法可能会造成 MAC 地址泄露，这个漏洞曾被用于寻找梅丽莎病毒的制作者位置。
  - ID 作为主键时在特定的环境会存在一些问题，比如做 DB 主键的场景下，UUID 就非常不适用：
      - MySQL 官方有明确的建议主键要尽量越短越好，36 个字符长度的 UUID 不符合要求。
      - 对 MySQL 索引不利：如果作为数据库主键，在 InnoDB 引擎下，UUID 的无序性可能会引起数据位置频繁变动，严重影响性能。

### snowflake方案

&emsp;&emsp;这种方案大致来说是一种以划分命名空间（UUID 也算，由于比较常见，所以单独分析）来生成 ID 的一种算法，这种方案把 64-bit 分别划分成多段，分开来标示机器、时间等，比如在 snowflake 中的 64-bit 分别表示如下图（图片来自网络）所示：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/uid/20201015161031476.png)

&emsp;&emsp;41-bit 的时间可以表示`(1L<<41)/(1000L*3600*24*365)=69`年的时间，10-bit 机器可以分别表示 1024 台机器。如果我们对 IDC 划分有需求，还可以将 10-bit 分 5-bit 给 IDC，分 5-bit 给工作机器。这样就可以表示 32 个 IDC，每个 IDC 下可以有 32 台机器，可以根据自身需求定义。12 个自增序列号可以表示2^12个 ID，理论上 snowflake 方案的 QPS 约为 409.6w/s，这种分配方式可以保证在任何一个 IDC 的任何一台机器在任意毫秒内生成的 ID 都是不同的。

- **优点：**
  - 毫秒数在高位，自增序列在低位，整个 ID 都是趋势递增的。
  - 不依赖数据库等第三方系统，以服务的方式部署，稳定性更高，生成 ID 的性能也是非常高的。
  - 可以根据自身业务特性分配 bit 位，非常灵活。
- **缺点：**
  - 强依赖机器时钟，如果机器上时钟回拨，会导致发号重复或者服务会处于不可用状态。

### 数据库生成

&emsp;&emsp;以 MySQL 举例，利用给字段设置`auto_increment_increment`和`auto_increment_offset`来保证 ID 自增，每次业务使用下列 SQL 读写 MySQL 得到 ID 号。

    begin;
    REPLACE INTO Tickets64 (stub) VALUES ('a');
    SELECT LAST_INSERT_ID();
    commit;

- **优点：**
  - 非常简单，利用现有数据库系统的功能实现，成本小，有 DBA 专业维护。
  - ID 号单调自增，可以实现一些对 ID 有特殊要求的业务。
- **缺点：**
  - 强依赖 DB，当 DB 异常时整个系统不可用，属于致命问题。配置主从复制可以尽可能的增加可用性，但是数据一致性在特殊情况下难以保证。主从切换时的不一致可能会导致重复发号。
  - ID 发号性能瓶颈限制在单台 MySQL 的读写性能。

## 美团Leaf方案实现

&emsp;&emsp;`Leaf`这个名字是来自德国哲学家、数学家莱布尼茨的一句话：`There are no two identical leaves in the world`，即**“世界上没有两片相同的树叶”**。

&emsp;&emsp;综合对比上述几种方案，每种方案都不完全符合我们的要求。所以 Leaf 分别在上述第二种和第三种方案上做了相应的优化，实现了 Leaf-segment 和 Leaf-snowflake 方案。

### Leaf-segment数据库方案

&emsp;&emsp;第一种 Leaf-segment 方案，在使用数据库的方案上，做了如下改变：

- 原方案每次获取 ID 都得读写一次数据库，造成数据库压力大。改为利用`proxy server`批量获取，每次获取一个`segment`（step 决定大小）号段的值。用完之后再去数据库获取新的号段，可以大大的减轻数据库的压力。
- 各个业务不同的发号需求用`biz_tag`字段来区分，每个biz-tag的 ID 获取相互隔离，互不影响。如果以后有性能需求需要对数据库扩容，不需要上述描述的复杂的扩容操作，只需要对biz_tag分库分表就行。


        Create Table: CREATE TABLE `leaf_alloc` (
        `biz_tag` varchar(128) NOT NULL DEFAULT '' comment '业务类型',
        `max_id` bigint(20) NOT NULL DEFAULT '1' comment '当前最大值',
        `step` int(11) NOT NULL comment '步长',
        `description` varchar(256) DEFAULT NULL,
        `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        PRIMARY KEY (`biz_tag`)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8

&emsp;&emsp;重要字段说明：`biz_tag`用来区分业务，`max_id`表示该biz_tag目前所被分配的 ID 号段的最大值，`step`表示每次分配的号段长度。原来获取 ID 每次都需要写数据库，现在只需要把step设置得足够大，比如 1000。那么只有当 1000 个号被消耗完了之后才会去重新读写一次数据库。读写数据库的频率从 1 减小到了`1/step`，大致架构如下图所示：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/uid/20201015164520311.png)

&emsp;&emsp;`test_tag`在第一台 Leaf 机器上是`1~1000`的号段，当这个号段用完时，会去加载另一个长度为`step=1000`的号段，假设另外两台号段都没有更新，这个时候第一台机器新加载的号段就应该是`3001~4000`。同时数据库对应的biz_tag这条数据的max_id会从 3000 被更新成 4000，更新号段的 SQL 语句如下：

    Begin
    UPDATE table SET max_id=max_id+step WHERE biz_tag=xxx
    SELECT tag, max_id, step FROM table WHERE biz_tag=xxx
    Commit

- **优点：**
  - Leaf 服务可以很方便的线性扩展，性能完全能够支撑大多数业务场景。
  - ID 号码是趋势递增的 8byte 的 64 位数字，满足上述数据库存储的主键要求。
  - 容灾性高：Leaf 服务内部有号段缓存，即使 DB 宕机，短时间内 Leaf 仍能正常对外提供服务。
  - 可以自定义max_id的大小，非常方便业务从原有的 ID 方式上迁移过来。
- **缺点：**
  - ID 号码不够随机，能够泄露发号数量的信息，不太安全。
  - TP999 数据波动大，当号段使用完之后还是会 hang 在更新数据库的 I/O 上，TP999 数据会出现偶尔的尖刺。
  - DB 宕机会造成整个系统不可用。

### 双buffer优化

&emsp;&emsp;对于第二个缺点，Leaf-segment 做了一些优化，简单的说就是：

- Leaf 取号段的时机是在号段消耗完的时候进行的，也就意味着号段临界点的 ID 下发时间取决于下一次从 DB 取回号段的时间，并且在这期间进来的请求也会因为 DB 号段没有取回来，导致线程阻塞。如果请求 DB 的网络和 DB 的性能稳定，这种情况对系统的影响是不大的，但是假如取 DB 的时候网络发生抖动，或者 DB 发生慢查询就会导致整个系统的响应时间变慢。

&emsp;&emsp;为此，我们希望 DB 取号段的过程能够做到无阻塞，不需要在 DB 取号段的时候阻塞请求线程，即当号段消费到某个点时就异步的把下一个号段加载到内存中。而不需要等到号段用尽的时候才去更新号段。这样做就可以很大程度上的降低系统的 TP999 指标。详细实现如下图所示：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/uid/20201015165005919.png)

&emsp;&emsp;采用双 buffer 的方式，Leaf 服务内部有两个号段缓存区segment。当前号段已下发 10% 时，如果下一个号段未更新，则另启一个更新线程去更新下一个号段。当前号段全部下发完后，如果下个号段准备好了则切换到下个号段为当前segment接着下发，循环往复。

&emsp;&emsp;每个biz-tag都有消费速度监控，通常推荐segment长度设置为服务高峰期发号 QPS 的 600 倍（10分钟），这样即使 DB 宕机，Leaf 仍能持续发号 10 ~ 20 分钟不受影响。

&emsp;&emsp;每次请求来临时都会判断下个号段的状态，从而更新此号段，所以偶尔的网络抖动不会影响下个号段的更新。

### Leaf-snowflake方案

&emsp;&emsp;Leaf-segment 方案可以生成趋势递增的 ID，同时 ID 号是可计算的，不适用于订单 ID 生成场景，比如竞对在两天中午 12 点分别下单，通过订单 ID 号相减就能大致计算出公司一天的订单量，这个是不能忍受的。面对这一问题，我们提供了 Leaf-snowflake 方案。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/uid/20201015161031476.png)

&emsp;&emsp;Leaf-snowflake 方案完全沿用 snowflake 方案的 bit 位设计，即是`1+41+10+12`的方式组装 ID 号。对于`workerID`的分配，当服务集群数量较小的情况下，完全可以手动配置。Leaf 服务规模较大，动手配置成本太高。所以使用 `Zookeeper` 持久顺序节点的特性自动对 snowflake 节点配置wokerID。Leaf-snowflake 是按照下面几个步骤启动的：

1. 启动 Leaf-snowflake 服务，连接 Zookeeper，在`leaf_forever`父节点下检查自己是否已经注册过（是否有该顺序子节点）。
2. 如果有注册过直接取回自己的`workerID`（ZK 顺序节点生成的int类型 ID 号），启动服务。
3. 如果没有注册过，就在该父节点下面创建一个持久顺序节点，创建成功后取回顺序号当做自己的`workerID`号，启动服务。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/uid/20201015165718157.png)

# 百度

&emsp;&emsp;`UidGenerator` 是百度开源的 Java 语言实现，基于 Snowflake 算法的唯一 ID 生成器。而且，它非常适合虚拟环境，比如：Docker。另外，它通过消费未来时间克服了雪花算法的并发限制。UidGenerator 提前生成 ID 并缓存在`RingBuffer`中。 检测结果显示，单个实例的 QPS 能超过 6000,000。

&emsp;&emsp;百度对snowflake组成部分稍微调整了一下：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/uid/20201015171435362.png)

&emsp;&emsp;由上图可知，UidGenerator 的时间部分只有 28 位，这就意味着 UidGenerator 默认只能承受 8.5 年（`2^28-1/86400/365`）。当然，根据你业务的需求，UidGenerator 可以适当调整`delta seconds、worker node id`和`sequence`占用位数。

&emsp;&emsp;接下来分析百度 UidGenerator 的实现。需要说明的是 UidGenerator 有两种方式提供：`DefaultUidGenerato`r和`CachedUidGenerator`。我们先分析比较容易理解的DefaultUidGenerator。

## DefaultUidGenerator

### delta seconds

&emsp;&emsp;这个值是指当前时间与`epoch`时间的时间差，且单位为秒。epoch时间就是指集成 UidGenerator 生成分布式 ID 服务第一次上线的时间，可配置，也一定要根据你的上线时间进行配置，因为默认的epoch时间可是`2016-09-20`，不配置的话，会浪费好几年的可用时间。

### worker id

&emsp;&emsp;接下来说一下 UidGenerator 是如何被worker id赋值的，搭建 UidGenerator 的话，需要创建一个表：

    DROP DATABASE IF EXISTS `xxxx`;
    CREATE DATABASE `xxxx` ;
    use `xxxx`;
    DROP TABLE IF EXISTS WORKER_NODE;
    CREATE TABLE WORKER_NODE
    (
    ID BIGINT NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
    HOST_NAME VARCHAR(64) NOT NULL COMMENT 'host name',
    PORT VARCHAR(64) NOT NULL COMMENT 'port',
    TYPE INT NOT NULL COMMENT 'node type: ACTUAL or CONTAINER',
    LAUNCH_DATE DATE NOT NULL COMMENT 'launch date',
    MODIFIED TIMESTAMP NOT NULL COMMENT 'modified time',
    CREATED TIMESTAMP NOT NULL COMMENT 'created time',
    PRIMARY KEY(ID)
    )
    COMMENT='DB WorkerID Assigner for UID Generator',ENGINE = INNODB;


&emsp;&emsp;UidGenerator 会在集成用它生成分布式 ID 的实例启动的时候，往这个表中插入一行数据，得到的 ID 值就是准备赋给workerId的值。由于workerId默认 22 位，那么，集成 UidGenerator 生成分布式 ID 的所有实例重启次数是不允许超过 4194303 次（即2^22-1），否则会抛出异常。

### sequence

- synchronized保证线程安全
- 如果时间有任何的回拨，那么直接抛出异常
- 如果当前时间和上一次是同一秒时间，那么sequence自增。如果同一秒内自增值超过2^13-1，那么就会自旋等待下一秒（getNextSecond）
- 如果是新的一秒，那么sequence重新从 0 开始

&emsp;&emsp;通过DefaultUidGenerator的实现可知，它对时钟回拨的处理比较简单粗暴。另外如果使用 UidGenerator 的DefaultUidGenerator方式生成分布式 ID，一定要根据你的业务的情况和特点，调整各个字段占用的位数：

    <property name="timeBits" value="29"/>
    <property name="workerBits" value="21"/>
    <property name="seqBits" value="13"/>
    <property name="epochStr" value="2016-09-20"/>

## CachedUidGenerator

&emsp;&emsp;CachedUidGenerator是 UidGenerator 的重要改进实现。它的核心利用了R`ingBuffer`，如下图所示，它本质上是一个数组，数组中每个项被称为`slot`。UidGenerator 设计了两个RingBuffer，一个保存唯一 ID，一个保存flag。RingBuffer的尺寸是2^n，n必须是正整数：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/uid/20201015172519634.png)

### RingBuffer Of Flag

&emsp;&emsp;其中，保存flag这个RingBuffer的每个slot的值都是 0 或者 1，0 是CANPUTFLAG的标志位，1 是CANTAKEFLAG的标识位。每个slot的状态要么是CANPUT，要么是CANTAKE。以某个slot的值为例，初始值为 0，即CANPUT。接下来会初始化填满这个RingBuffer，这时候这个slot的值就是 1，即CANTAKE。等获取分布式 ID 时取到这个slot的值后，这个slot的值又变为 0，以此类推。

### RingBuffer Of UID

&emsp;&emsp;保存唯一 ID 的RingBuffer有两个指针，Tail指针和Cursor指针。

&emsp;&emsp;Tail指针表示最后一个生成的唯一 ID。如果这个指针对上了Cursor指针，意味着RingBuffer已经满了。这时候，不允许再继续生成 ID 了。用户可以通过属性rejectedPutBufferHandler指定处理这种情况的策略。

&emsp;&emsp;Cursor指针表示最后一个已经给消费的唯一 ID。如果Cursor指针追上了Tail指针，意味着RingBuffer已经空了。这时候，不允许再继续获取 ID 了。用户可以通过属性rejectedTakeBufferHandler指定处理这种异常情况的策略。

&emsp;&emsp;另外，如果你想增强RingBuffer提升它的吞吐能力，那么需要配置一个更大的boostPower值：

    <property name="boostPower" value="3"></property>

## 百度小结

- **自增列：**UidGenerator的workerId在实例每次重启时初始化，且就是数据库的自增 ID，从而完美的实现每个实例获取到的workerId不会有任何冲突。
- **RingBuffer：**UidGenerator 不再在每次取 ID 时都实时计算分布式 ID，而是利用RingBuffer数据结构预先生成若干个分布式 ID 并保存。
- **时间递增：**传统的雪花算法实现都是通过System.currentTimeMillis()来获取时间并与上一次时间进行比较，这样的实现严重依赖服务器的时间。而 UidGenerator 的时间类型是AtomicLong，且通过incrementAndGet()方法获取下一次的时间，从而脱离了对服务器时间的依赖，也就不会有时钟回拨的问题。·