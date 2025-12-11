package main

import (
	"log"

	"github.com/robfig/cron/v3"
)

func TaskRun() {
	// cron test
	c := cron.New(cron.WithSeconds())
	log.Printf("collectDispose 启动，等待任务调度...")

	if config.Cron.Enable {
		// todo 使用配置文件  config.Cron.ScheduleDispos
		_, err := c.AddFunc("*/10 * * * * *", collectDispose)
		if err != nil {
			return
		}
	}
	c.Start()
}

// collectDispose 定时清理任务，清理过期数据、按时间段统计网络流量存储。
// 按小时，天，周，月统计网络流量数据，存储到对应的表中。
func collectDispose() {

}
