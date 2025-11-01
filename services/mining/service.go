package mining

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// New creates a new mining service
func New() *Service {
	return &Service{
		Miners: make(map[string]Miner),
	}
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Initialize miners
	s.Miners["xmrig"] = NewXMRigMiner()
	s.Miners["ttminer"] = NewTTMiner()

	// Initialize Gin router
	s.Router = gin.Default()
	s.setupRoutes()

	// Create and start the HTTP server
	s.Server = &http.Server{
		Addr:    ":8080",
		Handler: s.Router,
	}

	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", s.Server.Addr, err)
		}
	}()

	// Listen for context cancellation to gracefully shut down the server
	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Server.Shutdown(ctxShutdown); err != nil {
			log.Fatalf("server shutdown failed: %+v", err)
		}
	}()

	return nil
}

func (s *Service) minerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		minerName := c.Param("miner_name")
		miner, ok := s.Miners[minerName]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "miner not found"})
			c.Abort()
			return
		}
		c.Set("miner", miner)
		c.Next()
	}
}

func (s *Service) setupRoutes() {
	minerGroup := s.Router.Group("/miners/:miner_name")
	minerGroup.Use(s.minerMiddleware())
	{
		minerGroup.POST("/install", s.handleInstallMiner)
		minerGroup.POST("/start", s.handleStartMiner)
		minerGroup.POST("/stop", s.handleStopMiner)
		minerGroup.GET("/stats", s.handleGetMinerStats)
	}
}

func (s *Service) handleInstallMiner(c *gin.Context) {
	miner := c.MustGet("miner").(Miner)
	if err := miner.Install(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "installed"})
}

func (s *Service) handleStartMiner(c *gin.Context) {
	miner := c.MustGet("miner").(Miner)
	var config Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := miner.Start(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "started"})
}

func (s *Service) handleStopMiner(c *gin.Context) {
	miner := c.MustGet("miner").(Miner)
	if err := miner.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "stopped"})
}

func (s *Service) handleGetMinerStats(c *gin.Context) {
	miner := c.MustGet("miner").(Miner)
	stats, err := miner.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetMiner returns a miner by name
func (s *Service) GetMiner(name string) (Miner, error) {
	miner, ok := s.Miners[name]
	if !ok {
		return nil, fmt.Errorf("miner not found: %s", name)
	}
	return miner, nil
}
