package main

import (
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/handler"
	"wechat-dev/thinking/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// config init
	config.InitRedis()
	service.DoTemplate()

	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	var h handler.Handler
	// gin.Context，封装了request和response
	r.GET("/wx", h.Verify)
	r.POST("/wx", h.HandleMsgFromUser)
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8080")
}
