
> 转载自 https://www.cnblogs.com/xuwc/p/9123078.html

# 问题描述

&emsp;&emsp;某天A君突然发现自己的接口请求量突然涨到之前的10倍，没多久该接口几乎不可使用，并引发连锁反应导致整个系统崩溃。如何应对这种情况呢？生活给了我们答案：比如老式电闸都安装了保险丝，一旦有人使用超大功率的设备，保险丝就会烧断以保护各个电器不被强电流给烧坏。同理我们的接口也需要安装上“保险丝”，以防止非预期的请求对系统压力过大而引起的系统瘫痪，当流量过大时，可以采取拒绝或者引流等机制。

&emsp;&emsp;在大数据量高并发访问时，经常会出现服务或接口面对暴涨的请求而不可用的情况，甚至引发连锁反映导致整个系统崩溃。此时你需要使用的技术手段之一就是限流，当请求达到一定的并发数或速率，就进行等待、排队、降级、拒绝服务等。在限流时，常见的两种算法是漏桶和令牌桶算法算法，本文即对相关内容进行重点介绍。

# 常用的限流算法

&emsp;&emsp;常用的限流算法有两种：漏桶算法和令牌桶算法。

&emsp;&emsp;漏桶算法思路很简单，水（请求）先进入到漏桶里，漏桶以一定的速度出水，当水流入速度过大会直接溢出，可以看出漏桶算法能强行限制数据的传输速率。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/fenbu/081225378155003.png)

&emsp;&emsp;漏桶算法（Leaky Bucket）：主要目的是控制数据注入到网络的速率，平滑网络上的突发流量。漏桶算法提供了一种机制，通过它，突发流量可以被整形以便为网络提供一个稳定的流量。请求先进入到漏桶里，漏桶以一定的速度出水，当水请求过大会直接溢出，可以看出漏桶算法能强行限制数据的传输速率

&emsp;&emsp;对于很多应用场景来说，除了要求能够限制数据的平均传输速率外，还要求允许某种程度的突发传输。这时候漏桶算法可能就不合适了，令牌桶算法更为适合。如图2所示，令牌桶算法的原理是系统会以一个恒定的速度往桶里放入令牌，而如果请求需要被处理，则需要先从桶里获取一个令牌，当桶里没有令牌可取时，则拒绝服务。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/fenbu/081226107372877.png)

&emsp;&emsp;令牌桶算法（Token Bucket）：是网络流量整形（Traffic Shaping）和速率限制（Rate Limiting）中最常使用的一种算法。典型情况下，令牌桶算法用来控制发送到网络上的数据的数目，并允许突发数据的发送。

&emsp;&emsp;大小固定的令牌桶可自行以恒定的速率源源不断地产生令牌。如果令牌不被消耗，或者被消耗的速度小于产生的速度，令牌就会不断地增多，直到把桶填满。后面再产生的令牌就会从桶中溢出。最后桶中可以保存的最大令牌数永远不会超过桶的大小

> 两种算法的区别

&emsp;&emsp;两者主要区别在于“漏桶算法”能够强行限制数据的传输速率，而“令牌桶算法”在能够限制数据的平均传输速率外，还允许某种程度的突发传输。在“令牌桶算法”中，只要令牌桶中存在令牌，那么就允许突发地传输数据直到达到用户配置的门限，所以它适合于具有突发特性的流量。

# 限流工具类RateLimiter

