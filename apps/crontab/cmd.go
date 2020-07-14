package crontab

import "apiTools/utils"

// 运行定时任务APP
func RunCrontabApp() (err error) {
	// 初始化定时任务对象
	cronObj = utils.NewWithCron()

	// 创建定时任务
	go createCronTask()
	return
}
