/*
 * @Author: reel
 * @Date: 2023-06-16 06:04:12
 * @LastEditors: reel
 * @LastEditTime: 2023-11-09 06:26:10
 * @Description: 定义常用的模型用于快速开发
 */

package rdb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fbs-io/core/pkg/consts"
	"gorm.io/gorm"
)

// 通用的字段处理
type Model struct {
	ID        uint      `json:"id" gorm:"primaryKey;"`                    // id
	CreatedAT uint      `json:"created_at" gorm:"autoCreateTime:milli"`   // 创建时间
	CreatedBy string    `json:"created_by"`                               // 创建人
	UpdatedAT uint      `json:"updated_at" gorm:"autoUpdateTime:milli"`   // 修改时间
	UpdatedBy string    `json:"updated_by"`                               // 修改人
	DeletedAT DeletedAt `json:"deleted_at" gorm:"index;softDelete:milli"` // 删除时间
	DeletedBy string    `json:"deleted_by"`                               // 删除人
	Status    int8      `json:"status" gorm:"index;comment:状态;default:1"` // 是否启用 1 表示启用 -1表示失效
}

// 通用查询查询字段处理
type ModelQuery struct {
	ID     uint `json:"id"`
	Status int8 `json:"status"`
}

func (m *Model) getAuth(tx *gorm.DB) string {
	auth, ok := tx.Get(consts.CTX_AUTH)
	if ok {
		return auth.(string)
	}
	return ""
}
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.Status = 1
	m.CreatedBy = m.getAuth(tx)
	return nil
}

func (m *Model) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedBy = m.getAuth(tx)
	return nil
}

// func (m *Model) BeforeDelete(tx *gorm.DB) error {
// 	m.DeletedBy = m.getAuth(tx)
// 	return nil
// }

type ModeMapJson map[string]interface{}

func (j *ModeMapJson) Scan(value interface{}) error {

	bytes, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var mJson = make(ModeMapJson, 10)
	err := json.Unmarshal([]byte(bytes), &mJson)
	*j = mJson
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j ModeMapJson) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

type ModeListJson []interface{}

func (j *ModeListJson) Scan(value interface{}) error {

	bytes, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var mJson = make(ModeListJson, 0, 10)
	err := json.Unmarshal([]byte(bytes), &mJson)
	*j = mJson
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j ModeListJson) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

// 定义分区字段
// 用于表分区使用
type ShardingModel struct {
	ShadingKey string `json:"-" gorm:"column:sk;index;comment:分区"`
}

func (m *ShardingModel) TableName(table string) string {
	if m.ShadingKey == "" {
		return table
	}
	return fmt.Sprintf("%s_%s", table, m.ShadingKey)
}

func (m *ShardingModel) txSharding(tx *gorm.DB) *gorm.DB {
	sk, ok := tx.Get(consts.CTX_SHARDING_KEY)
	if ok {
		m.ShadingKey = sk.(string)
	}
	tableI, ok := tx.Get(TX_SHADING_TABLE_KEY)
	if ok && tableI != nil {
		tableList, _ := tableI.(map[string][]string)
		model, _ := tx.Get(TX_SHADING_MODEL_KEY)
		table := tx.Statement.Table
		if model.(int8) == SHADING_MODEL_TABLE && tableList[table] != nil {
			tx.Statement.Table = m.TableName(table)
		}
	}
	return tx
}

func (m *ShardingModel) BeforeCreate(tx *gorm.DB) error {
	m.txSharding(tx)
	return nil
}

func (m *ShardingModel) BeforeUpdate(tx *gorm.DB) error {
	m.txSharding(tx)
	return nil
}
func (m *ShardingModel) BeforeDelete(tx *gorm.DB) error {
	m.txSharding(tx)
	return nil
}

// 定义数据权限模型
// 用于有数据权限需求的查询, 创建等, 和分区表的原理类似, 但不涉及到分表操作
// 适用于string类型的数据
type DataPermissionStringModel struct {
	DataPermission string `json:"-" gorm:"column:dp;index;comment:数据权限"`
}

type DataPermissionStringCtx struct {
	DataPermissionType  int8 //
	DataPermission      string
	DataPermissionScope []string
}

// 定义数据权限模型
// 用于有数据权限需求的查询, 创建等, 和分区表的原理类似, 但不涉及到分表操作
// 适用于int类型的数据
type DataPermissionIntModel struct {
	DataPermission int64 `json:"-" gorm:"column:dp;index;comment:数据权限"`
}

type DataPermissionIntCtx struct {
	DataPermissionType  int8 //
	DataPermission      int64
	DataPermissionScope []int64
}
