
# 摘要

&emsp;&emsp;对于集合，我认为关注的点主要有四点：

|关注点|结论|
|---|---|
|ArrayList是否允许null|允许|
|ArrayList是否允许重复数据|允许|
|ArrayList是否有序|有序|
|ArrayList是否线程安全|非同步|

## 添加元素

    public boolean add(E e) {
        ensureCapacity(size + 1);
        elementData[size++] = e;
        return true;
    }

&emsp;&emsp;先不去管第2行的ensureCapacity方法，这个方法是扩容用的，底层实际上在调用add方法的时候只是给elementData的某个位置添加了一个数据而已，用一张图表示的话是这样的：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151125163027140-67064210.png)

## 扩容

&emsp;&emsp;构造ArrayList的时候，默认的底层数组大小是10:

    public ArrayList() {this(10);}

&emsp;&emsp;那么有一个问题来了，底层数组的大小不够了怎么办？答案就是扩容，这也就是为什么一直说ArrayList的底层是基于`动态数组`实现的原因，动态数组的意思就是指底层的数组大小并不是固定的，而是根据添加的元素大小进行一个判断，不够的话就动态扩容，扩容的代码就在**ensureCapacity**里面：

    public void ensureCapacity(int minCapacity) {
        modCount++;
        int oldCapacity = elementData.length;
        if (minCapacity > oldCapacity) {
            Object oldData[] = elementData;
            int newCapacity = (oldCapacity * 3) / 2 + 1;
            if (newCapacity < minCapacity) {
                newCapacity = minCapacity;
                elementData = Arrays.copyOf(elementData, newCapacity);
            }
        }
    }

&emsp;&emsp;看到扩容的时候把元素组大小先乘以3，再除以2，最后加1。(新版JDK11 中：`int newCapacity = oldCapacity + (oldCapacity >> 1);`)可能有些人要问为什么？我们可以想：

1. 如果一次性扩容扩得太大，必然造成内存空间的浪费
2. 如果一次性扩容扩得不够，那么下一次扩容的操作必然比较快地会到来，这会降低程序运行效率，要知道扩容还是比较耗费性能的一个操作
3. 所以扩容扩多少，是JDK开发人员在时间、空间上做的一个权衡，提供出来的一个比较合理的数值。最后调用到的是Arrays的copyOf方法，将元素组里面的内容复制到新的数组里面去：

        public static <T,U> T[] copyOf(U[] original, int newLength, Class<? extends T[] newType>) {
            T[] copy = ((Object)newType == (Object)Object[].class)
            ? (T[]) new Object[newLength]
            : (T[]) Array.newInstance(newType.getComponentType(), newLength);
            System.arraycopy(original, 0, copy, 0,
                            Math.min(original.length, newLength));
            return copy;
        }

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151125163041468-1531821237.png)

## 删除元素

&emsp;&emsp;ArrayList支持两种删除方式：

1. 按照下标删除
2. 按照元素删除，这会删除ArrayList中与指定要删除的元素匹配的第一个元素

&emsp;&emsp;对于ArrayList来说，这两种删除的方法差不多，都是调用的下面一段代码：

    int numMoved = size - index - 1;
    if (numMoved > 0) {
        System.arraycopy(elementData, index+1, elementData, index, numMoved);
    }
    elementData[--size] = null;

&emsp;&emsp;其实做的事情就是两件:

1. 把指定元素后面位置的所有元素，利用System.arraycopy方法整体向前移动一个位置
2. 最后一个位置的元素指定为null，这样让gc可以去回收它

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151125163053187-1795426088.png)

## 插入元素

&emsp;&emsp;看一下ArrayList的插入操作，插入操作调用的也是add方法。

    public static void main(String[] args){
        List<String> list = new ArrayList<String>();
        list.add("111");
        list.add("222");
        list.add("333");
        list.add("444");
        list.add("555");
        list.add("666");
        list.add("777");
        list.add("888");
        list.add(2, "000");
        System.out.println(list);
    }

&emsp;&emsp;有一个地方不要搞错了，第11行的add方法的意思是，往第几个元素后面插入一个元素，像第12行就是往第二个元素后面插入一个000。看一下运行结果也证明了这一点：

    [111,222,000,333,444,555,666,777,888]

