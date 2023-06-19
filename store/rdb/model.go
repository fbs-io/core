/*
 * @Author: reel
 * @Date: 2023-06-16 06:04:12
 * @LastEditors: reel
 * @LastEditTime: 2023-06-20 06:34:14
 * @Description: 定义常用的模型用于快速开发
 */

package rdb

import "gorm.io/gorm"

// 通用的字段处理
type Model struct {
    ID        uint `json:"id" gorm:"primaryKey;"`                 // id
    CreateAT  uint `json:"create_at" gorm:"autoCreateTime:milli"` // 创建时间
    CreateBy  uint `json:"create_by"`                             // 创建人
    UpdatedAT uint `json:"update_at" gorm:"autoUpdateTime:milli"` // 修改时间
    UpdatedBy uint `json:"update_by"`                             // 修改人
    DeletedAT uint `json:"delete_at" gorm:"index"`                // 删除时间
    DeleteBy  uint `json:"delete_by"`                             // 删除人
    Status    int8 `json:"status" gorm:"index;comment:状态"`        // 是否启用 1 表示启用 -1表示失效
}

// 通用查询查询字段处理
type ModelQuery struct {
    ID     uint `json:"id"`
    Status int8 `json:"status"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
    m.Status = 1
    return nil
}
