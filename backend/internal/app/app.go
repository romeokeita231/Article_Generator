package app

import (
	"fmt"
	"log"
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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
	RedisClient *redis.Client

	// Handlers
	UserHandler    *handler.UserHandler
	ArticleHandler *handler.ArticleHandler
	HealthHandler  *handler.HealthHandler
	StatisticsHandler *handler.StatisticsHandler

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

	// 初始化 Redis
	redisClient, err := initRedis(cfg)
	if err != nil {
		return nil, fmt.Errorf("init redis: %w", err)
	}
	// 初始化各层
	userStore := store.NewUserStore(db)
	articleStore := store.NewArticleStore(db)
	agentLogStore := store.NewAgentLogStore(db)

	// 初始化 SSE 管理器
	sseManager := common.NewSSEManager()

	// 初始化服务层
	userService := service.NewUserService(userStore)
	quotaService := service.NewQuotaService(userStore)
	agentLogService := service.NewAgentLogService(agentLogStore)
	statisticsService := service.NewStatisticsService(db, userStore, articleStore, redisClient)

	// COS 服务（判断是否已配置）
	cosEnabled := cfg.COS.Bucket != "" && cfg.COS.SecretID != "" && cfg.COS.SecretKey != ""
	var cosService *service.CosService
	if cosEnabled {
		cosService = service.NewCosService(cfg.COS)
	}

	// 初始化所有图片服务
	pexelsService := service.NewPexelsService(cfg)
	iconifyService := service.NewIconifyService(cfg.Iconify)
	mermaidService := service.NewMermaidService(cfg.Mermaid)
	nanoBananaService := service.NewNanoBananaService(cfg.NanoBanana)
	svgDiagramService := service.NewSVGDiagramService(cfg.SVGDiagram, cfg.AI)
	emojiPackService := service.NewEmojiPackService(cfg.EmojiPack)
	picsumService := service.NewPicsumService() // 降级服务

	// 初始化策略并注册所有服务
	imageStrategy := service.NewImageServiceStrategy(cosService, cosEnabled)
	imageStrategy.RegisterService(pexelsService)
	imageStrategy.RegisterService(iconifyService)
	imageStrategy.RegisterService(mermaidService)
	imageStrategy.RegisterService(nanoBananaService)
	imageStrategy.RegisterService(svgDiagramService)
	imageStrategy.RegisterService(emojiPackService)
	imageStrategy.RegisterService(picsumService)

	log.Println("图片服务策略初始化完成，已注册 7 个图片服务（含降级服务）")

	// 智能体服务（注入 agentLogService）
	agentService, err := service.NewArticleAgentService(cfg, imageStrategy, agentLogService, sseManager)
	if err != nil {
		return nil, fmt.Errorf("init agent service: %w", err)
	}



	articleService := service.NewArticleService(articleStore, agentService, quotaService, sseManager)

	// 处理器层
	userHandler := handler.NewUserHandler(userService)
	articleHandler := handler.NewArticleHandler(articleService, userService, agentLogService, sseManager)
	healthHandler := handler.NewHealthHandler()
	statisticsHandler := handler.NewStatisticsHandler(statisticsService)

	return &App{
		Config:         cfg,
		DB:             db,
		UserHandler:    userHandler,
		ArticleHandler: articleHandler,
		HealthHandler:  healthHandler,
		StatisticsHandler: statisticsHandler,
		UserService:    userService,
		RedisClient:    redisClient,
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

func initRedis(cfg *config.Config) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     cfg.Redis.GetRedisAddr(),
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })

    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("redis ping: %w", err)
    }
    return client, nil
}

// Close 关闭资源
func (a *App) Close() error {
    sqlDB, _ := a.DB.DB()
    if err := sqlDB.Close(); err != nil {
        return err
    }
    if a.RedisClient != nil {
        a.RedisClient.Close()
    }
    return nil
}

