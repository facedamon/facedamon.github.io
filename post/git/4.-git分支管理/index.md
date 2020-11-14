
# 摘要

&emsp;&emsp;案例，实际工作中你会用到类似的流程。

1. 开发某个网站
2. 为实现某个新的需求，创建分支
3. 在该分支上开展工作

&emsp;&emsp;正在此时，你接到通知，有一个很严重的bug需要紧急修复。你将按照如下方式来处理：

1. 切换到线上分支(production branch)
2. 为这个紧急任务新建一个新分支，并在其中修复它
3. 测试通过后，切换回线上分支，然后合并这个补丁分支，最后推送到线上
4. 切换回你最初的工作分支，继续工作

## 新建分支

&emsp;&emsp;我们就使用上述案例来讲解分支的创建和切换。假设，你当前的工作目录已经有一些提交对象了，如图所示。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/一个简单提交历史.png)

&emsp;&emsp;ok, 现在你需要解决问题编号#53,新建一个分支并同时切换到该分支。`git checkout -b [branchname]`


        $ git checkout -b iss53
        Switch to a new branch 'iss53'

# 上述命令等效于：
        $ git branch iss53
        $ git checkout iss53


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/创建一个新分支指针.png)

&emsp;&emsp;然后你在iss53分支上开始修复问题#53，此时iss53分支开始向前推进。


        $ vim index.html
        $ git commit -am 'added a nre footer [iss 53]'


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/分支随着工作的进展向前推进.png)

&emsp;&emsp;假设，你此时又接到通知，有一个新的问题需要紧急修复。那么我们站在Git分支的角度来思考，你仅仅需要做的是：提交你当前iss53分支的工作目录，然后再切换回master分支`git checkout master`，此时你的工作目录和你创建iss53分支时一摸一样，ok，现在你可以专心的修改新问题了。你需要创建一个针对新紧急问题的分支(hotfix branch):


        $ git checkout -b hotfix
        Switch to a new branch 'hotfix'
        $ vim index.html
        $ git commit -am 'fixed the broken email address'
        [hotfix 1fb7853] fixed the broken email address
            1 file changed, 2 insertions(+)


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/基于master分支的紧急问题分支.png)

&emsp;&emsp;假设你已经在hotfix分支完成了测试，现在需要合并到master分支并部署到线上，你可以使用`git merge`目录来达到此目的：


        $ git checkout master
        $ git merge hotfix
        Updating f42c576..3a0874c
        # 注意这里[快进]
        Fast-forward
            index.html | 2 ++
            1 file changed, 2 insertions(+)


&emsp;&emsp;注意上面的\[快进 fast-formard\]，由于当前master分支所指向的提交对象是你当前分支提交对象的父对象，所以merge时，只是简单的将master指针向前移动。当你试图合并两个分支时，如果顺着一个分支走下去能够到达另一个分支，那么Git在合并两者时，只会简单的将指针向前推进，因为这种情况下的合并操作没有需要解决的分歧。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/master被快进到hotfix.png)

&emsp;&emsp;在你修复了hotfix问题后，你准备回到被打断之前的工作中。然而，你应该先删除hotfix分支，因为你已经不需要它了，master分支已经指向了同一个位置。你可以使用`git branch -d`来删除分支.


        $ git branch -d hotfix
        Deleted branch hotfix (3a0874c)


        $ git checkout iss53
        $ Switched to branch 'iss53'
        $ vim index.html
        $ git commit -am 'finished the new footer [issue 53]'
        [iss53 ad82d7a] finished the new footer [issue 53]
            1 file changed, 1 insertion(+)


![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/继续在iss53分支上工作.png)

&emsp;&emsp;你在hotfix分支上所做的工作并没有合并到iss53分支中。如果你需要拉取hotfix所做的修改，你可以使用`git merge master`目录将master分支合并到iss53分支，或者你也可以等到iss53分支完成任务后，再将其合并回master分支。

## 分支合并

&emsp;&emsp;假设此时你已经修正了#53问题，打算将iss53分支合并到master分支。


        $ git checkout master
        $ Switch to branch 'master'
        $ git merge iss53
        Merge made by the 'recursive' strategy.
        index.html | 1 +
            1 file changed, 1 insertion(+)


