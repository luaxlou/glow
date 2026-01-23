package glowapp

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/luaxlou/glow/starter/glowapp/config"
)

var (
	name   string
	port   int
	domain string

	// Lifecycle management
	shutdownCh   = make(chan struct{})
	shutdownOnce sync.Once
	wg           sync.WaitGroup

	// Deprecated: kept for backward compatibility.
	// Apps no longer register/heartbeat to glow-server at runtime.
	skipRegistration bool
)

type Option func()

func WithNoRegistration() Option {
	return func() {
		skipRegistration = true
	}
}

func Init(appName string, opts ...Option) {
	config.AppIdentity = appName
	name = appName
	for _, opt := range opts {
		opt()
	}
}

// SetInfo updates the app info and resends registration
func SetInfo(p int, d string) {
	port = p
	domain = d
}

func Name() string {
	return name
}

// RegisterCleanup registers a function to be called when the app shuts down.
// It uses a goroutine to wait for the shutdown signal.
func RegisterCleanup(componentName string, cleanup func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-shutdownCh
		log.Printf("[%s] Starting cleanup...", componentName)
		cleanup()
		log.Printf("[%s] Cleanup finished.", componentName)
	}()
}

// WaitForShutdown blocks until an OS signal is received.
// It then broadcasts the shutdown signal to all registered components
// and waits for them to complete.
func WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received signal: %v. Initiating shutdown...", sig)

	shutdownOnce.Do(func() {
		close(shutdownCh)
	})

	// Wait for all cleanups with a timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All components shut down gracefully.")
	case <-time.After(10 * time.Second):
		log.Println("Shutdown timed out. Forcing exit.")
	}
}
