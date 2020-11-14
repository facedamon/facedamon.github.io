
# 摘要

&emsp;&emsp;今天我们要实现的是一致性哈希算法，一致性哈希算法是Geec从单节点走向分布式节点的一个重要的环节。

## 我该访问谁

&emsp;&emsp;对于分布式缓存来说，当一个节点接收到请求，如果该节点并没有存储缓存值，那么它面临的难题是，从谁那获取数据？自己，还是节点1，2，3，4...。假设包括自己在哪一共有10个节点，当一个节点接收到请求时，随机选择一个节点，由该节点从数据源获取数据。

&emsp;&emsp;假设第一次随机选取了节点1，节点1从数据源获取到数据的同时缓存该数据；那第二次，只有1/10的可能性再次选择节点1，有9/10的概率选择其它节点，如果选择了其它节点，就意味着需要再一次从数据源获取数据，一般来说，这个操作是很耗时的。这样做，一是缓存效率低，二是各个节点上存储着相同的数据，浪费了大量的存储空间。

&emsp;&emsp;那有什么办法，对于给定的key，每一次都选择同一个节点呢？使用hash算法也能够做到这一点。那把key的每一个字符的ACSII码加起来，再除以10取余数可以吗？当然可以，这可以认为是自定义的hash算法。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/geec/hash_select.jpg)

&emsp;&emsp;从上图可以看到，任意一个节点任意时刻请求查找键`Tom`对应的值，都会分配给节点2，有效地解决了上述问题。

## 节点数量变化了怎么办

&emsp;&emsp;简单求取Hash值解决了缓存性能问题，但是没有考虑节点数量变化的场景。假设，移除了其中一台节点，只剩下9个，那么之前`hash(key) % 10`变成了`hash(key) % 9`，也就意味着几乎缓存值对应的节点都发生了改变。即几乎所有的缓存值都失效了。节点在接收到对应的请求时，均需要重新去数据源获取数据，容易引起`缓存雪崩`

> 缓存雪崩：缓存在同一时刻全部失效，造成瞬间DB请求量大、压力骤增，引入雪崩。常因为缓存服务器宕机，或缓存设置了相同的过期时间引起。

&emsp;&emsp;**如何解决这个问题？**

# 算法原理

## 步骤

&emsp;&emsp;一致性哈希算法将key映射到2^32的空间中，将这个数字首位相连，形成一个环。

- 计算节点/机器(通常使用节点的名称、编号和IP地址)的哈希值，放置在环上
- 计算key的哈希值，放置在环上，**顺时针**寻找到的第一个节点，就是应选取的节点/机器

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/geec/add_peer.jpg)

&emsp;&emsp;环上有peer2, peer4, peer6三个节点， key11, key2, key27均映射到peer2, key23映射到peer4。此时，如果新增节点/机器peer8, 假设它新增为止如图所示，那么只有key27从peer2调整到peer8, 其余的映射均没有发生改变。

&emsp;&emsp;也就是说，一致性哈希算法，在新增/删除节点时，只需要重新定位该节点附近的一小部分数据，而不需要重新定位所有的节点，这就解决了上述的问题。

## 数据倾斜问题

&emsp;&emsp;如果服务器的节点过少，容易引起key的倾斜。例如上面例子中的peer2, peer4, peer6分布在环的上半部分，下半部分是空的。那么映射到环下半部分的key都会被分配给peer2，key过度向peer2倾斜，缓存节点间负载不均。

&emsp;&emsp;为了解决这个问题，引入了虚拟节点的概念，**一个真实节点对应多个虚拟节点**。

&emsp;&emsp;假设1个真是节点对应3个虚拟节点，那么peer1对应的虚拟节点是peer-1, peer1-2, peer1-3(通常以添加编号的方式实现)， 其余节点以相同的方式操作。

- 第一步，计算虚拟节点的hash值，放置在环上
- 第二步，计算key的hash值，在环上顺时针寻找到应选取的虚拟节点，例如是peer2-1，那么就对应真是节点peer2.

&emsp;&emsp;虚拟节点扩充了节点的数量，解决了节点较少的情况下数据容易倾斜的问题。而且代价非常小，只需要增加一个字典维护真是节点与虚拟节点的映射关系即可。

## Go实现

