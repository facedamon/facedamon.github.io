
# 摘要

&emsp;&emsp;乐观锁和悲观锁是两种思想，用于解决并发场景下的数据竞争问题。

- 乐观锁：乐观锁在操作数据时非常乐观，认为别人不会同时修改数据。因此乐观锁不会上锁，只是在执行更新的时候判断一下在此期间别人是否修改了数据：如果别人修改了数据则放弃操作，否则执行操作。
- 悲观锁：悲观锁在操作数据时比较悲观，认为别人会同时修改数据。因此操作数据时直接把数据锁住，知道操作完成后才会释放锁；上锁期间其它人不能修改数据。

## 实现方式

&emsp;&emsp;在说明实现方式之前，需要明确：**乐观锁和悲观锁是两种思想，它们的使用是非常广泛的，不局限于某种编程语言或数据库。**

&emsp;&emsp;悲观锁的实现方式是加锁，加锁既可以是对代码块加锁(如Java的synchronized关键字)，也可以是对数据加锁(如mysql中的排他锁)。

&emsp;&emsp;乐观锁的实现方式主要有两种：CAS机制和版本号机制，下面详细介绍。

### CAS(Compare And Swap)

&emsp;&emsp;CAS操作包括了3个操作数：

- 需要读写的内存位置(V)
- 进行比较的预期值(A)
- 拟写入的新值(B)

&emsp;&emsp;CAS操作逻辑如下：如果内存位置V的值等于预期的A值，则将该位置更新为新值B，否则不进行任何操作。许多CAS的操作是自旋的：如果操作不成功，会一直重试，直到操作成功为止。

&emsp;&emsp;这里引出一个新的问题，既然CAS包含了Compare和Swap两个操作，它又如何保证原子性呢？答案是：CAS是由CPU支持的原子操作，其原子性是在硬件层面进行保证的。

&emsp;&emsp;下面以Java中的自增操作(i++)为例，看一下悲观锁和CAS分别是如何保证线程安全的。我们知道，在Java中自增操作不是原子操作，它实际上包含三个独立的操作：（1）读取i值；（2）加1；（3）将新值写回i

&emsp;&emsp;因此，如果并发执行自增操作，可能导致计算结果的不准确。在下面的代码示例中：value1没有进行任何线程安全方面的保护，value2使用了乐观锁(CAS)，value3使用了悲观锁(synchronized)。运行程序，使用1000个线程同时对value1、value2和value3进行自增操作，可以发现：value2和value3的值总是等于1000，而value1的值常常小于1000。

    package org.example;

    import java.util.concurrent.atomic.AtomicInteger;

    public class Sync {
        // value1线程不安全
        private static int value1 = 0;
        // value2乐观锁
        private static AtomicInteger value2 = new AtomicInteger(0);
        // value3悲观锁
        private static int value3 = 0;
        private static synchronized void increaseValue3() {
            value3++;
        }

        public static void main(String[] args) throws Exception{
            // 开启1000 个线程，并执行自增操作
            for (int i = 0; i < 1000; i++) {
                new Thread(new Runnable() {
                    @Override
                    public void run() {
                        try {
                            Thread.sleep(100);
                        }catch (InterruptedException e) {
                            e.printStackTrace();
                        }
                        value1++;
                        value2.getAndIncrement();
                        increaseValue3();
                    }
                }).start();
            }
            //打印结果
            Thread.sleep(1000);
            System.out.println("线程不安全：" + value1);
            System.out.println("乐观锁(AtomicInteger)：" + value2);
            System.out.println("悲观锁(synchronized)：" + value3);
        }
    }

&emsp;&emsp;首先来介绍AtomicInteger。AtomicInteger是java.util.concurrent.atomic包提供的原子类，利用CPU提供的CAS操作来保证原子性；除了AtomicInteger外，还有AtomicBoolean、AtomicLong、AtomicReference等众多原子类。

