/*
 * @Author: reel
 * @Date: 2023-10-15 22:49:03
 * @LastEditors: reel
 * @LastEditTime: 2023-10-17 07:41:02
 * @Description: 分区相关
 */
package rdb

import (
	"fmt"
	"strings"

	"github.com/fbs-io/core/pkg/env"
	"github.com/fbs-io/core/pkg/errorx"
)

// 设置分区模式
//
// SHADING_MODEL_NOT : 不分区, 默认值, 数据都在一张表中,
//
// SHADING_MODEL_TABLE: 按表分区, 根据分区字段值, 设置表后缀
//
// TODO:SHADING_MODEL_DB: 按库(schema)分区, 根据分区字段值, 设置不同的库名(schema)后缀
//
// 该模式适用于使用cores 上下文ctx.TX()方式生成的 gorm.DB, 且在上下文中传入了分区字段, 会自动构建查询条件, 配合 ShardingModel使用,可以自动写入分区字段
//
// 如果直接使用gorm.DB查询, 该设置并不会生效
func (store *rdbStore) SetShardingModel(model int8, suffix []interface{}) {
	store.shardingModel = model
	store.shardingSuffixs = suffix
	if suffix == nil {
		store.shardingSuffixs = make([]interface{}, 0, 100)
	}

	store.db.Callback().Query().Before("*").Register("sharding", rdb.switchSharding)
	store.db.Callback().Row().Before("*").Register("sharding", rdb.switchSharding)
	store.db.Callback().Raw().Before("*").Register("sharding", rdb.switchSharding)
	store.db.Callback().Create().Before("*").Register("sharding", rdb.switchSharding)
	store.db.Callback().Update().Before("*").Register("sharding", rdb.switchSharding)
	store.db.Callback().Delete().Before("*").Register("sharding", rdb.switchSharding)
}

func (store *rdbStore) ShardingModel() (model int8) {
	return store.shardingModel
}
func (store *rdbStore) ShardingTable() map[string][]string {
	return store.shardingTable
}

// SHADING_MODEL_TABLE(按表分区) 模式下, 自定义添加可用于分区的表
func (store *rdbStore) AddShardingTable(tableName string) {
	store.shardingTable[tableName] = make([]string, 0, 100)
	for _, suffix := range store.shardingSuffixs {
		store.shardingTable[tableName] = append(store.shardingTable[tableName], fmt.Sprintf("%s_%s", tableName, suffix))
	}
}

// SHADING_MODEL_TABLE(按表分区) 模式下, 自定义添加可用于分区的表
func (store *rdbStore) AddMigrateList(fs ...func() error) {
	store.migrateList = append(store.migrateList, fs...)
}

// 设置分区后缀
//
// 同时自动创建分区表
func (store *rdbStore) AddShardingSuffixs(suffixs string) (err error) {
	store.shardingSuffixs = append(store.shardingSuffixs, suffixs)
	for tableName, _ := range store.shardingTable {
		store.AddShardingTable(tableName)
		err = store.AutoShardingTable(tableName)
		if err != nil {
			return
		}
	}
	return
}

// 迁移表
func (store *rdbStore) AutoShardingTable(tableName string) (err error) {
	tabler := store.tablers[tableName]
	if tabler == nil {
		return errorx.Errorf("无法获取表名为:%s的表结构:", tableName)
	}
	if strings.Contains(env.Active().DBInit(), tableName) ||
		env.Active().DBInit() == TABLE_INIT_ALL {
		// 分区表在重置主表时也全部重置
		if store.shardingModel == SHADING_MODEL_TABLE {
			for _, table := range store.shardingTable[tableName] {
				err = store.db.Table(table).Migrator().DropTable(tabler)
				if err != nil {
					return
				}
			}
		}

		if store.shardingModel == SHADING_MODEL_TABLE {
			for _, table := range store.shardingTable[tableName] {
				err = store.db.Table(table).Migrator().AutoMigrate(tabler)
				if err != nil {
					return
				}
			}
		}

	} else {
		if store.shardingModel == SHADING_MODEL_TABLE {
			for _, table := range store.shardingTable[tableName] {
				err = store.db.Table(table).Migrator().AutoMigrate(tabler)
				if err != nil {
					return
				}
			}
		}
	}
	return nil
}
