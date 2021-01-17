
# 摘要

&emsp;&emsp;之前的文章讲解了两种Map，分别是HashMap与LinkedHashMap，它们保证了以O(1)的时间复杂度进行增、删、改、查，从存储角度考虑，这两种数据结构是非常优秀的。另外，LinkedHashMap还额外地保证了Map的遍历顺序可以与put顺序一致，解决了HashMap本身无序的问题。

&emsp;&emsp;尽管如此，HashMap与LinkedHashMap还是有自己的局限性----**它们不具备统计性能，或者说它们的统计性能时间复杂度并不是很好才更准确，所有的统计必须遍历所有Entry，因此时间复杂度为O(N)。**比如Map的Key有1、2、3、4、5、6、7，我现在要统计：

1. 所有Key比3大的键值对有哪些
2. Key最小的和Key最大的是哪两个

&emsp;&emsp;就类似这些操作，HashMap和LinkedHashMap做得比较差，此时我们可以使用TreeMap。`TreeMap的Key按照自然顺序进行排序或者根据创建映射时提供的Comparator接口进行排序。`**TreeMap为增、删、改、查这些操作提供了log(N)的时间开销，**从存储角度而言，这比HashMap与LinkedHashMap的O(1)时间复杂度要差些；但是在统计性能上，TreeMap同样可以保证log(N)的时间开销，这又比HashMap与LinkedHashMap的O(N)时间复杂度好不少。

> 因此总结而言：如果只需要存储功能，使用HashMap与LinkedHashMap是一种更好的选择；如果还需要保证统计性能或者需要对Key按照一定规则进行排序，那么使用TreeMap是一种更好的选择。

&emsp;&emsp;对于集合，我认为关注的点主要有四点：

## 红黑树的一些基本概念

|关注点|结论|
|---|---|
|TreeMap是否允许null|Key不允许为空，和Value允许为空|
|TreeMap是否允许重复数据|Key重复会覆盖、Value允许重复|
|TreeMap是否有序|按照key的自然排序或者Comparator接口指定的排序算法进行排序|
|TreeMap是否线程安全|非线程安全|

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/collection/)