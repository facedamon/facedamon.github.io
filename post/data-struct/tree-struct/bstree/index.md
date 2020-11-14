
> 转载自 https://www.cnblogs.com/skywang12345/
> 转载《算法导论》
> 转载 https://github.com/emirpasic/gods 

# 概要
&emsp;&emsp;本节先对二叉树的相关理论知识进行介绍，然后给出Java的详细实现。关于二叉树的学习，需要说明的是：它并不难，不仅不难，而且它非常简单。初次接触树的时候，我也觉得它似乎很难；而之所以产生这种感觉是由于二叉树有一大堆陌生概念、性质等内容。而当我真正的实现了二叉树再回头来看它的相关概念和性质的时候，觉得原来它是如此简单！因此，建议在学习二叉树的时候：先对二叉树基本的概念、性质有个基本了解，遇到难懂的知识点，可以画图来帮助理解；在有个基本概念后，再亲自动手实现二叉查找树；最后再回头来总结以下二叉树的理论知识时，你会发现它的确很简单！在代码实践中，我以二叉查找树，而不是单纯的二叉树为例子进行说明，单纯的二叉树非常简单，实际使用很少。况且掌握了二叉查找树，二叉树也就自然掌握了。

# 树的介绍
## 树的定义
&emsp;&emsp;树是一种数据结构，它是由n(n >=1) 个有限结点组成一个具有层次关系的集合


&emsp;&emsp;**把它叫做树是因为它看起来像一颗倒挂的树，也就是说它是跟朝上，而叶朝下。它具有以下特点：**


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/树的示意图.jpg)

1. 每个结点有零个或多个子结点;
2. 没有父结点的结点成为根结点;
3. 每一个非根结点有且只有一个父结点;
4. 除了根结点外，每个子结点可以分为多个不相交的子树。

## 树的基本术语
&emsp;&emsp;若一个结点有子树，那么该结点称为子树根的"双亲"，子树的根是该结点的"孩子"。有相同双亲的结点互为"兄弟"。一个结点的所有子树上的任何结点都是该结点的后裔。从根结点到某个结点的路径上的所有结点都是该结点的祖先。

1. **结点的度：** 结点拥有的子树的数目

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/结点的度.jpg)

&emsp;&emsp;如图A结点的度为3，B结点的度为2，C结点的度为1，D结点的度为3.E,F,G,H,I以及J度为0

2. **叶子：**度为零的结点
3. **分支结点：**度不为零的结点
4. **树的度：**树中结点的最大的度
5. **层次：**根结点的层次为1，其余结点的层次等于该结点的双亲结点的层次加1
6. **树的高度：**树中结点的最大层次
7. **无序树：**如果树中结点的各子树之间的次序是不重要的，可以交换位置
8. **有序树：**如果树中结点的各个子树之间的次序是重要的，不可以交换位置。
9. **森林：**0个或多个不相交的树组成。对森林加上一个根，森林即成为树；删去根，树即成为森林

# 二叉树的介绍
## 二叉树的定义
&emsp;&emsp;二叉树是每个结点最多有两个子树的树结构。它有五种基本形态：二叉树可以是空集；根可以有空的左子树或右子树；或者左右树皆为空。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/二叉树的五种基本形态.jpg)

## 二叉树的性质
&emsp;&emsp;二叉树有以下几个性质：

1. 二叉树第i层上的结点数目最多为 $2^{i-1}(i\geq1)$

2. 深度为k的二叉树最多有 $2^k-1 (i\geq1)$个结点
3. 包含n个结点的二叉树的高度至少$\log_2(n+1)$
4. 在任意一颗二叉树中，若终端结点的个数为$n_0$，度为2的结点数为$n_2$，则$n_0=n_2+1$

### 性质1:二叉树第i层上的结点数目最多为$2^{i-1}(i\geq1)$
&emsp;&emsp;证明： 下面用“数学归纳法”进行证明。

1. 当i=1时，第i层的结点数目为$2^{i-1}=2^0=1$。因为第1层上只有一个跟结点，所以命题成立。
2. 假设当i>1，第i层的结点数目为$2^{i-1}$.这个是根据1推断出来的。下面根据这个假设，推断出"第(i+1)层的结点数目为$2^i$"即可。由于二叉树的每个结点至多有两个孩子，故"第(i+1)层上的结点数目"最多是"第i层的结点数目的2倍"。即，第(i+1)层上的结点数目最大值=$2\times2^{i-1}=2^i$.故假设成立，原命题成立。

