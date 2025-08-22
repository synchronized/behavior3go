package core

import (
	"fmt"

	b3 "github.com/magicsea/behavior3go"
	"github.com/magicsea/behavior3go/config"
)

type GetBehaviorTreeConfig func(id string) *config.BTTreeCfg

/**
 创建行为树上下文
**/
type Context struct {
	Maps    *b3.RegisterStructMaps
	ExtMaps *b3.RegisterStructMaps

	GetConfig GetBehaviorTreeConfig
}

func NewContext(maps *b3.RegisterStructMaps, extMaps *b3.RegisterStructMaps, getConfig GetBehaviorTreeConfig) *Context {
	ctx := &Context{
		Maps:      maps,
		ExtMaps:   extMaps,
		GetConfig: getConfig,
	}
	return ctx
}

/**
 * 创建节点
**/
func (this *Context) NewNode(nodeName string) IBaseNode {
	if this.ExtMaps != nil && this.ExtMaps.CheckElem(nodeName) {
		// Look for the name in custom nodes
		if tnode, err := this.ExtMaps.New(nodeName); err == nil {
			return tnode.(IBaseNode)
		}
	}
	if tnode, err2 := this.Maps.New(nodeName); err2 == nil {
		return tnode.(IBaseNode)
	}
	return nil
}

func CreateBehaviorTree(ctx *Context, id string) *BehaviorTree {
	return doCreateBehaviorTree(ctx, id, 0)
}

func doCreateBehaviorTree(ctx *Context, id string, depth int) *BehaviorTree {
	if depth > 1000 {
		panic(fmt.Sprintf("CreateBehaviorTree: 递归层数过多可能有地方死循环了 id:%v, depth:%v", id, depth))
	}

	var data = ctx.GetConfig(id)

	var this = NewBeTree()
	this.title = data.Title             //|| this.title;
	this.description = data.Description // || this.description;
	this.properties = data.Properties   // || this.properties;
	this.dumpInfo = data

	nodes := make(map[string]IBaseNode)

	// Create the node list (without connection between them)

	var newBaseNode = func(spec *config.BTNodeCfg) IBaseNode {
		if spec.Category == "tree" {
			return new(SubTree)
		} else {
			if newNode := ctx.NewNode(spec.Name); newNode != nil {
				return newNode
			}
		}
		return nil
	}

	// Connect the nodes
	var openNodes = []*config.BTNodeCfg{data.Nodes[data.Root]}
	var parentNodeMap = make(map[string]IBaseNode)
	for len(openNodes) > 0 {
		var currNodeSpec = openNodes[0]
		var spec = currNodeSpec
		openNodes = append(openNodes[:0], openNodes[1:]...)

		var node = newBaseNode(spec)

		if node == nil {
			// Invalid node name
			panic("CreateBehaviorTree: Invalid node name:" + spec.Name + ",title:" + spec.Title)
		}

		node.Ctor()
		node.Initialize(spec)
		node.SetBaseNodeWorker(node.(IBaseWorker))
		nodes[spec.Id] = node

		if node.GetCategory() == b3.TREE {
			var subTreeNode = node.(ISubTree)
			var subTree = doCreateBehaviorTree(ctx, subTreeNode.GetName(), depth+1)
			if subTree == nil {
				panic("CreateBehaviorTree: Invalid tree id:" + subTreeNode.GetName() + ",title:" + spec.Title)
			}
			subTreeNode.SetTree(subTree)
		}

		var parentNode = parentNodeMap[spec.Id]
		if parentNode != nil {
			switch parentNode.GetCategory() {
			case b3.COMPOSITE:
				compNode := parentNode.(IComposite)
				compNode.AddChild(node)
			case b3.DECORATOR:
				decNode := parentNode.(IDecorator)
				decNode.SetChild(node)
			}
		}

		switch node.GetCategory() {
		case b3.COMPOSITE:
			for i := 0; i < len(spec.Children); i++ {
				var cid = spec.Children[i]
				openNodes = append(openNodes, data.Nodes[cid])
				parentNodeMap[cid] = node
			}
		case b3.DECORATOR:
			if spec.Child != "" {
				var cid = spec.Child
				openNodes = append(openNodes, data.Nodes[cid])
				parentNodeMap[cid] = node
			}
		}
	}

	this.root = nodes[data.Root]

	this.root.SetDepth(1)
	return this
}
