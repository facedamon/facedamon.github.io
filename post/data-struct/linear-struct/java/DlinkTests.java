/**
 * 
 */
public class DlinkTests {

    public static void test_int() {
        System.out.println("\n----test_int----");
        DoubleLink<Integer> dlnk = new DoubleLink<>();
        dlnk.insert(0, 20);
        dlnk.appendLast(10);
        dlnk.insertFirst(30);

        System.out.printf("isEmpty()=%b\n", dlnk.isEmpty());
        System.out.printf("size()=%d\n", dlnk.size());
        System.out.println(dlnk.toString());
    }

    public static void test_string() {
        System.out.println("\n----test_string----");
        DoubleLink<String> dlnk = new DoubleLink<>();
        dlnk.insert(0, "ten");
        dlnk.appendLast("twenty");
        dlnk.insertFirst("thirty");

        System.out.printf("isEmpty()=%b\n", dlnk.isEmpty());
        System.out.printf("size()=%d\n", dlnk.size());
        System.out.println(dlnk.toString());
    }

    private static class Student {
        private int id;
        private String name;

        public Student(int id, String name) {
            this.id = id;
            this.name = name;
        }

        @Override
        public String toString() {
            return "["+id+", "+name+"]";
        }
    }

    private static Student[] students = new Student[]{
        new Student(10, "sky"),
        new Student(20, "joy"),
        new Student(30, "vic"),
        new Student(40, "dan"),
    };

    private static void test_object() {
        System.out.println("\n----test_object----");
        DoubleLink<Student> dlnk = new DoubleLink<>();
        dlnk.insert(0, students[1]);
        dlnk.appendLast(students[0]);
        dlnk.insertFirst(students[2]);

        System.out.printf("isEmpty()=%b\n", dlnk.isEmpty());
        System.out.printf("size()=%d\n", dlnk.size());
        System.out.println(dlnk.toString());
    }

    public static void main(String[] args) {
        //test_int();// isEmpty = false; size() = 3; [30, 20, 10]
        //test_string(); // isEmpty = false; size() = 3; [thirty, ten, twenty]
        test_object(); // isEmpty = false; size() = 3; [[30, vic], [20, joy], [10, sky]]
    }
    
}