### 性质2:深度为k的二叉树至多有$2^k-1(k\geq1)$个结点
&emsp;&emsp;证明：在具有相同深度的二叉树中，当每一层都含有最大结点数时，其树中结点树最多。利用性质1可知，深度为k的二叉树的结点数至多为$$2^0+2^1+...+2^{k-1}=2^K-1$$ 上述公式可以根据等比数列证明
$$
    f(n) = 
    \begin{cases}
    n_1, & \text{q=1}  \\\\[2ex]
    \frac {a_1\times(1-q^n)}{1-q}, & \text{$q\neq1$}
    \end{cases}
$$

### 性质3:包含n个结点的二叉树的高度至少为$log_2(n+1)$
&emsp;&emsp;证明：根据性质2可知，高度为h的二叉树最多有$2^h-1$个结点。反之，对于包含n个结点的二叉树的高度至少为$log_2(n+1)$

### 性质4:在任意一颗二叉树中，若终端结点的个数为$n_0$，度为2的结点数为$n_2$，则$n_0=n_2+1$
&emsp;&emsp;证明：因为二叉树中所以结点的度数均不大于2，所以结点总数(记为n)="0度结点数($n_0$)" + "1度结点数($n_1$)" + "2度结点数($n_2$)"。由此，得到等式一.$n=n_0+n_1+n_2$

另一方面，0度结点没有孩子，1度结点有一个孩子，2度结点有两个孩子，故二叉树中孩子结点总数是$n_1+2n_2$。此外，只有根不是任何结点的孩子。故二叉树中的结点总数又可表示为等式二：$n=n_1+2n_2+1$

由等式一和等式二计算得到：$n_0=n_2+1$。原命题得证

## 满二叉树，完全二叉树和二叉查找树
### 满二叉树
**定义**：高度为h，并且由$2^h-1$个结点的二叉树，被称为满二叉树

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/满二叉树示意图.jpg)

### 完全二叉树
**定义**：一颗二叉树中，只有最下面两层结点的度可以小于2，并且最下一层的叶结点集中在靠左的若干位置上。这样的二叉树被称为完全二叉树。
**特点**：叶子结点只能出现在最下层和次下层，且最下层的叶子结点集中在树的左部。显然，一颗满二叉树一定是完全二叉树，而完全二叉树不一定是满二叉树。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/完全二叉树示意图.jpg)

### 二叉查找树
**定义**：二叉查找树(Binary Search Tree), 又被称为二叉搜索树。设x为二叉查找树中的一个结点，x结点包含关键字key，结点x的key值记为key[x].如果y是x的左子树中的一个结点，则key[y]<=key[x];如果y是x的右子树的一个结点，则key[y]>=key[x]

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/二叉查找树示意图.jpg)

在二叉查找树中：
1. 若任意结点的左子树不空，则左子树上所有结点的值均小于它的根结点的值；
2. 任意结点的右子树不空，则右子树上所有结点的值均大于它的根结点的值；
3. 任意结点的左、右子树也分别为二叉查找树；
4. 没有键值相等的结点(no duplicate nodes)

> 在实际应用中，二叉查找树的使用比较多。下面，用Java语言实现二叉查找树

# 二叉查找树的Java实现
## 二叉查找树结点的定义

        public class BSTree<T extends Comparable<T>> {

            private BSTNode<T> mRoot;   //根结点

            public class BSTNode<T extends Comparable<T>> {
                T key;              // 键值
                BSTNode<T> left;    // 左结点
                BSTNode<T> right;   //右结点
                BSTNode<T> parent;  //父结点

                public BSTNode(T key, BSTNode<T> parent, BSTNode<T> left, BSTNode<T> right) {
                    this.key = key;
                    this.parent = parent;
                    this.left = left;
                    this.right = right;
                }
            }
        }

