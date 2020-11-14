
# 摘要

&emsp;&emsp;**之所以要将Git分支原理单独讲解，主要是想好好理解一下Git的分支模型，这对以后的分支管理学习大有裨益。**

## 简介

&emsp;&emsp;为了真正理解Git处理分支的方式，我们需要回顾以下Git是如何保存数据的。

&emsp;&emsp;获取你还记得第一节的内容，Git保存的不是文件的变化，而是一系列的文件快照。在提交操作时，Git会保存一个提交对象(`commit object`)。知道了Git保存数据的方式，我们可以很自然的想到---该提交对象会包含一个指向暂存内容快照的`指针`。不仅仅是指针，该对象还包含了作者、email、提交信息以及`指向它父对象的指针`。首次提交产生的提交对象没有父对象，普通提交操作产生的提交对象有一个父对象，而由多个分支合并产生的提交对象有多个父对象。

&emsp;&emsp;为了更加形象地说明，我们假设现在有一个工作目录，里面包含了三个将要被暂存和提交的文件。暂存操作会为每个文件计算hash，然后会把当前版本的文件快照保存到Git仓库中(Git使用blob对象保存它)，最终将hash加入到暂存区等待提交。


        $ git add README test.rb LICENSE
        $ git commit -m 'the initial commit of my project'


&emsp;&emsp;当使用`git commit`提交时，Git会先计算每个子目录(本例指含有一个根目录)的hash，然后在Git仓库中这些hash保存为`树对象`.之后，Git便会创建一个提交对象，它包含了上面提到的信息外，还包含`指向这个树对象的指针`。如此一来，Git就可以在需要的时候重现此次保存的快照。

&emsp;&emsp;ok,现在Git仓库中有五个对象：三个blob(保存文件快照)，一个树对象(记录目录结构和blob对象索引)以及一个提交对象(包含指向树对象和所有提交信息)


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/首次提交对象及树对象.png)

&emsp;&emsp;做些修改后再次提交，产生的提交对象会包含一个指向上次提交对象的指针。


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/提交对象及父对象.png)

&emsp;&emsp;**Git分支，其实本质上仅仅是指向提交对象的可变指针。**Git默认分支是master，在多次提交后，其实已经又一个指向最后那个提交对象的master分支，它会在每次提交的操作中自动向前移动。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/分支及其提交历史.png)

## 创建分支原理

&emsp;&emsp;Git是如何创建新分支的呢？很简单，它只是为你创建了一个可以移动的新的指针。比如，创建一个testing分支，`git branch`:


        $ git branch testing


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/两个指向相同提交历史的分支.png)

&emsp;&emsp;so，Git怎么知道当前在哪一个分支上？它有一个名为HEAD的特殊指针。在本例中，你仍然在master分支上，因为`git branch`仅仅创建了一个新分支，并不会自动切换过去.

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/HEAD指向当前所在的分支.png)

&emsp;&emsp;可以使用`git log --decorate`查看各个分区当前所指的对象。如下，master和testing分支均指向hash f30ab开头的提交对象


        $ git log --oneline --decorate
        f30ab (HEAD, master, testing) add feature


## 切换分支原理

&emsp;&emsp;切换分支`git checkout`.


        $ git checkout testing


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/HEAD指向当前所在分支.png)

&emsp;&emsp;so，这样的实现方式会带来什么好处呢？我们不妨再提交一次


        $ git commit -am 'made a change'


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/HEAD分支随着提交操作自动向前移动.png)

&emsp;&emsp;如图，testing分支向前移动了，但是master分支却没有，它仍然指向运行checkout之前的提交对象。我们切换回master分支看看：


        $ git checkout master


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/检出时HEAD随之移动.png)

&emsp;&emsp;这条命令做了两件事：

1. 使HEAD指向master分支
2. 将工作目录恢复成master分支所指向的快照内容

&emsp;&emsp;如果你现在做修改操作的化，项目将始较旧的版本。这就是忽略testing分支所做的事情，以便于向另一个方向开发。
