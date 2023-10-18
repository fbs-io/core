/*
 * @Author: reel
 * @Date: 2023-10-15 07:48:02
 * @LastEditors: reel
 * @LastEditTime: 2023-10-19 06:35:41
 * @Description: 回掉函数
 */
package rdb

import (
	"fmt"

	"github.com/fbs-io/core/pkg/consts"
	"gorm.io/gorm"
)

// TODO: 增加链路追踪
func (store *rdbStore) registerCallbacks() {
	store.db.Callback().Query().Before("*").Register("subQuery", rdb.subQuery)
	store.db.Callback().Row().Before("*").Register("subQuery", rdb.subQuery)
	store.db.Callback().Raw().Before("*").Register("subQuery", rdb.subQuery)
	store.db.Callback().Create().Before("*").Register("creates", rdb.setCreatesCallback)
	store.db.Callback().Update().Before("*").Register("updates", rdb.setUpdatesCallback)
}

// 目前只要用于表内分区和多表分区模式
//
// 支持单表分区和多表分区模式
//
// TODO:增加多数据库模式完善中
func (store *rdbStore) switchSharding(tx *gorm.DB) {
	// 过滤初始化操作时的数据库操作
	if tx.Statement == nil {
		return
	}
	if tx.Statement.Schema == nil {
		return
	}

	table := tx.Statement.Table
	if !store.shardingAllTable[table] {
		return
	}

	sk, ok := tx.Get(consts.CTX_SHARDING_KEY)
	if !ok || sk.(string) == "" {
		return
	}

	// TODO:细化条件查询
	tx.Where("sk = ? ", sk)
	if tx.Statement.BuildClauses != nil {
		switch tx.Statement.BuildClauses[0] {
		// case "SELECT":
		// 	tx.Where("sk = ? ", sk)
		case "UPDATE", "CREATE":
			tx.Statement.SetColumn("sk", sk, true)
		}
	}
	// 统一设置条件和设置字段
	// 设置字段值, 用于更新, 创建用
	// tx.Statement.SetColumn("sk", sk, true)
	// 增加查询设置查询条件, 用于更新, 删除, 查询用

	switch store.shardingModel {
	// case SHADING_MODEL_ONE:
	// 	// 暂无可设置
	// 	// tx.Where("sk = ? ", sk)
	case SHADING_MODEL_TABLE:
		// 有分区字段, 但是么有设置分区表
		if store.shardingTable[table] != nil {
			tx.Statement.Table = fmt.Sprintf("%s_%s", table, sk)
			tx.Table(tx.Statement.Table)

		}
	case SHADING_MODEL_DB:
		// TODO 完善DB逻辑
	default:
	}

}

// 子查询, 用于分页,由子查询的字段和构建条件时,可以设置子查询
// 可以自定义子查询主键字段,
func (store *rdbStore) subQuery(tx *gorm.DB) {
	sk, _ := tx.Get(consts.CTX_SHARDING_KEY)

	cbI, ok := tx.Get(TX_CONDITION_BUILD_KEY)
	if !ok || cbI == nil {
		return
	}
	columnI, ok := tx.Get(TX_SUB_QUERY_COLUMN_KEY)
	if columnI == nil || !ok {
		return
	}
	table := tx.Statement.Table
	col := columnI.(string)
	cb := cbI.(*Condition)
	sub := rdb.BuildQuery(cb).Table(table)

	subColumns := fmt.Sprintf("sub%s", col)
	sub = sub.Select(fmt.Sprintf("%s as %s ", col, subColumns))

	// 用于分区
	if tx.Statement.Schema.FieldsByName["ShadingKey"] != nil && store.shardingTable[table] == nil {
		sub = sub.Where("sk = ? ", sk)
	}
	tx.Table(table).Joins(fmt.Sprintf("join ( ? ) t1 on t1.%s = %s", subColumns, col), sub)
}

// 设置创建前的回掉函数
func (store *rdbStore) setCreatesCallback(tx *gorm.DB) {
	field := tx.Statement.Schema.FieldsByDBName["created_by"]
	if field != nil && field.OwnerSchema.String() == "github.com/fbs-io/core/store/rdb.Model" {
		auth, _ := tx.Get(consts.CTX_AUTH)
		tx.Statement.SetColumn("created_by", auth, true)

	}

}

// 设置更新前的回掉函数
func (store *rdbStore) setUpdatesCallback(tx *gorm.DB) {
	field := tx.Statement.Schema.FieldsByDBName["updated_by"]
	if field != nil && field.OwnerSchema.String() == "github.com/fbs-io/core/store/rdb.Model" {
		auth, _ := tx.Get(consts.CTX_AUTH)
		if auth == nil {
			return
		}
		tx.Statement.SetColumn("updated_by", auth, true)
	}

}

// 设置软删除的操作人
// func (store *rdbStore) setDeleteCallback(tx *gorm.DB) {
// 	field := tx.Statement.Schema.FieldsByDBName["deleted_by"]
// 	if field != nil && field.OwnerSchema.String() == "github.com/fbs-io/core/store/rdb.Model" {
// 		auth, _ := tx.Get(consts.CTX_AUTH)
// 		defer tx.Set(consts.CTX_AUTH, auth)
// 		if auth == nil {
// 			return
// 		}
// 		txN := store.db.Where("1=1")

// 		// 复制model
// 		txN.Statement.Model = tx.Statement.Model

// 		// 复制条件
// 		txN.Statement.Clauses = tx.Statement.Clauses
// 		txN.Table(tx.Statement.Table).Update("deleted_by", auth)
// 	}
// }
