package routers

import (
	"apiTools/controlers"
	"apiTools/libs/config"
	"apiTools/routers/middlreware"
	"apiTools/utils"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

var (
	Router *gin.Engine
)

// 初始化gin
func InitRouter() {
	if config.GetString("web::appMode") == "production" {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}
	//Router = gin.Default()
	Router = gin.New()

	// pprof 性能分析
	if config.GetBool("web::enablePprof") {
		pprof.Register(Router)
	}

	// 设置静态文件
	Router.LoadHTMLGlob(filepath.Join(utils.GetRootPath(), "views", "/*"))
	Router.Static("/static", filepath.Join(utils.GetRootPath(), "static"))

	// 设置全局中间件
	Router.Use(gin.Recovery())
	Router.Use(middlreware.AllowCors())
	Router.Use(middlreware.Logger())
	Router.Use(middlreware.ProApiDocs())

	if config.GetBool("web::enableIpLimiting") {
		Router.Use(middlreware.IpLimiting())
	}

	// 404错误处理
	Router.NoRoute(controlers.NoRouter)
	Router.NoMethod(controlers.NoRouter)

	// 加载路由
	initApiRouter()
	initControlers()
}
