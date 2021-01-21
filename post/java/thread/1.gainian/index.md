
# 摘要

&emsp;&emsp;多线程是Java中不可避免的一个重要主体。从本章开始，我们将展开对多线程的学习。接下来的内容，是对“JDK中新增JUC包”之前的Java多线程内容的讲解，涉及到的内容包括，Object类中的wait(), notify()等接口；Thread类中的接口；synchronized关键字。

注：JUC包是指，Java.util.concurrent包，它是由Java大师Doug Lea完成并在JDK1.5版本添加到Java中的。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/juc/18152411-a974ea82ebc04e72bd874c3921f8bfec.jpg)

&emsp;&emsp;说明：线程次共包括以下5种状态

1. 新建状态(New): 线程对象被创建后，就进入了新建状态。例如，Thread thread = new Thread().
2. 就绪状态(Runnable): 也被成为"可执行状态"。线程对象被创建后，其它线程调用了该对象的start()，从而来启动该线程。例如，thread.start().处于就绪状态的线程，随时可能被CPU调度执行
3. 运行状态(Running): 线程获取CPU权限进行执行。需要注意的是，线程只能从就绪状态进入到运行状态。
4. 阻塞状态(Blocked): 阻塞状态是线程因为某种原因放弃CPU使用权，暂停运行。直到线程进入就绪状态，才有机会转到运行状态。
   1. 等待阻塞--通过调用线程的wait方法，让线程等待某个工作的完成
   2. 同步阻塞--线程在获取synchronized同步锁失败，它会进入同步阻塞状态
   3. 其它阻塞--通过调用线程的sleep或join或发出了I/O请求时，线程会进入到阻塞状态。当sleep状态。当sleep状态超时、join等待线程终止或者超时、或者I/O处理完毕，线程重新转入就绪状态
5. 死亡状态(Dead): 线程执行完了或者因异常退出了run方法，该线程结束声明周期

&emsp;&emsp;这5种状态涉及到的内容包括Object类，Thread和synchronized关键字。这些内容我们会逐个学习。

- Object，定义了wait(), notify(), notifyAll()等休眠/唤醒函数
- Thread，定义了一些列的线程操作函数。例如，sleep()休眠函数，interrupt()中断函数，getName()获取线程名称等
- synchronized，是关键字，它区分为synchronized代码块和synchronized方法。synchronized的作用是让线程获取对象的同步锁

&emsp;&emsp;在后面详细介绍wait, notify等方法时，我们会分析为什么"wait, notify等方法要定义在Object类，而不是Thread类中"

## 常用的实现多线程的两种方式

&emsp;&emsp;之所以说是常用的，是因为还可以通过java.util.concurrent包中的线程池来实现多线程。关于线程池的内容，我们以后会详细介绍。现在，先对Thread和Runnable进行了解。

### Thread和Runnable简介

&emsp;&emsp;`Runnable`是一个接口，该接口中只包含了一个run()方法。

    public interface Runnable {
        public abstract void run();
    }

&emsp;&emsp;Runnable的作用，实现多线程。我们可以定义一个类A实现Runnable接口；然后，通过new Thread(new A())等方式新建线程.

&emsp;&emsp;`Thread`是一个类。Thread本身就实现了Runnable接口。它的声明如下：

    public class Thread implements Runnable {}

### Thread和Runnable的异同点

- 相同点：都是多线程的实现方式
- 不同点：
  - Thread是类，Runnable是接口
  - Thread本事实现了Runnable接口。我们直到“一个类只能有一个父类，但是却能实现多个接口”，因此Runnable具有更好的扩展性。
  - Runnable还可以用于“资源的共享”。即，多个线程都是基于某一个Runnable对象建立的，它们会共享Runnable对象上的资源

> 通常，建议痛殴Runnable实现多线程

### 示例

    public class MyThread extends Thread{
        private int ticket = 10;

        @Override
        public void run() {
            for (int i = 0; i < 20; i++) {
                if (this.ticket > 0) {
                    System.out.println(this.getName()+" 卖票：ticket"+this.ticket--);
                }
            }
        }

        static class Test {
            public static void main(String[] args) {
                MyThread t1 = new MyThread();
                MyThread t2 = new MyThread();
                MyThread t3 = new MyThread();
                t1.start();
                t2.start();
                t3.start();
            }
        }
    }

</br>

    Thread-2 卖票：ticket10
    Thread-2 卖票：ticket9
    Thread-0 卖票：ticket10
    Thread-1 卖票：ticket10
    Thread-0 卖票：ticket9
    Thread-2 卖票：ticket8
    Thread-0 卖票：ticket8
    Thread-1 卖票：ticket9
    Thread-0 卖票：ticket7
    Thread-2 卖票：ticket7
    Thread-0 卖票：ticket6
    Thread-1 卖票：ticket8
    Thread-0 卖票：ticket5
    Thread-2 卖票：ticket6
    Thread-0 卖票：ticket4
    Thread-1 卖票：ticket7
    Thread-0 卖票：ticket3
    Thread-2 卖票：ticket5
    Thread-0 卖票：ticket2
    Thread-0 卖票：ticket1
    Thread-1 卖票：ticket6
    Thread-2 卖票：ticket4
    Thread-1 卖票：ticket5
    Thread-1 卖票：ticket4
    Thread-1 卖票：ticket3
    Thread-2 卖票：ticket3
    Thread-1 卖票：ticket2
    Thread-2 卖票：ticket2
    Thread-2 卖票：ticket1
    Thread-1 卖票：ticket1

