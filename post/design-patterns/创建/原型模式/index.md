
# 摘要

&emsp;&emsp;原型模式(Propertype Pattern)是用于创建重复的对象，同时又能保证性能。

&emsp;&emsp;这种模式是实现了一个原型接口，`该接口用于创建当前对象的克隆。当直接创建对象的代价比较大时，则采用这种模式`。例如，一个对象需要在一个高代价的数据库操作之后被创建。我们可以缓存该对象，在下一个请求时返回它的克隆，在需要的时候更新数据库，以此来减少数据库调用。

# 介绍

- **意图**：用原型实例指定创建对象的种类，并且通过拷贝这些原型创建新的对象。
- **主要解决**：在运行期建立和删除原型
- **何时使用**：当一个系统应该独立于它的产品创建，构成和表示时
- **如何解决**：利用已有的一个原型对象，快速地生成和原型对象一样的实例
- **使用场景**：在实际项目中，原型模式很少单独出现，一般是和工厂方法模式一起出现，通过clone的方法创建一个对象，然后由工厂方法提供给调用者。原型模式已经与Java融为一体。
- **注意事项**：与通过对一个类进行实例化来构造新对象不同的是，原型模式是通过拷贝一个现有对象生成新对象。浅拷贝实现Cloneable，重写，深拷贝是通过实现Serializable读取二进制流

# 实现

&emsp;&emsp;我们将创建一个抽象类Shape和扩展了Shape类的实体类。下一步是定义类ShapeCache，该类把shape对象存储在一个Hashtable中，并在请求的时候返回它们的克隆。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/design-pattern/create/20201015-prototype-pattern.svg)


## 创建一个实现了Cloneable接口的抽象类

    public abstract class Shape implements Cloneable {
        private String id;
        protected String type;

        abstract void draw();

        public String getType(){
            return type;
        }

        public String getId() {
            return id;
        }

        public void setId(String id) {
            this.id = id;
        }

        public Object clone() {
            Object clone = null;
            try {
                clone = super.clone();
            } catch(CloneNotSupportedException e) {
                e.printStackTrace();
            }
            return clone;
        }
    }

## 创建扩展了抽象类的实体类

    public class Rectangle extends Shape {
        public Rectangle() {
            type = "Rectangle";
        }

        @Override
        public void draw() {
            sout("inside Rectangle:draw() method.");
        }
    }

    public class Square extends Shape {
        public Square() {
            type = "Square";
        }

        @Override
        public void draw() {
            sout("inside Square:draw() method.")
        }
    }

    public class Circle extends Shape {
        public Circle() {
            type = "Circle";
        }

        @Override
        public void draw() {
            sout("Inside Circle:draw() method.");
        }
    }

## 创建Cache，模拟从数据库获取实体类，并存储在Hashtable

    public class ShapeCache {
        private static Hashtable<String, Shape> shapeMap = new Hashtable<String, Shape>();

        public static Shape getShape(String shapeId) {
            Shape cachedShape = shapeMap.get(shapeId);
            return (Shape) cachedShape.clone();
        }

        // 模拟加载数据库查询
        public static void loadCache() {
            Circle circle = new Circle();
            circle.setId("1");
            shapeMap.put(circle.getId(), circle);

            Square square = new Square();
            square.setId("2");
            shapeMap.put(square.getId, square);

            Rectangle rectangle = new Rectangle();
            rectangle.setId("3");
            shapeMap.put(rectangle.getId(), rectangle);
        }
    }

## 获取克隆

    public class PrototyPatternDemo {
        public static void main(String[] args) {
            ShapeCache.loadCache();

            Shape clonedShape = (Shape) ShapeCache.getShape("1");
            sout("Shape: " + clonedShape.getType());

            Shape cloedShape2 = (Shape) ShapeCache.getShape("2");
            sout("Shape: " + cloedShape2.getType());

            Shape cloedShape3 = (Shape) ShapeCache.getShape("3");
            sout("Shape: " + cloedShape3.getType());
        }
    }

> Shape: Circle
</br>
> Shape: Square
</br>
> Shape: Rectangle

# Golang版本

    package prototype

    // cloneable是原型对象需要实现的接口
    type Cloneable interface {
        Clone() Cloneable
    }

    type PrototypeManager struct {
        prototypes map[string]Cloneable
    }

    func NewPrototypeManager() *PrototypeManager {
        return &PrototypeManager{
            prototypes: make(map[string]Coleable),
        }
    }

    func (p *PrototypeManager) Get(name string) Cloneable {
        return p.prototypes[name]
    }

    func (p *PrototypeManager) Set(name string, prototype Cloneable) {
        p.prototypes[name] = prototype
    } 


</br>

    package prototype

    var manager *PrototypeManager

    type 