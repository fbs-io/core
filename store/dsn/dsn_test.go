/*
 * @Author: reel
 * @Date: 2023-06-21 05:47:37
 * @LastEditors: reel
 * @LastEditTime: 2023-06-21 06:01:28
 * @Description: 测试DSN功能
 */
package dsn

import (
    "fmt"
    "testing"
)

func TestDsn(t *testing.T) {
    dbdsn := NewDBDsn()
    cachedsn := NewCacheDsn()

    fmt.Println("数据库默认连接信息", dbdsn.Link())
    fmt.Println("缓存默认连接信息", cachedsn.Link())

    SetName("core")(dbdsn)
    SetName("cache")(cachedsn)

    fmt.Println("数据库修改默认名称", dbdsn.Link())
    fmt.Println("缓存修改默认名称", cachedsn.Link())

    SetType("postgres")(dbdsn)
    SetType("redis")(cachedsn)

    fmt.Println("数据库修改数据库类型", dbdsn.Link())
    // redis 不直接使用dsn生成的link,
    fmt.Println("缓存修改数据库类型", cachedsn.Link())

}
