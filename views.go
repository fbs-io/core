package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/fbs-io/core/store/rdb"
)

const (
	jsonContent = "application/json"
	formContent = "application/x-www-form-urlencoded"

	ViewTypeForm     = "form"
	ViewTypeColumn   = "column"
	ViewTypeTable    = "table"
	ViewRoleShow     = "show"
	viewValueTypeStr = "string"
	viewValueTypeNum = "number"

	// 标签相关
	tagJson    = "json"
	tagForm    = "form"
	tagDesc    = "desc"
	tagBinding = "binding"
	tagDefault = "default"
	tagGorm    = "gorm"
	tagView    = "views"

	// views 标签
	viewSelect     = "select"      // 下拉菜单
	viewMultiple   = "multiple"    // 多选
	viewDisabled   = "disabled"    // 不可选
	viewKey        = "key"         // 主键
	viewSwitcth    = "switch"      // 开关
	viewFilter     = "filter"      // 过滤器, 用于查询条件
	viewSpan       = "span"        // 宽度
	viewsInput     = "input"       // 输入框, 默认值
	viewDepend     = "depend"      // 依赖, 用于查询条件
	viewRaido      = "radio"       // 选择框
	viewDate       = "date"        // 日期选择框
	viewHidden     = "hidden"      // 隐藏字段
	viewsCheckbox  = "checkbox"    // 多选框
	viewTextarea   = "textarea"    // 多行文本框
	viewsFile      = "file"        //文件
	viewWidth      = "width"       // 宽度
	viewHeight     = "height"      // 高度
	viewCalc       = "calc"        // 计算公式
	viewFormat     = "format"      // 格式化
	viewFormatType = "format_type" //格式化类型
	viewRel        = "rel"         // 关联其他值
	viewMoney      = "money"       // 金额格式化
	viewPer        = "%"           // 百分比格式化
	viewCode       = "code"        // 定义字段code, 用于前端显示
	viewSort       = "sort"        // 定义字段是否可以自定义排序, 用于前端显示
	viewFixed      = "fixed"       // 定义字段是否固定, 用于前端显示
	// viewName       = "name"        // 定义字段名称, 用于前端显示

	// 参数相关
	paramsValue      = "value"
	paramsValueType  = "value_type"
	valueInt         = "int"
	paramsValueBool  = "bool"
	valueFloat       = "float"
	valueTypeDecimal = "decimal"
)

// 根据参数结构体生成API参数,
//
// 当前仅支持 form 和 json 两种格式
//
// 根据参数第一个字段的标签判断content-type类型
//
// 如果参数结构体定义多个参数格式, 将其他参数无法正确使用
//
// TODO: 支持文件/多文件参数定义

func (r *router) genViews(method, pathName, viewType, ViewItemCode string, rt reflect.Type) (viewsList []*Views, contentType string) {
	if rt == nil {
		return
	}
	// 指针类型转换
	if rt.Kind() == reflect.Ptr {
		fmt.Println("指针类型")
		rt = rt.Elem()
	}
	viewsList = make([]*Views, 0, 100)
	// viewsMap := make(map[string]*Views, 100)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// 获取前端参数名称
		key := field.Tag.Get(tagForm)
		contentTypeCustom := formContent
		if key == "" {
			contentTypeCustom = jsonContent
			key = field.Tag.Get(tagJson)
		}
		// TODO: 增加其他类型检查

		// 如果没有获取到key, 说明该参数无效, 跳过不在录入
		if key == "" {
			continue
		}

		// 通过第一个获取到参数的结构体的类型作为整个请求的content_type
		if contentType == "" {
			contentType = contentTypeCustom
		}
		view := &Views{
			ResourceCode: r.resource.Code,
			ViewCode:     pathName,
			ViewItemCode: ViewItemCode,
			ViewType:     viewType,
			ViewRole:     method,
			ValueType:    viewValueTypeStr,
			Code:         key,
			FormatType:   viewsInput,
		}

		// 前端参数的数据类型
		typeStr := field.Type.String()
		// 后端参数类型转换为前端的参数类型
		if strings.Contains(typeStr, valueInt) {
			view.ValueType = viewValueTypeNum
		} else if strings.Contains(typeStr, valueFloat) {
			view.ValueType = viewValueTypeNum
		} else if strings.Contains(typeStr, valueTypeDecimal) {
			view.ValueType = viewValueTypeNum
		}

		// 用于前端API文档中的默认值
		view.Default = field.Tag.Get(tagDefault)

		if strings.Contains(typeStr, "[]") {
			view.ValueType = "[]:" + view.ValueType
			view.Default = "[]"
		}

		view.Name = field.Tag.Get(tagDesc)
		view.Rules = field.Tag.Get(tagBinding)

		r.genViewsTag(strings.Split(field.Tag.Get(tagView), ";"), view)

		r.core.Views = append(r.core.Views, view)
		// var tableViews []*Views
		// 如果参数是切片, 递归处理切片元素的结构体字段
		if field.Type.Kind() == reflect.Slice {
			elemType := field.Type.Elem()

			// 处理切片元素为指针的情况
			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}

			// 解析切片元素的结构体字段
			if elemType.Kind() == reflect.Struct {
				view.ViewType = ViewTypeTable

				r.genViews(method, pathName, ViewTypeTable, key, elemType)
			}
		}
	}
	return
}

