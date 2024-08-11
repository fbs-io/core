/*
 * @Author: reel
 * @Date: 2023-10-15 22:49:03
 * @LastEditors: reel
 * @LastEditTime: 2024-08-11 11:48:44
 * @Description: 分区相关
 */
package rdb

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fbs-io/core/pkg/env"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/dsn"
	"gorm.io/gorm"
)

// 设置分区模式
//
// SHADING_MODEL_NOT : 不分区, 默认值, 数据都在一张表中,
//
// SHADING_MODEL_TABLE: 按表分区, 根据分区字段值, 设置表后缀
//
// TODO:SHADING_MODEL_DB: 按库(不按schema)分区, 根据分区字段值, 设置不同的库名后缀
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
	store.setDBCallbackWithSharding(store.db)

}

func (store *rdbStore) setDBCallbackWithSharding(db *gorm.DB) {
	db.Callback().Query().Before("*").Register("sharding", store.switchSharding)
	db.Callback().Row().Before("*").Register("sharding", store.switchSharding)
	db.Callback().Raw().Before("*").Register("sharding", store.switchSharding)
	db.Callback().Create().Before("*").Register("sharding", store.switchSharding)
	db.Callback().Update().Before("*").Register("sharding", store.switchSharding)
	db.Callback().Delete().Before("*").Register("sharding", store.switchSharding)
}

func (store *rdbStore) ShardingModel() (model int8) {
	return store.shardingModel
}
func (store *rdbStore) ShardingTable() map[string][]string {
	return store.shardingTable
}

// 自定义分区表, 比如按时间等分区
// TODO:逻辑待完善
func (store *rdbStore) AddShardingTable(tableName string) {
	store.shardingTable[tableName] = make([]string, 0, 100)
	// for _, suffix := range store.shardingSuffixs {
	// 	store.shardingTable[tableName] = append(store.shardingTable[tableName], fmt.Sprintf("%s_%s", tableName, suffix))
	// }
}

// 项目启动时, 添加初始化执行的动作, 如迁移表等
func (store *rdbStore) AddMigrateList(fs ...func() error) {
	store.migrateList = append(store.migrateList, fs...)
}

// 设置分区后缀
//
// 同时自动创建分区表
func (store *rdbStore) AddShardingSuffixs(suffixs string) (err error) {
	store.shardingSuffixs = append(store.shardingSuffixs, suffixs)
	for tableName := range store.shardingAllTable {
		err = store.AutoShardingTable(tableName, suffixs)
		if err != nil {
			return
		}
	}
	return
}

