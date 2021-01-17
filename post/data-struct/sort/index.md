
# 摘要

&emsp;&emsp;排序是计算机内经常进行的一种操作，其目的是将一组“无序”的记录序列调整为“有序”记录。分内部排序和外部排序。

&emsp;&emsp;**若整个排序过程不需要访问外存便能完成，则称此类排序为内部排序。反之，若参加排序的记录数列很大，整个序列的排序过程不可能在内存中完成，则称此类排序为外部排序。**

&emsp;&emsp;将杂乱无章的数据元素，通过一定的方法按关键字顺序排列的过程叫做排序。**假定在待排序的记录序列中，存在多个具有相同的关键字的记录，若经过排序，这些记录的相对次序保持不变，即在原序列中，$R_i=R_j$，且$R_i$在$R_j$之前，而在排序后的序列中，$R_i$仍在$R_j$之前，则称这种排序算法是稳定的；否则称为不稳定。**

|排序方法|平均情况|最好情况|最坏情况|辅助空间|稳定性|类型|
|---|---|---|---|---|---|---|
|冒泡排序|O(n^2)|O(n)|O(n^2)|O(1)|稳定|交换排序|
|快速排序|O(nlogn)|O(nlogn)|O(n^2)|O(logn)~O(n)|不稳定|交换排序|
|直接插入排序|O(n^2)|O(n)|O(n^2)|O(1)|稳定|插入排序|
|希尔排序|O(nlogn)~O(n^2)|O(n^1.3)|O(n^2)|O(1)|不稳定|插入排序|
|简单选择排序|O(n^2)|O(n^2)|O(n^2)|O(1)|不稳定|选择排序|
|堆排序|o(nlogn)|O(nlogn)|O(nlogn)|O(1)|不稳定|选择排序|
|归并排序|O(nlogn)|O(nlogn)|O(nlogn)|O(n)|稳定|归并排序|
|计数排序|O(n+k)|O(n+k)|O(n+k)|O(k)|稳定|计数排序|

## 工具

    /**
     * 交换数组a中，位置i和位置j的数据
     * @param a
     * @param i
     * @param j
     */
    private void swap(int[] a, int i, int j) {
        if (a == null || a.length < 1) {
            return;
        }
        int temp = a[i];
        a[i] = a[j];
        a[j] = temp;
    }

</br>

    /**
     * 打印数组
     * @param a
     */
    private void printArr(int[] a) {
        if (a == null || a.length < 1) {
            return;
        }
        for (int k : a) {
            System.out.printf("%d\t", k);
        }
        System.out.println();
    }

</br>

    private long start = 0;

    @Before
    public void before() {
        this.start = System.currentTimeMillis();
    }

    @After
    public void after() {
        System.out.printf("运行时间:%d毫秒\n", (System.currentTimeMillis() - start));
    }

## 交换排序

### 冒泡排序

&emsp;&emsp;已知一组无序数据a[1]、a[2]、...a[n]，需将其按升序排序。首先比较a[1]与a[2]的值，若a[1]大于a[2]则交换两者的值，否则不变。再比较a[2]与a[3]的值，若a[2]大于a[3]则交换两者的值，否则不变。再比较a[3]与a[4]，以此类推，最后比较a[n-1]与a[n]的值。这样处理一轮后，a[n]的值一定是这组数据中最大的。再对a[1]~a[n-1]以相同方法处理一轮，则a[n-1]的值就一定是a[1]~a[n-1]中最大的。再对a[1]~a[n-2]以相同方法处理一轮，以此类推。共处理n-1轮后数组a就是以升序排列了。总的来讲，每一轮排序后最大或者小的数将移动到序列的最后，**理论上总共要进行n(n-1)/2次交换**。

    public void bubbleSort(int[] a) {
        if (a == null || a.length < 1) {
            return;
        }
        for (int i = 0; i < a.length - 1; i++) {
            for (int j = 0; j < a.length - 1 - i; j++) {
                if (a[j] > a[j + 1]) {
                    swap(a, j, j + 1);
                }
            }
        }
    }

    @Test
    public void testBubbleSort() {
        int[] a = {1, 2, 3, 4, 1, -1};
        bubbleSort(a);
        printArr(a);
        //-1	1	1	2	3	4	
        //运行时间:1毫秒
    }

