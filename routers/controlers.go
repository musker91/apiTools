package routers

import "apiTools/controlers"

// 初始化路由
func initControlers() {
	Router.GET("/", controlers.ApiIndex)
	Router.GET("/docs/:apiName", controlers.ApiDocs)
	Router.GET("/about", controlers.ApiAbout)
	Router.GET("/visitappli", controlers.VisitAppli)
	Router.POST("/visitappli", controlers.VisitAppli)
}
