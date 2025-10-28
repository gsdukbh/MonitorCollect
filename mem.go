package main

import "encoding/json"

// MemFields 表示内存使用情况统计
// 处理telegraf 采集的内存数据
// 对应样例中 name=mem 的字段
type MemFields struct {
	Active           int64   `json:"active"`            // 活跃内存（字节）
	Available        int64   `json:"available"`         // 可用内存（字节）
	AvailablePercent float64 `json:"available_percent"` // 可用内存百分比
	Buffered         int64   `json:"buffered"`          // 缓冲区内存（字节）
	Cached           int64   `json:"cached"`            // 缓存内存（字节）
	CommitLimit      int64   `json:"commit_limit"`      // 可分配的总内存（字节）
	CommittedAs      int64   `json:"committed_as"`      // 已分配的内存（字节）
	Dirty            int64   `json:"dirty"`             // 等待写回磁盘的内存（字节）
	Free             int64   `json:"free"`              // 空闲内存（字节）
	HighFree         int64   `json:"high_free"`         // 高位空闲内存（字节）
	HighTotal        int64   `json:"high_total"`        // 高位总内存（字节）
	HugePageSize     int64   `json:"huge_page_size"`    // 大页面大小（字节）
	HugePagesFree    int64   `json:"huge_pages_free"`   // 空闲大页面数量
	HugePagesTotal   int64   `json:"huge_pages_total"`  // 总大页面数量
	Inactive         int64   `json:"inactive"`          // 不活跃内存（字节）
	LowFree          int64   `json:"low_free"`          // 低位空闲内存（字节）
	LowTotal         int64   `json:"low_total"`         // 低位总内存（字节）
	Mapped           int64   `json:"mapped"`            // 映射内存（字节）
	PageTables       int64   `json:"page_tables"`       // 页表内存（字节）
	Shared           int64   `json:"shared"`            // 共享内存（字节）
	Slab             int64   `json:"slab"`              // Slab 内存（字节）
	Sreclaimable     int64   `json:"sreclaimable"`      // 可回收 Slab 内存（字节）
	Sunreclaim       int64   `json:"sunreclaim"`        // 不可回收 Slab 内存（字节）
	SwapCached       int64   `json:"swap_cached"`       // 交换缓存（字节）
	SwapFree         int64   `json:"swap_free"`         // 空闲交换空间（字节）
	SwapTotal        int64   `json:"swap_total"`        // 总交换空间（字节）
	Total            int64   `json:"total"`             // 总内存（字节）
	Used             int64   `json:"used"`              // 已用内存（字节）
	UsedPercent      float64 `json:"used_percent"`      // 内存使用百分比
	VmallocChunk     int64   `json:"vmalloc_chunk"`     // 最大 vmalloc 块（字节）
	VmallocTotal     int64   `json:"vmalloc_total"`     // 总 vmalloc 空间（字节）
	VmallocUsed      int64   `json:"vmalloc_used"`      // 已用 vmalloc 空间（字节）
	WriteBack        int64   `json:"write_back"`        // 正在写回的内存（字节）
	WriteBackTmp     int64   `json:"write_back_tmp"`    // 临时写回缓冲区（字节）
}