&emsp;&emsp;下面看一下AtomicInteger的源码，了解下它的自增操作getAndIncrement()是如何实现的（源码以Java7为例，Java8有所不同，但思想类似）。

    public class AtomicInteger extends Number implements java.io.Serializable {
        //存储整数值，volatile保证可视性
        private volatile int value;
        //Unsafe用于实现对底层资源的访问
        private static final Unsafe unsafe = Unsafe.getUnsafe();
    
        //valueOffset是value在内存中的偏移量
        private static final long valueOffset;
        //通过Unsafe获得valueOffset
        static {
            try {
                valueOffset = unsafe.objectFieldOffset(AtomicInteger.class.getDeclaredField("value"));
            } catch (Exception ex) { throw new Error(ex); }
        }
    
        public final boolean compareAndSet(int expect, int update) {
            return unsafe.compareAndSwapInt(this, valueOffset, expect, update);
        }
    
        public final int getAndIncrement() {
            for (;;) {
                int current = get();
                int next = current + 1;
                if (compareAndSet(current, next))
                    return current;
            }
        }
    }

&emsp;&emsp;源码分析说明：

1. getAndIncrement()实现的自增操作是自旋CAS操作：在循环中进行compareAndSet,如果执行成功则退出，否则一直执行。
2. 其中compareAndSet是CAS操作的核心，它是利用Unsafe对象实现的。
3. Unsafe又是何许人也？Unsafe是用来帮助Java访问操作系统底层资源的类(如可以分配内存、释放内存)，通过Unsafe，Java具有了底层操作能力，可以提升效率；强大的底层资源操作能力也带来了安全隐患(类的名字Unsafe也在提醒我们这一点)，因此正常情况下用户无法使用。AtomicInteger在这里使用了Unsafe提供的CAS功能。
4. valueOffset可以理解为value在内存中的偏移量，对应了CAS三个操作数(V/A/B)中的V；偏移量的获得也是通过Unsafe实现的。
5. value域的volatile修饰词：Java并发编程要保证线程安全，需要保证**原子性、可视性和有序性**；CAS操作可以保证原子性，而volatile可以保证可视性和一定程度的有序性；在AtomicInteger中，volatile和CAS一起保证了线程安全性。源于volatile作用原理的说明涉及到Java内存模型，这里不详细展开。

&emsp;&emsp;说完了AtomicInteger，再说synchronized。synchronized通过对代码块加锁来保证线程安全：在同一时刻，只能有一个线程可以执行代码块中的代码。synchronized是一个重量级的操作，不仅是因为加锁需要消耗额外的资源，还因为线程状态的切换会涉及操作系统核心态和用户态的转换；不过随着JVM对锁进行的一系列优化(如自旋锁、轻量级锁、锁粗化等)，synchronized的性能表现已经越来越好。

### 版本号机制

&emsp;&emsp;除了CAS，版本号机制也可以用来实现乐观锁。版本号机制的基本思路是在数据中增加一个字段version，表示该数据的版本号，每当数据被修改，版本号加1.当某个线程查询数据时，将该数据的版本号一起查出来；当该线程更新数据时，判断当前版本号与之前读取的版本号是否一致，如果一致才进行操作。

&emsp;&emsp;需要注意的是，这里使用了版本号作为判断数据变化的标记，实际上可以根据实际情况选用其它能够标记数据版本的字段，如时间戳等。

&emsp;&emsp;下面以“更新玩家金币数”为例（数据库为MySQL，其他数据库同理），看看悲观锁和版本号机制是如何应对并发问题的。

&emsp;&emsp;考虑这样一种场景：游戏系统需要更新玩家的金币数，更新后的金币数依赖于当前状态(如金币数、等级等)，因此更新前需要先查询玩家当前状态。

&emsp;&emsp;下面的实现方式，没有进行任何线程安全方面的保护。如果有其他线程在query和update之间更新了玩家的信息，会导致玩家金币数的不准确。

        @Transactional
        public void updateCoins(Integer playerId) {
            // 根据playerId查询玩家信息
            Player player = query("select coins, level from player where play_id = {0}", playerId);
            // 根据玩家当前信息及其它信息，计算新的金币数
            Long newCoins = _;
            // 更新金币数
            update("update player set coins = {0} where play_id = {1}", newCoins, playerId);
        }

&emsp;&emsp;为了避免这个问题，悲观锁通过加锁解决这个问题，代码如下所示。在查询玩家信息时，使用select …… for update进行查询；该查询语句会为该玩家数据加上**排它锁**，直到事务提交或回滚时才会释放排它锁；在此期间，如果其他线程试图更新该玩家信息或者执行select for update，会被阻塞。

        @Transactional
        public void updateCoins(Integer playerId) {
            // 根据playerId查询玩家信息(加排他锁)
            Player player = query("select coins, level from player where play_id = {0} for update", playerId);
            // 根据玩家当前信息及其它信息，计算新的金币数
            Long newCoins = _;
            // 更新金币数
            update("update player set coins = {0} where play_id = {1}", newCoins, playerId);
        }

