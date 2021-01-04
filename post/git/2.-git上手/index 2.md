
# 概要

&emsp;&emsp;希望通过这一节学习，能够**`快速`**掌握配置初始化仓库(repository)、开始或停止跟踪(track)文件、暂存(stage)、提交(commit)。Git忽略指定文件和文件模式，如何快速撤销错误操作，如何浏览你的项目历史版本及不同提交间的差异，如何向你的远程仓库推送(push)以及拉取(pull)文件。

## 创建仓库

- 本地导入


        $ git init
        $ git add LICENSE
        $ git commit -m 'initial project version'


&emsp;&emsp;init命令创建一个`.git`的子目录，包含你初始化的Git仓库中所有必须文件

- 远程克隆


        $ git clone https://github.com/libgit2/libgit2


&emsp;&emsp;Git克隆的是该Git仓库服务器的几乎所有数据，而不是仅仅复制完成你的工作所需的文件。默认拉取仓库中的每一个文件的每一个版本

## 状态变迁

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/文件状态变化周期.png)

&emsp;&emsp;初次克隆某个仓库，工作目录中的所有文件都属于以跟踪文件，并处于位修改状态。编辑某个文件，Git将它标记为已修改文件，依次将这些修改过的文件放入暂存区，提交所有暂存的文件，如此反复。

## 查看状态


        $ git status
        On branch master
        nothing to commit, working directory clean

        $ echo 'my project' > README
        $ git status 
        On branch master
        Untracked files:
            (use "git add <file>..." to include in what will be commited)
            README
            nothing added to cimmit but untracked
            files present (use "git add" to track)


&emsp;&emsp;可以看到新建的README文件出现在Untracked files下面。未跟踪的文件意味着Git在之前的快照中没有这些文件，你需要告诉它“我需要跟踪这些文件”。

## 跟踪文件


        $ git add README
        $ git status
        On branch master
        Changes to be commited:
            (use "git reset HEAD <file>..." to unstage)
            new file: README


&emsp;&emsp;只要在Changes to be commited 这行下面的，就说明是已暂存状态。`git add`命令使用文件或目录作为参数，如果参数是目录，该命令将递归跟踪该目录下的所有文件。

## 暂存已修改文件

&emsp;&emsp;修改已经被追踪的文件：CONTRIBUTING.md


        $ git status
        On branch master
        Changes to be committed:
            (use "git reset HEAD <file>..." to unstage)
            new file: README
        Changes not staged for commit:
            (use "git add <file>..." to update what will be committed)
            (use "git checkout -- <file>..." to discard丢弃 changes in working directory)
            modified: CONTRIBUTING.md


&emsp;&emsp;文件CONTRIBUTING.md出现在Changes not staged for commit下面，说明已跟踪的文件内容发生变化，但还没有放到暂存区。查看上图状态周期，git add命令：可以跟踪新文件，或者将已跟踪的文件放入暂存区。`可以理解该命令为添加内容到下一次提交中。`

## 状态简览

&emsp;&emsp;`git status`命令输出十分详细，你可以使用`git status -s`得到更加紧凑的格式输出。该命令工作中使用最广泛.


        $ git status -s
        M README
        MM Rakefile
        A  lib/git.rb
        M  lib/simplegit.rb
        ?? LICENSE.txt


        - ?? : 新添加的未追踪文件
        - A  : 新添加到暂存区中的文件
        - 右M : 该文件被修改了但是还没有放入暂存区
        - 左M : 该文件被修改了并放入了暂存区
        - MM : 该文件在工作区被修改并提交到暂存区后又在工作区被修改了，所以暂存区和工作区都有该文件被修改的记录 

## .gitignore


        # .a  注释
        # but do track lib.a, even though you are ignoring .a files above
        ## 不跟踪lib.a
        !lib.a

        # only ignore the TODO file in the current directory. not subdir/TODO
        # 忽略当前目录的TODO文件， 非子目录的TODO
        /TODO

        # ignore all files in the build/directory
        # 忽略build文件夹下的所有文件
        build/

        # ignore doc/notes.txt, but not doc/server/arch.txt
        # 忽略doc目录下的txt结尾的文件，非doc子目录下的txt文件
        doc/*.txt

        # ignore all .pdf files in the doc/ directory
        # 忽略doc下所有文件夹下的所有pdf文件
        doc/**/*.pdf