&emsp;&emsp;Google开源工具包Guava提供了限流工具类RateLimiter，该类基于令牌桶算法来完成限流，非常易于使用。RateLimiter类的接口描述请参考：[RateLimiter接口描述](http://docs.guava-libraries.googlecode.com/git-history/master/javadoc/com/google/common/util/concurrent/RateLimiter.html)。

&emsp;&emsp;RateLimiter 从概念上来讲，速率限制器会在可配置的速率下分配许可证，如果必要的话，每个acquire() 会阻塞当前线程直到许可证可用后获取该许可证，一旦获取到许可证，不需要再释放许可证。通俗的讲RateLimiter会按照一定的频率往桶里扔令牌，线程拿到令牌才能执行，比如你希望自己的应用程序QPS不要超过1000，那么RateLimiter设置1000的速率后，就会每秒往桶里扔1000个令牌。例如我们需要处理一个任务列表，但我们不希望每秒的任务提交超过两个，此时可以采用如下方式：

        final RateLimiter rateLimiter = RateLimiter.create(2.0)
        void submitTasks(List tasks, Executor executor) {
            for (Runnable task: tasks) {
                // 阻塞
                rateLimiter.acquire();
                executor.execute(task);
            }
        }

&emsp;&emsp;有一点很重要，那就是请求的许可数从来不会影响到请求本身的限制（调用acquire(1) 和调用acquire(1000) 将得到相同的限制效果，如果存在这样的调用的话），但会影响下一次请求的限制，也就是说，如果一个高开销的任务抵达一个空闲的RateLimiter，它会被马上许可，但是下一个请求会经历额外的限制，从而来偿付高开销任务。注意：RateLimiter 并不提供公平性的保证。

&emsp;&emsp;下面是主要源码

    public double acquire() {
            return acquire(1);
        }

    public double acquire(int permits) {
            checkPermits(permits);  //检查参数是否合法（是否大于0）
            long microsToWait;
            synchronized (mutex) { //应对并发情况需要同步
                microsToWait = reserveNextTicket(permits, readSafeMicros()); //获得需要等待的时间 
            }
            ticker.sleepMicrosUninterruptibly(microsToWait); //等待，当未达到限制时，microsToWait为0
            return 1.0 * microsToWait / TimeUnit.SECONDS.toMicros(1L);
        }

    private long reserveNextTicket(double requiredPermits, long nowMicros) {
            resync(nowMicros); //补充令牌
            long microsToNextFreeTicket = nextFreeTicketMicros - nowMicros;
            double storedPermitsToSpend = Math.min(requiredPermits, this.storedPermits); //获取这次请求消耗的令牌数目
            double freshPermits = requiredPermits - storedPermitsToSpend;

            long waitMicros = storedPermitsToWaitTime(this.storedPermits, storedPermitsToSpend)
                    + (long) (freshPermits * stableIntervalMicros); 

            this.nextFreeTicketMicros = nextFreeTicketMicros + waitMicros;
            this.storedPermits -= storedPermitsToSpend; // 减去消耗的令牌
            return microsToNextFreeTicket;
        }

    private void resync(long nowMicros) {
        // if nextFreeTicket is in the past, resync to now
        if (nowMicros > nextFreeTicketMicros) {
            storedPermits = Math.min(maxPermits,
                    storedPermits + (nowMicros - nextFreeTicketMicros) / stableIntervalMicros);
            nextFreeTicketMicros = nowMicros;
        }
    }

# 适用场景

&emsp;&emsp;并不能说明令牌桶一定比漏桶好，她们使用场景不一样。令牌桶可以用来保护自己，主要用来对调用者频率进行限流，为的是让自己不被打垮。所以如果自己本身有处理能力的时候，如果流量突发（实际消费能力强于配置的流量限制），那么实际处理速率可以超过配置的限制。而漏桶算法，这是用来保护他人，也就是保护他所调用的系统。主要场景是，当调用的第三方系统本身没有保护机制，或者有流量限制的时候，我们的调用速度不能超过他的限制，由于我们不能更改第三方系统，所以只有在主调方控制。这个时候，即使流量突发，也必须舍弃。因为消费能力是第三方决定的。

&emsp;&emsp;**总结起来：如果要让自己的系统不被打垮，用令牌桶。如果保证别人的系统不被打垮，用漏桶算法**

# 服务治理--限流(令牌桶算法)

&emsp;&emsp;最近在写一个分布式服务的框架，对于分布式服务的框架来说，除了远程调用，还要进行服务的治理。当进行促销的时候，所有的资源都用来完成重要的业务，就比如双11的时候，主要的业务就是让用户查询商品，以及购买支付，此时，金币查询、积分查询等业务就是次要的，因此要对这些服务进行服务的降级，典型的服务降级算法是采用令牌桶算法，因此在写框架的时候去研究了一下令牌桶算法

&emsp;&emsp;在实施QOS策略时，可以将用户的数据限制在特定的带宽，当用户的流量超过额定带宽时，超过的带宽将采取其它方式来处理。要衡量流量是否超过额定的带宽，网络设备并不是采用单纯的数字加减法来决定的，也就是说，比如带宽为100K，而用户发来的流量为110K，网络设备并不是靠110K减去100K等于10K，就认为用户超过流量10K。网络设备衡量流量是否超过额定带宽，需要使用令牌桶算法来计算。下面详细介绍令牌桶算法机制。

&emsp;&emsp;当网络设备衡量流量是否超过额定带宽时，需要查看令牌桶，而令牌桶中会放置一定数量的令牌，一个令牌允许接口发送或接收1bit数据（有时是1 Byte数据），当接口通过1bit数据后，同时也要从桶中移除一个令牌。当桶里没有令牌的时候，任何流量都被视为超过额定带宽,只有当桶中有令牌时，数据才可以通过接口。令牌桶中的令牌不仅仅可以被移除，同样也可以往里添加，所以为了保证接口随时有数据通过，就必须不停地往桶里加令牌，由此可见，往桶里加令牌的速度，就决定了数据通过接口的速度。

&emsp;&emsp;因此，我们通过控制往令牌桶里加令牌的速度从而控制用户流量的带宽。而设置的这个用户传输数据的速率被称为承诺信息速率（CIR），通常以秒为单位。比如我们设置用户的带宽为1000  bit每秒，只要保证每秒钟往桶里添加1000个令牌即可。

> 举例：

&emsp;&emsp;将CIR设置为8000  bit/s，那么就必须每秒将8000个令牌放入桶中，当接口有数据通过时，就从桶中移除相应的令牌，每通过1  bit，就从桶中移除1个令牌。当桶里没有令牌的时候，任何流量都被视为超出额定带宽，而超出的流量就要采取额外动作。

&emsp;&emsp;每秒钟往桶里加的令牌就决定了用户流量的速率，这个速率就是CIR，但是每秒钟需要往桶里加的令牌总数，并不是一次性加完的，一次性加进的令牌数量被称为Burst  size（Bc），如果Bc只是CIR的一半，那么很明显每秒钟就需要往桶里加两次令牌，每次加的数量总是Bc的数量。还有就是加令牌的时间，Time interval（Tc），Tc表示多久该往桶里加一次令牌，而这个时间并不能手工设置，因为这个时间可以靠CIR和Bc的关系计算得到，  Bc/ CIR= Tc。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/fenbu/682841-20161101192914955-808105864.png)

