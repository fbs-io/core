/*
 * @Description: 用于开启一个 http 服务
 * @Params: 可变参数, 可以指定端口号, 服务名称,
 * @Author: LenLee
 * @Date: 2022-06-24 21:59:45
 * @LastEditTime: 2023-06-14 21:36:31
 * @LastEditors: reel
 * @FilePath:
 */
package mux

import (
    "context"
    "fmt"

    "encoding/json"
    "strings"

    "github.com/fbs-io/core/pkg/consts"
    "github.com/fbs-io/core/pkg/errorx"
    "github.com/fbs-io/core/pkg/httpx"

    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

var _ Mux = (*httpServer)(nil)

type Mux interface {
    Name() string
    Status() int8
    Start() error
    Stop() error
    Engine() *gin.Engine
    SetAddr(addr string)
}

type httpServer struct {
    name      string
    status    int8
    client    httpx.Client
    server    *http.Server
    engine    *gin.Engine
    timeout   time.Duration
    statusUrl string
    opt       *opts
}

func New(optF ...OptFunc) (srv Mux, err error) {
    var opt = &opts{
        addr:          ":6018", // 后台管理地址
        maxHeaderSize: 2 << 32,
        maxReadTime:   60 * time.Second,
        maxWriteTIme:  60 * time.Second,
        timeout:       time.Second * 30,
    }
    for _, f := range optF {
        f(opt)
    }
    engine := gin.New()
    h := &httpServer{
        name:      opt.name,
        client:    httpx.New(),
        timeout:   opt.timeout,
        engine:    engine,
        opt:       opt,
        statusUrl: "/srv_status",
    }
    engine.GET(h.statusUrl, func(ctx *gin.Context) { ctx.JSON(200, 1) })

    return h, nil
}

func (h *httpServer) PowerP()      {}
func (h *httpServer) Name() string { return h.name }
func (h *httpServer) Status() int8 {
    // 0表示服务关闭, 1表示正常, -1表示服务不可用
    if h.status == consts.SERVER_IS_NULL {
        return consts.SERVER_IS_NULL
    }
    if h.server == nil {
        return consts.SERVER_IS_NULL
    }
    body, err := h.client.Do("http://127.0.0.1:" + strings.Split(h.server.Addr, ":")[1] + h.statusUrl)
    var status int8
    json.Unmarshal(body, &status)
    if err == nil && status == consts.SERVER_IS_RUN {
        return consts.SERVER_IS_RUN
    }
    return consts.SERVER_IS_DOWN
}
func (h *httpServer) Stop() error {
    ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
    defer cancel()
    return h.server.Shutdown(ctx)
}
func (h *httpServer) Start() error {
    var err error
    go func() {
        h.server = &http.Server{
            Addr:           h.opt.addr,
            Handler:        h.engine,
            MaxHeaderBytes: h.opt.maxHeaderSize,
            ReadTimeout:    h.opt.maxReadTime,
            WriteTimeout:   h.opt.maxWriteTIme,
        }
        h.server.SetKeepAlivesEnabled(true)
        if err = h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            return
        }
    }()
    return errorx.Wrap(err, fmt.Sprintf("%s发生错误", h.name))
}

func (h *httpServer) Engine() *gin.Engine {
    return h.engine
}
func (h *httpServer) SetAddr(addr string) {
    h.opt.addr = addr
}
