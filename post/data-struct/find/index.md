
> 转载自 https://www.cnblogs.com/maybe2030/p/4715035.html#_label0

# 摘要

&emsp;&emsp;查找是在大量的信息中寻找一个特定的信息元素，在计算机应用中，查找是常用的基本运算，例如编译程序中符号表的查找。本文简单概括性的介绍了常见的七种查找算法，说是七种，其实二分查找、插值查找以及斐波那契查找都可以归为一类——插值查找。插值查找和斐波那契查找是在二分查找的基础上的优化查找算法。树表查找和哈希查找会在后续的博文中进行详细介绍。

- 查找定义：根据给定的某个值，在查找表中确定一个其关键字等于给定值的数据元素
- 查找算法分类：

1. 静态查找和动态查找(静态或动态都是针对查找表而言的。动态表指查找表中有删除和插入操作的表)
2. 无序查找和有序查找
   1. 无序查找：被查找数列有序无序均可
   2. 有序查找：被查找数列必须为有序数列

- 平均查找长度(Average Search Length, ASL): 由于所费时间在关键字的比较上，所以把平均需要和待查找值比较的关键字次数称为平均查找长度。


$$
ASL=\sum_{i=1}^nP_iC_i
$$

&emsp;&emsp;其中n为查找表中元素个数，Pi为查找第i个元素的概率，通常假设每个元素查找概率相同，Pi=1/n，Ci是找到第i个元素的比较次数。

> 一个算法的ASL越大，说明时间性能差，反之，时间性能好.

## 顺序查找

- 说明：顺序查找适合于存储结构为顺序存储或链接存储的线性表。
- 基本思想：顺序查找也称为线形查找，属于无序查找算法。从数据结构线形表的一端开始，顺序扫描，依次将扫描到的结点关键字与给定值k相比较，若相等则表示查找成功；若扫描结束仍没有找到关键字等于k的结点，表示查找失败。
- 复杂度分析
  - 假设每个元素查找概率相同，$P_i=\frac1n$
  - Ci在于这个元素在查找表中的位置，如第0号元素就需要比较一次，第一号元素比较2次......第n号元素要比较n+1次。所以Ci=i；所以 
  
![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/1539112-20190108133107514-1052426163.png)

&emsp;&emsp;可以看出，顺序查找方法查找成功的平均 比较次数约为表长的一半

&emsp;&emsp;**当待查找元素不在查找表中时，也就是扫描整个表都没有找到，即比较了n次，查找失败O(n)**

    // 顺序查找
    int sequenceFind(int[] a, int value, int n) {
        if (a == null || a.length != n || n < 1) {
            return -1;
        }
        if (n == 1 && a[0] != value) {
            return -1;
        }
        for (int i = 0; i < n; i++) {
            if (a[i] == value) {
                return i;
            }
        }
        return -1;
    }

## 二分查找

- 说明：元素必须是有序的，如果是无序的则要先进行排序操作。
- 基本思想：也称为是折半查找，属于有序查找算法。用给定值k先与中间结点的关键字比较，中间结点把线形表分成两个子表，若相等则查找成功；若不相等，再根据k与该中间结点关键字的比较结果确定下一步查找哪个子表，这样递归进行，直到查找到或查找结束发现表中没有这样的结点。
- 复杂度分析：**最坏情况下，关键词比较次数为log2(n+1)，且期望时间复杂度为O(log2n)**

&emsp;&emsp;注：折半查找的前提条件是需要有序表顺序存储，对于静态查找表，一次排序后不再变化，折半查找能得到不错的效率。但对于需要频繁执行插入或删除操作的数据集来说，维护有序的排序会带来不小的工作量，那就不建议使用。

    // 二分查找
    int binaryFind(int[] a, int value, int n) {
        if (a == null || a.length != n || n < 1) {
            return -1;
        }
        if (n == 1 && a[0] != value) {
            return -1;
        }
        int low = 0;
        int high = n - 1;
        int mid;
        while (low <= high) {
            mid = (low + high) >> 1;
            if (a[mid] == value) {
                return mid;
            }
            if (a[mid] > value) {
                high = mid - 1;
            }
            if (a[mid] < value) {
                low = mid + 1;
            }
        }
        return -1;
    }

    @Test
    public void testBinaryFind() {
        int[] a = {1, 2, 3, 4, 5, 6};
        //int[] a = {};
        //int[] a = {1};
        //int[] a = {1,2};
        System.out.println(binaryFind(a, 3, a.length));
    }

    // 二分查找 递归
    int binaryFindRev(int[] a, int value, int low, int high) {
        if (a == null || a.length < 1) {
            return -1;
        }
        int mid = low + ((high - low) >> 1);
        if (a[mid] == value) {
            return mid;
        }
        if (mid >= high) {
            return -1;
        }
        if (a[mid] > value) {
            return binaryFindRev(a, value, low, mid - 1);
        }
        if (a[mid] < value) {
            return binaryFindRev(a, value, mid + 1, high);
        }
        return -1;
    }

    @Test
    public void testBinaryFindRev() {
        int[] a = {1, 2, 3, 4, 5, 6};
        //int[] a = {};
        //int[] a = {1};
        //int[] a = {1,2};
        System.out.println(binaryFindRev(a, 3, 0, a.length - 1));
    }

