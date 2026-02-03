package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/wingnut128/outlier-go/docs" // swagger docs
	"github.com/wingnut128/outlier-go/internal/config"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *gin.Engine
}

// NewServer creates a new HTTP server with the given configuration
func NewServer(cfg *config.Config) *Server {
	// Set Gin mode based on log level
	if cfg.Logging.Level == "debug" || cfg.Logging.Level == "trace" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(requestLogger(cfg))

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Body size limit (100MB)
	router.MaxMultipartMemory = 100 << 20 // 100 MB

	s := &Server{
		config: cfg,
		router: router,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	s.router.GET("/health", handleHealth)
	s.router.POST("/calculate", handleCalculate)
	s.router.POST("/calculate/file", handleCalculateFile)

	// Swagger documentation
	s.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// requestLogger creates a custom logger middleware based on config
func requestLogger(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		if raw != "" {
			path = path + "?" + raw
		}

		// Format log based on config
		switch cfg.Logging.Format {
		case "json":
			log.Printf(`{"time":"%s","status":%d,"latency":"%s","ip":"%s","method":"%s","path":"%s"}`,
				start.Format(time.RFC3339),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		default:
			log.Printf("[%s] %d %s %s %s %s",
				start.Format("2006-01-02 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.BindIP, s.config.Server.Port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Channel to listen for errors from the server
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		log.Printf("Outlier API server listening on http://%s\n", addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or error
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
	case sig := <-shutdown:
		log.Printf("Received signal %v, shutting down gracefully...\n", sig)

		// Graceful shutdown with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}
		log.Println("Server stopped gracefully")
	}

	return nil
}
