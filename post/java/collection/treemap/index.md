
# 摘要

&emsp;&emsp;之前的文章讲解了两种Map，分别是HashMap与LinkedHashMap，它们保证了以O(1)的时间复杂度进行增、删、改、查，从存储角度考虑，这两种数据结构是非常优秀的。另外，LinkedHashMap还额外地保证了Map的遍历顺序可以与put顺序一致，解决了HashMap本身无序的问题。

&emsp;&emsp;尽管如此，HashMap与LinkedHashMap还是有自己的局限性----**它们不具备统计性能，或者说它们的统计性能时间复杂度并不是很好才更准确，所有的统计必须遍历所有Entry，因此时间复杂度为O(N)。**比如Map的Key有1、2、3、4、5、6、7，我现在要统计：

1. 所有Key比3大的键值对有哪些
2. Key最小的和Key最大的是哪两个

&emsp;&emsp;就类似这些操作，HashMap和LinkedHashMap做得比较差，此时我们可以使用TreeMap。`TreeMap的Key按照自然顺序进行排序或者根据创建映射时提供的Comparator接口进行排序。`**TreeMap为增、删、改、查这些操作提供了log(N)的时间开销，**从存储角度而言，这比HashMap与LinkedHashMap的O(1)时间复杂度要差些；但是在统计性能上，TreeMap同样可以保证log(N)的时间开销，这又比HashMap与LinkedHashMap的O(N)时间复杂度好不少。

> 因此总结而言：如果只需要存储功能，使用HashMap与LinkedHashMap是一种更好的选择；如果还需要保证统计性能或者需要对Key按照一定规则进行排序，那么使用TreeMap是一种更好的选择。

## 红黑树的一些基本概念

&emsp;&emsp;在讲TreeMap前还是先说一下红黑树的一些基本概念，这样可以更好地理解之后TreeMap的源代码。

&emsp;&emsp;二叉查找树在生成的时候是非常容易失衡的，造成的最坏情况就是一边倒（即只有左子树/右子树），这样会导致树检索的效率大大降低。

&emsp;&emsp;红黑树是为了维护二叉查找树的平衡而产生的一种树，根据维基百科的定义，红黑树有五个特性，但我觉得讲得不太易懂，我自己总结一下，红黑树的特性大致有三个（换句话说，插入、删除节点后整个红黑树也必须满足下面的三个性质，如果不满足则必须进行旋转）：

1. 根节点与叶节点都是黑色节点，其中叶节点为NULL节点
2. 每个红色节点的两个子节点都是黑色节点，换句话说就是不能有连续两个红色节点
3. 从根节点到所有叶子节点上的黑色节点数量是相同的

&emsp;&emsp;上述的性质约束了红黑树的关键：**从根到叶子的最长可能路径不多于最短可能路径的两倍长**。得到这个理论的理由是：

1. 红黑树中最短的可能路径全部为黑色节点的路径
2. 红黑树中最长的可能路径是红黑相间的路径

&emsp;&emsp;此时2正好是1的两倍长。结果就是这个树大致上是平衡的，因为插入、删除和查找某个值这样的操作最坏情况都要求与树的高度成正比，这个高度的理论上限允许红黑树在最坏情况下都是高效的，而不同于普通的二叉查找树，最终保证了红黑树能够以O(log2N)的时间复杂度进行搜索、插入、删除。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/801753-20170517224514744-1577819946.jpg)

&emsp;&emsp;可以看到根节点到所有NULL叶子节点所经过的黑色节点都是2个。

&emsp;&emsp;另外从这张图中我们还能得到一个理论：**红黑树并不是高度的平衡树**。所谓平衡二叉树指的是一个空树或它的左右两个子树的高度差绝对值不超过1，但是我们看：

- 最左边的路径0026-->0017-->0012-->0010-->0003-->NULL LEAF，它的高度为5
- 最后边的路径0026-->0041-->0047-->NULL LEAF，它的高度为3

&emsp;&emsp;左右子树的高度差值为2，因此红黑树并不是高度平衡的，它放弃了高度平衡的特性而只追求部分平衡，这种特性降低了插入、删除时对树旋转的要求，从而提升了树的整体性能。而其他平衡树比如AVL树虽然查找性能为性能是O(logn)，但是为了维护其平衡特性，可能要在插入、删除操作时进行多次的旋转，产生比较大的消耗。

&emsp;&emsp;对于集合，我认为关注的点主要有四点：