// MemFieldsDb 用于数据库存储的 Mem 字段结构体
// 使用 gorm 标签定义数据库字段映射
type MemFieldsDb struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`           // 主键ID，自增
	Host      string `gorm:"type:varchar(100);not null;index" json:"host"` // 主机名
	Timestamp int64  `gorm:"not null;index" json:"timestamp"`              // 时间戳（纳秒或毫秒）

	// 内存统计字段
	Active           int64   `gorm:"not null" json:"active"`                               // 活跃内存（字节）
	Available        int64   `gorm:"not null" json:"available"`                            // 可用内存（字节）
	AvailablePercent float64 `gorm:"type:decimal(10,6);not null" json:"available_percent"` // 可用内存百分比
	Buffered         int64   `gorm:"not null" json:"buffered"`                             // 缓冲区内存（字节）
	Cached           int64   `gorm:"not null" json:"cached"`                               // 缓存内存（字节）
	CommitLimit      int64   `gorm:"not null" json:"commit_limit"`                         // 可分配的总内存（字节）
	CommittedAs      int64   `gorm:"not null" json:"committed_as"`                         // 已分配的内存（字节）
	Dirty            int64   `gorm:"not null" json:"dirty"`                                // 等待写回磁盘的内存（字节）
	Free             int64   `gorm:"not null" json:"free"`                                 // 空闲内存（字节）
	HighFree         int64   `gorm:"not null" json:"high_free"`                            // 高位空闲内存（字节）
	HighTotal        int64   `gorm:"not null" json:"high_total"`                           // 高位总内存（字节）
	HugePageSize     int64   `gorm:"not null" json:"huge_page_size"`                       // 大页面大小（字节）
	HugePagesFree    int64   `gorm:"not null" json:"huge_pages_free"`                      // 空闲大页面数量
	HugePagesTotal   int64   `gorm:"not null" json:"huge_pages_total"`                     // 总大页面数量
	Inactive         int64   `gorm:"not null" json:"inactive"`                             // 不活跃内存（字节）
	LowFree          int64   `gorm:"not null" json:"low_free"`                             // 低位空闲内存（字节）
	LowTotal         int64   `gorm:"not null" json:"low_total"`                            // 低位总内存（字节）
	Mapped           int64   `gorm:"not null" json:"mapped"`                               // 映射内存（字节）
	PageTables       int64   `gorm:"not null" json:"page_tables"`                          // 页表内存（字节）
	Shared           int64   `gorm:"not null" json:"shared"`                               // 共享内存（字节）
	Slab             int64   `gorm:"not null" json:"slab"`                                 // Slab 内存（字节）
	Sreclaimable     int64   `gorm:"not null" json:"sreclaimable"`                         // 可回收 Slab 内存（字节）
	Sunreclaim       int64   `gorm:"not null" json:"sunreclaim"`                           // 不可回收 Slab 内存（字节）
	SwapCached       int64   `gorm:"not null" json:"swap_cached"`                          // 交换缓存（字节）
	SwapFree         int64   `gorm:"not null" json:"swap_free"`                            // 空闲交换空间（字节）
	SwapTotal        int64   `gorm:"not null" json:"swap_total"`                           // 总交换空间（字节）
	Total            int64   `gorm:"not null" json:"total"`                                // 总内存（字节）
	Used             int64   `gorm:"not null" json:"used"`                                 // 已用内存（字节）
	UsedPercent      float64 `gorm:"type:decimal(10,6);not null" json:"used_percent"`      // 内存使用百分比
	VmallocChunk     int64   `gorm:"not null" json:"vmalloc_chunk"`                        // 最大 vmalloc 块（字节）
	VmallocTotal     int64   `gorm:"not null" json:"vmalloc_total"`                        // 总 vmalloc 空间（字节）
	VmallocUsed      int64   `gorm:"not null" json:"vmalloc_used"`                         // 已用 vmalloc 空间（字节）
	WriteBack        int64   `gorm:"not null" json:"write_back"`                           // 正在写回的内存（字节）
	WriteBackTmp     int64   `gorm:"not null" json:"write_back_tmp"`                       // 临时写回缓冲区（字节）

	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"` // 记录创建时间（Unix 时间戳）
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"` // 记录更新时间（Unix 时间戳）
}

// TableName 指定表名
func (MemFieldsDb) TableName() string {
	return "mem_metrics"
}

// FromMemFields 从 MemFields 和 tags 填充 MemFieldsDb
func (m *MemFieldsDb) FromMemFields(host string, timestamp int64, fields MemFields) {
	m.Host = host
	m.Timestamp = timestamp
	m.Active = fields.Active
	m.Available = fields.Available
	m.AvailablePercent = fields.AvailablePercent
	m.Buffered = fields.Buffered
	m.Cached = fields.Cached
	m.CommitLimit = fields.CommitLimit
	m.CommittedAs = fields.CommittedAs
	m.Dirty = fields.Dirty
	m.Free = fields.Free
	m.HighFree = fields.HighFree
	m.HighTotal = fields.HighTotal
	m.HugePageSize = fields.HugePageSize
	m.HugePagesFree = fields.HugePagesFree
	m.HugePagesTotal = fields.HugePagesTotal
	m.Inactive = fields.Inactive
	m.LowFree = fields.LowFree
	m.LowTotal = fields.LowTotal
	m.Mapped = fields.Mapped
	m.PageTables = fields.PageTables
	m.Shared = fields.Shared
	m.Slab = fields.Slab
	m.Sreclaimable = fields.Sreclaimable
	m.Sunreclaim = fields.Sunreclaim
	m.SwapCached = fields.SwapCached
	m.SwapFree = fields.SwapFree
	m.SwapTotal = fields.SwapTotal
	m.Total = fields.Total
	m.Used = fields.Used
	m.UsedPercent = fields.UsedPercent
	m.VmallocChunk = fields.VmallocChunk
	m.VmallocTotal = fields.VmallocTotal
	m.VmallocUsed = fields.VmallocUsed
	m.WriteBack = fields.WriteBack
	m.WriteBackTmp = fields.WriteBackTmp
}

// FromFieldsMap 填充 MemFields
func (m *MemFields) FromFieldsMap(fieldsMap map[string]interface{}) error {
	b, err := json.Marshal(fieldsMap)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, m)
}
