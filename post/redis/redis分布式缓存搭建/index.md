
> 转载自 https://www.cnblogs.com/esofar/p/10486621.html

# 摘要

&emsp;&emsp;Redis Cluster 即 Redis 集群，是 Redis 官方在 3.0 版本推出的一套分布式存储方案。完全去中心化，由多个节点组成，所有节点彼此互联。Redis 客户端可以直接连接任何一节点获取集群中的键值对，不需要中间代理，如果该节点不存在用户所指定的键值，其内部会自动把客户端重定向到键值所在的节点。

&emsp;&emsp;edis 集群是一个网状结构，每个节点都通过 TCP 连接跟其他每个节点连接。在一个有 N 个节点的集群中，每个节点都有 N-1 个流出的 TCP 连接，和 N-1 个流入的连接，这些 TCP 连接会永久保持。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/redis/cluster.png)

## 数据分区原理

&emsp;&emsp;Redis 集群会将用户数据分散保存至各个节点中，突破单机 Redis 内存最大存储容量。集群引入了 哈希槽slot的概念，其搭建完成后会生 16384 个哈希槽slot，同时会根据节点的数量大致均等的将 16384 个哈希槽映射到不同的节点上。当用户存储key-value时，集群会先对key进行 CRC16 校验然后对 16384 取模来决定key-value放置哪个槽，从而实现自动分割数据到不同的节点上。

## 数据冗余

&emsp;&emsp;Redis 集群支持主从复制和故障恢复。集群使用了主从复制模型，每个主节点master应至少有一个从节点slave。假设某个主节点故障，其所有子节点会广播一个数据包给集群里的其他主节点来请求选票，一旦某个从节点收到了大多数主节点的回应，那么它就赢得了选举，被推选为主节点，负责处理之前旧的主节点负责的哈希槽。

## 可以存储多少数据

&emsp;&emsp;在初始化Redis时，需要在redis.conf中的`maxmemory`配置内存大小,该选项告诉Redis当使用了多少物理内存后就开始拒绝后续的请求写入，该参数能很好的保护Redis不会因为使用了过多的物理内存而导致swap。

&emsp;&emsp;`官网说单例可以处理2.5亿个key，每个key最大占用内存为512M`

## 单节点安装

        Lunix版本：Linux A-80-5 3.10.0-1062.9.1.el7.x86_64 GNU/Linux
        Redis版本：6.0.6

1. 准备环境

        yum -y install gcc tcl
        yum -y install centos-release-scl
        yum -y install devtoolset-9-gcc devtoolset-9-gcc-c++ devtoolset-9-binutils
        scl enable devtoolset-9 bash
        #echo "source /opt/rh/devtoolset-9/enable" >> /etc/profile 永久生效

2. 下载Redis

        wget http://download.redis.io/releases/redis-6.0.6.tar.gz

3. 编译安装

        tar -zxvf redis-6.0.6.tar.gz -C /usr/local/
        cd /usr/local/redis-6.0.6
        make MALLOC=libc
        make test
        make install PREFIX=/usr/local/redis

## 搭建Redis分布式集群

&emsp;&emsp;依据 Redis Cluster 内部故障转移实现原理，Redis 集群至少需要 3 个主节点，而每个主节点至少有 1 从节点，因此搭建一个集群至少包含 6 个节点，三主三从，并且分别部署在不同机器上。

&emsp;&emsp;条件有限，测试环境下我们只能在一台机器上创建一个伪集群，通过不同的 TCP 端口启动多个 Redis 实例，组成集群。

&emsp;&emsp;目前 Redis Cluster 的搭建有两种方式：

1. 手动方式搭建，即手动执行 cluster 命令，一步步完成搭建流程。
2. 自动方式搭建，即使用官方提供的集群管理工具快速搭建。
</br>

&emsp;&emsp;两种方式原理一样，自动搭建方式只是将手动搭建方式中需要执行的 Redis 命令封装到了可执行程序。生产环境下推荐使用第二种方式，简单快捷，不易出错。不过本文实战演示两种方式都会提及。

## 手动方式搭建

### 启动节点

&emsp;&emsp;搭建集群的第一步就是要先把参与搭建集群的每个节点启动起来。由于我们这是在一台机器上模拟多个节点，可以预先规划下各个节点的属性：

