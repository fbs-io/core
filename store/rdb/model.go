/*
 * @Author: reel
 * @Date: 2023-06-16 06:04:12
 * @LastEditors: reel
 * @LastEditTime: 2023-08-22 06:59:41
 * @Description: 定义常用的模型用于快速开发
 */

package rdb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 通用的字段处理
type Model struct {
	ID        uint   `json:"id" gorm:"primaryKey;"`                 // id
	CreateAT  uint   `json:"create_at" gorm:"autoCreateTime:milli"` // 创建时间
	CreateBy  string `json:"create_by"`                             // 创建人
	UpdatedAT uint   `json:"update_at" gorm:"autoUpdateTime:milli"` // 修改时间
	UpdatedBy string `json:"update_by"`                             // 修改人
	DeletedAT uint   `json:"delete_at" gorm:"index"`                // 删除时间
	DeleteBy  string `json:"delete_by"`                             // 删除人
	Status    int8   `json:"status" gorm:"index;comment:状态"`        // 是否启用 1 表示启用 -1表示失效
}

// 通用查询查询字段处理
type ModelQuery struct {
	ID     uint `json:"id"`
	Status int8 `json:"status"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.Status = 1
	auth, ok := tx.Get("auth")
	if ok {
		m.CreateBy = auth.(string)
	}
	return nil
}

func (m *Model) BeforeUpdate(tx *gorm.DB) error {
	auth, ok := tx.Get("auth")
	if ok {
		m.UpdatedBy = auth.(string)
	}
	return nil
}

func (m *Model) BeforeDelete(tx *gorm.DB) error {
	auth, ok := tx.Get("auth")
	if ok {
		m.DeleteBy = auth.(string)
	}

	return nil
}

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
