
# 摘要

&emsp;&emsp;在代理模式（Proxy Pattern）中，一个类代表另一个类的功能。这种类型的设计模式属于结构型模式。在代理模式中，我们创建具有现有对象的对象，以便向外界提供功能接口。

## 介绍

- 意图：为其他对象提供一种代理以控制对这个对象的访问。
- 主要解决：在直接访问对象时带来的问题，比如说：要访问的对象在远程的机器上。在面向对象系统中，有些对象由于某些原因（比如对象创建开销很大，或者某些操作需要安全控制，或者需要进程外的访问），直接访问会给使用者或者系统结构带来很多麻烦，我们可以在访问此对象时加上一个对此对象的访问层。
- 何时使用：想在访问一个类时做一些控制。
- 如何解决：增加中间层。
- 关键代码：实现与被代理类组合。
- 应用实例： 1、Windows 里面的快捷方式。 2、猪八戒去找高翠兰结果是孙悟空变的，可以这样理解：把高翠兰的外貌抽象出来，高翠兰本人和孙悟空都实现了这个接口，猪八戒访问高翠兰的时候看不出来这个是孙悟空，所以说孙悟空是高翠兰代理类。 3、买火车票不一定在火车站买，也可以去代售点。 4、一张支票或银行存单是账户中资金的代理。支票在市场交易中用来代替现金，并提供对签发人账号上资金的控制。 5、spring aop。
- 优点： 1、职责清晰。 2、高扩展性。 3、智能化。
- 缺点： 1、由于在客户端和真实主题之间增加了代理对象，因此有些类型的代理模式可能会造成请求的处理速度变慢。 2、实现代理模式需要额外的工作，有些代理模式的实现非常复杂。
- 使用场景：按职责来划分，通常有以下使用场景： 1、远程代理。 2、虚拟代理。 3、Copy-on-Write 代理。 4、保护（Protect or Access）代理。 5、Cache代理。 6、防火墙（Firewall）代理。 7、同步化（Synchronization）代理。 8、智能引用（Smart Reference）代理。
- **注意事项**： 1、和适配器模式的区别：适配器模式主要改变所考虑对象的接口，而代理模式不能改变所代理类的接口。 2、和装饰器模式的区别：装饰器模式为了增强功能，而代理模式是为了加以控制。

## 实现

&emsp;&emsp;我们将创建一个 Image 接口和实现了 Image 接口的实体类。ProxyImage 是一个代理类，减少 RealImage 对象加载的内存占用。ProxyPatternDemo 类使用 ProxyImage 来获取要加载的 Image 对象，并按照需求进行显示。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/design-pattern/struct/1616210495.png)

        public interface Image {
            void display();
        }

        public class RealImage implements Image{

            private String fileName;

            public RealImage(String fileName) {
                this.fileName = fileName;
                loadFromDisk(fileName);
            }

            private void loadFromDisk(String fileName) {
                System.out.println("Loading " + fileName);
            }

            @Override
            public void display() {
                System.out.println("Displaying " + fileName);
            }
        }

        public class ProxyImage implements Image{

            private RealImage image;
            private String fileName;

            public ProxyImage(String fileName) {
                this.fileName = fileName;
            }

            @Override
            public void display() {
                if (image == null) {
                    image = new RealImage(fileName);
                }
                image.display();
            }
        }

        public class App {
            public static void main(String[] args) {
                Image image = new ProxyImage("test");
                image.display();
            }
        }

## JDK Proxy

&emsp;&emsp;核心api包含两个重要的接口和类：一个是`InvocationHandler`，另一个是`Proxy`。`static Object newProxyInstance(ClassLoader loader, Class[] interfaces, InvocationHandler h)`该方法用于为指定类装载器、一组接口及调用处理器生成动态代理类实例。

- loader：制定代理类的ClassLoader加载器
- interface：指定代理类要实现的所有接口
- h：表示的是当前这个动态代理对象在调用方法的时候，会关联到哪一个InvocationHandler对象上

