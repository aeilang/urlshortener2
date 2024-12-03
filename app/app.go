package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aeilang/urlshortener/config"
	"github.com/aeilang/urlshortener/database"
	"github.com/aeilang/urlshortener/internal/api"
	"github.com/aeilang/urlshortener/internal/cache"
	"github.com/aeilang/urlshortener/internal/service"
	"github.com/aeilang/urlshortener/pkg/shortcode"
	"github.com/aeilang/urlshortener/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Applcation struct {
	e                  *echo.Echo
	db                 *sql.DB
	redis              *cache.RedisCache
	shortCodeGenerator *shortcode.ShortCodeGenerator
	cfg                *config.Config
	urlHandler         *api.URLHandler
	urlService         *service.URLService
}

func NewApplication(filePath string) (*Applcation, error) {
	app := new(Applcation)
	err := app.init(filePath)
	return app, err
}

func (a *Applcation) Run() {
	go a.startServer()
	go a.clearnup()
	a.shutdown()
}

func (a *Applcation) init(filePath string) error {
	// 初始化配置
	cfg, err := config.LoadConfig(filePath)
	if err != nil {
		return err
	}
	a.cfg = cfg

	// 初始化数据库
	db, err := database.NewDB(cfg.Database)
	if err != nil {
		return err
	}
	a.db = db

	// 初始化redis
	redis, err := cache.NewReisClient(cfg.Redis)
	if err != nil {
		return err
	}
	a.redis = redis

	// 初始化shortCodeGenerator
	a.shortCodeGenerator = shortcode.NewShortCodeGenerator(cfg.ShortCode)

	// 初始化service
	a.urlService = service.NewURLService(db, a.shortCodeGenerator, cfg, redis)

	// 初始化handler
	a.urlHandler = api.NewURLHandler(a.urlService)

	// 初始化echo
	a.initEcho()

	// 初始化路由
	a.e.POST("/api/url", a.urlHandler.CreateURL)
	a.e.GET("/:code", a.urlHandler.RedirectURL)

	return nil
}

func (a *Applcation) clearnup() {
	ticker := time.NewTicker(a.cfg.App.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := a.urlService.Cleanup(context.Background()); err != nil {
			log.Printf("failed to clean expired URLs: %v", err)
		}
	}
}

func (a *Applcation) startServer() {
	if err := a.e.Start(a.cfg.Server.Address); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (a *Applcation) shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	defer func() {
		if err := a.db.Close(); err != nil {
			log.Printf("close db err: %v", err)
		}
	}()

	defer func() {
		if err := a.redis.Close(); err != nil {
			log.Printf("close redis err: %v", err)
		}
	}()

	// 优雅的关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func (a *Applcation) initEcho() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Server.WriteTimeout = a.cfg.Server.WriteTimeout
	e.Server.ReadTimeout = a.cfg.Server.ReadTimeout
	e.Validator = validator.NewCustomValidator()
	a.e = e
}