## 插值查找

&emsp;&emsp;在介绍插值查找之前，首先考虑一个新问题，为什么上述算法一定要是折半，而不是折四分之一或者折更多呢？

&emsp;&emsp;打个比方，在英文字典里面查“apple”，你下意识翻开字典是翻前面的书页还是后面的书页呢？如果再让你查“zoo”，你又怎么查？很显然，这里你绝对不会是从中间开始查起，而是有一定目的的往前或往后翻。

&emsp;&emsp;同样的，比如要在取值范围1 ~ 10000 之间 100 个元素从小到大均匀分布的数组中查找5， 我们自然会考虑从数组下标较小的开始查找。

&emsp;&emsp;经过以上分析，折半查找这种查找方式，**不是自适应的（也就是说是傻瓜式的）**。二分查找中查找点计算如下：

                                            // offset
    mid = (low + high) / 2; ==> mid = low + 1/2 * (high - low);

&emsp;&emsp;通过类比，我们可以将查找的点改进为如下：

                // offset
    mid = low + (key - a[low])/(a[high] - a[low]) * (high - low);

&emsp;&emsp;也就是将上述的比例参数1/2改进为自适应的，**根据关键字在整个有序表中所处的位置，让mid值的变化更靠近关键字key**，这样也就间接地减少了比较次数。

&emsp;&emsp;基本思想：基于二分查找算法，将查找点的选择改进为自适应选择，可以提高查找效率。当然，差值查找也属于有序查找。

&emsp;&emsp;**注：对于表长较大，而关键字分布又比较均匀的查找表来说，插值查找算法的平均性能比折半查找要好的多。反之，数组中如果分布非常不均匀，那么插值查找未必是很合适的选择。**

- 复杂度分析：**查找成功或者失败的时间复杂度均为O(log2(log2n))。**

</br>

    int insertionFind(int[] a, int value, int n) {
        if (a == null || a.length != n || n < 1) {
            return -1;
        }
        if (n == 1 && a[0] != value) {
            return -1;
        }
        int low = 0;
        int high = n - 1;
        int mid = 0;
        while (low <= high) {
            // 防止by zero 错误
            if (low == high) {
                if (a[low] == value) {
                    return low;
                }
                return -1;
            }
            mid = low + (value - a[low])/(a[high] - a[low]) * (high - low);
            if (a[mid] == value) {
                return mid;
            }
            if (a[mid] > value) {
                high = mid - 1;
            }
            if (a[mid] < value) {
                low = mid + 1;
            }
        }
        return -1;
    }

    @Test
    public void testInsertionFind() {
        int[] a = {1, 2, 3, 4, 5, 6};
        //int[] a = {};
        //int[] a = {1};
        //int[] a = {1,2,3};
        //int[] a = {1,2,4};
        System.out.println(insertionFind(a, 3, a.length));
    }

## 斐波那契查找

&emsp;&emsp;在介绍斐波那契查找算法之前，我们先介绍一下很它紧密相连并且大家都熟知的一个概念——黄金分割。

&emsp;&emsp;黄金比例又称黄金分割，是指事物各部分间一定的数学比例关系，即将整体一分为二，较大部分与较小部分之比等于整体与较大部分之比，其比值约为1:0.618或1.618:1。

&emsp;&emsp;0.618被公认为最具有审美意义的比例数字，这个数值的作用不仅仅体现在诸如绘画、雕塑、音乐、建筑等艺术领域，而且在管理、工程设计等方面也有着不可忽视的作用。因此被称为黄金分割。

&emsp;&emsp;大家记不记得斐波那契数列：1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89…….（从第三个数开始，后边每一个数都是前两个数的和）。然后我们会发现，随着斐波那契数列的递增，前后两个数的比值会越来越接近0.618，利用这个特性，我们就可以将黄金比例运用到查找技术中。

- 基本思想：也是二分查找的一种提升算法，通过运用黄金比例的概念在数列中选择查找点进行查找，提高查找效率。同样地，斐波那契查找也属于一种有序查找算法。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/20150323100632467.jpeg)

&emsp;&emsp;斐波那契查找与折半查找很相似，他是根据斐波那契序列的特点对有序表进行分割的。他**要求开始表中记录的个数为某个斐波那契数小1，及n=F(k)-1**;