- 优点：稳定
- 缺点：慢，每次只能移动相邻的两个元素。

### 快速排序

&emsp;&emsp;设要排序的数组是A[0]……A[N-1]，首先任意选取一个数据（通常选用数组的第一个数）作为关键数据，然后将所有比它小的数都放到它前面，所有比它大的数都放到它后面，这个过程称为一趟快速排序。值得注意的是，快速排序不是一种稳定的排序算法，也就是说，多个相同的值的相对位置也许会在算法结束时产生变动。

&emsp;&emsp;一趟快速排序的算法是：

1. 设置两个变量i、j，排序开始的时候：i=0，j=N-1；
2. 以第一个数组元素作为关键数据，赋值给key，即key=A[0]；
3. 从j开始向前搜索，即由后开始向前搜索(j--)，找到第一个小于key的值A[j]，将A[j]和A[i]互换；
4. 从i开始向后搜索，即由前开始向后搜索(i++)，找到第一个大于key的A[i]，将A[i]和A[j]互换；
5. 重复第3、4步，直到i=j； (3,4步中，没找到符合条件的值，即3中A[j]不小于key,4中A[i]不大于key的时候改变j、i的值，使得j=j-1，i=i+1，直至找到为止。找到符合条件的值，进行交换的时候i， j指针位置不变。另外，i==j这一过程一定正好是i+或j-完成的时候，此时令循环结束）。

&emsp;&emsp;下面以数列a={30,40,60,10,20,50}为例，演示它的快速排序过程(如下图)。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/121659127078460.jpg)

1. 在第1趟中，设置x=a[i]，即x=30。
2. 从"右 --> 左"查找小于x的数：找到满足条件的数a[j]=20，此时j=4；然后将a[j]赋值a[i]，此时i=0；接着从左往右遍历。
3. 从"左 --> 右"查找大于x的数：找到满足条件的数a[i]=40，此时i=1；然后将a[i]赋值a[j]，此时j=4；接着从右往左遍历。
4. 从"右 --> 左"查找小于x的数：找到满足条件的数a[j]=10，此时j=3；然后将a[j]赋值a[i]，此时i=1；接着从左往右遍历。
5. 从"左 --> 右"查找大于x的数：找到满足条件的数a[i]=60，此时i=2；然后将a[i]赋值a[j]，此时j=3；接着从右往左遍历。
6. 从"右 --> 左"查找小于x的数：没有找到满足条件的数。当i>=j时，停止查找；然后将x赋值给a[i]。此趟遍历结束！
7. **按照同样的方法，对子数列进行递归遍历。最后得到有序数组！**

> 1~6 步

    private int partition(int[] a, int left, int right) {
        if (a == null || a.length < 1 || left < 0 || left > right) {
            return -1;
        }
        int x = a[left];
        int i = left;
        int j = right;
        while (i < j) {
            while (i < j && a[j] > x) {
                j--;
            }
            if (i < j) {
                a[i++] = a[j];
            }
            while (i < j && a[i] < x) {
                i++;
            }
            if (i < j) {
                a[j--] = a[i];
            }
        }
        a[i] = x;
        return i;
    }

> 递归版本

    public void quickSort(int[] a, int left, int right) {
        if (a == null || a.length < 1 || a.length < right || left < 0 || left > right) {
            throw new RuntimeException("quickSort param err.");
        }
        int k = partition(a, left, right);
        if (k > left) {
            quickSort(a, left, k - 1);
        }
        if (k < right) {
            quickSort(a, k + 1, right);
        }
    }

> 非递归版本

    public void quickSort2(int[] a, int left, int right) {
        if (a == null || a.length < 1 || a.length < right || left < 0 || left > right) {
            throw new RuntimeException("quickSort param err.");
        }
        Deque<Integer> stack = new ArrayDeque<Integer>();
        // 先押入右下标，再押入左下标
        stack.push(right);
        stack.push(left);
        // i表示left，j表示right
        int i, j;
        while (!stack.isEmpty()) {
            i = stack.pop();
            j = stack.pop();
            if (i < j) {
                int k = partition(a, i, j);
                if (k > i) {
                    // 如果k > i，说明k的左侧都是比自身小的，则排序左侧
                    // 同样先押入右下标，再押入左下标
                    stack.push(k - 1);
                    stack.push(i);
                }
                if (k < j) {
                    stack.push(j);
                    stack.push(k + 1);
                }
            }
        }
    }

