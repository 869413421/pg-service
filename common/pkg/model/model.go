package model

import (
	"fmt"
	"github.com/869413421/pg-service/common/pkg/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
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
	return model.CreatedAt.Format("2006-01-02")
}

func (model BaseModel) UpdatedAtDate() string {
	return model.UpdatedAt.Format("2006-01-02")
}

var db *gorm.DB

func connectDB() *gorm.DB {
	// 从系统环境变量获取数据库信息
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	DBName := os.Getenv("DB_DATABASE")
	password := os.Getenv("DB_PASSWORD")

	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			user, password, host, DBName,
		),
	)

	if err != nil {
		panic(fmt.Sprintf("connection to db error %v", err))
	}

	return db
}

func setupDB() {
	//1.连接数据库
	conn := connectDB()
	conn.Set("gorm:table_options", "ENGINE=InnoDB")
	conn.Set("gorm:table_options", "Charset=utf8")
	sqlDB := conn.DB()

	//2.设置最大连接数
	dbMaxConnections, _ := types.StringToInt(os.Getenv("DB_MAX_CONNECTIONS"))
	sqlDB.SetMaxOpenConns(dbMaxConnections)

	//3.设置最大空闲连接数
	dbMaxIdeConnections, _ := types.StringToInt(os.Getenv("DB_MAX_IDE_CONNECTIONS"))
	sqlDB.SetMaxIdleConns(dbMaxIdeConnections)

	//4. 设置每个链接的过期时间
	dbConnectionMaxLifeTime, _ := types.StringToInt(os.Getenv("DB_CONNECTIONS_MAX_LIFE_TIME"))
	sqlDB.SetConnMaxLifetime(time.Duration(dbConnectionMaxLifeTime) * time.Minute)
	db = conn
	fmt.Println("setting")
	fmt.Println(dbMaxConnections)
	fmt.Println(dbMaxIdeConnections)
	fmt.Println(dbConnectionMaxLifeTime)
}

// GetDB 开放给外部获得db连接
func GetDB() *gorm.DB {
	if db == nil {
		fmt.Println("connect setup")
		setupDB()
	} else {
		fmt.Println("not connect")
	}
	sqlDB := db.DB()
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		setupDB()
	}

	return db
}
