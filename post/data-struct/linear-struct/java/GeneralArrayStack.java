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