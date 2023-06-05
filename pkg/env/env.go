/*
 * @Author: reel
 * @Date: 2023-05-11 19:52:08
 * @LastEditors: reel
 * @LastEditTime: 2023-05-17 07:53:17
 * @Description: 用于配置app环境变量
 */
package env

import (
    "flag"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm/logger"
)

const (
    ENV_APP_NAME         = "DDNS"
    ENV_MODE_DEV         = "dev"
    ENV_MODE_FAT         = "fat"
    ENV_MODE_PRO         = "pro"
    ENV_MSC_DEPAULT_PORT = ":6018"
)

var (
    active Environment
    dev    Environment = &environment{
        value:        ENV_MODE_DEV,
        ginMode:      gin.DebugMode,
        logLevel:     logrus.DebugLevel,
        project:      ENV_APP_NAME,
        gormLogLevel: logger.Info,
        then:         time.Now().Unix(),
    }
    fat Environment = &environment{
        value:        ENV_MODE_FAT,
        ginMode:      gin.ReleaseMode,
        logLevel:     logrus.DebugLevel,
        project:      ENV_APP_NAME,
        gormLogLevel: logger.Warn,
        then:         time.Now().Unix(),
    }
    pro Environment = &environment{
        value:        ENV_MODE_PRO,
        ginMode:      gin.ReleaseMode,
        logLevel:     logrus.InfoLevel,
        project:      ENV_APP_NAME,
        gormLogLevel: logger.Error,
        then:         time.Now().Unix(),
    }
)

var _ Environment = (*environment)(nil)

type environment struct {
    value        string          // 环境变量配置
    super        bool            // 是否为超级用户
    mscAddr      string          // 管理中心http监听地址
    ginMode      string          // gin运行模式
    project      string          // 项目名称
    dbInit       string          // DB初始化, 用于开发快速重置数据库
    dataPath     string          // 用于制定数据存储位置
    then         int64           // 系统启动时间
    logLevel     logrus.Level    // 日志等级
    gormLogLevel logger.LogLevel // gorm日志等级
}

type Environment interface {
    Value() string
    Super() bool
    MscAddr() string
    Mode() string
    DBInit() string
    Then() int64
    Project() string
    DataPath() string
    setSupper(is bool)
    setDBInit(table string)
    setMscAddr(addr string)
    setDataPath(path string)
    Level() logrus.Level
    GormLogLevel() logger.LogLevel
}

// gin框架运行模式
func (env *environment) Mode() string {
    return env.ginMode
}

// APP运行模式
func (env *environment) Value() string {
    return env.value
}

// 是否为超级用户
func (env *environment) Super() bool {
    return env.super
}

// 系统启动时间
func (env *environment) Then() int64 {
    return env.then
}

// 用于重置DB某张表或所有表重置, 主要用于开发时方便快速更新表结构
func (env *environment) DBInit() string {
    return env.dbInit
}

// 数据存放路径
func (env *environment) DataPath() string {
    return env.dataPath
}

// 日志等级
func (env *environment) Level() logrus.Level {
    return env.logLevel
}

// gorm日志等级
func (env *environment) GormLogLevel() logger.LogLevel {
    return env.gormLogLevel
}

// 服务监听地址
func (env *environment) MscAddr() string {
    return env.mscAddr
}

// 当前工程名称
func (env *environment) Project() string {
    return env.project
}

// 设置监听地址
func (env *environment) setMscAddr(addr string) {
    env.mscAddr = addr
}

// 设置是否为超级用户
func (env *environment) setSupper(is bool) {
    env.super = is
}

// 设置DB重置的表
func (env *environment) setDBInit(table string) {
    env.dbInit = table
}

// 设置数据路径
func (env *environment) setDataPath(path string) {
    env.dataPath = path
}

func Init() {
    env := flag.String("env", "", "请输入运行环境:\n dev:开发环境\n fat:测试环境\n pro:正式环境\n")
    mscAddr := flag.String("mscaddr", ENV_MSC_DEPAULT_PORT, "请输入端口及地址") // 管理端口, 查看当前系统状态
    super := flag.String("super", "", "")                               // 超级管理员地址
    dbInit := flag.String("dbinit", "", "")                             // 重置数据库
    dataPath := flag.String("data_path", "", "")                        // 数据存放路径
    flag.Parse()
    mscAddrStr := strings.ToLower(strings.TrimSpace(*mscAddr))
    superStr := strings.TrimSpace(*super)
    dbInitStr := strings.TrimSpace(*dbInit)
    dataPathStr := strings.TrimSpace(*dataPath)
    switch strings.ToLower(strings.TrimSpace(*env)) {
    case ENV_MODE_DEV:
        active = dev
    case ENV_MODE_FAT:
        active = fat
    case ENV_MODE_PRO:
        active = pro
    default:
        active = pro
    }
    active.setMscAddr(mscAddrStr)
    if superStr == ENV_APP_NAME+"1qwer$#@" {
        active.setSupper(true)
    }
    if dbInitStr != "" {
        active.setDBInit(dbInitStr)
    }
    if dataPathStr != "" {
        active.setDataPath(dataPathStr)
    }

}

func Active() Environment {
    return active
}
