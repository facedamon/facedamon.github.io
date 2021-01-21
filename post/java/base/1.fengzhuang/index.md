
# 摘要

&emsp;&emsp;封装从字面上来理解就是包装的意思，就是隐藏信息，是指**利用抽象数据类型将数据和基于数据的操作封装在一起，使其构成一个不可分割的独立实体**，数据被保护在抽象数据类型的内部，尽可能地隐藏内部细节，只保留一些对外接口使之与外部发生联系。系统的其它对象只能通过包裹在数据外面的已经授权的操作来与这个封装的对象进行交流。也就是说，**用户是无需知道对象内部的细节，但可以通过该对象对外提供的接口来访问对象**。

## 封装

&emsp;&emsp;对于封装而言，一个对象它所封装的是自己的属性和方法，所以它是不需要以来其它对象就可以完成自己的操作。

- 封装的好处

1. 良好的封装能够减少耦合
2. 类内部的结构可以自由修改
3. 可以对成员进行更精确的控制
4. 隐藏信息，实现细节

&emsp;&emsp;首先我们先来看两个类:`Husband.java`和`Wife.java`

    public class Husband {
        /*
        * 对象的封装
        * 一个人的姓名、性别、年龄、妻子都是这个人的私有属性
        *
        */
        private String name;
        private String sex;
        private int age;
        private Wife wife;

        /*
        * setter()、getter()是该对象对外开发的接口
        */
        public String getName() {
            return name;
        }

        public void setName(String name) {
            this.name = name;
        }

        public String getSex() {
            return sex;
        }

        public void setSex(String sex) {
            this.sex = sex;
        }

        public int getAge() {
            return age;
        }

        public void setAge(int age) {
            this.age = age;
        }

        public void setWife(Wife wife) {
            this.wife = wife;
        }
    }

    public class Wife{
        private String name;
        private int age;
        private String sex;
        private Huaband husband;
        
        public String getName() {
            return name;
        }

        public void setName(String name) {
            this.name = name;
        }

        public String getSex() {
            return sex;
        }

        public void setSex(String sex) {
            this.sex = sex;
        }

        public void setAge(int age) {
            this.age = age;
        }

        public void setHusband(Husband husband) {
            this.husband = husband;
        }

        public Husband getHusband() {
            return husband;
        }
    }

&emsp;&emsp;从上面两个实例我们可以看出Husband里面wife引用是没有getter()的，同时wife的age也是没有getter()的。男人金屋藏娇，同时女人不愿意别人知道她的年龄。

&emsp;&emsp;所以封装把一个对象的属性私有化，同时提供一些可以被外界访问的属性的方法，如果不想被外界访问，我们大可不必提供方法给外界访问。但是如果一个类没有提供给外界访问的方法，那么这个类也没有什么意义了。比如我们将一个房子看作是一个对象，里面漂亮的装饰，如沙发、电视、空调、茶桌等等都是该房子的私有属性，但是如果我们没有那些墙遮挡，是不是别人就会一览无余呢？没有一点儿隐私，就是存在那个遮挡的墙，我们既能够有自己的隐私而且我们可以随意的更改里面的摆设而不会影响到其它。但是如果没有门窗，一个包裹的严严实实的黑盒子，又有什么存在的意义呢？所以通过门窗别人也能够看到里面的风景。所以说门窗就是房子对象留给外界的接口。

&emsp;&emsp;通过这个我们还不能真正体会封装的好处。现在我们从程序的角度来分析封装带来的好处。如果我们不使用封装，那么该对象就没有setter和getter，那么Husband类应该是这样：

    public class Husband {
        public String name;
        public String sex;
        public int age;
        public Wife wife;
    }

&emsp;&emsp;此时，我们应该这样来使用它：

    Husband husband = new Husband();
    husband.age = 30;
    husband.name = "tom";
    husband.sex = "male";

&emsp;&emsp;但是那天如果我们需要修改Husband，例如将age修改为String类型的呢？你只有一处使用好了这个类还好，如果你有几十个甚至上百个这个的地方，你是不是要改到崩溃。如果使用了封装，我们完全不需要做任何修改，只需要稍微改变下Husband的setAge方法即可。

    public class Husband {
        private String name;
        private String sex;
        // 改成string
        private String age;
        private Wife wife;

        public void setAge(int age) {
            this.age = String.valueOf(age);
        }
        ...
    }

&emsp;&emsp;其它的地方依然那样引用husband.setAge(22)保持不变

&emsp;&emsp;到了这里我们确实可以看出，**封装确实可以使我们容易地修改类地内部实现，而无需修改使用了该类地客户代码。**

&emsp;&emsp;我们再看这个好处：可以对成员变量进行更精确地控制。

&emsp;&emsp;还是那个Husband，一般来说我们在引用这个对象地时候是不容易出错地，但是有时你迷糊了，写成了这样：

    Husband husband = new Husband();
    husband.setAge(300);

&emsp;&emsp;也许是你粗心大意，毕竟谁见过300岁地人。使用封装就可以避免这个问题，我们对age地访问入口做一些控制：

    public void setAge(int age) {
        if (age > 120) {
            System.out.printlnn("ERROR:error age input...");
        } else {
            this.age = age;
        }
    }

&emsp;&emsp;上面都是对setter方法地控制，其实通过使用封装我们也能够对对象地出口作出很好地控制。例如性别我们在数据库中一般都是1、0方式来存储地，但是在前台我们又不能展示1、0，这里我们只需要在getter方法里做一些转换即可。

    public String getSexName() {
        if ("0".equals(sex)) {
            sexName = "女";
        } else if ("1".equals(sex)) {
            sexName = "男";
        } else {
            sexName = "其它";
        }
        return sexName;
    }

&emsp;&emsp;在使用地时候我们只需要使用sexName即可实现正确地性别显示。同理也可以用于针对不同地状态作出不同地操作。