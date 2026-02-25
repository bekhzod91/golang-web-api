package core

import (
	"fmt"
	"github.com/go-chi/cors"
	"github.com/hzmat24/api/pkg/aws"
	"github.com/hzmat24/api/pkg/multi_tenency"
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	goredis "github.com/redis/go-redis/v9"

	configPkg "github.com/hzmat24/api/config"
	loggerPkg "github.com/hzmat24/api/pkg/logger"
	redisPkg "github.com/hzmat24/api/pkg/redis"
)

var BRAND = `
 _   _  ______ __   __   ___   _____  ____      _   
| | | |(___  /|  \ /  | / _ \ (_   _)(___ \   /  |  
| |_| |   / / |   v   || |_| |  | |    __) ) / o |_ 
|  _  |  / /  | |\_/| ||  _  |  | |   / __/ /__   _)
| | | | / /__ | |   | || | | |  | |  | |___    | |  
|_| |_|/_____)|_|   |_||_| |_|  |_|  |_____)   |_|
`

type App struct {
	mux         IMux
	config      configPkg.Config
	logger      *httplog.Logger
	db          *multi_tenency.DB
	redisClient *goredis.Client
	awsClient   *aws.Client
}

func NewApp() *App {
	mux := NewMux()
	config := configPkg.NewConfig()
	logger := loggerPkg.NewLogger(config)

	redisClient, err := redisPkg.NewClient(config)
	if err != nil {
		panic(err)
	}

	db, err := multi_tenency.NewDB(config)
	if err != nil {
		panic(err)
	}

	awsClient, err := aws.NewClient(config)
	if err != nil {
		panic(err)
	}

	app := &App{
		mux:         mux,
		config:      config,
		logger:      logger,
		redisClient: redisClient,
		db:          db,
		awsClient:   awsClient,
	}

	app.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"https://*", "http://*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowOriginFunc:    func(r *http.Request, origin string) bool { return true },
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-Tenant"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
		MaxAge:             300, // Maximum value not ignored by any of major browsers
	}))

	app.mux.Use(AppMiddleware(app))
	app.mux.Use(httplog.RequestLogger(logger))
	app.mux.Use(chimiddleware.RealIP)
	app.mux.Use(chimiddleware.Recoverer)

	return app
}

func (a *App) MountPublicRouter(fn func(IMux)) {
	fn(a.mux)
}

func (a *App) MountTenantRouter(fn func(IMux)) {
	a.mux.Group(func(r IMux) {
		r.Use(multi_tenency.MultiTenancy(a.db))
		fn(r)
	})
}

func (a *App) Run() {
	fmt.Println(BRAND)
	fmt.Printf("ENVIRONMENT: %s\n", a.config.Environment)
	fmt.Printf("-----------------------------------------------\n")
	fmt.Printf("POSTGRES_HOST: %s\n", a.config.PostgresHost)
	fmt.Printf("POSTGRES_PORT: %d\n", a.config.PostgresPort)
	fmt.Printf("POSTGRES_USER: %s\n", a.config.PostgresUser)
	fmt.Printf("POSTGRES_PASSWORD: ******************\n")
	fmt.Printf("POSTGRES_DATABASE: %s\n", a.config.PostgresDatabase)
	fmt.Printf("POSTGRES_SCHEMA: %s\n", a.config.PostgresSchema)
	fmt.Printf("POSTGRES_TLS: %t\n", a.config.PostgresTLSEnabled)
	fmt.Printf("-----------------------------------------------\n")
	fmt.Printf("REDIS_HOST: %s\n", a.config.RedisHost)
	fmt.Printf("REDIS_PORT: %d\n", a.config.RedisPort)
	fmt.Printf("REDIS_USER: %s\n", a.config.RedisUser)
	fmt.Printf("REDIS_PASSWORD: ******************\n")
	fmt.Printf("REDIS_TLS: %t\n", a.config.RedisTLSEnabled)
	fmt.Printf("-----------------------------------------------\n")
	fmt.Printf("Server run: http://%s:%d\n", a.config.APIHost, a.config.APIPort)

	addr := fmt.Sprintf("%s:%d", a.config.APIHost, a.config.APIPort)
	err := http.ListenAndServe(addr, a.mux)
	if err != nil {
		panic(err)
	}
}
