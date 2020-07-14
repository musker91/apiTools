package routers

import (
	"apiTools/controlers"
)

// 初始化api路由
func initApiRouter() {
	apiGroup := Router.Group("/api")
	// whois routers
	{
		apiGroup.GET("/whoisquery", controlers.WhoisQuery)
		apiGroup.POST("/whoisquery", controlers.WhoisQuery)
	}
	// short routers
	{
		// 长链接转换为短链接
		apiGroup.GET("/toshorturl", controlers.ShortToShortUrl)
		apiGroup.POST("/toshorturl", controlers.ShortToShortUrl)
		// 短链接解析回长链接
		apiGroup.GET("/parseshorturl", controlers.ShortParseShortUrl)
		apiGroup.POST("/parseshorturl", controlers.ShortParseShortUrl)
		// 标准官方短链接转换
		apiGroup.GET("/tooffshorturl", controlers.ShortToOfficial)
		apiGroup.POST("/tooffshorturl", controlers.ShortToOfficial)
		// 标准官方短链接解析
		apiGroup.GET("/parseoffshorturl", controlers.ShortParseOfficial)
		apiGroup.POST("/parseoffshorturl", controlers.ShortParseOfficial)

	}
	// ip query
	{
		apiGroup.GET("/ipv4query", controlers.Ipv4Query)
		apiGroup.POST("/ipv4query", controlers.Ipv4Query)
	}
	// proxy pool query
	{
		apiGroup.GET("/proxypool", controlers.ProxyPoolQuery)
		apiGroup.POST("/proxypool", controlers.ProxyPoolQuery)
	}
	// bank card info
	{
		apiGroup.GET("/bankcard", controlers.BankCardInfo)
		apiGroup.POST("/bankcard", controlers.BankCardInfo)
	}
	// icp info query
	{
		apiGroup.GET("/icpquery", controlers.ICPQueryInfo)
		apiGroup.POST("/icpquery", controlers.ICPQueryInfo)
	}
	// mobile telephone query
	{
		apiGroup.GET("/telquery", controlers.MobileTelQueryInfo)
		apiGroup.POST("/telquery", controlers.MobileTelQueryInfo)
	}
	// text to audio
	{
		apiGroup.GET("/text_to_audio", controlers.BdTextToAudio)
		apiGroup.POST("/text_to_audio", controlers.BdTextToAudio)
	}
	// qq info query
	{
		apiGroup.GET("/qqinfo", controlers.QQInfoQueryInfo)
		apiGroup.POST("/qqinfo", controlers.QQInfoQueryInfo)
	}
	// tencent domain check
	{
		apiGroup.GET("/dcheck", controlers.TxDomainCheck)
		apiGroup.POST("/dcheck", controlers.TxDomainCheck)
	}
	// segmentation cut
	{
		apiGroup.GET("/segcut", controlers.TextSegCut)
		apiGroup.POST("/segcut", controlers.TextSegCut)
	}

	// short video parse
	{
		apiGroup.GET("/svp", controlers.ShortVideoParse)
		apiGroup.POST("/svp", controlers.ShortVideoParse)
	}
}
