package glowhttp

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luaxlou/glow/starter/glowapp"
)

var (
	engine      *gin.Engine
	once        sync.Once
	initialized bool
	port        string
)

func Init(p int) {
	port = ":" + strconv.Itoa(p)
	// Initialize if explicitly called, though Run() will check too
	// Maybe we can load port here if we want to fail fast
}

// Router returns the singleton Gin engine.
// It initializes the engine on the first call.
func Router() *gin.Engine {
	once.Do(func() {
		// Set Gin mode based on env or default to Release
		if os.Getenv("GIN_MODE") == "" {
			gin.SetMode(gin.ReleaseMode)
		}
		engine = gin.Default()

		// Add default middleware if needed
		// engine.Use(gin.Logger())
		// engine.Use(gin.Recovery())

		initialized = true
		log.Println("HTTP Starter (Gin) initialized.")
	})
	return engine
}

// Run starts the HTTP server on the configured port.
// It returns immediately, running the server in a goroutine.
// Lifecycle is managed by glowapp.
func Run() {
	if !initialized {
		// Initialize if not already done (e.g. user didn't add any routes but just wants to start)
		Router()
	}

	if port == "" {
		// Priority: OP_APP_PORT > PORT > 8080
		p := os.Getenv("OP_APP_PORT")
		if p == "" {
			p = os.Getenv("PORT")
		}
		if p == "" {
			p = "8080"
		}
		// If p is just a number, prefix with :
		if _, err := strconv.Atoi(p); err == nil {
			port = ":" + p
		} else {
			port = p
		}
	}

	// Report to glowapp
	// Extract port number
	portNumStr := strings.TrimPrefix(port, ":")
	if portNum, err := strconv.Atoi(portNumStr); err == nil {
		domain := os.Getenv("OP_APP_DOMAIN")
		if domain == "" {
			domain = os.Getenv("DOMAIN")
		}
		glowapp.SetInfo(portNum, domain)
	}

	srv := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	go func() {
		log.Printf("Server starting on %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	glowapp.RegisterCleanup("HTTP Starter", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		}
		log.Println("HTTP Server exiting")
	})
}
