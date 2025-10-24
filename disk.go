package main

import "encoding/json"

// DiskFields 表示磁盘使用情况统计
// 处理telegraf 采集的磁盘数据
// 对应样例中 name=disk 的字段
type DiskFields struct {
	Free              int64   `json:"free"`                // 可用空间（字节）
	InodesFree        int64   `json:"inodes_free"`         // 可用 inode 数量
	InodesTotal       int64   `json:"inodes_total"`        // 总 inode 数量
	InodesUsed        int64   `json:"inodes_used"`         // 已用 inode 数量
	InodesUsedPercent float64 `json:"inodes_used_percent"` // inode 使用百分比
	Total             int64   `json:"total"`               // 总空间（字节）
	Used              int64   `json:"used"`                // 已用空间（字节）
	UsedPercent       float64 `json:"used_percent"`        // 空间使用百分比
}

// DiskFieldsDb 用于数据库存储的 Disk 字段结构体
// 使用 gorm 标签定义数据库字段映射
type DiskFieldsDb struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`             // 主键ID，自增
	Device    string `gorm:"type:varchar(100);not null;index" json:"device"` // 设备名称，如 "mmcblk0p2"
	Fstype    string `gorm:"type:varchar(50);not null" json:"fstype"`        // 文件系统类型，如 "ext4"
	Host      string `gorm:"type:varchar(100);not null;index" json:"host"`   // 主机名
	Mode      string `gorm:"type:varchar(20);not null" json:"mode"`          // 挂载模式，如 "rw"
	Path      string `gorm:"type:varchar(255);not null;index" json:"path"`   // 挂载路径，如 "/"
	Timestamp int64  `gorm:"not null;index" json:"timestamp"`                // 时间戳（纳秒或毫秒）

	// 磁盘空间统计字段
	Free        int64   `gorm:"not null" json:"free"`                            // 可用空间（字节）
	Total       int64   `gorm:"not null" json:"total"`                           // 总空间（字节）
	Used        int64   `gorm:"not null" json:"used"`                            // 已用空间（字节）
	UsedPercent float64 `gorm:"type:decimal(10,6);not null" json:"used_percent"` // 空间使用百分比

	// Inode 统计字段
	InodesFree        int64   `gorm:"not null" json:"inodes_free"`                            // 可用 inode 数量
	InodesTotal       int64   `gorm:"not null" json:"inodes_total"`                           // 总 inode 数量
	InodesUsed        int64   `gorm:"not null" json:"inodes_used"`                            // 已用 inode 数量
	InodesUsedPercent float64 `gorm:"type:decimal(10,6);not null" json:"inodes_used_percent"` // inode 使用百分比

	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"` // 记录创建时间（Unix 时间戳）
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"` // 记录更新时间（Unix 时间戳）
}

// TableName 指定表名
func (DiskFieldsDb) TableName() string {
	return "disk_metrics"
}

// FromDiskFields 从 DiskFields 和 tags 填充 DiskFieldsDb
func (d *DiskFieldsDb) FromDiskFields(device, fstype, host, mode, path string, timestamp int64, fields DiskFields) {
	d.Device = device
	d.Fstype = fstype
	d.Host = host
	d.Mode = mode
	d.Path = path
	d.Timestamp = timestamp
	d.Free = fields.Free
	d.Total = fields.Total
	d.Used = fields.Used
	d.UsedPercent = fields.UsedPercent
	d.InodesFree = fields.InodesFree
	d.InodesTotal = fields.InodesTotal
	d.InodesUsed = fields.InodesUsed
	d.InodesUsedPercent = fields.InodesUsedPercent
}

// FromFieldsMap 填充 DiskFields
func (d *DiskFields) FromFieldsMap(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}
