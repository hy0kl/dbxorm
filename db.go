package dbdao

import (
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/go-xorm/xorm"
	"github.com/hy0kl/gconfig"
	"github.com/spf13/cast"
)

var initOnce sync.Once

type DBDao struct {
	Engine *xorm.EngineGroup
}

var (
	dbInstance   map[string]*DBDao
	curDbPoses   map[string]*uint64 // 当前选择的数据库
	showSql      bool
	showExecTime bool
	slowDuration time.Duration

	maxConn int = 100
	maxIdle int = 30
)

func newDBDaoWithParams(hosts []string, driver string) (Db *DBDao) {
	Db = new(DBDao)
	//engine, err := xorm.NewEngine(driver, host)
	engine, err := xorm.NewEngineGroup(driver, hosts)

	Db.Engine = engine

	// TODO: 增加存活检查
	if err != nil {
		log.Fatal(err)
	}

	Db.Engine.SetMaxOpenConns(maxConn)
	Db.Engine.SetMaxIdleConns(maxIdle)
	Db.Engine.SetConnMaxLifetime(time.Second * 3000)
	Db.Engine.ShowSQL(showSql)
	Db.Engine.ShowExecTime(showExecTime)
	Db.Engine.SetLogger(dbLogger)

	return
}

func Init() {
	initOnce.Do(func() {
		initDb()
	})
}

func initDb() {
	dbInstance = make(map[string]*DBDao, 0)
	curDbPoses = make(map[string]*uint64)

	var dbConfig = gconfig.GetConfStringMap("DbConfig")
	showSql = dbConfig["showSql"] == "true"
	showExecTime = dbConfig["showExecTime"] == "true"
	slowDuration = time.Duration(cast.ToInt(dbConfig["slowQueryDuration"])) * time.Millisecond

	maxConnConfig := cast.ToInt(dbConfig["maxConn"])
	if maxConnConfig > 0 {
		maxConn = maxConnConfig
	}

	maxIdleConfig := cast.ToInt(dbConfig["maxIdle"])
	if maxIdleConfig > 0 {
		maxIdle = maxIdleConfig
	}

	if maxIdle > maxConn {
		maxIdle = maxConn
	}

	for cluster, hosts := range gconfig.GetConfArrayMap("DbCluster") {
		instance := cluster
		dbInstance[instance] = newDBDaoWithParams(hosts, cast.ToString(dbConfig["driver"]))
		curDbPoses[instance] = new(uint64)
	}
}

func GetDbInstance(db string) *DBDao {
	Init()

	if instances, ok := dbInstance[db]; ok {
		return instances
	} else {
		panic(fmt.Sprintf(`can not use database, db: %s`, db))
	}
}

func (r *DBDao) Close() error {
	return r.Engine.Close()
}
