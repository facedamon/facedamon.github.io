
> 转载自 https://www.cnblogs.com/skywang12345/

# 概要
&emsp;&emsp;本节会现对栈的原理进行介绍，然后分别通过Java/Golang三种语言来演示栈的实现示例。

# 栈的介绍
&emsp;&emsp;栈是一种线性存储结构，它有以下几个特点：

1. 栈中数据是按照"先进后出 LIFO Last In First Out" 方式进出栈的。

2. 向栈中添加/删除数据时，只能从栈顶进行操作。

&emsp;&emsp;栈通常包括三种操作: push, peek, pop.

&emsp;&emsp;push -- 向栈中添加元素。

&emsp;&emsp;peek -- 返回栈顶元素。

&emsp;&emsp;push -- 返回并删除栈顶元素。

-  栈示意图

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/栈示意图.jpg)

&emsp;&emsp;栈中数据依次是30-->20-->10

-  出栈

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/出栈.jpg)

&emsp;&emsp;出栈前：栈顶元素是30.此时，栈中的元素依次是30-->30-->10
&emsp;&emsp;出栈后：30出栈后，栈顶元素变成20.此时，栈中的元素依次是20-->10

- 入栈

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/data-struct/linear-struct/入栈.jpg)

&emsp;&emsp;入栈前：栈顶元素是20.此时，栈中的元素依次是20-->10
&emsp;&emsp;入栈后：40入栈后，栈顶元素变成40.此时，栈中的元素依次是40-->20-->10

# Java实现

        import java.lang.reflect.Array;

        public class GeneralArrayStack<T> {
            private static final int DEFAULT_SIZE = 12;
            private T[] mArray;
            private int count;

            public GeneralArrayStack(Class<T> type){
                    this(type, DEFAULT_SIZE);
            }
            public GeneralArrayStack(Class<T> type, int size) {
                mArray = (T[])Array.newInstance(type, size);
                count = 0;
            }

            public void push(T val) {
                mArray[count++] = val;
            }

            public T peek() {
                return mArray[count-1];
            }

            public  T pop() {
                T ret = mArray[count-1];
                count--;
                return ret;
            }

            public int size () {
                return count;
            }

            public boolean isEmpty() {
                return size() == 0;
            }

            @Override
            public String toString() {
                if (isEmpty()) {
                    return "stack is empty\n";
                }
                StringBuffer sBuffer = new StringBuffer();
                int i = size() - 1;
                while (i >= 0) {
                    sBuffer.append(mArray[i]);
                    i--;
                }
                return sBuffer.toString;
            }
        }


# JDK自带


        public class StackTest {
            public static void main(String[] args) {
                Stack<Integer> stack = new Stack<Integer>();
                stack.push(10);
                stack.push(20);
                stack.push(30);

                int tmp = stack.pop();
                tmp = stack.peek();
                stack.push(40);

                while(!stack.empty()) {
                    tmp = stack.pop();
                    sout("tmp=%d\n", tmp);
                }
            }
        }


# Golang实现

- CopyLeft https://github.com/emirpasic/gods

        package arraystack

        import (
            "container/list"
            "github.com/facedamon/gods/lists/arraylist"
        )

        // Stack holds elements in an array-list
        type Stack struct {
            list *arraylist.List
        }

        // New instantiates a new empty stack
        func New() *Stack {
            return &Stack{list: arraylist.New()}
        }

        // Push adds a value onto the top of the stack
        func (stack *Stack) Push(value interface{}) {
            stack.list.Add(value)
        }

        // Pop removes top element on stack and returns it, or nil if stack is empty.
        // Second return parameter is true, unless the stack was empty and there was nothing to top.
        func (stack *Stack) Pop() (value interface{}, ok bool) {
            value, ok = stack.list.Get(stack.list.Size() - 1)
            stack.list.Remove(stack.list.Size() - 1)
            return
        }

        // Peek returns top element on the stack without removing it, or nil if stack is empty
        // Second return parameter is true, unless the stack was empty and there was nothing to peek.
        func (stack *Stack) Peek() (value interface{}, ok bool) {
            return stack.list.Get(stack.list.Size() - 1)
        }

        // Empty returns true if stack does not contain any element
        func (stack *Stack) Empty() bool {
            return stack.list.Empty()
        }

        // Size returns number of elements within the stack.
        func (stack *Stack) Size() int {
            return stack.list.Size()
        }

        // Clear removes all elements from the stack
        func (stack *Stack) Clear() {
            stack.list.Clear()
        }

        func (stack *Stack) Values() []interface{} {
            size := stack.list.Size()
            elements := make([]interface]{}, size, size)
            for i := 1; i <= size; i++ {
                elements[size - i], _ = stack.list.Get(i - 1)
            }
            return elements
        }

        // String returns a string representation of container
        func (stack *Stack) String() string {
            str := "ArrayStack\n"
            values := []string{}
            for _, value := range stack.list.Values() {
                values = append(values, fmt.Sprintf("%v", value))
            }
            str += strings.Join(values, ", ")
            return str
        }

        // Check that the index is within dounds of the list
        func (stack *Stack) withinRange(index int) bool {
            return index >= 0 && index < stack.list.Size()
        }
