
# 摘要

&emsp;&emsp;面向对象编程有三大特性：封装、继承、多态。

&emsp;&emsp;封装隐藏了类的内部实现机制，可以在不影响使用的情况下改变类的内部结构，同时也保护了数据。对外界而已它的内部细节是隐藏的，暴露给外界的只是它的访问方法。

&emsp;&emsp;继承是为了重用父类代码。两个类若存在IS-A的关系就可以使用继承。同时继承也为实现多态做了铺垫。那么什么是多态呢？多态的实现机制又是什么？

&emsp;&emsp;所谓多态就是指**程序中定义的引用变量所指向的具体类型和通过该引用变量发出的方法调用在编程时并不确定，而是在程序运行期间才确定。因为在程序运行时才确定具体的类，这样，不用修改源程序代码，就可以让引用变量绑定到各种不同的类实现上，从而导致该引用调用的具体方法随之改变，即不修改程序代码就可以改变程序运行时所绑定的具体代码，让程序可以选择多个运行状态，这就是多态性**。

&emsp;&emsp;比如你是一个酒神，对酒情有独钟。某日回家发现桌上有几个杯子里面都装了白酒，从外面看我们是不可能知道这是些什么酒，只有喝了之后才能够猜出来是何种酒。你一喝，这是剑南春、再喝这是五粮液、再喝这是酒鬼酒….在这里我们可以描述成如下：

    酒 a = 剑南春

    酒 b = 五粮液

    酒 c = 酒鬼酒

&emsp;&emsp;这里所表现的的就是多态。剑南春、五粮液、酒鬼酒都是酒的子类，我们只是通过酒这一个父类就能够引用不同的子类，这就是多态——**我们只有在运行的时候才会知道引用变量所指向的具体实例对象**。

## 多态

&emsp;&emsp;诚然，要理解多态我们就必须要明白什么是“向上转型”。在继承中我们简单介绍了向上转型，这里就在啰嗦下：在上面的喝酒例子中，酒（Wine）是父类，剑南春（JNC）、五粮液（WLY）、酒鬼酒（JGJ）是子类。我们定义如下代码：

    JNC a = new JNC();

&emsp;&emsp;对于这个代码我们非常容易理解无非就是实例化了一个剑南春的对象嘛！但是这样呢？

    Wine a = new JNC();

&emsp;&emsp;在这里我们这样理解，这里定义了一个Wine 类型的a，它指向JNC对象实例。由于JNC是继承与Wine，所以JNC可以自动向上转型为Wine，所以a是可以指向JNC实例对象的。这样做存在一个非常大的好处，在继承中我们知道子类是父类的扩展，它可以提供比父类更加强大的功能，如果我们定义了一个指向子类的父类引用类型，那么它除了能够引用父类的共性外，还可以使用子类强大的功能。

&emsp;&emsp;但是向上转型存在一些缺憾，那就是它必定会导致一些方法和属性的丢失，而导致我们不能够获取它们。所以**父类类型的引用可以调用父类中定义的所有属性和方法，对于只存在与子类中的方法和属性它就望尘莫及了**。

    public class Wine {
        public void fun1() {
            System.out.println("Wine`s fun1");
            fun2();
        }

        public void fun2() {
            System.out.println("Wine`s fun2");
        }
    }

</br>

    public class JNC extends Wine {
        // 重载
        public void fun1(String a) {
            System.out.println("JNC`s fun1");
            fun2();
        }

        // 重写
        public void fun2() {
            System.out.println("JNC`s fun2");
        }
    }

</br>

    public class Test {
        public static void main(String[] args) {
            Wine a = new JNC();
            a.fun1();
        }
    }

    Output:
    Wine`s fun1
    JNC`s fun2

&emsp;&emsp;从程序的运行结果中我们发现，a.fun1()首先是运行父类Wine中的fun1().然后再运行子类JNC中的fun2()。

&emsp;&emsp;分析：在这个程序中子类JNC重载了父类Wine的方法fun1()，重写fun2()，而且重载后的fun1(String a)与 fun1()不是同一个方法，由于父类中没有该方法，向上转型后会丢失该方法，所以执行JNC的Wine类型引用是不能引用fun1(String a)方法。而子类JNC重写了fun2() ，那么指向JNC的Wine引用会调用JNC中fun2()方法。

&emsp;&emsp;所以对于多态我们可以总结如下：

&emsp;&emsp;**指向子类的父类引用由于向上转型了，它只能访问父类中拥有的方法和属性，而对于子类中存在而父类中不存在的方法，该引用是不能使用的，尽管是重载该方法。若子类重写了父类中的某些方法，在调用该些方法的时候，必定是使用子类中定义的这些方法（动态连接、动态调用）**。

&emsp;&emsp;对于面向对象而已，多态分为`编译时多态`和`运行时多态`。其中编辑时多态是静态的，主要是指方法的重载，它是根据参数列表的不同来区分不同的函数，通过编译之后会变成两个不同的函数，在运行时谈不上多态。而运行时多态是动态的，它是通过动态绑定来实现的，也就是我们所说的多态性。

## 多态的实现

### 实现条件

&emsp;&emsp;在刚刚开始就提到了继承在为多态的实现做了准备。子类Child继承父类Father，我们可以编写一个指向子类的父类类型引用，该引用既可以处理父类Father对象，也可以处理子类Child对象，当相同的消息发送给子类或者父类对象时，该对象就会根据自己所属的引用而执行不同的行为，这就是多态。即多态性就是相同的消息使得不同的类做出不同的响应。

> Java实现多态有三个必要条件：继承、重写、向上转型