&emsp;&emsp;版本号机制则是另一种思路，它为玩家信息增加一个字段：version。在初次查询玩家信息时，同时查询出version信息；在执行update操作时，校验version是否发生了变化，如果version变化，则不进行更新。

        @Transactional
        public void updateCoins(Integer playerId) {
            // 根据playerId查询玩家信息,包含version信息
            Player player = query("select coins, level, version from player where play_id = {0}", playerId);
            // 根据玩家当前信息及其它信息，计算新的金币数
            Long newCoins = _;
            // 更新金币数
            update("update player set coins = {0}, version = version + 1 where play_id = {1}", newCoins, playerId);
        }

# 优缺点和适用场景

&emsp;&emsp;乐观锁和悲观锁并没有优劣之分，它们有各自适合的场景；下面从两个方面进行说明。

## 功能限制

&emsp;&emsp;与悲观锁相比，乐观锁适用的场景受到了更多的限制，无论是CAS还是版本号机制。

&emsp;&emsp;例如，CAS只能保证单个变量操作的原子性，当涉及到多个变量时，CAS是无能为力的，而synchronized则可以通过对整个代码块加锁来处理。再比如版本号机制，如果query的时候是针对表1，而update的时候是针对表2，也很难通过简单的版本号来实现乐观锁。

## 竞争激烈程度

&emsp;&emsp;如果悲观锁和乐观锁都可以使用，那么选择就要考虑竞争的激烈程度：

- 当竞争不激烈 (出现并发冲突的概率小)时，乐观锁更有优势，因为悲观锁会锁住代码块或数据，其他线程无法同时访问，影响并发，而且加锁和释放锁都需要消耗额外的资源。
- 当竞争激烈(出现并发冲突的概率大)时，悲观锁更有优势，因为乐观锁在执行更新时频繁失败，需要不断重试，浪费CPU资源。

# 乐观锁加锁吗？

1. 乐观锁本身是不加锁的，只是在更新时判断一下数据是否被其它线程更新；AtomicInteger便是一个例子。
2. 有时乐观锁可能与加锁操作合作，例如，在前述updateCoins例子中，MySQL在执行updaye时会加排他锁。但这只是乐观锁与加锁合作的例子，不能改变“乐观锁本身不加锁”这一事实。

# CAS有那些缺点

- ABA问题

&emsp;&emsp;假设有两个线程，两个线程按照顺序进行以下操作：

1. 线程1读取内存中数据为A；
2. 线程2将该数据修改为B；
3. 线程2将该数据修改为A；
4. 线程1对数据进行CAS操作；

&emsp;&emsp;在第四步中，由于内存中数据仍然为A，因此CAS操作成功，但实际上该数据已经被线程2修改过了。这就是ABA问题。

&emsp;&emsp;在AtomicInteger的例子中，ABA似乎没有什么危害。但是在某些场景下，ABA却会带来隐患，例如栈顶问题：一个栈的栈顶经过两次(或多次)变化又恢复了原值，但是栈可能已经发生了变化。

&emsp;&emsp;对于ABA问题，比较有效的方案是引入版本号，内存中的值每发生一次变化，版本号都+1；在仅从CAS操作时，不仅比较内存中的值，也会比较版本号，只有当二者都没有变化时，CAS才能执行成功。Java中的AtomicStampedReference类便是使用版本号来解决ABA问题的。

- 高竞争下的开销问题

&emsp;&emsp;在并发冲突概率大的竞争环境下，如果CAS一直失败，会一直重试，CPU开销较大。针对这个问题的一个思路是引入退出机制，如重试次数超过一定阈值后失败退出。当然，更重要的是避免在高竞争环境下使用乐观锁。

- 功能限制

&emsp;&emsp;CAS的功能是比较受限的，例如CAS只能保证单个变量操作的原子性，这意味着：

1. 原子性不一定能保证线程安全，例如在Java中需要与volatile配合来保证线程安全
2. 当涉及到多个变量时，CAS也无能为力

&emsp;&emsp;除此之外，CAS的实现需要硬件层面处理器的支持，在Java中普通用户无法直接使用，只能借助atomic包下的原子类，灵活性受到限制。