|节点编号|IP地址|TCP端口|节点类型|从节点|启动配置|
|:-:|:-:|:-:|:-:|:-:|:-:|
|A|127.0.0.1|7001|主|D|/usr/local/redis-cluster/7001/redis.conf|
|B|127.0.0.1|7002|主|E|/usr/local/redis-cluster/7002/redis.conf|
|C|127.0.0.1|7003|主|F|/usr/local/redis-cluster/7003/redis.conf|
|D|127.0.0.1|8001|从|/|/usr/local/redis-cluster/8001/redis.conf|
|E|127.0.0.1|8002|从|/|/usr/local/redis-cluster/8002/redis.conf|
|F|127.0.0.1|8003|从|/|/usr/local/redis-cluster/8003/redis.conf|

&emsp;&emsp;根据上述规划，可以先通过如下命令创建各个节点启动配置文件的存放目录。

        mkdir -p /usr/local/redis-cluster && cd /usr/local/redis-cluster
        mkdir -p 7001 7002 7003 8001 8002 8003

&emsp;&emsp;开放端口：

        firewall-cmd --add-port=7001/tcp --zone=public --permanent
        firewall-cmd --add-port=7002/tcp --zone=public --permanent
        firewall-cmd --add-port=7003/tcp --zone=public --permanent
        firewall-cmd --add-port=8001/tcp --zone=public --permanent
        firewall-cmd --add-port=8003/tcp --zone=public --permanent
        firewall-cmd --add-port=8003/tcp --zone=public --permanent
        # 通信端口
        firewall-cmd --add-port=17001/tcp --zone=public --permanent
        firewall-cmd --add-port=17002/tcp --zone=public --permanent
        firewall-cmd --add-port=17003/tcp --zone=public --permanent
        firewall-cmd --add-port=18001/tcp --zone=public --permanent
        firewall-cmd --add-port=18003/tcp --zone=public --permanent
        firewall-cmd --add-port=18003/tcp --zone=public --permanent

        firewall-cmd --reload

&emsp;&emsp;顺序执行如下行命令，进入 Redis 源码包目录并将默认配置文件redis.conf分别复制到六个节点配置存放目录中，作为各自节点启动配置文件。

        cd /usr/local/redis
        cp redis.conf /usr/local/redis-cluster/7001 
        cp redis.conf /usr/local/redis-cluster/7002
        cp redis.conf /usr/local/redis-cluster/7003 
        cp redis.conf /usr/local/redis-cluster/8001
        cp redis.conf /usr/local/redis-cluster/8002 
        cp redis.conf /usr/local/redis-cluster/8003

&emsp;&emsp;接下来需要分别修改每个节点的配置文件。下面贴的是节点 A 的配置文件/usr/local/redis-cluster/7001/redis.conf中启用或修改的一些必要参数。其他节点 B、C、D、E、F 参照修改，注意把涉及端口的地方修改成各自节点预先规划的即可。

    #bind 192.168.83.128                         # 当前节点主机公网IP
    port 7001                                   # 客户端连接监听端口
    protected-mode no                           # 保护模式
    #masterauth 123456                           # 启用密码认证
    #requirepass 123456                          # masterauth 和 requirepass 必须同时设置
    maxmemory 1GB                               # 最大物理内存
    maxmemory-policy allkeys-lru                # 采用lru缓存清除策略
    databases 1                                 # 集群模式只能使用0号数据库
    pidfile /var/run/redis_7001.pid             # Redis实例pid文件
    daemonize yes                               # 以守护进程运行Redis实例
    loglevel notice                             # debug/verbose/notice/warning
    #save 300 10                                # 300s内有10个key发生变化，就持久化
    #dbfilename dump.rdb                        # 快照文件名
    logfile /data/module/redis/log/redis.log    # 日志文件
    #dir /data/module/redis/data                # 快照文件存放路径
    cluster-enabled yes                         # 启用集群模式
    cluster-node-timeout 15000                  # 当前节点连接超时毫秒数
    cluster-config-file node-7001.conf          # 当前节点集群配置文件路径