&emsp;&emsp;BSTree是二叉树，它保护了二叉树的根结点mRoot；mRoot是BSTNode类型，而BSTNode是二叉查找树的结点，它是BSTree的内部类。BSTNode包含二叉查找树的几个基本信息：
1. key -- 关键字，是用来对二叉查找树的结点进行排序的。
2. left -- 它指向当前结点的左孩子。
3. right -- 它指向当前结点的右孩子。
4. parent -- 它指向当前结点的父结点。

## 遍历
&emsp;&emsp;这里讲解前序遍历、中序遍历、后序遍历3种方式
### 前序遍历(先根遍历)
&emsp;&emsp;若二叉树非空，则执行以下操作：
1. 访问根结点；
2. 先序遍历左子树；
3. 先序遍历右子树；


        private void preOrder(BSTNode<T> tree){
            if (null != tree) {
                sysout(tree.key+" ");
                preOrder(tree.left);
                preOrder(tree.right);
            }
        }

        public void preOrder(){
            preOrder(mRoot);
        }


### 中序遍历(中根遍历)
&emsp;&emsp;若二叉树非空，则执行以下操作：
1. 中序遍历左子树；
2. 访问根结点；
3. 中序遍历右子树。


        private void inOrder(BSTNode<T> tree){
            if (null != tree) {
                inOrder(tree.left);
                sysout(tree.key+" ");
                inOrder(tree.right);
            }
        }

        public void inOrder() {
            inOrder(mRoot);
        }


### 后序遍历(后根遍历)
&emsp;&emsp;若二叉树非空，则执行以下操作：
1. 后序遍历左子树；
2. 后序遍历右子树；
3. 访问根结点


        private void postOrder(BSTNode<T> tree) {
            if (null != tree) {
                postOrder(tree.left);
                postOrder(tree.right);
                sysout(tree.key+" ");
            }
        }

        public void postOrder() {
            postOrder(mRoot);
        }


> 看看下面这棵树的各种遍历方式：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/遍历示例.jpg)

&emsp;&emsp;对于上面的二叉树而言：
1. 前序遍历结果：3 1 2 5 4 6
2. 中序遍历结果：1 2 3 4 5 6
3. 后序遍历结果：2 1 4 6 5 3

## 查找
### 递归版本的代码


        private BSTNode<T> search(BSTNode<T> x, T key) {
            if (null == x) {
                return x;
            }
            int cmp = key.compareTo(x.key);
            if (cmp < 0) {
                return search(x.left, key);
            }
            if (cmp > 0) {
                return search(x.right, key);
            }
            return x;
        }

        public BSTNode<T search(T key) {
            return search(mRoot, key);
        }


### 非递归版本的代码


        private BSTNode<T iterativeSearch(BSTNode<T> x, T key) {
            while (null != x) {
                int cmp = key.compareTo(x.key);

                if (cmp < 0){
                    x = x.left;
                }else if (cmp > 0) {
                    x = x.right;
                }else {
                    return x;
                }
            }
            return x;
        }

        public BSTNode<T> iterativeSearch(T key) {
            return iterativeSearch(mRoot, key);
        }


## 最大值和最小值
### 查找最大值的代码


        private BSTNode<T> maximum(BSTNode<T> tree) {
            if (null == tree) {
                return null;
            }
            while(null != tree.right) {
                tree = tree.right;
            }
            return tree;
        }

        public T maximum() {
            BSTNode<T> p = maximum(mRoot);
            if (null != p) {
                return p.key;
            }
            return null;
        }


### 查找最小值的代码


        private BSTNode<T> minimum(BSTNode<T> tree) {
            if (null != tree) {
                return null;
            }
            while (null != tree.left) {
                tree = tree.left;
            }
            return tree;
        }

        public T minimum() {
            BSTNode<T> p = minimum(mRoot);
            if (null != p) {
                return p.key;
            }
            return null;
        }


## 前驱和后继

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/前驱后继.png)

- 结点的前驱： 是该结点的左子树中的最大结点
&emsp;&emsp; 如图结点4的前驱结点是3，2的前驱结点是1，6的前驱结点是5
- 结点的后继： 是该结点的右子树中的最小结点
&emsp;&emsp; 如图结点7的后继结点是8，5的后继是6，2的后继是3