// 迁移表
//
// 对表分区和库分区分别处理
//
// 支持自定义用于分区迁移的表
func (store *rdbStore) AutoShardingTable(tableName, suffixs string) (err error) {
	tabler := store.tablers[tableName]
	if tabler == nil {
		return errorx.Errorf("无法获取表名为:%s的表结构:", tableName)
	}

	// 通过反射获取模型中是否包含分区字段用于创建分区
	rt := reflect.TypeOf(tabler).Elem()
	rtModel, ok1 := rt.FieldByName("ShardingModel")
	rtKey, ok2 := rt.FieldByName("ShadingKey")
	// 通过多重判断, 确定模型中包含了分区字段
	if ok1 && ok2 &&
		rtKey.Name == "ShadingKey" &&
		rtModel.Name == "ShardingModel" &&
		strings.Contains(rtKey.Tag.Get("gorm"), "column:sk") {

		store.shardingAllTable[tabler.TableName()] = true
	}
	// 增加数据权限字段的判断
	rtModel, ok1 = rt.FieldByName("DataPermissionStringModel")
	rtKey, ok2 = rt.FieldByName("DataPermission")
	if ok1 && ok2 &&
		rtKey.Name == "DataPermission" &&
		rtModel.Name == "DataPermissionStringModel" &&
		strings.Contains(rtKey.Tag.Get("gorm"), "column:dp") {

		store.dataPermissionTable[tabler.TableName()] = true
		store.dataPermissionTable["DataPermissionStringModel"] = true
	}
	// 增加数据权限字段的判断
	rtModel, ok1 = rt.FieldByName("DataPermissionIntModel")
	rtKey, ok2 = rt.FieldByName("DataPermission")
	if ok1 && ok2 &&
		rtKey.Name == "DataPermission" &&
		rtModel.Name == "DataPermissionIntModel" &&
		strings.Contains(rtKey.Tag.Get("gorm"), "column:dp") {

		store.dataPermissionTable[tabler.TableName()] = true
		store.dataPermissionTable["DataPermissionIntModel"] = true
	}

	// 处理表迁移
	switch store.shardingModel {
	// 处理表分区时的表迁移
	case SHADING_MODEL_TABLE:
		for _, suffix := range store.shardingSuffixs {
			nweTable := fmt.Sprintf("%s_%v", tableName, suffix)
			if suffixs != "" && suffixs != suffix.(string) {
				continue
			}
			if strings.Contains(env.Active().DBInit(), tableName) ||
				env.Active().DBInit() == TABLE_INIT_ALL {
				// 分区表在重置主表时也全部重置
				err = store.db.Table(nweTable).Migrator().DropTable(tabler)
				if err != nil {
					return
				}
			}
			err = store.db.Table(nweTable).Migrator().AutoMigrate(tabler)
			if err != nil {
				return
			}
		}

	// 处理库分区时的表迁移
	case SHADING_MODEL_DB:
		for _, suffix := range store.shardingSuffixs {
			if suffixs != "" && suffixs != suffix.(string) {
				continue
			}
			db := store.dbPool[suffix.(string)]
			if db == nil {
				suffixDsn := dsn.CopyDsn(store.dsn)
				suffixDsn.Name = fmt.Sprintf("%s_%s", suffixDsn.Name, strings.ToLower(suffix.(string)))
				if suffixDsn.Type == dsn.DSN_TYPE_PGSQL {
					datname := ""
					store.db.Table("pg_database").Where(" datname = ?", suffixDsn.Name).Select("datname").Find(&datname)
					if datname != suffixDsn.Name {
						store.db.Exec(fmt.Sprintf(`%s %s ;`, "create database ", suffixDsn.Name))
					}
				} else {
					store.db.Exec(fmt.Sprintf(`%s %s ;`, "create database if not exists", suffixDsn.Name))
				}
				db, err = genDBWithDsn(suffixDsn)
				if err != nil {
					return err
				}
				store.setDBCallbackWithSharding(db)
				store.dbPool[suffix.(string)] = db
			}
			if strings.Contains(env.Active().DBInit(), tableName) ||
				env.Active().DBInit() == TABLE_INIT_ALL {
				// 分区表在重置主表时也全部重置
				err = db.Table(tableName).Migrator().DropTable(tabler)
				if err != nil {
					return
				}
			}
			err = db.Table(tableName).Migrator().AutoMigrate(tabler)
			if err != nil {
				return
			}
		}
	}
	return nil
}

// 设置分区后缀
//
// 同时自动创建分区表
func (store *rdbStore) AddShardingSuffixsWithTX(tx *gorm.DB, suffixs string) (err error) {
	store.shardingSuffixs = append(store.shardingSuffixs, suffixs)
	for tableName := range store.shardingAllTable {
		err = store.AutoShardingTableWithTX(tx, tableName, suffixs)
		if err != nil {
			return
		}
	}
	return
}