&emsp;&emsp;`InvocationHandler`负责连接代理类和委托类的中间类必须实现的接口。它自定义了一个invoke方法，用于集中处理在动态代理对象上的方法调用，通常在该方法中实现对委托类的代理访问。


        public class Target implements InvocationHandler {

            private Object delegate;
            public Target(Object delegate) {
                this.delegate = delegate;
            }

            @Override
            public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
                System.out.println("------before------");
                Object ret = method.invoke(delegate, args);
                System.out.println("------after------");
                return ret;
            }
        }

        public class App {
            static class A implements Appendable{
                private StringBuilder builder;
                public A() {
                    builder = new StringBuilder();
                }
                @Override
                public Appendable append(CharSequence csq) throws IOException {
                    System.out.println("------appending------");
                    return builder.append(csq);
                }

                @Override
                public Appendable append(CharSequence csq, int start, int end) throws IOException {
                    System.out.println("------appending------");
                    return builder.append(csq, start, end);
                }

                @Override
                public Appendable append(char c) throws IOException {
                    System.out.println("------appending------");
                    return builder.append(c);
                }
                @Override
                public String toString() {
                    return builder.toString();
                }
            }
            public static void main(String[] args) {
                //创建被代理对象
                Appendable a = new A();
                //创建代理器
                InvocationHandler handler = new Target(a);
                //生成代理类
                Appendable appendable = (Appendable) Proxy.newProxyInstance(A.class.getClassLoader(), A.class.getInterfaces(), handler);
                try {
                    appendable = appendable.append('a');
                    System.out.println(appendable.toString());
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }
        //output:
        ------before------
        ------appending------
        ------after------
        a

## CGLIB

- 原理：动态生成一个要代理类的子类，子类重写要代理的类的所有不是final的方法。在子类中采用方法拦截的技术拦截所有父类方法的调用，顺势织入横切逻辑。它比使用java反射的JDK动态代理要快。
- 底层：使用字节码处理框架ASM，来转换字节码并生成新的类。不鼓励直接使用ASM，因为它要求你必须对JVM内部结构包括class文件的格式和指令集都很熟悉。
- 缺点：对于final方法，无法进行代理。

> 为什么使用CGLIB

&emsp;&emsp;CGLIB代理主要通过对字节码的操作，为对象引入间接级别，以控制对象的访问。我们知道Java中有一个动态代理也是做这个事情的，那我们为什么不直接使用Java动态代理，而要使用CGLIB呢？答案是CGLIB相比于JDK动态代理更加强大，JDK动态代理虽然简单易用，但是其有一个致命缺陷是，只能对接口进行代理。如果要代理的类为一个普通类、没有接口，那么Java动态代理就没法使用了。

        // 代理目标类
        public class TargetObject {
            public String method1(String paramName) {
                return paramName;
            }

            public int method2(int count) {
                return count;
            }

            public int method3(int count) {
                return count;
            }

            @Override
            public String toString() {
                return "TargetObject []"+ getClass();
            }
        }

        // 目标对象拦截器
        public class TargetInterceptor implements MethodInterceptor {
            @Override
            public Object intercept(Object o, Method method, Object[] objects, MethodProxy methodProxy) throws Throwable {
                System.out.println("调用前");
                Object result = methodProxy.invokeSuper(o, objects);
                System.out.println(" 调用后");
                return result;
            }
        }

        @Test
        public void testCGLIB() {
            Enhancer enhancer = new Enhancer();
            //设置父类
            enhancer.setSuperclass(TargetObject.class);
            //设置拦截器
            enhancer.setCallback(new TargetInterceptor());
            //生成代理类并向上转型
            TargetObject targetObject = (TargetObject) enhancer.create();
            System.out.println(targetObject);
            System.out.println(targetObject.method1("mmm"));
            System.out.println(targetObject.method2(100));
            System.out.println(targetObject.method3(200));
        }

        //output:
        调用前
        调用后
        TargetObject []class org.facedamon.TargetObject$$EnhancerByCGLIB$$d0737a56
        调用前
        调用后
        mmm
        调用前
        调用后
        100
        调用前
        调用后
        200

### 回调过滤器CallbackFilter

&emsp;&emsp;在CGLib回调时可以设置对不同方法执行不同的回调逻辑，或者根本不执行回调。在JDK动态代理中并没有类似的功能，对InvocationHandler接口方法的调用对代理类内的所以方法都有效。

        public class TargetMethodCallbackFilter implements CallbackFilter {
            /**
            * @param method 过滤方法
            * @return 数字，表示callback数组中的索引位置，要用到的callback
            */
            @Override
            public int accept(Method method) {
                if (method.getName().equals("method1")) {
                    System.out.println("filter method1 == 0");
                    return 0;
                }
                if (method.getName().equals("method2")) {
                    System.out.println("filter method2 == 1");
                    return 1;
                }
                if (method.getName().equals("method3")) {
                    System.out.println("filter method3 == 2");
                    return 2;
                }
                return 0;
            }
        }

        /**
        * 表示锁定方法返回值，无论被代理类的方法返回什么值，回调方法都返回固定值
        */
        public class TargetResultFixed implements FixedValue {
            @Override
            public Object loadObject() throws Exception {
                System.out.println("锁定结果");
                Object obj = 999;
                return obj;
            }
        }

        @Test
        public void testCGLIBCallbackFilter() {
            Enhancer enhancer = new Enhancer();
            enhancer.setSuperclass(TargetObject.class);
            CallbackFilter callbackFilter = new TargetMethodCallbackFilter();
            //no operator 什么操作也不做，代理类直接调用被代理的方法不进行拦截
            Callback noopCb = NoOp.INSTANCE;
            //方法拦截器
            Callback callback1 = new TargetInterceptor();
            //锁定方法返回值
            Callback fixedValue = new TargetResultFixed();
            Callback[] cbArray = new Callback[]{callback1, noopCb, fixedValue};
            enhancer.setCallbacks(cbArray);
            enhancer.setCallbackFilter(callbackFilter);
            TargetObject targetObject = (TargetObject) enhancer.create();
            //这里的method调用就会触发CallbackFilter的accept方法，根据返回值调用cbArray
            System.out.println(targetObject.method1("mmm"));
            System.out.println(targetObject.method2(100));
            System.out.println(targetObject.method3(100));
            System.out.println(targetObject.method3(200));
        }

        //output:
        filter method1 == 0
        filter method2 == 1
        filter method3 == 2
        调用前
        调用后
        mmm
        100
        锁定结果
        999
        锁定结果
        999

## Javassist