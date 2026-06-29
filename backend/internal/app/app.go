package app

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"

	"github.com/romeokeita231/Article_Generator/internal/common"
	"github.com/romeokeita231/Article_Generator/internal/config"
	"github.com/romeokeita231/Article_Generator/internal/handler"
	"github.com/romeokeita231/Article_Generator/internal/service"
	"github.com/romeokeita231/Article_Generator/internal/store"
)

// App 应用程序
type App struct {
	Config *config.Config
	DB     *gorm.DB

	// Handlers
	UserHandler   *handler.UserHandler
    ArticleHandler *handler.ArticleHandler
	HealthHandler *handler.HealthHandler

	// Services (用于中间件)
	UserService *service.UserService
}

// New 创建应用实例
func New(cfg *config.Config) (*App, error) {
	// 初始化数据库
	db, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	// 初始化各层
	userStore := store.NewUserStore(db)
	articleStore := store.NewArticleStore(db)
	sseManager := common.NewSSEManager()

	userService := service.NewUserService(userStore)
	quotaService := service.NewQuotaService(userStore)
	pexelsService := service.NewPexelsService(cfg)
	cosService := service.NewCosService()

    
	userHandler := handler.NewUserHandler(userService)
	healthHandler := handler.NewHealthHandler()

	agentService, err := service.NewArticleAgentService(cfg, pexelsService, cosService, sseManager)
	if err != nil {
		return nil, fmt.Errorf("init agent service: %w", err)
	}

    articleService := service.NewArticleService(articleStore, agentService, quotaService, sseManager)

	return &App{
		Config:        cfg,
		DB:            db,
		UserHandler:   userHandler,
        ArticleHandler: handler.NewArticleHandler(articleService, userService, sseManager),
		HealthHandler: healthHandler,
		UserService:   userService,
	}, nil
}

// initDB 初始化数据库
func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	log.Println("database connected")
	return db, nil
}

// Close 关闭资源
func (a *App) Close() error {
	sqlDB, err := a.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
