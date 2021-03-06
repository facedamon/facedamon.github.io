
# 摘要

&emsp;&emsp;第一次请求时将一些耗时操作的结果暂存，以后遇到相同的请求，直接返回暂存的数据。我想这是大部分人对于缓存的理解。在计算机系统中，缓存无处不在，比如我们访问一个网页，网页和引用的JS/CSS等静态文件，根据不同的策略，会缓存在浏览器本地或是CDN服务器，那么在第二次访问时，就会感觉网页加载的速度快了不少；比如微博的点赞数量，不可能每个人每次访问，都从数据库中查找所有点赞的记录再统计，数据库的操作相对来说是很耗时的，很难维持较大的流量，所以一点点赞这类数据是缓存在Redis这类内存缓存集群中的。

> 架构世界，缓存为王

&emsp;&emsp;缓存中最简单的莫过于存储在内存中的键值对缓存了。说到键值对，很容易想到的是字典(dict)类型，键值对缓存有什么问题？

1. **内存不够了怎么办?**

&emsp;&emsp;随机删掉好呢?还是按照时间顺序好呢?或者是没有其他更好的淘汰策略?不同数据的访问频率是不一样的，优先删除访问频率低的数据是不是更好呢?数据的访问频率有可能随着时间变化，那优先删除最近最少访问的数据可能是一个更好的选择。我们需要实现一个合理的**淘汰策略**

2. **并发写入冲突了怎么办?**

&emsp;&emsp;对缓存的访问，一般不可能是串行的。map是没有并发保护的，应对并发的场景，修改操作需要加锁。

3. **单机性能不够怎么办?**

&emsp;&emsp;单机资源是有限的，计算、存储等都是有限的。随着业务量和访问量的增加，单机很容易遇到瓶颈。如果**利用多台计算机的资源，并行处理提醒性能就要缓存应用能够支持分布式，这成为水平扩展**。与水平扩展相对应的是垂直扩展，即通过增加单个节点的计算、存储、带宽来提高系统的性能，硬件的成本和性能并非呈线性关系，大部分情况下，分布式系统是一个更优的选择。


## 关于geec

&emsp;&emsp;设计一个分布式缓存，需要考虑资源控制、淘汰策略、并发、分布式节点通信等方面的问题。而且，针对不同的应用场景，还需要在不同的特性之间权衡，例如，是否需要支持缓存更新?还是假定缓存在淘汰之前是不允许改变的。不同的权衡对应着不同的实现。

&emsp;&emsp;groupcache是golang版memcached，目的是在某些特定场合替代memcached。groupcache和memcache是同一个作者。无论是了解单机缓存还是分布式缓存，深入学习这个库的实现都是非常有意义的。

> 需要完成的基本功能

- 单机缓存和基于HTTP的分布式缓存
- 最近最少访问(Least Recently Used, LRU)缓存策略
- 防止缓存击穿
- 使用一致性哈希选择节点，实现负载均衡
- 使用protobuf优化节点间二进制通信