- a. 按特定的速率向令牌桶投放令牌

- b. 根据预设的匹配规则先对报文进行分类，不符合匹配规则的报文不需要经过令牌桶的处理，直接发送；

- c. 符合匹配规则的报文，则需要令牌桶进行处理。当桶中有足够的令牌则报文可以被继续发送下去，同时令牌桶中的令牌 量按报文的长度做相应的减少；

- d. 当令牌桶中的令牌不足时，报文将不能被发送，只有等到桶中生成了新的令牌，报文才可以发送。这就可以限制报文的流量只能是小于等于令牌生成的速度，达到限制流量的目的。


        package org.example;

        import java.io.*;
        import java.util.Random;
        import java.util.concurrent.*;
        import java.util.concurrent.locks.ReentrantLock;

        /**
        * Hello world!
        *
        */
        public class TokenBucket {
            // 默认桶容量 即最大瞬时流量是64MB
            // 1024B = 1KB
            // 1MB = 1024KB = 1KB * 1024 = 1024B * 1024
            // 64MB = 64 * 1MB
            private static final int DEFAULT_BUCKET_SIZE = 1024 * 1024 * 64;

            // 每个桶数据的单位是1字节 1Byte
            private int everyTokenSize = 1;

            // 瞬时最大流量
            private int maxFlowRate;

            // 平均流量
            private int avgFlowRate;

            // 使用队列来装载桶数据
            // 最大的流量峰值就是 = everyTokenSize * DEFAULT_BUCKET_SIZE
            private ArrayBlockingQueue<Byte> tokenQueue =
                    new ArrayBlockingQueue<>(DEFAULT_BUCKET_SIZE);

            private ScheduledExecutorService scheduledExecutorService =
                    Executors.newSingleThreadScheduledExecutor();
            private volatile boolean isStart = false;
            private ReentrantLock lock = new ReentrantLock(true);
            private static final byte A_CHAR = 'a';

            public TokenBucket(){}

            public TokenBucket(int maxFlowRate, int avgFlowRate) {
                this.maxFlowRate = maxFlowRate;
                this.avgFlowRate = avgFlowRate;
            }

            public TokenBucket(int everyTokenSize, int maxFlowRate, int avgFlowRate) {
                this.everyTokenSize = everyTokenSize;
                this.maxFlowRate = maxFlowRate;
                this.avgFlowRate = avgFlowRate;
            }

            public void addToken(Integer tokenNum) {
                // 若队列(桶)已经满了，就不再入新的令牌
                for (int i = 0; i < tokenNum; i++) {
                    tokenQueue.offer(Byte.valueOf(A_CHAR));
                }
            }

            // 获取令牌个数
            public boolean getTokens(byte[] dataSize) {
                // 传输率 = 传输内容大小对应的桶个数
                int needTokenNum = dataSize.length / everyTokenSize + 1;
                final ReentrantLock lock = this.lock;
                lock.lock();
                try {
                    // 是否存在足够的桶
                    boolean result = needTokenNum <= tokenQueue.size();
                    if (!result) {
                        return false;
                    }
                    int tokenCount = 0;
                    for (int i = 0; i < needTokenNum; i++) {
                        Byte poll = tokenQueue.poll();
                        if (poll != null) {
                            tokenCount++;
                        }
                    }
                    return tokenCount == needTokenNum;
                }finally {
                    lock.unlock();
                }
            }

            public TokenBucket build() {
                start();
                return this;
            }

            public static TokenBucket newBuilder() {
                return new TokenBucket();
            }

            public TokenBucket everyTokenSize(int everyTokenSize) {
                this.everyTokenSize = everyTokenSize;
                return this;
            }

            public TokenBucket maxFlowRate(int maxFlowRate) {
                this.maxFlowRate = maxFlowRate;
                return this;
            }

            public TokenBucket avgFlowRate(int avgFlowRate) {
                this.avgFlowRate = avgFlowRate;
                return this;
            }

            private String stringCopy(String data, int copyNum) {
                StringBuilder sbuilder = new StringBuilder(data.length() * copyNum);
                for (int i = 0; i < copyNum; i++) {
                    sbuilder.append(data);
                }
                return sbuilder.toString();
            }

            public void start() {
                // 初始化桶队列大小
                if (maxFlowRate != 0) {
                    tokenQueue = new ArrayBlockingQueue<Byte>(maxFlowRate);
                }
                // 初始化令牌生产者
                TokenProducer tokenProducer = new TokenProducer(avgFlowRate, this);
                scheduledExecutorService.scheduleAtFixedRate(tokenProducer, 0, 1, TimeUnit.SECONDS);
                isStart = true;
            }

            public void stop() {
                isStart = false;
                scheduledExecutorService.shutdown();
            }

            public boolean isStart() {
                return isStart;
            }

            // 令牌生产者线程
            class TokenProducer implements Runnable {
                private int avgFlowRate;
                private TokenBucket tokenBucket;

                public TokenProducer(int avgFlowRate, TokenBucket tokenBucket) {
                    this.avgFlowRate = avgFlowRate;
                    this.tokenBucket = tokenBucket;
                }

                @Override
                public void run() {
                    tokenBucket.addToken(avgFlowRate);
                }
            }

            private static void tokenTest() throws IOException, InterruptedException {
                TokenBucket tokenBucket = TokenBucket.newBuilder().avgFlowRate(512).maxFlowRate(1024).build();
                BufferedWriter bufferedWriter = new BufferedWriter(
                        new OutputStreamWriter(new FileOutputStream("/Users/facedamon/Desktop/t.txt"))
                );
                // 四个byte
                String data = "xxxx";
                // 模拟1000次请求
                // 每次请求获取i1字节个数据
                for (int i = 1; i <= 1000; i++) {
                    Random random = new Random();
                    int i1 = random.nextInt(100);
                    boolean tokens = tokenBucket.getTokens(tokenBucket.stringCopy(data, i1).getBytes());
                    TimeUnit.MILLISECONDS.sleep(100);
                    if (tokens) {
                        bufferedWriter.write("token pass --- index:" + i1);
                    } else {
                        bufferedWriter.write("token rejuect --- index:" + i1);
                    }
                    bufferedWriter.newLine();
                    bufferedWriter.flush();
                }
                bufferedWriter.close();
                tokenBucket.stop();
            }

            public static void main(String[] args) throws IOException, InterruptedException {
                tokenTest();
            }
        }

