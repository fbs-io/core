/*
 * @Author: reel
 * @Date: 2023-06-15 06:55:41
 * @LastEditors: reel
 * @LastEditTime: 2025-10-12 20:18:29
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
	eqKey        = "eq"    // 等于
	neKey        = "ne"    // 不等于
	ltKey        = "lt"    // 小于
	leKey        = "le"    // 小于或等于
	qeKey        = "qe"    // 大于或等于
	qtKey        = "qt"    // 大于
	inKey        = "in"    // 包含
	niKey        = "ni"    // 不包含
	likeKey      = "like"  // 模糊查询
	likeLeftKey  = "like%" //	左侧模糊查询
	likeRightKey = "%like" //	右侧模糊查询

	eqVal        = "= ?"
	neVal        = "<> ?"
	ltVal        = "< ?"
	leVal        = "<= ?"
	qeVal        = ">= ?"
	qtVal        = "> ?"
	inVal        = "in ?"
	niVal        = "not in (?)"
	like         = "like ? "
	likeLeftVal  = "like ? "
	likeRightVal = "like ? "

	and = "&"
	or  = "|"
)

var (
	ckeys = map[string]string{
		eqKey:        eqVal,
		neKey:        neVal,
		ltKey:        ltVal,
		leKey:        leVal,
		qeKey:        qeVal,
		qtKey:        qtVal,
		inKey:        inVal,
		niKey:        niVal,
		likeKey:      like,
		likeRightKey: likeRightVal,
		likeLeftKey:  likeLeftVal,
	}
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
		values := make([]any, 0, 100)
		switch value.Kind() {
		// 对切片处理
		case reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				val := value.Index(i)
				values = append(values, val.Interface())
			}
			tx = tx.Where(key, values)
		default:
			count := strings.Count(key, "?")
			valuesList := make([]any, 0, count)
			for i := 0; i < count; i++ {
				valuesList = append(valuesList, value.Interface())

			}
			tx = tx.Where(key, valuesList...)

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
			ckey := "%s %s %s"
			ckeyOr := "%s or %s %s "
			ckeyAnd := "%s or %s %s "

			conditions := strings.Split(tag.Get("conditions"), "=")
			// 判断是否有查询条件, 如果没有
			ck := conditions[0]
			if len(conditions) >= 2 {
				key = conditions[1]
			}
			// 不生成查询条件
			if ck == "-" {
				continue
			}
			condition := ckeys[ck]
			if condition == "" {
				condition = ckeys[eqKey]
			}
			// // 对模糊查询的值单独处理
			switch ck {
			case likeKey:
				valueType.SetString(fmt.Sprintf("%%%v%%", valueType.Interface()))
			case likeRightKey:
				valueType.SetString(fmt.Sprintf("%%%v", valueType.Interface()))
			case likeLeftKey:
				valueType.SetString(fmt.Sprintf("%v%%", valueType.Interface()))
			}
			// 对or和and的值单独处理
			keys := strings.Split(key, or)
			keyAll := ""
			keyAll = fmt.Sprintf(ckey, keyAll, keys[0], condition)

			if len(keys) > 1 {
				for _, k := range keys[1:] {
					keyAll = fmt.Sprintf(ckeyOr, keyAll, k, condition)
				}
			} else {
				keys = strings.Split(key, and)
				for _, k := range keys[1:] {
					keyAll = fmt.Sprintf(ckeyAnd, keyAll, k, condition)
				}
			}

			cb.Where[keyAll] = valueType
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
