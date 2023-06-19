/*
 * @Author: reel
 * @Date: 2023-06-15 06:55:41
 * @LastEditors: reel
 * @LastEditTime: 2023-06-15 07:06:33
 * @Description: 根据条件结构体, 自动构建查询语句, 并返回gorm.DB, 用于扩展
 */
package rdb

import (
    "strings"

    "gorm.io/gorm"
)

type Condition struct {
    Limit      int
    Offset     int
    IsQueryAll bool // 是否查询全部, 默认不被允许
    Column     string
    TableName  string
    Orders     []string
    Where      map[string][]interface{}
}

func (rdb *rdbStore) BuildQuery(cb *Condition) (tx *gorm.DB) {
    tx = rdb.db
    sub := rdb.db
    for k, v := range cb.Where {
        if len(v) == 0 {
            sub = sub.Where(k)
        } else {
            sub = sub.Where(k, v...)
        }
    }

    // 子查询用于快速分页查询
    sub = sub.Table(cb.TableName).Select("id as subid")
    sub = sub.Where("status = ? ", 1).Order(strings.Join(cb.Orders, ", "))
    if !cb.IsQueryAll {
        sub = sub.Limit(cb.Limit).Offset(cb.Offset)
    }
    tx = tx.Table(cb.TableName).Joins("join ( ? ) t1 on t1.subid = id", sub)

    // 设置查询的名称
    tx = tx.Select(cb.Column)
    return tx
}