- MyThread继承于Thread，它是自定义个线程。每个MyThread都会卖出10张票。
- 主线程main创建并启动3个MyThread子线程。每个子线程都各自卖出了10张票。

    public class MyThread implements Runnable {
        private int ticket = 10;

        @Override
        public void run() {
            for (int i = 0; i < 20; i++) {
                if (this.ticket > 0) {
                    System.out.println(Thread.currentThread().getName() + " 卖票：ticket" + this.ticket--);
                }
            }
        }

        static class Test {
            public static void main(String[] args) {
                MyThread mt = new MyThread();
                // 启动3个线程t1,t2,t3(它们共用一个Runnable对象)，这3个线程一共卖10张票！
                Thread t1 = new Thread(mt);
                Thread t2 = new Thread(mt);
                Thread t3 = new Thread(mt);

                t1.start();
                t2.start();
                t3.start();
            }
        }
    }

</br>

    Thread-2 卖票：ticket8
    Thread-0 卖票：ticket10
    Thread-1 卖票：ticket9
    Thread-2 卖票：ticket7
    Thread-1 卖票：ticket5
    Thread-0 卖票：ticket6
    Thread-1 卖票：ticket3
    Thread-2 卖票：ticket4
    Thread-1 卖票：ticket1
    Thread-0 卖票：ticket2

- 和上面“MyThread继承于Thread”不同；这里的MyThread实现了Thread接口。
- 主线程main创建并启动3个子线程，而且这3个子线程都是基于“mt这个Runnable对象”而创建的。运行结果是这3个子线程一共卖出了10张票。这说明它们是共享了MyThread接口的。

> 问题：

1. Runnable的例子中可能会多处10张票

&emsp;&emsp;是因为t1, t2, t3共用一个任务，同时对ticket进行操作；可能导致并发问题。例如，t1和t2在开始时读取到的ticket的值都是10，而t1卖了一张票之后，ticket=9；此时t2进行卖票，而t2之前读到的票的张数是10，所以，t2卖完之后也是9张。 这样的话，问题就出现了：原本10张票，t1和t2各自都买了一张票之后，却还剩下9张！这样就导致了"实际卖出的票可能会多于10张"！

&emsp;&emsp;怎么办？

    private volatile int ticket = 10;

    @Override
    public synchronized void run() {
        for (int i = 0; i < 20; i++) {
            if (this.ticket > 0) {
                System.out.println(Thread.currentThread().getName() + " 卖票：ticket" + this.ticket--);
            }
        }
    }

- 给run()方法添加synchronized同步关键字。这保证了多个线程对是同一个MyThread对象的run()是互斥操作的。
- 将ticket修改成volatile。这样保证了，每个线程修改了ticket之后，其他线程都能读取到最新的ticket值。

2. 加入synchronized之后，只有thread-0在工作。

&emsp;&emsp;在主线程中创建了t1,t2和t3共三个线程，它们共用一个MyThread任务对象。thread-0启动之后，一直占着"MyThread同步锁"，而一个对象有且只有一个同步锁；因此，只有thread-0在工作。

> 如何解决同步问题，后期锁会讲到

## Thread中start和run的区别

- start：它的作用是启动一个新线程，新线程会执行相应的run方法。start不能被重复调用
- run：run就和普通的成员方法一样，可以被重复调用。单独调用run的话，会在当前线程中执行run，而并不会启动新线程

    public class MyThread extends Thread {
        public MyThread(String name) {
            super(name);
        }

        @Override
        public void run() {
            System.out.println(Thread.currentThread().getName() + "is running");
        }

        static class Test {
            public static void main(String[] args) {
                MyThread mt = new MyThread("mythread");
                System.out.println(Thread.currentThread().getName() + "call mythread.run");
                mt.run();
                System.out.println(Thread.currentThread().getName() + "call mythread.start");
                mt.start();
            }
        }
    }

    // Output:
    maincall mythread.run
    mainis running
    maincall mythread.start
    mythreadis running

- Thread.currentThread().getName()是用于获取“当前线程”的名字。当前线程是指正在cpu中调度执行的线程。
- mythread.run()是在“主线程main”中调用的，该run()方法直接运行在“主线程main”上。
- mythread.start()会启动“线程mythread”，“线程mythread”启动之后，会调用run()方法；此时的run()方法是运行在“线程mythread”上。

### 相关源码

    public syncchronized void start() {
        // 如果线程不是就绪状态，则抛出异常
        if (threadStatus != 0) {
            throw new IllegalThreadStateException();
        }
        // 将线程添加到ThreadGroup中
        group.add(this);

        boolean started = false;
        try {
            // 通过start0启动线程
            start0();
            // 设置started标记
            started = true;
        } finally {
            try {
                if (!started) {
                    group.threadStartFailed(this);
                }
            }catch (Throwable ignore){}
        }
    }