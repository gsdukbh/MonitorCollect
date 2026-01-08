package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func TaskRun() {
	// cron test
	c := cron.New()
	log.Printf("collectDisposeHour 启动，等待任务调度...")

	if config.Cron.Enable {
		//  使用配置文件  config.Cron.ScheduleDispos
		_, err := c.AddFunc(config.Cron.ScheduleDispos, collectDisposeHour)
		if err != nil {
			return
		}
		_, err = c.AddFunc("2 0 * * *", clearDisk)
		if err != nil {
			return
		} // 每天凌晨清理磁盘数据
		_, err = c.AddFunc("3 1 * * *", clearCpu)
		if err != nil {
			return
		} // 每天凌晨1点清理CPU数据
		_, err = c.AddFunc("4 2 * * *", clearMem)
		if err != nil {
			return
		} // 每天凌晨2点清理内存数据
	}
	c.Start()
}

// aggKey 用于聚合数据的键
type aggKey struct {
	Host      string
	Interface string
	Hour      time.Time
}

// trafficStats 用于临时存储统计数据
type trafficStats struct {
	MinRecv int64
	MaxRecv int64
	MinSent int64
	MaxSent int64
	Count   int64
}

// collectDisposeHour 定时清理任务，清理过期数据、按时间段统计网络流量存储。
// 按小时统计网络流量数据，存储到对应的表中。
// 流量存储单位为 MB，速度为 MB/s
func collectDisposeHour() {
	log.Printf("collectDisposeHour 执行中...")
	// 获取过去 30 天的数据处理起止时间
	columbina := -24 * 1 * time.Hour
	startTime := time.Now().Add(columbina).Unix()

	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		log.Printf("Failed to start transaction: %v", tx.Error)
		return
	}

	// 1. 获取数据
	rawData, err := fetchRawData(tx, startTime)
	if err != nil {
		log.Printf("Failed to fetch raw data: %v", err)
		tx.Rollback()
		return
	}

	if len(rawData) == 0 {
		log.Println("No data to process.")
		tx.Rollback()
		return
	}

	// 2. 聚合统计
	statsMap := aggregateTrafficStats(rawData)

	// 3. 转换为目标数据结构
	hourData := prepareHourData(statsMap)

	// 4. 保存统计结果
	if len(hourData) > 0 {
		if err := saveHourData(tx, hourData); err != nil {
			log.Printf("Failed to insert hour data: %v", err)
			tx.Rollback()
			return
		}
	}

	// 5. 删除原始数据
	if err := deleteRawData(tx, startTime); err != nil {
		log.Printf("Failed to delete processed data: %v", err)
		tx.Rollback()
		return
	}

	// 6. 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return
	}

	log.Printf("collectDisposeHour 成功: 处理了 %d 条原始记录，生成 %d 条小时记录。", len(rawData), len(hourData))
}

// fetchRawData 从数据库获取原始数据
func fetchRawData(tx *gorm.DB, startTime int64) ([]NetInterfaceFieldsDb, error) {
	var rawData []NetInterfaceFieldsDb
	if err := tx.Where("timestamp > ?", startTime).Find(&rawData).Error; err != nil {
		return nil, err
	}
	return rawData, nil
}

// aggregateTrafficStats 对原始数据进行聚合计算
func aggregateTrafficStats(rawData []NetInterfaceFieldsDb) map[aggKey]*trafficStats {
	statsMap := make(map[aggKey]*trafficStats)

	for _, record := range rawData {
		ts := time.Unix(record.Timestamp, 0)
		hourTime := time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), 0, 0, 0, ts.Location())

		k := aggKey{
			Host:      record.Host,
			Interface: record.Interface,
			Hour:      hourTime,
		}

		if _, exists := statsMap[k]; !exists {
			statsMap[k] = &trafficStats{
				MinRecv: record.BytesRecv,
				MaxRecv: record.BytesRecv,
				MinSent: record.BytesSent,
				MaxSent: record.BytesSent,
				Count:   0,
			}
		}

		s := statsMap[k]
		if record.BytesRecv < s.MinRecv {
			s.MinRecv = record.BytesRecv
		}
		if record.BytesRecv > s.MaxRecv {
			s.MaxRecv = record.BytesRecv
		}
		if record.BytesSent < s.MinSent {
			s.MinSent = record.BytesSent
		}
		if record.BytesSent > s.MaxSent {
			s.MaxSent = record.BytesSent
		}

		s.Count++
	}
	return statsMap
}

