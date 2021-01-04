
# 摘要

&emsp;&emsp;抽象工厂模式(Abstract Factory Pattern)用于生成产品族的工厂。该工厂又称为其它工厂的工厂。在抽象工厂模式中，接口是负责创建一个相关对象的工厂，不需要显式指定它们的类。每个生成的工厂都能按照工厂模式提供对象。

## 介绍

- **意图**： 提供一个创建一系列相关或相互依赖对象的接口，而无需指定它们的具体的类。
- **主要解决**: 主要解决接口选择问题。
- **何时使用**: 系统的产品有多于一个的产品族，而系统只消费其中某一族的产品。
- **如何解决**: 在一个产品族里，定义多个产品。
- **关键代码**: 在一个工厂里聚合多个同类产品。
- **应用实例**: 

&emsp;&emsp;工作了，为了参加一些聚会，肯定有两套或多套衣服吧，比如说有商务装（成套，一系列具体产品）、时尚装（成套，一系列具体产品），甚至对于一个家庭来说，可能有商务女装、商务男装、时尚女装、时尚男装，这些也都是成套的，即一系列具体产品。假设一种情况（现实中是不存在的，但有利于说明抽象工厂模式），在您的家中，某一个衣柜（具体工厂）只能存放某一种这样的衣服（成套，一系列具体产品），每次拿这种成套的衣服时也自然要从这个衣柜中取出了。用 OOP 的思想去理解，所有的衣柜（具体工厂）都是衣柜类的（抽象工厂）某一个，而每一件成套的衣服又包括具体的上衣（某一具体产品），裤子（某一具体产品），这些具体的上衣其实也都是上衣（抽象产品），具体的裤子也都是裤子（另一个抽象产品）。

## 实现

&emsp;&emsp;我们将创建`Shape`和`Color`接口和实现这些接口的实体类。下一步是创建抽象工厂类`AbstractFactory`。接着定义工厂类`ShapeFactory`和`ColorFactory`，这两个工厂类都是扩展了`AbstractFactory`。然后创建一个工厂生成器`FactoryProducer`，演示类使用`FactoryProducer`来获取`AbstractFactory`对象，它将向`AbstractFactory`传递形状信息Shape，以便获取它所需对象的类型。同时它还向`AbstractFactory`传递颜色信息Color，以便获取它所需对象的类型。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/design-pattern/create/3E13CDD1-2CD2-4C66-BD33-DECBF172AE03.jpg)


### 定义接口

        public interface Shape {
            void draw();
        }

        public interface Color {
            void fill();
        }

### 创建接口的实现

        public class Rectangle implements Shape {

            @Override
            public void draw() {
                System.out.println("Inside Rectangle::draw() method.")
            }
        }

        public class Square implements Shape {

            @Override
            public void draw() {
                System.out.println("Inside Square::draw() method.")
            }
        }

        public class Circle implements Shape {

            @Override
            public void draw() {
                System.out.println("Inside Circle::draw() method.")
            }
        }

        public class Red implements Color {

            @Override
            public void fill() {
                System.out.println("Inside Red::fill() method.")
            }
        }

        public class Green implements Color {

            @Override
            public void fill() {
                System.out.println("Inside Green::fill() method.")
            }
        }

        public class Blue implements Color {

            @Override
            public void fill() {
                System.out.println("Inside Blue::fill() method.")
            }
        }

### 定义抽象工厂

        public abstract class AbstractFactory {
            public abstract Color getColor(String color);
            public abstract Shape getShape(String shape);
        }

