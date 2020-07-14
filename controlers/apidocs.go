package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Api文档首页
func ApiIndex(c *gin.Context) {
	allApiInfo, err := modles.GetAllApiInfo()
	if err != nil {
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"info":    allApiInfo,
		}).Error("access api index page fail")
		c.HTML(http.StatusBadGateway, "error.html", gin.H{"errorCode": "503", "errorMsg": "系统异常请稍后重试"})
		return
	}
	c.HTML(http.StatusOK, "index.html", allApiInfo)
}

// Api文档详情页
func ApiDocs(c *gin.Context) {
	var apiDocFileName, urlPath, countKey string

	// 获取数据
	reqApiName := c.Param("apiName")
	urlPath = c.GetString("urlPath")
	if reqApiName != urlPath {
		urlPath = ""
	}
	apiDocFileName = c.GetString("docFile")
	countKey = c.GetString("countKey")
	mainten := c.GetBool("mainten")

	if apiDocFileName == "" || urlPath == "" || countKey == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}
	apiDocInfo, err := modles.GetApiDocInfo(apiDocFileName, urlPath, countKey, mainten)
	if err != nil {
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"info":    apiDocInfo,
		}).Error("access api doc page fail")
		c.HTML(http.StatusBadGateway, "error.html", gin.H{"errorCode": "503", "errorMsg": "系统异常请稍后重试"})
		return
	}
	c.HTML(http.StatusOK, "docs.html", apiDocInfo)
}

// 关于ApiTools
func ApiAbout(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", nil)
}

// api 访问白名单限制申请
func VisitAppli(c *gin.Context) {
	method := c.Request.Method
	switch method {
	case "GET":
		c.HTML(http.StatusOK, "token_apply.html", nil)
	case "POST":
		var form modles.UserVisitAppli
		err := c.ShouldBind(&form)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "request fail"})
			return
		}
		err = modles.CreateUserVisit(&form)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "request fail"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": ""})
		return
	}
}
