
# 摘要

&emsp;&emsp;LinkedList是基于链表实现的，所以先讲解一下什么是链表。链表原先是C/C++的概念，是一种线性的存储结构，意思是将要存储的数据存在一个存储单元里面，这个存储单元里面除了存放有待存储的数据以外，还存储有其下一个存储单元的地址（下一个存储单元的地址是必要的，有些存储结构还存放有其前一个存储单元的地址），每次查找数据的时候，通过某个存储单元中的下一个存储单元的地址寻找其后面的那个存储单元。

&emsp;&emsp;LinkedList是一种`双向链表`，双向链表我认为有两点含义：

1. 链表中任意一个存储单元都可以通过向前或者向后寻址的方式获取到其前一个存储单元和其后一个存储单元
2. 链表的尾节点的后一个节点是链表的头结点，链表的头结点的前一个节点是链表的尾节点

&emsp;&emsp;LinkedList既然是一种双向链表，必然有一个存储单元，看一下LinkedList的基本存储单元，它是LinkedList中的一个内部类：

    private static class Node<E> {
        E item;
        Node<E> next;
        Node<E> prev;

        Node(Node<E> prev, E element, Node<E> next) {
            this.item = element;
            this.next = next;
            this.prev = prev;
        }
    }

&emsp;&emsp;看到LinkedList的Node中的"E item"，就是它真正存储的数据。"Node<E> next"和"Node<E> prev"表示的就是这个存储单元的前一个存储单元的引用地址和后一个存储单元的引用地址。

&emsp;&emsp;对于集合，我认为关注的点主要有四点：

|关注点|结论|
|---|---|
|LinkedList是否允许null|允许|
|LinkedList是否允许重复数据|允许|
|LinkedList是否有序|有序|
|LinkedList是否线程安全|非线程安全|

## 获取元素

&emsp;&emsp;我们知道随机读取元素不是LinkedList所擅长的，读取效率比起ArrayList也低得多，那么我来看一下为什么

        public E get(int index) {
            checkElementIndex(index);
            return node(index).item;
        }

        // 二分法查找
        Node node(int index) {
            if (index < (size >> 1)) {
                Node<E> x = first;
                for (int i = 0; i < index; i++) {
                    x = x.next;
                }
                return x;
            } else {
                Node<E> x = last;
                for (int i = size - 1; i > index; i--) {
                    x = x.prev;
                }
                return x;
            }
        }

&emsp;&emsp;从上述代码中我们可以看到get(int index)方法是通过node(int index)来实现的，它的实现机制是：

&emsp;&emsp比较传入的索引参数index与集合长度size/2，如果是index小，那么从第一个顺序循环，直到找到为止；如果index大，那么从最后一个倒序循环，直到找到为止。也就是说越靠近中间的元素，调用get(int index)方法遍历的次数越多，效率也就越低，而且随着集合的越来越大，get(int index)执行性能也会指数级降低。**因此在使用LinkedList的时候，我们不建议使用这种方式读取数据，可以使用getFirst()，getLast()方法，将直接用到类中的first和last变量。**

## 添加元素

&emsp;&emsp;LinkedList插入、删除基本都在首或者尾部进行，效率非常高。

    public boolean add(E e) {
        linkLast(e);
        return true;
    }

    void linkLast(E e) {
        final Node<E> l = last;
        final Node<E> newNode = new Node<E>(l, e, null);
        last = newNode;
        if (l == null) {
            first = newNode;
        } else {
            l.next = newNode;
        }
        size++;
        mmodCount++;
    }

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/1610022167831.jpg)

&emsp;&emsp;同理，我们看看在指定位置如何插入一个元素。

    public void add(int index, E element) {
        checkPositionIndex(index);
        if (index == size) {
            linkLast(element);
        } else {
            linkBefore(element, node(index));
        }
    }

    void linkBefore(E e, Node<E> succ) {
        final Node<E> pred = succ.prev;
        final Node<E> newNode = newNode(pred, e, succ);
        succ.prev = newNode;
        if (pred == null) {
            first = newNode;
        } else {
            pred.next = newNode;
        }
        size++;
        modCount++;
    }

&emsp;&emsp;其实从代码中看到和add(E e)的代码实现没有本质区别，都是通过新建一个Node实体，同时指定其prev和next来实现，不同点在于需要调用node(int index)通过传入的index来定位到要插入的位置，这个也是比较耗时的.

> LinkedList插入效率高是相对的

## 删除元素

&emsp;&emsp;我们看一下指定位置和指定元素删除:

    // 删除指定元素
    public boolean remove(Object o) {
        if (o == null) {
            for (Node<E> x = first; x != null; x = x.next) {
                if (x.item == null) {
                    unlink(x);
                    return true;
                }
            }
        } else {
            for (Node<E> x = first; x != null; x = x.next) {
                if (o.equals(x.item)) {
                    unlink(x);
                    return true;
                }
            }
        }
        return false;
    }

    // 删除指定位置元素
    public E remmove(int index) {
        checkElementIndex(index);
        return unlink(node(index));
    }

    // 删除某节点，并将该节点的上一个节点和下一个节点关联起来
    E unlink(Node<E> x) {
        final E element = x.item;
        final Node<E> next = x.next;
        final Node<E> prev = x.prev;

        // 头节点
        if (prev == null) {
            first = next;
        } else {
            prev.next = next;
            x.prev = null;
        }

        // 尾节点
        if (next == null) {
            last = prev;
        } else {
            next.prev = prev;
            x.next = null;
        }
        x.item = null;
        size--;
        modCount++;
        return element;
    }

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/1610023398927.jpg)

## 遍历

> 普通for非常慢，foreach和Iterator差不多，因为LinkedList内部foreach也调用了Iterator

&emsp;&emsp;**LinkedList其实是一个集栈、队列于一体的数据结构**。特别是最新的JDK，其中改动最大的是，取消了双向循环链表，该用普通双向链表，并将栈、队列常用API加入其中，毕竟链表这种结构对首部和尾部的操作是最高效的。那么插入、删除并非绝对的高效，而是相对高效，因为内部调用了node函数。

&emsp;&emsp;关于链表的详细介绍，请参考数据结构章节。

## 多线程

&emsp;&emsp;单链表有两个选择：`ConcurrentLinkedQueue`和`LinkedBlockingQueue`。

&emsp;&emsp;ConcurrentLinkedQueue是通过CAS实现同步，LinkedBlockingQueue通过ReentrantLock锁实现同步。