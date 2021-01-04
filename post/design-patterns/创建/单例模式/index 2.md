
# 摘要

&emsp;&emsp;单例模式(Singleton Pattern) 是Java中最简单的设计模式之一。它提供了一种创建对象的最佳方法。

&emsp;&emsp;这种模式涉及到一个单一的类，该类负责创建自己的对象，同时确保只有单个对象被创建。这个类提供了一种访问其唯一的对象的方式，可以直接访问，不需要实例化该类的对象。

- **注意**：
1. 单例类只能有一个实例。
2. 单例类必须自己创建的唯一实例。
3. 单例类必须给所有其它对象提供这一实例。

## 介绍

- **意图**：保证一个类仅有一个实例，并提供一个访问它的全局访问点。
- **主要解决**: 一个全局使用的类频繁地创建于销毁。
- **何时使用**：当您想控制实例数目，节省系统资源的时候。
- **如何解决**: 判断系统是否已经有这个单例，如果有则返回，如果没有则创建。
- **关键代码**: 构造函数是私有的。
- **应用实例**：
1. 一个班级只有一个班主任。
2. Windows是多进程多线程的，在操作一个文件的时候，就不可避免地出现多个进程或线程同时操作一个文件的现象，所有所有文件的处理必须通过唯一的实例来进行。
3. 一些设备管理器常常设计为单例模式，比如一个电脑有两台打印机，在输出的时候就要处理不能两台打印机打印同一个文件。

## 实现

&emsp;&emsp;我们将创建一个SingleObject类。SingleObject类有它的私有构造函数和本身的一个静态实例。

&emsp;&emsp;SingleObject类提供了一个静态方法，供外界获取它的静态实例。SingletonPatternDemo类使用SingleObject类来获取SingleObject对象。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/design-pattern/create/62576915-36E0-4B67-B078-704699CA980A.jpg)

### 步骤1

&emsp;&emsp;

    package default;

    public class SingleObject {
        // 创建SingleObject的一个对象
        private static SingleObject instance = new SingleObject();
        //让构造函数为private，这样该类就不会被实例化
        private SingleObject(){}

        //获取唯一可用的对象
        public static SingleObject getInstance() {
            return instance;
        }

        public void showMessage(){
            System.out.println("Hello World!);
        }
    }

### 步骤2

&emsp;&emsp;从singleton类获取唯一的对象

    package default;

    public class SingletonPatternDemo {
        public static main(String[] args) {
            SingleObject object = SingleObject.getInstance();
            object.showMessage();
        }
    }


## 单例模式的几种实现方式

### 懒汉式，线程不安全

&emsp;&emsp;这种方式是最基本的实现方式，这种实现最大的问题就是不支持多线程。因为没有加锁synchronized, 所以严格意义上它并不算单例模式。这种方式 lazy loading很明显，不要求线程安全，在多线程下不能正常工作

    package default;

    public class Singleton {
        private static Singleton instance;
        private SIngleton() {}

        public static Singleton getInstance() {
            if (instance == null) {
                instance = new Singleton();
            }
            return instance;
        }
    }

### 懒汉式，线程安全

&emsp;&emsp;这种方式具备很好的lazy loading, 能够在多线程中很好的工作，但是，效率底下，99%情况下不需要同步。优点：第一次调用才初始化，避免内存浪费。缺点：必须加锁才能保证单例，但加锁会影响效率。`加锁其实只需要在第一次初始化的时候用到，之后的调用都没必要再加锁。`

        package default;

        public class Singleton {
            private static Singleton instance;
            private Singleton() {}

            public static synchronized Singleton getInstance() {
                if (instance == null) {
                    instance = new Singleton();
                }
                return instance;
            }
        }

### 饿汉式

&emsp;&emsp;这种方式比较常用，但很容易产生垃圾。优点：没有加锁，执行效率高。缺点：类加载时就初始化，浪费内存。它基于classloader机制避免了多线程的同步问题，不过，instance在类装载时就实例化，虽然导致类装载的原因有很多种，在单例模式中大多数都是调用getInstance方法，但是也不能确定有其它的方式导致类装载，这时候初始化instance显然没有达到lazy loading的效果。

    package default;

    public class Singleton() {
        private static Singleton instance = new Singleton();
        private Singleton() {}

        public static Singleton getInstance() {
            return instance;
        }
    }

### 双重校验锁(double-checked locking)

&emsp;&emsp;新判断对象是否已经被初始化，再决定要不要加锁。执行双重检查是因为，如果多个线程同时通过了第一次检查，并且其中一个线程首先通过了第二次检查并实例化了对象，那么剩余通过了第一次检查的线程就不会再去实例化对象。这样，除了初始化的时候会出现加锁的情况，后续的所有调用都会避免加锁而直接返回，解决了性能消耗的问题。

&emsp;&emsp;实例化对象步骤：

1. 分配内存空间
2. 初始化对象
3. 将对象指向刚分配的内存空间

&emsp;&emsp;有些编译器为了性能，会将第二步和第三步进行**重排序**。使用`volatile`关键字后，重排序被禁止，所有的写操作都将发生在读操作之前。

    package default;

    public class Singleton {
        private volatile static Singlrton singleton;
        private Singleton() {}

        public static Singleton getSingle() {
            if (null == singleton) {
                synchronized(Singleton.class) {
                    if (null == singleton) {
                        singleton = new Singleton();
                    }
                }
            }
            return singleton;
        }
    }

### 静态内部类

&emsp;&emsp;线程安全，懒加载。这种方式能达到双检锁方式一样的功效，但实现更简单。对静态域使用延迟初始化，应使用这种方式而不是双检锁方式。这种方式只适用于静态域的情况，双检锁方式可在实例域需要延迟初始化时使用。

    package default;

    public class Singleton {
        private static class SingletonHolder {
            private static final Singleton INSTANCE = new Singleton();
        }
        private Singleton(){}
        public static final Singleton getInstance() {
            return SingletonHoler.INSTANCE;
        }
    }

&emsp;&emsp;这种方式同样利用了classloader机制来保证初始化instance时只有一个线程，它跟懒汉式不同的是：懒汉式只要Singleton类被装载了，那么instance就会被实例化(没有达到lazy loading效果)，而这种方式是Singleton类被装载，instance不一定被初始化。因为SingletonHolder类没有被使用，只有通过显示调用getInstance方法时，才会显示装载SingletonHolder类，从而实例化instance。

### 枚举

&emsp;&emsp;这种方式是实现单例模式的最佳方法。简洁，自动支持序列化，绝对防止多次实例化。

    package default;

    public enum Singleton {
        INSTANCE;
        public void whateverMethod() {}
    }

> 一般情况下建议使用饿汉式，只有在明确实现lazy loading效果时，才会使用静态内部类方式。如果涉及反序列化创建对象时，可以尝试枚举。如果有其它特殊需求，可以考虑使用双检锁。


# Golang版本

&emsp;&emsp;直接使用`sync`包中的once，线程安全。

    package singleton

    import "sync"

    type Singletonn struct{}

    var singleton *Singleton
    var once sync.Once

    func GetInstance() *Singleton {
        once.Do(func(){
            singleton = &Singleton{}
        })
        return singleton
    }