package core

import (
	"fmt"

	b3 "github.com/magicsea/behavior3go"
	"github.com/magicsea/behavior3go/config"
)

/**
 * The BehaviorTree class, as the name implies, represents the Behavior Tree
 * structure.
 *
 * There are two ways to construct a Behavior Tree: by manually setting the
 * root node, or by loading it from a data structure (which can be loaded
 * from a JSON). Both methods are shown in the examples below and better
 * explained in the user guide.
 *
 * The tick method must be called periodically, in order to send the tick
 * signal to all nodes in the tree, starting from the root. The method
 * `BehaviorTree.tick` receives a target object and a blackboard as
 * parameters. The target object can be anything: a game agent, a system, a
 * DOM object, etc. This target is not used by any piece of Behavior3JS,
 * i.e., the target object will only be used by custom nodes.
 *
 * The blackboard is obligatory and must be an instance of `Blackboard`. This
 * requirement is necessary due to the fact that neither `BehaviorTree` or
 * any node will store the execution variables in its own object (e.g., the
 * BT does not store the target, information about opened nodes or number of
 * times the tree was called). But because of this, you only need a single
 * tree instance to control multiple (maybe hundreds) objects.
 *
 * Manual construction of a Behavior Tree
 * --------------------------------------
 *
 *     var tree = new b3.BehaviorTree();
 *
 *     tree.root = new b3.Sequence({children:[
 *       new b3.Priority({children:[
 *         new MyCustomNode(),
 *         new MyCustomNode()
 *       ]}),
 *       ...
 *     ]});
 *
 *
 * Loading a Behavior Tree from data structure
 * -------------------------------------------
 *
 *     var tree = new b3.BehaviorTree();
 *
 *     tree.load({
 *       'title'       : 'Behavior Tree title'
 *       'description' : 'My description'
 *       'root'        : 'node-id-1'
 *       'nodes'       : {
 *         'node-id-1' : {
 *           'name'        : 'Priority', // this is the node type
 *           'title'       : 'Root Node',
 *           'description' : 'Description',
 *           'children'    : ['node-id-2', 'node-id-3'],
 *         },
 *         ...
 *       }
 *     })
 *
 *
 * @module b3
 * @class BehaviorTree
**/
type BehaviorTree struct {

	/**
	 * The tree id, must be unique. By default, created with `b3.createUUID`.
	 * @property {String} id
	 * @readOnly
	**/
	id string

	/**
	 * The tree title.
	 * @property {String} title
	 * @readonly
	**/
	title string

	/**
	 * Description of the tree.
	 * @property {String} description
	 * @readonly
	**/
	description string

	/**
	 * A dictionary with (key-value) properties. Useful to define custom
	 * variables in the visual editor.
	 *
	 * @property {Object} properties
	 * @readonly
	**/
	properties map[string]interface{}

	/**
	 * The reference to the root node. Must be an instance of `b3.BaseNode`.
	 * @property {BaseNode} root
	**/
	root IBaseNode

	/**
	 * The reference to the debug instance.
	 * @property {Object} debug
	**/
	debug interface{}

	dumpInfo *config.BTTreeCfg
}

func NewBeTree() *BehaviorTree {
	tree := &BehaviorTree{}
	tree.Initialize()
	return tree
}

/**
 * Initialization method.
 * @method Initialize
 * @construCtor
**/
func (this *BehaviorTree) Initialize() {
	this.id = b3.CreateUUID()
	this.title = "The behavior tree"
	this.description = "Default description"
	this.properties = make(map[string]interface{})
	this.root = nil
	this.debug = nil
}

func (this *BehaviorTree) GetID() string {
	return this.id
}

func (this *BehaviorTree) GetTitile() string {
	return this.title
}

func (this *BehaviorTree) SetDebug(debug interface{}) {
	this.debug = debug
}

func (this *BehaviorTree) GetRoot() IBaseNode {
	return this.root
}