&emsp;&emsp;完成上述工作就可以通过一下几组命令启动待搭建集群中的6个节点了。

    /usr/local/bin/redis-server /usr/local/redis-cluster/7001/redis.conf
    /usr/local/bin/redis-server /usr/local/redis-cluster/7002/redis.conf
    /usr/local/bin/redis-server /usr/local/redis-cluster/7003/redis.conf
    /usr/local/bin/redis-server /usr/local/redis-cluster/8001/redis.conf
    /usr/local/bin/redis-server /usr/local/redis-cluster/8002/redis.conf
    /usr/local/bin/redis-server /usr/local/redis-cluster/8003/redis.conf

&emsp;&emsp;最后通过ps -ef|grep redis命令确认各个节点服务是否已经正常运行。

    [root@192.168.83.128 bin]# ps -fe | grep redis
    root       5613      1  0 04:25 ?        00:00:00 /usr/local/bin/redis-server 127.0.0.1:7001 [cluster]
    root       5650      1  0 04:26 ?        00:00:00 /usr/local/bin/redis-server 127.0.0.1:7002 [cluster]
    root       5661      1  0 04:26 ?        00:00:00 /usr/local/bin/redis-server 127.0.0.1:7003 [cluster]
    root       5672      1  0 04:27 ?        00:00:00 /usr/local/bin/redis-server 127.0.0.1:8001 [cluster]
    root       5681      1  0 04:27 ?        00:00:00 /usr/local/bin/redis-server 127.0.0.1:8002 [cluster]
    root       5690      1  0 04:27 ?        00:00:00 /usr/local/bin/redis-server 127.0.0.1:8003 [cluster]
    root       5731   1311  0 04:28 pts/0    00:00:00 grep --color=auto redis

&emsp;&emsp;如上输出可以看书上面规划的6个节点都成功启动。

### 节点握手

&emsp;&emsp;虽然上面 6 个节点都启用了群集支持，但默认情况下它们是不相互信任或者说没有联系的。节点握手就是在各个节点之间创建链接（每个节点与其他节点相连），形成一个完整的网格，即集群。
&emsp;&emsp;节点握手的命令如下：

    cluster meet ip port

&emsp;&emsp;但为了创建群集，不需要发送形成完整网格所需的所有 `cluster meet` 命令。只要能发送足够的cluster meet消息，可以让每个节点都可以通过一系列已知节点到达每个其他节点，缺失的链接将被自动创建。

&emsp;&emsp;例如，如果我们通过`cluster meet`将节点 A 与节点 B 连接起来，并将 B 与 C 连接起来，则 A 和 C 会自己找到握手方式并创建链接。

&emsp;&emsp;我们的创建的 6 个节点可以通过 redis-cli 连接到 A 节点执行如下五组命令完成握手，生产环境需要将 IP 127.0.0.1替换成外网 IP。

    cluster meet 127.0.0.1 7002
    cluster meet 127.0.0.1 7003
    cluster meet 127.0.0.1 8001
    cluster meet 127.0.0.1 8002
    cluster meet 127.0.0.1 8003

&emsp;&emsp;如上述命令正常执行输出结果如下。

    [root@192.168.83.128 bin]# /usr/local/bin/redis-cli -p 7001
    127.0.0.1:7001> cluster meet 127.0.0.1 7002
    OK
    127.0.0.1:7001> cluster meet 127.0.0.1 7003
    OK
    127.0.0.1:7001> cluster meet 127.0.0.1 8001
    OK
    127.0.0.1:7001> cluster meet 127.0.0.1 8002
    OK
    127.0.0.1:7001> cluster meet 127.0.0.1 8003
    OK