&emsp;&emsp;GitHub有一个十分详细的针对数种项目及语言的.gitignore文件列表，你可以[点击这里](https://github.com/github/gitignore)找到。

## 查看已暂存和未暂存的修改

- **当前做了哪些更新还没暂存？**
- **哪些已经暂存准备下次提交？**

&emsp;&emsp;查看尚未暂存文件更新了哪些部分，不加参数直接输入`git diff`，此命令比较的事工作目录种当前文件和暂存区快照之间的差异。


        $ git diff
        diff --git a/CONTRIBUTING.md
        b/CONTRIBUTING.md
        index 8ebb991..643e24d 100644
        --- a/CONTRIBUTING.md
        +++ b/CONTRIBUTING.md
        @@ -65,7 +65,8 @@ branch directly, things can get messy.


&emsp;&emsp;查看已暂存的将要添加到下次里的内容，可以使用`git diff --staged`


        $ git diff --staged
        diff --git a/README b/README
        new file mode 100644
        index 000000..03902a1
        --- /dev/null
        +++ b/README
        @@ -0,0 +1 @@
        +My Project


## 提交更新

&emsp;&emsp;暂存区准备好后，就可以提交了。提交之前用`git status -s`看一下是不是都暂存起来了。


        $ git commit -m 'Story 182: Fix benchmarks for speed'
        [master 463dc4f] Story 182: Fix benchmarks for speed
            2 files changed, 2 insertions(+)
            create mode 10064 README


&emsp;&emsp;提交后，显示当前在哪个分支(master)提交的，本次提交的完整SHA-1校验(463dc4f), 以及在本次提交种，有多少文件修订过，多少行添加和修改过。

> 跳过使用暂存区

&emsp;&emsp;Git提供了一个跳过使用暂存区的方式，只要在提交的时候，`git commit -am 'xxx'`.添加`-a`选项，Git就会自动把所有已跟踪过的文件暂存起来一并提交，从而跳过git add步骤。

## 删除文件

- **删除已暂存文件，一并删除工作区文件**


        $ git rm -f PROJECTS.md
        rm -r 'PROJECTS.md'


- **删除已暂存文件，保留工作区文件**


        $ git rm --cached log/\*.log


&emsp;&emsp;第二种情况非常常见，比如项目一开始时，忘记添加.gitignore文件，不小心把一大堆日志文件或编译文件添加到暂存区，这一做法尤其有用。

## 移动文件


        $ git mv README.md README
        $ git status
        On branch master
        Changes to be committed:
        (use "git reset HEAD <file>"... to unstage)

        renamed: README.md -> README


&emsp;&emsp;其实，运行`git mv`就相当于运行了下面三条命令


        $ mv README.md README
        $ git rm README.md
        $ git add README


## 查看提交历史

&emsp;&emsp;`git log`会按提交时间列出所有的更新，最近的更新排在最上面。这个命令会列出每个提交的SHA-1校验和、作者的名字、电子邮件地址、提交时间以及提交说明。

&emsp;&emsp;常用选项`-p`，用来显示每次提交的内容差异。你可以加上`-2`来仅仅显示最近两次提交。当进行代码审查，或者快速浏览某个搭档提交的commit所带来的变化时，这个参数就非常有用了。


        $ git log -p -2
        commit
        ca82a6dff817ec66f44342007202690a93763949
        AUthor: Scott Chacon <schacon@geemail.com>
        Date: Mon Mar 17 21:52:11 2008 -0700

            changed the version number

        diff --git a/Rakefile b/Rakefile
        index a874b73..8f94139 100644
        --- a/Rakefile
        +++ b/Rakefile
        @@ -5,7 +5,7 @@ require
        'rake/gempackagetask'
        “spec = Gem::Specification.new do |s|
            s.platform  =   Gem::Platform::RUBY
            s.name      =   "simplegit"
        -    s.version   =   "0.1.0"
        +    s.version   =   "0.1.1"
            s.author    =   "Scott Chacon"
            s.email     =   "schacon@geemail.com”
            s.summary   =   "A simple gem for using Git in Ruby code.”
        
        “commit 085bb3bcb608e1e8451d4b2432f8ecbe6306e7e7
        Author: Scott Chacon <schacon@gee-mail.com>
        Date:   Sat Mar 15 16:40:33 2008 -0700

            removed unnecessary test

        diff --git a/lib/simplegit.rb b/lib/simplegit.rb
        index a0a60ae..47c6340 100644
        --- a/lib/simplegit.rb
        +++ b/lib/simplegit.rb
        @@ -18,8 +18,3 @@ class SimpleGit
            end

        end
        -
        -if $0 == __FILE__
        -  git = SimpleGit.new”
        “-  puts git.show
        -end
        \ No newline at end of file”



&emsp;&emsp;你还可以使用总结性选项，如果你想看每次提交的简略统计信息，可以使用`--stat`选项。列出所有被修改过的文件、有多少文件被修改了以及被修改的文件哪些行被移除或是添加了，在每次提交的最后还有一个总结。


        $ git log --stat
        commit
        ca82a6dff817ec66f44342007202690a93763949
        Author: Scott Chacon <schacon@geemail.com>
        Date: Mon Mar 17 21:52:11 2008 -0700

            changed the version number

        Rakefile | 2 +-
        1 file changed, 1 insertion(+), 1 deletion(-)

        commit
        a11bef06a3f659402fe7563abf99ad00de2209e6
        Author: Scott Chacon <schacon@geemail.com>
        Date: Sat Mar 15 10:31:28 2008 -0700

            first commit

        README  | 6 ++++++
        Rakefile | 23
        +++++++++++++++++++++++
        lib/simplegit.rc | 25
        +++++++++++++++++++++++++
        3 files changed, 54 insertions(+)


&emsp;&emsp;常用选项`--pretty`。看英文名字就可以猜出它是指定格式展示提交历史的。这个选项有一些内建的自选项供你使用。`oneline`将每个提交放在一行显示，查看的提交数很大时非常有用。


        $ git log --pretty=oneline

        ca82a6dff817ec66f44342007202690a93763949 changed the version number
        085bb3bcb608e1e8451d4b2432f8ecbe6306e7e7 removed unnecessary test
        a11bef06a3f659402fe7563abf99ad00de2209e6 first commit


&emsp;&emsp;内建的子选项，最灵活的莫过于`format`,它可以定制要显示的记录格式


        $ git log --pretty=format:"%h - %an, %ar : %s"

        ca82a6d - Scott Chacon, 6 years ago:
        changed the version number
        085bb3b - Scott Chacon, 6 years ago:
        removed unnecessary test
        a11bef0 - Scott Chacon, 6 years ago:
        first commit


> git log \--pretty=format 常用占位符

|占位符|说明|
|:--|:--:|
|%H|提交对象的完整哈希|
|%h|提交对象的简短哈希|
|%T|树对象的完整哈希|
|%t|树对象的简短哈希|
|%P|父对象的完整哈希|
|%p|父对象的简短哈希|
|%an|作者author名字|
|%ae|作者email|
|%ad|作者修订日期|
|%ar|作者修订日期，按多久以前的方式显示|
|%cn|提交者committer名字|
|%ce|提交者email|
|%cd|提交日期|
|%cr|提交日期，按多久以前的方式显示|
|%s|提交说明|

&emsp;&emsp;作者和提交者的区别：当你为某个开源项目发布补丁，然后某个核心成员将你的补丁并入项目，你就是作者，核心成员就是提交者。

&emsp;&emsp;子选项`format`与 选项`--graph`结合，可以展示一些人性化的ASCII字符，这种输出类型在学习分支与合并以后会变得非常有趣。


        $ git log --pretty=format:"%h %s" --graph
        * 2d3acf9 ignore errors from SIGCHLD on trap
        *  5e3ee11 Merge branch 'master' of git://github.com/dustin/grit
        |\
        | * 420eac9 Added a method for getting the current branch.
        * | 30e367c timeout code and tests
        * | 5a09431 add timeout protection to grit
        * | e1193f8 support for heads with slashes in them
        |/
        * d6016bc require time for xmlschema
        *  11d191e Merge branch 'defunkt' into local


> git log 常用选项

|选项|说明|
|:--|:--:|
|-p|按补丁格式显示每个更新之间的差异|
|\--stat|显示每次更新的文件修改统计信息|
|\--shortstat|只显示--stat中最后的行数修改添加移除统计
|\--name-only|仅在提交信息下显示已修改的文件清单|
|\--name-status|显示新增、修改、删除的文件清单|
|\--abbrev-commit|仅显示SHA-1的前几个字符|
|\--relative-date|使用较短的相对时间显示|
|\--graph|显示ASCII图形表示分支合并历史|
|\--pretty|使用其它格式显示历史提交信息。可用的子选项包括oneline，short，format|

## 撤销操作

&emsp;&emsp;提交后发现忘记暂存某些需要修改的文件，可以使用`--amend`选项，最终你只会有一个提交，第二次提交将替代第一次提交的结果。


        $ git commit -m 'initial commit'
        $ git add forgotten_file
        $ git commit -amend


&emsp;&emsp;如果你意外的执行`git add *`暂存了不需要的文件，想要撤销，通常`git status`会给你一些提示`git reset HEAD <file>...`来取消暂存。


        $ git reset HEAD CONTRIBUTING.md
        Unstaged changes after reset:
        M CONTRIBUTING.md


&emsp;&emsp;假如你并不想保留对CONTRIBUTING.md文件的修改，想将它还原成上次提交时的样子，通常`git status`也给了你一些提示`git checkout -- <file>...`来取消暂存。


        $ git checkout -- CONTRIBUTING.md


&emsp;&emsp;如果你仍然想保留对这个文件的修改，但是现在仍然需要撤销操作，那么[Git分支]通常是更好的做法。

&emsp;&emsp;你需要记住，在Git中任何已提交的东西几乎总是可以恢复的。甚至哪些被删除的分支中的提交或使用`--amend`选项的提交也可以恢复。但是，任何你为提交的东西丢失后很可能再也找不到了。


## 远程仓库

&emsp;&emsp;远程仓库是指托管在互联网中的你的项目的版本库。你可以有好几个远程仓库，通常有些仓库对你只读，有些则可以读写。与他人协作涉及管理远程仓库以及根据需要推送和拉取数据。

### 查看远程仓库

&emsp;&emsp;如果你的远程仓库不止一个，`git remote -v`会将它们全部列出


$ git remote -v


### 添加远程仓库

&emsp;&emsp;运行`git remote add <shortname><url>`添加一个新的远程Git仓库，同时指定一个你可以轻松引用的简写。


        $ git remote
        origin
        $ git remote add pb https://github.com/paulboone/ticgit


### 拉取与抓取

&emsp;&emsp;`git fetch [remote-name]`该命令会访问远程仓库，从中拉取所有你还没有的数据。执行完后，你将会拥有那个远程仓库中所有分支的引用，它并不会自动合并，你需要手动合并。

&emsp;&emsp;如上，你可以使用字符串pb来代替整个URL。


        git fetch pb
        remote: Counting objects: 43, done
        remote: Compressing objects: 100% (36/36), done.
        remote: Total 43 (delta 10), reused 31 (delta 5)
        Unpacking objects: 100% (43/43), done.
        From https://github.com/paulboone/ticgit
        * [new branch] master -> pb/master
        * [new branch] ticgit -> pb/ticgit


&emsp;&emsp;如果你有一个分支跟踪一个远程分支，可以使用`git pull`命令来自动抓取然后合并远程分支到当前分支。

### 推送

&emsp;&emsp;当你想将maste分支推送到origin服务器时。`git push [remote-name] [branch-name]`. 当你和其他人协作时，他们先推送到上游你再推送到上游，你的推送就会被拒绝。你必须先将他们的工作拉取下来并将其合并后才能推送。


        $ git push origin master


### 查看远程仓库

&emsp;&emsp;`git remote show [remote-name]`。它会列出远程仓库的URL与跟踪分支的信息。它会告诉你正处于master分支，并且运行git pull，就会抓取所有的远程引用，然后将远程master分支合并到本地master分支。


        $ git remote show origin
        * remote origin
        Fetch URL:
        https://github.com/schacon/ticgit
        Push URL:
        https://github.com/schacon/ticgit
        HEAD branch: master
        Remote branches:
            master
        tracked
            dev-branch
        tracked
            Local branch configured for 'git pull':
                master merges with remote master
            Local ref configured for 'git push':
                master pushes to master (up to date)


### 远程仓库移除与重命名

&emsp;&emsp;如果你想要重命名引用的名字可以运行`git remote rename`修改一个远程仓库的简写。例如，你想将pb重命名为paul,可以如下这么做。注意，这也会将远程的pb/master修改为paul/master


        $ git remote rename pb paul
        $ git remote
        origin
        paul


&emsp;&emsp;如果你想移除一个远程仓库，不再使用某个特定镜像了：


        $ git remote rm paul
        $ git remote
        origin


## 打标签

&emsp;&emsp;所谓标签，即是将某个分支标识为具有代表性，最常见的就是版本序列号。

### 列出标签

&emsp;&emsp;列出符合某些模式的标签


        $ git tag -l 'v1.8.5*'
        v1.8.5
        v1.8.5-rc0
        v1.8.5-rc1
        v1.8.5-rc2
        v1.8.5-rc3


### 创建标签

&emsp;&emsp;附注标签是存储在Git数据库中的一个完整对象。它们是可以被校验的；其中包含打标签的名字、email、日期时间、还有标签信息。可以使用GNU Privacy Guard 签名与验证.


        $ git tag -a v1.4 -m 'my version 1.4'
        $ git tag
        v0.1
        v1.3
        v1.4


&emsp;&emsp;可以使用`git show`命令可以看到标签信息与对应的提交信息.


        $ git show v1.4
        tag v1.4
        Tagger: Ben Straub <ben@straub.cc>
        Date: Sat May 3 20:19:12 2014 -0700

        my version 1.4

        commit
        ca82a6dff817ec66f44342007202690a93763949
        Author: Scott Chacon <schacon@geemail.com>
        Date: Mon Mar 17 21:52:11 2008 -0700

            changed the version number


### 后期打标签

&emsp;&emsp;如果你已经提交了多次版本，想给之前的某次提交打标签，你需要找到该版本的hash值或部分hash。


        $ git log --pretty=oneline
        ...
        9fceb02d0ae598e95dc970b74767f19372d61af8 updated rakefile
        ...

        $ git tag -a v1.2 9fceb02


### push标签

&emsp;&emsp;你需要显式地推送标签到服务器。`git push origin [tagname]`


        $ git push origin master v1.5
        Couting objects: 14, done.
        Dalta comoression using up to 8 threads.
        Compressing objects: 100% (12/12). done.
        Writing objects: 100% (14/14), 2.05Kin | 0 bytes/s, done.
        Total 14 (delta 3), reused 0 (delta 0)
        To git@github.com:schacon/simplegit.git
        * [new tag] v1.5 -> v1.5


&emsp;&emsp;如果想要一次性推送多个标签，可以使用`--tags`选项。它将会把所有不在远程仓库的标签全部推送上去。


        $ git push origin --tags


### 删除标签

&emsp;&emsp;删除本地标签，删除远程仓库标签`git push <remote> :refs/tags/<tagname>`


        $ git tag -d v1.4
        Deleted tag 'v1.4'
        $ git push origin :refs/tags/v1.4
        To /git@github.com:schacon/simplegit.git
        - [deleted] v1.4


## Git别名

&emsp;&emsp;通过`git config`文件可以为每个命令设置一个别名。


        $ git config --global alias.co checkout
        $ git config --global alias.br branch
        $ git config --global alias.ci commit
        $ git config --global alias.st status
        $ git config --global alias.unstage 'reset HEAD --'


        $ git unstage filA
        // 等效
        $ git reset HEAD -- fileA


> 常用last, 可以查看最近一次提交的信息


        $ git config --global alias.last 'log -1 HEAD'
