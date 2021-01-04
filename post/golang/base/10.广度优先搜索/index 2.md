
# 广度优先算法
- 为爬虫实战做好准备
- 应用广泛，综合性强
- 面试常见

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/migong.png)

&emsp;&emsp;上图中是一个6*5的迷宫，0 代表可以走，1 代表墙，从左上角进右下角出，只能上左下右，不能对角线走。我们希望能直接给出最短路线，不会出现走死路重新走的情况，并且不会出现绕路的情况。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/zero.png)

&emsp;&emsp;起点用0表示，花了0步走到这个点。探索上左下右。那么我们说这个0是我们发现但没有探索过的点。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/one.png)

&emsp;&emsp;探索的结果就是0周围有四个1.这四个1又是我们已经发现但还没探索的点，这时候中间的0就是已经探索过了的点。

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/top.png)

&emsp;&emsp;最上面的1一开始是发现为探索，开始探索，发现了三个2.我们要等到1探索完了才能轮到2.

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/two.png)

&emsp;&emsp;1探索完了，这时才轮到2.

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/three.png)

&emsp;&emsp;我们发现所有的点，会经历三种状态: 以发现未探索，以发现以探索，未发现。

&emsp;&emsp;接下来来看一下迷宫

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/first.png)

&emsp;&emsp;左上角是迷宫对照图，管道是一个队列。从(0,0)开始，发现了(1,0),放入队列

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/second.png)

&emsp;&emsp;探索队列(1,0)，发现了(2,0),(1,1),开始探索(2,0)

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/second_1.png)

&emsp;&emsp;探索(2,0)发现都走不通，开始队列(1,1)

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/third.png)

&emsp;&emsp;探索(1,1)发现了(1,2),加入队列，再探索(1,2)

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/fourth.png)

&emsp;&emsp;探索(1,2),发现了(0,2),(2,2)，加入队列，开始探索(0,2)。。。走到最后的结果是什么呢？

![avatar](https://cdn.jsdelivr.net/gh/facedamon/MarkDownPhotos@master/golang/finally.png)

&emsp;&emsp;那么最终的路线就画出来了，从后往前倒过来认路，因为倒过来每个数字只有一个小于自己点。那这样的广度优先算法什么时候结束呢？
- 走到终点
- 走到队列为空(死路)

# 实现
&emsp;&emsp;迷宫的样式存在文件maze.in里。
		6 5
		0 1 0 0 0
		0 0 0 1 0
		0 1 0 1 0
		1 1 1 0 0
		0 1 0 0 1
		0 1 0 0 0

		package main
		
		func readMaze(filename string) [][]int {
		    file, err := os.Open(filename)
		    if err != nil {
		        panic(err)
		    }
		
		    var row, col int
		    fmt.Fscanf(file, "%d %d", &row, &col)
		    maze := make([][]int, row)
		    for i := range maze {
		        maze[i] = make([]int, col)
		        for j := range col {
		            fmt.Fscanf(file, "%d", &maze[i][j])
		        }
		    }
		
		    return maze
		}
		
		type point struct {
		    i, j int
		}
		
		// 上左下右
		var dirs = [4]point{
		    {-1, 0},
		    {0, -1},
		    {1, 0},
		    {0, 1},
		}
		
		func (p point) add(r point) point {
		    return point{p.i + r.i, p.j + r.j}
		}
		
		// 计算next在grid上的坐标是否越界
		func (p point) at(grid [][]int) (grid, bool) {
		    // 往上，下越界 =行越界
		    if p.i < 0 || p.i >= len(grid) {
		        return 0, false
		    }
		    // 列越界
		    if p.j < 0 || p.j >= len(grid[p.i]) {
		        return 0, false
		    }
		
		    return grid[p.i][p.j], true
		}
		
		func walk(maze [][]int, start, end point) [][]int {
		    // 从start开始走了多少步
		    // steps和maze一样大，它表示图片中的数字
		    steps := make([][]int, len(maze))
		    for i := range steps {
		        steps[i] = make([]int, len(maze[i]))
		    }
		    // 将开始点加入探索队列
		    Q := []point{start}
		
		    for len(Q) > 0 {
		        // 探索队列头部
		        cur := Q[0]
		        Q = Q[1:]
		
		        if cur == end {
		            break
		        }
		
		        for _, dir := range dirs {
		            // 新发现节点 = 当前节点+方向
		            next := cur.add(dir)
		            
		            // maze at next is 0        下一个点是空白
		            // and steps at next is 0   下一个点没走过
		            // and next != start
		            val, ok := if next.at(maze)
		            if !ok || val == 1 {
		                continue
		            }
		
		            val, ok = next.at(steps) 
		            if !ok || val != 0 {
		                continue
		            }
		            
		            if next == start {
		                continue
		            }
		
		            // 填入下一个探索steps
		            // 加入下一个探索队列
		            curSteps, _ := cur.at(steps)
		            steps[next.i][next.j] = curSteps + 1
		
		            Q = append(Q, next)
		
		        }
		
		    }
		    return steps
		}
		
		func main() {
		    maze := readMaze("maze/maze.in")
		
		    for _, row := range maze {
		        for _, val := range row {
		            fmt.Printf("%d ", val)
		        }
		        fmt.Println()
		    }
		
		    steps := walk(maze, point{0,0}, point{len(maze) - 1, len(maze[0]) - 1})
		    for -, row := range steps {
		        for _, val := range row {
		            fmt.Printd("%3d", val)
		        }
		        fmt.Println()
		    }
		}
- 用循环创建二维slice
- 使用slice实现队列
- 用Fscanf读取文件
- 对point抽象
