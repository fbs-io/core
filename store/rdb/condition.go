/*
 * @Author: reel
 * @Date: 2023-06-15 06:55:41
 * @LastEditors: reel
 * @LastEditTime: 2024-07-07 13:50:05
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
	eq = "="
	ne = "<>"
	lt = "<"
	le = "<="
	qe = ">="
	qt = ">"
	in = "in"

	notin     = "ni"
	like      = "like"
	likeRight = "like%"
	likeLeft  = "%like"
)

type Condition struct {
	PageSize    int
	PageNumber  int
	Columns     string
	TableName   string
	Orders      string
	Where       map[string]reflect.Value
	QryDelete   bool
	IsSharding  bool
	ShardingKey string
}

func NewCondition() *Condition {
	return &Condition{
		PageSize:   10,
		PageNumber: 0,
		Where:      make(map[string]reflect.Value, 0),
	}
}

// 通过传入条件, 自动完成gorm的语句生成
//
// 不适用于大表的翻页查询, 大表查询请优化表结构
func (store *rdbStore) BuildQuery(cb *Condition) (tx *gorm.DB) {
	tx = store.DB()
	for key, value := range cb.Where {
		values := make([]interface{}, 0, 100)
		switch value.Kind() {
		// 对切片处理
		case reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				val := value.Index(i)
				values = append(values, val.Interface())
			}
			tx = tx.Where(key, values)
		default:
			tx = tx.Where(key, value.Interface())

		}
	}

	// 限定最大获取输了
	if cb.PageSize > 1000 {
		cb.PageSize = 1000
	}
	// 限定最大获取输了
	cb.PageNumber = cb.PageNumber - 1
	if cb.PageNumber < 0 {
		cb.PageNumber = 0
	}
	if cb.Orders != "" {
		tx = tx.Order(cb.Orders)
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
		if paramsType.Field(i).Name == "ShardingModel" {
			cb.IsSharding = true
		}
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
			cb.Orders = valueType.String()
		case "coloums":
			cb.Columns = valueType.String()
		default:
			ckey := "%s %s"

			// 处理查询在某个范围, 如 1< age <10
			conditions := strings.Split(tag.Get("conditions"), " ")
			condition := conditions[0]
			if len(conditions) >= 2 {
				key = conditions[1]
			}
			// 不生成查询条件
			if condition == "-" {
				continue
			}

			// 对模糊查询的单独处理
			switch condition {
			// 条件为空, 默认是等于
			case "":
				condition = fmt.Sprintf("%s (?)", eq)
			// not in
			case notin:
				condition = "not in (?)"
			case like:
				condition = "like (?)"
				valueType.SetString(fmt.Sprintf("%%%v%%", valueType.Interface()))
			case likeLeft:
				condition = "like (?)"
				valueType.SetString(fmt.Sprintf("%%%v", valueType.Interface()))
			case likeRight:
				condition = "like (?)"
				valueType.SetString(fmt.Sprintf("%v%%", valueType.Interface()))
			default:
				condition = fmt.Sprintf("%s (?)", condition)
			}

			cb.Where[fmt.Sprintf(ckey, key, condition)] = valueType
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
//
// 请注意, 该方法不适用于多库分区的查询构建
func (store *rdbStore) BuildQueryWithParams(params reflect.Value) *gorm.DB {
	cb := GenConditionWithParams(params)
	tx := store.BuildQuery(cb)
	return tx
}
