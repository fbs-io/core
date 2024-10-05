/*
 * @Author: reel
 * @Date: 2023-10-15 07:48:02
 * @LastEditors: reel
 * @LastEditTime: 2024-10-05 00:20:28
 * @Description: 回掉函数
 */
package rdb

import (
	"fmt"
	"strings"

	"github.com/fbs-io/core/pkg/consts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: 增加链路追踪
func (store *rdbStore) registerCallbacks() {
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
	var (
		sks   string
		skDBs string
	)
	sub := store.buildSubQuery(tx)
	store.dataPermissonCallback(tx, sub)

	// 分区表不存在不再直接查询,TODO: 增加锁
	table := tx.Statement.Table
	if !store.shardingAllTable[table] {
		return
	}
	// 分区字段不存在,不再处理
	sk, ok := tx.Get(consts.CTX_SHARDING_KEY)
	if !ok || sk.(string) == "" {
		return
	}
	sks = sk.(string)

	// 分区数据库标识
	skDB, ok := tx.Get(consts.CTX_SHARDING_DB)
	if ok {
		skDBs = skDB.(string)
	}

	if tx.Statement.BuildClauses != nil {
		switch tx.Statement.BuildClauses[0] {
		case "SELECT":
			if sub != nil {
				sub.Where("sk = ?", sk)
			} else {
				tx.Where("sk = ?", sk)
			}
		case "UPDATE":
			store.setUpdatesCallback(tx)
			tx.Statement.SetColumn("sk", sk, true)
			tx.Where("sk = ? ", sk)
		case "INSERT":
			store.setCreatesCallback(tx)
			tx.Statement.SetColumn("sk", sk, true)

		}

		switch store.shardingModel {

		case SHADING_MODEL_TABLE:
			// 有分区字段, 但是么有设置分区表
			if store.shardingTable[table] != nil {
				tx.Statement.Table = fmt.Sprintf("%s_%s", table, sk)
				tx.Table(tx.Statement.Table)

			}
		case SHADING_MODEL_DB:
			db := store.dbPool[sk.(string)]
			if db != nil && sks != skDBs {
				tx.Statement.ConnPool = db.Config.ConnPool
			}
			if sub != nil {
				sub.Statement.ConnPool = tx.Statement.ConnPool
			}
		default:
		}
	}

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

const (
	DATA_PERMISSION_ONESELF       int8 = iota + 1 //本人可见
	DATA_PERMISSION_ALL                           //全部可见
	DATA_PERMISSION_ONLY_DEPT                     //所在部门可见
	DATA_PERMISSION_ONLY_DEPT_ALL                 //所在部门及子级可见
	DATA_PERMISSION_ONLY_CUSTOM                   //选择的部门可见
)

// 对有数据权限的表进行操作
func (store *rdbStore) dataPermissonCallback(tx *gorm.DB, subTx *gorm.DB) *gorm.DB {
	// 过滤初始化操作时的数据库操作
	if tx.Statement == nil {
		return tx
	}
	if tx.Statement.Schema == nil {
		return tx
	}
	table := tx.Statement.Table
	if !store.dataPermissionTable[table] {
		return tx
	}
	dpi, ok := tx.Get(consts.CTX_DATA_PERMISSION_KEY)
	if !ok && dpi != nil {
		return tx
	}

	if store.dataPermissionTable["DataPermissionIntModel"] {
		dp := dpi.(*DataPermissionIntCtx)
		tx = switchMoel[int64](tx, subTx, dp.DataPermissionType, dp.DataPermission, dp.DataPermissionScope)
	} else if store.dataPermissionTable["DataPermissionStringModel"] {
		dp := dpi.(*DataPermissionStringCtx)
		tx = switchMoel[string](tx, subTx, dp.DataPermissionType, dp.DataPermission, dp.DataPermissionScope)
	}
	return tx
}

// 使用泛型完成该方法
func switchMoel[T int64 | string](tx *gorm.DB, subTx *gorm.DB, dataPermissionType int8, dataPermission T, dataPermissionScope []T) *gorm.DB {

	if tx.Statement.BuildClauses != nil {
		auth, _ := tx.Get(consts.CTX_AUTH)
		switch tx.Statement.BuildClauses[0] {
		case "SELECT", "UPDATE", "DELETE":
			switch dataPermissionType {
			//本人可见
			case DATA_PERMISSION_ONESELF:
				auth, _ := tx.Get(consts.CTX_AUTH)
				if subTx != nil {
					subTx.Where("created_by = ? ", auth)
				} else {
					tx.Where("created_by = ? ", auth)
				}
			//全部可见
			case DATA_PERMISSION_ALL:
				// 不添加过滤条件
			// 只可见当前部门, 传入当前部门
			case DATA_PERMISSION_ONLY_DEPT:
				if subTx != nil {
					subTx.Where("(dp = ? or created_by = ?)", dataPermission, auth)
				} else {
					tx.Where("(dp = ? or created_by = ?)", dataPermission, auth)
				}
			// 所在部门及子级可见, 自定义部门
			case DATA_PERMISSION_ONLY_DEPT_ALL, DATA_PERMISSION_ONLY_CUSTOM:
				if subTx != nil {
					subTx.Where("(dp in (?) or created_by = ? )", dataPermissionScope, auth)
				} else {
					tx.Where("(dp in (?) or created_by = ? )", dataPermissionScope, auth)
				}
			// 默认权限只有本人可见
			default:
				tx = tx.Where("created_by = ? ", auth)
				if subTx != nil {
					subTx.Where("created_by = ? ", auth)
				} else {
					tx.Where("created_by = ? ", auth)
				}
			}
		case "INSERT":
			tx.Statement.SetColumn("dp", dataPermission, true)
		}

	}
	return tx
}

// 子查询, 用于分页,由子查询的字段和构建条件时,可以设置子查询
// 可以自定义子查询主键字段,
func (store *rdbStore) buildSubQuery(tx *gorm.DB) (sub *gorm.DB) {
	cbI, ok := tx.Get(TX_CONDITION_BUILD_KEY)
	if !ok || cbI == nil {
		return
	}
	columnI, ok := tx.Get(TX_SUB_QUERY_COLUMN_KEY)
	if columnI == nil || !ok {
		return
	}
	// 通过该方法判断是否查询count(*)总数方法,
	// 查询总数去除可能存在的子查询
	if strings.Contains(strings.ToLower(fmt.Sprintf("%v", tx.Statement.Clauses["SELECT"].Expression)), "count(*)") {
		tx.Statement.Clauses["FROM"] = clause.Clause{}
		return

	}
	table := tx.Statement.Table
	col := columnI.(string)
	cb := cbI.(*Condition)

	sub = rdb.BuildQuery(cb).Table(table)
	sub.Where("deleted_at = 0 or deleted_at is null")
	subColumns := fmt.Sprintf("sub%s", col)
	sub = sub.Select(fmt.Sprintf("%s as %s ", col, subColumns))
	orders := "id"
	if cb.Orders != "" {
		orders = cb.Orders
	}
	tx.Table(table).Joins(fmt.Sprintf("join ( ? ) t1 on t1.%s = %s", subColumns, col), sub).Order(orders)
	return

}
