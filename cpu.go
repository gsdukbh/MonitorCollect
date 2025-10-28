package main

import (
	"encoding/json"
	"log"
)

// CPUFields 表示 CPU 使用情况统计
// 处理telegraf 采集的CPU数据
// 对应样例中 name=cpu 的字段
type CPUFields struct {
	UsageActive    float64 `json:"usage_active"`     // CPU 活跃时间百分比
	UsageGuest     float64 `json:"usage_guest"`      // 运行虚拟 CPU 的时间百分比
	UsageGuestNice float64 `json:"usage_guest_nice"` // 运行低优先级虚拟 CPU 的时间百分比
	UsageIdle      float64 `json:"usage_idle"`       // CPU 空闲时间百分比
	UsageIowait    float64 `json:"usage_iowait"`     // 等待 I/O 完成的时间百分比
	UsageIrq       float64 `json:"usage_irq"`        // 处理硬件中断的时间百分比
	UsageNice      float64 `json:"usage_nice"`       // 运行低优先级进程的时间百分比
	UsageSoftirq   float64 `json:"usage_softirq"`    // 处理软件中断的时间百分比
	UsageSteal     float64 `json:"usage_steal"`      // 虚拟化环境中被其他虚拟机占用的时间百分比
	UsageSystem    float64 `json:"usage_system"`     // 内核态时间百分比
	UsageUser      float64 `json:"usage_user"`       // 用户态时间百分比
}

// CPUFieldsDb 用于数据库存储的 CPU 字段结构体
// 使用 gorm 标签定义数据库字段映射
type CPUFieldsDb struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`           // 主键ID，自增
	CPU       string `gorm:"type:varchar(50);not null;index" json:"cpu"`   // CPU 标识符，如 "cpu0", "cpu1", "cpu-total"
	Host      string `gorm:"type:varchar(100);not null;index" json:"host"` // 主机名
	Timestamp int64  `gorm:"not null;index" json:"timestamp"`              // 时间戳（纳秒或毫秒）

	// CPU 使用情况字段
	UsageActive    float64 `gorm:"type:decimal(10,6);not null" json:"usage_active"`     // CPU 活跃时间百分比
	UsageGuest     float64 `gorm:"type:decimal(10,6);not null" json:"usage_guest"`      // 运行虚拟 CPU 的时间百分比
	UsageGuestNice float64 `gorm:"type:decimal(10,6);not null" json:"usage_guest_nice"` // 运行低优先级虚拟 CPU 的时间百分比
	UsageIdle      float64 `gorm:"type:decimal(10,6);not null" json:"usage_idle"`       // CPU 空闲时间百分比
	UsageIowait    float64 `gorm:"type:decimal(10,6);not null" json:"usage_iowait"`     // 等待 I/O 完成的时间百分比
	UsageIrq       float64 `gorm:"type:decimal(10,6);not null" json:"usage_irq"`        // 处理硬件中断的时间百分比
	UsageNice      float64 `gorm:"type:decimal(10,6);not null" json:"usage_nice"`       // 运行低优先级进程的时间百分比
	UsageSoftirq   float64 `gorm:"type:decimal(10,6);not null" json:"usage_softirq"`    // 处理软件中断的时间百分比
	UsageSteal     float64 `gorm:"type:decimal(10,6);not null" json:"usage_steal"`      // 虚拟化环境中被其他虚拟机占用的时间百分比
	UsageSystem    float64 `gorm:"type:decimal(10,6);not null" json:"usage_system"`     // 内核态时间百分比
	UsageUser      float64 `gorm:"type:decimal(10,6);not null" json:"usage_user"`       // 用户态时间百分比

	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"` // 记录创建时间（Unix 时间戳）
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"` // 记录更新时间（Unix 时间戳）
}

// TableName 指定表名
func (CPUFieldsDb) TableName() string {
	return "cpu_metrics"
}

// FromCPUFields 从 CPUFields 和 tags 填充 CPUFieldsDb
func (c *CPUFieldsDb) FromCPUFields(cpu string, host string, timestamp int64, fields CPUFields) {
	c.CPU = cpu
	c.Host = host
	c.Timestamp = timestamp
	c.UsageActive = fields.UsageActive
	c.UsageGuest = fields.UsageGuest
	c.UsageGuestNice = fields.UsageGuestNice
	c.UsageIdle = fields.UsageIdle
	c.UsageIowait = fields.UsageIowait
	c.UsageIrq = fields.UsageIrq
	c.UsageNice = fields.UsageNice
	c.UsageSoftirq = fields.UsageSoftirq
	c.UsageSteal = fields.UsageSteal
	c.UsageSystem = fields.UsageSystem
	c.UsageUser = fields.UsageUser
}

// FromFieldsMap 填充 CPUFields
func (c *CPUFields) FromFieldsMap(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, c)
}

// SaveCPUToDB 保存 CPU 数据到数据库
func SaveCPUToDB(metric *TelegrafJson) {
	var cpuFields CPUFields
	if err := cpuFields.FromFieldsMap(metric.Fields); err != nil {
		log.Printf("解析 CPU 字段出错: %v", err)
		return
	}
	// 转换为数据库实体
	var cpuDb CPUFieldsDb
	cpuDb.FromCPUFields(
		metric.Tags["cpu"],  // CPU 标识
		metric.Tags["host"], // 主机名
		metric.Timestamp,    // 时间戳
		cpuFields,           // CPU 指标
	)
	// 保存到数据库
	if err := db.Create(&cpuDb).Error; err != nil {
		log.Printf("保存 CPU 数据到数据库出错: %v", err)
		return
	}

}
