
> 转载自 https://segmentfault.com/a/1190000003036452

# 概述

&emsp;&emsp;`sql.DB`不是一个连接，它是数据库的抽象接口。它可以根据driver打开关闭数据库连接，管理连接池。正在使用的连接被标记为繁忙，用完后回到连接吃等待下次使用。所以，如果你没有把连接释放回连接池，会导致过多连接使系统资源耗尽。

## 使用DB
### 导入driver

		import (
		    "database/sql"
		    "log"
		    _ "github.com/go-sql-driver/mysql"
		)

### 连接DB

		var db *sql.DB
		
		func main(
		    db, err := sql.Open("mysql", "user:passwd@tcp(127.0.0.1:3306)/hello")
		    if err != nil {
		        log.Fatal(err)
		    }
		    defer db.Close()
		    err = db.Ping()
		    if err != nil {
		        // do something here
		    }
		)

&emsp;&emsp;`sql.Open`的第一个参数树driver名称，第二个参数是driver连接数据库的信息，各个driver可能不同。DB不是连接，并且只有当需要使用时才会创建连接，如果像立即验证连接，需要使用`Ping`.

&emsp;&emsp;sql.DB的设计就是用来作为长连接使用的。不要频繁的Open,Close.比较好的做法是，为每个不同的**datastore**创建对应的DB对象，保持这些对象Open。如果你需要短连接，可以将DB作为func的参数传入，而不要在func中Open，Close

### 读取DB

&emsp;&emsp;如果方法包含`Query`，那么这个方法是用于查询并返回rows的。其他情况应该用`Exec`。

		var (
		    id      int
		    name    string
		)
		
		rows, err := db.Query("select id, name from users where id = ?", 1)
		if err != nil {
		    log.Fatal(err)
		}
		defer rows.Close()
		
		for rows.Next() {
		    err := rows.Scan(&id, &name)
		    if err != nil {
		        log.Fatal(err)
		    }
		    log.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
		    log.Fatal(err)
		}

&emsp;&emsp;上面代码的过程为:`db.Query`表示向数据库发送一个query，`defer rows.Close`非常重要，遍历rows使用`rows.Next`，把遍历到的数据存入变量使用`rows.Scan`，遍历完后检查err，有几点需要注意：
1. 检查遍历是否有error
2. 结果集rows未关闭前，底层的连接处于繁忙状态。当遍历读到最后一条记录时，会发生一个内部EOF错误，自动调用rows.Close,但是如果提前退出循环，rows不会关闭，连接不会回到连接池中，连接也不会关闭。所以手动关闭非常重要。rows.Close可以多次调用，是无害操作。

### 单行Query

&emsp;&emsp;如果err在`Scan`后才产生，所以可以：

		var name string
		err = db.QueryRow("select name from users where id = ?", 1).Scan(&name)
		if err != nil {
		    log.Fatal(err)
		}
		fmt.Println(name)

## 处理Error
### 循环Rows的Error

&emsp;&emsp;如果循环中发生错误会自动调用`rows.Close`，用`rows.Err`接收这个错误，循环之后判断error是非常有必要的

		for rows.Nexr() {
		    // ...
		}
		if err = rows.Err(); err != nil {
		    // handle the err here
		}

### 关闭Resultsets时的error

&emsp;&emsp;如果你在rows遍历结束之前退出循环，必须手动关闭Resultset,并接收error

		for rows.Next() {
		    // ...
		    break;
		}
		
		// do the usual 'if err = rows.Err'
		// it`s always safe to close here
		if err = rows.Close(); err != nil {
		    // but what should we do if there`s an error ?
		    log.Println(err)
		}

### QueryRow的error

		var name string
		err = db.QueryRow("select name from users where id = ?", 1).Scan(&name)
		if err != nil {
		    log.Fatal(err)
		}
		fmt.Println(name)

&emsp;&emsp;如果id为1的记录不存在，err为`sql.ErrNoRows`,一般应用中不存在的情况都需要单独处理。此外，Query返回的错误都会延迟到Scan被调用，所以应该写成这样：

		var name string
		err = db.QueryRow("select name from users where id = ?", 1).Scan(&name)
		if err != nil {
		    if err == sql.ErrNoRows {
		        // there were no rows, but otherwise no error occurred
		    } else {
		        log.Fatal(err)
		    }
		}
		fmt.Println(name)

### 分析数据库Error

&emsp;&emsp;各个数据库处理方式不太一样，mysql为例：

		if driverErr, ok := err.(*mysql.MySQLError); ok {
		    if driverErr.Number == 1045 {
		        // handle the permission-denied error
		    }
		}

