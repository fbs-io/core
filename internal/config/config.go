/*
 * @Author: reel
 * @Date: 2023-05-16 20:17:56
 * @LastEditors: reel
 * @LastEditTime: 2023-06-05 07:57:28
 * @Description: 系统配置相关操作
 */
package config

import (
    "encoding/json"

    "github.com/fbs-io/core/internal/pem"
    "github.com/fbs-io/core/pkg/encrypt"
    "github.com/fbs-io/core/pkg/errorx"
)

type Config struct {
    IsLoad   bool   `json:"-"`
    Port     string `json:"port" binding:"required"`
    User     string `json:"user" binding:"required"`
    Pwd      string `json:"password" binding:"required"`
    Pwd2     string `json:"password2" binding:"required,eqfield=Pwd"`
    DataPath string `json:"data_path" binding:"required"` // 项目数据路径, 包括日志, 本地缓存, 本地db, 本地其他文件等
    // Email    string `json:"email" binding:"required"`
    // DB配置
    DbType string `json:"db_type" binding:"required"`
    DBName string `json:"db_name" binding:"required"`
    DBHost string `json:"db_host" `
    DBPort string `json:"db_port" `
    DBUser string `json:"db_user" `
    DBPwd  string `json:"db_pwd" `
    // 缓存配置
    CacheType string `json:"cache_type" binding:"required"`
    CacheName string `json:"cache_name" binding:"required"`
    CacheHost string `json:"cache_host"`
    CachePort string `json:"cache_port"`
    CacheUser string `json:"cache_user"`
    CachePwd  string `json:"cache_pwd"`
}

// 加载配置
// 配置文件获取及字符转换封装在函数内
func (c *Config) Load() error {
    pems, err := pem.GetPems()
    if err != nil {
        return errorx.Wrap(err, "无法加载配置文件")
    }
    confB, err := encrypt.InternalDecode(pems)
    if err != nil {
        return errorx.Wrap(err, "无法加载正确的配置文件, 请删除配置文件后重新设置")
    }

    err = json.Unmarshal(confB, &c)
    if err != nil {
        return errorx.Wrap(err, "无法加载正确的配置文件, 请删除配置文件后重新设置")
    }
    return nil
}

func (c *Config) Dump() error {
    data, err := json.Marshal(c)
    if err != nil {
        return errorx.Wrap(err, "配置转JSON错误")
    }
    str, err := encrypt.InternalEncode(data)
    if err != nil {

        return errorx.Wrap(err, "配置转码错误")
    }
    err = pem.UpdatePems(str)
    if err != nil {

        return errorx.Wrap(err, "配置更新到文件失败")
    }
    return nil
}