&emsp;&emsp;接下来可以通过 cluster nodes 命令查看节点之间 的链接状态。我随机找了两个节点 B 和 F 测试，输出结果如下所示。

    [root@192.168.83.128 bin]# /usr/local/bin/redis-cli -p 7002 cluster nodes
    61e8c4ed8d1ff2a765a4dd2c3d300d8121d26e12 127.0.0.1:7001@17001 master - 0 1552220691885 4 connected
    a8a41694f22977fda78863bdfb3fc03dd1fab1bd 127.0.0.1:8002@18002 master - 0 1552220691000 5 connected
    51987c4b5530c81f2845bb9d521daf6d3dce3659 127.0.0.1:8001@18001 master - 0 1552220690878 3 connected
    1b4b3741945d7fed472a1324aaaa6acaa1843ccb 127.0.0.1:7002@17002 myself,master - 0 1552220690000 1 connected
    19147f56e679767bcebb8653262ff7f56ca072a8 127.0.0.1:7003@17003 master - 0 1552220691000 2 connected
    ed6fd72e61b747af3705b210c7164bc68739303e 127.0.0.1:8003@18003 master - 0 1552220690000 0 connected

    [root@192.168.83.128 bin]# /usr/local/bin/redis-cli -p 8002 cluster nodes
    1b4b3741945d7fed472a1324aaaa6acaa1843ccb 127.0.0.1:7002@17002 master - 0 1552220700255 1 connected
    ed6fd72e61b747af3705b210c7164bc68739303e 127.0.0.1:8003@18003 master - 0 1552220703281 0 connected
    19147f56e679767bcebb8653262ff7f56ca072a8 127.0.0.1:7003@17003 master - 0 1552220700000 2 connected
    a8a41694f22977fda78863bdfb3fc03dd1fab1bd 127.0.0.1:8002@18002 myself,master - 0 1552220701000 5 connected
    61e8c4ed8d1ff2a765a4dd2c3d300d8121d26e12 127.0.0.1:7001@17001 master - 0 1552220702275 4 connected
    51987c4b5530c81f2845bb9d521daf6d3dce3659 127.0.0.1:8001@18001 master - 0 1552220701265 3 connected

&emsp;&emsp;可以看到，节点B和节点F都已经分别和其它5个节点建立链接。
&emsp;&emsp;至此，节点握手完成。

### 分配槽位

&emsp;&emsp;此时Redis集群还没有处理上线状态，可以在任意一节点上执行`cluster info` 命令来查看目前集群的运行状态。

    [root@192.168.83.128 bin]# /usr/local/bin/redis-cli -p 7001 cluster info
    cluster_state:fail

&emsp;&emsp;上面输出`cluster_state:fail`表示当前集群处于下线状态。因为只有给集群中所有主节点分配好槽位（即哈希槽slot，本文第一小节有提及）集群才能上线。

&emsp;&emsp;分配槽位的命令如下：

    cluster addslots slot [slot...]

&emsp;&emsp;根据预先规划，这一步需要使用 `cluster addslots` 命令手动将 16384 个哈希槽大致均等分配给主节点 A、B、C。 

    /usr/local/bin/redis-cli -p 7001 cluster addslots {0..5461}
    /usr/local/bin/redis-cli -p 7002 cluster addslots {5462..10922}
    /usr/local/bin/redis-cli -p 7003 cluster addslots {10923..16383}

&emsp;&emsp;上面三组命令执行完毕，可以再次查看目前集群的一些运行参数

    [root@192.168.83.128 bin]# /usr/local/bin/redis-cli -p 7001 cluster info
    cluster_state:ok
    cluster_slots_assigned:16384
    cluster_slots_ok:16384
    cluster_slots_pfail:0
    cluster_slots_fail:0
    cluster_known_nodes:6
    cluster_size:3
    cluster_current_epoch:5
    cluster_my_epoch:4
    cluster_stats_messages_ping_sent:11413
    cluster_stats_messages_pong_sent:10509
    cluster_stats_messages_meet_sent:11
    cluster_stats_messages_sent:21933
    cluster_stats_messages_ping_received:10509
    cluster_stats_messages_pong_received:10535
    cluster_stats_messages_received:21044

&emsp;&emsp;如上输出cluster_state:ok证明 Redis 集群成功上线。

### 主从复制

&emsp;&emsp;Redis 集群成功上线，不过还没有给主节点指定从节点，此时如果有一个节点故障，那么整个集群也就挂了，也就无法实现高可用。

&emsp;&emsp;集群中需要使用 cluster replicate 命令手动给从节点配置主节点。

&emsp;&emsp;集群复制命令如下：

    cluster replicate node-id

