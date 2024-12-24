package node

type File struct {
	node        *Node
	contentType string
	data        string
}

func NewFile(node *Node, contentType string, data string) *File {
	return &File{
		node:        node,
		contentType: contentType,
		data:        data,
	}
}

func (file *File) GetNode() *Node {
	return file.node
}

func (file *File) GetIdentifier() uint64 {
	return file.node.GetIdentifier()
}

func (file *File) GetName() string {
	return file.node.GetName()
}
