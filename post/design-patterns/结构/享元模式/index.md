
# 摘要

&emsp;&emsp;享元模式（Flyweight Pattern）主要用于减少创建对象的数量，以减少内存占用和提高性能。这种类型的设计模式属于结构型模式，它提供了减少对象数量从而改善应用所需的对象结构的方式。

&emsp;&emsp;享元模式尝试重用现有的同类对象，如果未找到匹配的对象，则创建新对象。我们将通过创建 5 个对象来画出 20 个分布于不同位置的圆来演示这种模式。由于只有 5 种可用的颜色，所以 color 属性被用来检查现有的 Circle 对象。

## 介绍

- 意图：运用共享技术有效地支持大量细粒度的对象。
- 主要解决：在有大量对象时，有可能会造成内存溢出，我们把其中共同的部分抽象出来，如果有相同的业务请求，直接返回在内存中已有的对象，避免重新创建。
- 何时使用： 1、系统中有大量对象。 2、这些对象消耗大量内存。 3、这些对象的状态大部分可以外部化。 4、这些对象可以按照内蕴状态分为很多组，当把外蕴对象从对象中剔除出来时，每一组对象都可以用一个对象来代替。 5、系统不依赖于这些对象身份，这些对象是不可分辨的。
- 如何解决：用唯一标识码判断，如果在内存中有，则返回这个唯一标识码所标识的对象。
- 关键代码：用 HashMap 存储这些对象。
- 应用实例： 1、JAVA 中的 String，如果有则返回，如果没有则创建一个字符串保存在字符串缓存池里面。 2、数据库的数据池。

- 优点：大大减少对象的创建，降低系统的内存，使效率提高。
- 缺点：提高了系统的复杂度，需要分离出外部状态和内部状态，而且外部状态具有固有化的性质，不应该随着内部状态的变化而变化，否则会造成系统的混乱。
- 使用场景： 1、系统有大量相似对象。 2、需要缓冲池的场景。
- **注意事项**： 1、注意划分外部状态和内部状态，否则可能会引起线程安全问题。 2、这些类必须有一个工厂对象加以控制。

## 实现

&emsp;&emsp;我们将创建一个 Shape 接口和实现了 Shape 接口的实体类 Circle。下一步是定义工厂类 ShapeFactory。

&emsp;&emsp;ShapeFactory 有一个 Circle 的 HashMap，其中键名为 Circle 对象的颜色。无论何时接收到请求，都会创建一个特定颜色的圆。ShapeFactory 检查它的 HashMap 中的 circle 对象，如果找到 Circle 对象，则返回该对象，否则将创建一个存储在 hashmap 中以备后续使用的新对象，并把该对象返回到客户端。

&emsp;&emsp;FlyWeightPatternDemo 类使用 ShapeFactory 来获取 Shape 对象。它将向 ShapeFactory 传递信息（red / green / blue/ black / white 内部状态），以便获取它所需对象的颜色。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/design-pattern/struct/1616398933.jpg)

        public interface Shape {
            void draw();
        }

        public enum Color {
            RED,
            Green,
            Blue,
            White,
            Black
        }

        public class Circle implements Shape {

            private Color color;
            @Setter
            private int x;
            @Setter
            private int y;
            @Setter
            private int radius;

            public Circle(Color color) {
                this.color = color;
            }

            @Override
            public void draw() {
                System.out.println(
                        String.format("Circle:Draw() [Color : %s, x: %d, y: %d, radius: %d",
                                color, x, y, radius));
            }
        }

        public class ShapeFactory {
            private static final HashMap<Color, Shape> circleMap = new HashMap<>(8);
            public static Shape getCircle(Color color) {
                Shape circle = circleMap.get(color);
                if (circle == null) {
                    circle = new Circle(color);
                    circleMap.put(color, circle);
                    System.out.println("Creating circle of color: " + color);
                }
                return circle;
            }
        }

        public class FlyweightPatternDemo {
            private Color getRandomColor() {
                return Color.values()[(int)(Math.random() * Color.values().length)];
            }
            private int getRandomX() {
                return (int)(Math.random() * 100);
            }
            private int getRandomY() {
                return (int)(Math.random() * 100);
            }

            @Test
            public void testFlyweight() {
                for (int i = 0; i < 20; i++) {
                    Circle circle = (Circle) ShapeFactory.getCircle(getRandomColor());
                    circle.setX(getRandomX());
                    circle.setY(getRandomY());
                    circle.setRadius(100);
                    circle.draw();
                }
            }
        }

        //output:
        Creating circle of color: Blue
        Circle:Draw() [Color : Blue, x: 90, y: 20, radius: 100
        Creating circle of color: White
        Circle:Draw() [Color : White, x: 36, y: 33, radius: 100
        Creating circle of color: Green
        Circle:Draw() [Color : Green, x: 94, y: 66, radius: 100
        Creating circle of color: RED
        Circle:Draw() [Color : RED, x: 66, y: 15, radius: 100
        Circle:Draw() [Color : RED, x: 63, y: 56, radius: 100
        Circle:Draw() [Color : RED, x: 5, y: 74, radius: 100
        Creating circle of color: Black
        Circle:Draw() [Color : Black, x: 36, y: 85, radius: 100
        Circle:Draw() [Color : Black, x: 70, y: 2, radius: 100
        Circle:Draw() [Color : Green, x: 65, y: 54, radius: 100
        Circle:Draw() [Color : White, x: 9, y: 34, radius: 100
        Circle:Draw() [Color : Green, x: 61, y: 34, radius: 100
        Circle:Draw() [Color : RED, x: 65, y: 64, radius: 100
        Circle:Draw() [Color : White, x: 62, y: 18, radius: 100
        Circle:Draw() [Color : Black, x: 71, y: 69, radius: 100
        Circle:Draw() [Color : RED, x: 60, y: 49, radius: 100
        Circle:Draw() [Color : White, x: 33, y: 54, radius: 100
        Circle:Draw() [Color : Blue, x: 60, y: 90, radius: 100
        Circle:Draw() [Color : RED, x: 96, y: 12, radius: 100
        Circle:Draw() [Color : RED, x: 41, y: 72, radius: 100
        Circle:Draw() [Color : White, x: 38, y: 80, radius: 100