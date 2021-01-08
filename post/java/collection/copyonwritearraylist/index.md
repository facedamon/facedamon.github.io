
# 摘要

1. CopyOnWriteArrayList位于java.util.concurrent包下，可想而知，这个类是为并发而设计的。
2. CopyOnWriteArrayList，顾名思义，Write的时候总是要Copy，也就是说**对于CopyOnWriteArrayList，任何可变操作(add,set,remove等）都是伴随复制这个动作的，**后面会解读底层实现机制。

&emsp;&emsp;对于集合，我认为关注的点主要有四点：

|关注点|结论|
|---|---|
|CopyOnWriteArrayList是否允许null|允许|
|CopyOnWriteArrayList是否允许重复数据|允许|
|CopyOnWriteArrayList是否有序|有序|
|CopyOnWriteArrayList是否线程安全|线程安全|

## 添加元素

    public class CopyOnWriteArrayList<E>
        implements List<E>, RandomAccess, Cloneable, java.io.Serializable {
        private static final long serialVersionUID = 8673264195747942595L;

        /** The lock protecting all mutators */
        transient final ReentrantLock lock = new ReentrantLock();

        /** The array, accessed only via getArray/setArray. */
        private volatile transient Object[] array;
        ...

        public CopyOnWriteArrayList() {
            setArray(new Object[0]);
        }

        final void setArray(Object[] a) {
            array = a;
        }
    }

&emsp;&emsp;看到，底层还是一个Object[] array,**默认大小为0**。来看看add的源码：

    public boolean add(E e) {
        final ReentrantLock = this.lock;
        lock.lock();
        try {
            Object[] elements = getArray();
            int len = elements.length;
            Object[] newElements = Arrays.copyOf(elements, len + 1);
            newElements[len] = e;
            setArray(newElements);
            return true;
        } finally {
            lock.unlock();
        }
    }

&emsp;&emsp;看一下逻辑图：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151206205856722-1595736205.png)

&emsp;&emsp;由于array被`transient, volatile`修饰，保证了内存的非重排性。

1. 加锁
2. 拿到原数组，得到新数组大小(原数组大小+1)，实例化出一个新的数组来
3. 把原数组的元素复制到新数组中去
4. 新数组最后一个位置为待添加的元素(因为新数组的大小是按照原数组大小+1来的)
5. 把Object array引用指向新数组
6. 解锁

&emsp;&emsp;另外，插入、删除、修改操作也都是一样的，每一次的操作都是以对Object[] array进行一次复制为基础的，如果上面的流程看懂了，那么研究插入、删除、修改的源代码就不难。

|优点|缺点|
|---|---|
|绝对的线程安全|修改代价十分昂贵，每次修改都伴随着一次数组复制|

&emsp;&emsp;CopyOnWriteArrayList这个并发组件，其实反映的是两个十分重要的分布式理念：

1. **读写分离**

&emsp;&emsp;我们读取CopyOnWriteArrayList的时候读取的是CopyOnWriteArrayList中的Object[] array，但是修改的时候，操作的是一个新的Object[] array，读和写操作的不是同一个对象，这就是读写分离。这种技术数据库用的非常多，在高并发下为了缓解数据库的压力，即使做了缓存也要对数据库做读写分离，读的时候使用读库，写的时候使用写库，然后读库、写库之间进行一定的同步，这样就避免同一个库上读、写的IO操作太多

2. **最终一致**

&emsp;&emsp;对CopyOnWriteArrayList来说，线程1读取集合里面的数据，未必是最新的数据。因为线程2、线程3、线程4四个线程都修改了CopyOnWriteArrayList里面的数据，但是线程1拿到的还是最老的那个Object[] array，新添加进去的数据并没有，所以线程1读取的内容未必准确。不过这些数据虽然对于线程1是不一致的，但是对于之后的线程一定是一致的，它们拿到的Object[] array一定是三个线程都操作完毕之后的Object array[]，这就是最终一致。最终一致对于分布式系统也非常重要，它通过容忍一定时间的数据不一致，提升整个分布式系统的可用性与分区容错性。当然，最终一致并不是任何场景都适用的，像火车站售票这种系统用户对于数据的实时性要求非常非常高，就必须做成强一致性的。

&emsp;&emsp;最后总结一点，随着CopyOnWriteArrayList中元素的增加，CopyOnWriteArrayList的修改代价将越来越昂贵，因此，

> CopyOnWriteArrayList适用于读操作远多于修改操作的并发场景中。