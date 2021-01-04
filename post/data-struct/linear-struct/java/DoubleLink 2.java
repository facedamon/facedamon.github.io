/**
 * Java实现的双向链表
 * 注： Java自带的集合包中有实现双向链表，路径是: java.util.LinkedList
 */
public class DoubleLink<T> {

    //双向链表"节点"对应的结构体
    private class DNode<T> {
        public DNode prev;
        public DNode next;
        public T value;

        public DNode(T value, DNode prev, DNode next) {
            this.value = value;
            this.prev = prev;
            this.next = next;
        }
    }

    //表头
    private DNode<T> mHead;
    //节点个数
    private int mCount;

    /**
     * 创建表头，表头没有存储数据
     * 初始化节点个数为0
     */
    public DoubleLink() {
        mHead = new DNode<T>(null, null, null);
        mHead.prev = mHead.next = mHead;
        mCount = 0;
    }

    /**
     * 
     * @return 节点数目
     */
    public int size() {
        return mCount;
    }

    /**
     * 
     * @return 链表是否为空
     */
    public boolean isEmpty() {
        return mCount == 0;
    }

    /**
     * 
     * @param index
     * @return index位置的节点
     */
    public DNode<T> getNode(int index) {
        if (index < 0 || index >= mCount) {
            throw new IndexOutOfBoundsException();
        }
        //正向查找
        if (index <= mCount / 2) {
            DNode<T> node = mHead.next;
            for (int i = 0; i < index; i++) {
                node = node.next;
            }
            return node;
        }

        //反向查找
        DNode<T> rnode = mHead.prev;
        int rindex = mCount - index - 1;
        for (int j = 0; j < rindex; j++) {
            rnode = rnode.prev;
        }
        return rnode;
    }

    /**
     * 
     * @param index
     * @return index位置的节点的值
     */
    public T get(int index) {
        return getNode(index).value;
    }

    /**
     * 
     * @return 获取第一个节点
     */
    public T getFirst() {
        return getNode(0);
    }

    /**
     * 
     * @return 获取最后一个节点的值
     */
    public T getLast() {
        return getNode(mCount - 1).value;
    }

    /**
     * 将节点插入到第index位置之前
     * @param index
     * @param t
     */
    public void insert(int index, T t) {
        //表头
        if (index == 0) {
            DNode<T> node = new DNode<T>(t, mHead, mHead.next);
            mHead.next.prev = node;
            mHead.next = node;
            mCount++;
            return;
        }

        DNode<T> inode = getNode(index);
        DNode<T> tnode = new DNode<T>(t, inode.prev, inode);
        inode.prev.next = tnode;
        inode.prev = tnode;
    }

    /**
     * 将节点插入第一个节点处
     * @param t
     */
    public void insertFirst(T t) {
        insert(0, t);
    }

    /**
     * 将节点追加到链表的末尾
     * @param t
     */
    public void appendLast(T t) {
        DNode<T> node = new DNode<T>(t, mHead.prev, mHead);
        mHead.prev.next = node;
        mHead.prev = node;
        mCount++;
    }

    /**
     * 删除index位置的节点
     * @param index
     */
    public void del(int index) {
        DNode<T> inode = getNode(index);
        inode.prev.next = inode.next;
        inode.next.prev = inode.prev;
        inode = null;
        mCount--;
    }

    /**
     * 删除第一个节点
     */
    public void delFirst() {
        del(0);
    }

    /**
     * 删除最后一个节点
     */
    public void delLast() {
        del(mCount - 1);
    }

    /**
     * 打印链表
     */
    @Override
    public String toString() {
        StringBuffer str = new StringBuffer("[");
        for (int i = 0; i < mCount; i++) {
            str.append(get(i)).append(",");
        }
        str.deleteCharAt(str.length()-1).append("]");
        return str.toString();
    }
}