&emsp;&emsp;你会发现和上次的合并信息不太一样。在这种情况下，你的开发历史从一个更早的地方开始分叉开来，因为，master分支所在提交并不是iss53分支所在提交的直接父对象，Git不得不做一些额外工作。出现这种情况时，Git会使用两个分支的末端所指的快照(C4和C5)以及这两个分支的共同工作父对象(C2)，做一个简单的`三方合并`.

&emsp;&emsp;和之前将分支指针向前推进所不同的是，Git将此次三方合并的结果做了一个新快照并且自动创建一个新的提交对象指向它。这个被称作一次合并提交，`它的特别之处在于它有不止一个父对象`。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/一个合并提交.png)

&emsp;&emsp;既然你的修改已经合并进来了，就不再需要iss53分支了。现在你可以删除这个分支了.


        $ git branch -d iss53


## 遇到冲突时的分支合并

&emsp;&emsp;**如果你在两个不同的分支中，对同一个文件的同一部分进行了不同的修改，Git就没法顺利的合并它们。**如果你对#53问题的修改和有关hotfix的修改都车技到同一个文件的同一处，在合并时就会产生冲突：


        $ git merge iss53
        Auto-mergeing index.html
        CONFLICT (content): Merge conflict in index.html
        Automatic merge failed; fix conflicts and
        then commit the result.


&emsp;&emsp;此时Git做了合并，但没有自动创建新的合并提交。Git暂停下来，等待你去解决合并产生的冲突。你可以在合并冲突后使用`git status`来查看那些因为包含合并冲突而处于未合并(unmerged)状态的文件：


        $ git status
        On branch master
        You have unmerged paths
            (fix conflicts and run 'git commit')
        Unmerged paths:
            (use 'git add <file>...' to mark resolution)
            both modified: index.html
        no changes added to commit (use 'git add' and/or 'git commit -a')


&emsp;&emsp;任何因包含合并冲突而待解决的文件，都会以未合并状态标识出来。`Git会在有冲突的文件中加入标准的冲突解决标记，这样你可以打开这些包含冲突的文件手动解决。`


        <<<<<<< HEAD index.html
        <div id="footer">contract :
        email.support@github.com</div>
        =======
        <div id="footer">
        please contact us at support@github.com
        </div>
        >>>>>>> iss53: index.html


&emsp;&emsp;上述表示HEAD所指示的版本(也就是你master分支所在的位置，因为你在运行merge目录时已经检出到了该分支)在(======)的上半部分，而iss53分支所指示的版本在======的下半部分。为了解决冲突，你必须选择使用其中一个，或者自行合并这些内容。比如，你可以通过把这段内容替换成下面的样子以解决冲突：


        <div id="footer">
        please contact us at email.support@github.com
        </div>


&emsp;&emsp;在你解决了所有文件里的冲突后，对每个文件使用`git add`目录将其标记未已解决，你可以再次运行`git status`来确认所有的合并冲突都以解决。


        $ git status
        On branch master
        All conflicts fixed but you are still merging.
            (use 'git commit' to conclude merge)
        Changes to be committed:
            modified: index.html


&emsp;&emsp;如果你对结果满意，并确定之前冲突文件都以暂存，此时你可以使用`git commit`来完成合并提交：


        Merge branch 'iss53'
        Conflicts:
            index.html
        #
        # It looks like you may be committing a merge
        # If this is not correct, please remove the file
        #   .git/MERGE_HEAD
        # and try again

        # Please enter the commit message for your changes. Lines starting
        # with '#' will be ignored, and an empty message aborts the commmit.
        # On branch master
        # All conflicts flixed but you are still merging.
        # 
        # Changes to be committed:
        #   modified: index.html
        #


# 分支管理

&emsp;&emsp;`git branch` 命令不只是可以创建与删除分支，如果不加任何参数，会得到所有分支列表：


        $ git branch
        iss53
        * master
        testing


&emsp;&emsp;master前面的*字符，它代表现在HEAD指针检出的那一个分支。如果你需要查看每一个分支的最后一次提交，可以运行`git branch -v`


        $ git branch -v
            iss53   93b412c fix javascript issue
        *   master  7a98805 Merge branch 'iss53'
            testing 782fd34 add scott to the author list in the readmes


&emsp;&emsp;如果你想查看哪些分支已经合并到当前分支，可以运行`git branch --merged`


        $ git branch --merged
            iss53
        *   master


&emsp;&emsp;同理，如果你想查看所有包含未合并工作的分支，可以运行`git branch --no-merged`


        $ git branch --no-merged


# 分支工作流

