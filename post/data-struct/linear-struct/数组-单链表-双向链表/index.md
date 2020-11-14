
> 转载自 https://www.cnblogs.com/skywang12345/

# 概要
&emsp;&emsp;线性表是一种线性结构，它是具有相同类型的n(n>=0)个元素的有序序列。本节介绍线性表的几个基本组成部分：数组、单向链表、双向链表。

# 数组
&emsp;&emsp;数组是有上届和下届的元素在上下届内是连续的。存储10,20,30,40,50的数组的示意图如下：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/数组存储示例.jpg)

- 数组的特点是: 
    - 连续的
    - 随机访问快
&emsp;&emsp;数组中稍微复杂一点的是多维数组和动态数组。对于C语言而言，多维数组本质上也是通过一维数组实现的。至于动态数组，是指数组的容量能动态增长的数组；对于C语言而言若要提供动态数组，需要手动实现；而对于Java而言，Collection集合中提供了ArrayList和Vector；对于Golang而言，默认提供了Slice。

# 单向链表
&emsp;&emsp;单向链表是链表的一种，它由节点组成，每个节点都包含下一个节点的指针。单链表的示意图如下:

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/单链表示例.jpg)

&emsp;&emsp;表头为空，表头的后继节点是"节点10"，"节点10"的后继节点是"节点20"。"节点20"的前驱节点是"节点10"

- 单链表删除节点

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/单链表删除节点示例.jpg)

&emsp;&emsp;删除"节点30"：

&emsp;&emsp;删除之前："节点20"的后继为"节点30"，而"节点30"的后继节点为"节点40"。

&emsp;&emsp;删除之后："节点20"的后继节点为"节点40"

- 单链表添加节点

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/单链表添加节点示例.jpg)

&emsp;&emsp;在"节点10"与"节点20"之间添加"节点15"

&emsp;&emsp;**添加之前**："节点10"的后继节点为"节点20"。

&emsp;&emsp;**添加之后**: "节点10"的后继节点为"节点15"，而"节点15"的后继节点为"节点20".

&emsp;&emsp;单链表的特点是: `节点的链接方向是单向的；相对于数组来说，单链表的随机访问速度慢，但是单链表删除、添加数据的效率很高`.

# 双向链表
&emsp;&emsp;双向链表是链表的一种。和单链表一样，双链表也是由节点组成，它的每个数据节点中都有两个指针，分别指向后继和前驱所以，从双向链表中的任意一个节点开始，都可以很方便地昂问它的前驱节点和后继节点，一般我们都构造双向循环链表。

&emsp;&emsp;双链表的示意图如下：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/双链表示例.jpg)

&emsp;&emsp;表头为空，表头的后继节点为"节点10"；"节点10"的后继节点是"节点20"，"节点20"的前驱节点是"节点10"；"节点20"的后继节点是"节点30"，"节点30"的前驱节点是"节点20"，。。。；末尾节点的后继节点是表头。

- 双链表删除节点

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/双链表删除节点示例.jpg)

&emsp;&emsp;**删除"节点30"**：

&emsp;&emsp;**删除之前**："节点20"的后继节点为"节点30"，"节点30"的前驱节点为"节点20"。"节点30"的后继节点为"节点40"，"节点40"的前驱节点为"节点30"。

&emsp;&emsp;**删除之后**："节点20"的后继节点为"节点40"，"节点40"的前驱节点为"节点20"

- 双链表添加节点

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/双链表添加节点.jpg)

&emsp;&emsp;在"节点10"与"节点20"之间添加"节点15"

&emsp;&emsp;**添加之前**："节点10"的后继节点为"节点20"，"节点20"的前驱节点为"节点10"

&emsp;&emsp;**添加之后**："节点10"的后继节点为"节点15"，"节点15"的前驱节点为"节点10"。"节点15"的后继节点为"节点20"，"节点20"的前驱节点为"节点15"

> 下面介绍双链表的实现，分别介绍Java/Golang两种实现


