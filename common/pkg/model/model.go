package model

import (
	"fmt"
	"github.com/869413421/pg-service/common/pkg/config"
	"github.com/869413421/pg-service/common/pkg/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uint64    "gorm:column:id;primaryKey;autoIncrement;not null"
	CreatedAt time.Time `gorm:"column:created_at;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;index"`
}

func (model BaseModel) GetStringID() string {
	return types.UInt64ToString(model.ID)
}

func (model BaseModel) CreatedAtDate() string {
	return model.CreatedAt.Format("2006-01-02 15:04:05")
}

func (model BaseModel) UpdatedAtDate() string {
	return model.UpdatedAt.Format("2006-01-02 15:04:05")
}

var db *gorm.DB
var dbConfig *config.Db

func connectDB() *gorm.DB {
	// 从系统环境变量获取数据库信息
	serviceConfig := config.LoadConfig()
	dbConfig = &serviceConfig.Db

	db, err := gorm.Open(mysql.Open(fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Address, dbConfig.Database, dbConfig.Charset,
	)), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("connection to db error %v", err))
	}

	//db.Use(dbresolver.Register(dbresolver.Config{
	//	// `db2` 作为 sources，`db3`、`db4` 作为 replicas
	//	Sources:  []gorm.Dialector{mysql.Open("db2_dsn")},
	//	Replicas: []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")},
	//	// sources/replicas 负载均衡策略
	//	Policy: dbresolver.RandomPolicy{},
	//}).Register(dbresolver.Config{
	//	// `db1` 作为 sources（DB 的默认连接），对于 `User`、`Address` 使用 `db5` 作为 replicas
	//	Replicas: []gorm.Dialector{mysql.Open("db5_dsn")},
	//}, &User{}, &Address{}).Register(dbresolver.Config{
	//	// `db6`、`db7` 作为 sources，对于 `orders`、`Product` 使用 `db8` 作为 replicas
	//	Sources:  []gorm.Dialector{mysql.Open("db6_dsn"), mysql.Open("db7_dsn")},
	//	Replicas: []gorm.Dialector{mysql.Open("db8_dsn")},
	//}, "orders", &Product{}, "secondary"))

	return db
}

func setupDB() {
	//1.连接数据库
	conn := connectDB()
	conn.Set("gorm:table_options", "ENGINE=InnoDB")
	conn.Set("gorm:table_options", "Charset=utf8")
	sqlDB, err := conn.DB()
	if err != nil {
		panic(err)
	}

	//2.设置最大连接数
	sqlDB.SetMaxOpenConns(dbConfig.MaxConnections)

	//3.设置最大空闲连接数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdeConnections)

	//4. 设置每个链接的过期时间
	sqlDB.SetConnMaxLifetime(dbConfig.ConnectionMaxLifeTime * time.Minute)
	db = conn
}

// GetDB 开放给外部获得db连接
func GetDB() *gorm.DB {
	if db == nil {
		setupDB()
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		setupDB()
	}

	return db
}
