package crontab

import (
	"apiTools/libs/config"
	"apiTools/libs/logger"
	"apiTools/modles"
	"github.com/robfig/cron/v3"
)

var (
	cronObj *cron.Cron
)

const (
	extractProxyToRedisSpec = "*/30 * * * * ?"
)

func createCronTask() {
	var err error
	_, err = cronObj.AddJob(extractProxyToRedisSpec, &extractProxyToRedisTask{})
	if err != nil {
		logger.Echo.Errorf("create extracting proxy to redis task fail, err: %s", err)
	}
	logger.Echo.Info("create extracting proxy to redis task success")

	cronObj.Start()
	logger.Echo.Info("proxy app: cron task coroutine is running...")
}

// 提取代理数据到redis数据库
type extractProxyToRedisTask struct {
}

func (t extractProxyToRedisTask) Run() {
	redisProxyPools := config.GetRedisProxyPools()
	for _, redisConf := range redisProxyPools {
		go func(keyName string, checkUrl string) {
			for i := 0; i < 3; i++ {
				logger.Echo.Debugf("run extracting proxy to redis task(keyName: %s, checkUrl: %s)", keyName, checkUrl)
				status := modles.ExtractProxyToRedis(keyName, checkUrl)
				if status {
					break
				}
			}
		}(redisConf.KeyName, redisConf.CheckUrl)
	}
}
