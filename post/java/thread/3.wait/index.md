
# 摘要

&emsp;&emsp;在Object.java中，定义了wait(), notify()和notifyAll()等接口。**wait()的作用是让当前线程进入等待状态，同时，wait()也会让当前线程释放它所持有的锁**。而notify()和notifyAll()的作用，则是**唤醒当前对象上的等待线程**；notify()是唤醒单个线程，而notifyAll()是唤醒所有的线程。

- notify 唤醒在此对象监视器上等待的单个线程
- notifyAll 唤醒在此对象监视器上等待的所有线程
- wait 让‘当前线程’处于“等待(阻塞)状态”，直到其它线程调用此对象的notify或notifyAll方法，当前线程被唤醒进入就绪状态
- wait(long timeout) 让当前线程处于“等待(阻塞状态)”，直到其它线程调用此对象的notify或notifyAll方法，或者超过指定的时间量，当前线程被唤醒进入就绪状态
- wait(long timeout, int nanos) 让当前线程处于“等待(阻塞)状态”，直到其它线程调用此对象的notify或notifyAll方法，或者其它某个线程中断当前线程，或者已经超过某个实际时间量，当前线程被唤醒进入就绪状态

## wait和notify示例

    // WaitTest.java的源码
    class ThreadA extends Thread{

        public ThreadA(String name) {
            super(name);
        }

        public void run() {
            synchronized (this) {
                System.out.println(Thread.currentThread().getName()+" call notify()");
                // 唤醒当前的wait线程
                notify();
            }
        }
    }

    public class WaitTest {

        public static void main(String[] args) {

            ThreadA t1 = new ThreadA("t1");

            synchronized(t1) {
                try {
                    // 启动“线程t1”
                    System.out.println(Thread.currentThread().getName()+" start t1");
                    t1.start();

                    // 主线程等待t1通过notify()唤醒。
                    System.out.println(Thread.currentThread().getName()+" wait()");
                    t1.wait();

                    System.out.println(Thread.currentThread().getName()+" continue");
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        }
    }

</br>

    main start t1
    main wait()
    t1 call notify()
    main continue

&emsp;&emsp;如下图，说明了主线程和线程1的流程

1. 图中"主线程"代表"主线程main"，"线程t1"代表WaitTest中启动的"线程t1"，而"锁"代表"t1这个对象的同步锁"
2. "主线程"通过new ThreadA("t1")新建"线程t1"，随后通过synchronized(t1)获取"t1对象的同步锁"，然后调用t1.start启动"线程t1"
3. "主线程"执行t1.wait释放"t1对象的锁"并进入"等待阻塞状态"。等待t1对象的线程通过notify将其唤醒
4. "线程t1"运行后，通过synchronized(this)获取"当前对象的锁"，接着调用notify唤醒"当前对象上的等待线程"，也就是唤醒"主线程"
5. "线程t1"运行完毕之后，释放"当前对象的锁"。紧接着，"主线程"获取"t1对象的锁"，然后接着运行。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/juc/18183712-f04899f92aaa43b6a33a85fecfa60a9d.png)

> t1.wait()应该是让“线程t1”等待；但是，为什么却是让“主线程main”等待了呢？

&emsp;&emsp;“当前线程”在调用wait时，必须拥有该对象的同步锁。该线程调用wait之后，会释放该锁；然后一直等待直到“其它线程”调用对象的同步锁notify或notifyAll方法。然后，该线程继续等待直到它重新获取"该对象的同步锁"，就可以接着运行

&emsp;&emsp;**jdk的解释中，说wait()的作用是让“当前线程”等待，而“当前线程”是指正在cpu上运行的线程！**

&emsp;&emsp;这也意味着，虽然t1.wait()是通过“线程t1”调用的wait()方法，但是调用t1.wait()的地方是在“主线程main”中。而主线程必须是“当前线程”，也就是运行状态，才可以执行t1.wait()。所以，此时的“当前线程”是“主线程main”！因此，t1.wait()是让“主线程”等待，而不是“线程t1”！

## wait(long timeout) 和 notify

    // WaitTimeoutTest.java的源码
    class ThreadA extends Thread{

        public ThreadA(String name) {
            super(name);
        }

        public void run() {
            System.out.println(Thread.currentThread().getName() + " run ");
            // 死循环，不断运行。
            while(true)
                ;
        }
    }

    public class WaitTimeoutTest {

        public static void main(String[] args) {

            ThreadA t1 = new ThreadA("t1");

            synchronized(t1) {
                try {
                    // 启动“线程t1”
                    System.out.println(Thread.currentThread().getName() + " start t1");
                    t1.start();

                    // 主线程等待t1通过notify()唤醒 或 notifyAll()唤醒，或超过3000ms延时；然后才被唤醒。
                    System.out.println(Thread.currentThread().getName() + " call wait ");
                    t1.wait(3000);

                    System.out.println(Thread.currentThread().getName() + " continue");
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        }
    }

</br>

    main start t1
    main call wait 
    t1 run                  // 大约3秒之后...输出“main continue”
    main continue

&emsp;&emsp;如下图，说明了主线程和线程1的流程

1. 图中"主线程"代表"主线程main"，"线程t1"代表WaitTest中启动的"线程t1"，而"锁"代表"t1这个对象的同步锁"
2. 主线程main执行t1.start()启动“线程t1”。
3. 主线程main执行t1.wait(3000)，此时，主线程进入“阻塞状态”。需要“用于t1对象锁的线程通过notify() 或者 notifyAll()将其唤醒” 或者 “超时3000ms之后”，主线程main才进入到“就绪状态”，然后才可以运行。
4. “线程t1”运行之后，进入了死循环，一直不断的运行。
5. 超时3000ms之后，主线程main会进入到“就绪状态”，然后接着进入“运行状态”。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/juc/18183848-d4c16bf2760847afa4fede6a9d959083.png)

## wait() 和 notifyAll()

    public class NotifyAllTest {

        private static Object obj = new Object();
        public static void main(String[] args) {

            ThreadA t1 = new ThreadA("t1");
            ThreadA t2 = new ThreadA("t2");
            ThreadA t3 = new ThreadA("t3");
            t1.start();
            t2.start();
            t3.start();

            try {
                System.out.println(Thread.currentThread().getName()+" sleep(3000)");
                Thread.sleep(3000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }

            synchronized(obj) {
                // 主线程等待唤醒。
                System.out.println(Thread.currentThread().getName()+" notifyAll()");
                obj.notifyAll();
            }
        }

        static class ThreadA extends Thread{

            public ThreadA(String name){
                super(name);
            }

            public void run() {
                synchronized (obj) {
                    try {
                        // 打印输出结果
                        System.out.println(Thread.currentThread().getName() + " wait");

                        // 唤醒当前的wait线程
                        obj.wait();

                        // 打印输出结果
                        System.out.println(Thread.currentThread().getName() + " continue");
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }
            }
        }
    }

</br>

    t1 wait
    main sleep(3000)
    t3 wait
    t2 wait
    main notifyAll()
    t2 continue
    t3 continue
    t1 continue

1. 主线程中新建并且启动了3个线程"t1", "t2"和"t3"。
2. 主线程通过sleep(3000)休眠3秒。在主线程休眠3秒的过程中，我们假设"t1", "t2"和"t3"这3个线程都运行了。以"t1"为例，当它运行的时候，它会执行obj.wait()等待其它线程通过notify()或额nofityAll()来唤醒它；相同的道理，"t2"和"t3"也会等待其它线程通过nofity()或nofityAll()来唤醒它们。
3. 主线程休眠3秒之后，接着运行。执行 obj.notifyAll() 唤醒obj上的等待线程，即唤醒"t1", "t2"和"t3"这3个线程。 紧接着，主线程的synchronized(obj)运行完毕之后，主线程释放“obj锁”。这样，"t1", "t2"和"t3"就可以获取“obj锁”而继续运行了！

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/juc/18183923-95275c066212410f96181704a681f453.png)

# 为什么notify, wait等函数定义在Object，而不是Thread中

&emsp;&emsp;Object中的wait(), notify()等函数，和synchronized一样，会对“对象的同步锁”进行操作。

&emsp;&emsp;wait会使“当前线程”等待，因为线程进入等待状态，所以线程应该释放它持有的“同步锁”，否则其它线程获取不到该“同步锁”而无法运行。那么，线程调用wait之后，会释放它所持有的“同步锁”；而且，根据前面的介绍，我们知道：等待线程可以被notify或notifyAll唤醒。现在，思考一个问题：**notify是根据什么唤醒等待线程的？或者说，wait等待线程和notify之间是通过什么关联起来的？**。答案是：“对象同步锁”。

&emsp;&emsp;负责唤醒等待线程的那个线程(我们称之为唤醒线程)，它只有在获取“该对象的同步锁”(这里的同步锁必须和等待线程的同步锁是同一个)，并且调用notify或notifyAll方法之后，才能唤醒等待线程。虽然，等待线程被唤醒；但是，它不能立即执行，因为唤醒线程还持有“该对象的同步锁”。必须等待唤醒线程释放“对象的同步锁”之后，等待线程才能获取到“对象的同步锁”而继续执行

&emsp;&emsp;总之，notify, wait依赖于"同步锁"，而"同步锁"是对象锁持有，并且每个对象有且仅有一个，所以，它们被定义在Object中。

# yield和wait的比较

- wait是让当先线程由"运行状态"进入到“等待阻塞状态”，而yield是让线程由"运行状态"进入到"就绪状态"
- wait会让当前线程释放它所持有对象的同步锁，而yield方法不会释放锁

# wait与sleep的比较

- wait是让当前线程由“运行状态”进入到“等待(阻塞)状态”，而sleep是让线程由“运行状态”进入到“休眠(阻塞)状态”
- wait会释放对象的同步锁，而sleep则不会

# join

&emsp;&emsp;join的作用：让"主线程"等待"自线程"结束后才能继续运行。我们通过简单例子理解一下：

    // 主线程
    public class Father extends Thread {
        public void run() {
            Son s = new Son();
            s.start();
            s.join();
            ...
        }
    }
    // 子线程
    public class Son extends Thread {
        public void run() {
            ...
        }
    }

&emsp;&emsp;上面的有两个类Father(主线程类)和Son(子线程类)。因为Son是在Father中创建并启动的，所以，Father是主线程类，Son是子线程类。在Father主线程中，通过new Son()新建“子线程s”。接着通过s.start()启动“子线程s”，并且调用s.join()。在调用s.join()之后，Father主线程会一直等待，直到“子线程s”运行完毕；在“子线程s”运行完毕之后，Father主线程才能接着运行。 这也就是我们所说的“join()的作用，是让主线程会等待子线程结束之后才能继续运行”！

## join源码

    public final void join() throws InterruptedException {
        join(0);
    }

    public final synchronized void join(long millis)
    throws InterruptedException {
        long base = System.currentTimeMillis();
        long now = 0;

        if (millis < 0) {
            throw new IllegalArgumentException("timeout value is negative");
        }

        if (millis == 0) {
            while (isAlive()) {
                wait(0);
            }
        } else {
            while (isAlive()) {
                long delay = millis - now;
                if (delay <= 0) {
                    break;
                }
                wait(delay);
                now = System.currentTimeMillis() - base;
            }
        }
    }

&emsp;&emsp;从代码中，我们可以发现。当millis==0时，会进入while(isAlive())循环；**即只要子线程是活的，主线程就不停的等待**。我们根据上面解释join()作用时的代码来理解join()的用法！

&emsp;&emsp;问题：虽然s.join()被调用的地方是发生在“Father主线程”中，但是s.join()是通过“子线程s”去调用的join()。那么，join()方法中的isAlive()应该是判断“子线程s”是不是Alive状态；对应的wait(0)也应该是“让子线程s”等待才对。但如果是这样的话，s.join()的作用怎么可能是“让主线程等待，直到子线程s完成为止”呢，应该是让"子线程等待才对(因为调用子线程对象s的wait方法嘛)"？

&emsp;&emsp;答案：wait()的作用是让“当前线程”等待，而这里的“当前线程”是指当前在CPU上运行的线程。所以，虽然是调用子线程的wait()方法，但是它是通过“主线程”去调用的；所以，休眠的是主线程，而不是“子线程”！

## join示例

    // JoinTest.java的源码
    public class JoinTest{ 

        public static void main(String[] args){ 
            try {
                ThreadA t1 = new ThreadA("t1"); // 新建“线程t1”

                t1.start();                     // 启动“线程t1”
                t1.join();                        // 将“线程t1”加入到“主线程main”中，并且“主线程main()会等待它的完成”
                System.out.printf("%s finish\n", Thread.currentThread().getName()); 
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        } 

        static class ThreadA extends Thread{

            public ThreadA(String name){ 
                super(name); 
            } 
            public void run(){ 
                System.out.printf("%s start\n", this.getName()); 

                // 延时操作
                for(int i=0; i <1000000; i++)
                ;

                System.out.printf("%s finish\n", this.getName()); 
            } 
        } 
    }

</br>

    t1 start
    t1 finish
    main finish

&emsp;&emsp;说明：

1. 在“主线程main”中通过 new ThreadA("t1") 新建“线程t1”。 接着，通过 t1.start() 启动“线程t1”，并执行t1.join()。
2. 执行t1.join()之后，“主线程main”会进入“阻塞状态”等待t1运行结束。“子线程t1”结束之后，会唤醒“主线程main”，“主线程”重新获取cpu执行权，继续运行。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/juc/18184312-a72a58e2bda54b17bf669f325ecda377.png)

# interrupt

&emsp;&emsp;interrupt()的作用是中断本线程。

&emsp;&emsp;本线程中断自己是被允许的；其它线程调用本线程的interrupt()方法时，会通过checkAccess()检查权限。这有可能抛出SecurityException异常。

&emsp;&emsp;如果本线程是处于阻塞状态：调用线程的wait(), wait(long)或wait(long, int)会让它进入等待(阻塞)状态，或者调用线程的join(), join(long), join(long, int), sleep(long), sleep(long, int)也会让它进入阻塞状态。若线程在阻塞状态时，调用了它的interrupt()方法，那么它的“中断状态”会被清除并且会收到一个InterruptedException异常。例如，线程通过wait()进入阻塞状态，此时通过interrupt()中断该线程；调用interrupt()会立即将线程的中断标记设为“true”，但是由于线程处于阻塞状态，所以该“中断标记”会立即被清除为“false”，同时，会产生一个InterruptedException的异常。

&emsp;&emsp;如果线程被阻塞在一个Selector选择器中，那么通过interrupt()中断它时；线程的中断标记会被设置为true，并且它会立即从选择操作中返回。

&emsp;&emsp;如果不属于前面所说的情况，那么通过interrupt()中断线程时，它的中断标记会被设置为“true”。

&emsp;&emsp;中断一个“已终止的线程”不会产生任何操作。

## 终止线程的方式

&emsp;&emsp;下面，我先分别讨论线程在“阻塞状态”和“运行状态”的终止方式，然后再总结出一个通用的方式。

### 终止处于阻塞状态的线程

&emsp;&emsp;当线程由于被调用了sleep(), wait(), join()等方法而进入阻塞状态；若此时调用线程的interrupt()将线程的中断标记设为true。由于处于阻塞状态，中断标记会被清除，同时产生一个InterruptedException异常。将InterruptedException放在适当的位置就能终止线程，形式如下：

    @Override
    public void run() {
        try {
            while (true) {
                // 执行任务...
            }
        } catch (InterruptedException ie) {  
            // 由于产生InterruptedException异常，退出while(true)循环，线程终止！
        }
    }

### 终止处于运行状态的线程

&emsp;&emsp;通常，我们通过“标记”方式终止处于“运行状态”的线程。其中，包括“中断标记”和“额外添加标记”。

- 通过“中断标记”终止线程

    @Override
    public void run() {
        while (!isInterrupted()) {
            // 执行任务...
        }
    }

&emsp;&emsp;说明：isInterrupted()是判断线程的中断标记是不是为true。当线程处于运行状态，并且我们需要终止它时；可以调用线程的interrupt()方法，使用线程的中断标记为true，即isInterrupted()会返回true。此时，就会退出while循环。

&emsp;&emsp;注意：**interrupt()并不会终止处于“运行状态”的线程！它会将线程的中断标记设为true**。

- 通过“额外添加标记”。

</br>

    private volatile boolean flag= true;
    protected void stopTask() {
        flag = false;
    }

    @Override
    public void run() {
        while (flag) {
            // 执行任务...
        }
    }

&emsp;&emsp;说明：线程中有一个flag标记，它的默认值是true；并且我们提供stopTask()来设置flag标记。当我们需要终止该线程时，调用该线程的stopTask()方法就可以让线程退出while循环。

&emsp;&emsp;注意：将flag定义为volatile类型，是为了保证flag的可见性。即其它线程通过stopTask()修改了flag之后，本线程能看到修改后的flag的值。

> 综合线程处于“阻塞状态”和“运行状态”的终止方式，比较通用的终止线程的形式如下：

    @Override
    public void run() {
        try {
            // 1. isInterrupted()保证，只要中断标记为true就终止线程。
            while (!isInterrupted()) {
                // 执行任务...
            }
        } catch (InterruptedException ie) {  
            // 2. InterruptedException异常保证，当InterruptedException异常产生时，线程被终止。
        }
    }

### 终止线程的示例

    // Demo1.java的源码
    class MyThread extends Thread {

        public MyThread(String name) {
            super(name);
        }

        @Override
        public void run() {
            try {
                int i=0;
                while (!isInterrupted()) {
                    Thread.sleep(100); // 休眠100ms
                    i++;
                    System.out.println(Thread.currentThread().getName()+" ("+this.getState()+") loop " + i);
                }
            } catch (InterruptedException e) {
                System.out.println(Thread.currentThread().getName() +" ("+this.getState()+") catch InterruptedException.");
            }
        }
    }

    public class Demo1 {

        public static void main(String[] args) {
            try {
                Thread t1 = new MyThread("t1");  // 新建“线程t1”
                System.out.println(t1.getName() +" ("+t1.getState()+") is new.");

                t1.start();                      // 启动“线程t1”
                System.out.println(t1.getName() +" ("+t1.getState()+") is started.");

                // 主线程休眠300ms，然后主线程给t1发“中断”指令。
                Thread.sleep(300);
                t1.interrupt();
                System.out.println(t1.getName() +" ("+t1.getState()+") is interrupted.");

                // 主线程休眠300ms，然后查看t1的状态。
                Thread.sleep(300);
                System.out.println(t1.getName() +" ("+t1.getState()+") is interrupted now.");
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }
    }

</br>

    t1 (NEW) is new.
    t1 (RUNNABLE) is started.
    t1 (RUNNABLE) loop 1
    t1 (RUNNABLE) loop 2
    t1 (TIMED_WAITING) is interrupted.
    t1 (RUNNABLE) catch InterruptedException.
    t1 (TERMINATED) is interrupted now.

- 主线程main中通过new MyThread("t1")创建线程t1，之后通过t1.start()启动线程t1。
- t1启动之后，会不断的检查它的中断标记，如果中断标记为“false”；则休眠100ms。
- t1休眠之后，会切换到主线程main；主线程再次运行时，会执行t1.interrupt()中断线程t1。t1收到中断指令之后，会将t1的中断标记设置“false”，而且会抛出InterruptedException异常。在t1的run()方法中，是在循环体while之外捕获的异常；因此循环被终止。

&emsp;&emsp;额外添加标记：

    public class MyThread extends Thread {

        private volatile boolean flag = true;

        private void stopTask() {
            flag = false;
        }

        public MyThread(String name) {
            super(name);
        }

        @Override
        public void run() {
            synchronized (this) {
                try {
                    int i = 0;
                    while (flag) {
                        // 不释放锁
                        Thread.sleep(100);
                        i++;
                        System.out.println(Thread.currentThread().getName() + " (" + this.getState() + ") loop " + i);
                    }
                    System.out.println(Thread.currentThread().getName() + "(" + this.getState() + ") down.");
                } catch (InterruptedException e) {
                    System.out.println(Thread.currentThread().getName() + " (" + this.getState() + ") catch InterruptedException.");
                }
            }
        }

        static class Test {
            public static void main(String[] args) {
                try {
                    MyThread t1 = new MyThread("t1");
                    System.out.println(t1.getName() + " (" + t1.getState() + ") is new.");

                    t1.start();
                    System.out.println(t1.getName() + " (" + t1.getState() + ") is started.");

                    Thread.sleep(300);
                    t1.stopTask();

                    Thread.sleep(300);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        }
    }

</br>

    t1 (NEW) is new.
    t1 (TIMED_WAITING) is started.
    t1 (RUNNABLE) loop 1
    t1 (RUNNABLE) loop 2
    t1 (RUNNABLE) loop 3
    t1(RUNNABLE) down.

## interrupted()和isInterrupted()

&emsp;&emsp;interrupted() 和 isInterrupted()都能够用于检测对象的“中断标记”。

&emsp;&emsp;区别是，interrupted()除了返回中断标记之外，它还会清除中断标记(即将中断标记设为false)；而isInterrupted()仅仅返回中断标记。

# 线程优先级

&emsp;&emsp;Java中的线程优先级番位是1～10，默认的优先级是5.高优先级线程会优于低优先级线程执行。

&emsp;&emsp;java 中有两种线程：用户线程和守护线程。可以通过isDaemon()方法来区别它们：如果返回false，则说明该线程是“用户线程”；否则就是“守护线程”。

&emsp;&emsp;用户线程一般用户执行用户级任务，而守护线程也就是“后台线程”，一般用来执行后台任务。需要注意的是：Java虚拟机在“用户线程”都结束后会后退出。

&emsp;&emsp;当Java虚拟机启动时，通常有一个单一的非守护线程（该线程通过是通过main()方法启动）。JVM会一直运行直到下面的任意一个条件发生，JVM就会终止运行：

1. 调用了exit()方法，并且exit()有权限被正常执行。
2. 所有的“非守护线程”都死了(即JVM中仅仅只有“守护线程”)。