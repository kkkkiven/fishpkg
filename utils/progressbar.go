package utils

import (
	"fmt"
	"time"
)

// 进度条定义
type ProgressBar struct {
	Name string // 名称
	Total Counter // 总任务数
	Current *Counter // 当前已完成数量
	startTime time.Time // 任务开始时间
	exit chan struct{} // 是否退出进度条
	alreadyOutput bool
}

// 完成指定数量的任务
func (pgb *ProgressBar) Done(i ...int64) {
	if len(i) == 1 {
		pgb.Current.Add(i[0])
	} else {
		pgb.Current.Add(1)
	}
}

// 退出进度条
func (pgb *ProgressBar) Exit() {
	pgb.exit <- struct{}{}
	if pgb.alreadyOutput {
		fmt.Println("")
	}
}

// 启动并显示进度条
func (pgb *ProgressBar) Run() {
	for {
		select {
		case <-pgb.exit:
			return
		case <-time.After(time.Second):
			pgb.print()
		}
	}
}

// 输出进度
func (pgb *ProgressBar) print() {
	var percent float64 = 0
	var per, surplus int64

	cur := pgb.Current.Get()
	if pgb.Total.Get() > 0 {
		percent = float64(cur) / float64(pgb.Total.Get())
	}
	useTime := int64(time.Now().Sub(pgb.startTime) / time.Second)
	if useTime > 0 {
		per = cur / useTime
	}
	if per > 0 {
		surplus = (pgb.Total.Get() - cur) / per
	}

	fmt.Printf("\r\x1b[36m%s[%.2f%%]  %d/%d  %d/ps  used:%s  surplus:%s\x1b[0m", pgb.Name, percent * 100, cur, pgb.Total.Get(), per, HumanTimeSecond(useTime), HumanTimeSecond(surplus))

	if !pgb.alreadyOutput {
		pgb.alreadyOutput = true
	}
}

// 创建一个进度条
func NewPGBar(total int64, name ...string) *ProgressBar {
	var n string
	if len(name) == 1 && name[0] != "" {
		n = name[0] + " "
	}
	return &ProgressBar{
		Name: n,
		Total: Counter{v : total},
		Current: &Counter{},
		startTime: time.Now(),
		exit: make(chan struct{}, 1),
	}
}

// 生成一个新的进度条并运行
func NewAndRunPGBar(total int64, name ...string) *ProgressBar {
	pgb := NewPGBar(total, name...)
	go pgb.Run()
	return pgb
}