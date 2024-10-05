/*
 * @Author: reel
 * @Date: 2024-08-14 07:24:57
 * @LastEditors: reel
 * @LastEditTime: 2024-10-05 00:01:14
 * @Description: 管理表结构
 */
package rdb

import "time"

const (
	TABLE_SYSTEM_CORE_ENTITYINFO = "e_sys_core_entityinfo"
	TABLE_SYSTEM_CORE_SHARDING   = "e_sys_core_sharding"
)

type EntityInfo struct {
	TableName          string        // 表名
	TableComment       string        // 表描述
	IsSharding         bool          // 是否是分区
	ShardingModel      int8          // 分区模式, 表分区, 库分区等
	IsDataPermission   bool          // 是否有数据权限
	DataPermissionType string        // 数据权限字段是数字类型还是字符串类型
	IsPartition        bool          // 是否是分区表, 同一张表中, 按业务字段进行分区, 一般时业务日期
	PartitionModel     string        // 分区模式, 日期, id等, 暂时只支持日期
	CacheTTL           time.Duration // 数据缓存时间, -1 表示永久缓存, 单位秒
}

func (t *EntityInfo) TableaName() string {
	return TABLE_SYSTEM_CORE_ENTITYINFO
}

// 分区后缀
type Sharding struct {
	Suffix string `gorm:"unique"` // 分区后缀

}

func (t *Sharding) TableName() string {
	return TABLE_SYSTEM_CORE_SHARDING
}
