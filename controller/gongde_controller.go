package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gongde/config"
	"gongde/public/common"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func GongDeUpdate(c *gin.Context) {
	ttl := int64(common.RDB.TTL(c, config.RedisKey).Val().Seconds())
	if ttl == 0 {
		c.JSON(200, gin.H{"count": ttl})
		return
	}
	num := common.RDB.Incr(c, config.RedisKey).Val()
	go func() {
		if ttl < 3600 { //功德小于1小时更新ttl
			common.RDB.Expire(c, config.RedisKey, time.Duration(ttl+num)*time.Second)
			common.RDB.DecrBy(c, config.RedisKey, num)
		}
	}()
	c.JSON(200, gin.H{"count": ttl + num, "msg": "success"})
}

func Ws(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) // error ignored for sake of simplicity
	if err != nil {
		c.JSON(400, gin.H{"msg": "ws连接失败"})
	}
	defer conn.Close()
	for {
		// Read message from browser
		//msgType, msg, err := conn.ReadMessage()
		time.Sleep(time.Second)
		ttl := int64(common.RDB.TTL(c, config.RedisKey).Val().Seconds())
		num, _ := common.RDB.Get(c, config.RedisKey).Int64()

		// Print the message to the console
		//fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = conn.WriteMessage(websocket.TextMessage, []byte(strconv.FormatInt(ttl+num, 10))); err != nil {
			return
		}
	}
}