&emsp;&emsp;我们在geec目录下新建package `consistenthash`, 用来实现一致性hash算法

    package consistenthash

    import "hash/crc32"
    // Hash maps bytes to uint32
    type Hash func(data []byte) uint32
    // Map constains all hashed keys
    type Map struct {
            hash     Hash
            replicas int
            keys     []int
            hashMap  map[int]string
    }

    // New creates a Map instance
    func New(replicas int, fn Hash) *Map {
            m := &Map{
                    replicas: replicas,
                    hash:     fn,
                    hashMap:  make(map[int]string),
            }
            if m.hash == nil {
                    m.hash = crc32.ChecksumIEEE
            }
            return m
    }

- 定义了函数类型`Hash`, 采取依赖注入的方式，允许用于替换成自定义的Hash函数，也方便测试时替换，默认为`crc32.CHecksumIEEE`算法。
- `Map`是一致性哈希算法的主数据结构，包含4个成员变量：Hash函数，虚拟节点倍数replicas, 哈希环keys， 虚拟节点与真是节点的映射表hashMap，键是虚拟节点的哈希值，值是真实节点的名称。
- 构造函数`New()`允许自定义虚拟节点倍数和Hash函数

&emsp;&emsp;接下来，实现添加真实节点/机器的add方法

        func (m *Map) Add(keys ...string) {
                for _, key := range keys {
                        for i := 0; i < m.replicas; i++ {
                                hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
                                m.keys = append(m.keys, hash)
                                m.hashMap[hash] = key
                        }
                }
                sort.Ints(m.keys)
        }

- Add函数允许传入0或多个真实节点的名称
- 对每一个真实节点key，对应创建m.replicas个虚拟节点，虚拟节点的名称是：`strconv.Itoa(i) + key`,即通过添加编号的方式区分不同虚拟节点
- 使用m.hash计算虚拟节点的哈希值，使用`append(m.keys, hash)`添加到环上
- 在hashMap中增加虚拟节点和真实节点的映射关系
- 最后一步，环上的哈希值排序

        // Get gets the closest item in the hash to the provided key
        func (m *Map) Get(key string) string {
                if len(m.keys) == 0 {
                        return ""
                }
                hash := int(m.hash([]byte(key)))
                // Binary search for appropriate replica.
                idx := sort.Search(len(m.keys), func(i int) bool {
                        return m.keys[i] >= hash
                })
                return m.hashMap[m.keys[idx%len(m.keys)]]
        }

- 选择节点，第一步，计算key的哈希值
- 顺时针找到第一个匹配的虚拟节点的下标idx，从m.keys中获取到对应的哈希值。如果`idx == len(m.keys)`，说明应选择`m.keys[0]`因为m.keys是一个环状结构，所以采取取余数的方式处理这种情况
- 通过hashMap映射得到真实的节点

# 测试

        func TestHashing(t *testing.T) {
                hash := New(3, func(key []byte) uint32 {
                        i, _ := strconv.Atoi(string(key))
                        return uint32(i)
                })

                // 6 16 26 4 14 24 2 12 22
                hash.Add("6,", "4", "2")
                testCases := map[string]string{
                        "2":  "2",
                        "11": "2",
                        "23": "4",
                        "27": "2",
                }

                for k, v := range testCases {
                        if hash.Get(k) != v {
                                t.Errorf("Asking for %s, should have yield %s", k, v    )
                        }
                }

                // 8 18 28
                hash.Add("8")

                // 27 should now map to 8
                testCases["27"] = "8"

                for k, v := range testCases {
                        if hash.Get(k) != v {
                                t.Errorf("Asking for %s, should have yield %s", k, v    )
                        }
                }

&emsp;&emsp;如果要进行测试，那么我们需要明确地知道每一个传入地key地哈希值，那使用默认地`crc32.CHecksumIEEE`算法显然达不到目的。所以在这里使用了自定义地Hash算法。自定义地Hash算法只处理数字，传入字符串表示地数字，返回对应的数字即可。


- 一开始，有2/4/6三个真实节点，对应的虚拟节点的哈希值是02/12/22、04/14/24、06/16/26
- 那么用例2/11/23/27选择的虚拟节点分别是02/12/24/02，也就是真实节点2/2/4/2
- 添加一个真实节点8，对应的虚拟节点的哈希值是08/18/28，此时，用例27对应的虚拟节点从02变更为28，即真实节点8.