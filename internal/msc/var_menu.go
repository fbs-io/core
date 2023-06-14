/*
 * @Author: reel
 * @Date: 2023-06-09 06:02:04
 * @LastEditors: reel
 * @LastEditTime: 2023-06-12 07:11:54
 * @Description: 请填写简介
 */
package msc

var menu = []map[string]interface{}{
    {
        "name":      "overview",
        "path":      "/overview",
        "component": "overview",
        "meta": map[string]interface{}{
            "icon":  "el-icon-menu",
            "title": "服务器概览",
            "affix": false,
        },
    },
    {
        "name":      "service",
        "path":      "/service",
        "component": "service",
        "meta": map[string]interface{}{
            "icon":  "el-icon-menu",
            "title": "服务管理",
            "affix": false,
        },
    },
}