&emsp;&emsp;集群中各个节点的`node-id`可以用`cluster nodes`命令查看，如下输出`1b4b3741945d7fed472a1324aaaa6acaa1843ccb`即是主节点 B 的`node-id`。

    [root@192.168.83.128 bin]# /usr/local/bin/redis-cli -p 8002 cluster nodes
    1b4b3741945d7fed472a1324aaaa6acaa1843ccb 127.0.0.1:7002@17002 master - 0 1552220700255 1 connected
    ed6fd72e61b747af3705b210c7164bc68739303e 127.0.0.1:8003@18003 master - 0 1552220703281 0 connected
    19147f56e679767bcebb8653262ff7f56ca072a8 127.0.0.1:7003@17003 master - 0 1552220700000 2 connected
    a8a41694f22977fda78863bdfb3fc03dd1fab1bd 127.0.0.1:8002@18002 myself,master - 0 1552220701000 5 connected
    61e8c4ed8d1ff2a765a4dd2c3d300d8121d26e12 127.0.0.1:7001@17001 master - 0 1552220702275 4 connected
    51987c4b5530c81f2845bb9d521daf6d3dce3659 127.0.0.1:8001@18001 master - 0 1552220701265 3 connected

&emsp;&emsp;根据预先规划，A主D从；B主E从；C主F从。执行如下三组命令分别为从节点 D、E、F 指定其主节点，使群集可以自动完成主从复制。

    /usr/local/bin/redis-cli -p 8001 cluster replicate 61e8c4ed8d1ff2a765a4dd2c3d300d8121d26e12
    /usr/local/bin/redis-cli -p 8002 cluster replicate 1b4b3741945d7fed472a1324aaaa6acaa1843ccb
    /usr/local/bin/redis-cli -p 8003 cluster replicate 19147f56e679767bcebb8653262ff7f56ca072a8

&emsp;&emsp;命令执行成功后，我们便算以手动方式成功搭建了一个 Redis 集群。

&emsp;&emsp;最后，再来查看一下集群中的节点信息。

    [root@192.168.83.128 bin]# /usr/local/bin/redis-cli -p 8002 cluster nodes
    1b4b3741945d7fed472a1324aaaa6acaa1843ccb 127.0.0.1:7002@17002 master - 0 1552233328337 1 connected 5462-10922
    ed6fd72e61b747af3705b210c7164bc68739303e 127.0.0.1:8003@18003 slave 19147f56e679767bcebb8653262ff7f56ca072a8 0 1552233327000 2 connected
    19147f56e679767bcebb8653262ff7f56ca072a8 127.0.0.1:7003@17003 master - 0 1552233325000 2 connected 10923-16383
    a8a41694f22977fda78863bdfb3fc03dd1fab1bd 127.0.0.1:8002@18002 myself,slave 1b4b3741945d7fed472a1324aaaa6acaa1843ccb 0 1552233327000 5 connected
    61e8c4ed8d1ff2a765a4dd2c3d300d8121d26e12 127.0.0.1:7001@17001 master - 0 1552233327327 4 connected 0-5461
    51987c4b5530c81f2845bb9d521daf6d3dce3659 127.0.0.1:8001@18001 slave 61e8c4ed8d1ff2a765a4dd2c3d300d8121d26e12 0 1552233326320 4 connected


## 自动搭建方式

&emsp;&emsp;Redis 3.0 版本之后官方发布了一个集群管理工具 `redis-trib.rb`，集成在 Redis 源码包的src目录下。其封装了 Redis 提供的集群命令，使用简单、便捷。

不过 redis-trib.rb 是 Redis 作者使用 Ruby 语言开发的，故使用该工具之前还需要先在机器上安装 Ruby 环境。后面作者可能意识到这个问题，Redis 5.0 版本开始便把这个工具集成到 `redis-cli` 中，以`--cluster`参数提供使用，其中create命令可以用来创建集群。

### 启动节点

&emsp;&emsp;使用集群管理工具搭建集群之前，也是需要先把各个节点启动起来的。节点的启动方式请参见本文「手动方式创建」-「启动节点」一节，此处不再赘述。

### 集群管理工具搭建

    /usr/local/bin/redis-cli --cluster create 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:8001 127.0.0.1:8002 127.0.0.1:8003 --cluster-replicas 1

&emsp;&emsp;主节点在前，从节点在后。其中--cluster-replicas参数用来指定一个主节点带有的从节点个数，如上--cluster-replicas 1即表示 1 个主节点有 1 个从节点。

