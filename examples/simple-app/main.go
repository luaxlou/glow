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
	// 1. Initialize & Register App Identity
	glowapp.Init("simple-app")

	glowhttp.Init(33203)
	glowmysql.Init("simple_app_db")

	fmt.Printf("App %s starting\n", glowapp.Name())

	// 2. Implicit wiring: usage triggers configuration & connection
	// The MySQL component calls config, which calls app (identity) to resolve configuration.
	db, err := glowmysql.DB()
	if err != nil {
		log.Printf("Warning: MySQL not available: %v. Running without DB.", err)
	} else {
		fmt.Println("MySQL is ready and connected!")
		// Verify connection
		if err := db.Ping(); err != nil {
			log.Printf("Failed to ping DB: %v", err)
		}
	}

	r := glowhttp.Router()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello from Implicit-wiring app with Gin!")
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

	glowhttp.Run()
	glowapp.WaitForShutdown()

}