## 插入排序

### 直接插入排序

&emsp;&emsp;直接插入排序(Straight Insertion Sort)的基本思想是：把n个待排序的元素看成为一个有序表和一个无序表。开始时有序表中只包含1个元素，无序表中包含有n-1个元素，排序过程中每次从无序表中取出第一个元素，将它插入到有序表中的适当位置，使之成为新的有序表，重复n-1次可完成排序过程。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/121755260214693.jpg)

    public void insertSort(int[] a) {
        if (a == null || a.length < 1) {
            return;
        }
        int i, j, k;
        int n = a.length;
        for (i = 1; i < n; i++) {
            // 为a[i]在前面的a[0...n-1]有序区间中找一个合适的位置j
            for (j = i - 1; j >= 0; j--) {
                if (a[j] < a[i]) {
                    break;
                }
            }
            //找到了一个合适的位置j
            if (j != i - 1) {
                //将比a[i]大的数据向后移动一位，腾出j的位置
                int temp = a[i];
                for (k = i - 1; k > j; k--) {
                    a[k + 1] = a[k];
                }
                //当循环结束时k=j
                //将a[i]放入k+1的位置
                a[k + 1] = temp;
            }
        }
    }

### 希尔排序

&emsp;&emsp;希尔排序(Shell Sort)是插入排序的一种，它是针对直接插入排序算法的改进。该方法又称缩小增量排序，因DL．Shell于1959年提出而得名。

&emsp;&emsp;希尔排序实质上是一种分组插入方法。它的基本思想是：对于n个待排序的数列，取一个小于n的整数gap(gap被称为步长)将待排序元素分成若干个组子序列，**所有距离为gap的倍数的记录放在同一个组中**；然后，对各组内的元素进行直接插入排序。 这一趟排序完成之后，每一个组的元素都是有序的。然后减小gap的值，并重复执行上述的分组和排序。重复这样的操作，当gap=1时，整个数列就是有序的。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/122325434993142.jpg)

- 时间复杂度：希尔排序的时间复杂度与步长gap选取有关，当gap=1，希尔排序退化为直接插入排序，此时时间复杂度为O(n^2).而Hibbard增量序列的取法为$𝐷𝑘=2^𝑘−1$：{1, 3, 7, 15, 31, 63, 127, 255, 511, 1023, 2047, 4095, 8191...}
最坏时间复杂度为𝑂(𝑁^3/2)；平均时间复杂度约为𝑂(𝑁^5/4).

</br>

    /**
     * 𝐷𝑘=2^𝑘−1
     * ：{1, 3, 7, 15, 31, 63, 127, 255, 511, 1023, 2047, 4095, 8191...}
     * Hibbard增量
     *
     * @param len
     * @return
     */
    private List<Integer> getHibbard(int len) {
        if (len < 1) {
            return null;
        }
        List<Integer> arr = new ArrayList<>();
        int k = 1;
        while (true) {
            int temp = (1 << k++) - 1;
            if (temp <= len) {
                arr.add(temp);
            } else {
                break;
            }
        }
        return arr;
    }

    /**
     * 对希尔排序中的单组进行直接插入排序
     *
     * @param a    待排序的数组
     * @param left 数组的起始位置
     * @param gap  步长
     */
    private void gapSort(int[] a, int left, int gap) {
        if (a == null || a.length < 1 || a.length < left) {
            return;
        }
        int len = a.length;
        for (int j = left + gap; j < len; j += gap) {
            if (a[j] < a[j - gap]) {
                int temp = a[j];
                int k = j - gap;
                while (k >= 0 && a[k] > temp) {
                    a[k + gap] = a[k];
                    k -= gap;
                }
                a[k + gap] = temp;
            }
        }
    }

    public void shellSort(int[] a) {
        if (a == null || a.length < 1) {
            return;
        }
        int len = a.length;
        List<Integer> gaps = getHibbard(len);
        ListIterator<Integer> lt = gaps.listIterator(gaps.size());
        while (lt.hasPrevious()) {
            Integer gap = lt.previous();
            // 共gap个组，对每一组都执行直接插入排序
            for (int i = 0; i < gap; i++) {
                gapSort(a, i, gap);
            }
        }
    }