### 定义工厂类

        public class ShapeFactory extends AbstractFactory {
            
            @Override
            public Shape getShape(String shapeType) {
                if (null == shapeType) {
                    return null;
                }
                if (shapeType.equalsIgnoreCase("CIRCLE")) {
                    return new Circle();
                }
                if (shapeType.equalsIgnoreCase("RECTANGLE")) {
                    return new Rectangle();
                }
                if (shapeType.equalsIgnoreCase("SQUARE")) {
                    return new Square();
                }
                return null;
            }

            @Override
            public Color getColor(String color) {
                return null;
            }
        }

        public class ColorFactory extends AbstractFactory {

            @Override
            public Shape getShape(String shapeType) {
                return null;
            }

            @Override
            public Coloe getColor(String color) {
                if (null == color) {
                    return null;
                }
                if (color.equalsIgnoreCase("RED")) {
                    return new Red();
                }
                if (color.equalsIgnoreCase("GREEN")) {
                    return new Green();
                }
                if (color.equalsIgnoreCase("BLUE")) {
                    return new Blue();
                }
                return null;
            }
        }

### 定义工厂生成器(关联产品族)

        public class FactoryProducer {
            public statuc AbstractFactory getFactory(String choice) {
                if (choice.equalsIgnoreCase("SHAPE")) {
                    return new ShapeFactory();
                }
                if (choice.equalsIgnoreCase("COLOR")) {
                    return new ColorFactory();
                }
                return null;
            }
        }

### 使用(传递类型信息获取实体类)

        public class AbstractFactoryPatternDemo {
            public static void main(String[] args) {
                AbstractFactory shapeFactory = FactoryProducer.getFactory("SHAPE");
                Shape shape1 = shapeFactory.getShape("CIRCLE");
                shape1.draw();

                AbstractFactory colorFactory = FactoryProducer.getFactory("COLOR");
                Color color1 = colorFactory.getColor("RED");
                color1.fill();
            }
        }

## Golang 版本

&emsp;&emsp;将上述行为使用Golang版本实现一边，因为Golang没有继承，我们采用组合的方式。

### 定义接口

        package data

        import "fmt"

        type Shape interface {
            Draw() error
        }

        type Color interface {
            Fill() error
        }

### 创建接口的实现

        type ShapeType int
        type ColorType int

        const (
            Rectangle ShapeType = 1 << iota
            Square ShapeType
            Circle ShapeType

            Red ColorType
            Green ColorType
            Blue ColorType
        )

        func newRectangle() Shape {
            // ...
        }
        func newSquare() Shape {
            // ...
        }
        func newCirclr() Shape {
            // ...
        }

        func newRed() Color {
            // ...
        }
        func newGreen() Color() {
            // ...
        }
        func newBlue() Color() {
            // ...
        }

### 定义抽象工厂

        type AbstractFactory interface {
            getColor(ColorType color) Color
            getShape(ShapeType shape) Shape
        }

### 定义工厂类

        type ShapeFactory struct{}
        func (*ShapeFactory) getShape(i ShapeType) Shape {
            switch i {
                case Rectangle:
                    return newRectangle()
                case Square:
                    return newSquare()
                default:
                    return newCircle()
            }
        }
        func (*ShapeFactory) getColor(i ColorType) Color {
            return nil
        }

        type ColorFactory struct{}
        func (*ColorFactory) getColor(i ColorType) Color {
            switch i {
                case Red:
                    return newRed()
                case Green:
                    return newGreen()
                default:
                    return newBlue()
            }
        }
        func (*ColorFactory) getShape(i ShapeType) Shape {
            return nil
        }

### 定义工厂生成器(关联产品族)

        type choice int
        const (
            ShapeChoice choice = 1 << iota
            ColorChoice choice
        )

        func NewFactory(i choice) AbstractFactory {
            switch i {
                case ShapeChoice:
                    return ShapeFactory{}
                default ColorChoice:
                    return ColorFactory{}
            }
        }

## 优点

&emsp;&emsp;当一个产品族中的多个对象被设计成一起工作时，它能保证客户端始终只使用同一个产品族中的对象。

## 缺点

&emsp;&emsp;产品族扩展非常困难，要增加一个系列的某一产品，既要在抽象的Creator里加代码，又要在具体的里面加代码。