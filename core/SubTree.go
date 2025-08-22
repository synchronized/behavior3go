package core

import (
	b3 "github.com/magicsea/behavior3go"
	. "github.com/magicsea/behavior3go/config"
)

var _ ISubTree = (*SubTree)(nil)

type ISubTree interface {
	IBaseNode
	SetTree(tree *BehaviorTree)
	GetTree() *BehaviorTree
	GetChild() IBaseNode
}

//子树，通过Name关联树ID查找
type SubTree struct {
	Action
	tree *BehaviorTree
}

func (this *SubTree) Ctor() {
	this.category = b3.TREE
}

func (this *SubTree) Initialize(setting *BTNodeCfg) {
	this.Action.Initialize(setting)
}

/**
 *执行子树
 *使用sTree.Tick(tar, tick.Blackboard)的方法会导致每个树有自己的tick。
 *如果子树包含running状态，同时复用了子树会导致歧义。
 *改为只使用一个树，一个tick上下文。
**/
func (this *SubTree) OnTick(tick *Tick) b3.Status {

	//使用子树，必须先SetSubTreeLoadFunc
	//子树可能没有加载上来，所以要延迟加载执行
	if this.tree == nil {
		return b3.ERROR
	}

	if tick.GetTarget() == nil {
		panic("SubTree tick.GetTarget() nil !")
	}

	//tar := tick.GetTarget()
	//return sTree.Tick(tar, tick.Blackboard)

	tick.pushSubtreeNode(this)
	ret := this.tree.GetRoot().Execute(tick)
	tick.popSubtreeNode()
	return ret
}

func (this *SubTree) SetDepth(depth int) {
	this.BaseNode.SetDepth(depth)
	{
		var child = this.GetChild()
		child.SetDepth(depth + 1)
	}
}

func (this *SubTree) String() string {
	return "SBT_" + this.GetTitle()
}

func (this *SubTree) SetTree(tree *BehaviorTree) {
	this.tree = tree
}
func (this *SubTree) GetTree() *BehaviorTree {
	return this.tree
}

func (this *SubTree) GetChild() IBaseNode {
	return this.GetTree().GetRoot()
}
