package orm

import (
	"context"
	"sync"
	"time"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// LadeGorm 数据库
type LadeGorm struct {
	container framework.Container // 服务容器

	configPath string
	dbs        map[string]*gorm.DB //多个DB连接，key为dsn
	config     *gorm.Config        // gorm配置文件，可以修改

	lock *sync.RWMutex // 读写dbs锁
}

// NewLadeGorm ，启动Gorm服务，参数只有container
func NewLadeGorm(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	dbs := make(map[string]*gorm.DB)
	lock := &sync.RWMutex{}
	return &LadeGorm{
		container: container,
		dbs:       dbs,
		lock:      lock,
	}, nil

}

// GetDB 根据选项配置，获取对应的DB对象
func (app *LadeGorm) GetDB(option ...contract.DBOption) (*gorm.DB, error) {
	// 用于再次函数内写日志
	logger := app.container.MustMake(contract.LogKey).(contract.Log)
	// 读取默认配置
	config := GetBaseConfig(app.container)
	// 用于存入orm的config中
	logService := app.container.MustMake(contract.LogKey).(contract.Log)

	// 设置Logger
	ormLogger := NewOrmLogger(logService)
	config.Config = &gorm.Config{
		Logger: ormLogger,
	}

	// 遍历option对config进行修改
	for _, opt := range option {
		if err := opt(app.container, config); err != nil {
			return nil, err
		}
	}

	// 最终的config没有dsn，就生成dsn
	if config.Dsn == "" {
		dsn, err := config.FormatDsn()
		if err != nil {
			return nil, err
		}
		config.Dsn = dsn
	}

	// 已有dsn对应的数据库连接实例，则返回该实例
	app.lock.RLock()
	if db, ok := app.dbs[config.Dsn]; ok {
		app.lock.RUnlock()
		return db, nil
	}
	app.lock.RUnlock()

	// 创建新DB实例
	app.lock.Lock()
	defer app.lock.Unlock()

	var db *gorm.DB
	var err error
	switch config.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Dsn), config)
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Dsn), config)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Dsn), config)
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(config.Dsn), config)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(config.Dsn), config)
	}

	// 获取标准库的sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	// 设置连接池参数
	if config.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
	}
	if config.ConnMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
	}
	if config.ConnMaxLifetime != "" {
		lifeTime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			logger.Error(context.Background(), "conn max life time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxLifetime(lifeTime)
		}

	}
	// TODO: 是否应该==nil才存入dbs
	if err != nil {
		app.dbs[config.Dsn] = db
	}

	return db, err
}