- 继承：在多态中必须存在有继承关系的子类和父类。
- 重写：子类对父类中某些方法进行重新定义，在调用这些方法时就会调用子类的方法。
- 向上转型：在多态中需要将子类的引用赋给父类对象，只有这样该引用才能够具备技能调用父类的方法和子类的方法。

&emsp;&emsp;只有满足了上述三个条件，我们才能够在**同一个继承结构中使用统一的逻辑实现代码处理不同的对象，从而达到执行不同的行为**。

&emsp;&emsp;对于Java而言，它多态的实现机制遵循一个原则：当超类对象引用变量引用子类对象时，被引用对象的类型而不是引用变量的类型决定了调用谁的成员方法，但是这个被调用的方法必须是在超类中定义过的，也就是说被子类覆盖的方法。

### 实现形式

> 在Java中有两种形式可以实现多态。继承和接口

#### 基于继承实现多态

    public class Wine {
        private String name;
        // name settet getter

        public Wine(){}

        public String drink() {
            return "喝的是 " + name;
        }

        @Override
        public String toString() {
            return "";
        }
    }

</br>

    public class JNC extends Wine {
        public JNC() {
            setName("JNC");
        }

        @Override
        public String drink() {
            return "喝的是 " + getName();
        }

        @Override
        public Strinng toString() {
            return "Wine : " + getName();
        }
    }

</br>

    public class JGJ extends Wine {
        public JGJ() {
            setName("JGJ");
        }

        @Override
        public String drink() {
            return "喝的是 " + getName();
        }

        @Override
        public String toString() {
            return "Wine : " + getName();
        }
    }

</br>

    public class Test {
        public static void main(String[] args) {
            Wine[] wines = new Wine[2];
            Wine jnc = new JNC();
            Wine jgj = new JGJ();
            wines[0] = jnc;
            wines[1] = jgj;

            for (Wine w : wines) {
                System.out.println(w.toString() + "---" + w.drink());
            }
        }
    }
    Output:
    Wine : JNC--喝的是 JNC
    WIne : JGJ--喝的是 JGJ

&emsp;&emsp;在上面的代码中JNC、JGJ继承Wine，并且重写了drink()、toString()方法，程序运行结果是调用子类中方法，输出JNC、JGJ的名称，这就是多态的表现。不同的对象可以执行相同的行为，但是他们都需要通过自己的实现方式来执行，这就要得益于向上转型了。

&emsp;&emsp;对于引用子类的父类类型，在处理该引用时，它适用于继承该父类的所有子类，子类对象的不同，对方法的实现也就不同，执行相同动作产生的行为也就不同。

&emsp;&emsp;如果父类是抽象类，那么子类必须要实现父类中所有的抽象方法，这样该父类所有的子类一定存在统一的对外接口，但其内部的具体实现可以各异。这样我们就可以使用顶层类提供的统一接口来处理该层次的方法。

#### 基于接口实现多态

&emsp;&emsp;继承是通过重写父类的同一方法的几个不同子类来体现的，那么接口就是通过实现接口并覆盖接口中同一方法的几不同的类体现的。

&emsp;&emsp;在接口的多态中，指向接口的引用必须是指定这实现了该接口的一个类的实例程序，在运行时，根据对象引用的实际类型来执行对应的方法。

&emsp;&emsp;继承都是单继承，只能为一组相关的类提供一致的服务接口。但是接口可以是多实现，它能够利用一组相关或者不相关的接口进行组合与扩充，能够对外提供一致的服务接口。所以它相对于继承来说有更好的灵活性。

### 经典实例

    public class A {
        // 实际开发中没有人会这么些
        public String show(D obj) {
            return ("A and D");
        }

        public String show(A obj) {
            return ("A and A");
        }
    }

    public class B extends A {
        public String show(B obj) {
            return ("B and B");
        }
        // 重写
        public String show(A obj) {
            return ("B and A");
        }
    }

    public class C extends B {

    }

    public class D extends B {

    }

</br>

    public class Test {
        public static void main(String[] args) {
            A a1 = new A();
            A a2 = new B();
            B b = new B();
            C c = new C();
            D d = new D();

            System.out.println("1--" + a1.show(b));
            System.out.println("2--" + a1.show(c));
            System.out.println("3--" + a1.show(d));
            System.out.println("4--" + a2.show(b));
            System.out.println("5--" + a2.show(c));
            System.out.println("6--" + a2.show(d));
            System.out.println("7--" + b.show(b));
            System.out.println("8--" + b.show(c));
            System.out.println("9--" + b.show(d));   
        }
    }

    Output:
    1--A and A
    2--A and A
    3--A and D
    4--B and A
    5--B and A
    6--A and D
    7--B and B
    8--B and B
    9--A and D

&emsp;&emsp;在继承链中对象方法的调用存在一个优先级:**this.show(o), super.show(o), this.show((super)o), super.show((super)o)**。这就是典型的向上转型。

&emsp;&emsp;我们分析第5个，a2.show(c)，a2是A类型的引用，所以this代表A，a2.show(c)，它在A类中没有找到，于是super.show(o),由于A没有超类，所以到this.show((super)o),this表示A，(super)o表示B和A(向上转型)，在A中找到了show(A),同时由于a2的引用是B，且B类重写了show(A)，因此最终会调用子类B类的show(A).

&emsp;&emsp;我们分析第8个，b.show(c)，b的this是B，它在B类中没有找到，于是super.show(o),此时super是A，在A中也没有找到，跳到this.show((super)o)，此时this还是B，C向上转型为B，找到了B.show(B).

&emsp;&emsp;我们分析就9个，b.show(d)，乍看上去，可第8个一样，但是在第super.show(o)中找到了A.show(D).所以执行A.show(D)。但是，我们一般在项目中是不会存在这种继承关系的。A是顶级父类，不会在抽象方法中引用子类。