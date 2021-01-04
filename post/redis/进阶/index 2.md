
## 事务

&emsp;&emsp;在微博中，用户之间是"关注"和"被关注"的关系。如果要使用Redis存储这样的关系可以使用集合类型。思路是对每个用户使用两个集合类型键，分别名为`user:用户ID:followers`和`user:用户ID:following`，用来存储关注该用户的用户集合该用户关注的用户集合。

    func follow(currentUser, targetUser string) {
        sadd user:currentUser:following targetUser
        sadd user:targetUser:followers currentUser
    }

&emsp;&emsp;如ID为1的用户A想关注ID为2的用户B，只需要执行follow(1, 2)即可。

&emsp;&emsp;然而在实现该功能的时候有一个问题：完成关注操作需要依次执行两条Redis命令，如果在第一条命令执行完后因为某种原因导致第二条命令没有执行，就会出现一个奇怪的现象：A查看自己关注的用户列表时会发现B，而B查看关注自己的用户列表时却没有A。

### 概述

&emsp;&emsp;Redis的事务是一组命令的集合。事务同命令一样都是Redis的最小执行单位，一个事务中的命令要么都执行，要么都不执行。事务的原理是先将属于一个事务的命令发送给Redis，然后再让Redis依次执行这些命令。

        multi
        sadd user:1:following 2
        queued
        sadd user:2followers 1
        queued
        exec
        (integer) 1
        (integer) 2

&emsp;&emsp;可以看到Redis没有立刻执行这些命令，而是返回`queued`表示这两条ml已经进入等待执行的事务队列中了。当把所有要在同一个事务中执行的命令都发给Redis后，我们使用`exec`命令告诉Redis将等待执行的事务队列中的所有命令按照发送顺序依次执行。`exec`命令返回值就是这些命令的返回值组成的列表，返回值顺序和命令的顺序相同。

&emsp;&emsp;除此之外，Redis事务还能保证一个事务内的命令依次执行而不被其它命令插入。

### 错误处理

&emsp;&emsp;如果一个事务中的某个命令执行出错，Redis会怎么处理？

1. **语法错误**。语法错误指命令不存在或者命令参数不对

        multi
        set key value
        queued
        set key
        (error) err wrong number of arguments for 'set' command
        errorcommand key
        (error) err unknown command 'errorcommand'
        exec
        (error) execabort transaction disarded because of previoud errors

&emsp;&emsp;跟在`multi`命令后执行了3个命令：一个是正确的，成功的加入事务队列；其余两个命令都有语法错误，只要有一个命令有语法错误，执行exec命令后Redis就会直接返回错误，连语法正确的命令也不会执行。

2. **运行错误**。运行错误指在命令执行时出现错误，比如使用散列类型的命令操作集合类型的键，这种错误在实际执行之前Redis是无法发现的，所以在事务里这样的命令是会被Redis接受并执行的。如果事务里的一条命令出现了运行错误，事务里其它的命令依然会继续执行。

        multi
        set key 1
        queued
        sadd key 2
        queued
        set key 3
        queued
        exec
        ok
        (error) err operation against a key holding the wrong kind of value
        ok
        get key
        "3"

&emsp;&emsp;可见，虽然`sadd key 2`出现了错误，但是`set key 3`依然执行了。Redis事务没有关系数据库事务提供的回滚功能，为此开发者必须在事务执行出错后自己收拾剩下的摊子。


### watch

&emsp;&emsp;我们已经知道在一个事务中只有当所有都依次执行完后才能饿到每个结果的返回值。可是有些情况下需要先获得一条命令的返回值，然后再根据这个值执行下一条命令。

&emsp;&emsp;肯定会有很多读者想到可以用事务来实现`incr`以防止竞态条件，可是因为事务中的每个命令的执行结果狗屎最后一起返回的，所以无法将前一条命令的结果作为下一条命令的参数。

&emsp;&emsp;为了解决这个问题，我们需要换一个思路。