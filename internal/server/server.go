package server

import (
	"context"
	"fmt"
	"net/http"

	_ "job_solition/docs"
	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/middleware"
	"job_solition/internal/routes"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config     *config.Config
	router     *gin.Engine
	httpServer *http.Server
	postgres   *db.PostgreSQL
}

func NewServer(cfg *config.Config) *Server {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	postgres, err := db.NewPostgreSQL(cfg.PostgreSQL)
	if err != nil {
		panic(fmt.Errorf("ошибка при подключении к PostgreSQL: %w", err))
	}

	if err := postgres.InitDatabase(); err != nil {
		panic(fmt.Errorf("ошибка при инициализации базы данных: %w", err))
	}

	srv := &Server{
		config:   cfg,
		router:   router,
		postgres: postgres,
		httpServer: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: router,
		},
	}

	srv.setupServer()

	return srv
}

func (s *Server) setupServer() {
	utils.SetupValidators()

	s.router.Use(gin.Recovery())
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.CORS())

	s.router.Use(middleware.RateLimit(s.config.RateLimit.Requests, s.config.RateLimit.Duration))

	routes.SetupAllRoutes(s.router, s.postgres, s.config)

	swaggerConfig := ginSwagger.Config{
		URL:          "/swagger/doc.json",
		DocExpansion: "none",
		DeepLinking:  true,
	}

	s.router.GET("/swagger/*any", ginSwagger.CustomWrapHandler(&swaggerConfig, swaggerFiles.Handler))
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.postgres.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии соединения с PostgreSQL: %w", err)
	}

	return s.httpServer.Shutdown(ctx)
}
