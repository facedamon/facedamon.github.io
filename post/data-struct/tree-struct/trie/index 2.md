
# 摘要

&emsp;&emsp;前缀树，又称`字典树`。之前在[手写http框架(路由)](https://facedamon.github.io/post/golang/http/golang-%E6%89%8B%E5%86%99http%E6%A1%86%E6%9E%B6-router/)中已经略有介绍。这次我们专门做一个关于前缀树的详细介绍及应用。

&emsp;&emsp;前缀树主要应用于**统计，排序和保存大量的字符。** 最常用的莫过于HTTP路由分组和搜索引擎词频统计。

> 思考

&emsp;&emsp;`给出n个单词和m个询问，每次询问一个单词，回答这个单词是否在单词表中出现过？
`

&emsp;&emsp;很多人的第一想法是使用map，这样足够简单。

假设有100W个单词，难道要把100W个单词都存到内存里吗？然后还要祈求map不会出现碰撞，在1s内查询出结果。这显然不可能。

&emsp;&emsp;而且这样做还有个问题，假设单词只包含英文字母，那么就只会出现26个字符的不同组合，如果一股脑存储的话，重复率极高。我们可以把每个单词拆分成若干字母，总共有26中可能，如果一个字母已经存在，则不再存储而是存储下一个字母，并把前缀字母与自己关联起来。这就是我们今天要说的`前缀树`。

## 引入

&emsp;&emsp;假设现在有5个单词：ab, ad, be, cd, cg, 我们把他们拆分成单个字母，并从上往下深度优先构建这棵树。(与之相反的是从左往右广度优先)


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/trie.png)


&emsp;&emsp;从这可树上可以看出，一下几个要点：

1. **有一个root节点，并且root节点没有元素。**
2. **root节点可以有多个子节点。**
3. **非root节点存储一个元素。**
4. **非root节点也可以有多个子节点。**
5. **所有节点都是自顶向下。**

&emsp;&emsp;OK, 到此我们可以定义我们的Node节点属性了。在此之前我们还要抉择一个事情，就是一个节点的子节点是使用`数组`呢还是`Map`来存储, 像这样：

        childs []*node;
        childs map[key]node;

&emsp;&emsp;我们都知道数组是一段连续内存地址，在知道下标的前提下，查询速度非常快，但是插入和删除在元素量很多的情况下会愈来愈慢，因为要重排整个数组。Map内部是一个十字双向链表(就是key为数组，value为数组指针指向的链表)，在插入和删除元素上非常快，在查询上会很慢，因为要从head开始遍历每个next指针。

&emsp;&emsp;对于trie树来说，如果数据量非常大(十万及以上)就采用数组，否则采用map。在使用上没有太大的区别。我们今天的应用案例是`敏感词过滤`，就是从一堆字符中命中查询的字符。

## 定义

        type node struct {
            isRootNode bool
            isPathEnd bool
            character rune
            children map[rune]*node
        }

        type trie struct {
            root *node
        }

        // 构造函数
        // 其实root是有元素的，是一个特殊元素。
        func NewTrie() *Trie {
            return &Trie{
                root: NewRootNode(0)
            }
        }

        // 构建root节点
        func NewRootNode(character rune) *Node {
            return &Node{
                isRootNode: true,
                character: character,
                children: make(map[rune]*node, 0),
            }
        }

        // 构建非root节点
        func NewNode(character rune) *Node {
            return &Node{
                character: character,
                children: make(map[rune]*node, 0),
            }
        }

&emsp;&emsp;`node`节点就是上图中的每个节点对象。`trie`是root节点。children 是一个map对象，存储自身和自己点的映射关系，映射关键词是character元素, 此处采用**rune**类型，rune是golang特有的针对unicode编码的实现。

## 添加

