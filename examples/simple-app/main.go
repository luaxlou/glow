package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	gw "github.com/gorilla/websocket"
	"github.com/luaxlou/glow/starter/glowapp"
	"github.com/luaxlou/glow/starter/glowhttp"
	"github.com/luaxlou/glow/starter/glowmysql"
	"github.com/luaxlou/glow/starter/glowwebsocket"
)

func main() {
	// 1. 初始化应用身份
	glowapp.Init("simple-app")

	// 2. 初始化 HTTP 服务
	glowhttp.Init(33203)

	// 3. 声明需要使用的 MySQL 数据库
	// 注意：不再调用 ProvisionResource()，资源由 glow apply 配置
	glowmysql.Init("simple_app_db")

	fmt.Printf("App %s starting\n", glowapp.Name())

	// 4. 使用 MySQL（从本地配置文件读取 DSN）
	// 配置文件由 'glow apply -f app.yaml' 生成
	// 位置: /var/lib/glow-server/apps/simple-app/simple-app_local_config.json
	db, err := glowmysql.Gorm()
	if err != nil {
		log.Printf("Warning: MySQL not available: %v. Running without DB.", err)
	} else {
		fmt.Println("MySQL is ready and connected!")
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Failed to get sql.DB: %v", err)
		} else if err := sqlDB.Ping(); err != nil {
			log.Printf("Failed to ping DB: %v", err)
		}
	}

	// 5. 设置 HTTP 路由
	r := glowhttp.Router()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello from Glow app with Gin!")
	})

	r.GET("/ws", func(c *gin.Context) {
		glowwebsocket.Handle(c, func(conn *gw.Conn) {
			for {
				messageType, p, err := conn.ReadMessage()
				if err != nil {
					return
				}
				if err := conn.WriteMessage(messageType, p); err != nil {
					return
				}
			}
		})
	})

	// 6. 启动服务
	glowhttp.Run()
	glowapp.WaitForShutdown()
}
