package middlreware

import (
	"apiTools/libs/config"
	"apiTools/libs/logger"
	"apiTools/modles"
	"apiTools/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

// 允许跨域中间件
func AllowCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

// 自定义日志输出中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		//执行时间
		useTime := time.Since(start)
		// 日志打印。分别是状态码，执行时间,请求客户端ip,请求方法,请求url
		logger.Echo.WithFields(logrus.Fields{
			"statusCode": c.Writer.Status(),
			"useTime":    fmt.Sprintf("%v", useTime),
			"clientIp":   c.ClientIP(),
			"method":     c.Request.Method,
			"urlPath":    c.Request.URL.Path,
			"queryParam": c.Request.URL.RawQuery,
		}).Info()
	}
}

func ProApiDocs() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlpath := c.Request.URL.Path
		urlPathSlice := strings.Split(urlpath, "/")
		// url不匹配
		if len(urlPathSlice) < 3 {
			c.Next()
			return
		} else if len(urlPathSlice) == 3 && strings.HasSuffix(urlpath, "/") {
			c.Next()
			return
		}
		// 获取字段
		rk := urlPathSlice[1]
		urlPath := urlPathSlice[2]
		// 获取json data
		jsonData, ok := modles.JsonData.(map[string]interface{})
		if !ok {
			c.Next()
			return
		}
		// 获取json data中api data
		apiData, ok := jsonData[urlPath].(map[string]interface{})
		// 没有匹配到，说明不是访问的api相关，直接进入路由函数
		if !ok {
			if rk == "api" || rk == "docs" {
				c.String(http.StatusNotFound, "404 page not found")
				c.Abort()
			} else {
				c.Next()
				return
			}
		}

		// 判断当前是否被禁用, 禁止就终止，返回404
		if apiData["enable"] == false {
			c.String(http.StatusNotFound, "404 page not found")
			c.Abort()
		}
		c.Set("urlPath", urlPath)
		countKey, countKeyOk := apiData["countKey"]
		docFile, docFileOk := apiData["docFile"]
		// 根据路由处理
		if rk == "api" {
			// 判断当前接口是否在维护中
			if apiData["mainten"] == true {
				c.JSON(http.StatusBadRequest, gin.H{"code": 1, "data": gin.H{}, "msg": "Interface maintenance"})
				c.Abort()
			}
			if countKeyOk {
				c.Set("countKey", countKey)
				c.Set("docFile", docFile)
				// redis自增
				countKeyName := fmt.Sprintf("apiCount_%s", countKey)
				redisClient := modles.RedisPool.Get()
				defer redisClient.Close()
				_, _ = redisClient.Do("INCR", countKeyName)
			}
		} else if rk == "docs" {
			if docFileOk {
				c.Set("countKey", countKey)
				c.Set("docFile", docFile)
				c.Set("mainten", apiData["mainten"])
			}
			titleName, ok := apiData["titleName"]
			if ok {
				c.Set("titleName", titleName)
			}
		}
		c.Next()
	}
}

// api接口访问限流
func IpLimiting() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		urlPathSlice := strings.Split(urlPath, "/")
		var rk string
		if len(urlPathSlice) >= 2 {
			rk = urlPathSlice[1]
		}
		if rk != "api" {
			c.Next()
			return
		}
		accessDenied := gin.H{
			"code": 1,
			"data": gin.H{},
			"msg":  "access denied",
		}
		// 验证token
		token := c.Query("token")
		if token != "" {
			tokens, err := modles.GetTokensFromCache()
			if err == nil {
				if utils.IsInSlice(token, tokens) {
					c.Next()
					return
				} else {
					c.JSON(http.StatusForbidden, accessDenied)
					c.Abort()
					return
				}
			}
		}
		// 无token
		var ipLimitingTimeSeconds = config.GetInt("web::ipLimitingTimeSeconds") // IP限流时间段, 秒
		var ipLimitingCount = config.GetInt("web::ipLimitingCount")             // IP限流时间段位内请求不能超过的次数
		var liftIpLimiting = config.GetInt("web::liftIpLimiting")               // 解除ip限流的时间, 秒
		clientIp := c.ClientIP()
		redisClient := modles.RedisPool.Get()
		defer redisClient.Close()
		// 获取redis中存储的访问次数
		count, _ := redis.Int(redisClient.Do("GET", clientIp))
		if count >= ipLimitingCount {
			if count == ipLimitingCount {
				_, _ = redisClient.Do("SET", clientIp, "999", "EX", liftIpLimiting)
			}
			c.JSON(http.StatusForbidden, accessDenied)
			c.Abort()
		} else {
			if count == 0 {
				_, _ = redisClient.Do("SET", clientIp, 1, "EX", ipLimitingTimeSeconds)
			} else {
				_, _ = redisClient.Do("INCR", clientIp)
			}
		}
		c.Next()
	}
}