&emsp;&emsp;命令执行成功会有类似如下输出。

    [root@192.168.83.128 bin]# redis-cli --cluster create 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:8001 127.0.0.1:8002 127.0.0.1:8003 --cluster-replicas 1
    >>> Performing hash slots allocation on 6 nodes...
    Master[0] -> Slots 0 - 5460
    Master[1] -> Slots 5461 - 10922
    Master[2] -> Slots 10923 - 16383
    Adding replica 127.0.0.1:8001 to 127.0.0.1:7001
    Adding replica 127.0.0.1:8002 to 127.0.0.1:7002
    Adding replica 127.0.0.1:8003 to 127.0.0.1:7003
    >>> Trying to optimize slaves allocation for anti-affinity
    [WARNING] Some slaves are in the same host as their master
    M: 32f9819fc7d561bfa2b7189182200e86d9901b8a 127.0.0.1:7001
    slots:[0-5460] (5461 slots) master
    M: cca0fbfa374bc175d481e68ee9ed13b65453e967 127.0.0.1:7002
    slots:[5461-10922] (5462 slots) master
    M: 964cfa1c2dcfe36b6d3c63637f0d57ccb568354e 127.0.0.1:7003
    slots:[10923-16383] (5461 slots) master
    S: 1b47b9e6e7a79523579b8d2ddcd5e708583ed317 127.0.0.1:8001
    replicates 32f9819fc7d561bfa2b7189182200e86d9901b8a
    S: aba9330f3e70f26a8af4ced1b672fbcc7bc62d78 127.0.0.1:8002
    replicates cca0fbfa374bc175d481e68ee9ed13b65453e967
    S: 254db0830cd764e075aa793144572d5fa3a398f0 127.0.0.1:8003
    replicates 964cfa1c2dcfe36b6d3c63637f0d57ccb568354e
    Can I set the above configuration? (type 'yes' to accept): yes
    >>> Nodes configuration updated
    >>> Assign a different config epoch to each node
    >>> Sending CLUSTER MEET messages to join the cluster
    Waiting for the cluster to join
    ...
    >>> Performing Cluster Check (using node 127.0.0.1:7001)
    M: 32f9819fc7d561bfa2b7189182200e86d9901b8a 127.0.0.1:7001
    slots:[0-5460] (5461 slots) master
    1 additional replica(s)
    S: aba9330f3e70f26a8af4ced1b672fbcc7bc62d78 127.0.0.1:8002
    slots: (0 slots) slave
    replicates cca0fbfa374bc175d481e68ee9ed13b65453e967
    S: 1b47b9e6e7a79523579b8d2ddcd5e708583ed317 127.0.0.1:8001
    slots: (0 slots) slave
    replicates 32f9819fc7d561bfa2b7189182200e86d9901b8a
    S: 254db0830cd764e075aa793144572d5fa3a398f0 127.0.0.1:8003
    slots: (0 slots) slave
    replicates 964cfa1c2dcfe36b6d3c63637f0d57ccb568354e
    M: cca0fbfa374bc175d481e68ee9ed13b65453e967 127.0.0.1:7002
    slots:[5461-10922] (5462 slots) master
    1 additional replica(s)
    M: 964cfa1c2dcfe36b6d3c63637f0d57ccb568354e 127.0.0.1:7003
    slots:[10923-16383] (5461 slots) master
    1 additional replica(s)
    [OK] All nodes agree about slots configuration.
    >>> Check for open slots...
    >>> Check slots coverage...
    [OK] All 16384 slots covered.

> OK，搭建完成！一条命令搞定。

&emsp;&emsp;建立连接后，启用密码认证，放开：

    masterauth
    requirepass
    #都要配置

## 监控

### 下载exporter

    wget https://github.com/oliver006/redis_exporter/releases/download/v1.9.0/redis_exporter-v1.9.0.linux-amd64.tar.gz

### 启动指标监控

    ./redis_exporter -redis.addr=redis://127.0.0.1:7001 -redis.password=xxx &
    # 切记，如果密码中有转义字符，一定要使用\转义，exporter监听了一个节点就会采集所有节点数据。
    # 默认是9121端口

### 开放端口

    firewall-cmd --add-port=9121/tcp --zone=public --permanent
    firewall-cmd --reload