&emsp;&emsp;开始将k值与第F(k-1)位置的记录比较,mid=low+F(k-1)+1,比较结果分为三种：

1. 相等，返回mid
2. 大于，low=mid+1，k-=2；(说明：low=mid+1说明待查找的元素在[mid+1,high]范围内，k-=2 说明范围[mid+1,high]内的元素个数为n-(F(k-1))= Fk-1-F(k-1)=Fk-F(k-1)-1=F(k-2)-1个，所以可以递归的应用斐波那契查找。)
3. 小于，high=mid-1,k-=1;(说明：low=mid+1说明待查找的元素在[low,mid-1]范围内，k-=1 说明范围[low,mid-1]内的元素个数为F(k-1)-1个，所以可以递归 的应用斐波那契查找。)

- 复杂度分析：**最坏情况下，时间复杂度为O(log2n)，且其期望复杂度也为O(log2n)**。

</br>

    // 斐波那契数组的长度
    private static final int MAX_SIZE = 20;

    // 构造一个斐波那契数组
    int[] fibonacci() {
        int[] f = new int[MAX_SIZE];
        f[0] = 0;
        f[1] = 1;
        for (int i = 2; i < MAX_SIZE; i++) {
            f[i] = f[i-1]+f[i-2];
        }
        return f;
    }

    // 斐波那契查找
    int fibonacciFind(int a[], int value) {
        if (a == null || a.length < 1) {
            return -1;
        }
        int low = 0;
        int high = a.length - 1;
        // fibonacci index
        int k = 0;
        int mid = 0;
        int[] f = fibonacci();
        while (f[k] - 1 < a.length) {
            k++;
        }
        // 扩展到f[k]长度
        int[] temp = Arrays.copyOf(a, f[k]);
        // 扩展的空间复制最后一位的数字
        for (int i = high+1; i < temp.length; i++) {
            temp[i] = a[high];
        }

        while (low <= high) {
            mid = low + f[k-1]-1;
            if (temp[mid] > value) {
                high = mid - 1;
                k--;
            } else if (temp[mid] < value) {
                low = mid + 1;
                k-=2;
            } else {
                //查找值的下标在arr数组额范围内，直接返回
                if (mid < a.length) {
                    return mid;
                } else {
                    //不在就返回right,为什么？因为后面几位的值和right的值是一样的，说明查找的值就是right
                    return high;
                }
            }
        }
        return -1;
    }

    @Test
    public void testFibonacciFind() {
        int[] a = {0,16,24,35,47,59,62,88,99};
        //int[] a = {0,16,24,35,47,59,62,88};
        //int[] a = {0,16};
        //int[] a = {0};
        //int[] a = {99};
        //int[] a = {};
        System.out.println(fibonacciFind(a, 99));
    }

## 二叉查找树

