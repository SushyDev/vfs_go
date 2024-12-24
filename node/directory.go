package node

type Directory struct {
	node *Node
}

func NewDirectory(node *Node) *Directory {
	return &Directory{
		node: node,
	}
}

func (directory *Directory) GetNode() *Node {
	return directory.node
}

func (directory *Directory) GetIdentifier() uint64 {
	return directory.node.GetIdentifier()
}

func (directory *Directory) GetName() string {
	return directory.node.GetName()
}
