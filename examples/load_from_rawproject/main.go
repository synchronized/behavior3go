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
)

func main() {
	projectConfig, ok := config.LoadRawProjectCfg("example.b3")
	if !ok {
		fmt.Println("LoadRawProjectCfg err")
		return
	}

	//自定义节点注册
	maps := b3.NewRegisterStructMaps()
	maps.Register("Log", new(LogTest))

	var firstTree = CreateBevTreeFromConfig(maps, projectConfig.Data.Trees...)
	firstTree.Print()

	//输入板
	board := NewBlackboard()
	//循环每一帧
	for i := 0; i < 5; i++ {
		firstTree.Tick(i, board)
	}
}