&emsp;&emsp;在整个项目开发周期的不同阶段，你可以同时拥有多个开放的分支，你可以定期把某些分支合并入其它分支。许多Git开发者都喜欢使用这种方式来工作，比如只在master分支上保留完全稳定的代码。还有一些名为develop或者next的平行分支，被用来做后续开发或者测试稳定性，一旦这些分支达到稳定状态，就可以被合并入master分支。

&emsp;&emsp;通常把它们想像成流水线可能更好理解一点，哪些经过测试考验的提交会被遴选到更加稳定的流水线上去。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/渐进稳定分支流水线视图.png)

> 举个例子

&emsp;&emsp;假设，你在master分支上工作到C1，这时为了解决一个问题新建iss91分支，在iss91分支上工作到C4，然而对于那个问题你又有了新的想法，于是你再次新建iss91v2分支视图用另一种方法解决那个问题，接着你回到master分支工作了一会，你又冒出了一个不太确定的想法，你便在C10的时候新建dumbidea分支，并在上面开始了实验。那么你的提交历史看起来像下面这样：

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/拥有多个特性分支的提交历史.png)

&emsp;&emsp;假设，你决定使用iss91v2分支方案来解决这个问题，另外，你将dumbidea分支拿给你的同事看后，结果发现是惊人之举。这是你选择抛弃iss91分支(丢弃C5, C6)，然后把另外两个分支合并进master主分支。最终你的提交历史看起来像下面这样:

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/合并dumbidea和iss91v2分支之后的提交历史.png)

> 我们将会在**分布式Git**中揭示更多有关分支工作流的细节

# 远程分支

&emsp;&emsp;假设你的网络里有一个在git.ourcompany.com的Git服务器。如果你从这里克隆，Git的clone命令会为你自动将其命名为origin，拉取它的所有数据，创建一个指向它的master分支指针，并且在本地将其命名为origin/master。Git也会给你一个与origin的master分支在指向同一个地方的本地master分支，这样你就有工作的基础了。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/克隆之后的服务器与本地仓库.png)

&emsp;&emsp;如果你在本地master分支做了一些工作，于此同时，其他人推送提交到git.outcompary.com并更新了它的maste分支，那么你的提交历史将向不同的方向前进。另外，只要你不与origin服务器连接，你的origin/master指针就不会移动.

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/本地与远程的工作可以分叉.png)

&emsp;&emsp;如果要同步你的工作，运行`git fetch origin`命令。这个命令查找origin是哪一个服务器，从中抓取本地没有的数据，并且更新本地数据库，移动origin/master指针指向新的、更新后的位置。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/fetch更新你的远程仓库引用.png)

&emsp;&emsp;假设你还有一台内部Git服务器，仅仅用于sprint小组的开发工作，地址：git.team1.outcompary.com。你可以运行`git remote add`命令添加一个新的远程仓库。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/添加另一个远程仓库.png)

&emsp;&emsp;现在，运行`git fetch teamone`来抓取远程仓库teamone有而本地没有的数据。因为那台服务器上现有的数据是origin服务器上的一个子集，所以Git并不会抓取数据而是会设置远程跟踪分支teamone/master指向teamone的master分支

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/远程跟踪分支teamone.png)

&emsp;&emsp;综上，**远程跟踪分支是远程分支状态的引用.**它们是你不能移动的本地引用，当你做任何网络通信操作时，它们会自动移动。远程跟踪分支像是你上次连接到远程仓库时，哪些分支所处状态的书签。

&emsp;&emsp;**它们以(remote)/(branch)形式命名。**如果你想看你最后一次与远程仓库origin通信时master分支的状态，你可以查看origin/master分支。你与同时合作解决了问题并且它们推送了iss53分支，你可能自己本地有自己的iss53分支，在服务器上的分支会指向origin/iss53的提交。

# 推送

&emsp;&emsp;你可以将不愿意分享的内容放到私人分支上，而将需要和别人协作的内容推送到共来分支。

&emsp;&emsp;如果你希望和别人一起在名为serverfix的分支上工作，你可以像推送第一个分支那样，运行`git push (remote) (branch)`:


        $ git push origin serverfix
        Counting objects: 24, done.
        Delta compression using up to 8 threads.
        Compressing objects: 100%(15/15), done.
        Writing objects: 100% (24/24), 1.91kib | 0 bytes/s, done.
        Total 24 (delta 2), reused 0 (delta 0) To https://github.com/schacon/simplegit
        * [new branch] serverfix -> serverfix