func (r *router) genViewsTag(viewTags []string, view *Views) {
	for _, val := range viewTags {
		kvs := strings.Split(val, ":")
		k := kvs[0]
		if len(k) == 0 {
			continue
		}

		switch kvs[0] {
		case viewWidth:
			i, _ := strconv.Atoi(kvs[1])
			if i > 0 {
				view.Width = int16(i)
			}
		case viewHeight:
			i, _ := strconv.Atoi(kvs[1])
			if i > 0 {
				view.Height = int16(i)
			}
		case viewSelect:
			view.FormatType = viewSelect
			view.Format = view.Code
			if len(kvs) > 1 {
				view.Format = kvs[1]
			}
		case viewMultiple:
			view.Multiple = "Y"
		case viewDisabled:
			view.Disabled = "Y"
		case viewKey:
			view.Key = "Y"
			view.Disabled = "Y"
		case viewFormat:
			if len(kvs) > 1 {
				view.Format = kvs[1]
			} else {
				view.Format = view.Code
			}
		case viewHidden:
			view.Hidden = 1
		case viewMoney:
			view.Format = viewMoney
		case viewPer:
			view.Format = viewPer
		case viewSwitcth:
			// 如果文字类型, 设置默认值Y,N
			view.Format = view.Code
			view.FormatType = viewSwitcth
			view.CustomValue = []any{"Y", "N"}
			view.Default = "N"
			// 如果数字类型, 设置默认值1,-1
			if view.ValueType == viewValueTypeNum {
				view.CustomValue = []any{1, -1}
				view.Default = "-1"
			}

			// 支持自定义, 格式: switch:1,2
			if len(kvs) > 1 {
				vaList := strings.Split(kvs[1], ",")
				view.CustomValue = make(rdb.ModeListJson, 0, len(vaList))
				for _, val := range vaList {
					if view.ValueType == viewValueTypeNum {
						num, _ := strconv.Atoi(val)
						view.CustomValue = append(view.CustomValue, num)
					} else {
						view.CustomValue = append(view.CustomValue, val)
					}
				}
			}
		case viewRaido:
			view.FormatType = viewRaido
			view.Format = view.Code
			if len(kvs) > 1 {
				view.Format = kvs[1]
			}
		case viewDate:
			view.FormatType = viewDate
			view.Format = "YYYY-MM-DD"
			if len(kvs) > 1 {
				view.Format = kvs[1]
			}
		case viewFilter:
			view.Filter = kvs[1]
		case viewSpan:
			span, _ := strconv.Atoi(kvs[1])
			view.Width = int16(span)
		case viewDepend:
			view.Depend = kvs[1]
		case viewTextarea:
			view.FormatType = viewTextarea
		case viewCalc:
			view.Calc = kvs[1]
		case viewRel:
			view.Related = kvs[1]
		case tagDefault:
			view.Default = kvs[1]
		case view.Code:
			if len(kvs) > 1 {
				view.Code = kvs[1]
			}
		case viewSort:
			view.IsOrder = 1
		case viewFixed:
			view.Fixed = "Y"
		}
	}

	if view.Name == "" {
		sub := r.core.ViewsMap[fmt.Sprintf("%s:%s", view.ResourceCode, view.Code)]
		if sub != nil {
			view.Name = sub.Name
		}
	}
	if view.Name == "" {
		view.Name = strings.ToUpper(view.Code)
	}
	if view.FormatType == viewsInput && view.ValueType == viewValueTypeNum {
		view.FormatType = viewValueTypeNum
	}
}

