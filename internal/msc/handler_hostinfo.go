/*
 * @Author: reel
 * @Date: 2023-06-10 20:16:56
 * @LastEditors: reel
 * @LastEditTime: 2023-09-12 06:20:32
 * @Description: 请填写简介
 */
package msc

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fbs-io/core/pkg/env"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

const (
	mb        = uint64(1024 * 1024) // MB
	Minute    = 60
	Hour      = 60 * Minute
	Day       = 24 * Hour
	Week      = 7 * Day
	Month     = 30 * Day
	Year      = 12 * Month
	seconds   = 1000 * 1000 * 1000
	miSeconds = 1000 * 1000
)

func (h *handler) hostInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		hosts, _ := host.Info()
		cpus, _ := cpu.Info()
		osday1, osday2 := timeSincePro(int64(hosts.BootTime))
		m := new(runtime.MemStats)
		appday1, appday2 := timeSincePro(env.Active().Then())
		v, _ := mem.VirtualMemory()
		lastgc := int64(m.LastGC)
		if int64(m.LastGC) == 0 {
			lastgc = env.Active().Then()
		}
		lastgc = (time.Now().UnixNano()/seconds - lastgc)
		data := map[string]interface{}{
			"os":           hosts.Platform + hosts.PlatformVersion,                                         // 系统信息
			"cpu":          fmt.Sprintf("%s (%d cores) X %d", cpus[0].ModelName, cpus[0].Cores, len(cpus)), // cpu信息
			"mem":          v.Total / mb,                                                                   // 内存总量MB
			"host":         hosts.Hostname,                                                                 // 系统名称
			"pause":        fmt.Sprintf("%.3fms", float64(m.PauseNs[(m.NumGC+255)%256])/miSeconds),         // 上次 GC 时间 单位: 毫秒
			"numgc":        m.NumGC,                                                                        // GC次数
			"nextgc":       m.NextGC / mb,                                                                  // 下次 GC
			"lastgc":       lastgc,                                                                         // 上次 GC 时间, 单位:秒
			"version":      env.Active().Version(),                                                         // 内核版本
			"osrunday":     fmt.Sprintf("%s-%s", osday1, osday2),                                           // 系统运行时间
			"goversion":    runtime.Version(),                                                              // golang版本
			"apprunday":    fmt.Sprintf("%s-%s", appday1, appday2),                                         // app运行时间
			"pausetotal":   fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/seconds),                          // GC 总时间 单位:秒
			"appversion":   env.Active().AppVersion(),                                                      // app版本
			"numgoroutine": runtime.NumGoroutine(),                                                         // 协程个数
		}

		ctx.JSON(200, errno.ERRNO_OK.ToMapWithData(data))
	}
}

func (h *handler) sysInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := map[string]interface{}{}
		for k, v := range h.sysinfo {
			data[k] = v
		}

		data["processinfo"] = h.procinfos
		ctx.JSON(200, errno.ERRNO_OK.ToMapWithData(data))
	}
}

type processInfo struct {
	PID        int32   `json:"pid"`
	PName      string  `json:"pname"`
	CpuPercent float64 `json:"cpupercent"`
	MemInfo    uint64  `json:"meminfo"`
	IO         int32   `json:"io"`
}

type processInfos []processInfo

func (p processInfos) Len() int           { return len(p) }
func (p processInfos) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p processInfos) Less(i, j int) bool { return p[i].CpuPercent > p[j].CpuPercent }

// 获取进程相关信息
func (m *handler) getProcessInfo() processInfos {
	var procinfos = make(processInfos, 0, 100)
	processes, _ := process.Processes()
	for _, p := range processes {
		cpuPercent, _ := p.CPUPercent()
		if cpuPercent < 0.01 {
			continue
		}

		var procinfo = processInfo{}
		mem, _ := p.MemoryInfo()
		procinfo.PID = p.Pid
		procinfo.IO, _ = p.IOnice()
		procinfo.PName, _ = p.Name()
		procinfo.MemInfo = mem.RSS / mb
		procinfo.CpuPercent = cpuPercent
		procinfos = append(procinfos, procinfo)
	}

	sort.Sort(procinfos)
	return procinfos[:20]
}

