/*
 * @Author: reel
 * @Date: 2023-06-15 06:55:41
 * @LastEditors: reel
 * @LastEditTime: 2023-06-23 08:05:02
 * @Description: 根据条件结构体, 自动构建查询语句, 并返回gorm.DB, 用于扩展
 */
package rdb

import (
    "fmt"
    "reflect"
    "strings"

    "gorm.io/gorm"
)

const (
    eq  = "="
    ne  = "<>"
    lt  = "<"
    le  = "<="
    qe  = ">="
    qt  = ">"
    in  = "in"

    notin     = "ni"
    like      = "like"
    likeRight = "like%"
    likeLeft  = "%like"
)

type Condition struct {
    PageSize   int
    PageNumber int
    Columns    string
    TableName  string
    Orders     []string
    Where      map[string]interface{}
}

func NewCondition() *Condition {
    return &Condition{
        PageSize:   10,
        PageNumber: 0,
        Orders:     make([]string, 0, 10),
        Where:      make(map[string]interface{}, 0),
    }
}

// 通过传入条件, 自动完成gorm的语句生成
//
// 此方法适用于表中有ID(主键)的字段, 优化了翻页查询性能
func (rdb *rdbStore) BuildQueryByID(cb *Condition) (tx *gorm.DB) {
    tx = rdb.db
    sub := rdb.db
    for k, v := range cb.Where {
        sub = sub.Where(k, v)
    }
    // 限定最大获取输了
    if cb.PageSize > 1000 {
        cb.PageSize = 1000
    }
    // 子查询用于快速分页查询
    sub = sub.Table(cb.TableName).Select("id as subid")
    sub = sub.Where("status = ? ", 1).Order(strings.Join(cb.Orders, ", "))
    sub = sub.Limit(cb.PageSize).Offset(cb.PageNumber * cb.PageSize)

    tx = tx.Table(cb.TableName).Joins("join ( ? ) t1 on t1.subid = id", sub)

    // 设置查询的名称
    if cb.Columns != "" {
        tx = tx.Select(cb.Columns)
    }
    return tx
}

// 通过传入条件, 自动完成gorm的语句生成
//
// 不适用于大表的翻页查询, 大表查询请优化表结构
func (rdb *rdbStore) BuildQuery(cb *Condition) (tx *gorm.DB) {
    tx = rdb.db
    for k, v := range cb.Where {
        tx = tx.Where(k, v)
    }

    // 限定最大获取输了
    if cb.PageSize > 1000 {
        cb.PageSize = 1000
    }

    tx = tx.Limit(cb.PageSize).Offset(cb.PageNumber * cb.PageSize)

    tx = tx.Table(cb.TableName)

    // 设置查询的名称
    if cb.Columns != "" {
        tx = tx.Select(cb.Columns)
    }
    return tx
}

// 根据请求参数构建查询条件
//
// 其中表名和返回值需要手动添加
//
// 仅适用单表的简单where-and条件查询, 不适用于复杂关联查询
//
// 复杂业务查询须手动处理或构建查询视图
func GenConditionWithParams(params reflect.Value) *Condition {
    cb := NewCondition()
    paramsType := params.Type().Elem()
    for i := 0; i < params.Elem().NumField(); i++ {
        if params.Elem().Field(i).IsZero() {
            continue
        }

        // 判断json或form字段
        tag := paramsType.Field(i).Tag
        key := tag.Get("json")
        if key == "" {
            key = tag.Get("form")
        }
        if key == "" {
            continue
        }

        valueType := params.Elem().Field(i)
        switch key {
        case "page_size":
            cb.PageSize = int(valueType.Int())
        case "page_num":
            cb.PageNumber = int(valueType.Int())
        case "orders":
            cb.Orders = strings.Split(valueType.String(), ",")
        case "coloums":
            cb.Columns = valueType.String()
        default:
            value := valueType.Interface()
            ckey := "%s %s"

            // 处理查询在某个范围, 如 1< age <10
            conditions := strings.Split(tag.Get("conditions"), " ")
            condition := conditions[0]
            if len(conditions) >= 2 {
                key = conditions[1]
            }
            switch condition {
            // 条件为空, 默认是等于
            case "":
                condition = fmt.Sprintf("%s (?)", eq)
            // not in
            case notin:
                condition = "not in (?)"
            case like:
                condition = fmt.Sprintf("like '%%%v%%' ", value)
                value = nil
            case likeLeft:
                condition = fmt.Sprintf("like '%%%v' ", value)
                value = nil
            case likeRight:
                condition = fmt.Sprintf("like '%v%%' ", value)
                value = nil
            default:
                condition = fmt.Sprintf("%s (?)", condition)
            }
            cb.Where[fmt.Sprintf(ckey, key, condition)] = value
        }
    }
    return cb
}

// 根据请求参数构建查询条件
//
// 其中表名和返回值需要手动添加
//
// 仅适用单表的简单where-and条件查询, 不适用于复杂关联查询
//
// 复杂业务查询须手动处理或构建查询视图
func (db *rdbStore) BuildQueryWithParams(params reflect.Value) *gorm.DB {
    paramsType := params.Type().Elem()
    tx := db.DB()
    pageNumber := 0
    pageSize := 10
    for i := 0; i < params.Elem().NumField(); i++ {
        if params.Elem().Field(i).IsZero() {
            continue
        }
        // 判断json或form字段
        tag := paramsType.Field(i).Tag
        key := tag.Get("json")
        if key == "" {
            key = tag.Get("form")
        }
        if key == "" {
            continue
        }

        valueType := params.Elem().Field(i)
        switch key {
        case "page_size":
            pageSize = int(valueType.Int())
        case "page_num":
            pageNumber = int(valueType.Int())
        case "orders":
            tx = tx.Order(strings.Split(valueType.String(), ","))
        case "coloums":
            tx = tx.Select(valueType.String())
        default:
            value := valueType.Interface()
            ckey := "%s %s"
            // 处理查询在某个范围, 如 1< age <10
            conditions := strings.Split(tag.Get("conditions"), " ")
            condition := conditions[0]
            if len(conditions) >= 2 {
                key = conditions[0]
                condition = conditions[1]
            }

            switch condition {
            // 条件为空, 默认是等于
            case "":
                condition = fmt.Sprintf("%s (?)", eq)
            // not in
            case notin:
                condition = "not in (?)"
                value = interface{}(strings.Split(valueType.String(), ","))
            case in:
                condition = " in (?)"
                value = interface{}(strings.Split(valueType.String(), ","))
            case like:
                condition = fmt.Sprintf("like '%%%v%%' ", value)
                value = nil
            case likeLeft:
                condition = fmt.Sprintf("like '%%%v' ", value)
                value = nil
            case likeRight:
                condition = fmt.Sprintf("like '%v%%' ", value)
                value = nil
            default:
                condition = fmt.Sprintf("%s (?)", condition)
            }
            tx = tx.Where(fmt.Sprintf(ckey, key, condition), value)

        }
    }
    if pageSize > 1000 {
        pageSize = 1000
    }
    tx = tx.Limit(pageSize).Offset(pageNumber * pageSize)
    return tx
}
