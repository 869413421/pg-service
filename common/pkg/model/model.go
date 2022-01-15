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

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	// 从系统环境变量获取数据库信息
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	DBName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	DB, err := gorm.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			user, password, host, DBName,
		),
	)

	if err != nil {
		panic(fmt.Sprintf("connection to db error %v", err))
	}

	return DB
}
