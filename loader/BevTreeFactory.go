package loader

import (
	_ "fmt"
	_ "reflect"

	b3 "github.com/magicsea/behavior3go"
	. "github.com/magicsea/behavior3go/actions"
	. "github.com/magicsea/behavior3go/composites"
	. "github.com/magicsea/behavior3go/config"
	. "github.com/magicsea/behavior3go/core"
	. "github.com/magicsea/behavior3go/decorators"
)

func CreateBaseStructMaps() *b3.RegisterStructMaps {
	st := b3.NewRegisterStructMaps()
	//actions
	st.Register("Error", &Error{})
	st.Register("Failer", &Failer{})
	st.Register("Runner", &Runner{})
	st.Register("Succeeder", &Succeeder{})
	st.Register("Wait", &Wait{})
	st.Register("Log", &Log{})
	//composites
	st.Register("MemPriority", &MemPriority{})
	st.Register("MemSequence", &MemSequence{})
	st.Register("Priority", &Priority{})
	st.Register("Sequence", &Sequence{})

	//decorators
	st.Register("Inverter", &Inverter{})
	st.Register("Limiter", &Limiter{})
	st.Register("MaxTime", &MaxTime{})
	st.Register("Repeater", &Repeater{})
	st.Register("RepeatUntilFailure", &RepeatUntilFailure{})
	st.Register("RepeatUntilSuccess", &RepeatUntilSuccess{})
	return st
}

func CreateBevTreeFromConfig(extMaps *b3.RegisterStructMaps, configs ...*BTTreeCfg) *BehaviorTree {
	var treeConfigMap = make(map[string]*BTTreeCfg)
	for _, v := range configs {
		treeConfigMap[v.ID] = v
	}
	var firstConfigId = configs[0].ID
	var baseMaps = CreateBaseStructMaps()
	var context = NewContext(baseMaps, extMaps, func(id string) *BTTreeCfg {
		return treeConfigMap[id]
	})

	return CreateBehaviorTree(context, firstConfigId)
}
