
# 摘要

&emsp;&emsp;上一篇文章了解了HashMap的实现原理，讲到了HashMap不是线程安全的。那么HashMap在多线程环境下又会有什么问题呢？先来看一段代码：

    public class HashMapThread extends Thread {
        private static AtomicInteger ai = new AtomicInteger(0);
        private static Map<Integer, Integer> map = new HashMap<Integer, Integer>(1);

        @Override
        public void run() {
            while (ai.get() < 100000) {
                map.put(ai.get(), ai.get());
                ai.incrementAndGet();
            }
        }
    }

&emsp;&emsp;这个线程的作用很简单，给AtomicInteger不断自增并写入HashMap中，其中AtomicInteger和HashMap都是全局共享的，也就说所有线程操作的都是同一个AtomicInteger和HashMap。开5个线程操作一下run方法中的代码：

    public static void main(String[] args) {
        HashMapThread hmt0 = new HashMapThread();
        HashMapThread hmt1 = new HashMapThread();
        HashMapThread hmt2 = new HashMapThread();
        HashMapThread hmt3 = new HashMapThread();
        HashMapThread hmt4 = new HashMapThread();
        hmt0.start();
        hmt1.start();
        hmt2.start();
        hmt3.start();
        hmt4.start();
    }

&emsp;&emsp;多运行几次死循环就出来了，大概运行7、8次的样子，其中有几次是数组下标越界异常。这里要提一点，多线程环境下代码会出现问题并不意味着多线程环境下一定会出现问题，但是只要出现了问题，或者是死锁、或者是死循环，那么你的项目除了重启就没有什么别的办法了，所以代码的线程安全性在开发、评审的时候必须要考虑到。

## 正常的扩容过程

&emsp;&emsp;先来看一下HashMap一次正常的扩容过程。简单一点，假设有三个经过了最终rehash得到的数字，分别是5 7 3，HashMap的table也只有2，那么HashMap把这三个数字put进数据结构之后应该是这个样子：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212170803809-1228904061.png)

&emsp;&emsp;这应该很好理解。看一下resize的代码：

    void addEntry(int hash, K key, V value, int bucketIndex) {
        Entry<K,V> e = table[bucketIndex];
        table[bucketIndex] = new Entry<K,V>(hash, key, value, e);
        if (size++ >= thresold) {
            resize(2 * table.length);
        }
    }

    void resize(int newCapacity) {
        Entry[] oldTable = table;
        int oldCapacity = oldTable.lenngth
        if (oldCapacity == MAXIMUM_CAPACITY) {
            threshold = Integer.MAX_VALUE;
            return;
        }
        Entry[] newTable = new Table[newCapacity];
        transfer(newTable);
        table = newTable;
        threshold = (int)(newCapacity * loadFactory);
    }

    void transfer(Entry[] newTable) {
        Entry[] src = table;
        int newCapacity = newTable.length;
        for (int j = 0; j < src.length; j++) {
            Entry<K,V> e = src[j];
            if (e != null) {
                src[j] = null;
                do {
                    Entry<K,V> next = e.next;
                    int i = indexFor(e.hash, newCapacity);
                    e.next = newTable[i];
                    // 就是这里
                    newTable[i] = e;
                    e = next;
                } while (e != null);
            }
        }
    }

&emsp;&emsp;总结一下这三段代码，HashMap一次扩容的过程应该是：

1. 取当前table的2被作为新table的大小
2. 根据算出的新table的大小new出一个新的Entry数组来，名为newTable
3. 轮询原table的每一个位置，将每个位置上连接的Entry，算出在新table上的位置，并以链表形式连接
4. 原table上的所有Entry全部轮询完毕之后，意味着原table上面的所有Entry预警迁移到了新的table上，HashMap中的table指向newTable

&emsp;&emsp;这样就完成了一次扩容，用图表示是这样的：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212171001169-218100779.png)

## 扩容导致的死循环

&emsp;&emsp;既然是扩容导致的死循环，那么继续看扩容的代码的话。现在两个线程，线程A和线程B。假设第30行执行完毕，线程A切换，那么对于线程A而言，是这样的：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212184247794-1773995400.png)

&emsp;&emsp;CPU切换到线程B运行，线程B将整个扩容过程全部执行完毕，于是就形成了：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212184340794-575118916.png)

&emsp;&emsp;此时CPU切换到线程A上，执行第29行~第36行的do...while...循环，首先放置3这个Entry：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212185029137-1615739989.png)

&emsp;&emsp;我们必须知道，**由于线程B已经执行完毕，因此根据Java内存模型，现在table里面所有的Entry都是最新的，也就是7的next是3，3的next是null。**3放置到table[3]的位置上了，下面的步骤是：

1. e=next,即e=7
2. 判断e不等于null，循环继续
3. next=e.next,即next=7的next，也就是3
4. 放置7这个Entry

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212185651591-1282375682.png)

&emsp;&emsp;放置完7之后，继续运行代码：

1. e=next,也就是e=3;
2. 判断e不等于null，循环继续
3. next=e.next,即3的next，也就是null
4. 放置3这个Entry

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212185956028-948946343.png)

&emsp;&emsp;3移到table[3]上去了，3的next指向7，由于原先7的next指向3，这样就形成了一个死循环。

&emsp;&emsp;此时执行35行的e=next,那么e=null,循环终止。尽管此次循环确实结束了，但是后面的操作，只要涉及轮询HashMap数据结构的，无论是迭代还是扩容，都将在table[3]这个链表处出现死循环。这也就是前面的死循环堆栈出现的原因。

### 3 5 7 又会怎么样

&emsp;&emsp;可能有人觉得上面的数字5 7 3太巧了，像是专门为了产生HashMap的死循环而故意选择的数字。

&emsp;&emsp;这个问题，可以引用《分布式一致性原理》的结论：`任何在多线程下可能发生的错误场景最终一定会发生`。

&emsp;&emsp;现在看看正常的顺序3 5 7会发生什么问题。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212201914200-984546904.png)

&emsp;&emsp;这是扩容前数据结构中的内容，扩容之后正常的应该是：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212202024200-201515113.png)

&emsp;&emsp;现在在多线程下遇到问题了，某个线程先放7:

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212202127091-1874097507.png)

&emsp;&emsp;再接着放5:

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151212202205653-1833003565.png)

&emsp;&emsp;由于5的next此时为null，因为扩容操作结束， 3 5 7造成的结果就是元素丢失。

# 如何解决

&emsp;&emsp;把一个线程非安全的集合作为全局共享，本身就是一种错误的做法，并发下一定会产生错误。所以，解决这个问题的办法很简单，有两种：

1. 使用Collections.synchronizedMap(Map<K,V>m)方法把HashMap变成一个线程安全的Map
2. 使用ConcurrentHashMap这个线程安全的Map