## 选择排序

### 简单选择排序

&emsp;&emsp;选择排序(Selection sort)是一种简单直观的排序算法。它的基本思想是：首先在未排序的数列中找到最小(or最大)元素，然后将其存放到数列的起始位置；接着，再从剩余未排序的元素中继续寻找最小(or最大)元素，然后放到已排序序列的末尾。以此类推，直到所有元素均排序完毕。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/130003333397123.jpg)

1. i=0。找出a[1...5]中的最小值a[3]=10，然后将a[0]和a[3]互换。 数列变化：20,40,30,10,60,50 -- > 10,40,30,20,60,50
2. i=1。找出a[2...5]中的最小值a[3]=20，然后将a[1]和a[3]互换。 数列变化：10,40,30,20,60,50 -- > 10,20,30,40,60,50
3. i=2。找出a[3...5]中的最小值，由于该最小值大于a[2]，该趟不做任何处理。 
4. i=3。找出a[4...5]中的最小值，由于该最小值大于a[3]，该趟不做任何处理。 
5. i=4。交换a[4]和a[5]的数据。 数列变化：10,20,30,40,60,50 -- > 10,20,30,40,50,60

</br>

    public void selectSort(int[] a) {
        if (a == null || a.length < 1) {
            return;
        }
        // 有序区间的末尾
        int i;
        // 无序区间的起始
        int j;
        // 无序区间中最小元素
        int min;
        int len = a.length;
        for (i = 0; i < len; i++) {
            min = i;
            //找出a[i+1]~a[n]之间的最小元素
            for (j = i + 1; j < len; j++) {
                if (a[j] < a[min]) {
                    min = j;
                }
            }
            // 若min != i, 说明无序区间存在小元素
            // 交换之后，保证了a[0]~a[i]之间的元素是有序的
            if (min != i) {
                swap(a, i, min);
            }
        }
    }

### 堆排序

> 未完待续

## 归并排序

&emsp;&emsp;将两个的有序数列合并成一个有序数列，我们称之为"归并"。

&emsp;&emsp;从下往上的归并排序：将待排序的数列分成若干个长度为1的子数列，然后将这些数列两两合并；得到若干个长度为2的有序数列，再将这些数列两两合并；得到若干个长度为4的有序数列，再将它们两两合并；直接合并成一个数列为止。这样就得到了我们想要的排序结果。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/151857010277686.jpg)

1. 将数组{80,30,60,40,20,10,50,70}看作由8个有序的子数组{80},{30},{60},{40},{20},{10},{50}和{70}组成。
2. 将这8个有序的子数列两两合并。得到4个有序的子树列{30,80},{40,60},{10,20}和{50,70}。
3. 将这4个有序的子数列两两合并。得到2个有序的子树列{30,40,60,80}和{10,20,50,70}。
4. 将这2个有序的子数列两两合并。得到1个有序的子树列{10,20,30,40,50,60,70,80}。

</br>

    /**
     * 将一个数组中的两个相邻有序区间合并成一个
     * @param a 包含两个有序区间
     * @param left 第一个有序区间的起始下标
     * @param mid 第一个有序区间的结束下标，也是第二个有序区间的起始下标
     * @param right 第二个有序区间的结束下标
     */
    private void merge(int[] a, int left, int mid, int right) {
        if (a == null || a.length < 1 || left < 0 || right >= a.length) {
            return;
        }
        //汇总2个有序区间的临时区域
        int[] temp = new int[right - left + 1];
        //第1个有序区间索引
        int i = left;
        //第2个有序区间索引
        int j = mid + 1;
        //临时区间索引
        int k = 0;
        while (i <= mid && j <= right) {
            if (a[i] <= a[j]) {
                temp[k++] = a[i++];
            } else {
                temp[k++] = a[j++];
            }
        }
        while (i <= mid) {
            temp[k++] = a[i++];
        }
        while (j <= right) {
            temp[k++] = a[j++];
        }
        //将排序后的元素全部整合到数组a中
        for (i = 0; i < k; i++) {
            a[left + i] = temp[i];
        }
        temp = null;
    }

    /**
     * 对数组a做若干次合并：
     * 数组a的总长度为len，将它分为若干个长度为gap的子数组
     * 将每2个相邻的子数组进行合并排序
     * @param a
     * @param gap
     */
    private void mergeGroups(int[] a, int gap) {
        if (a == null || a.length < 1 || gap > a.length) {
            return;
        }
        int i;
        int k = gap << 1;
        int len = a.length;
        //将每2个相邻的子数组进行合并排序
        for (i = 0; i + k < len; i += k) {
            merge(a, i, i + gap - 1, i + k - 1);
        }
        // 剩余一个子数组没有被匹配
        if (i + gap - 1 < len - 1) {
            merge(a, i, i + gap - 1, len - 1);
        }
    }

    public void mergeSort(int[] a) {
        if (a == null || a.length < 1) {
            return;
        }
        int len = a.length;
        for (int i = 1; i < len; i <<= 1) {
            mergeGroups(a, i);
        }
    }

