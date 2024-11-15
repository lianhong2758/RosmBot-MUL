package rosm

import (
	"time"
	"unsafe"

	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/rate"
)

// defaultLimiterManager 默认限速器管理
//
//	每 10s 5次触发
var defaultLimiterManager = rate.NewManager[string](time.Second*10, 5)

type fakeLM struct {
	limiters unsafe.Pointer
	interval time.Duration
	burst    int
}

// SetDefaultLimiterManagerParam 设置默认限速器参数
//
//	每 interval 时间 burst 次触发
func SetDefaultLimiterManagerParam(interval time.Duration, burst int) {
	f := (*fakeLM)(unsafe.Pointer(defaultLimiterManager))
	f.interval = interval
	f.burst = burst
}

// LimitByUser 默认限速器 每 10s 5次触发
//
//	按用户限制
func LimitByUser(ctx *Ctx) *rate.Limiter {
	return defaultLimiterManager.Load(ctx.Being.User.ID)
}

// LimitByGroup 默认限速器 每 10s 5次触发
//
//	按群号限制
func LimitByGroup(ctx *Ctx) *rate.Limiter {
	return defaultLimiterManager.Load(tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID))
}

// LimiterManager 自定义限速器管理
type LimiterManager struct {
	m *rate.LimiterManager[string]
}

// NewLimiterManager 新限速器管理
func NewLimiterManager(interval time.Duration, burst int) (m LimiterManager) {
	m.m = rate.NewManager[string](interval, burst)
	return
}

// LimitByUser 自定义限速器
//
//	按用户限制
func (m LimiterManager) LimitByUser(ctx *Ctx) *rate.Limiter {
	return m.m.Load(ctx.Being.User.ID)
}

// LimitByGroup 自定义限速器
//
//	按群号限制
func (m LimiterManager) LimitByGroup(ctx *Ctx) *rate.Limiter {
	return m.m.Load(tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID))
}
