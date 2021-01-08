
# 集合框架图

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/249993-20170119151138921-1550766941.jpg)

- 所有集合类都位于`java.util`包下。Java集合类主要由两个接口派生：Collection和Map，Collection和Map是Java集合框架的根接口，这两个接口又包含了一些子接口或实现。
- Collection接口是一组允许重复的对象，是高度抽象出来的。它包含了集合的基本操作和属性
- Set接口继承Collection，集合元素不重复。Set的实现类有HashSet、TreeSet。HashSet依赖于HashMap，TreeSet依赖于TreeMap。
- List接口继承Collection，允许重复，维护元素插入顺序。List的实现类有LinkedList,ArrayList,Vector,Stack.
- Map接口是键-值对象，AbstractMap是一个抽象类，它实现了Map接口中的大部分API。而HashMap、TreeMap、WeakHashMap都继承于AbstractMap。Hashtable虽然继承于Dictionary，但它实现了Map接口
- Set、List和Map可以看作集合的三大基类：
  - List集合是有序集合，集合中的元素可以重复，访问集合中的元素可以根据元素的索引来访问
  - Set集合是无序集合，集合中的元素不可以重复，访问集合中的元素只能根据元素本身来访问
  - Map集合中保存k-v形式的元素，访问时只能根据每项元素的key来访问其value
- Iterator接口。它是遍历集合的工具，即我们通常通过Iterator迭代器来遍历集合。一般集合类中都有一个内部类来实现Iterator接口(实现iterator()方法，返回iterator对象)，外部类实现Iterable接口。
- Arrays和Collections。它们是操作数组、集合的两个工具类。

## Collection接口

&emsp;&emsp;Collection接口是处理对象集合的根接口，其中定义了很多对元素进行操作的方法。Collection接口有三个主要的子接口：`List`、`Set`、`Queue`。**Map不是Collection的子接口**

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/1609929325038.jpg)

- add()添加一个元素到集合中
- addAll()将指定集合中的所有元素添加到集合中
- contains()检测集合中是否包含指定的元素
- toArray()返回一个表示集合的数组
- iterator返回一个Iterator迭代器接口

### List接口

&emsp;&emsp;List代表一个有序的集合，集合中每个元素都有其对应的顺序索引。List允许使用重复元素，可以通过索引来访问指定位置的元素。

- ArrayList

&emsp;&emsp;ArrayList是一个`动态数组`，也是我们最常用的集合。它允许任何符合规则的元素插入甚至包括null。每一个ArrayList都有一个初始容量(10)，该容量代表了数组的大小。随着容器中的元素不断增加，容器的大小也会随着增加。在每次向容器中增加元素的同时都会进行容器检查，当快溢出时，就会进行扩容操作((oldCapacity * 3)/2 + 1)。所以如果我们明确所插入元素的多少，最好指定一个初始容量值，避免过多的进行扩容操作而浪费时间、空间效率。