// 迁移表
//
// 对表分区和库分区分别处理
//
// 支持自定义用于分区迁移的表
func (store *rdbStore) AutoShardingTableWithTX(tx *gorm.DB, tableName, suffixs string) (err error) {
	tabler := store.tablers[tableName]
	if tabler == nil {
		return errorx.Errorf("无法获取表名为:%s的表结构:", tableName)
	}

	// 通过反射获取模型中是否包含分区字段用于创建分区
	rt := reflect.TypeOf(tabler).Elem()
	rtModel, ok1 := rt.FieldByName("ShardingModel")
	rtKey, ok2 := rt.FieldByName("ShadingKey")
	// 通过多重判断, 确定模型中包含了分区字段
	if ok1 && ok2 &&
		rtKey.Name == "ShadingKey" &&
		rtModel.Name == "ShardingModel" &&
		strings.Contains(rtKey.Tag.Get("gorm"), "column:sk") {

		store.shardingAllTable[tabler.TableName()] = true
	}
	// 增加数据权限字段的判断
	rtModel, ok1 = rt.FieldByName("DataPermissionStringModel")
	rtKey, ok2 = rt.FieldByName("DataPermission")
	if ok1 && ok2 &&
		rtKey.Name == "DataPermission" &&
		rtModel.Name == "DataPermissionStringModel" &&
		strings.Contains(rtKey.Tag.Get("gorm"), "column:dp") {

		store.dataPermissionTable[tabler.TableName()] = true
		store.dataPermissionTable["DataPermissionStringModel"] = true
	}
	// 增加数据权限字段的判断
	rtModel, ok1 = rt.FieldByName("DataPermissionIntModel")
	rtKey, ok2 = rt.FieldByName("DataPermission")
	if ok1 && ok2 &&
		rtKey.Name == "DataPermission" &&
		rtModel.Name == "DataPermissionIntModel" &&
		strings.Contains(rtKey.Tag.Get("gorm"), "column:dp") {

		store.dataPermissionTable[tabler.TableName()] = true
		store.dataPermissionTable["DataPermissionIntModel"] = true
	}

	// 处理表迁移
	switch store.shardingModel {
	// 处理表分区时的表迁移
	case SHADING_MODEL_TABLE:
		for _, suffix := range store.shardingSuffixs {
			nweTable := fmt.Sprintf("%s_%v", tableName, suffix)
			if suffixs != "" && suffixs != suffix.(string) {
				continue
			}
			if strings.Contains(env.Active().DBInit(), tableName) ||
				env.Active().DBInit() == TABLE_INIT_ALL {
				// 分区表在重置主表时也全部重置
				err = store.db.Table(nweTable).Migrator().DropTable(tabler)
				if err != nil {
					return
				}
			}
			err = tx.Table(nweTable).Migrator().AutoMigrate(tabler)
			if err != nil {
				return
			}
		}

	// 处理库分区时的表迁移
	case SHADING_MODEL_DB:
		for _, suffix := range store.shardingSuffixs {
			if suffixs != "" && suffixs != suffix.(string) {
				continue
			}
			db := store.dbPool[suffix.(string)]
			if db == nil {
				suffixDsn := dsn.CopyDsn(store.dsn)
				suffixDsn.Name = fmt.Sprintf("%s_%s", suffixDsn.Name, strings.ToLower(suffix.(string)))
				if suffixDsn.Type == dsn.DSN_TYPE_PGSQL {
					datname := ""
					store.db.Table("pg_database").Where(" datname = ?", suffixDsn.Name).Select("datname").Find(&datname)
					if datname != suffixDsn.Name {
						store.db.Exec(fmt.Sprintf(`%s %s ;`, "create database ", suffixDsn.Name))
					}
				} else {
					store.db.Exec(fmt.Sprintf(`%s %s ;`, "create database if not exists", suffixDsn.Name))
				}
				db, err = genDBWithDsn(suffixDsn)
				if err != nil {
					return err
				}
				store.dbPool[suffix.(string)] = db
			}

			tx.Statement.ConnPool = db.Config.ConnPool
			if strings.Contains(env.Active().DBInit(), tableName) ||
				env.Active().DBInit() == TABLE_INIT_ALL {
				// 分区表在重置主表时也全部重置
				err = tx.Table(tableName).Migrator().DropTable(tabler)
				if err != nil {
					return
				}
			}
			err = tx.Table(tableName).Migrator().AutoMigrate(tabler)
			if err != nil {
				return
			}
		}
	}
	return nil
}