## 计数排序

&emsp;&emsp;计数排序的核心在于将输入的数据值转化为键存储在额外开辟的数组空间中。作为一种线性时间复杂度的排序，计数排序要求输入的数据必须是有确定范围的`正整数`。

&emsp;&emsp;假设待排序的数组a中共有N个整数，并且已知数组a中数据的范围[0, MAX)。在排序时，创建容量为MAX的数组r，并将数组元素都初始化为0；在排序时，逐个遍历数组a，将数组a的值，作为"辅助数组r"的下标。当a中数据被读取时，就将辅助数组的值加1。例如，读取到数组a[3]=5，则将r[5]的值+1。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/152240225909832.jpg)

    private int getMaxValue(int[] a) {
        if (a == null || a.length < 1) {
            return -1;
        }
        int max = a[0];
        for (int v : a) {
            if (v < 0) {
                throw new IllegalArgumentException("number must grant than zero");
            }
            if (max < v) {
                max = v;
            }
        }
        return max;
    }

    /**
     * 计数排序
     * @param a
     */
    public void countSort(int[] a) {
        if (a == null || a.length < 1) {
            return;
        }
        int tempLen = getMaxValue(a) + 1;
        int[] temp = new int[tempLen];
        for (int v : a) {
            temp[v]++;
        }
        int i = 0;
        for (int j = 0; j < tempLen; j++) {
            while (temp[j] > 0) {
                a[i++] = j;
                temp[j]--;
            }
        }
    }

## 桶排序

&emsp;&emsp;桶排序是计数排序的升级版。它利用了函数的映射关系，高效与否的关键就在于这个映射函数的确定。为了使桶排序更加高效，我们需要做到这两点：

1. 在额外空间充足的情况下，尽量增大桶的数量
2. 使用的映射函数能够将输入的 N 个数据均匀的分配到 K 个桶中

&emsp;&emsp;当输入的数据可以均匀的分配到每一个桶中，最快。当输入的数据被分配到了同一个桶中，最慢。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/Bucket_sort_1.svg_.png)
![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/Bucket_sort_2.svg_.png)

    /**
     * 桶排序
     * @param a
     * @param bucketSize 桶数
     */
    public void bucketSort(int[] a, int bucketSize) {
        if (a == null || a.length < 1 || bucketSize < 1) {
            return;
        }
        int min;
        int max;
        min = max = a[0];
        for (int v : a) {
            if (v < min) {
                min = v;
            } else if (v > max) {
                max = v;
            }
        }
        // 每个桶装载的数据量
        int bucketCount = (int) (Math.floor(max - min) / bucketSize) + 1;
        // 二维数组，一位表示桶，二位表示桶中数据
        int[][] buckets = new int[bucketCount][0];
        // 利用映射函数将数据分配到各个桶中
        int len = a.length;
        for (int i = 0; i < len; i++) {
            int index = (int) (Math.floor(a[i] - min) / bucketSize);
            buckets[index] = append(buckets[index], a[i]);
        }
        int index = 0;
        for (int[] bucket: buckets) {
            if (bucket.length <= 0) {
                continue;
            }
            // 每个桶中数据进行插入排序
            insertSort(bucket);
            for (int v : bucket) {
                a[index++] = v;
            }
        }
    }

    /**
     * 扩容
     * @param bucket
     * @param i
     * @return
     */
    private int[] append(int[] bucket, int i) {
        bucket = Arrays.copyOf(bucket, bucket.length + 1);
        bucket[bucket.length - 1] = i;
        return bucket;
    }

## 基数排序