/**
 * This method dump the current BT into a data structure.
 *
 * Note: This method does not record the current node parameters. Thus,
 * it may not be compatible with load for now.
 *
 * @method dump
 * @return {Object} A data object representing this tree.
**/
func (this *BehaviorTree) dump() *config.BTTreeCfg {
	return this.dumpInfo
}

/**
 * Propagates the tick signal through the tree, starting from the root.
 *
 * This method receives a target object of any type (Object, Array,
 * DOMElement, whatever) and a `Blackboard` instance. The target object has
 * no use at all for all Behavior3JS components, but surely is important
 * for custom nodes. The blackboard instance is used by the tree and nodes
 * to store execution variables (e.g., last node running) and is obligatory
 * to be a `Blackboard` instance (or an object with the same interface).
 *
 * Internally, this method creates a Tick object, which will store the
 * target and the blackboard objects.
 *
 * Note: BehaviorTree stores a list of open nodes from last tick, if these
 * nodes weren't called after the current tick, this method will close them
 * automatically.
 *
 * @method tick
 * @param {Object} target A target object.
 * @param {Blackboard} blackboard An instance of blackboard object.
 * @return {Constant} The tick signal state.
**/
func (this *BehaviorTree) Tick(target interface{}, blackboard *Blackboard) b3.Status {
	if blackboard == nil {
		panic("The blackboard parameter is obligatory and must be an instance of b3.Blackboard")
	}

	/* CREATE A TICK OBJECT */
	var tick = NewTick()
	tick.debug = this.debug
	tick.target = target
	tick.Blackboard = blackboard
	tick.tree = this

	/* TICK NODE */
	var state = this.root._execute(tick)

	/* CLOSE NODES FROM LAST TICK, IF NEEDED */
	var lastOpenNodes = blackboard._getTreeData(this.id).OpenNodes
	var currOpenNodes []IBaseNode
	currOpenNodes = append(currOpenNodes, tick._openNodes...)

	l := len(lastOpenNodes)
	if l == len(currOpenNodes) {
		if l == 0 || lastOpenNodes[l-1] == currOpenNodes[l-1] {
			return state
		}
	}

	// does not close if it is still open in this tick
	var start = 0
	for i := 0; i < b3.MinInt(len(lastOpenNodes), len(currOpenNodes)); i++ {
		start = i
		if lastOpenNodes[i] != currOpenNodes[i] {
			break
		}
	}

	// close the nodes
	for i := len(lastOpenNodes) - 1; i >= start; i-- {
		lastOpenNodes[i]._close(tick)
	}

	/* POPULATE BLACKBOARD */
	blackboard._getTreeData(this.id).OpenNodes = currOpenNodes
	blackboard.SetTree("nodeCount", tick._nodeCount, this.id)

	return state
}

func (this *BehaviorTree) Print() {
	printNode(this.root, 0)
}

func printNode(root IBaseNode, blk int) {

	//fmt.Println("new node:", root.Name, " children:", len(root.Children), " child:", root.Child)
	for i := 0; i < blk; i++ {
		fmt.Print(" ") //缩进
	}

	//fmt.Println("|—<", root.Name, ">") //打印"|—<id>"形式
	fmt.Println(fmt.Sprintf("|—(%v)[%v]%v", root.GetDepth(), root.GetCategory(), root.GetTitle()))

	switch root.GetCategory() {

	case b3.COMPOSITE:
		comp := root.(IComposite)
		for i := 0; i < comp.GetChildCount(); i++ {
			printNode(comp.GetChild(i), blk+3)
		}

	case b3.DECORATOR:
		dec := root.(IDecorator)
		if dec.GetChild() != nil {
			//fmt.Print("=>")
			printNode(dec.GetChild(), blk+3)
		}

	case b3.TREE:
		subt := root.(ISubTree)
		if subt.GetChild() != nil {
			//fmt.Print("=>")
			printNode(subt.GetChild(), blk+3)
		}
	}

}
