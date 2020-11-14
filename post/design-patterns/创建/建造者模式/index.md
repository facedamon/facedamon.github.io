
# 摘要

&emsp;&emsp;Builder Pattern使用多个简单的对象一步一步构建一个复杂的对象。

# 介绍

- **意图**：将一个复杂的构建与其表示相分离，使得同样的构建过程可以创建不同的表示。
- **主要解决**：主要解决在软件系统中，有时候面临着"一个复杂对象"的创建工作，其通常由各个不分的子对象用一定的算法构成；由于需求的变化，这个复杂对象的各个部分疆场面临着剧烈的变化，但是将它们组合在一起的算法却相对稳定。
- **关键代码**：`建造者`：创建和提供实例，`导演`：管理建造出来的比例的依赖关系。
- **应用实例**：去肯德基、汉堡、可乐、薯条、炸鸡翅等是不变的，而其组合是经常变化的，产出很多套餐
- **优点**：建造者独立、易扩展。便于控制细节风险
- **缺点**：产品必须有共同点，范围有限制。如内部变化复杂，会有很多的建造类
- **使用场景**：需要生成的对象具有复杂的内部结构；需要生成的对象内部属性本身相互依赖


# 实现

&emsp;&emsp;我们假设一个快餐店的商业案例，其中，一个典型的套餐可以是汉堡(Burger)和一杯冷饮(COld drink)。汉堡可以是素食汉堡(Veg Burger)或鸡肉汉堡(Chicken Burger)，它们是包在纸盒中。冷饮可以是可口可乐(coke)或百事可乐(pspsi)，它们是装载瓶子中。

&emsp;&emsp;我们将创建一个表示食物条目的item接口和实现item接口的实体类，以及表示食物包装的Packing接口和实现Packing接口的实体类，汉堡是包在纸盒中，冷饮是装在瓶子中。

&emsp;&emsp;创建一个Meal套餐类，表示item食物条目的集合，创建一个MealBuilder建造者，通过不同的item菜单来构建不同的meal套餐。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/markdownps2@master/design-pattern/create/20201015-builder-pattern.svg)

## 创建一个表示食物条目和食物包装的接口

    package interface;

    public interface Item {
        String name();
        Packing packing();
        float price();
    }

    public interface Packing {
        String pack();
    }

## 创建实现Packing的实体类

    package impl;

    public class Wrapper implements Packing {
        @Override
        public String packing() {
            return "Wrapper";
        }
    }

    public class Bottle implements Packing {
        @Override
        public String packing() {
            return "Bottle";
        }
    }

## 创建实现Item条目的抽象类

    public abstract class Burger implements Item {
        @Override
        public Packing packing() {
            return new Wrapper();
        }

        @Override
        public abstract float price();
    }

    public abstract class ColdDrink implements Item {
        @Override
        public Packing packing() {
            return new Bottle();
        }

        @Override
        public abstract float price();
    }

## 创建Burger和ColdDrink的具体实现

    package impl;

    public class VegBurger extends Burger {
        @Override
        public float price {
            return 25.0f;
        }

        @Override
        public String name() {
            return "Veg Burger";
        }
    }

    public class ChickenBueger extends Burger {
        @Override
        public float price() {
            return 50.5f;
        }

        @Override
        public String name() {
            return "Chicken Burger";
        }
    }

    public class Coke extends ColdDrink {
        @Override
        public float price() {
            return 30.0f;
        }

        @Override
        public String name() {
            return "Coke";
        }
    }

    public class Pepsi exntends ColdDrink {
        @Override
        public float price() {
            return 35.0f;
        }

        @Override
        public String name() {
            return "Pepsi";
        }
    }

## 创建Meal套餐

    package default;

    import java.util.ArrayList;
    import java.util.List;

    public class Meal {
        private List<Item> items = new ArrayList<Item>();

        public void addItem(Item item) {
            items.add(item);
        }

        public float getCost() {
            float cost = 0.0f;
            for (Item item : items) {
                cost += item.price();
            }
            return cost;
        }

        public void showItems() {
            for (Item item: items) {
                sout(String.format("Item: %s, Packing: %s, Price: %f", item.name(), item.packing().pack(), item.price()));
            }
        }
    }

## 创建一个MealBuilder构建者，负责创建Meal套餐

    package default;

    public class MealBuilder {
        public Meal prepareVegMeal() {
            Meal meal = new Meal();
            meal.addItem(new VegBurger());
            meal.addItem(new Coke());
            return meal;
        }

        public Meal prepareNonVegMeal() {
            Meal meal = new Meal();
            meal.addItem(new ChickenBurger);
            meal.addItem(new Pepsi());
            return meal;
        }
    }

## 创建demo

    package demo;

    public class BuilderPatternDemo {
        public static void main(String[] args) {
            MealBuilder mealBuilder = new MealBuilder();
            Meal vegMeal = mealBuilder.prepareVegMeal();
            sout("Veg Meal: %s, total cost: %f", vegMeal.showItems(), vegMeal.getCost());

            Meal nonVegMeal = mealBuilder.prepareNonVegMeal();
            sout("non-veg meal: %s, total cost: %f", nonVegMeal.showItems(), nonVegMeal.getCost());
        }
    }

## 输出

&emsp;&emsp;`Veg Meal`

&emsp;&emsp;`Item: Veg Burger, Packing: Wrapper, Price: 25.0`

&emsp;&emsp;`Item: Coke, Packing: Bottle, Price: 30.0`

&emsp;&emsp;`total cost: 55.0`

</br>

&emsp;&emsp;`non-veg meal`

&emsp;&emsp;`Item: Chicken Burger, Packing: Wrapper, Price: 50.5`

&emsp;&emsp;`Item: Packing: Bottle, Price: 35.0`

&emsp;&emsp;`total cost: 85.5`

# Golang版本

    package car

    type Speed float64

    const(
        MPH Speed = 1
        KPH       = 1.60934
    )

    type Color string

    const (
        BlueColor Color = "blue"
        GreenColor      = "green"
        RedColor        = "red"
    )

    type Wheels string

    const (
        SportsWHeels Wheels = "sports"
        SteelWheels         = "steel"
    )

    type Builder interface {
        Color(Color) Builder
        Wheels(Wheels) Builder
        TopSpeed(Speed) Builder
        Build() Interface
    }

    type Interface interface {
        Drive() error
        Stop() error
    }

    type bu struct {
         c Color
         w Wheels
         t Speed
    }

    func NewBuilder() *bu {
            return &bu{}
    }

    func (b *bu) Color(c Color) Builder {
            b.c = c
            return b
    }

    func (b *bu) Wheels(w Wheels) Builder {
            b.w = w
            return b
    }

    func (b *bu) TopSpeed(t Speed) Builder {
            b.t = t
            return b
    }

    func (b *bu) Build() Interface {
            return b
    }

    func (b *bu) Drive() error {
            fmt.Printf("%s\n", "i am driving")
            return nil
    }

    func (b *bu) Stop() error {
            return nil
    }

    func (b *bu) Paint(c Color) Builder {
            b.c = c
            return b
    }

## 使用

    assembly := car.NewBuilder().Paint(car.RedColor)

    familyCar := assembly.Wheels(car.SportsWheels).TopSpeed(50 * car.MPH).Build()
    familyCar.Drive()

    sportsCar := assembly.Wheels(car.SteelWheels).TopSpeed(150 * car.MPH).Build()
    sportsCar.Drive()