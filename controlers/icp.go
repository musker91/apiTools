package controlers

import (
"apiTools/libs/logger"
"apiTools/modles"
"github.com/gin-gonic/gin"
"github.com/sirupsen/logrus"
"net/http"
)

func ICPQueryInfo(c *gin.Context) {
	var requestForm modles.ICPForm
	err := c.Bind(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "param has fail", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("get param ICPQueryInfo fields from icp info query")
		return
	}
	result, msg, err := modles.QueryICPInfo(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": msg, "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("query domain icp info fail")
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "data": result})
}