- 基本思想：二叉查找树是先对待查找的数据进行生成树，确保树的左分支的值小于右分支的值，然后在就行和每个节点的父节点比较大小，查找最适合的范围。 这个算法的查找效率很高，但是如果使用这种查找方法要首先创建树。 
- 有关二叉查找树的详细说明，请移步[二叉查找树](http://facedamon.github.io/http://127.0.0.1:1313/post/data-struct/tree-struct/bstree/)

- 复杂度分析：**它和二分查找一样，插入和查找的时间复杂度均为O(logn)，但是在最坏的情况下仍然会有O(n)的时间复杂度**。原因在于插入和删除元素的时候，树没有保持平衡。我们追求的是在最坏的情况下仍然有较好的时间复杂度，这就是平衡查找树设计的初衷。
- 基于二叉查找树进行优化，进而可以得到其他的树表查找算法，如平衡树、红黑树等高效算法。

## 分块查找

&emsp;&emsp;分块查找又称索引顺序查找，它是顺序查找的一种改进方法。

- 将n个数据元素"按块有序"划分为m块（m ≤ n）。每一块中的结点不必有序，但块与块之间必须"按块有序"；即第1块中任一元素的关键字都必须小于第2块中任一元素的关键字；而第2块中任一元素又都必须小于第3块中的任一元素，……
- 算法流程
  - 先选取各块中的最大关键字构成一个索引表
  - 查找分两个部分：先对索引表进行二分查找或者顺序查找，以确定待查记录在哪一块中；然后，在已确定的块中用顺序法进行查找

</br>

    private static class BlockFind {
        // 索引表
        private int[] index;
        // 块
        private ArrayList<Integer>[] lists;

        // 构建索引表
        public BlockFind(int[] index) {
            if (index != null && index.length != 0) {
                this.index = index;
                this.lists = new ArrayList[index.length];
                for (int i = 0; i < lists.length; i++) {
                    lists[i] = new ArrayList<Integer>();
                }
            } else {
                throw new RuntimeException("index cannot be null or empty.");
            }
        }

        /**
         * 二分查找定位索引位置
         * @param data 要插入的值
         * @return
         */
        private int binaryFind(int data) {
            int low = 0;
            int high = index.length - 1;
            int mid = 0;
            while (low <= high) {
                mid = (low + high) >> 1;
                if (index[mid] > data) {
                    high = mid - 1;
                } else {
                    // 如果相等，也插入在后面
                    low = mid + 1;
                }
            }
            return low;
        }

        // 打印每块元素
        public void printAll() {
            for (int i = 0; i < lists.length; i++) {
                ArrayList list = lists[i];
                System.out.println("ArrayList " + i + ":");
                System.out.println(lists[i].toString());
            }
        }

        // 插入元素
        public void insert(int value) {
            int i = binaryFind(value);
            lists[i].add(value);
        }

        // 查找元素
        public String find(int data) {
            int i = binaryFind(data);
            for (int j = 0; j < lists[i].size(); j++) {
                if (data == lists[i].get(j)) {
                    return String.format("在第%d块，第%d个",i+1, j+1);
                }
            }
            return String.format("%d 不存在", data);
        }
    }

    @Test
    public void testBlockFind() {
        int[] index = {10, 20, 30};
        BlockFind blockFind = new BlockFind(index);
        blockFind.insert(1);
        blockFind.insert(12);
        blockFind.insert(22);

        blockFind.insert(9);
        blockFind.insert(18);
        blockFind.insert(23);

        blockFind.insert(5);
        blockFind.insert(15);
        blockFind.insert(27);

        blockFind.printAll();

        System.out.println(blockFind.find(18));
        System.out.println(blockFind.find(29));
    }

### 特点

&emsp;&emsp;分块查找的特点其实显而易见，那就是分块查找拥有顺序查找和二分查找的双重优势，即顺序查找不需要有序，二分查找的速度快。

&emsp;&emsp;分块查找由于只需要索引表有序，所以特别适合用于在动态变化的数据元素序列中查找。但是如何分块比较复杂。如果分块过于稀疏，则可能导致每一块的内容过多，在顺序查找时效率很低；如果分块过密，则又会导致块数很多，无论是插入还是删除数据，都会频繁地进行二分查找；如果块数特别多，则基本上和直接二分查找的动态插入数据类似，这样分块查找就没有意义了。

&emsp;&emsp;所以对于分块查找来说，可以根据数据量的大小及数据的区间来进行对分块的选择。二分查找的平均查找长度近似 log2(n+1)-1，这里的n是块数；顺序查找的平均查找长度为 (n+1)/2，这里的 n 是每块的个数。

&emsp;&emsp;尽量等分为固定的块，假设块数为 a，每个块内的元素数量为 b，则 b=n/a，那么接下来就好办了，如果给定一个数据量 n 进行分块，则总的平均查找长度为 (b+1)/2+log2(a+1)-1，这样就可以解出 a 和 b 分别为多少了。

## 哈希查找

- 什么是哈希表（Hash）？

&emsp;&emsp;我们使用一个下标范围比较大的数组来存储元素。可以设计一个函数（哈希函数， 也叫做散列函数），使得每个元素的关键字都与一个函数值（即数组下标）相对应，于是用这个数组单元来存储这个元素；也可以简单的理解为，按照关键字为每一个元素"分类"，然后将这个元素存储在相应"类"所对应的地方。但是，不能够保证每个元素的关键字与函数值是一一对应的，因此极有可能出现对于不同的元素，却计算出了相同的函数值，这样就产生了"冲突"，换句话说，就是把不同的元素分在了相同的"类"之中。后面我们将看到一种解决"冲突"的简便做法。

> 总的来说，"直接定址"与"解决冲突"是哈希表的两大特点。

- 什么是哈希函数

&emsp;&emsp;哈希函数的规则是：通过某种转换关系，使关键字适度的分散到指定大小的的顺序结构中，越分散，则以后查找的时间复杂度越小，空间复杂度越高。

- 算法思想：哈希的思路很简单，如果所有的键都是整数，那么就可以使用一个简单的无序数组来实现：将键作为索引，值即为其对应的值，这样就可以快速访问任意键的值。这是对于简单的键的情况，我们将其扩展到可以处理更加复杂的类型的键。
- 算法流程
  - 用给定的哈希函数构造哈希表；
  - 根据选择的冲突处理方法解决地址冲突；
  - 在哈希表的基础上执行哈希查找。
- 复杂度分析
  - 单纯论查找复杂度：对于无冲突的Hash表而言，查找复杂度为O(1)（注意，在查找之前我们需要构建相应的Hash表）。

&emsp;&emsp;有关哈希表的详细说明，请移步[散列表](http://facedamon.github.io/post/data-struct/hash/)