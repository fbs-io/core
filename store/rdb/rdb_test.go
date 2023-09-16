/*
 * @Author: reel
 * @Date: 2023-06-20 07:03:06
 * @LastEditors: reel
 * @LastEditTime: 2023-09-05 19:26:19
 * @Description: 测试通过条件结构体自动完成查询条件的设置
 */
package rdb

import (
	"fmt"
	"reflect"
	"testing"
)

type OrgBase struct {
	OrgName       string `conditions:"like%" json:"org_name" gorm:"comment:组织名称;index"`
	OrgCode       string `conditions:"like" json:"org_code" gorm:"comment:组织代码;uniqueIndex"`
	OrgParentCode string `conditions:"<>" json:"org_parent_code" gorm:"comment:上级组织code;index"`
	OrgPath       string `conditions:"%like" json:"org_path" gorm:"comment:组织路径;index"`
	OrgDeep       int8   `conditions:">=" json:"org_deep" gorm:"comment:组织层级;index"` // 方便定位数据
}

func (org *OrgBase) TableName() string { return "org" }

var (
	data = []*OrgBase{
		{"总经理", "001", "", "001", 1},
		{"财务部", "002", "", "002", 1},
		{"出纳课", "002001", "", "002/002001", 2},
		{"决算课", "002002", "", "002/002002", 2},
		{"资金组", "002001001", "", "002/002001/002001001", 3},
		{"营业部", "003", "", "003", 1},
	}
)

// 配置功能测试 以及注册功能测试
func TestConfig(t *testing.T) {
	// 测试配置功能
	err := rdb.SetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 测试表注册功能
	rdb.Register(&OrgBase{},
		func() error { return rdb.CreateInBatches(data) },
	)
	// DB状态显示测试
	fmt.Println("DB Status: ", rdb.Status())

	// 测试开启功能
	err = rdb.Start()
	if err != nil {
		fmt.Println(err)
	}
	defer rdb.Stop()

	fmt.Println("DB : ", rdb.statTab)

	// DB状态显示测试
	fmt.Println("DB Status: ", rdb.Status())

}

func TestCondition(t *testing.T) {

	// 测试开启功能

	err := rdb.Start()
	if err != nil {
		fmt.Println(err)
	}
	defer rdb.Stop()

	// 测试状态
	fmt.Println("DB Status: ", rdb.Status())

	cb := &Condition{
		PageSize:   10,
		PageNumber: 0,
		Columns:    "org_code",
		Orders:     "org_code",
		Where: map[string]interface{}{
			"org_code like (?)": []string{"002"},
		},
	}
	// 测试根据参数反射类型自动生成参数
	// 模拟参数
	rt := reflect.TypeOf(OrgBase{})
	rv := reflect.New(rt)
	p := rv.Interface().(*OrgBase)
	p.OrgCode = "001"
	p.OrgName = "总经理"
	p.OrgParentCode = "0"
	p.OrgPath = "002"
	p.OrgDeep = 2

	// 根据反射值生成查询条件
	cb2 := GenConditionWithParams(rv)
	fmt.Println(cb2)

	// 根据条件结构体生成查询语句完成查询
	tx := rdb.BuildQuery(cb)
	res := make([]map[string]interface{}, 0, 10)
	err = tx.Table(p.TableName()).Find(&res).Error

	fmt.Println(err, res)

}

func TestDBBuild(t *testing.T) {

	err := rdb.Start()
	if err != nil {
		fmt.Println(err)
	}
	defer rdb.Stop()

	// 模糊查询测试:
	type code struct {
		OrgCode string `conditions:"%like" json:"org_code" `
	}
	rt := reflect.TypeOf(code{})
	rv := reflect.New(rt)
	codep := rv.Interface().(*code)
	codep.OrgCode = "002"
	tx := rdb.BuildQueryWithParams(rv)
	res := make([]map[string]interface{}, 0, 10)
	p := &OrgBase{}

	// 自己输入表名和结果集
	err = tx.Table(p.TableName()).Find(&res).Error
	fmt.Println("模糊查询测试: ", err, res)
	fmt.Println("       ")

	// 测试等于查询
	type code2 struct {
		OrgCode string `json:"org_code" `
	}

	rt = reflect.TypeOf(code2{})
	rv = reflect.New(rt)
	codep2 := rv.Interface().(*code2)
	codep2.OrgCode = "002"
	tx = rdb.BuildQueryWithParams(rv)
	res = make([]map[string]interface{}, 0, 10)
	p = &OrgBase{}

	err = tx.Table(p.TableName()).Find(&res).Error
	fmt.Println("精确查询测试: ", err, res)
	fmt.Println("       ")

	// 测试范围in/notin查询
	type code3 struct {
		OrgCode string `json:"org_code" conditions:"ni"`
	}

	rt = reflect.TypeOf(code3{})
	rv = reflect.New(rt)
	codep3 := rv.Interface().(*code3)
	codep3.OrgCode = "002,001"
	tx = rdb.BuildQueryWithParams(rv)
	res = make([]map[string]interface{}, 0, 10)
	p = &OrgBase{}

	err = tx.Table(p.TableName()).Find(&res).Error
	fmt.Println("范围查询in/notin测试: ", err, res)
	fmt.Println("       ")

	type code4 struct {
		OrgCode string `json:"org_code" conditions:">"`
	}

	// 测试范围<查询
	rt = reflect.TypeOf(code4{})
	rv = reflect.New(rt)
	codep4 := rv.Interface().(*code4)
	codep4.OrgCode = "002"
	tx = rdb.BuildQueryWithParams(rv)
	res = make([]map[string]interface{}, 0, 10)
	p = &OrgBase{}

	err = tx.Table(p.TableName()).Find(&res).Error
	fmt.Println("范围查询</>/>=/<=测试: ", err, res)
	fmt.Println("       ")

	type code5 struct {
		StartDeep int `json:"star_deep" conditions:"org_deep >="`
		EndDeep   int `json:"end_deep" conditions:"org_deep <="`
	}

	// 测试范围<查询
	rt = reflect.TypeOf(code5{})
	rv = reflect.New(rt)
	codep5 := rv.Interface().(*code5)
	codep5.StartDeep = 2
	codep5.EndDeep = 3
	tx = rdb.BuildQueryWithParams(rv)
	res = make([]map[string]interface{}, 0, 10)
	p = &OrgBase{}

	err = tx.Table(p.TableName()).Find(&res).Error
	fmt.Println("范围查询< and > 测试: ", err, res)

}
