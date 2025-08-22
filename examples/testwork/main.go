/*
从原生工程文件加载
*/
package main

import (
	"fmt"
	b3 "github.com/magicsea/behavior3go"
	"github.com/magicsea/behavior3go/config"
	. "github.com/magicsea/behavior3go/core"
	. "github.com/magicsea/behavior3go/examples/share"
	. "github.com/magicsea/behavior3go/loader"
	"sync"
	"time"
)

// 所有的树管理
var mapTreesByID = sync.Map{}
var maps = b3.NewRegisterStructMaps()

func init() {
	//自定义节点注册
	maps.Register("Log", new(LogTest))
	maps.Register("SetValue", new(SetValue))
	maps.Register("IsValue", new(IsValue))
}

func main() {
	projectConfig, ok := config.LoadRawProjectCfg("testwork.b3")
	if !ok {
		fmt.Println("LoadRawProjectCfg err")
		return
	}

	var firstTree = CreateBevTreeFromConfig(maps, projectConfig.Data.Trees...)
	firstTree.Print()

	time.Sleep(time.Second)
	//输入板
	board := NewBlackboard()
	//循环每一帧
	for i := 0; i < 40; i++ {
		firstTree.Tick(i, board)
		time.Sleep(time.Millisecond * 100)
	}
}
