package proxyPool

import (
	"apiTools/libs/logger"
	"apiTools/utils"
)

const (
	checkDBProxyTaskSpec = "0 */3 * * * ?"
)

var (
	cronObj = utils.NewWithCron()
)

// 检测数据库代理信息任务
type checkDBProxyTask struct {
}

func (t checkDBProxyTask) Run() {
	logger.Echo.Debug("run checkDBProxyTask")
	checkDBProxy()
}

func runCrontab() {
	_, err := cronObj.AddJob(checkDBProxyTaskSpec, &checkDBProxyTask{})
	if err != nil {
		logger.Echo.Errorf("create check db proxy task fail, err: %s", err)
		return
	}
	cronObj.Start()
	logger.Echo.Info("create check db proxy task success")
	return
}
