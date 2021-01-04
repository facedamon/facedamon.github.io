
# 淘汰算法简介

&emsp;&emsp;geec的缓存全部存储在内存中，内存是有限的，因此不可能无限制添加数据。假定，我们设置缓存能够使用的内存大小为N,那么在某一个时间点，添加了某一条缓存记录后，占用内存超过了N，这个时候就需要从缓存中移除某一条或多条数据了。那移除谁呢？我们肯定希望尽可能移除"没用"的数据，那如何判定数据"有用"还是"没用"呢?

1. **FIFO(First In First Out)**

&emsp;&emsp;先进先出，也就是淘汰缓存中最老的记录。FIFO认为，最早添加的记录，其不再被使用的可能性比刚添加的可能性大。这种算法的实现也非常简单，创建一个队列，新增记录添加到队尾，每次内存不够时，淘汰队首。但是很多场景下，部分记录虽然是最早添加但也最常被访问，而不得不因为呆的时间太久而被淘汰。这类数据会被频繁地添加进缓存，又被淘汰出去，导致缓存命中率降低。

2. **LFU(Least Frequently Used)**

&emsp;&emsp;最少使用，也就是淘汰缓存中访问频率最低地记录。LFU认为，如果数据过去被访问多次，那么将来被访问的频率也更高。LFU的实现需要维护一个按照访问次数排序的队列，每次访问，访问次数加1，队列重新排序，淘汰时选择访问次数最少的即可。LFU算法的命中率是比较高的，但缺点也非常明显，维护每个记录的访问次数，对内存的消耗很高；另外如果数据的访问模式发生变化，LFU需要较长的时间去适应，也就是说**LFU算法受历史数据的影响比较大**。例如某个数据历史访问次数很高，但在某个时间点之后几乎不再被访问，但因为历史访问次数过高，而迟迟不能被淘汰。

3. **LRU(Least Recently Used)**

&emsp;&emsp;最近最少使用，相对于仅考虑时间因素的FIFO和仅考虑访问频率的LFU，LRU算法可以认为是相对平衡的一种淘汰算法。LRU认为，如果数据最近被访问过，那么将来被访问的概率也会更高。LRU算法的实现非常简单，维护一个队列，如果某条记录被访问了，则移动到队尾，那么队首则是最近最少访问的数据，淘汰该条记录即可。

# LRU算法实现

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/lru.jpg)

&emsp;&emsp;这张图很好的表示了LRU算法最核心的2个数据结构
- 绿色是字典，存储键和值的映射关系。这样根据某个键查找对应的值的复杂度是O(1)，在字典中插入一条记录的复杂度也是O(1)。
- 红色是双向链表实现的队列。将所有的值放到双向链表中，这样，当访问到某个值时，将其移动到队尾的复杂度是O(1)，在队尾新增一条记录以及删除一条记录的复杂度均为O(1)

&emsp;&emsp;接下来我们创建一个包含字典和双向链表的结构体类型Cache，方便实现后续的操作。

		package lru
		
		import "container/list"
		
		// Cache is a LRU cache.
		// It is not safe for concurrent access.
		type Cache struct {
			maxBytes int64
			nbytes int64
			ll *list.List
			cache map[string]*list.Element
			// optional and executed when entry is purged
			OnEvicted func(key string, value Value)
		}
		
		type entry struct {
			key string
			value Value
		}
		
		// Value use Len to count how many bytes it takes
		type Value interface {
			Len() int
		}

- 这里我们直接使用Go标准库实现的双向链表`list.List`。
- 字典定义`map[Key]*list.Element`，键是任意字符，值是双向链表中对应节点的指针。
- `maxBytes`是允许使用的最大字节(0表示没有限制)，`nbytes`是当前已经使用的字节数，`OnEvicted`是回调函数。
- 键值对`entry`是双向链表节点的数据类型，在链表中保存每个值对应的key的好处在于，淘汰head节点时，需要用key从字典中删除对应的映射。

		// New is the Constructor of Cache
		func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
			return &Cache{
				maxBytes: maxBytes,
				ll: list.New(),
				cache: make(map[string]*list.Element),
				OnEvicted: onEvicted,
			}
		}

## 查找

&emsp;&emsp;查找主要有2个步骤，第一步是从字典中找到对应的双向链表的节点，第二步，将节点移动到队尾。

		// Get look ups a key`s value
		func (c *Cache) Get(key string) (value Value, ok bool) {
			if ele, ok := c.cache[key]; ok {
				c.ll.MoveToFront(ele)
				return ele.Value.(*entry).value, true
			}
			return
		}}

- 如果键对应的链表节点存在，则将对应节点移动到队尾，并返回查找到的值。

## 删除

		func (c *Cache) RemoveOldest() {
			ele := c.ll.Back()
			if ele != nil {
				c.ll.Remove(ele)
				kv := ele.Value.(*entry)
				delete(c.cache, kv.key)
				c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
				if c.OnEvicted != nil {
					c.OnEvicted(kv.key, kv.value)
				}
			}
		}}

- `c.ll.Back()`取到队首节点，从链表中删除。
- `delete(c.cache, kv.key)`，从字典`c.cache`中删除节点映射关系。
- 执行回调函数。

## 新增/修改

		func (c *Cache) Add(key string, value Value) {
			if ele, ok := c.cache[key]; ok {
				c.ll.MoveToFront(ele)
				kv := ele.Value.(*entry)
				c.nbytes += int64(value.Len()) - int64(kv.value.Len())
				kv.value = value
			} else {
				ele := c.ll.PushFront(&entry{key, value})
				c.cache[key] = ele
				c.nbytes += int64(len(key)) + int64(value.Len())
			}
			for c.maxBytes != 0 && c.maxBytes < c.nbytes {
				c.RemoveOldest()
			}
		}}

- 如果键存在，则更新对应节点的值，并将该节点移到队尾。
- 不存在，首先队尾添加新节点`&entry{key, value}`，再字典添加映射关系。
- 如果超过了设定的最大值`c.maxBytes`，则移除最少访问的节点。

&emsp;&emsp;为了测试方便，我们实现`Len()`用来获取添加多少条数据。

		func (c *Cache) Len() int {
		        return c.ll.Len()
		}

## Test

- **表格驱动测试**

- Get测试

		func (c *Cache) Add(key string, value Value) {
			if ele, ok := c.cache[key]; ok {
				c.ll.MoveToFront(ele)
				kv := ele.Value.(*entry)
				c.nbytes += int64(value.Len()) - int64(kv.value.Len())
				kv.value = value
			} else {
				ele := c.ll.PushFront(&entry{key, value})
				c.cache[key] = ele
				c.nbytes += int64(len(key)) + int64(value.Len())
			}
			for c.maxBytes != 0 && c.maxBytes < c.nbytes {
				c.RemoveOldest()
			}
		}}


- LRU测试

		func (c *Cache) Add(key string, value Value) {
			if ele, ok := c.cache[key]; ok {
				c.ll.MoveToFront(ele)
				kv := ele.Value.(*entry)
				c.nbytes += int64(value.Len()) - int64(kv.value.Len())
				kv.value = value
			} else {
				ele := c.ll.PushFront(&entry{key, value})
				c.cache[key] = ele
				c.nbytes += int64(len(key)) + int64(value.Len())
			}
			for c.maxBytes != 0 && c.maxBytes < c.nbytes {
				c.RemoveOldest()
			}
		}
