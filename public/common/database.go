package common

import (
	models "gongde/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

// 全局mysql数据库变量
var DB *gorm.DB

func InitMysql(datasource string) {
	engine, err := gorm.Open(mysql.Open(datasource), &gorm.Config{})
	// engine, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		log.Printf("gorm New Engine Error:%v", err)
		os.Exit(1)
	}
	DB = engine
	dbAutoMigrate(DB)
}

// 自动迁移表结构
func dbAutoMigrate(DB *gorm.DB) {
	_ = DB.AutoMigrate(
		&models.GongdeBasic{},
	)
}
