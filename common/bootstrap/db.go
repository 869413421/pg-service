package bootstrap

import (
	"github.com/869413421/pg-service/common/pkg/model"
	"github.com/869413421/pg-service/common/pkg/types"
	"os"
	"time"
)

func SetupDB() {
	//1.连接数据库
	db := model.ConnectDB()
	sqlDB := db.DB()

	//2.设置最大连接数
	dbMaxConnections, _ := types.StringToInt(os.Getenv("DB_MAX_CONNECTIONS"))
	sqlDB.SetMaxOpenConns(dbMaxConnections)

	//3.设置最大空闲连接数
	dbMaxIdeConnections, _ := types.StringToInt(os.Getenv("DB_MAX_IDE_CONNECTIONS"))
	sqlDB.SetMaxIdleConns(dbMaxIdeConnections)

	//4. 设置每个链接的过期时间
	dbConnectionMaxLifeTime, _ := types.StringToInt(os.Getenv("DB_CONNECTIONS_MAX_LIFE_TIME"))
	sqlDB.SetConnMaxLifetime(time.Duration(dbConnectionMaxLifeTime) * time.Minute)
}
