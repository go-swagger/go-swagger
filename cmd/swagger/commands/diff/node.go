package diff

type Node struct {
	Field     string
	TypeName  string
	IsArray   bool
	ChildNode *Node
}

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

func (n *Node) AddLeafNode(toAdd *Node) *Node {

	if n.ChildNode == nil {
		n.ChildNode = toAdd
	} else {
		n.ChildNode.AddLeafNode(toAdd)
	}

	return n
}

func (n Node) Copy() *Node {
	newNode := n

	if newNode.ChildNode != nil {
		n.ChildNode = newNode.ChildNode.Copy()
	}
	return &newNode
}