&emsp;&emsp;一颗树想要使用必然是现有创建树，就是插入元素。我们先设想一下，现在有一份包含汉字的字符串要插入这棵树，按照我们之前的理论，肯定是先将这串字符转换成UNICODE码数组，这样大家都有了相同的前缀标准了。假设这棵树已经有一些元素了，一个字符串想要进入，就拿自己的UNICODE码一个一个的去匹配节点，而且是深度优先的匹配，什么意思呢？就是一旦map里有值，那就说明有了相同的前缀了，继续往下查，直到找到map没值的情况或叶子节点。

            // 批量添加关键词 public
            func (t *trie) Add(words ...string) {
                for _, word := range words {
                    t.add(world)
                }
            }

            // 单词添加关键词 private
            func (t *trie) add(word string) {
                var current = t.root
                var runes = []rune(word)
                for position := 0; position < len(runes); position++ {
                    r := runes[position]
                    if next, ok := current.children[r]; ok {
                        current = next
                    } else {
                        // 没查到则新建
                        newNode := NewNode(r)
                        current.children[r] = newNode
                        current = newNode
                    }
                    // 如果到最后一个字符了
                    if position == len(runes) - 1 {
                        current.isPathEnd = true
                    }
                }
            }

## 删除

&emsp;&emsp;有添加肯定有删除。说到删除，做过关系型数据库设计的人肯定有过这样的经历，一开始认为删除只是简单的delete，但是后来发现我们给误删了，没法了，或者这些数据本身是有价值的，不能删除。这时候就只能状态转化了，转化成客户看不到的状态，等到非删不可的时候再去删除。OK，我们的trie就采用这种思维。

&emsp;&emsp;还有个问题，怎么删除？你总不能把整个单词都删了吧，这样其它单词的前缀怎么办，成了空指针了。所以我们采用删除(状态转换)最后一个UNICODE码的方式。

        // 批量删除关键词 public
        func (t *trie) Del(words ...string) {
            for _, word := range words {
                t.del(word)
            }
        }

        // 单词删除关键词 private
        func (t *trie) del(word string) {
            var current = t.root
            var runes = []rune(word)
            for position := 0; position < len(runes); position++ {
                r := runes[position]
                if next, ok := current.children[r]; !ok {
                    // 没查到
                    return
                } else {
                    current = next
                }
                // 查到了最后一个UNICODE码
                if position == len(runes) - 1 {
                    current.softDel()
                }
            }
        }

        // 软删除
        func (n *node) softDel() {
            n.isPathEnd = false
        }

## 查找(校验)

&emsp;&emsp;已经把字符给软删除了，我们想测试一下，那最好的办法就是查找。我们先来看看几种可能的情况：

- 敏感词：`fucking`

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/tree/trie/sensitive.png)

1. 查询词为 `read the fucking code`

&emsp;&emsp;敏感词在**中间**， **横向**匹配字母`r`没找到，切换到下一个字母`e`，继续这个操作。当匹配到字母`f`时，开始**纵向**查找，完整的匹配到了`fucking`关键词，这时查询词自身还没有结束，但是搜索树已经到了**叶子节点**，那我们可以确定`fucking`这个词肯定是敏感词，但是`code`是不是就不知道了，直接返回第一个敏感词。

2. 查询词为 `fucking asshole`

&emsp;&emsp;敏感词在**开头**，**横向**匹配字母`f`找到了，切换到下一个字母`u`开始**纵向**查找，完整的匹配到了`fucking`关键词，同样，查询词自身还没有结束，但是搜索树已经到了**叶子节点**，那我们可以确定`fucking`这个词肯定是敏感词, 直接返回第一个敏感词。

3. 查询词为 `what fucking`

&emsp;&emsp;敏感词在**结尾**， **横向**匹配字母`w`没找到，切换到下一个字母`h`, 继续这个操作。当匹配到字母`f`时，开始**纵向查找**，完整的匹配到了`fucking`关键词，这时查询词已经结束了， 就不需要再判断树是否到达叶子节点了，直接返回第一个敏感词。

4. 查询词为 `ee fuc gt`