&emsp;&emsp;看看源码:

    public void add(int index, E element) {
        if (index > size || index < 0) {
            throw new IndexOutOfBoundException("Index: " + index+", Size: "+size);
        }
        ensureCapacity(size+1);
        System.arraycopy(elementData,index,elementData,index+1,size-index);
        elementData[index] = element;
        size++;
    }

&emsp;&emsp;看到插入的时候，按照指定位置，把从指定位置开始的所有元素利用System.arraycopy方法做一个整体的复制，**向后移动一个位置**（当然先要用ensureCapacity方法进行判断，加了一个元素之后数组会不会不够大），然后指定位置的元素设置为需要插入的元素，完成了一次插入的操作。用图表示这个过程是这样的：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20151125163125077-1317332042.png)

# 优缺点

|优点|缺点|
|---|---|
|ArrayList底层以数组实现，是一种随机访问模式，再加上它实现了RandomAccess接口，因此查找也就是get的时候非常快|删除元素的时候，涉及到一次元素复制，如果要复制的元素很多，那么就会比较耗费性能|
|ArrayList在顺序添加一个元素的时候非常方便，只是往数组里面添加了一个元素而已|插入元素的时候，涉及到一次元素复制，如果要复制的元素很多，那么就会比较耗费性能|

> 因此，ArrayList比较适合顺序添加、随机访问的场景。

# 为什么ArrayList的elementData是用transient修饰的？

    transient Object[] elementData;

&emsp;&emsp;不知道大家有没有想过，为什么elementData是使用transient修饰的呢？关于这个问题，说说我的看法。我们看一下ArrayList的定义：

    public class ArrayList<E> extends AbstractList<E>
            implements List<E>, RandomAccess, Cloneable, java.io.Serializable

&emsp;&emsp;看到ArrayList实现了Serializable接口，这意味着ArrayList是可以被序列化的，用transient修饰elementData意味着我不希望elementData数组被序列化。这是为什么？因为序列化ArrayList的时候，ArrayList里面的elementData未必是满的，比方说elementData有10的大小，但是我只用了其中的3个，那么是否有必要序列化整个elementData呢？显然没有这个必要，因此ArrayList中重写了writeObject方法：

    private void writeObejct(java.io.ObjectOutputStream s) throws java.io.IOException {
        int expectedModCount = modCount;
        s.defaultWriteObject();
        s.writeInt(elementData.length);
        for (int i = 0; i < size; i++) {
            s.writeObject(elementData[i]);
        }
        if (modCCount != expectedModCount) {
            throw new ConcurrentModificationnException();
        }
    }

&emsp;&emsp;每次序列化的时候调用这个方法，先调用defaultWriteObject()方法序列化ArrayList中的非transient元素，elementData不去序列化它，然后遍历elementData，只序列化那些有的元素，这样：

1. 加快了序列化的速度
2. 减小了序列化之后的文件大小

> 不失为一种聪明的做法，如果以后开发过程中有遇到这种情况，也是值得学习、借鉴的一种思路。

# 线程安全处理

- Collections.synchronizedList: 最常用的方法是通过 Collections 的 synchronizedList 方法将 ArrayList 转换成线程安全的容器后再使用。

        List<Object> list = Collections.synchronizedList(new ArrayList<Object>);

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/1610014324810.jpg)

&emsp;&emsp;从源码可以看出，该方法主要是通过加锁的方式实现同步。

- CopyOnWriteArrayList

        List<Object> list1 = new CopyOnWriteArrayList<Object>();

> CopyOnWriteArrayList实现了读写分离

&emsp;&emsp;关于CopyOnWriteArrayList的详细介绍，请参见[[Java | CopyOnWriteArrayList]](http://facedamon.github.io/post/java/collection/CopyOnWriteArrayList)

&emsp;&emsp;CopyOnWriteArrayList和Collections.synchronized是实现线程安全列表的两种方式。两种实现方式分别针对不同情况有不同的性能表现，其中CopyOnWriteArrayList的写操作性能较差，而多线程的读操作性能较好。而Collections.synchronizedList的写操作性能比CopyOnWriteArrayList在多线程操作下要好很多，而读操作因为采用了synchronzied关键字，其读操作性能并不如CopyOnWriteArrayList.