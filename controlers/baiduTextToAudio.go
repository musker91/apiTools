package controlers

import (
	"apiTools/libs/logger"
	"apiTools/modles"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func BdTextToAudio(c *gin.Context) {
	var requestForm modles.AudioForm
	err := c.Bind(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "param has fail", "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("get param failed from baidu text to audio info query")
		return
	}
	resp, err := modles.BdTextToAudio(&requestForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": resp.Msg, "data": gin.H{}})
		logger.Echo.WithFields(logrus.Fields{
			"routers": c.Request.URL.Path,
			"err":     err,
			"query":   requestForm,
		}).Error("query bank card info fail")
		return
	}
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s.mp3"`, resp.FileName),
	}
	if resp.ConTextType == "audio/x-bd-bv" {
		c.DataFromReader(http.StatusOK, resp.ContentLength, resp.ConTextType, bytes.NewReader(resp.Data), extraHeaders)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": resp.Msg})
	}
}