// prepareHourData 将聚合结果转换为数据库模型
func prepareHourData(statsMap map[aggKey]*trafficStats) []NetInterfaceCollectHour {
	var hourData []NetInterfaceCollectHour
	for k, s := range statsMap {
		totalBytes := (s.MaxRecv - s.MinRecv) + (s.MaxSent - s.MinSent)
		// 如果计数器重置，避免负值
		if totalBytes < 0 {
			totalBytes = 0
		}

		// 计算平均速度 (bits/s)
		// 1 byte = 8 bits
		// 1 hour = 3600 seconds
		avgBps := float64(totalBytes*8) / 3600.0

		// 格式化速度字符串
		speedStr := formatNetSpeed(avgBps)

		// 存入 MB 供参考
		totalMB := totalBytes / 1024 / 1024

		item := NetInterfaceCollectHour{
			Host:      k.Host,
			Interface: k.Interface,
			Hour:      k.Hour,
			Total:     totalMB,
			Speed:     float64(int64(avgBps/1000000.0*100+0.5)) / 100.0, // 保留两位小数 (Mbps)
			SpeedStr:  speedStr,
		}

		hourData = append(hourData, item)
	}
	return hourData
}

// saveHourData 批量保存小时统计数据
func saveHourData(tx *gorm.DB, data []NetInterfaceCollectHour) error {
	// 按批写入，避免一次性插入大量数据导致内存或事务压力
	return tx.CreateInBatches(data, 100).Error
}

// deleteRawData 删除已处理的原始数据,今天之前的数据
func deleteRawData(tx *gorm.DB, startTime int64) error {
	return tx.Where("timestamp < ?", startTime).Delete(&NetInterfaceFieldsDb{}).Error
}

// formatNetSpeed 将 bits/s 转换为人类可读的字符串 (Kbps, Mbps, Gbps)
func formatNetSpeed(bps float64) string {
	if bps >= 1000*1000*1000 {
		return fmt.Sprintf("%.2f Gbps", bps/1000/1000/1000)
	} else if bps >= 1000*1000 {
		return fmt.Sprintf("%.2f Mbps", bps/1000/1000)
	} else if bps >= 1000 {
		return fmt.Sprintf("%.2f Kbps", bps/1000)
	}
	return fmt.Sprintf("%.2f bps", bps)
}

// clearDisk 清理过期磁盘数据
func clearDisk() {
	log.Printf("clearDisk 执行中...")
	columbina := -24 * 30 * time.Hour
	startTime := time.Now().Add(columbina).Unix()

	// 删除过期数据
	if err := db.Where("timestamp < ?", startTime).Delete(&DiskFieldsDb{}).Error; err != nil {
		log.Printf("Failed to clear old disk data: %v", err)
		return
	}
}

// clearCpu 清理过期 CPU 数据
func clearCpu() {
	log.Printf("clearCpu 执行中...")
	columbina := -24 * 30 * time.Hour
	startTime := time.Now().Add(columbina).Unix()

	// 删除过期数据
	if err := db.Where("timestamp < ?", startTime).Delete(&CPUFieldsDb{}).Error; err != nil {
		log.Printf("Failed to clear old CPU data: %v", err)
		return
	}
}

// clearMem 清理过期内存数据
func clearMem() {
	log.Printf("clearMem 执行中...")
	columbina := -24 * 30 * time.Hour
	startTime := time.Now().Add(columbina).Unix()

	// 删除过期数据
	if err := db.Where("timestamp < ?", startTime).Delete(&MemFieldsDb{}).Error; err != nil {
		log.Printf("Failed to clear old Mem data: %v", err)
		return
	}
}
