/*
 * @Author: reel
 * @Date: 2023-05-16 21:37:52
 * @LastEditors: reel
 * @LastEditTime: 2023-06-11 15:00:25
 * @Description: 设置项目唯一标识, 用于记录项目是否启动
 */
package pem

import (
    "io"
    "os"
    "path"
    "runtime"

    "github.com/fbs-io/core/pkg/env"
    "github.com/fbs-io/core/pkg/errorx"
    "github.com/fbs-io/core/pkg/filex"
)

func GetPems() (pems string, err error) {

    filePath, err := GetPemPath()
    if err != nil {
        return
    }
    file, err := os.OpenFile(path.Join(filePath, "pems"), os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        return "", err
    }
    defer file.Close()

    pemsByte, err := io.ReadAll(file)
    if err != nil {
        return "", err
    }

    if len(pemsByte) > 0 {
        return string(pemsByte), nil
    }
    return "", errorx.New("pems为空")
}

func UpdatePems(data string) (err error) {
    filePath, err := GetPemPath()
    if err != nil {
        return
    }
    file, err := os.OpenFile(path.Join(filePath, "pems"), os.O_RDWR|os.O_TRUNC, 0666)
    if err != nil {
        return err
    }
    defer file.Close()
    _, err = file.Write([]byte(data))
    return
}

func GetPemPath() (filePath string, err error) {
    userHome, err := os.UserHomeDir()
    switch runtime.GOOS {
    case "windows":
        if err != nil {
            return "", err
        }

        filePath = path.Join(userHome, "AppData", env.Active().AppName(), "pems")
        err = filex.CreatDir(filePath)
        if err != nil {
            return "", err
        }
    default:
        filePath = path.Join(userHome, "."+env.Active().AppName(), "pems")
        err := filex.CreatDir(filePath)
        if err != nil {
            return "", err
        }
    }
    return
}
