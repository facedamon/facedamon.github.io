import java.lang.reflect.Array;

public class ArrayQueue<T> {
    private T[] mArray;
    private int count;

    public ArrayQueue(Class<T> type, int size) {
        mArray = (T[])Array.newInstance(type, size);
    }

    public void add(T val) {
        mArray[count++] = val;
    }

    public T front() {
        return mArray[0];
    }

    public T pop() {
        T ret = mArray[0];
        count--;
        for (int i = 1; i <= count; i++) {
            mArray[i-1] = mArray[i];
        }
        return ret;
    }

    public int size() {
        return count;
    }

    public boolean isEmpty(){
        return size() == 0;
    }

    public static void main(String[] args) {
        ArrayQueue<Integer> queue = new ArrayQueue(Integer.class, 12);
        queue.add(10);
        queue.add(20);
        queue.add(30);

        int tmp = queue.pop();
        System.out.printf("tmp=%d\n", tmp);

        tmp = queue.front();
        System.out.printf("tmp=%d\n", tmp);

        queue.add(40);

        while(!queue.isEmpty()) {
            System.out.printf("value=%d\n", queue.pop());
        }
    }
}