|关注点|结论|
|---|---|
|TreeMap是否允许null|Key不允许为空，和Value允许为空|
|TreeMap是否允许重复数据|Key重复会覆盖、Value允许重复|
|TreeMap是否有序|按照key的自然排序或者Comparator接口指定的排序算法进行排序|
|TreeMap是否线程安全|非线程安全|

## TreeMap基本数据结构

    public class TreeMap<K,V>
        extends AbstractMap<K,V>
            implements NavigableMap<K,V>, Cloneable, java.io.Serializable {}

&emsp;&emsp;TreeMap基于红黑树实现，既然是红黑树，那么每个节点中除了Key-->Value映射之外，必然存储了红黑树节点特有的一些内容，它们是：

1. 父节点引用
2. 左子节点引用
3. 右子节点引用
4. 节点颜色

    static final class Entry<K,V> implements Map.Entry<K,V> {
        K key;
        V value;
        Entry<K,V> left = null;
        Entry<K,V> right = null;
        Entry<K,V> parent;
        boolean color = BLACK;
        ...
    }

&emsp;&emsp;由于颜色只有红色和黑色两种，因此颜色可以使用布尔类型（boolean）来表示，黑色表示为true，红色为false。

### 添加

&emsp;&emsp;首先看一下测试代码：

    public class MapTest {
        @Test
        public void testTreeMap() {
            TreeMap<Integer,String> treeMap = new TreeMap<Integer,String>();
            treeMap.put(10,"10");
            treeMap.put(85,"85");
            treeMap.put(15,"15");
            treeMap.put(70,"70");
            treeMap.put(20，"20");
            treeMap.put(60,"60");
            treeMap.put(30,"30");
            treeMap.put(50,"50");

            for (Map.Entry<Integer,String> entry : treeMap.entrySet()) {
                sout(entry.getKey() + ":" + entry.getValue());
            }
        }
    }

&emsp;&emsp;本文接下来的内容会给出插入每条数据之后红黑树的数据结构是什么样子的。首先看一下treeMap的put方法的代码实现：

    public V put(K key, V value) {
        Entry<K,V> t = root;
        if (t == null) {
            // type check
            compare(key, key);
            root = new Entry(key, value, null);
            size = 1;
            modCount++;
            return null;
        }
        int cmp;
        Entry<K,V> parent;
        //split comparator and comparable paths
        Comparator<? super K> cpr = comparator;
        if (cpr != null) {
            do {
                parent = t;
                cmp = cpr.compare(key, t.key);
                if (cmp < 0) 
                    t = t.left;
                else if (cmp > 0)
                    t = t.right;
                else
                    return t.setValue(value);
            } while (t != null);
        }else {
         if (key == null)
             throw new NullPointerException();
         Comparable<? super K> k = (Comparable<? super K>) key;
         do {
             parent = t;
             cmp = k.compareTo(t.key);
             if (cmp < 0)
                 t = t.left;
             else if (cmp > 0)
                 t = t.right;
             else
                 return t.setValue(value);
            } while (t != null);
        }
        Entry<K,V> e = new Entry<>(key, value, parent);
        if (cmp < 0)
            parent.left = e;
        else
            parent.right = e;
        fixAfterInsertion(e);
        size++;
        modCount++;
        return null;
    }

&emsp;&emsp;从这段代码，先总结一下TreeMap添加数据的几个步骤：

1. 获取根节点，根节点为空，产生一个根节点，将其着色为黑色，退出余下流程
2. 获取比较器，如果传入的Comparator接口不为空，使用传入的Comparator接口实现类进行比较；如果传入的Comparator接口为空，将Key强转为Comparable接口进行比较
3. 从根节点开始逐一依照规定的排序算法进行比较，取比较值cmp，如果cmp=0，表示插入的Key已存在；如果cmp>0，取当前节点的右子节点；如果cmp<0，取当前节点的左子节点
4. 排除插入的Key已存在的情况，第（3）步的比较一直比较到当前节点t的左子节点或右子节点为null，此时t就是我们寻找到的节点，cmp>0则准备往t的右子节点插入新节点，cmp<0则准备往t的左子节点插入新节点
5. new出一个新节点，默认为黑色，根据cmp的值向t的左边或者右边进行插入
6. 插入之后进行修复，包括左旋、右旋、重新着色这些操作，让树保持平衡性

&emsp;&emsp;第1~第5步都没有什么问题，红黑树最核心的应当是第6步插入数据之后进行的修复工作，对应的Java代码是TreeMap中的fixAfterInsertion方法，下面看一下put每个数据之后TreeMap都做了什么操作，借此来理清TreeMap的实现原理。

