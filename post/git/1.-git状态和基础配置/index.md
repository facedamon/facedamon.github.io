
# 三种状态

&emsp;&emsp;Git有三种状态：已提交(committed), 已修改(modified), 已暂存(staged)。已提交表示数据已经安全的保存在本地数据恐库中。已修改表示修改了文件，但还没保存到数据库中。已暂存表示对文件已修改的当前保本做了标记，使之包含在下次提交的快照中。

&emsp;&emsp;以此引入Git的三个工作区域：Git仓库，工作目录以及暂存区域。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/books/git/git工作区域.png)

## 基本的Git工作流程：

1. 在工作目录中修改文件
2. 暂存文件，将文件的快照放入暂存区域
3. 提交更新，找到暂存区域的文件，将快照永久性存储到Git仓库

# 配置

&emsp;&emsp;Git自带一个git config的工具来帮助设置Git外观和行为的配置变量，这些变量存储在三个不同的位置。

1. `/etc/gitconfig`: 包含系统上每一个用户及他们仓库的通用配置。如果使用带有 `--system`选项的git config时，它会从此文件读取配置变量。
2. `~/.gitconfig 或 ~/.config/git/config`: 只针对当前用户。使用带有`--global`选项的git config时，从此文件读取配置变量。
3. `.git/config`: 只针对所在目录仓库

&emsp;&emsp;每一个级别覆盖上一级别的配置，所以.git/config覆盖/etc/gitconfig

## 配置个人信息


$ git config --global user.name "John"
$ git config --global user.email john@example.com
$ git config --global core.editor vim


&emsp;&emsp;使用了`--global`选项，那么该命令控制个人目录下的git配置文件。如果你想针对特定目录下的仓库使用不同的用户名和邮件，可以在那个项目下运行没有--global选项的命令来配置。

## 检查配置信息

        $ git config --list
        user.name=John
        user.email=John@example.com
        color.status=auto
        color.branch=auto
        core.editor=vim


&emsp;&emsp;你应该看到了重复的变量名，因为Git会从不同的文件中读取同一个配置，这种情况下，Git会使用它找到的每一个变量的最后一个配置。

&emsp;&emsp;你也可以通过`git config <key>`来检查Git的某一项配置：


        $ git config user.name
        John


# 获取帮助


        $ git help <verb>
        $ git <verb> --help
        $ man git-<verb>


&emsp;&emsp;例如想获取config或log命令手册：


        $ git help config
        $ man git-log