&emsp;&emsp;`MySQLError, Number`都是DB特异的，别的数据库可能是别的类型或字段。这里的数字可以替换成常量。[MySQL error numbers maintained by VividCortex](https://github.com/VividCortex/mysqlerr)

## 连接错误
### NULL值处理

&emsp;&emsp;简单说就是设计数据库时不要出现null，处理起来非常费力。Null的type很有限，例如没有`sql.NullUnit64`;null值没有默认零值。

		for rows.Next() {
		    var s sql.NullString
		    err := rows.Scan(&s)
		    // check err
		    if s.Valid {
		        // use s.String
		    } else {
		        // NULL value
		    }
		}

### 未知Column

&emsp;&emsp;`rows.Columns`的使用，用于处理不能得知结果字段个数或类型的情况。

		cols, err := rows.Columns()
		if err != nil {
		    // handle the error
		} else {
		    dest := []interface{} {
		        // Standard MySQL columns
		        new(unit64),    // id
		        new(string),    // host
		        new(string),    // user
		        new(string),    // db
		        new(string),    // command
		        new(uint32),    // time
		        new(string),    // state
		        new(string),    // info
		    }
		    if len(cols) == 11 {
		        // 启动 Server
		    } else if len(cols) > 8 {
		        // handle this case
		    }
		    err = rows.Scan(dest...)
		    // work with the values in dest
		}
		cols, err := rows.Columns()
		vals := make([]interface{}, len(cols))
		for i, _ := range cols {
		    vals[i] = new(sql.RawBytes)
		}
		for rows.Next() {
		    err = rows.Scan(vals...)
		    // 现在，你可以检查每一个vals元素为nil的情况
		    // 并且，你可以对每个column使用断言来映射到指定的字段类型上
		}

## 修改数据

&emsp;&emsp;一般用Prepared Statements和`Exec`完成**INSERT, UPDATE, DELETE** operation.

		stmt, err := db.Prepare("insert into users(name) values(?)")
		if err != nil {
		    log.Fatal(err)
		}
		res, err := stmt.Exec("Dolly)
		if err != nil {
		    log.Fatal(err)
		}
		lastId, err := res.LastInsertId()
		if err != nil {
		    log.Fatal(err)
		}
		rowCnt, err := res.RowsAffected()
		if err != nil {
		    log.Fatal(err)
		}
		log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

## 事物

&emsp;&emsp;`db.Begin`开始事物，`commit, Rollback`关闭事物。`Tx`从连接池中取出一个连接，在关闭之前都是使用这个连接。**Tx不能和DB层的Begin，Commit混合使用。**

&emsp;&emsp;如果你需要通过多条语句修改连接状态，你必须使用Tx，例如：
- 创建仅对单个连接可见的临时表
- 设置变量，例如`SET @var := somevalue`
- 改变连接选项，例如字符集，超时

### Prepared Statements

&emsp;&emsp;在数据库层面，PS是和单个数据库连接绑定的。客户端发送一个有占位符的statement到服务器，服务器返回一个statement Id，然后客户端发送ID和参数来执行statement。

&emsp;&emsp;在Go中，连接不能直接暴露，你不能为连接绑定statement，只能为DB或TX绑定。database/sql包有自动重试功能。当你生产一个PS时：
1. 自动在连接池中绑定到一个空闲连接
2. stmt对象记住了绑定了哪个连接
3. 执行stmt时，尝试使用该连接。如果不可用，例如连接被关闭或繁忙，会自动re-prepare,绑定到另一个连接

&emsp;&emsp;这就导致在高并发的场景，过度使用statement可能导致statement泄漏，statement持续重复prepare和re-prepare的过程，甚至会达到服务器statement数量上限。

### 在TX中使用PS

&emsp;&emsp;PS在TX中唯一绑定一个连接，**不会re-prepare**

&emsp;&emsp;TX和statement不能分离，在DB中创建的statement也不能在TX中使用，因为他们必定不是使用的同一个连接，必须十分小心：

		tx, err := db.Begin()
		if err != nil {
		    log.Fatal(err)
		}
		defer tx.Rollback()
		stmt, err := tx.Prepare("insert into foo values (?)")
		if err != nil {
		    log.Fatal(err)
		}
		defer stmt.Close()      // danger
		for i := 0; i < 10; i++ {
		    _, err = stmt.Exec(i)
		    if err != nil {
		        log.Fatal(err)
		    }
		}
		err = tx.Commit()
		if err != nil {
		    log.Fatal(err)
		}
		// stmt.Close() runs here

&emsp;&emsp;`*sql.Tx`一旦释放，连接就回到连接池中，这里stmt在关闭时就无法找到连接。所以必须在tx commit或rollback之前就关闭stmt。

## GoStyle TX

&emsp;&emsp;确保Begin, Commit, Rollback出现在同一个函数中，它使事物更容易跟踪，并允许你通过defer来确保事物正确关闭。

		func (s service) dosomething() (err error) {
		    tx, err := s.db.Begin()
		    if err != nil {return}
		    defer func() {
		        if err != nil {
		            tx.Rollback()
		            return
		        }
		        err = tx.Commit()
		    }()
		    if _, err = tx.Exec(...); err != nil {
		        return
		    }
		    if _, err = tx.Exec(...); err != nil {
		        return
		    }
		}

&emsp;&emsp;**另一种方法是使用事物处理程序包装事物：**

		func transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
		    tx, err := db.Begin()
		    if err != nil {return}
		    defer func() {
		        if p := recover(); p != nil {
		            tx.Rollback()
		            panic(p)
		        } else if err != nil {
		            tx.Rollback()
		        } else {
		            err = tx.Commit()
		        }
		    }()
		    err = txFunc(tx)
		    return err
		}

&emsp;&emsp;使用上面的事物处理程序，你可以做到一下几点：

		func (s service) dosomething() error {
		    return transact(s.db, func (tx *sql.Tx) error {
		        if _, err := tx.Exec(...); err != nil {
		            return err
		        }
		        if _, err := tx.Exec(...); err != nil {
		            return err
		        }
		        return nil
		    })
		}

&emsp;&emsp;这样可以使你的事物更加简洁，并确保事物得到正确的处理。