&emsp;&emsp;你也可以运行`git push origin serverfix:serverfix`，它会推送本地的serverfix分支，将其作为远程仓库的serverfix分支。如果你并不想让远程仓库上的分支叫serverfix，可以运行git push origin serverfix:awesomebranch来将本地的serverfix分支推送到远程仓库上的awesomebranch分支。

&emsp;&emsp;下一次其他协作者从服务器上抓取数据时，他们会在本地生成一个远程分支origin/serverfix,指向服务器的servefix分支的引用。


        $ git fetch origin
        remote: Counting objects: 7, done.
        remote: Compressing objects: 100% (2/2), done.
        remote: Total 3 (delta 0), reused 3 (delta 0)
        Unpacking objects: 100% (3/3), done.
        From https://github.com/schacon/simplegit
        * [new branch]      serverfix    -> origin/serverfix


&emsp;&emsp;**注意**，当抓取到新的远程跟踪分支时，本地不会自动生成一份可编辑的副本。这种情况下，**不会有一个新的serverfix分支，只有一个不可以修改的origin/servefix指针。**

&emsp;&emsp;那么，你可以运行`git merge origin/serverfix`将这些工作合并到当前分支。如果你想在自己的serverfix分支上工作，可以将其建立在远程跟踪分支之上,这回给你一个用于工作的本地分支，并且起点位于origin/serverfix.


        $ git checkout -b serverfix
        origin/serverfix
        Branch serverfix set up to track remote
        branch serverfix from origin.
        Switched to a new branch 'serverfix'


# 跟踪分支

&emsp;&emsp;从一个远程跟踪分支检出一个本地分支会自动创建“跟踪分支”(它跟踪的分支叫做'上游分支')。跟踪分支是与远程分支有直接关系的本地分支。

&emsp;&emsp;当克隆一个仓库时，它通常会自动创建一个跟踪origin/master的master分支。如果你想设置其他跟踪分支，可以运行`git checkout -b [branch] [remotename]/[branch]`.这是一个十分常用的操作，所以Git提供了--track快捷方式：


        $ git checkout --track origin/serverfix
        Branch serverfix set up to track remote
        branch serverfix from origin.
        Switched to a new branch 'serverfix'


&emsp;&emsp;如果想要将本地分支与远程分支设置为不同名，可以轻松地使用上一个命令增加一个不同名字的本地分支：


        $ git checkout -b sf origin/serverfix
        Branch sf set up to track remote branch
        serverfix from origin
        Switch to a new branch 'sf'


&emsp;&emsp;如果你本地本来就有一个本地分支，你想要跟踪上游分支：


        $ git branch -u origin/serverfix
        Branch serverfix set up to track remote
        branch serverfix from origin


&emsp;&emsp;如果你想查看跟踪分支，可以使用`git branch -vv`.这会将所有的本地分支列出来并且包含更多的信息，例如每一个分支正在跟踪哪个远程分支与本地分支是否是领先、落后或是都有。


        $ git branch -vv
        iss53     7e424c3 [origin/iss53: ahead 2] forgot the brackets
        master    1ae2a45 [origin/master] deploying index fix
        * serverfix f8674d9 [teamone/server-fix-good: ahead 3, behind 1] this should do it
        testing   5ea463a trying something new


&emsp;&emsp;我们解析一下上面的含义。iss53分支正在跟踪origin/iss53并且ahead是2，意味着本地有两个提交还没有推送到服务器上。master分支正在跟踪origin/master分支并且进度是最新的。serverfix分支正在跟踪teamone服务器上的server-fix-good分支并且领先3落后1，意味着服务器上又一次提交还没有合并，同时本地有三次提交还没有推送。testing分支没有跟踪任何分支。

# 拉取

&emsp;&emsp;当`git fetch`命令从服务器上抓取本地没有的数据时，它并不会修改工作目录中的内容。它并不会修改工作目录中的内容，它只会获取数据然后让你自己合并。然而，`git pull`大多数情况下它是git fetch紧接着git merge命令。**由于git pull的魔法经常令人困惑所以通常单独显式地使用fetch与merge命令会更好一些**

# 删除远程分支

&emsp;&emsp;假设你和协作者已经完成了一个特性并且将其合并到了远程仓库的master分支。可以运行`--delete`选项的git push命令来删除一个远程分支。


        $ git push origin --delete serverfix
        To https://github.com/schacon/simplegit
            - [deleted]     serverfix
