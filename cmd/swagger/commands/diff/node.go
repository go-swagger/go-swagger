package diff

// Node defines a diff position in a spec
type Node struct {
	Field     string
	TypeName  string
	IsArray   bool
	ChildNode *Node
}

// String std string render method
func (n *Node) String() string {
	name := n.Field
	if n.IsArray {
		name = "array["+n.TypeName+"]"
	}
	
	if n.ChildNode != nil {
		return name + "." + n.ChildNode.String()
	}
	if len(n.TypeName)>0{
		return name +" : "+ n.TypeName
	}
	return name
}

// AddLeafNode (recursive) finds the nil leaf and replaces it with the node specified
func (n *Node) AddLeafNode(toAdd *Node) *Node {

	if n.ChildNode == nil {
		n.ChildNode = toAdd
	} else {
		n.ChildNode.AddLeafNode(toAdd)
	}

	return n
}

// Copy returns a deep copy of the Node
func (n Node) Copy() *Node {
	newNode := n

	if newNode.ChildNode != nil {
		n.ChildNode = newNode.ChildNode.Copy()
	}
	return &newNode
}
