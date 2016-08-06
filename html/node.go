// ==================================================
// Copyright (c) 2016 tacigar
// https://github.com/tacigar/Go-HTML
// ==================================================

package html

type NodeType int

const (
	TextNode NodeType = iota
	ElementNode
	DocumentNode // Root
)

type Node struct {
	Parent     *Node
	Children   []*Node
	Type       NodeType
	Data       string
	Attributes map[string]string
}