&emsp;&emsp;关于ArrayList的详细介绍，请参见[[Java | ArrayList]](http://facedamon.github.io/post/java/collection/ArrayList)

- LinkedList

&emsp;&emsp;LinkedList是一个`双向链表`。它除了有ArrayList的基本操作方法之外还额外提供了get, remove, insert方法在LinkedList的首部和尾部。由于实现方式不同，`LinkedList不能随机访问`，它所有的操作都是要按照双向链表的需要执行。在列表中索引将从开头或结尾遍历列表。这样做的好处就是可以通过较低的代价在List中进行插入和删除。

&emsp;&emsp;关于LinkedList的详细介绍，请参见[[Java | LinkedList]](http://facedamon.github.io/post/java/collection/LinkedList)

- Vector/Stack 已经被历史淘汰。堆栈和队列，如果不考虑线程安全问题，请使用LinkedList。

> LinkedList和ArrayList的对比

1. 顺序插入速度ArrayList会比较快，因为ArrayList是基于数组实现的，数组是事先new好的，只要往指定位置塞一个数据就好了；LinkedList则不同，每次顺序插入的时候LinkedList将new一个对象出来，如果对象较大，那么new的时间势必会长一点，再加上一些引用赋值操作，所以顺序插入LinkedList必然慢于ArrayList
2. 基于上一点，因为LinkedList里面不仅维护了待插入的元素，还维护了Node的前置和后继Node，如果一个LinkedList中的Node非常多，那么LinkedList将比ArrayList更耗费一些内存
3. 数据遍历的速度：`使用各自遍历效率最高的方式，ArrayList的遍历效率会比LinkedList的遍历效率高一些`。**ArrayList使用普通for循环遍历时间为6ms｜ArrayList使用foreach循环遍历时间为12ms｜LinkedList使用普通for循环遍历时间为38482ms｜LinkedList使用foreach循环遍历时间为11ms。**有了JDK API的解释，这个结果并不让人感到意外，最最想要提出的一点是：如果使用普通for循环遍历LinkedList，在大数据量的情况下，其遍历速度将慢得令人发指.

4. 有些说法认为LinkedList做插入和删除更快，这种说法其实是不准确的：
   1. LinkedList做插入、删除的时候，慢在寻址，快在只需要改变前后Node的引用
   2. ArrayList做插入、删除的时候，满在数组元素的批量copy，快在寻址

&emsp;&emsp;**所以，如果待插入、删除的元素是在数据结构的前半段尤其是非常靠前的位置的时候，LinkedList的效率将大大快过ArrayList，因为ArrayList将批量copy大量的元素；越往后，对于LinkedList来说，因为它是双向链表，所以在第2个元素后面插入一个数据和在倒数第2个元素后面插入一个元素在效率上基本没有差别，但是ArrayList由于要批量copy的元素越来越少，操作速度必然追上乃至超过LinkedList。**

&emsp;&emsp;从这个分析看出，如果你十分确定你插入、删除的元素是在前半段，那么就使用LinkedList；如果你十分确定你插入、删除的元素在比较靠后的位置，那么可以考虑使用ArrayList。如果你不能确定你要做的插入、删除是在哪儿呢？那还是建议你使用LinkedList吧，因为一来LinkedList整体插入、删除的执行效率比较稳定，没有ArrayList这种越往后越快的情况；二来插入元素的时候，弄得不好ArrayList就要进行一次扩容，记住，**ArrayList底层数组扩容是一个既消耗时间又消耗空间的操作.**


### Set接口

&emsp;&emsp;Set是一种`不包括重复元素`的集合。它维持自己的内部排序，所以随机访问没有任何意义。与List一样，它同样允许null的存在但是仅有一个。由于Set接口的特殊性，所有传入Set集合中的元素都必须不同，同时要注意任何可变对象。如果在对集合中元素进行操作时，导致e1.equals(e2)==true，则必定会产生某些问题。Set接口有三个具体实现类：HashSet、LinkedHashSet、TreeSet

&emsp;&emsp; **Set是一种不包含重复的元素的Collection，无序，即任意的两个元素e1和e2都有e1.equals(e2)=false，Set最多有一个null元素。需要注意的是:虽然Set中元素没有顺序，但是元素在set中的位置是由该元素的HashCode决定的，其具体位置其实是固定的。**

- HashSet

&emsp;&emsp;HashSet是一个`没有重复元素`的集合，它由HashMap实现，不保证元素的顺序(元素插入的顺序与输出顺序不一致)，HashSet允许使用null元素，HashSet按照hash算法来存储集合元素。HashSet通过一个HashMap存储元素，元素存放在HashMap的key中，而value统一使用一个object对象。

&emsp;&emsp;关于HashSet的详细介绍，请参见[[Java | HashSet]](http://facedamon.github.io/post/java/collection/HashSet)

- LinkedHashSet

&emsp;&emsp;LinkedHashSet，其底层是`基于LinkedHashMap`来实现的，`有序、非同步`。LinkedHashSet集合同样是根据元素的hashCode值来决定元素的存储位置，但是它同时`使用链表维护了元素的次序`。这样使得元素看起来像是以插入顺序保存的。

&emsp;&emsp;关于LinkedHashSet的详细介绍，请参见[[Java | HashSet]](http://facedamon.github.io/post/java/collection/LinkedHashSet)

- TreeSet

&emsp;&emsp;TreeSet是一个`有序集合`。其底层是`基于TreeMap`来实现的，`非同步`。TreeSet支持两种排序方式，`自然排序和Compare`，其中自然排序为默认排序方式。当我们构造TreeSet时，若使用不带参数的构造函数，则TreeSet使用自然比较器；若用户需要只用自定义的比较器，则需要使用带比较器的构造函数。

&emsp;&emsp;关于TreeSet的详细介绍，请参见[[Java | HashSet]](http://facedamon.github.io/post/java/collection/TreeSet)

### Map接口

&emsp;&emsp;Map与List、Set接口不同，它是由一系列键值对组成的集合，提供了key到value的映射。同时它也`没有继承Collection`。在Map中它保证了key与value之间的一一对应关系。所以它不能存在相同的key，当然value可以相同。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/249993-20170117144843208-531675510.png)

- HashMap

&emsp;&emsp;以[哈希表数据结构](http://facedamon.github.io/post/data-struct/hash)实现，查找对象时通过哈希函数计算其位置，它是为快速查询而设计的，其内部定义了一个hash表数组(Entry[] table),元素会通过哈希转换函数将元素的哈希地址转换成数组中存放的索引，如果有冲突，则使用散列表的形式将所有相同哈希地址的元素串起来，可通过查看HashMap.Enntry源码它是一个单链表结构。

&emsp;&emsp;关于HashMap的详细介绍，请参见[[Java | HashSet]](http://facedamon.github.io/post/java/collection/HashMap)

- LinkedHashMap

&emsp;&emsp;LinkedHashMap是HashMap的一个子类，它保留`插入的顺序`，如果需要输出的顺序和输入时的相同，那么就选用LinkedHashMap.

&emsp;&emsp;LinkedHashMap实现与HashMap的不同之处在于，后者维护着一个运行于所有条目的双向链表。此链表定义了迭代顺序，该迭代顺序可以是`插入顺序或者是访问顺序`。

&emsp;&emsp;根据链表中元素的顺序可以分为：按插入顺序和按访问顺序(调用get方法)的链表。默认是插入顺序排序，如果指定按访问顺序排序，那么调用get方法后，会将这次访问的元素移至链表尾部，不断访问可以形成按访问顺序跑许的链表。

&emsp;&emsp;由于LinkedHashMap需要维护元素的插入顺序，因此性能略低于HashMap的性能，但在迭代访问Map里的全部元素时将有很好的性能，因为它以链表来维护内部顺序。

&emsp;&emsp;关于LinkedHashMap的详细介绍，请参见[[Java | HashSet]](http://facedamon.github.io/post/java/collection/LinkedHashMap)

- TreeMap

&emsp;&emsp;**TreeMap是一个有序的k-v集合，非同步，基于红黑树实现。**每一个k-v节点作为红黑树的一个节点。TreeMap存储时会进行排序，会根据key来对k-v键值对进行排序，其中排序方式也是分为两种，一种是自然排序，一种是定制排序，具体取决于使用的构造方法。

&emsp;&emsp;自然排序：TreeMap中所有的`key必须实现Comparable接口`，并且所有的key都应该是同一个类的对象。

&emsp;&emsp;定制排序：定义TreeMap时，`创建一个ccomparator对象`，该对象对所有的treeMap中所有的key值进行排序，采用定制排序的时候不需要TreeMap中所有的key必须实现Comparable接口。

&emsp;&emsp;TreeMap判断两个元素相等的标准：两个key通过compareTo()方法返回0，则认为这两个key相等。

&emsp;&emsp;关于TreeMap的详细介绍，请参见[[Java | HashSet]](http://facedamon.github.io/post/java/collection/TreeMap)

## Iterator与ListIterator

### Iterator

    public interface Iterator<E> {}

&emsp;&emsp;Iterator是一个接口，它是集合的迭代器。集合可以通过Iterator去遍历集合中的元素。

- boolean hasNext()：判断集合里是否存在下一个元素。如果有，hasNext()方法返回 true。
- E next()：返回集合里下一个元素。
- void remove()：删除集合里上一次next方法返回的元素。

&emsp;&emsp;注意：

1. Iterator只能单向移动。
2. Iterator.remove()是唯一安全的方式来在迭代过程中修改集合；如果在迭代过程中以任何其它的方式修改了基本集合将会产生未知的行为。而且每调用一次next()方法，remove()方法只能被调用一次，如果违反这个规则将抛出一个异常

### ListIterator

    public interface ListIterator<E> extends Iterator<E> {
      boolean hasNext();
      E next();
      boolean hasPrevious();
      E previous();
      int nextIndex();
      int previous();
      innt previousIndex();
      void remove();
      void set(E e);
      void add(E e);
    }

&emsp;&emsp;由以上定义我们可以推出ListIterator可以:

1. 双向移动（向前/向后遍历)
2. 产生相对于迭代器在列表中指向的当前位置的前一个和后一个元素的索引
3. 可以使用set()方法替换它访问过的最后一个元素
4. 可以使用add()方法在next()方法返回的元素之前或previous()方法返回的元素之后插入一个元素