# Java实现双链表


        /**
        * Java实现的双向链表
        * 注： Java自带的集合包中有实现双向链表，路径是: java.util.LinkedList
        */
        public class DoubleLink<T> {

            //双向链表"节点"对应的结构体
            private class DNode<T> {
                public DNode prev;
                public DNode next;
                public T value;

                public DNode(T value, DNode prev, DNode next) {
                    this.value = value;
                    this.prev = prev;
                    this.next = next;
                }
            }

            //表头
            private DNode<T> mHead;
            //节点个数
            private int mCount;

            /**
            * 创建表头，表头没有存储数据
            * 初始化节点个数为0
            */
            public DoubleLink() {
                mHead = new DNode<T>(null, null, null);
                mHead.prev = mHead.next = mHead;
                mCount = 0;
            }

            /**
            * 
            * @return 节点数目
            */
            public int size() {
                return mCount;
            }

            /**
            * 
            * @return 链表是否为空
            */
            public boolean isEmpty() {
                return mCount == 0;
            }

            /**
            * 
            * @param index
            * @return index位置的节点
            */
            public DNode<T> getNode(int index) {
                if (index < 0 || index >= mCount) {
                    throw new IndexOutOfBoundsException();
                }
                //正向查找
                if (index <= mCount / 2) {
                    DNode<T> node = mHead.next;
                    for (int i = 0; i < index; i++) {
                        node = node.next;
                    }
                    return node;
                }

                //反向查找
                DNode<T> rnode = mHead.prev;
                int rindex = mCount - index - 1;
                for (int j = 0; j < rindex; j++) {
                    rnode = rnode.prev;
                }
                return rnode;
            }

            /**
            * 
            * @param index
            * @return index位置的节点的值
            */
            public T get(int index) {
                return getNode(index).value;
            }

            /**
            * 
            * @return 获取第一个节点
            */
            public T getFirst() {
                return getNode(0);
            }

            /**
            * 
            * @return 获取最后一个节点的值
            */
            public T getLast() {
                return getNode(mCount - 1).value;
            }

            /**
            * 将节点插入到第index位置之前
            * @param index
            * @param t
            */
            public void insert(int index, T t) {
                //表头
                if (index == 0) {
                    DNode<T> node = new DNode<T>(t, mHead, mHead.next);
                    mHead.next.prev = node;
                    mHead.next = node;
                    mCount++;
                    return;
                }

                DNode<T> inode = getNode(index);
                DNode<T> tnode = new DNode<T>(t, inode.prev, inode);
                inode.prev.next = tnode;
                inode.prev = tnode;
            }

            /**
            * 将节点插入第一个节点处
            * @param t
            */
            public void insertFirst(T t) {
                insert(0, t);
            }

            /**
            * 将节点追加到链表的末尾
            * @param t
            */
            public void appendLast(T t) {
                DNode<T> node = new DNode<T>(t, mHead.prev, mHead);
                mHead.prev.next = node;
                mHead.prev = node;
                mCount++;
            }

            /**
            * 删除index位置的节点
            * @param index
            */
            public void del(int index) {
                DNode<T> inode = getNode(index);
                inode.prev.next = inode.next;
                inode.next.prev = inode.prev;
                inode = null;
                mCount--;
            }

            /**
            * 删除第一个节点
            */
            public void delFirst() {
                del(0);
            }

            /**
            * 删除最后一个节点
            */
            public void delLast() {
                del(mCount - 1);
            }

            /**
            * 打印链表
            */
            @Override
            public String toString() {
                StringBuffer str = new StringBuffer("[");
                for (int i = 0; i < mCount; i++) {
                    str.append(get(i)).append(",");
                }
                str.deleteCharAt(str.length()-1).append("]");
                return str.toString();
            }
        }



        /**
        * 
        */
        public class DlinkTests {

            public static void test_int() {
                System.out.println("\n----test_int----");
                DoubleLink<Integer> dlnk = new DoubleLink<>();
                dlnk.insert(0, 20);
                dlnk.appendLast(10);
                dlnk.insertFirst(30);

                System.out.printf("isEmpty()=%b\n", dlnk.isEmpty());
                System.out.printf("size()=%d\n", dlnk.size());
                System.out.println(dlnk.toString());
            }

            public static void test_string() {
                System.out.println("\n----test_string----");
                DoubleLink<String> dlnk = new DoubleLink<>();
                dlnk.insert(0, "ten");
                dlnk.appendLast("twenty");
                dlnk.insertFirst("thirty");

                System.out.printf("isEmpty()=%b\n", dlnk.isEmpty());
                System.out.printf("size()=%d\n", dlnk.size());
                System.out.println(dlnk.toString());
            }

            private static class Student {
                private int id;
                private String name;

                public Student(int id, String name) {
                    this.id = id;
                    this.name = name;
                }

                @Override
                public String toString() {
                    return "["+id+", "+name+"]";
                }
            }

            private static Student[] students = new Student[]{
                new Student(10, "sky"),
                new Student(20, "joy"),
                new Student(30, "vic"),
                new Student(40, "dan"),
            };

            private static void test_object() {
                System.out.println("\n----test_object----");
                DoubleLink<Student> dlnk = new DoubleLink<>();
                dlnk.insert(0, students[1]);
                dlnk.appendLast(students[0]);
                dlnk.insertFirst(students[2]);

                System.out.printf("isEmpty()=%b\n", dlnk.isEmpty());
                System.out.printf("size()=%d\n", dlnk.size());
                System.out.println(dlnk.toString());
            }

            public static void main(String[] args) {
                //test_int();// isEmpty = false; size() = 3; [30, 20, 10]
                //test_string(); // isEmpty = false; size() = 3; [thirty, ten, twenty]
                test_object(); // isEmpty = false; size() = 3; [[30, vic], [20, joy], [10, sky]]
            }
            
        }