### 查找前驱结点代码

        // 找结点x的前驱结点，即，查找：二叉树中数据小于该结点的最大结点
        public BSTNode<T> predecessor(BSTNode<T> x) {
            // 如果x存在左子树，则x的前驱结点为：以其左孩子为根的子树的最大结点
            if (null != x.left) {
                return maximum(x.left);
            }

            // 如果没有左子树，则有以下两种可能：
            // 1. x是一个右孩子，则x的前驱结点为它的父结点
            // 2. x是一个左孩子，则查找x的最低的父结点，并且该父结点要具有右孩子
            BSTNode<T> y = x.parent;
            // 一直找
            while(null != y && y.left == x) {
                x = y;
                y = y.parent;
            }
            return y;
        }


### 查找后继结点的代码


        // 找结点x的后继结点，即，查找二叉树中数据大于该结点的最小结点
        public BSTNode<T> successor(BSTNode<T> x) {
            // 如果x存在右子树，则x的后继结点为以其右孩子为根结点的子树的最小结点
            if (null != x.right) {
                return minimum(x.right);
            }
            // 如果x没有右孩子。则x右以下两种可能：
            // 1. x是一个左孩子，则x的后继结点为它的父结点
            // 2. x是一个右孩子，则查找x的最低的父结点，并且该父结点要具有左孩子，找到的这个最低的父结点就是x的后继结点
            BSTNode<T> y = x.parent;
            while(null != y && y.right == x) {
                x = y;
                y = y.parent;
            }
            return y;
        }


## 重点：插入结点


        /**
        * bst 二叉树
        * z 要插入的结点
        *
        /
        private void insert(BSTNode<T> bst, BSTNode<T> z) {
            int cmp;
            BSTNode<T> y = null;
            BSTNode<T> x = bst.mRoot;

            // 查找z的插入位置
            while (null != x) {
                y = x;
                cmp = z.key.compareTo(x.key);
                if (cmp < 0) {
                    x = x.left;
                } else {
                    x = x.right;
                }
            }
            z.parent = y;
            if (null == y) {
                bst.mRoot = z;
            } else {
                // 放左边还是右边
                cmp = z.key.compareTo(y.key);
                if (cmp < 0) {
                    y.left = z;
                } else {
                    y.right = z;
                }
            }
        }

        /**
        *
        * 新建结点key，并将其插入到二叉树中
        *
        */
        public void insert(T key) {
            BSTNode<T> z = new BSTNode<T>(key, null, null, null);
            if (null != z) {
                insert(this, z);
            }
        }


## 重点：删除结点
1. 如果z没有孩子结点，那么只是简单地将它删除，并修改它的父结点，用NULL来替换z
2. 如果z只有一个孩子，那么将这个孩子提升到树中z的位置，并修改z的父结点，用z的孩子来替换z
3. 如果z有两个孩子，那么找到z的后继y(一定在z的右子树中)，并让y占据树中z的位置。z原来右子树部分成为y的新的右子树，并且z的左子树成为y的新的左子树。这种情况稍显麻烦，因为还与y是否为z的右孩子相关

如下图所示：

- 如果z没有左孩子，那么用其右孩子来替换z，这个右孩子可以是NIL，也可以不是。当z的右孩子是NIL时，此时这种情况归为z没有孩子结点的情形。当z的右孩子非NIL时，这种情况就是z仅有一个孩子结点的情形，该孩子是其右孩子.


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/删除结点1.png)


- 如果z仅有一个孩子且为其左孩子，那么用其左孩子来替换z


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/删除结点2.png)


- 否则，z既有一个左孩子又有一个右孩子。我们要查找z的后继y，这个后继位于z的右子树中并且没有左孩子。现在需要将y移出原来的的位置进行拼接，并替换树中的z
- 如果y是z的右孩子，那么用y替换z，并仅留下y的右孩子


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/删除结点3.png)

- 负责，y位于z的右子树中但并不是z的右孩子，在这种情况下，先用y的右孩子替换y，然后再用y替换z。


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/BSTree/删除结点4.png)

