package core

import (
	b3 "github.com/magicsea/behavior3go"
	. "github.com/magicsea/behavior3go/config"
)

type IDecorator interface {
	IBaseNode
	SetChild(child IBaseNode)
	GetChild() IBaseNode
}

type Decorator struct {
	BaseNode
	BaseWorker
	child IBaseNode
}

func (this *Decorator) Ctor() {

	this.category = b3.DECORATOR
}

/**
 * Initialization method.
 *
 * @method Initialize
 * @construCtor
**/
func (this *Decorator) Initialize(params *BTNodeCfg) {
	this.BaseNode.Initialize(params)
	//this.BaseNode.IBaseWorker = this
}

func (this *Decorator) SetDepth(depth int) {
	this.BaseNode.SetDepth(depth)
	if this.child != nil {
		var child = this.GetChild()
		child.SetDepth(depth + 1)
	}
}

//GetChild
func (this *Decorator) GetChild() IBaseNode {
	return this.child
}

func (this *Decorator) SetChild(child IBaseNode) {
	this.child = child
}