### prometheus 配置

     - job_name: 'redis_exporter_targets'
    static_configs:
      - targets:
        - redis://172.18.11.139:7001
        - redis://172.18.11.139:7002
        - redis://172.18.11.140:7003
        - redis://172.18.11.140:8001
        - redis://172.18.11.141:8002
        - redis://172.18.11.141:8003
    metrics_path: /scrape
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 172.18.11.139:9121
    ## config for scraping the exporter itself
    - job_name: 'redis_exporter'
        static_configs:
          - targets:
            - 172.18.11.139:9121

### 规则配置

    groups:
    - name:  Redis
    rules: 
        - alert: RedisDown
        expr: redis_up  == 0
        for: 5m
        labels:
            severity: error
        annotations:
            summary: "Redis down (instance {{ $labels.instance }})"
            description: "Redis 挂了啊，mmp\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}"
        - alert: MissingBackup
        expr: time() - redis_rdb_last_save_timestamp_seconds > 60 * 60 * 24
        for: 5m
        labels:
            severity: error
        annotations:
            summary: "Missing backup (instance {{ $labels.instance }})"
            description: "Redis has not been backuped for 24 hours\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}"       
        - alert: OutOfMemory
        expr: redis_memory_used_bytes / redis_total_system_memory_bytes * 100 > 90
        for: 5m
        labels:
            severity: warning
        annotations:
            summary: "Out of memory (instance {{ $labels.instance }})"
            description: "Redis is running out of memory (> 90%)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}"
        - alert: ReplicationBroken
        expr: delta(redis_connected_slaves[1m]) < 0
        for: 5m
        labels:
            severity: error
        annotations:
            summary: "Replication broken (instance {{ $labels.instance }})"
            description: "Redis instance lost a slave\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}"
        - alert: TooManyConnections
        expr: redis_connected_clients > 1000
        for: 5m
        labels:
            severity: warning
        annotations:
            summary: "Too many connections (instance {{ $labels.instance }})"
            description: "Redis instance has too many connections\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}"       
        - alert: NotEnoughConnections
        expr: redis_connected_clients < 5
        for: 5m
        labels:
            severity: warning
        annotations:
            summary: "Not enough connections (instance {{ $labels.instance }})"
            description: "Redis instance should have more connections (> 5)\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}"
        - alert: RejectedConnections
        expr: increase(redis_rejected_connections_total[1m]) > 0
        for: 5m
        labels:
            severity: error
        annotations:
            summary: "Rejected connections (instance {{ $labels.instance }})"
            description: "Some connections to Redis has been rejected\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}"

### grafana 面板

    ID: 763


## 数据结构化

&emsp;&emsp;Redis最常见的使用场景便是**存储用户信息**。我们简单举个实例描述下Hash的应用场景。

&emsp;&emsp;用户ID为查找的key，存储的value用户对象包含姓名，年龄，生日等信息，如果用普通的key/value结构来存储，主要有以下2种存储方式：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/redis/1597211448.jpg)

&emsp;&emsp;第一种方式将用户ID作为查找key,把其他信息封装成一个对象以序列化的方式存储，这种方式的缺点是，增加了序列化/反序列化的开销，并且在需要修改其中一项信息时，需要把整个对象取回，并且修改操作需要对并发进行保护，引入CAS等复杂问题。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/redis/1597211590.jpg)

&emsp;&emsp;第二种方法是这个用户信息对象有多少成员就存成多少个key-value对儿，用用户ID+对应属性的名称作为唯一标识来取得对应属性的值，虽然省去了序列化开销和并发问题，但是用户ID为重复存储，如果存在大量这样的数据，内存浪费还是非常可观的。

&emsp;&emsp;那么Redis提供的Hash很好的解决了这个问题，Redis的Hash实际是内部存储的Value为一个HashMap，并提供了直接存取这个Map成员的接口，如下图：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/redis/1597211638.jpg)


&emsp;&emsp;也就是说，Key仍然是用户ID, value是一个Map，这个Map的key是成员的属性名，value是属性值，这样对数据的修改和存取都可以直接通过其内部Map的Key(Redis里称内部Map的key为field), 也就是通过 key(用户ID) + field(属性标签) 就可以操作对应属性数据了，既不需要重复存储数据，也不会带来序列化和并发修改控制的问题。很好的解决了问题