// 获取系统资源占用信息
func (m *handler) getSysInfo() map[string]interface{} {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	// cpu,内存
	cpuPercents, _ := cpu.Percent(time.Millisecond*200, false)
	virMem, _ := mem.VirtualMemory()
	appmem := memStats.Alloc / mb

	// 网络相关
	netios, _ := net.IOCounters(true)
	var sends, recv uint64
	for _, netio := range netios {
		if netio.Name == "lo0" {
			continue
		}
		sends += netio.BytesSent / mb
		recv += netio.BytesRecv / mb
	}

	// 磁盘相关
	disks, _ := disk.Usage("./")
	diskios, _ := disk.IOCounters()
	var read, write uint64
	for _, disk := range diskios {

		read += disk.ReadBytes / mb
		write += disk.WriteBytes / mb
	}

	// cpu使用率处理, 如果多核cpu取平均值
	var cpup float64
	for _, cpupercent := range cpuPercents {
		cpup += cpupercent
	}
	cpup /= float64(len(cpuPercents))

	data := map[string]interface{}{
		"cpu":     fmt.Sprintf("%.1f", cpup),
		"memp":    fmt.Sprintf("%.1f", virMem.UsedPercent),
		"appmemp": fmt.Sprintf("%.1f", float64(appmem)/float64(virMem.Total/mb)*100),
		"appmem":  fmt.Sprintf("%.1f", float64(appmem)),
		"disk":    fmt.Sprintf("%.1f", disks.UsedPercent),
		"netio": map[string]interface{}{
			"send": sends,
			"recv": recv,
		},
		"diskio": map[string]interface{}{
			"read":  read,
			"write": write,
		},
	}
	m.sysinfo = data
	return data

}

func computeTimeDiff(diff int64) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = "当前"
	case diff < 2:
		diff = 0
		diffStr = "1 秒"
	case diff < 1*Minute:
		diffStr = fmt.Sprintf("%d 秒", diff)
		diff = 0

	case diff < 2*Minute:
		diff -= 1 * Minute
		diffStr = "1 分钟"
	case diff < 1*Hour:
		diffStr = fmt.Sprintf("%d 分钟", diff/Minute)
		diff -= diff / Minute * Minute

	case diff < 2*Hour:
		diff -= 1 * Hour
		diffStr = "1 小时"
	case diff < 1*Day:
		diffStr = fmt.Sprintf("%d 小时", diff/Hour)
		diff -= diff / Hour * Hour

	case diff < 2*Day:
		diff -= 1 * Day
		diffStr = "1 天"
	case diff < 1*Week:
		diffStr = fmt.Sprintf("%d 天", diff/Day)
		diff -= diff / Day * Day

	case diff < 2*Week:
		diff -= 1 * Week
		diffStr = "1 周"
	case diff < 1*Month:
		diffStr = fmt.Sprintf("%d 周", diff/Week)
		diff -= diff / Week * Week

	case diff < 2*Month:
		diff -= 1 * Month
		diffStr = "1 月"
	case diff < 1*Year:
		diffStr = fmt.Sprintf("%d 月", diff/Month)
		diff -= diff / Month * Month

	case diff < 2*Year:
		diff -= 1 * Year
		diffStr = "1 年"
	default:
		diffStr = fmt.Sprintf("%d 年", diff/Year)
		diff = 0
	}
	return diff, diffStr
}

func timeSincePro(then int64) (string, string) {
	now := time.Now()
	diff := now.Unix() - then

	if diff < 0 {
		return "future", ""
	}

	var diffStr string
	var timeStrs = make([]string, 0, 10)
	for {
		if diff == 0 {
			break
		}

		diff, diffStr = computeTimeDiff(diff)
		timeStrs = append(timeStrs, diffStr)
	}
	if len(timeStrs) < 2 {
		return strings.Join(timeStrs, ", "), ""
	}

	return timeStrs[0], strings.Join(timeStrs[1:], ", ")
}
