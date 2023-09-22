package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"used-car-deal-gobackend/base/config"
	"used-car-deal-gobackend/base/datasource"
	"used-car-deal-gobackend/base/logger"
	"used-car-deal-gobackend/base/security"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/controller"
	"used-car-deal-gobackend/repo"
	"used-car-deal-gobackend/service"
)

func NewServer() (*Server, error) {
	server := &Server{}
	return server, nil
}

type Server struct {
	transactionManger *datasource.TransactionManger
	datasource        datasource.DBCommon
	beans             map[string]any
	engine            *gin.Engine
}

func (s *Server) init() {
	s.beans = make(map[string]any)
	//1.初始化配置
	s.loadConfig()
	//2.初始化数据库
	s.initDataSource()
	//4.初始化路由
	s.initRouter()
	s.registerApiRouter()

}
func (s *Server) Start() {
	s.init()
	s.engine.Run(":" + config.AppConfig.Port)
}

func (s *Server) loadConfig() {
	config.InitConfig()
	//初始化日志处理
	err := logger.InitLogger()
	if err != nil {
		panic(fmt.Sprintf("init logger failed, err:%v\n", err))
	}
}

func (s *Server) initDataSource() {
	dsn := config.MySQLConfig.Username + ":" + config.MySQLConfig.Password + "@" + config.MySQLConfig.Url
	var err error
	s.datasource, err = datasource.NewMysqlConn(dsn)
	if err != nil {
		panic(fmt.Errorf("初始化数据源错误:%v", err))
	}

	s.transactionManger = datasource.NewTransactionManger(s.datasource)

}

func (s *Server) initRouter() {
	s.engine = gin.New()

	s.engine.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 统一异常处理
	s.engine.Use(web.ErrorHandleMiddleware)
	s.engine.Use(datasource.TransactionMiddleware)

	//跨域请求
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://foo.com"},                        //允许跨域发来请求的网站
		AllowMethods:  []string{"GET", "POST", "DELETE", "PUT", "OPTION"}, //允许请求的方法
		AllowHeaders:  []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
	}))

}

func (s *Server) registerApiRouter() {
	//用户
	userRepo := repo.NewUserRepo(s.transactionManger)
	userService := service.NewUserService(userRepo, s.transactionManger)
	userCtrl := controller.NewUserController(userService)
	//登录接口
	loginController := controller.NewLoginController(userService)
	s.engine.POST("/login", loginController.Login)

	api := s.engine.Group("/api/", security.JWTAuthMiddleware())
	{
		//用户

		api.GET("/users", userCtrl.Query)
		api.GET("/user/:userId", userCtrl.Find)
		api.POST("/user", userCtrl.Create)
		api.PUT("/user", userCtrl.Update)
		api.DELETE("/user", userCtrl.Delete)

	}

}
