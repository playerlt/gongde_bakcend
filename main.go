package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"gongde/config"
	"gongde/controller"
	models "gongde/model"
	"gongde/public/common"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	common.InitRedis(config.RedisAddr)
	common.InitMysql(config.MysqlDataSource)
	config.Counter = getCounter()
	common.RDB.Set(context.Background(), config.RedisKey, "0", time.Duration(config.Counter)*time.Second)
	go doTask()

	router := gin.Default()
	router.Use(common.Cors()) //配置跨域
	router.POST("/gongde/update", controller.GongDeUpdate)
	router.GET("/getCount", controller.Ws)
	router.Run(":" + config.ServerPort)
	//执行定时任务
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server")
}

func doTask() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	s := gocron.NewScheduler(timezone)
	// 每小时执行一次
	_, err := s.Every(1).Hour().Do(func() {
		go func() { //redis备份到数据库
			ttl := int64(common.RDB.TTL(context.Background(), config.RedisKey).Val().Seconds())
			num, _ := common.RDB.Get(context.Background(), config.RedisKey).Int64()
			err := common.DB.Table("gongde_basic").
				Exec("UPDATE gongde_basic SET count = ? where id = 1", ttl+num).Error
			if err != nil {
				log.Fatal("备份失败: ", err)
			}
		}()
	})
	if err != nil {
		log.Fatal("备份失败: ", err)
		return
	}
	s.StartBlocking()
}

//从数据库中同步计数器到redis
func getCounter() int {
	var gongde models.GongdeBasic
	err := common.DB.Where("id = 1").First(&gongde).Error
	if err != nil {
		log.Fatal("从数据库导入到redis失败: ", err)
	}
	return gongde.Count
}
