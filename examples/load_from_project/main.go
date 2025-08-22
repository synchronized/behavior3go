/*
从导出的工程文件加载
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
	projectConfig, ok := config.LoadProjectCfg("project.json")
	if !ok {
		fmt.Println("LoadTreeCfg err")
		return
	}

	//自定义节点注册
	maps := b3.NewRegisterStructMaps()
	maps.Register("Log", new(LogTest))

	var firstTree = CreateBevTreeFromConfig(maps, projectConfig.Trees...)
	firstTree.Print()

	//输入板
	board := NewBlackboard()
	//循环每一帧
	for i := 0; i < 5; i++ {
		firstTree.Tick(i, board)
	}
}
