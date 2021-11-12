# dbxorm
基于 `go-xorm/xorm` 封装，将 `xorm` 实例工厂化，便于接入项目使用。  

## 安装
```shell
go get github.com/hy0kl/dbxorm
```

## 数据库配置

```ini
[DbConfig]
; 数据库类型
driver=mysql
; 是否打开执行sql记录
showSql=true
; 是否记录sql执行时间，前提需打开showSql
showExecTime=true
; 慢执行日志，前提打开showExecTime，单位500s，默认为0，打全部执行耗时日志
slowQueryDuration=500
; 最大连接数，默认100
maxConn=50
; 最大空闲连接数，默认30
maxIdle=30

[DbCluster]
; group = 主库 从库1 从库2 ... 从库n
blog = blog_rw:123456@tcp(localhost:3306)/blog_test blog_ro:123456@tcp(localhost:3306)/blog_test
shop = shop_rw:123456@tcp(localhost:3306)/shop_test shop_ro:123456@tcp(localhost:3306)/shop_test
```
> `blog` 与 `shop` 分别为两组不同的实例; 支持一主多从，配置至少 一主一从。

## 开始

```go
// 获取实例组
db := dbxorm.GetDbInstance("blog")

// 返回 Master 数据库引擎
masterEngine := db.Engine.Master()

// 依据给定的负载策略返回一个 Slave 数据库引擎
slaveEngine := db.Engine.Slave()

type Users struct {
    Id int64
    Name string
    Salt string
    Age int
    Passwd string `xorm:"varchar(200)"`
    Created time.Time `xorm:"created"`
    Updated time.Time `xorm:"updated"`
}

// 主库插入
affected, err := masterEngine.Insert(&user)
// INSERT INTO struct () values ()

affected, err := masterEngine.Insert(&user1, &user2)
// INSERT INTO struct1 () values ()
// INSERT INTO struct2 () values ()

affected, err := masterEngine.Insert(&users)
// INSERT INTO struct () values (),(),()

affected, err := masterEngine.Insert(&user1, &users)
// INSERT INTO struct1 () values ()
// INSERT INTO struct2 () values (),(),()

// 从库查询
has, err := slaveEngine.Get(&user)
// SELECT * FROM user LIMIT 1

has, err := slaveEngine.Where("name = ?", name).Desc("id").Get(&user)
// SELECT * FROM users WHERE name = ? ORDER BY id DESC LIMIT 1
```

## 详细使用介绍：
https://gobook.io/read/gitea.com/xorm/manual-zh-CN/

## 源码包
https://github.com/go-xorm/xorm