# Golang实现双链表

- CopyLeft https://github.com/emirpasic/gods


        package doublylinkedlist

        import (
            "fmt"
            "strings"
            "github.com/facedamon/gods/utils"
        )

        type element struct {
            value interface{}
            prev  *element
            next  *element
        }

        // List holds the elements, where each element points to
        // the next and previous element
        type List struct {
            first *element
            last  *element
            size  int
        }

        // New instantiates a new list and adds the passed values, if any, to the list
        func New(values ...interface{}) *List {
            list := &List{}
            if len(values) > 0 {
                list.Add(values...)
            }
            return list
        }

        // Add appends a value at the end of the list
        func (list *List) Add(values ...interface{}) {
            for _, value := range values {
                newElement := &element{value: value, prev: list.last}
                if list.size == 0 {
                    list.first = newElement
                    list.last = newElement
                } else {
                    list.last.next = newElement
                    list.last = newElement
                }
                list.size++
            }
        }

        // Prepend prepends a values
        func (list *List) Prepend(values ...interface{}) {
            for v := len(values) - 1; v >= 0; v-- {
                newElement := &element{value: values[v], next: list.first}
                if list.size == 0 {
                    list.first = newElement
                    list.last = newElement
                } else {
                    list.first.prev = newElement
                    list.first = newElement
                }
                list.size++
            }
        }

        func (list *List) withinRange(index int) bool {
            return index >= 0 && index < list.size
        }

        // Clear removes all elements from the list
        func (list *List) Clear() {
            list.size = 0
            list.last = nil
            list.first = nil
        }

        // Get returns the element at index.
        // Second return parameter is true if index is within bounds of the array and array is not empty
        func (list *List) Get(index int) (interface{}, bool) {
            if !list.withinRange(index) {
                return nil, false
            }
            if list.size-index < index {
                element := list.last
                for e := list.size - 1; e != index; e, element = e-1, element.prev {
                }
                return element.value, true
            }
            element := list.first
            for e := 0; e != index; e, element = e+1, element.next {
            }
            return element.value, true
        }

        // Remove removes the element at the given index from the list.
        func (list *List) Remove(index int) {
            if !list.withinRange(index) {
                return
            }
            if list.size == 1 {
                list.Clear()
                return
            }
            var element *element
            //二分法
            if list.size-index < index {
                element = list.last
                for e := list.size - 1; e != index; e, element = e-1, element.prev {
                }
            } else {
                element = list.first
                for e := 0; e != index; e, element = e+1, element.next {
                }
            }
            if element == list.first {
                list.first = element.next
            }
            if element == list.last {
                list.last = element.prev
            }
            if element.prev != nil {
                element.prev.next = element.next
            }
            if element.next != nil {
                element.next.prev = element.prev
            }

            element = nil
            list.size--
        }

        // Contains check if values (one or more) are present int the set.
        // All values have to be present in the set for the method to return true.
        // Performance time complexity of n^2.
        // Returns true if no arguments are passed at all.
        func (list *List) Contains(values ...interface{}) bool {
            if len(values) == 0 {
                return true
            }
            if list.size == 0 {
                return false
            }
            for _, value := range values {
                found := false
                for element := list.first; element != nil; element = element.next {
                    if element.value = value {
                        found = true
                        break
                    }
                }
                if !found {
                    return false
                }
            }
            return true
        }

        // Values returns all elements in the list.
        func (list *List) Values() []interface{} {
            values := make([]interface{}, list.size, list.size)
            for e, element := 0, list.first; element != nil; e, element = e+1, element.next {
                values[e] = element.value
            }
            return values
        }

        // IndexOf returns index of provided element
        func (list *List) IndexOf(value interface{}) int {
            if list.size == 0 {
                return -1
            }
            for index, element := range list.Values() {
                if element == value {
                    return index
                }
            }
            return -1
        }

        // Empty returns true if list does not contain any elements.
        func (list *List) Empty() bool {
            return list.size == 0
        }

        // Sort sortts values using
        func (list *List) Sort(comparator utils.comparator) {
            if list.size < 2 {
                return
            }
            values := list.Values()
            utils.Sort(values, comparator)
            list.Clear()
            list.Add(values...)
        } 

        // Insert inserts values at specified index position shifting the value at that position and any subsequent elements to the right.
        // Does not do anything if position is negative or bigger than list`s size
        // Note: position equal to list`s size is valid
        func (list *List) Insert(index int, values ...interface{}) {
            if !list.withinRange((index)) {
                if index == list.size {
                    list.Add(values...)
                }
                return
            }
            list.size += len(values)

            var beforeElement *element
            var foundElement *element

            // determine traversal direction, last to first or first to last
            if list.size - index < index {
                foundElement = list.last
                for e := list.size - 1; e != index; e, foundElement = e-1, foundElement.prev {
                    beforeElement = foundElement.prev
                }
            } else {
                foundElement = list.first
                for e := 0; e != index; e, foundElement = e+1, foundElement.next {
                    beforeElement = foundElement
                }
            }

            if foundElement == list.first {
                oldNextElement := list.first
                for i, v := range values {
                    newElement := &element{value: v}
                    if i == 0 {
                        list.first = newElement
                    } else {
                        newElement.prev = beforeElement
                        beforeElement.next = newElement
                    }
                    beforeElement = newElement
                }
                oldNextElement.prev = beforeElement
                beforeElement.next = oldNextElement
            } else {
                oldNextElement := beforeElement.next
                for _, v := range values {
                    newElement := &element{value: v}
                    newElement.prev = beforeElement
                    beforeElement.next = newElement
                    beforeElement = newElement
                }
                oldNextElement.prev = beforeElement
                beforeElement.next = oldNextElement
            }
        }

        // String returns a string representation of container
        func (list *List) String() string {
            str := "DoubleLinkedList\n"
            values := []string{}
            for element := list.first; element != nil; element = element.next {
                values = append(values, fmt.Sprintf("%f", element.value))
            }
            str += strings.Join(values, ",")
            return str
        }