func (r *router) genViewColumns(rt reflect.Type, viewType, viewItemCode, viewRole string, opt *SetViewOptions) (view *Views) {

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		// 如果参数是切片, 递归处理切片元素的结构体字段

		view = r.genViewCol(field, viewType, viewItemCode, viewRole, opt)
		if view == nil {
			continue
		}
		key := fmt.Sprintf("%s:%s", view.ResourceCode, view.Code)

		if r.core.ViewsMap[key] == nil {
			r.core.ViewsMap[key] = view
		}
		r.core.Views = append(r.core.Views, view)

		// 如果参数是切片, 递归处理切片元素的结构体字段
		if field.Type.Kind() == reflect.Slice {
			elemType := field.Type.Elem()

			// 处理切片元素为指针的情况
			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}

			// 解析切片元素的结构体字段
			if elemType.Kind() == reflect.Struct {
				view.ViewType = ViewTypeTable
				r.genViewColumns(elemType, ViewTypeTable, view.Code, viewRole, opt)

			}
		}

	}

	return
}

func (r *router) genViewCol(field reflect.StructField, viewType, viewItemCode, viewRole string, opt *SetViewOptions) (view *Views) {
	view = &Views{
		ResourceCode: r.resource.Code,
		ViewCode:     r.resource.Name,
		ViewType:     viewType,
		ViewItemCode: viewItemCode,
		ViewRole:     viewRole,
		ValueType:    "string",
		Width:        120,
		Height:       0,
		IsOrder:      1,
		Filter:       "",
		Fixed:        "N",
		Account:      "system",
	}
	// 获取前端参数名称
	view.Code = field.Tag.Get(tagJson)
	// TODO: 增加其他类型检查

	// 处理前端参数名称, 从gorm标签或者自定义的desc标签获取
	view.Name = field.Tag.Get(tagDesc)
	if view.Name == "" {
		for _, tag := range strings.Split(field.Tag.Get(tagGorm), ";") {
			if strings.Contains(tag, "comment:") {
				view.Name = strings.Split(tag, ":")[1]
				break
			}
		}
	}
	typeStr := field.Type.String()
	// 后端参数类型转换为前端的参数类型
	if strings.Contains(typeStr, valueInt) {
		view.ValueType = viewValueTypeNum
	} else if strings.Contains(typeStr, valueFloat) {
		view.ValueType = viewValueTypeNum
	} else if strings.Contains(typeStr, paramsValueBool) {
		view.ValueType = paramsValueBool
	} else if strings.Contains(typeStr, paramsValueBool) {
		view.ValueType = paramsValueBool
	}
	if view.ValueType == viewValueTypeNum {
		view.FormatType = "number"
	}
	viewsTag := field.Tag.Get(tagView)
	if viewsTag != "" {
		r.genViewsTag(strings.Split(viewsTag, ";"), view)
	}
	// 对view标签进行处理

	if view.Code == "" {
		return nil
	}
	if opt.ViewCode != "" {
		view.ViewCode = opt.ViewCode
	}

	if opt.ColumnWidth > 0 {
		view.Width = opt.ColumnWidth
	}

	if opt.ColumnHeight > 0 {
		view.Height = opt.ColumnHeight
	}

	if opt.ColumnIsHidden > 0 {
		view.Hidden = opt.ColumnIsHidden
	}
	if opt.ColumnIsOrder > 0 {
		view.IsOrder = opt.ColumnIsOrder
	}
	if opt.ColumnFilter != "" {
		view.Filter = opt.ColumnFilter
	}
	if opt.ColumnFixed != "" {
		view.Fixed = opt.ColumnFixed
	}
	if opt.Account != "" {
		view.Account = opt.Account
	}
	return
}