&emsp;&emsp;从一个二叉搜索树中删除结点z。结点z可以是根，可以是结点q的一个左孩子，也可以是q的一个右孩子。(a)结点z没有左孩子。用其右孩子r来替换z，其中r可以是NIL，也可以不是。(b)结点z有一个左孩子l但没有右孩子。用l来替换z。(c)结点z有两个孩子，其左孩子是结点l，其右孩子y还是其后继，y的右孩子是结点x。用y替换z，修改使l成为y的左孩子，但保留x仍为y的右孩子。(d)结点z有两个孩子(左孩子l和右孩子r)，并且z的后继y$\neq$ r位于以r为根的子树中。用y自己的右孩子x来替代y，并且y为r的双亲。然后，再置y为q的孩子和l的双亲。

&emsp;&emsp;为了在二叉搜索树内移动子树，定义一个子过程transplant,它是用另一课子树替换一颗子树并成为其双亲的孩子结点。当transplant用一颗以v为根的子树来替换一颗以u为根的子树时，结点u的双亲就变为结点v的双亲，并且最后v成为u的双亲的相应孩子.

        transplant(BSTNode<T> u, BSTNode<T> v){
            if (u.parent == NULL) {
                mRoot = v;
            } else if (u == u.parent.left) {
                u.parent.left = v;
            } else {
                u.parent.right = v;
            }
            if (v != NULL) {
                v.parent = u.parent;
            }
        }


&emsp;&emsp;注意到，transplant并没有处理v.left和v.right的更新；这些更新都是有transplant的调用者来负责的。下面是从二叉搜索树中删除结点z的删除过程：


        /**
        * 删除结点
        *
        */
        public void remove(T key) {
            BSTNode<T> z,y;

            if ((z = search(mRoot, key)) != null) {
                if (z.left == NULL) {
                    transplant(z, z.right);
                } else if (z.right == NULL) {
                    transplant(z, z.left);
                } else {
                    y = minmum(z.right)
                    if (y.parent != z) {
                        transplant(y, y.right);
                        y.right = z.right;
                        y.right.parent = y;
                    }
                    transplant(z,y);
                    y.left = z.left;
                    y.left.parent = y;
                }
            }
        }


&emsp;&emsp;remove过程处理4中情况如下。9-10行处理结点z没有左孩子的情况，第11-12行处理z有一个左孩子但没有右孩子情况。第13-22行处理剩下的两种情况，也就是z有两个孩子的情形。第14行查找结点y，它是z的后继。因为z的右子树非空，这样后继一定是这个子树中具有最小关键字的结点，因此就调用minmum(z.right)。如前所述，y没有左孩子。将y移出它原来位置进行拼接，并替换树中的z。如果y是z的右孩子，那么第20-22行用y替换z并成为z的双亲的一个孩子，用z的左孩子替换y的左孩子。如果y不是z的左孩子，第16-18行用y的右孩子替换y并成为y的双亲的一个孩子，然后将z的右孩子转变为y的右孩子，最后第20-22行用y替换y成为z的双亲的一个孩子，再用z的左孩子替换为y的左孩子。

&emsp;&emsp;除了第14行调用minmum之外，remove的每一行，包括调用translpant，都只花费常数时间。因此，在一颗高度为h的树上，remove的运行时间为O(h).

## 打印


        /*
        *   打印"二叉查找树"
        *
        *   key         结点的值
        *   direction   0，表示该结点是根结点
        *               -1，表示该节点是它的父节点左孩子
        *               1，表示该节点是他的父节点右孩子
        */
        private void print(BSTNode<T> tree, T key, int direction) {
            if (null != tree) {
                if (0 == direction) {
                    sout("%2d is root\n", tree.key);
                } else {
                    sout("%2d is %2d`s %6s child\n", tree.key, key, direction==1?"right":"left");
                }
                print(tree.left, tree.key, -1);
                print(tree.right, tree.key, 1);
            }
        }

        public void print() {
            if (null != mRoot) {
                print(mRoot, mRoot.key, 0);
            }
        }


## 销毁


        private void destroy(BSTNode<T> tree) {
            if (null == tree) {
                return;
            }
            if (null != tree.left) {
                destroy(tree.left);
            }
            if (null != tree.right) {
                destroy(tree.right);
            }
            tree = null;
        }

        public void clear() {
            destroy(mRoot);
            mRoot = null;
        }