&emsp;&emsp;这里面没有敏感词，我们试试看看它怎么走过这棵树。**横向**匹配字母`e`找到了，切换到下一个字母`e`开始**纵向**查找， 查询到`f`时发现没有匹配的，这时候，**就要跳出这个分支了**，当前查询字母已经变成了`f`了，但是我们应该从第二个单词`e`开始继续**横向**查找，找到了，切换到下一个字母`e`开始**纵向**查找， 查询到`f`时发现没有匹配的，这时候继续**跳出这个分支**，当前查询字母已经变成了`f`，我们应该从第三个单词`f`继续**横向**查找，找到了，切换到下一个字母`u`开始**纵向**查找，查询到`g`时发现没有匹配，这时候继续**跳出这个分支**，当前查询字母已经是`g`了， 但是我们应该从第四个单词'u'继续**横向**查找。。。。


1. 命中敏感词的一个必要条件是： 搜索到了叶子节点**
2. 跳出分支的必要条件是：没找到或者查询词已经结束了但是还没到叶子节点
3. 跳出分支时做了什么： 查询字母是上次查询字母的下一个字母
4. 如何定位敏感词：标记第一个命中字母和`连续`的最后一个命中字母
5. 在每次跳出分支时的开始查询字母都与上次连续(标记第一个命中字母)

> 题外思考，这种与上次字母连续的查找是不是可以优化，参考一下KMP算法，本次不做讨论。


        // Validate 验证字符串是否合法， 如果不合法则返回false和检测到的第一个敏感词
        func (t *trie) Validate(text string) (bool, string) {
            var (
                // 跳出分支的开始节点
                parent = t.root
                // 查询词UNICODE码
                runes = []rune(text)
                // 命中单词的第一个字母下标
                left = 0
            )

            for position := 0; position < len(runes); position++ {
                // 横向查找
                current, ok :=  parent.children(runes[position])
                // 在没有找到或者查询词已经结束的情况下需要跳出当前分支
                if !ok || (position == len(runes) - 1 && !current.IsPathEnd()){
                    parent = t.root
                    position = left
                    left++
                    continue
                }

                // 到叶子节点了
                if current.IsPathEnd() && left <= position {
                    return false, string(runes[left: position + 1])
                }
                //纵向查找
                parent = current
            }
        }

        func (n *node) IsPathEnd() bool {
            return n.IsPathEnd
        }
## 测试

        func TestTrie_Validate(t *testing.T) {
	        tree := New()
            words := []string{"FUCK", "fuck"}
	        err := tree.Add(words)
	        if err != nil {
		        t.Error(err)
	        }
	        t.Log(tree.Validate("你好FUCK的fuck了"))
        }
        // output: false FUCK
</br>
</br>

&emsp;&emsp;至此，其实我们已经基本完成了字符串匹配的核心算法了。上面的方法只能找出第一个被命中的敏感词，如果想要把查询词中的敏感词都查出来呢？其实没什么区别，只要把每个命中的单词都存储起来就行了。

        func (t *trie) FindAll(text string) []string {
            var (
                matches []string
                parent = t.root
                runes = []rune(text)
                left = 0
            )

            for position := ; position < len(runes); position++ {
                current, ok := parent.children[runes[position]]
                if !ok || (!current.IsPathEnd() && position == len(runes) - 1) {
                    parent = t.root
                    position = left
                    left++
                    continue
                }

                if current.IsPathEnd() && left <= position {
                    matches = append(matches, string(runes[left: position + 1]))
                }

                parent = current
            }
            return matches
        }


## 测试

        func TestTrie_FindAll(t *testing.T) {
	        tree := New()
            words := []string{"FUCK", "fuck"}
            err := tree.Add(words)
	        if err != nil {
	        	t.Error(err)
	        }
	        t.Log(tree.FindAll("F你好FUCK的fuck了"))
        }
        // output: [FUCK fuck]

## 过滤

