package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
		//AllowAllOrigins: true,
		AllowOrigins: []string{"*"},                                       //允许跨域发来请求的网站
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}, //允许请求的方法
		//AllowHeaders:  []string{"Origin", "Authorization", "X-Requested-With", "Accept", "Content-Type"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

}

func (s *Server) registerApiRouter() {
	fileService := &service.FileService{}
	//静态文件
	s.engine.StaticFS(config.FileUploadConfig.StaticFsPath, http.Dir(config.FileUploadConfig.FileStorage))
	//repo
	userRoleRepo := repo.NewUserRoleRepo(s.transactionManger)
	roleResRepo := repo.NewRoleResRepo(s.transactionManger)
	userRepo := repo.NewUserRepo(s.transactionManger)
	roleRepo := repo.NewRoleRepo(s.transactionManger)
	resRepo := repo.NewResRepo(s.transactionManger)
	itemRepo := repo.NewItemRepo(s.transactionManger)
	itemSpecRepo := repo.NewItemSpecRepo(s.transactionManger)
	itemCommonRepo := repo.NewItemCommonRepo(s.transactionManger)
	paramRepo := repo.NewParamRepo(s.transactionManger)

	//Service
	resService := service.NewResService(resRepo, s.transactionManger)
	roleService := service.NewRoleService(userRepo, roleRepo, resRepo, userRoleRepo, roleResRepo, s.transactionManger)
	userService := service.NewUserService(userRepo, roleRepo, userRoleRepo, fileService, s.transactionManger)
	itemService := service.NewItemService(itemRepo, fileService, s.transactionManger)
	authService := service.NewAuthService(s.transactionManger, userRepo, roleRepo, resRepo, userRoleRepo, roleResRepo)
	authHelper := security.NewAuthHelper(authService)
	hasPerms := authHelper.HasPermission
	itemSpecService := service.NewItemSpecService(itemSpecRepo, fileService, s.transactionManger)
	paramService := service.NewParamService(paramRepo, s.transactionManger)
	itemCommonService := service.NewItemCommonService(itemCommonRepo, fileService, s.transactionManger)
	//Controller
	resCtrl := controller.NewResController(resService)
	roleCtrl := controller.NewRoleController(roleService)
	userCtrl := controller.NewUserController(userService, fileService)
	itemCtrl := controller.NewItemController(itemService)
	fileUploadCtr := controller.NewFileUploadController()
	itemSpecCtrl := controller.NewItemSpecController(itemSpecService, fileService)
	paramCtrl := controller.NewParamController(paramService)
	itemCommonCtrl := controller.NewItemCommonController(itemCommonService)
	//登录接口
	authController := controller.NewAuthController(userService, authService, authHelper)

	admin := s.engine.Group("/api/admin")
	admin.POST("/login", authController.Login)

	adminApi := admin.Group("/svc", security.JWTAuthMiddleware())
	{
		adminApi.POST("/file-upload", fileUploadCtr.UploadFile)
		//当前用户
		adminApi.GET("/profile/current-user", authController.CurrentUser)
		adminApi.POST("/profile/update-profile", authController.UpdateUserProfile)
		adminApi.POST("/profile/change-password", authController.ChangePassword)
		//用户

		adminApi.GET("/user", hasPerms("system:user"), userCtrl.Query)
		adminApi.GET("/user/:userId", hasPerms("system:user"), userCtrl.Find)
		adminApi.POST("/user", hasPerms("system:user:add"), userCtrl.Create)
		adminApi.PUT("/user", hasPerms("system:user:update"), userCtrl.Update)
		adminApi.DELETE("/user", hasPerms("system:user:delete"), userCtrl.Delete)
		adminApi.POST("/user/enable", hasPerms("system:user:enable"), userCtrl.EnableUser)
		adminApi.POST("/user/password", hasPerms("system:user:password"), userCtrl.ChangePwd)

		//资源
		adminApi.GET("/res", hasPerms("system:res"), resCtrl.Query)
		adminApi.GET("/res/:resId", hasPerms("system:res"), resCtrl.Find)
		adminApi.POST("/res", hasPerms("system:res:add"), resCtrl.Create)
		adminApi.PUT("/res", hasPerms("system:res:update"), resCtrl.Update)
		adminApi.DELETE("/res", hasPerms("system:res:delete"), resCtrl.Delete)

		//角色

		adminApi.GET("/role", hasPerms("system:role"), roleCtrl.Query)
		adminApi.GET("/role/:roleId", hasPerms("system:role"), roleCtrl.Find)
		adminApi.POST("/role", hasPerms("system:role:add"), roleCtrl.Create)
		adminApi.PUT("/role", hasPerms("system:role:update"), roleCtrl.Update)
		adminApi.DELETE("/role", hasPerms("system:role:delete"), roleCtrl.Delete)
		adminApi.POST("/role-enable", hasPerms("system:role:enable"), roleCtrl.EnabledRole)
		adminApi.POST("/role-user", hasPerms("system:role:user:add"), roleCtrl.CreateUsersOfRole)
		adminApi.DELETE("/role-user", hasPerms("system:role:user:delete"), roleCtrl.DeleteUserRole)
		adminApi.GET("/role-user", hasPerms("system:role:user"), roleCtrl.QueryUserOfRole)

		//商品
		adminApi.GET("/item", hasPerms("prod:item"), itemCtrl.QueryItem)
		adminApi.GET("/item/:itemId", hasPerms("prod:item"), itemCtrl.FindItem)
		adminApi.POST("/item", hasPerms("prod:item:add"), itemCtrl.CreateItem)
		adminApi.PUT("/item", hasPerms("prod:item:update"), itemCtrl.UpdateItem)
		adminApi.DELETE("/item", hasPerms("prod:item:delete"), itemCtrl.DeleteItems)

		//商品规格
		adminApi.GET("/item-spec", hasPerms("prod:spec"), itemSpecCtrl.QuerySpecs)
		adminApi.GET("/item-spec/:specId", hasPerms("prod:spec"), itemSpecCtrl.FindSpec)
		adminApi.POST("/item-spec", hasPerms("prod:spec:add"), itemSpecCtrl.CreateSpec)
		adminApi.PUT("/item-spec", hasPerms("prod:spec:update"), itemSpecCtrl.UpdateSpec)
		adminApi.DELETE("item-spec", hasPerms("prod:spec:delete"), itemSpecCtrl.DeleteSpecs)
		adminApi.GET("/item-spec-color", hasPerms("prod:spec:color"), itemSpecCtrl.GetColorsOfSpec)
		adminApi.GET("/item-spec-color/:colorId", hasPerms("prod:spec:color"), itemSpecCtrl.FindSpecColor)
		adminApi.POST("/item-spec-color", hasPerms("prod:spec:color:add"), itemSpecCtrl.CreateSpecColor)
		adminApi.PUT("/item-spec-color", hasPerms("prod:spec:color:update"), itemSpecCtrl.UpdateSpecColor)
		adminApi.DELETE("/item-spec-color", hasPerms("prod:spec:color:delete"), itemSpecCtrl.DeleteSpecColor)
		adminApi.POST("/item-spec-media", hasPerms("prod:spec:media:add"), itemSpecCtrl.UploadSpecMedia)
		adminApi.DELETE("/item-spec-media", hasPerms("prod:spec:media:delete"), itemSpecCtrl.DeleteSpecMedia)
		//系统参数
		adminApi.GET("/sys-param", hasPerms("system:param"), paramCtrl.QueryParam)
		adminApi.GET("/sys-param/:paramId", hasPerms("system:param"), paramCtrl.FindParam)
		adminApi.GET("/sys-param-all", hasPerms("system:param:all"), paramCtrl.GetAllParamsWithItems)
		adminApi.GET("/sys-param-with-item/:paramKey", hasPerms("system:param:item:key"), paramCtrl.GetParamWithItemsByKey)
		adminApi.POST("/sys-param", hasPerms("system:param:add"), paramCtrl.CreateParam)
		adminApi.PUT("/sys-param", hasPerms("system:param:update"), paramCtrl.UpdateParam)
		adminApi.DELETE("/sys-param", hasPerms("system:params:delete"), paramCtrl.DeleteParams)
		adminApi.GET("/sys-param-item", hasPerms("system:param:item"), paramCtrl.GetItemOfParam)
		adminApi.GET("/sys-param-item/:itemId", hasPerms("system:param:item"), paramCtrl.FindItem)
		adminApi.POST("/sys-param-item", hasPerms("system:param:item:add"), paramCtrl.CreateItem)
		adminApi.PUT("/sys-param-item", hasPerms("system:param:item:update"), paramCtrl.UpdateItem)
		adminApi.DELETE("/sys-param-item", hasPerms("system:params:item:delete"), paramCtrl.DeleteItems)
		//品牌管理
		adminApi.GET("/brand", hasPerms("brand"), itemCommonCtrl.QueryBrands)
		adminApi.GET("/brand/:brandId", hasPerms("brand"), itemCommonCtrl.FindBrand)
		adminApi.POST("/brand", hasPerms("brand:add"), itemCommonCtrl.CreateBrand)
		adminApi.PUT("/brand", hasPerms("brand:update"), itemCommonCtrl.UpdateBrand)
		adminApi.DELETE("/brand", hasPerms("brand:delete"), itemCommonCtrl.DeleteBrands)
		adminApi.GET("/brand-all", hasPerms("brand"), itemCommonCtrl.GetAllBrandsWithSeries)
		//车系管理
		adminApi.GET("/series", hasPerms("series"), itemCommonCtrl.QuerySeries)
		adminApi.GET("/series/:seriesId", hasPerms("series"), itemCommonCtrl.FindSeries)
		adminApi.POST("/series", hasPerms("series:add"), itemCommonCtrl.CreateSeries)
		adminApi.PUT("/series", hasPerms("series:update"), itemCommonCtrl.UpdateSeries)
		adminApi.DELETE("/series", hasPerms("series:delete"), itemCommonCtrl.DeleteSeries)

	}

}