## golang实现二叉搜索树


        package btree

        import (
            "bytes"
            "fmt"
            "github.com/emirpasic/gods/trees"
            "github.com/emirpasic/gods/utils"
            "strings"
        )

        // Entry represents(代表) the k-v pair(对) contained within nodes
        type Entry struct {
            Key     interface{}
            Value   interface{}
        }

        // Node is a single element within the tree
        type Node struct {
            Parent      *Node
            Entries     []*Entry    // Contained keys in node
            Children    []*Node    // Children nodes
        }

        // Tree holds elements of the B-tree
        type Tree struct {
            Root        *Node               // Root node
            Comparator  utils.Comparator    // key comparator
            size        int                 // Total number of keys in the tree     实际树中结点个数量
            m           int                 // order (maximum number of children) 预期树中孩子的最大数量
        }

        // NewWith instantiates a B-tree with the order and a custom key comparator
        func NewWith(order int, comparator utils.Comparator) *Tree {
            if order < 3 {
                panic("Invalid order, should be at least 3")
            }
            return &Tree{m: order, Comparator: comparator}
        }

        func NewWithIntComparator(order int) *Tree {
            return NewWith(order, utils.IntComparator)
        }

        func NewWIthStringComparator(order int) *Tree {
            return NewWith(order, utils.StringComparator)
        }

        // Put inserts k-v pair node into the tree.
        // If k already exists, then its value is updated with the new v
        // K should adhere to the comparator`s type assertion, otherwise method panics（key 必须是comparator的断言类型，否则会异常中断）
        func (tree *Tree) Put(key interface{}, value interface{}) {
            entry := &Entry{key, value}
            if tree.Root == nil {
                tree.Root = &Node{Entries: []*Entry{entry}, Children: []*Node{}}
                tree.size++
                return
            }
            if tree.insert(tree.Root, entry) {
                tree.size++
            }
        }

        func (tree *Tree) insert(node *Node, entry *Entry) (inserted bool) {
            if tree.isLeaf(node) {
                return tree.isertIntoLeaf(node, entry)
            }
            return tree.insertIntoInternal(node, entry)
        }

        func (tree *Tree) isLeaf(node *Node) bool {
            return len(node.Children) == 0
        }

        func (tree *Tree) insertIntoLeaf(node *Node, entry *Entry) (inserted bool) {
            insertPosition, found := tree.search(node, entry.Key)
            if found {
                // update v
                node.Entries[insertPosition] = entry
                return false
            }
            // insert entry`s k in the middle of the node
            node.Entries = append(node.Entries, nil)
            // 后移一位
            copy(node.Entries[insertPosition+1:], node.Entries[insertPosition:])
            node.Entries[insertPosition] = entry
            tree.split(node)
            return true
        }

        // search searches only within the single node among its entries
        func (tree *Tree) search(node *Node, key interface{}) (index int, found bool) {
            low, high := 0, len(node.Emtries) - 1
            var mid int
            for low <= high {
                mid = (high + low) / 2      // mid = (high + low) >> 1
                compare := tree.Comparator(key, node.Entries[mid].key)
                switch {
                    case compare > 0: 
                        low = mid + 1
                    case compare < 0:
                        high = mid - 1
                    case compare == 0: 
                        return mid, true
                }
            }
            return low, false
        }

        func (tree *Tree) split(node *Node) {
            if !tree.shouldSplit(node) {
                return
            }
            if node == tree.Root {
                tree.splitRoot()
                return
            }
            tree.splitNonRoot(node)
        }

        func (tree *Tree) shouldSplit(node *Node) bool {
            return len(node.Entries) > tree.maxEntries()
        }

        func (tree *Tree) middle() int {
            return (tree.m - 1) / 2
        }

        func (tree *Tree) splitRoot() {
            middle := tree.middle()

            left := &Node{Entries: append([]*Entry(nil), tree.Root.Entries[:middle]...)}
            right := &Node{ENtries: append([]*ENtry(nil), tree.Root.Entries[middle+1:]...)}

            // move children from the node to be split into left and right nodes
            if !tree.isLeaf(tree.Root) {
                left.Children = append([]*Node(nil), tree.Root.Children(:middle+1)...)

            }
        }