&emsp;&emsp;做了校验和查询，那接下来肯定当属过滤了，在正常的业务系统中除了要应用校验功能，展示阶段肯定也就进行数据过滤。

&emsp;&emsp;在我们前面已经探讨了查询的逻辑后，其实，已经基本涵盖了整个算法。思考一下，最终的过滤效果是什么样，像这样: `niuFUCKwofuck的nn六4gh` >> `niuwo的nngh` ，没错，这就是我们想要的最终结果。常规思路想一下，肯定是要初始化一个返回的数组的，那这个数组的元素肯定是只存储命中词之外的词，那每次命中之外是什么时候？**就是跳出分支的时候。**

        func (t *trie) Filter(text string) string {
            var (
                parent = t.root
                left = 0
                runes = []rune(text)
                rs = make([]rune, 0, len(runes))
            )

            for position := 0; position < len(runes); position++ {
                current, ok "= parent.children[runes[position]]
                if !ok || (!current.IsPathEnd() && position == len(runes) - 1) {
                    // 就是这里了
                    rs = append(rs, runes[left])
                    parent = t.root
                    position = left
                    left++
                    continue
                }

                if current.IsPathEnd() {
                    // 这里不再存储命中词了
                    // 需要换一个分支了
                    left = position + 1
                    parent = t.root
                } else {
                    parent = current
                }
            }
            return string(rs)
        }

## 测试

            func TestTrie_Filter(t *testing.T) {
	            tree := New()
                words := []string{"FUCK", "fuck", "六4"}
                err := tree.Add(words)
	            if err != nil {
	            	t.Error(err)
	            }
	            t.Log(tree.Filter("niuFUCKwofuck的nn六4gh"))
            }
            //output: niuwo的nngh

## 替换

&emsp;&emsp;还有最后一个功能性函数，就是替换，就像你在页面上输入密码一样，敏感词被命中后会直接被替换成指定字符。like this: `niuFUCKwofuck的nn六4gh` >> `niu####wo####的nn##gh`

&emsp;&emsp;回顾一下**FindAll**函数是怎么做的，它是从runes里面分片存储的，那我们只要在存储的时候给它偷梁换柱不就行了。

        func (t *trie) Replace(text string, character rune) string {
            var (
                parent = t.root
                runes = []rune(text)
                left = 0
            )

            for position := 0; position <len(runes); position++ {
                current, ok := parent.children[runes[position]]
                if !ok || (!current.IsPathEnd() && position == len(runes) - 1) {
                    parent = t.root
                    position = left
                    left++
                    continue
                }

                if current.IsPathEnd() && left <= position {
                    // 就是这里了
                    for i := left; i <= position; i++ {
                        runes[i] = character
                    }
                }

                parent = current
            }
            return string(runes)
        }

## 测试


            func TestTrie_Filter(t *testing.T) {
	            tree := New()
                words := []string{"FUCK", "fuck", "六4"}
                err := tree.Add(words)
	            if err != nil {
	            	t.Error(err)
	            }
	            t.Log(tree.Replace("niuFUCKwofuck的nn六4gh", '#'))
            }
            //output: niu####wo####的nn##gh

&emsp;&emsp;至此，字典树的核心功能都已经完成了，只要数据量在十万级一下，就能够做到秒级命中，十万级以上就得改用数组了。

# 总结

&emsp;&emsp;不知道大家有没有发现，一个有趣的地方：任何复杂的算法，其实都是手工分析的表达，你看别人的代码，其实是在读作者的思想，没有任何一个算法会凭空出现，都是有理论指导的。还有一个地方值得我们学习，在学习一门新的算法的时候，不要急于求成，盲目抄代码，虽然抄代码是工程师进步的必经之路，但是一定要过一遍手工过程，只有自己手工过了，心里才有底，才明白这个地方为什么要这样写， 否则你就像农夫山泉一样，只是个搬运工。

&emsp;&emsp;我在文章中留下了一个悬念，那就是使用**KMP**算法来优化`left`效率低下的问题，只能等到下次总结了。
