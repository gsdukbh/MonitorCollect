package main

import (
	"encoding/json"
	"log"
	"time"
)

// 处理telegraf 采集的网络数据

// NetInterfaceFields 表示每个网卡(interface)的基础网络统计
// 对应样例中 name=net, tag.interface=end0 的字段
type NetInterfaceFields struct {
	BytesRecv   int64 `json:"bytes_recv"`   // 接收的总字节数
	BytesSent   int64 `json:"bytes_sent"`   // 发送的总字节数
	DropIn      int64 `json:"drop_in"`      // 接收时丢弃的数据包数
	DropOut     int64 `json:"drop_out"`     // 发送时丢弃的数据包数
	ErrIn       int64 `json:"err_in"`       // 接收时的错误数
	ErrOut      int64 `json:"err_out"`      // 发送时的错误数
	PacketsRecv int64 `json:"packets_recv"` // 接收的数据包总数
	PacketsSent int64 `json:"packets_sent"` // 发送的数据包总数
	Speed       int64 `json:"speed"`        // 网卡速度（Mbps）
}

// NetInterfaceFieldsDb 是用于存储网络接口统计数据的 GORM 模型
type NetInterfaceFieldsDb struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`         // 数据库主键
	Host      string `gorm:"type:varchar(100);not null;index"` // 主机名
	Interface string `gorm:"type:varchar(50);not null;index"`  // 网卡接口名
	Timestamp int64  `gorm:"not null;index"`                   // 数据采集时间戳

	BytesRecv   int64 `gorm:"column:bytes_recv"`   // 接收的总字节数
	BytesSent   int64 `gorm:"column:bytes_sent"`   // 发送的总字节数
	DropIn      int64 `gorm:"column:drop_in"`      // 接收时丢弃的数据包数
	DropOut     int64 `gorm:"column:drop_out"`     // 发送时丢弃的数据包数
	ErrIn       int64 `gorm:"column:err_in"`       // 接收时的错误数
	ErrOut      int64 `gorm:"column:err_out"`      // 发送时的错误数
	PacketsRecv int64 `gorm:"column:packets_recv"` // 接收的数据包总数
	PacketsSent int64 `gorm:"column:packets_sent"` // 发送的数据包总数
	Speed       int64 `gorm:"column:speed"`        // 网卡速度（Mbps）

	CreatedAt int64 `gorm:"autoCreateTime"` // 记录创建时间
	UpdatedAt int64 `gorm:"autoUpdateTime"` // 记录更新时间
}

// TableName 指定 NetInterfaceFieldsDb 的表名
func (NetInterfaceFieldsDb) TableName() string {
	return "net_interface_metrics"
}

// NetInterfaceCollectHour net Interface 按小时存储网络流量数据信息。
type NetInterfaceCollectHour struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`         // 数据库主键
	Host      string    `gorm:"type:varchar(100);not null;index"` // 主机名
	Interface string    `gorm:"type:varchar(50);not null;index"`  // 网卡接口名
	Hour      time.Time // 小时时间戳（格式：YYYYMMDDHH）
	Total     int64     // 小时内总流量（MB）
	Speed     float64   `gorm:"type:decimal(10,2)"` // 小时内平均速度（Mbps）保留两位小数
	SpeedStr  string    `gorm:"type:varchar(20)"`   // 格式化后的平均速度（e.g., "1.5 Mbps", "500 Kbps"）
	CreatedAt time.Time // 记录创建时间
}

// TableName 指定 NetInterfaceCollectHour 的表名
func (NetInterfaceCollectHour) TableName() string {
	return "net_interface_collect_hours"
}

// FromNetInterfaceFields 从 NetInterfaceFields 和标签填充 NetInterfaceFieldsDb
func (db *NetInterfaceFieldsDb) FromNetInterfaceFields(host, iface string, timestamp int64, fields NetInterfaceFields) {
	db.Host = host
	db.Interface = iface
	db.Timestamp = timestamp
	db.BytesRecv = fields.BytesRecv
	db.BytesSent = fields.BytesSent
	db.DropIn = fields.DropIn
	db.DropOut = fields.DropOut
	db.ErrIn = fields.ErrIn
	db.ErrOut = fields.ErrOut
	db.PacketsRecv = fields.PacketsRecv
	db.PacketsSent = fields.PacketsSent
	db.Speed = fields.Speed
}

// FromFieldsMap 填充 NetInterfaceFields
func (n *NetInterfaceFields) FromFieldsMap(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, n)
}

// SaveNetToDB 根据 tags 区分并保存网络数据到相应的数据库表
func SaveNetToDB(metric *TelegrafJson) {
	iface, isInterfaceMetric := metric.Tags["interface"]

	if isInterfaceMetric && iface != "all" {
		// 处理单个网卡的数据
		var netFields NetInterfaceFields
		if err := netFields.FromFieldsMap(metric.Fields); err != nil {
			log.Printf("解析网络接口字段出错: %v", err)
			return
		}
		var netDb NetInterfaceFieldsDb
		netDb.FromNetInterfaceFields(
			metric.Tags["host"],
			iface,
			metric.Timestamp,
			netFields,
		)
		// 此处应调用 gorm.DB.Create(&netDb) 来保存数据
		log.Printf("准备保存网络接口数据: %+v", netDb)
		if err := db.Create(&netDb).Error; err != nil {
			log.Printf("保存内存数据到数据库出错: %v", err)
			return
		}

	}
}
