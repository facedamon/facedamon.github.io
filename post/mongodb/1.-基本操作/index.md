
## 基本操作

- 创建document

        #use dbname
        use chapter;

- 删除document

        #db.dropDatabase()
        use chapter;
        db.dropDatabase()

- 创建collection

        #db.createCollection(name, options)
        # example 集合空间大小614299KB,文档最大个数为1000个
        db.createCollection('example', {capped:true, autoIndexId: true, size: 6142800, max: 10000})

- 删除collection

        #db.collection.drop()
        use chapter;
        db.example.drop()

- 查询全部数据

        db.getCollection('example').find({})

- 插入单条数据

        db.getCollection('example').insertOne({'name':'张小二', 'age':22, 'address':'浙江'})

- 插入单条数据，数据格式不同

        db.getCollection('example').insertOne({'name':18030, 'age':'十岁', 'address': 3.5})

- 插入多条数据

        data_list = [
            {'name':'主小三','age':20,'address':'北京'},
            {'name':'刘小厮','age':21,'address':'上海'},
            {'name':'冯小五','age':22,'address':'山东'}
        ]
        db.getCollection('example').insertMany(data_list)

- 查询特定数据

        db.getCollection('example').find({'age':'22'})
        # 多个参数是and关系
        db.getCollection('example').find({'age':'22','name':'冯小五'})

- 查询范围数据

|操作符|意义|
|:-:|:-:|
|$gt|great than|
|$gte|great than and equal|
|$lt|less than|
|$lte|less than and equal|
|$ne|not equal|

        db.getCollection('example').find({'age':{'$gte':20}})
        db.getCollection('example').find({'age':{'$lt':25, '$gt':21}})
        db.getCollection('example').find({'age':{'$lt':25, '$gt':21}, 'name':{'$ne':'刘小三'}})

- 限定返回字段

        # 0表示不显示， 1表示显示
        # 0和1只能出现一个
        db.getCollection('example').find({}, {'address':0, 'age':0})
        db.getCollection('example').find({}, {'age':1, 'name':1})

- count

        db.getCollection('example').find({'age': {'$gt':21}}).count()

- limit

        # 如果限制返回的条数为一个数字，则表示最多返回这么多条记录。如果超过限定条数，则只返回限定的条数
        db.getCollection('example').find({'age':{'$gt':21}}).limit(4)

- sort

        # 1: 正序；-1:倒序
        db.getCollection('example').find({'age':{'$gt':21}}).sort({'name': -1})

- 修改数据

        #updateOne 只更新第一条满足要求的数据
        #updateMany 更新所有满足要求的数据
        db.getCollection('example').updateMany({'name':'王小六'}, {'$set': {'address': '苏州', 'work': '工程师'}})

- 删除数据

        #deleteOne 只删除第一条满足要求的数据
        #deleteMany 更新所有满足要求的数据
        db.getCollection('example').deleteMany({'age':22})
        #慎用删除功能，一般工程上使用deleted字段来表示是否删除

- 数据去重

        #第一个参数表示对哪一个字段去重
        #第二个参数表示筛选条件，相当于find的第一个参数
        db.getCollection('example').distinct('age', {'age':{'$gte':21}})