# Semaphore进行并发流控

&emsp;&emsp;Java 并发库的Semaphore 可以很轻松完成信号量控制，Semaphore可以控制某个资源可被同时访问的个数，通过 acquire() 获取一个许可，如果没有就等待，而 release() 释放一个许可。单个信号量的Semaphore对象可以实现互斥锁的功能，并且可以是由一个线程获得了“锁”，再由另一个线程释放“锁”，这可应用于死锁恢复的一些场合。下面的Demo中申明了一个只有5个许可的Semaphore，而有20个线程要访问这个资源，通过acquire()和release()获取和释放访问许可：

        public class SemaphoreTest {
            public static void main(String[] args) throws Exception {
                ExecutorService executor = Executors.newCachedThreadPool();
                final Semaphore semaphore = new Semaphore(5);
                for (int i = 0; i < 20; i++) {
                    final int no = i;
                    Runnable run = new Runnable() {
                        @Override
                        public void run() {
                            try{
                                semaphore.acquire();
                                sout("Accessing: " + no);
                                Thread.sleep((long)(Math.random() * 10000));
                                semaphore.release();
                                sout(semaphore.availablePermits());
                            }catch(InterruptedException e) {
                                e.printStackTrace();
                            }
                        }
                    };
                    executor.execute(run);
                }
                executor.shutdown();
            }
        }